package model

import (
	"testing"
)

func TestGoType(t *testing.T) {
	tests := []struct {
		name    string
		pgType  string
		want    string
		wantErr bool
	}{
		{
			name:   "Should generate int2 type",
			pgType: TypeInt2,
			want:   "int",
		},
		{
			name:   "Should generate int4 type",
			pgType: TypeInt4,
			want:   "int",
		},
		{
			name:   "Should generate int8 type",
			pgType: TypeInt8,
			want:   "int64",
		},
		{
			name:   "Should generate numeric type",
			pgType: TypeNumeric,
			want:   "float32",
		},
		{
			name:   "Should generate float4 type",
			pgType: TypeFloat4,
			want:   "float32",
		},
		{
			name:   "Should generate float8 type",
			pgType: TypeFloat8,
			want:   "float64",
		},
		{
			name:   "Should generate text type",
			pgType: TypeText,
			want:   "string",
		},
		{
			name:   "Should generate varchar type",
			pgType: TypeVarchar,
			want:   "string",
		},
		{
			name:   "Should generate uuid type",
			pgType: TypeUuid,
			want:   "string",
		},
		{
			name:   "Should generate char type",
			pgType: TypeBpchar,
			want:   "string",
		},
		{
			name:   "Should generate bytea type",
			pgType: TypeBytea,
			want:   "[]byte",
		},
		{
			name:   "Should generate bool type",
			pgType: TypeBool,
			want:   "bool",
		},
		{
			name:   "Should generate time type",
			pgType: TypeTimestamp,
			want:   "time.Time",
		},
		{
			name:   "Should generate interval type",
			pgType: TypeInterval,
			want:   "time.Duration",
		},
		{
			name:   "Should generate json type",
			pgType: TypeJSON,
			want:   "map[string]interface{}",
		},
		{
			name:   "Should generate hstore type",
			pgType: TypeHstore,
			want:   "map[string]string",
		},
		{
			name:   "Should generate ip type",
			pgType: TypeInet,
			want:   "net.IP",
		},
		{
			name:   "Should generate cidr type",
			pgType: TypeCidr,
			want:   "net.IPNet",
		},
		{
			name:    "Should not generate unknown type",
			pgType:  "unknown",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoSimpleType(tt.pgType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoSimpleType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.want {
				t.Errorf("GoSimpleType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoSliceType(t *testing.T) {
	type args struct {
		pgType     string
		dimensions int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Should generate multi-dimension array",
			args: args{TypeInt4, 3},
			want: "[][][]int",
		},
		{
			name: "Should generate int2 array",
			args: args{TypeInt2, 1},
			want: "[]int",
		},
		{
			name: "Should generate int4 array",
			args: args{TypeInt4, 1},
			want: "[]int",
		},
		{
			name: "Should generate int8 array",
			args: args{TypeInt8, 1},
			want: "[]int64",
		},
		{
			name: "Should generate numeric array",
			args: args{TypeNumeric, 1},
			want: "[]float32",
		},
		{
			name: "Should generate float4 array",
			args: args{TypeFloat4, 1},
			want: "[]float32",
		},
		{
			name: "Should generate float8 array",
			args: args{TypeFloat8, 1},
			want: "[]float64",
		},
		{
			name: "Should generate text array",
			args: args{TypeText, 1},
			want: "[]string",
		},
		{
			name: "Should generate varchar array",
			args: args{TypeVarchar, 1},
			want: "[]string",
		},
		{
			name: "Should generate uuid array",
			args: args{TypeUuid, 1},
			want: "[]string",
		},
		{
			name: "Should generate char array",
			args: args{TypeBpchar, 1},
			want: "[]string",
		},
		{
			name: "Should generate bytea array",
			args: args{TypeBytea, 1},
			want: "[][]byte",
		},
		{
			name: "Should generate bool array",
			args: args{TypeBool, 1},
			want: "[]bool",
		},
		{
			name:    "Should not generate not supported type array",
			args:    args{TypeTimetz, 1},
			wantErr: true,
		},
		{
			name:    "Should not generate unknown type array",
			args:    args{"unknown", 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoSliceType(tt.args.pgType, tt.args.dimensions)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoSliceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.want {
				t.Errorf("GoSliceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoNullType(t *testing.T) {
	tests := []struct {
		name          string
		pgType        string
		avoidPointers bool
		want          string
		wantErr       bool
	}{
		{
			name:   "Should generate int2 type",
			pgType: TypeInt2,
			want:   "*int",
		},
		{
			name:   "Should generate int4 type",
			pgType: TypeInt4,
			want:   "*int",
		},
		{
			name:   "Should generate int8 type",
			pgType: TypeInt8,
			want:   "*int64",
		},
		{
			name:   "Should generate numeric type",
			pgType: TypeNumeric,
			want:   "*float32",
		},
		{
			name:   "Should generate float4 type",
			pgType: TypeFloat4,
			want:   "*float32",
		},
		{
			name:   "Should generate float8 type",
			pgType: TypeFloat8,
			want:   "*float64",
		},
		{
			name:   "Should generate text type",
			pgType: TypeText,
			want:   "*string",
		},
		{
			name:   "Should generate varchar type",
			pgType: TypeVarchar,
			want:   "*string",
		},
		{
			name:   "Should generate uuid type",
			pgType: TypeUuid,
			want:   "*string",
		},
		{
			name:   "Should generate char type",
			pgType: TypeBpchar,
			want:   "*string",
		},
		{
			name:   "Should generate bytea type",
			pgType: TypeBytea,
			want:   "*[]byte",
		},
		{
			name:   "Should generate bool type",
			pgType: TypeBool,
			want:   "*bool",
		},
		{
			name:   "Should generate time type",
			pgType: TypeTimestamp,
			want:   "pg.NullTime",
		},
		{
			name:   "Should generate interval type",
			pgType: TypeInterval,
			want:   "*time.Duration",
		},
		{
			name:   "Should generate json type",
			pgType: TypeJSON,
			want:   "map[string]interface{}",
		},
		{
			name:   "Should generate hstore type",
			pgType: TypeHstore,
			want:   "map[string]string",
		},
		{
			name:   "Should generate ip type",
			pgType: TypeInet,
			want:   "*net.IP",
		},
		{
			name:   "Should generate cidr type",
			pgType: TypeCidr,
			want:   "*net.IPNet",
		},
		{
			name:    "Should not generate unknown type",
			pgType:  "unknown",
			wantErr: true,
		},
		{
			name:          "Should generate int2 type avoiding pointers to sql.NullInt64",
			pgType:        TypeInt2,
			avoidPointers: true,
			want:          "sql.NullInt64",
		},
		{
			name:          "Should generate varchar type avoiding pointers to sql.NullInt64",
			pgType:        TypeVarchar,
			avoidPointers: true,
			want:          "sql.NullString",
		},
		{
			name:          "Should generate uuid type avoiding pointers to sql.NullInt64",
			pgType:        TypeUuid,
			avoidPointers: true,
			want:          "sql.NullString",
		},
		{
			name:          "Should generate bool type avoiding pointers to sql.NullBool",
			pgType:        TypeBool,
			avoidPointers: true,
			want:          "sql.NullBool",
		},
		{
			name:          "Should generate float64 type avoiding pointers to sql.NullFloat64",
			pgType:        TypeFloat8,
			avoidPointers: true,
			want:          "sql.NullFloat64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoNullType(tt.pgType, tt.avoidPointers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoNullType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.want {
				t.Errorf("GoNullType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoImport(t *testing.T) {
	type args struct {
		pgTypes       []string
		nullable      bool
		avoidPointers bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should not generate import for simple type",
			args: args{
				pgTypes: []string{
					TypeInt2, TypeInt4, TypeInt8, TypeNumeric, TypeFloat4, TypeFloat8, TypeBool, TypeText, TypeVarchar, TypeUuid, TypeBpchar,
				},
			},
			want: "",
		},
		{
			name: "Should not generate import for unknown type",
			args: args{
				pgTypes: []string{"unknown"},
			},
			want: "",
		},
		{
			name: "Should generate time import for interval type",
			args: args{
				pgTypes: []string{TypeInterval},
			},
			want: "time",
		},
		{
			name: "Should generate net import for net types",
			args: args{
				pgTypes: []string{
					TypeInet, TypeCidr,
				},
			},
			want: "net",
		},
		{
			name: "Should generate net import for json types",
			args: args{
				pgTypes: []string{
					TypeJSONB, TypeJSON,
				},
			},
			want: "",
		},
		{
			name: "Should generate sql import for nullable simple types avoiding pointer",
			args: args{
				pgTypes: []string{
					TypeInt2, TypeInt4, TypeInt8, TypeNumeric, TypeFloat4, TypeFloat8, TypeBool, TypeText, TypeVarchar, TypeUuid, TypeBpchar,
				},
				nullable:      true,
				avoidPointers: true,
			},
			want: "database/sql",
		},
		{
			name: "Should not generate sql import for nullable simple types",
			args: args{
				pgTypes: []string{
					TypeInt2, TypeInt4, TypeInt8, TypeNumeric, TypeFloat4, TypeFloat8, TypeBool, TypeText, TypeVarchar, TypeUuid, TypeBpchar,
				},
				nullable:      true,
				avoidPointers: false,
			},
			want: "",
		},
		{
			name: "Should generate go-pg import for nullable date time types",
			args: args{
				pgTypes: []string{
					TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz,
				},
				nullable: true,
			},
			want: "github.com/go-pg/pg",
		},
		{
			name: "Should generate time import for nullable date time types",
			args: args{
				pgTypes: []string{
					TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz,
				},
			},
			want: "time",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, pgType := range tt.args.pgTypes {
				if got := GoImport(pgType, tt.args.nullable, false, 0, tt.args.avoidPointers); got != tt.want {
					t.Errorf("GoImport() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
