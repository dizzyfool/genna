package model

import (
	"reflect"
	"testing"
)

func Test_goType(t *testing.T) {
	tests := []struct {
		name    string
		pgTypes []string
		want    string
		wantErr bool
	}{
		{
			name:    "Should not get unknown type",
			pgTypes: []string{"unknown"},
			wantErr: true,
		},
		{
			name:    "Should get int",
			pgTypes: []string{TypePGInt2, TypePGInt4},
			want:    TypeInt,
		},
		{
			name:    "Should get int64",
			pgTypes: []string{TypePGInt8},
			want:    TypeInt64,
		},
		{
			name:    "Should get float32",
			pgTypes: []string{TypePGFloat4},
			want:    TypeFloat32,
		},
		{
			name:    "Should get float64",
			pgTypes: []string{TypePGNumeric, TypePGFloat8},
			want:    TypeFloat64,
		},
		{
			name:    "Should get string",
			pgTypes: []string{TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint},
			want:    TypeString,
		},
		{
			name:    "Should get byte",
			pgTypes: []string{TypePGBytea},
			want:    TypeByte,
		},
		{
			name:    "Should get bool",
			pgTypes: []string{TypePGBool},
			want:    TypeBool,
		},
		{
			name:    "Should get time.Time",
			pgTypes: []string{TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz},
			want:    TypeTime,
		},
		{
			name:    "Should get duration",
			pgTypes: []string{TypePGInterval},
			want:    TypeDuration,
		},
		{
			name:    "Should get map[string]interface{}",
			pgTypes: []string{TypePGJSONB, TypePGJSON},
			want:    TypeMapInterface,
		},
		{
			name:    "Should get map[string]string",
			pgTypes: []string{TypePGHstore},
			want:    TypeMapString,
		},
		{
			name:    "Should get netIP",
			pgTypes: []string{TypePGInet},
			want:    TypeIP,
		},
		{
			name:    "Should get netIPNet",
			pgTypes: []string{TypePGCidr},
			want:    TypeIPNet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, typ := range tt.pgTypes {
				got, err := GoType(typ)
				if (err != nil) != tt.wantErr {
					t.Errorf("GoType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GoType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_goSlice(t *testing.T) {
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
			args: args{TypePGInt4, 3},
			want: "[][][]int",
		},
		{
			name: "Should generate int2 array",
			args: args{TypePGInt2, 1},
			want: "[]int",
		},
		{
			name: "Should generate int4 array",
			args: args{TypePGInt4, 1},
			want: "[]int",
		},
		{
			name: "Should generate int8 array",
			args: args{TypePGInt8, 1},
			want: "[]int64",
		},
		{
			name: "Should generate numeric array",
			args: args{TypePGNumeric, 1},
			want: "[]float64",
		},
		{
			name: "Should generate float4 array",
			args: args{TypePGFloat4, 1},
			want: "[]float32",
		},
		{
			name: "Should generate float8 array",
			args: args{TypePGFloat8, 1},
			want: "[]float64",
		},
		{
			name: "Should generate text array",
			args: args{TypePGText, 1},
			want: "[]string",
		},
		{
			name: "Should generate varchar array",
			args: args{TypePGVarchar, 1},
			want: "[]string",
		},
		{
			name: "Should generate uuid array",
			args: args{TypePGUuid, 1},
			want: "[]string",
		},
		{
			name: "Should generate char array",
			args: args{TypePGBpchar, 1},
			want: "[]string",
		},
		{
			name: "Should generate bool array",
			args: args{TypePGBool, 1},
			want: "[]bool",
		},
		{
			name: "Should generate json array",
			args: args{TypePGJSON, 1},
			want: "[]map[string]interface{}",
		},
		{
			name: "Should generate jsonb array",
			args: args{TypePGJSONB, 1},
			want: "[]map[string]interface{}",
		},
		{
			name: "Should generate point array",
			args: args{TypePGPoint, 1},
			want: "[]string",
		},
		{
			name:    "Should not generate not supported type array",
			args:    args{TypePGTimetz, 1},
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
			got, err := GoSlice(tt.args.pgType, tt.args.dimensions)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("GoSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goNullable(t *testing.T) {
	tests := []struct {
		name          string
		pgType        string
		avoidPointers bool
		want          string
		wantErr       bool
	}{
		{
			name:   "Should generate int2 type",
			pgType: TypePGInt2,
			want:   "*int",
		},
		{
			name:   "Should generate int4 type",
			pgType: TypePGInt4,
			want:   "*int",
		},
		{
			name:   "Should generate int8 type",
			pgType: TypePGInt8,
			want:   "*int64",
		},
		{
			name:   "Should generate numeric type",
			pgType: TypePGNumeric,
			want:   "*float64",
		},
		{
			name:   "Should generate float4 type",
			pgType: TypePGFloat4,
			want:   "*float32",
		},
		{
			name:   "Should generate float8 type",
			pgType: TypePGFloat8,
			want:   "*float64",
		},
		{
			name:   "Should generate text type",
			pgType: TypePGText,
			want:   "*string",
		},
		{
			name:   "Should generate varchar type",
			pgType: TypePGVarchar,
			want:   "*string",
		},
		{
			name:   "Should generate uuid type",
			pgType: TypePGUuid,
			want:   "*string",
		},
		{
			name:   "Should generate char type",
			pgType: TypePGBpchar,
			want:   "*string",
		},
		{
			name:   "Should generate bool type",
			pgType: TypePGBool,
			want:   "*bool",
		},
		{
			name:   "Should generate time type",
			pgType: TypePGTimestamp,
			want:   "*time.Time",
		},
		{
			name:   "Should generate interval type",
			pgType: TypePGInterval,
			want:   "*time.Duration",
		},
		{
			name:   "Should generate json type",
			pgType: TypePGJSON,
			want:   "map[string]interface{}",
		},
		{
			name:   "Should generate hstore type",
			pgType: TypePGHstore,
			want:   "map[string]string",
		},
		{
			name:   "Should generate ip type",
			pgType: TypePGInet,
			want:   "*net.IP",
		},
		{
			name:   "Should generate cidr type",
			pgType: TypePGCidr,
			want:   "*net.IPNet",
		},
		{
			name:   "Should generate point type",
			pgType: TypePGPoint,
			want:   "*string",
		},
		{
			name:    "Should not generate unknown type",
			pgType:  "unknown",
			wantErr: true,
		},
		{
			name:          "Should generate int2 type avoiding pointers to sql.NullInt64",
			pgType:        TypePGInt2,
			avoidPointers: true,
			want:          "sql.NullInt64",
		},
		{
			name:          "Should generate varchar type avoiding pointers to sql.NullInt64",
			pgType:        TypePGVarchar,
			avoidPointers: true,
			want:          "sql.NullString",
		},
		{
			name:          "Should generate uuid type avoiding pointers to sql.NullInt64",
			pgType:        TypePGUuid,
			avoidPointers: true,
			want:          "sql.NullString",
		},
		{
			name:          "Should generate bool type avoiding pointers to sql.NullBool",
			pgType:        TypePGBool,
			avoidPointers: true,
			want:          "sql.NullBool",
		},
		{
			name:          "Should generate float64 type avoiding pointers to sql.NullFloat64",
			pgType:        TypePGFloat8,
			avoidPointers: true,
			want:          "sql.NullFloat64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoNullable(tt.pgType, tt.avoidPointers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoNullable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("GoNullable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goImport(t *testing.T) {
	type args struct {
		pgTypes       []string
		nullable      bool
		avoidPointers bool
		ver           int
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
					TypePGInt2, TypePGInt4, TypePGInt8, TypePGNumeric, TypePGFloat4, TypePGFloat8, TypePGBool, TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar,
				},
				ver: 8,
			},
			want: "",
		},
		{
			name: "Should not generate import for unknown type",
			args: args{
				pgTypes: []string{"unknown"},
				ver: 8,
			},
			want: "",
		},
		{
			name: "Should generate time import for interval type",
			args: args{
				pgTypes: []string{TypePGInterval},
				ver: 8,
			},
			want: "time",
		},
		{
			name: "Should generate net import for net types",
			args: args{
				pgTypes: []string{
					TypePGInet, TypePGCidr,
				},
				ver: 8,
			},
			want: "net",
		},
		{
			name: "Should generate net import for json types",
			args: args{
				pgTypes: []string{
					TypePGJSONB, TypePGJSON,
				},
				ver: 8,
			},
			want: "",
		},
		{
			name: "Should generate sql import for nullable simple types avoiding pointer",
			args: args{
				pgTypes: []string{
					TypePGInt2, TypePGInt4, TypePGInt8, TypePGNumeric, TypePGFloat4, TypePGFloat8, TypePGBool, TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar,
				},
				ver: 8,
				nullable:      true,
				avoidPointers: true,
			},
			want: "database/sql",
		},
		{
			name: "Should not generate sql import for nullable simple types",
			args: args{
				pgTypes: []string{
					TypePGInt2, TypePGInt4, TypePGInt8, TypePGNumeric, TypePGFloat4, TypePGFloat8, TypePGBool, TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar,
				},
				ver: 8,
				nullable:      true,
				avoidPointers: false,
			},
			want: "",
		},
		{
			name: "Should generate time import for nullable date time types",
			args: args{
				pgTypes: []string{
					TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz,
				},
				ver: 8,
				nullable: true,
			},
			want: "time",
		},
		{
			name: "Should generate go-pg import for nullable date time types",
			args: args{
				pgTypes: []string{
					TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz,
				},
				ver: 8,
				nullable:      true,
				avoidPointers: true,
			},
			want: "github.com/go-pg/pg",
		},
		{
			name: "Should generate go-pg import for nullable date time types",
			args: args{
				pgTypes: []string{
					TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz,
				},
				ver: 9,
				nullable:      true,
				avoidPointers: true,
			},
			want: "github.com/go-pg/pg/v9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, pgType := range tt.args.pgTypes {
				if got := GoImport(pgType, tt.args.nullable, tt.args.avoidPointers, tt.args.ver); got != tt.want {
					t.Errorf("GoImport() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_fixIsArray(t *testing.T) {
	type args struct {
		pgType     string
		isArray    bool
		dimensions int
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 int
	}{
		{
			name: "Should fix array for Bytea type",
			args: args{
				pgType:     TypePGBytea,
				isArray:    false,
				dimensions: 0,
			},
			want:  true,
			want1: 1,
		},
		{
			name: "Should fix array for Bytea array type",
			args: args{
				pgType:     TypePGBytea,
				isArray:    true,
				dimensions: 1,
			},
			want:  true,
			want1: 2,
		},
		{
			name: "Should not fix type",
			args: args{
				pgType:     TypeInt,
				isArray:    false,
				dimensions: 0,
			},
			want:  false,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := fixIsArray(tt.args.pgType, tt.args.isArray, tt.args.dimensions)
			if got != tt.want {
				t.Errorf("fixIsArray() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("fixIsArray() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
