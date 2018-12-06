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
		TargetTable string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{"users"},
			want:   "*User",
		},
		{
			name:   "Should generate from non-countable",
			fields: fields{"audio"},
			want:   "*Audio",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{"user_orders"},
			want:   "*UserOrder",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{"userOrders"},
			want:   "*UserOrder",
		},
		{
			name:   "Should generate from plural in first place",
			fields: fields{"usersWithOrders"},
			want:   "*UserWithOrders",
		},
		{
			name:   "Should generate from plural in last place",
			fields: fields{"usersWithOrders"},
			want:   "*UserWithOrders",
		},
		{
			name:   "Should generate from abracadabra",
			fields: fields{"abracadabra"},
			want:   "*Abracadabra",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				TargetTable: tt.fields.TargetTable,
			}
			if got := r.StructFieldType(); got != tt.want {
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
