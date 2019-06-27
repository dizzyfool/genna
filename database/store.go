package database

import (
	"github.com/dizzyfool/genna/model"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
)

// Store is database helper
type Store struct {
	db orm.DB
}

// NewStore creates store
func NewStore(db orm.DB) *Store {
	return &Store{db: db}
}

// columnRow stores raw column information
type columnRow struct {
	SchemaName string   `sql:"schema"`
	TblName    string   `sql:"table"`
	ColumnName string   `sql:"column"`
	IsNullable bool     `sql:"nullable"`
	IsArray    bool     `sql:"array"`
	Dimensions int      `sql:"dims"`
	Type       string   `sql:"type"`
	Default    string   `sql:"default"`
	IsPK       bool     `sql:"pk"`
	IsFK       bool     `sql:"fk"`
	MaxLen     int      `sql:"len"`
	Enum       []string `sql:"enum,array"`
}

// Table converts row to Table model
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
		MaxLen:     r.MaxLen,
		Enum:       r.Enum,
	}
}

// columnRow stores raw relation information
type relationRow struct {
	Constraint        string   `sql:"constraint"`
	SchemaName        string   `sql:"schema"`
	TblName           string   `sql:"table"`
	ColumnsName       []string `sql:"columns,array"`
	TargetSchemaName  string   `sql:"targetSchema"`
	TargetTblName     string   `sql:"targetTable"`
	TargetColumnsName []string `sql:"targetColumns,array"`
}

// Relation converts row to Relation model
func (r *relationRow) Relation() model.Relation {
	return model.Relation{
		// TODO HasMany relation
		Type:          model.HasOne,
		SourceSchema:  r.SchemaName,
		SourceTable:   r.TblName,
		SourceColumns: r.ColumnsName,
		TargetSchema:  r.TargetSchemaName,
		TargetTable:   r.TargetTblName,
		TargetColumns: r.TargetColumnsName,
	}
}

func (s *Store) queryTables(schema []string) ([]columnRow, error) {
	rows := make([]columnRow, 0)
	query := `
		with
		    enums as (
		        select distinct true                   as "is_enum",
		                        sch.nspname            as "table_schema",
		                        tb.relname             as "table_name",
		                        col.attname            as "column_name",
                                array_agg(e.enumlabel) as "enum_values"
		        from pg_class tb
		        left join pg_namespace sch on sch.oid = tb.relnamespace
		        left join pg_attribute col on col.attrelid = tb.oid
		        inner join pg_enum e on e.enumtypid = col.atttypid
				group by 1, 2, 3, 4
		    ),
		    arrays as (
		        select sch.nspname  as "table_schema",
		               tb.relname   as "table_name",
		               col.attname  as "column_name",
		               col.attndims as "array_dims"
		        from pg_class tb
		        left join pg_namespace sch on sch.oid = tb.relnamespace
		        left join pg_attribute col on col.attrelid = tb.oid
		        where col.attndims > 0
		    ),
		    info as (
				select distinct
				 	kcu.table_schema as "table_schema",
					kcu.table_name   as "table_name",
					kcu.column_name  as "column_name",
					array_agg((
						select constraint_type::text 
						from information_schema.table_constraints tc 
						where tc.constraint_name = kcu.constraint_name 
							and tc.constraint_schema = kcu.constraint_schema 
							and tc.constraint_catalog = tc.constraint_catalog
					)) as "constraint_types"
				from information_schema.key_column_usage kcu
				group by kcu.table_schema, kcu.table_name, kcu.column_name
		    )
		select distinct c."table_schema"                       as "schema",
		                c."table_name"                         as "table",
		                c."column_name"                        as "column",
		                case
		                when i.constraint_types is null
		                then false
		                else 'PRIMARY KEY'=any (i.constraint_types)
		                end                                    as "pk",
		                'FOREIGN KEY'=any (i.constraint_types) as "fk",
		                c."is_nullable" = 'YES'                as "nullable",
		                c."data_type" = 'ARRAY'                as "array",
		                coalesce(a.array_dims, 0)              as "dims",
		                case
		                when e.is_enum = true
		                then 'varchar'
		                else ltrim(c."udt_name", '_')
		                end                                    as "type",
		                c.column_default                       as "default",
                        c.character_maximum_length             as "len",
						e.enum_values 						   as "enum"
		from information_schema.tables t
		left join information_schema.columns c using (table_name, table_schema)
		left join info i using (table_name, table_schema, column_name)
		left join arrays a using (table_name, table_schema, column_name)
		left join enums e using (table_name, table_schema, column_name)
		where t."table_schema" in (?)
		  and t.table_type = 'BASE TABLE'
		order by 1, 2, 4 desc nulls last
	`

	_, err := s.db.Query(&rows, query, pg.In(schema))

	return rows, err
}

// Tables get basic tables information among selected schemas
func (s *Store) Tables(schema []string) ([]model.Table, error) {

	rows, err := s.queryTables(schema)
	if err != nil {
		return nil, errors.Wrap(err, "getting table info error")
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
		select distinct 
		       co.conname            as "constraint",
		       ss.nspname            as "schema",
		       s.relname             as "table",
		       array_agg(sc.attname) as "columns",
		       ts.nspname            as "targetSchema",
		       t.relname             as "targetTable",
		       array_agg(tc.attname) as "targetColumns"
		from pg_constraint co
		left join tables s on co.conrelid = s.oid
		left join schemas ss on s.relnamespace = ss.oid
		left join columns sc on s.oid = sc.attrelid and sc.attnum = any (co.conkey)
		left join tables t on co.confrelid = t.oid
		left join schemas ts on t.relnamespace = ts.oid
		left join columns tc on t.oid = tc.attrelid and tc.attnum = any (co.confkey)
		where co.contype = 'f'
		  and co.conrelid in (select oid from pg_class c where c.relkind = 'r')
		  and array_position(co.conkey, sc.attnum) = array_position(co.confkey, tc.attnum)
		  and ss.nspname = ?
		  and s.relname = ?
		group by "constraint", schema, "table", "targetSchema", "targetTable"
	`

	// TODO HasMany relation "or (ts.nspname = ? and t.relname = ?)"

	rows := make([]relationRow, 0)
	if _, err := s.db.Query(&rows, query, schema, table); err != nil {
		return nil, errors.Wrap(err, "getting relations info error")
	}

	relations := make([]model.Relation, len(rows))
	for i, row := range rows {
		relations[i] = row.Relation()
	}

	return relations, nil
}
