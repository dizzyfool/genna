package model

import (
	"testing"
)

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
				SourceColumn: tt.fields.SourceColumn,
			}
			if got := r.StructFieldTag(); got != tt.want {
				t.Errorf("Relation.StructFieldTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelation_Import(t *testing.T) {
	type fields struct {
		SourceSchema string
		TargetSchema string
	}
	type args struct {
		importPath  string
		publicAlias string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Should not generate import inside same package",
			fields: fields{
				SourceSchema: "geo",
				TargetSchema: "geo",
			},
			args: args{},
			want: "",
		},
		{
			name: "Should generate import with default name",
			fields: fields{
				SourceSchema: "geo",
				TargetSchema: PublicSchema,
			},
			args: args{},
			want: "model",
		},
		{
			name: "Should generate import with custom name",
			fields: fields{
				SourceSchema: "geo",
				TargetSchema: PublicSchema,
			},
			args: args{
				publicAlias: "test",
			},
			want: "test",
		},
		{
			name: "Should generate import with custom name and prefix",
			fields: fields{
				SourceSchema: "geo",
				TargetSchema: PublicSchema,
			},
			args: args{
				publicAlias: "test",
				importPath:  "model",
			},
			want: "model/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Relation{
				SourceSchema: tt.fields.SourceSchema,
				TargetSchema: tt.fields.TargetSchema,
			}
			if got := r.Import(tt.args.importPath, tt.args.publicAlias); got != tt.want {
				t.Errorf("Relation.Import() = %v, want %v", got, tt.want)
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
	type args struct {
		withSchema  bool
		publicAlias string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Should generate from simple word",
			fields: fields{
				SourceSchema: PublicSchema,
				TargetSchema: PublicSchema,
				TargetTable:  "users",
			},
			args: args{},
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
			name: "Should generate without package",
			fields: fields{
				SourceSchema: PublicSchema,
				TargetSchema: "information_schema",
				TargetTable:  "users",
			},
			want: "*information_schema.User",
		},
		{
			name: "Should generate with schema",
			fields: fields{
				SourceSchema: PublicSchema,
				TargetSchema: "information_schema",
				TargetTable:  "users",
			},
			args: args{
				withSchema: true,
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
			args: args{
				withSchema: true,
			},
			want: "*User",
		},
		{
			name: "Should generate with package",
			fields: fields{
				SourceSchema: "information_schema",
				TargetSchema: PublicSchema,
				TargetTable:  "users",
			},
			want: "*model.User",
		},
		{
			name: "Should generate with package alias",
			fields: fields{
				SourceSchema: "information_schema",
				TargetSchema: PublicSchema,
				TargetTable:  "users",
			},
			args: args{
				publicAlias: "geo",
			},
			want: "*geo.User",
		},
		{
			name: "Should generate with ignored package alias",
			fields: fields{
				SourceSchema: PublicSchema,
				TargetSchema: "information_schema",
				TargetTable:  "users",
			},
			args: args{
				publicAlias: "ignored",
			},
			want: "*information_schema.User",
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
			if got := r.StructFieldType(tt.args.withSchema, tt.args.publicAlias); got != tt.want {
				t.Errorf("Relation.StructFieldType() = %v, want %v", got, tt.want)
			}
		})
	}
}
