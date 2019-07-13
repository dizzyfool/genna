package model

import (
	"testing"

	"github.com/dizzyfool/genna/util"
)

func TestRelation_GoName(t *testing.T) {
	type fields struct {
		SourceColumns []string
		TargetSchema  string
		TargetTable   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should generate simple name",
			fields: fields{
				SourceColumns: []string{"locationId"},
				TargetSchema:  util.PublicSchema,
				TargetTable:   "locations",
			},
			want: "Location",
		},
		{
			name: "Should generate multiple name",
			fields: fields{
				SourceColumns: []string{"city", "locationId"},
				TargetSchema:  util.PublicSchema,
				TargetTable:   "locations",
			},
			want: "CityLocation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRelation(tt.fields.SourceColumns, tt.fields.TargetSchema, tt.fields.TargetTable)
			if got := r.GoName; got != tt.want {
				t.Errorf("Relation.GoName = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelation_GoType(t *testing.T) {
	type fields struct {
		TargetSchema string
		TargetTable  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should generate from simple word",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "users",
			},
			want: "User",
		},
		{
			name: "Should generate from non-countable",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "audio",
			},
			want: "Audio",
		},
		{
			name: "Should generate from underscored",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "user_orders",
			},
			want: "UserOrder",
		},
		{
			name: "Should generate from camelCased",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "userOrders",
			},
			want: "UserOrder",
		},
		{
			name: "Should generate from plural in last place",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "usersWithOrders",
			},
			want: "UsersWithOrder",
		},
		{
			name: "Should generate from abracadabra",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "abracadabra",
			},
			want: "Abracadabra",
		},
		{
			name: "Should generate from numbers in first place",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "123-abc",
			},
			want: "T123Abc",
		},
		{
			name: "Should generate from name with dash & underscore",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "abc-123_abc",
			},
			want: "Abc123Abc",
		},
		{
			name: "Should generate with schema",
			fields: fields{
				TargetSchema: "information_schema",
				TargetTable:  "users",
			},
			want: "InformationSchemaUser",
		},
		{
			name: "Should generate without schema",
			fields: fields{
				TargetSchema: util.PublicSchema,
				TargetTable:  "users",
			},
			want: "User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRelation([]string{"ID"}, tt.fields.TargetSchema, tt.fields.TargetTable)
			if got := r.GoType; got != tt.want {
				t.Errorf("Relation.GoType = %v, want %v", got, tt.want)
			}
		})
	}
}
