package database

import (
	"fmt"
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

// Row stores raw column information
type Row struct {
	SchemaName string `sql:"schema"`
	TblName    string `sql:"table"`
	ColumnName string `sql:"column"`
	IsNullable bool   `sql:"nullable"`
	IsArray    bool   `sql:"array"`
	Type       string `sql:"type"`
	Default    string `sql:"default"`
	IsPK       bool   `sql:"pk"`
	IsFK       bool   `sql:"fk"`
}

// Column converts row to Table model
func (r *Row) Table() model.Table {
	return model.Table{
		Schema: r.SchemaName,
		Name:   r.TblName,
	}
}

// Column converts row to Column model
func (r *Row) Column() model.Column {
	return model.Column{
		Name: r.ColumnName,
		Type: r.Type,
		// TODO dimensions!
		IsArray:    r.IsArray,
		IsNullable: r.IsNullable,
		IsPK:       r.IsPK,
		IsFK:       r.IsFK,
	}
}

// Tables get basic tables information among selected schemas
func (s *Store) Tables(schema []string) ([]model.Table, error) {
	query := `
		select
		       c."table_schema"                  as "schema",
		       c."table_name"                    as "table",
		       c."column_name"                   as "column",
		       c."is_nullable" = 'YES'           as "nullable",
		       c."data_type" = 'ARRAY'           as "array",
		       ltrim(c."udt_name", '_')          as "type",
		       c.column_default                  as "default",
		       f.constraint_type = 'PRIMARY KEY' as "pk",
		       f.constraint_type = 'FOREIGN KEY' as "fk"
		from information_schema.columns c
		left join information_schema.key_column_usage k using (table_name, table_schema, column_name)
		left join information_schema.table_constraints f using (table_name, table_schema, constraint_name)
		where c."table_schema" in (?)
		order by 1, 2;
	`

	rows := make([]Row, 0)
	if _, err := s.db.Query(&rows, query, pg.In(schema)); err != nil {
		return nil, err
	}

	var current *model.Table
	tables := make([]model.Table, 0)
	for _, row := range rows {
		if current == nil || current.Schema != row.SchemaName || current.Name != row.TblName {
			tables = append(tables, row.Table())
			current = &(tables[len(tables)-1])
		}
		current.Columns = append(current.Columns, row.Column())
	}

	return tables, nil
}

// Relations gets relations of a selected table
// TODO implement it!
func (s *Store) Relations(schema, table string) ([]model.Relation, error) {
	a := `
		with
		    foreign_keys as (
		        select o.conname                                                     as constraint_name,
		               (select nspname from pg_namespace where oid = m.relnamespace) as table_schema,
		               m.relname                                                     as table_name,
		               (select a.attname
		                from pg_attribute a
		                where a.attrelid = m.oid
		                  and a.attnum = o.conkey[1]
		                  and a.attisdropped = false)                                as column_name,
		               (select nspname from pg_namespace where oid = f.relnamespace) as target_schema,
		               f.relname                                                     as target_table,
		               (select a.attname
		                from pg_attribute a
		                where a.attrelid = f.oid
		                  and a.attnum = o.confkey[1]
		                  and a.attisdropped = false)                                as target_column
		        from pg_constraint o
		        left join pg_class f on f.oid = o.confrelid
		        left join pg_class m on m.oid = o.conrelid
		        where o.contype = 'f'
		          and o.conrelid in (select oid from pg_class c where c.relkind = 'r')
		    ),
		    primary_keys as (
		        select o.conname                                                     as constraint_name,
		               (select nspname from pg_namespace where oid = m.relnamespace) as table_schema,
		               m.relname                                                     as table_name,
		               (select a.attname
		                from pg_attribute a
		                where a.attrelid = m.oid
		                  and a.attnum = o.conkey[1]
		                  and a.attisdropped = false)                                as column_name
		        from pg_constraint o
		        left join pg_class m on m.oid = o.conrelid
		        where o.contype = 'p'
		    )
		select c."table_name",
		       c."column_name",
		       c."is_nullable",
		       c."data_type",
		       c."udt_name",
		       c."character_maximum_length",
		       c."numeric_precision",
		       k.constraint_name is not null as "is_pk",
		       f.constraint_name is not null as "is_fk",
		       f.target_schema as "foreign_key_target_schema",
		       f.target_table as "foreign_key_target_table",
		       f.target_column as "foreign_key_target_columt"
		from information_schema.columns c
		left join primary_keys k using (table_name, table_schema, column_name)
		left join foreign_keys f using (table_name, table_schema, column_name)
		where c."table_schema" = 'public'
		  and c."table_name" = 'top3cpc';
`

	fmt.Print(a)

	return nil, nil
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
