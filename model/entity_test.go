package model

import (
	"testing"

	"github.com/dizzyfool/genna/util"
)

func TestTable_GoName(t *testing.T) {
	type fields struct {
		Name   string
		Schema string
	}
	tests := []struct {
		name       string
		fields     fields
		withSchema bool
		want       string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{Name: "users"},
			want:   "User",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{Name: "audio"},
			want:   "Audio",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{Name: "user_orders"},
			want:   "UserOrder",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{Name: "userOrders"},
			want:   "UserOrder",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{Name: "usersWithOrders"},
			want:   "UsersWithOrder",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{Name: "abracadabra"},
			want:   "Abracadabra",
		},
		{
			name:       "Should generate from simple word with public schema",
			fields:     fields{Name: "users", Schema: "public"},
			withSchema: true,
			want:       "User",
		},
		{
			name:       "Should generate from simple word with custom schema",
			fields:     fields{Name: "users", Schema: "users"},
			withSchema: true,
			want:       "UsersUser",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := NewEntity(tt.fields.Schema, tt.fields.Name, nil, nil)
			if got := tbl.GoName; got != tt.want {
				t.Errorf("Entity.GoName = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntity_AddColumn(t *testing.T) {
	entity := NewEntity(util.PublicSchema, "test", nil, nil)

	t.Run("Should add column", func(t *testing.T) {
		column1 := NewColumn("name", TypePGText, false, false, false, 0, false, false, 0, []string{}, 9)
		column2 := NewColumn("name_", TypePGText, false, false, false, 0, false, false, 0, []string{}, 9)
		column3 := NewColumn("timeout", TypePGInterval, false, false, false, 0, false, false, 0, []string{}, 9)
		column4 := NewColumn("duration", TypePGInterval, false, false, false, 0, false, false, 0, []string{}, 9)

		t.Run("Should add first column", func(t *testing.T) {
			entity.AddColumn(column1)
			if len(entity.Columns) != 1 {
				t.Errorf("Entity.Columns = %v, want %v", len(entity.Columns), 1)
			}
		})

		t.Run("Should add second column with same name", func(t *testing.T) {
			entity.AddColumn(column2)
			if len(entity.Columns) != 2 {
				t.Errorf("Entity.Columns = %v, want %v", len(entity.Columns), 2)
			}
			if entity.Columns[1].GoName != "Name1" {
				t.Errorf("Entity.Columns[1].GoName = %v, want %v", entity.Columns[1].GoName, "Name1")
			}
		})

		t.Run("Should add column with import", func(t *testing.T) {
			entity.AddColumn(column3)
			if len(entity.Imports) != 1 {
				t.Errorf("Entity.Imports = %v, want %v", len(entity.Imports), 1)
			}
		})

		t.Run("Should add column without import", func(t *testing.T) {
			entity.AddColumn(column4)
			if len(entity.Imports) != 1 {
				t.Errorf("Entity.Imports = %v, want %v", len(entity.Imports), 1)
			}
		})
	})
}

func TestEntity_AddRelation(t *testing.T) {
	column1 := NewColumn("test", TypePGText, false, false, false, 0, false, false, 0, []string{}, 9)
	relation1 := NewRelation([]string{"userId"}, util.PublicSchema, "users")

	entity := NewEntity(util.PublicSchema, "test", []Column{column1}, []Relation{relation1})

	t.Run("Should add column", func(t *testing.T) {
		relation2 := NewRelation([]string{"locationId"}, util.PublicSchema, "locations")
		relation3 := NewRelation([]string{"testId"}, util.PublicSchema, "tests")
		relation4 := NewRelation([]string{"testId"}, util.PublicSchema, "tests_")

		t.Run("Should add second relation", func(t *testing.T) {
			entity.AddRelation(relation2)
			if len(entity.Relations) != 2 {
				t.Errorf("Entity.Relations = %v, want %v", len(entity.Relations), 2)
			}
		})

		t.Run("Should add third relation with same name", func(t *testing.T) {
			entity.AddRelation(relation3)
			if len(entity.Relations) != 3 {
				t.Errorf("Entity.Relations = %v, want %v", len(entity.Relations), 3)
			}
			if entity.Relations[2].GoName != "TestRel" {
				t.Errorf("Entity.Relations[2].GoName = %v, want %v", entity.Relations[2].GoName, "Test1")
			}
		})

		t.Run("Should add forth relation with same name", func(t *testing.T) {
			entity.AddRelation(relation4)
			if len(entity.Relations) != 4 {
				t.Errorf("Entity.Relations = %v, want %v", len(entity.Relations), 4)
			}
			if entity.Relations[3].GoName != "TestRel1" {
				t.Errorf("Entity.Relations[3].GoName = %v, want %v", entity.Relations[3].GoName, "TestRel1")
			}
		})
	})
}

func TestEntity_HasMultiplePKs(t *testing.T) {
	entity := NewEntity(util.PublicSchema, "test", nil, nil)

	t.Run("Should add column", func(t *testing.T) {
		column1 := NewColumn("userId", TypePGText, false, false, false, 0, true, false, 0, []string{}, 9)
		column2 := NewColumn("locationId", TypePGText, false, false, false, 0, true, false, 0, []string{}, 9)

		t.Run("Should check for one key", func(t *testing.T) {
			entity.AddColumn(column1)
			if v := entity.HasMultiplePKs(); v {
				t.Errorf("Entity.HasMultiplePKs() = %v, want %v", v, false)
			}
		})

		t.Run("Should check for several keys", func(t *testing.T) {
			entity.AddColumn(column2)
			if v := entity.HasMultiplePKs(); !v {
				t.Errorf("Entity.HasMultiplePKs() = %v, want %v", v, true)
			}
		})
	})
}
