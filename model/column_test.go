package model

import (
	"testing"
)

func TestColumn_GoName(t *testing.T) {
	tests := []struct {
		name   string
		pgName string
		want   string
	}{
		{
			name:   "Should generate from simple word",
			pgName: "title",
			want:   "Title",
		},
		{
			name:   "Should generate from underscored",
			pgName: "short_title",
			want:   "ShortTitle",
		},
		{
			name:   "Should generate from camelCased",
			pgName: "shortTitle",
			want:   "ShortTitle",
		},
		{
			name:   "Should generate with underscored_id",
			pgName: "location_id",
			want:   "LocationID",
		},
		{
			name:   "Should generate with camelCasedId",
			pgName: "locationId",
			want:   "LocationID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewColumn(tt.pgName, TypePGText, false, false, false, 0, false, false, 0, []string{}, CustomTypeMapping{})
			if c.GoName != tt.want {
				t.Errorf("Column.Name = %v, want %v", c.GoName, tt.want)
			}
		})
	}
}

func TestColumn_GoType(t *testing.T) {
	type fields struct {
		pgType   string
		array    bool
		dims     int
		nullable bool
		sqlNulls bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should generate int2 type",
			fields: fields{
				pgType:   TypePGInt2,
				array:    false,
				dims:     0,
				nullable: false,
			},
			want: "int",
		},
		{
			name: "Should generate int2 array type",
			fields: fields{
				pgType:   TypePGInt2,
				array:    true,
				dims:     2,
				nullable: false,
			},
			want: "[][]int",
		},
		{
			name: "Should generate int2 nullable type",
			fields: fields{
				pgType:   TypePGInt2,
				array:    true,
				dims:     2,
				nullable: true,
			},
			want: "[][]int",
		},
		{
			name: "Should generate struct type",
			fields: fields{
				pgType:   TypePGTimetz,
				array:    false,
				dims:     0,
				nullable: true,
			},
			want: "*time.Time",
		},
		{
			name: "Should generate struct type",
			fields: fields{
				pgType:   TypePGTimetz,
				array:    false,
				dims:     0,
				nullable: true,
				sqlNulls: true,
			},
			want: "bun.NullTime",
		},
		{
			name: "Should generate interface for unknown type",
			fields: fields{
				pgType:   "unknown",
				array:    false,
				dims:     0,
				nullable: true,
			},
			want: "interface{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewColumn("test", tt.fields.pgType, tt.fields.nullable, tt.fields.sqlNulls, tt.fields.array, tt.fields.dims, false, false, 0, []string{}, CustomTypeMapping{})
			if got := c.Type; got != tt.want {
				t.Errorf("Column.Type = %v, want %v", got, tt.want)
			}
		})
	}
}
