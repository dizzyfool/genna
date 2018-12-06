package model

import "testing"

func TestColumn_StructFieldName(t *testing.T) {
	type fields struct {
		Name string
		IsPK bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should generate from simple word",
			fields: fields{Name: "title"},
			want:   "Title",
		},
		{
			name:   "Should generate from underscored",
			fields: fields{Name: "short_title"},
			want:   "ShortTitle",
		},
		{
			name:   "Should generate from camelCased",
			fields: fields{Name: "shortTitle"},
			want:   "ShortTitle",
		},
		{
			name:   "Should generate with underscored id",
			fields: fields{Name: "location_id"},
			want:   "LocationID",
		},
		{
			name:   "Should generate with camelCased id",
			fields: fields{Name: "locationId"},
			want:   "LocationID",
		},
		{
			name:   "Should generate primary key as ID",
			fields: fields{Name: "locationId", IsPK: true},
			want:   "ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{
				Name: tt.fields.Name,
				IsPK: tt.fields.IsPK,
			}
			if got := c.StructFieldName(); got != tt.want {
				t.Errorf("Column.StructFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumn_StructFieldType(t *testing.T) {
	type fields struct {
		Type       string
		IsArray    bool
		Dimensions int
		IsNullable bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// See TestGoType, TestGoSliceType, TestGoNullType for full test cases
		{
			name: "Should generate int2 type",
			fields: fields{
				Type:       TypeInt2,
				IsArray:    false,
				Dimensions: 0,
				IsNullable: false,
			},
			want: "int",
		},
		{
			name: "Should generate int2 array type",
			fields: fields{
				Type:       TypeInt2,
				IsArray:    true,
				Dimensions: 2,
				IsNullable: false,
			},
			want: "[][]int",
		},
		{
			name: "Should generate int2 nullable type",
			fields: fields{
				Type:       TypeInt2,
				IsArray:    true,
				Dimensions: 2,
				IsNullable: true,
			},
			want: "*[][]int",
		},
		{
			name: "Should generate struct type",
			fields: fields{
				Type:       TypeTimetz,
				IsArray:    false,
				Dimensions: 0,
				IsNullable: true,
			},
			want: "pg.NullTime",
		},
		{
			name: "Should generate interface for unknown type",
			fields: fields{
				Type:       "unknown",
				IsArray:    false,
				Dimensions: 0,
				IsNullable: true,
			},
			want: "interface{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{
				Type:       tt.fields.Type,
				IsArray:    tt.fields.IsArray,
				Dimensions: tt.fields.Dimensions,
				IsNullable: tt.fields.IsNullable,
			}
			if got := c.StructFieldType(); got != tt.want {
				t.Errorf("Column.StructFieldType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumn_StructFieldTag(t *testing.T) {
	type fields struct {
		Name       string
		Type       string
		IsArray    bool
		IsNullable bool
		IsPK       bool
		IsFK       bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should generate simple column",
			fields: fields{
				Name: "title",
				Type: TypeVarchar,
			},
			want: `sql:"title,notnull"`,
		},
		{
			name: "Should generate primary key",
			fields: fields{
				Name:       "userId",
				Type:       TypeInt2,
				IsPK:       true,
				IsNullable: true, // should ignore that
			},
			want: `sql:"userId,pk"`,
		},
		{
			name: "Should generate nullable column",
			fields: fields{
				Name:       "createdAt",
				Type:       TypeTimetz,
				IsNullable: true,
			},
			want: `sql:"createdAt"`,
		},
		{
			name: "Should generate array nullable column",
			fields: fields{
				Name:       "flags",
				Type:       TypeInt2,
				IsArray:    true,
				IsNullable: true,
			},
			want: `sql:"flags,array"`,
		},
		{
			name: "Should generate array column",
			fields: fields{
				Name:    "flags",
				Type:    TypeInt2,
				IsArray: true,
			},
			want: `sql:"flags,array,notnull"`,
		},
		{
			name: "Should generate hstore column",
			fields: fields{
				Name: "flags",
				Type: TypeHstore,
			},
			want: `sql:"flags,hstore,notnull"`,
		},
		{
			name: "Should generate unknown ignored column",
			fields: fields{
				Name: "flags",
				Type: "unknown",
			},
			want: `sql:"-"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Column{
				Name:       tt.fields.Name,
				Type:       tt.fields.Type,
				IsArray:    tt.fields.IsArray,
				IsNullable: tt.fields.IsNullable,
				IsPK:       tt.fields.IsPK,
				IsFK:       tt.fields.IsFK,
			}
			if got := c.StructFieldTag(); got != tt.want {
				t.Errorf("Column.StructFieldTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
