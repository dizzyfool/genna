package genna

import (
	"reflect"
	"testing"

	"github.com/dizzyfool/genna/model"

	"github.com/go-pg/pg/v9"
)

func prepareStore() (*store, error) {
	db, err := newDatabase(prepareReq())
	if err != nil {
		return nil, err
	}

	return newStore(db), nil
}

func Test_format(t *testing.T) {

	tests := []struct {
		name    string
		pattern string
		values  []interface{}
		want    string
	}{
		{
			name:    "Should format pg.Multi",
			pattern: "(id, name) in (?)",
			values: []interface{}{
				[]string{"1", "test"},
				[]string{"2", "test"},
			},
			want: "(id, name) in (('1','test'),('2','test'))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format(tt.pattern, pg.InMulti(tt.values...)); got != tt.want {
				t.Errorf("format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_table_Entity(t *testing.T) {
	type fields struct {
		Schema string
		Name   string
	}
	tests := []struct {
		name   string
		fields fields
		want   model.Entity
	}{
		{
			name: "Should create entity",
			fields: fields{
				Schema: "public",
				Name:   "users",
			},
			want: model.NewEntity("public", "users", nil, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := table{
				Schema: tt.fields.Schema,
				Name:   tt.fields.Name,
			}
			if got := z.Entity(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("table.Entity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_relation_Relation(t *testing.T) {
	type fields struct {
		Constraint    string
		SourceSchema  string
		SourceTable   string
		SourceColumns []string
		TargetSchema  string
		TargetTable   string
		TargetColumns []string
	}
	tests := []struct {
		name   string
		fields fields
		want   model.Relation
	}{
		{
			name: "Should create relation",
			fields: fields{
				Constraint:    "test",
				SourceSchema:  "public",
				SourceTable:   "users",
				SourceColumns: []string{"locationId"},
				TargetSchema:  "geo",
				TargetTable:   "locations",
				TargetColumns: []string{"locationId"},
			},
			want: model.NewRelation([]string{"locationId"}, "geo", "locations"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := relation{
				Constraint:    tt.fields.Constraint,
				SourceSchema:  tt.fields.SourceSchema,
				SourceTable:   tt.fields.SourceTable,
				SourceColumns: tt.fields.SourceColumns,
				TargetSchema:  tt.fields.TargetSchema,
				TargetTable:   tt.fields.TargetTable,
				TargetColumns: tt.fields.TargetColumns,
			}
			if got := r.Relation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("relation.Relation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_relation_Target(t *testing.T) {
	type fields struct {
		Constraint    string
		SourceSchema  string
		SourceTable   string
		SourceColumns []string
		TargetSchema  string
		TargetTable   string
		TargetColumns []string
	}
	tests := []struct {
		name   string
		fields fields
		want   table
	}{
		{
			name: "Should create target table",
			fields: fields{
				Constraint:    "test",
				SourceSchema:  "public",
				SourceTable:   "users",
				SourceColumns: []string{"locationId"},
				TargetSchema:  "geo",
				TargetTable:   "locations",
				TargetColumns: []string{"locationId"},
			},
			want: table{
				Schema: "geo",
				Name:   "locations",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := relation{
				Constraint:    tt.fields.Constraint,
				SourceSchema:  tt.fields.SourceSchema,
				SourceTable:   tt.fields.SourceTable,
				SourceColumns: tt.fields.SourceColumns,
				TargetSchema:  tt.fields.TargetSchema,
				TargetTable:   tt.fields.TargetTable,
				TargetColumns: tt.fields.TargetColumns,
			}
			if got := r.Target(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("relation.Target() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_column_Column(t *testing.T) {
	type fields struct {
		Schema     string
		Table      string
		Name       string
		IsNullable bool
		IsArray    bool
		Dimensions int
		Type       string
		Default    string
		IsPK       bool
		IsFK       bool
		MaxLen     int
		Values     []string
	}
	tests := []struct {
		name   string
		fields fields
		want   model.Column
	}{
		{
			name: "Should create column",
			fields: fields{
				Schema:     "public",
				Table:      "users",
				Name:       "userId",
				IsNullable: false,
				IsArray:    false,
				Dimensions: 0,
				Type:       model.TypePGInt8,
				Default:    "",
				IsPK:       true,
				IsFK:       false,
				MaxLen:     0,
				Values:     []string{},
			},
			want: model.NewColumn("userId", model.TypePGInt8, false, false, false, 0, true, false, 0, []string{}, 9),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := column{
				Schema:     tt.fields.Schema,
				Table:      tt.fields.Table,
				Name:       tt.fields.Name,
				IsNullable: tt.fields.IsNullable,
				IsArray:    tt.fields.IsArray,
				Dimensions: tt.fields.Dimensions,
				Type:       tt.fields.Type,
				Default:    tt.fields.Default,
				IsPK:       tt.fields.IsPK,
				IsFK:       tt.fields.IsFK,
				MaxLen:     tt.fields.MaxLen,
				Values:     tt.fields.Values,
			}
			if got := c.Column(false, 9); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("column.Column() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_store_Tables(t *testing.T) {
	store, err := prepareStore()
	if err != nil {
		t.Errorf("prepare Store error = %v", err)
		return
	}

	t.Run("Should get all tables from test DB", func(t *testing.T) {
		tables, err := store.Tables([]string{"public.*", "geo.*"})
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		if ln := len(tables); ln != 3 {
			t.Errorf("len(Store.Tables()) = %v, want %v", ln, 3)
			return
		}
	})

	t.Run("Should get specific table from test DB", func(t *testing.T) {
		tables, err := store.Tables([]string{"public.users"})
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		if ln := len(tables); ln != 1 {
			t.Errorf("len(Store.Tables()) = %v, want %v", ln, 1)
			return
		}
	})

	t.Run("Should get specific & geo tables from test DB", func(t *testing.T) {
		tables, err := store.Tables([]string{"public.users", "geo.*"})
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		if ln := len(tables); ln != 2 {
			t.Errorf("len(Store.Tables()) = %v, want %v", ln, 2)
			return
		}
	})
}

func Test_store_Relations(t *testing.T) {
	store, err := prepareStore()
	if err != nil {
		t.Errorf("prepare Store error = %v", err)
		return
	}

	t.Run("Should get all relations from test DB", func(t *testing.T) {
		tables, err := store.Tables([]string{"public.*"})
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		relations, err := store.Relations(tables)
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		if ln := len(relations); ln != 1 {
			t.Errorf("len(Store.Relations()) = %v, want %v", ln, 1)
			return
		}
	})
}

func Test_store_Columns(t *testing.T) {
	store, err := prepareStore()
	if err != nil {
		t.Errorf("prepare Store error = %v", err)
		return
	}

	t.Run("Should get all columns from test DB", func(t *testing.T) {
		tables, err := store.Tables([]string{"public.*"})
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		columns, err := store.Columns(tables)
		if err != nil {
			t.Errorf("get tables error = %v", err)
			return
		}

		if ln := len(columns); ln != 10 {
			t.Errorf("len(Store.Columns()) = %v, want %v", ln, 10)
			return
		}
	})
}
