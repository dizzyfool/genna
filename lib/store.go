package bungen

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/LdDl/bungen/model"
	"github.com/LdDl/bungen/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/schema"
)

var formatter = schema.NewFormatter(pgdialect.New())

func format(pattern string, values ...interface{}) string {
	return formatter.FormatQuery(pattern, values...)
}

type table struct {
	Schema string `bun:"table_schema"`
	Name   string `bun:"table_name"`
}

func (t table) Entity() model.Entity {
	return model.NewEntity(t.Schema, t.Name, nil, nil)
}

type relation struct {
	Constraint    string   `bun:"constraint_name"`
	SourceSchema  string   `bun:"schema_name"`
	SourceTable   string   `bun:"table_name"`
	SourceColumns []string `bun:"columns,array"`
	TargetSchema  string   `bun:"target_schema"`
	TargetTable   string   `bun:"target_table"`
	TargetColumns []string `bun:"target_columns,array"`
}

func (r relation) Relation() model.Relation {
	return model.NewRelation(r.SourceColumns, r.TargetSchema, r.TargetTable, r.TargetColumns)
}

func (r relation) Target() table {
	return table{
		Schema: r.TargetSchema,
		Name:   r.TargetTable,
	}
}

type column struct {
	bun.BaseModel

	Schema     string   `bun:"schema_name"`
	Table      string   `bun:"table_name"`
	Name       string   `bun:"column_name"`
	IsNullable bool     `bun:"nullable"`
	IsArray    bool     `bun:"is_array"`
	Dimensions int      `bun:"dims"`
	Type       string   `bun:"type"`
	Default    string   `bun:"def"`
	IsPK       bool     `bun:"is_pk"`
	IsFK       bool     `bun:"is_fk"`
	MaxLen     int      `bun:"len"`
	Values     []string `bun:"enum,array"`
}

func (c column) Column(useSQLNulls bool, customTypes model.CustomTypeMapping) model.Column {
	return model.NewColumn(c.Name, c.Type, c.IsNullable, useSQLNulls, c.IsArray, c.Dimensions, c.IsPK, c.IsFK, c.MaxLen, c.Values, customTypes)
}

// Store is database helper
type store struct {
	db *bun.DB
}

// NewStore creates Store
func newStore(db *bun.DB) *store {
	return &store{db: db}
}

