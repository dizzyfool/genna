package model

import (
	"testing"
)

func TestRelation_StructFieldName(t *testing.T) {
	type fields struct {
		SourceColumns []string
		TargetTable   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate simple name",
			fields: fields{[]string{"locationId"}, "locations"},
			want:   "Location",
		},
		{
			name:   "Should generate multiple name",
			fields: fields{[]string{"cityId", "locationId"}, "locations"},
			want:   "CityLocation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				SourceColumns: tt.fields.SourceColumns,
				TargetTable:   tt.fields.TargetTable,
			}
			if got := r.StructFieldName(); got != tt.want {
				t.Errorf("Relation.StructFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelation_StructFieldTag(t *testing.T) {
	type fields struct {
		SourceColumn string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate simple name",
			fields: fields{"locationId"},
			want:   `pg:"fk:locationId"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				SourceColumns: []string{tt.fields.SourceColumn},
			}
			if got := r.StructFieldTag(); got != tt.want {
				t.Errorf("Relation.StructFieldTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelation_StructFieldType(t *testing.T) {
	type fields struct {
		Type         int
		SourceSchema string
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
				SourceSchema: PublicSchema,
				TargetSchema: PublicSchema,
				TargetTable:  "users",
			},
			want: "*User",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{TargetTable: "audio"},
			want:   "*Audio",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{TargetTable: "user_orders"},
			want:   "*UserOrder",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{TargetTable: "userOrders"},
			want:   "*UserOrder",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{TargetTable: "usersWithOrders"},
			want:   "*UsersWithOrder",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{TargetTable: "abracadabra"},
			want:   "*Abracadabra",
		},
		{
			name:   "Should generate from numbers in first place",
			fields: fields{TargetTable: "123-abc"},
			want:   "*T123Abc",
		},
		{
			name:   "Should generate from name with dash & underscore",
			fields: fields{TargetTable: "abc-123_abc"},
			want:   "*Abc123Abc",
		},
		{
			name: "Should generate with schema",
			fields: fields{
				SourceSchema: PublicSchema,
				TargetSchema: "information_schema",
				TargetTable:  "users",
			},
			want: "*InformationSchemaUser",
		},
		{
			name: "Should generate without schema",
			fields: fields{
				SourceSchema: "information_schema",
				TargetSchema: PublicSchema,
				TargetTable:  "users",
			},
			want: "*User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				Type:         tt.fields.Type,
				SourceSchema: tt.fields.SourceSchema,
				TargetSchema: tt.fields.TargetSchema,
				TargetTable:  tt.fields.TargetTable,
			}
			if got := r.StructFieldType(); got != tt.want {
				t.Errorf("Relation.StructFieldType() = %v, want %v", got, tt.want)
			}
		})
	}
}
