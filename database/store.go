package database

import (
	"strings"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/dizzyfool/genna/model"
)

type Store struct {
	db orm.DB
}

func NewStore(db orm.DB) *Store {
	return &Store{db: db}
}

// columnRow stores raw column information
type columnRow struct {
	SchemaName string `sql:"schema"`
	TblName    string `sql:"table"`
	ColumnName string `sql:"column"`
	IsNullable bool   `sql:"nullable"`
	IsArray    bool   `sql:"array"`
	Dimensions int    `sql:"dims"`
	Type       string `sql:"type"`
	Default    string `sql:"default"`
	IsPK       bool   `sql:"pk"`
	IsFK       bool   `sql:"fk"`
}

// Column converts row to Table model
func (r *columnRow) Table() model.Table {
	return model.Table{
		Schema: r.SchemaName,
		Name:   r.TblName,
	}
}

// Column converts row to Column model
func (r *columnRow) Column() model.Column {
	return model.Column{
		Name:       r.ColumnName,
		Type:       r.Type,
		IsArray:    r.IsArray,
		Dimensions: r.Dimensions,
		IsNullable: r.IsNullable,
		IsPK:       r.IsPK,
		IsFK:       r.IsFK,
	}
}

// columnRow stores raw relation information
type relationRow struct {
	Constraint       string `sql:"constraint"`
	SchemaName       string `sql:"schema"`
	TblName          string `sql:"table"`
	ColumnName       string `sql:"column"`
	TargetSchemaName string `sql:"targetSchema"`
	TargetTblName    string `sql:"targetTable"`
	TargetColumnName string `sql:"targetColumn"`
}

// Column converts row to Column model
func (r *relationRow) Relation() model.Relation {
	return model.Relation{
		Type:         model.HasOne,
		SourceColumn: r.ColumnName,
		TargetSchema: r.TargetSchemaName,
		TargetTable:  r.TargetTblName,
		TargetColumn: r.TargetColumnName,
	}
}

// Tables get basic tables information among selected schemas
func (s *Store) Tables(schema []string) ([]model.Table, error) {
	query := `
		with
		    arrays as (
		        select sch.nspname  as "table_schema",
		               tb.relname   as "table_name",
		               col.attname  as "column_name",
		               col.attndims as "array_dims"
		        from pg_class tb
		        left join pg_namespace sch on sch.oid = tb.relnamespace
		        left join pg_attribute col on col.attrelid = tb.oid
		        where tb.relname = 'mapping'
		          and sch.nspname = 'public'
		          and col.attndims > 0
		    )
		select
		       c."table_schema"                  as "schema",
		       c."table_name"                    as "table",
		       c."column_name"                   as "column",
		       c."is_nullable" = 'YES'           as "nullable",
		       c."data_type" = 'ARRAY'           as "array",
		       coalesce(a.array_dims, 0)         as "dims",
		       ltrim(c."udt_name", '_')          as "type",
		       c.column_default                  as "default",
		       f.constraint_type = 'PRIMARY KEY' as "pk",
		       f.constraint_type = 'FOREIGN KEY' as "fk"
		from information_schema.columns c
		left join information_schema.key_column_usage k using (table_name, table_schema, column_name)
		left join information_schema.table_constraints f using (table_name, table_schema, constraint_name)
		left join arrays a using (table_name, table_schema, column_name)
		where c."table_schema" in (?)
		order by 1, 2;
	`

	rows := make([]columnRow, 0)
	if _, err := s.db.Query(&rows, query, pg.In(schema)); err != nil {
		return nil, err
	}

	var current = -1
	tables := make([]model.Table, 0)
	for _, row := range rows {
		if current == -1 ||
			tables[current].Schema != row.SchemaName ||
			tables[current].Name != row.TblName {

			table := row.Table()

			// filling relations for table
			relations, err := s.Relations(table.Schema, table.Name)
			if err != nil {
				return nil, err
			}
			table.Relations = relations

			tables = append(tables, table)
			current = len(tables) - 1
		}

		// filling columns for table
		tables[current].Columns = append(tables[current].Columns, row.Column())
	}

	return tables, nil
}

// Relations gets relations of a selected table
func (s *Store) Relations(schema, table string) ([]model.Relation, error) {
	query := `
		with
		    schemas as (
		        select nspname, oid
		        from pg_namespace
		    ),
		    tables as (
		        select oid, relnamespace, relname, relkind
		        from pg_class
		    ),
		    columns as (
		        select attrelid, attname, attnum
		        from pg_attribute a
		        where a.attisdropped = false
		    )
		select co.conname as "constraint",
		       ss.nspname as "schema",
		       s.relname  as "table",
		       sc.attname as "column",
		       ts.nspname as "targetSchema",
		       t.relname  as "targetTable",
		       tc.attname as "targetColumn"
		from pg_constraint co
		left join tables s on co.conrelid = s.oid
		left join schemas ss on s.relnamespace = ss.oid
		left join columns sc on s.oid = sc.attrelid and co.conkey[1] = sc.attnum
		left join tables t on co.confrelid = t.oid
		left join schemas ts on t.relnamespace = ts.oid
		left join columns tc on t.oid = tc.attrelid and co.confkey[1] = tc.attnum
		where co.contype = 'f'
		  and co.conrelid in (select oid from pg_class c where c.relkind = 'r')
		  and ss.nspname = ?
		  and s.relname = ?
	`

	// TODO HasMany relation "or (ts.nspname = ? and t.relname = ?)"

	rows := make([]relationRow, 0)
	if _, err := s.db.Query(&rows, query, schema, table); err != nil {
		return nil, err
	}

	relations := make([]model.Relation, len(rows))
	for i, row := range rows {
		relations[i] = row.Relation()
	}

	return relations, nil
}

func Schemas(tables []string) (schemas []string) {
	index := map[string]struct{}{}
	for _, table := range tables {
		schema := "public"

		d := strings.Split(table, ".")
		if len(d) >= 2 {
			schema = d[0]
		}

		if _, ok := index[schema]; ok {
			continue
		}

		index[schema] = struct{}{}
		schemas = append(schemas, schema)
	}

	return
}