func (s *store) Schemas() ([]string, error) {
	query := `select nspname from pg_catalog.pg_namespace`
	var result []string

	err := s.db.NewRaw(query).Scan(context.Background(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting schemas info error: %w", err)
	}

	return result, nil
}

func (s *store) Tables(selected []string) ([]table, error) {
	var schemas []string
	var tables []interface{}

	for _, s := range selected {
		schema, table := util.Split(s)
		if table == "*" {
			schemas = append(schemas, schema)
		} else {
			tables = append(tables, []string{schema, table})
		}
	}

	var where []string
	if len(schemas) > 0 {
		where = append(where, format("(table_schema) in (?)", bun.In(schemas)))
	}
	if len(tables) > 0 {
		where = append(where, format("(table_schema, table_name) in (?)", bun.In(tables)))
	}

	query := `
        select 
            table_schema,
            table_name
        from information_schema.tables
        where 
            table_type = 'BASE TABLE' and 
            (
                ` + strings.Join(where, "or \n") + `
            )`

	var result []table

	err := s.db.NewRaw(query).Scan(context.Background(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting tables info error: %w", err)
	}
	return result, nil
}

// Relations gets relations of a selected table
func (s *store) Relations(tables []table) ([]relation, error) {
	ts := make([]interface{}, len(tables))
	for i, t := range tables {
		ts[i] = []string{t.Schema, t.Name}
	}

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
		       co.conname            as constraint_name,
		       ss.nspname            as schema_name,
		       s.relname             as table_name,
		       array_agg(sc.attname) as columns,
		       ts.nspname            as target_schema,
		       t.relname             as target_table,
		       array_agg(tc.attname) as target_columns
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
		  and (ss.nspname, s.relname) in (?)
		group by constraint_name, schema_name, table_name, target_schema, target_table
	`

	var relations []relation

	err := s.db.NewRaw(query, bun.In(ts)).Scan(context.Background(), &relations)
	if err != nil {
		return nil, fmt.Errorf("getting relations info error: %w", err)
	}

	return relations, nil
}

func (s store) Columns(tables []table) ([]column, error) {
	ts := make([]interface{}, len(tables))
	for i, t := range tables {
		ts[i] = []string{t.Schema, t.Name}
	}
	query := `
		with
		    enums as (
		        select distinct true                   as is_enum,
		                        sch.nspname            as table_schema,
		                        tb.relname             as table_name,
		                        col.attname            as column_name,
                                array_agg(e.enumlabel) as enum_values
		        from pg_class tb
		        left join pg_namespace sch on sch.oid = tb.relnamespace
		        left join pg_attribute col on col.attrelid = tb.oid
		        inner join pg_enum e on e.enumtypid = col.atttypid
				group by 1, 2, 3, 4
		    ),
		    arrays as (
		        select sch.nspname  as table_schema,
		               tb.relname   as table_name,
		               col.attname  as column_name,
		               col.attndims as array_dims
		        from pg_class tb
		        left join pg_namespace sch on sch.oid = tb.relnamespace
		        left join pg_attribute col on col.attrelid = tb.oid
		        where col.attndims > 0
		    ),
		    info as (
				select distinct
				 	kcu.table_schema as table_schema,
					kcu.table_name   as table_name,
					kcu.column_name  as column_name,
					array_agg((
						select constraint_type::text 
						from information_schema.table_constraints tc 
						where tc.constraint_name = kcu.constraint_name 
							and tc.constraint_schema = kcu.constraint_schema 
							and tc.constraint_catalog = kcu.constraint_catalog
						limit 1
					)) as constraint_types
				from information_schema.key_column_usage kcu
				group by kcu.table_schema, kcu.table_name, kcu.column_name
		    )
		select distinct c.table_schema = 'public' as is_public,
                        c.table_schema            as schema_name,
		                c.table_name              as table_name,
		                c.column_name             as column_name,
                        c.ordinal_position        as ordinal,
		                case
		                when i.constraint_types is null
		                then false
		                else 'PRIMARY KEY'=any (i.constraint_types)
		                end                                    as is_pk,
		                'FOREIGN KEY'=any (i.constraint_types) as is_fk,
		                c.is_nullable = 'YES'                  as nullable,
		                c.data_type = 'ARRAY'                  as is_array,
		                coalesce(a.array_dims, 0)              as dims,
		                case
		                when e.is_enum = true
		                then 'varchar'
		                else ltrim(c.udt_name, '_')
		                end                         as type,
		                c.column_default            as def,
                        c.character_maximum_length  as len,
						e.enum_values 				as enum
		from information_schema.tables t
		left join information_schema.columns c using (table_name, table_schema)
		left join info i using (table_name, table_schema, column_name)
		left join arrays a using (table_name, table_schema, column_name)
		left join enums e using (table_name, table_schema, column_name)
		where (t.table_schema, t.table_name) in (?)
		  and t.table_type = 'BASE TABLE'
		order by 1 desc, 2, 3, 5 asc, 6 desc nulls last
	`

	var columns []column
	err := s.db.NewRaw(query, bun.In(ts)).Scan(context.Background(), &columns)
	if err != nil {
		return nil, fmt.Errorf("getting columns info error: %w", err)
	}

	return columns, nil
}

// Sort sorts table by schema and name (public tables always first)
func Sort(tables []table) []table {
	sort.Slice(tables, func(i, j int) bool {
		ti := tables[i]
		tj := tables[j]

		if ti.Schema == tj.Schema {
			return ti.Name < tj.Name
		}

		if ti.Schema == util.PublicSchema {
			return true
		}
		if tj.Schema == util.PublicSchema {
			return false
		}

		return util.Join(ti.Schema, ti.Name) < util.Join(tj.Schema, tj.Name)
	})

	return tables
}
