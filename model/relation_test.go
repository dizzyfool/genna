package model

import "testing"

func TestRelation_StructFieldName(t *testing.T) {
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
			fields: fields{"location"},
			want:   "Location",
		},
		{
			name:   "Should generate simple name with Id",
			fields: fields{"locationId"},
			want:   "Location",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{"location_id"},
			want:   "Location",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				SourceColumn: tt.fields.SourceColumn,
			}
			if got := r.StructFieldName(); got != tt.want {
				t.Errorf("Relation.StructFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelation_StructFieldType(t *testing.T) {
	type fields struct {
		TargetSchema string
		TargetTable  string
	}
	tests := []struct {
		name       string
		fields     fields
		withSchema bool
		want       string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{TargetTable: "users"},
			want:   "*User",
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
			name:   "Should generate from plural in first place",
			fields: fields{TargetTable: "usersWithOrders"},
			want:   "*UserWithOrders",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{TargetTable: "usersWithOrders"},
			want:   "*UserWithOrders",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{TargetTable: "abracadabra"},
			want:   "*Abracadabra",
		},
		{
			name:       "Should generate with schema",
			fields:     fields{TargetSchema: "information_schema", TargetTable: "users"},
			withSchema: true,
			want:       "*InformationSchemaUser",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				TargetSchema: tt.fields.TargetSchema,
				TargetTable:  tt.fields.TargetTable,
			}
			if got := r.StructFieldType(tt.withSchema); got != tt.want {
				t.Errorf("Relation.StructFieldType() = %v, want %v", got, tt.want)
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
			want:   `sql:"fk:locationId"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				SourceColumn: tt.fields.SourceColumn,
			}
			if got := r.StructFieldTag(); got != tt.want {
				t.Errorf("Relation.StructFieldTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
