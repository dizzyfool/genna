package model_old

import (
	"go/types"
	"reflect"
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
			pgType: TypePGInt2,
			want:   "int",
		},
		{
			name:   "Should generate int4 type",
			pgType: TypePGInt4,
			want:   "int",
		},
		{
			name:   "Should generate int8 type",
			pgType: TypePGInt8,
			want:   "int64",
		},
		{
			name:   "Should generate numeric type",
			pgType: TypePGNumeric,
			want:   "float64",
		},
		{
			name:   "Should generate float4 type",
			pgType: TypePGPloat4,
			want:   "float32",
		},
		{
			name:   "Should generate float8 type",
			pgType: TypePGPloat8,
			want:   "float64",
		},
		{
			name:   "Should generate text type",
			pgType: TypePGPext,
			want:   "string",
		},
		{
			name:   "Should generate varchar type",
			pgType: TypePGParchar,
			want:   "string",
		},
		{
			name:   "Should generate uuid type",
			pgType: TypePGPuid,
			want:   "string",
		},
		{
			name:   "Should generate char type",
			pgType: TypePGPpchar,
			want:   "string",
		},
		{
			name:   "Should generate bytea type",
			pgType: TypePGPytea,
			want:   "[]byte",
		},
		{
			name:   "Should generate bool type",
			pgType: TypePGPool,
			want:   "bool",
		},
		{
			name:   "Should generate time type",
			pgType: TypePGPimestamp,
			want:   "time.Time",
		},
		{
			name:   "Should generate interval type",
			pgType: TypePGPnterval,
			want:   "time.Duration",
		},
		{
			name:   "Should generate json type",
			pgType: TypePGPSON,
			want:   "map[string]interface{}",
		},
		{
			name:   "Should generate jsonb type",
			pgType: TypePGPSONB,
			want:   "map[string]interface{}",
		},
		{
			name:   "Should generate hstore type",
			pgType: TypePGPstore,
			want:   "map[string]string",
		},
		{
			name:   "Should generate ip type",
			pgType: TypePGPnet,
			want:   "net.IP",
		},
		{
			name:   "Should generate cidr type",
			pgType: TypePGPidr,
			want:   "net.IPNet",
		},
		{
			name:   "Should generate point type",
			pgType: TypePGPoint,
			want:   "string",
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
			args: args{TypePGPloat4, 1},
			want: "[]float32",
		},
		{
			name: "Should generate float8 array",
			args: args{TypePGPloat8, 1},
			want: "[]float64",
		},
		{
			name: "Should generate text array",
			args: args{TypePGPext, 1},
			want: "[]string",
		},
		{
			name: "Should generate varchar array",
			args: args{TypePGParchar, 1},
			want: "[]string",
		},
		{
			name: "Should generate uuid array",
			args: args{TypePGPuid, 1},
			want: "[]string",
		},
		{
			name: "Should generate char array",
			args: args{TypePGPpchar, 1},
			want: "[]string",
		},
		{
			name: "Should generate bytea array",
			args: args{TypePGPytea, 1},
			want: "[][]byte",
		},
		{
			name: "Should generate bool array",
			args: args{TypePGPool, 1},
			want: "[]bool",
		},
		{
			name: "Should generate json array",
			args: args{TypePGPSON, 1},
			want: "[]map[string]interface{}",
		},
		{
			name: "Should generate jsonb array",
			args: args{TypePGPSONB, 1},
			want: "[]map[string]interface{}",
		},
		{
			name: "Should generate point array",
			args: args{TypePGPoint, 1},
			want: "[]string",
		},
		{
			name:    "Should not generate not supported type array",
			args:    args{TypePGPimetz, 1},
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
			pgType: TypePGPloat4,
			want:   "*float32",
		},
		{
			name:   "Should generate float8 type",
			pgType: TypePGPloat8,
			want:   "*float64",
		},
		{
			name:   "Should generate text type",
			pgType: TypePGPext,
			want:   "*string",
		},
		{
			name:   "Should generate varchar type",
			pgType: TypePGParchar,
			want:   "*string",
		},
		{
			name:   "Should generate uuid type",
			pgType: TypePGPuid,
			want:   "*string",
		},
		{
			name:   "Should generate char type",
			pgType: TypePGPpchar,
			want:   "*string",
		},
		{
			name:   "Should generate bytea type",
			pgType: TypePGPytea,
			want:   "*[]byte",
		},
		{
			name:   "Should generate bool type",
			pgType: TypePGPool,
			want:   "*bool",
		},
		{
			name:   "Should generate time type",
			pgType: TypePGPimestamp,
			want:   "pg.NullTime",
		},
		{
			name:   "Should generate interval type",
			pgType: TypePGPnterval,
			want:   "*time.Duration",
		},
		{
			name:   "Should generate json type",
			pgType: TypePGPSON,
			want:   "map[string]interface{}",
		},
		{
			name:   "Should generate hstore type",
			pgType: TypePGPstore,
			want:   "map[string]string",
		},
		{
			name:   "Should generate ip type",
			pgType: TypePGPnet,
			want:   "*net.IP",
		},
		{
			name:   "Should generate cidr type",
			pgType: TypePGPidr,
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
			pgType:        TypePGParchar,
			avoidPointers: true,
			want:          "sql.NullString",
		},
		{
			name:          "Should generate uuid type avoiding pointers to sql.NullInt64",
			pgType:        TypePGPuid,
			avoidPointers: true,
			want:          "sql.NullString",
		},
		{
			name:          "Should generate bool type avoiding pointers to sql.NullBool",
			pgType:        TypePGPool,
			avoidPointers: true,
			want:          "sql.NullBool",
		},
		{
			name:          "Should generate float64 type avoiding pointers to sql.NullFloat64",
			pgType:        TypePGPloat8,
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
					TypePGInt2, TypePGInt4, TypePGInt8, TypePGNumeric, TypePGPloat4, TypePGPloat8, TypePGPool, TypePGPext, TypePGParchar, TypePGPuid, TypePGPpchar,
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
				pgTypes: []string{TypePGPnterval},
			},
			want: "time",
		},
		{
			name: "Should generate net import for net types",
			args: args{
				pgTypes: []string{
					TypePGPnet, TypePGPidr,
				},
			},
			want: "net",
		},
		{
			name: "Should generate net import for json types",
			args: args{
				pgTypes: []string{
					TypePGPSONB, TypePGJSON,
				},
			},
			want: "",
		},
		{
			name: "Should generate sql import for nullable simple types avoiding pointer",
			args: args{
				pgTypes: []string{
					TypePGInt2, TypePGInt4, TypePGInt8, TypePGNumeric, TypePGFloat4, TypePGFloat8, TypePGBool, TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar,
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
					TypePGInt2, TypePGInt4, TypePGInt8, TypePGNumeric, TypePGFloat4, TypePGFloat8, TypePGBool, TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar,
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
					TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz,
				},
				nullable: true,
			},
			want: "github.com/go-pg/pg",
		},
		{
			name: "Should generate time import for nullable date time types",
			args: args{
				pgTypes: []string{
					TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz,
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

func TestGoImportFromType(t *testing.T) {
	tests := []struct {
		name  string
		types []types.Type
		want  string
	}{
		{
			name: "Should not get import for std types",
			types: []types.Type{
				types.Typ[types.Int],
				types.Typ[types.Int64],
				types.Typ[types.Float32],
				types.Typ[types.Float64],
				types.Typ[types.String],
				types.Typ[types.Bool],
				types.NewMap(types.Typ[types.String], intrfce{}),
				types.NewMap(types.Typ[types.String], types.Typ[types.String]),
				types.NewSlice(types.Typ[types.Int]),
			},
			want: "",
		},
		{
			name: "Should not get import for time",
			types: []types.Type{
				dateTime{},
				&dateTime{},
				interval{},
				&interval{},
			},
			want: "time",
		},
		{
			name: "Should not get import for net",
			types: []types.Type{
				netIP{},
				&netIP{},
				netIPNet{},
				&netIPNet{},
			},
			want: "net",
		},
		{
			name: "Should not get import for nullable pg helper",
			types: []types.Type{
				pgNullTime{},
				&pgNullTime{},
			},
			want: "github.com/go-pg/pg",
		},
		{
			name: "Should not get import for nullable sql helper",
			types: []types.Type{
				sqlNullInt64{},
				&sqlNullInt64{},
				sqlNullFloat64{},
				&sqlNullFloat64{},
				sqlNullBool{},
				&sqlNullBool{},
				sqlNullString{},
				&sqlNullString{},
			},
			want: "database/sql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, typ := range tt.types {
				if got := GoImportFromType(typ); got != tt.want {
					t.Errorf("GoImportFromType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGoPointerType(t *testing.T) {
	type args struct {
		pgTypes    []string
		array      bool
		dimensions int
	}
	tests := []struct {
		name    string
		args    args
		want    types.Type
		wantErr bool
	}{
		{
			name: "Should not get unknown type",
			args: args{
				pgTypes: []string{"unknown"},
			},
			wantErr: true,
		},
		{
			name: "Should get int",
			args: args{
				pgTypes: []string{TypePGInt2, TypePGInt4},
			},
			want: types.NewPointer(types.Typ[types.Int]),
		},
		{
			name: "Should get int64",
			args: args{
				pgTypes: []string{TypePGInt8},
			},
			want: types.NewPointer(types.Typ[types.Int64]),
		},
		{
			name: "Should get float32",
			args: args{
				pgTypes: []string{TypePGFloat4},
			},
			want: types.NewPointer(types.Typ[types.Float32]),
		},
		{
			name: "Should get float64",
			args: args{
				pgTypes: []string{TypePGNumeric, TypePGFloat8},
			},
			want: types.NewPointer(types.Typ[types.Float64]),
		},
		{
			name: "Should get string",
			args: args{
				pgTypes: []string{TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint},
			},
			want: types.NewPointer(types.Typ[types.String]),
		},
		{
			name: "Should get []byte",
			args: args{
				pgTypes: []string{TypePGBytea},
			},
			want: bytea{},
		},
		{
			name: "Should get bool",
			args: args{
				pgTypes: []string{TypePGBool},
			},
			want: types.NewPointer(types.Typ[types.Bool]),
		},
		{
			name: "Should get time.Time",
			args: args{
				pgTypes: []string{TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz},
			},
			want: types.NewPointer(dateTime{}),
		},
		{
			name: "Should get duration",
			args: args{
				pgTypes: []string{TypePGInterval},
			},
			want: types.NewPointer(interval{}),
		},
		{
			name: "Should get map[string]interface{}",
			args: args{
				pgTypes: []string{TypePGJSONB, TypePGJSON},
			},
			want: types.NewMap(types.Typ[types.String], intrfce{}),
		},
		{
			name: "Should get map[string]string",
			args: args{
				pgTypes: []string{TypePGHstore},
			},
			want: types.NewMap(types.Typ[types.String], types.Typ[types.String]),
		},
		{
			name: "Should get netIP",
			args: args{
				pgTypes: []string{TypePGInet},
			},
			want: types.NewPointer(netIP{}),
		},
		{
			name: "Should get netIPNet",
			args: args{
				pgTypes: []string{TypePGCidr},
			},
			want: types.NewPointer(netIPNet{}),
		},
		// arrays
		{
			name: "Should not get unknown array type",
			args: args{
				pgTypes:    []string{"unknown"},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should get int array",
			args: args{
				pgTypes:    []string{TypePGInt2, TypePGInt4},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Int]),
		},
		{
			name: "Should get int64 array",
			args: args{
				pgTypes:    []string{TypePGInt8},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Int64]),
		},
		{
			name: "Should get float32 array",
			args: args{
				pgTypes:    []string{TypePGFloat4},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Float32]),
		},
		{
			name: "Should get float64 array",
			args: args{
				pgTypes:    []string{TypePGNumeric, TypePGFloat8},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Float64]),
		},
		{
			name: "Should get string array",
			args: args{
				pgTypes:    []string{TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.String]),
		},
		{
			name: "Should get []byte array",
			args: args{
				pgTypes:    []string{TypePGBytea},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(bytea{}),
		},
		{
			name: "Should get bool array",
			args: args{
				pgTypes:    []string{TypePGBool},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Bool]),
		},
		{
			name: "Should not get time.Time array",
			args: args{
				pgTypes:    []string{TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should not get duration array",
			args: args{
				pgTypes:    []string{TypePGInterval},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should get map[string]interface{} array",
			args: args{
				pgTypes:    []string{TypePGJSONB, TypePGJSON},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.NewMap(types.Typ[types.String], intrfce{})),
		},
		{
			name: "Should not get map[string]string array",
			args: args{
				pgTypes:    []string{TypePGHstore},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should not get netIP array",
			args: args{
				pgTypes:    []string{TypePGInet},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should not get netIPNet array",
			args: args{
				pgTypes:    []string{TypePGCidr},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, typ := range tt.args.pgTypes {
				got, err := GoPointerType(typ, tt.args.array, tt.args.dimensions)
				if (err != nil) != tt.wantErr {
					t.Errorf("GoPointerType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GoPointerType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGoSimpleType(t *testing.T) {
	tests := []struct {
		name    string
		pgTypes []string
		want    types.Type
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
			want:    types.Typ[types.Int],
		},
		{
			name:    "Should get int64",
			pgTypes: []string{TypePGInt8},
			want:    types.Typ[types.Int64],
		},
		{
			name:    "Should get float32",
			pgTypes: []string{TypePGFloat4},
			want:    types.Typ[types.Float32],
		},
		{
			name:    "Should get float64",
			pgTypes: []string{TypePGNumeric, TypePGFloat8},
			want:    types.Typ[types.Float64],
		},
		{
			name:    "Should get string",
			pgTypes: []string{TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint},
			want:    types.Typ[types.String],
		},
		{
			name:    "Should get []byte",
			pgTypes: []string{TypePGBytea},
			want:    bytea{},
		},
		{
			name:    "Should get bool",
			pgTypes: []string{TypePGBool},
			want:    types.Typ[types.Bool],
		},
		{
			name:    "Should get time.Time",
			pgTypes: []string{TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz},
			want:    dateTime{},
		},
		{
			name:    "Should get duration",
			pgTypes: []string{TypePGInterval},
			want:    interval{},
		},
		{
			name:    "Should get map[string]interface{}",
			pgTypes: []string{TypePGJSONB, TypePGJSON},
			want:    types.NewMap(types.Typ[types.String], intrfce{}),
		},
		{
			name:    "Should get map[string]string",
			pgTypes: []string{TypePGHstore},
			want:    types.NewMap(types.Typ[types.String], types.Typ[types.String]),
		},
		{
			name:    "Should get netIP",
			pgTypes: []string{TypePGInet},
			want:    netIP{},
		},
		{
			name:    "Should get netIPNet",
			pgTypes: []string{TypePGCidr},
			want:    netIPNet{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, typ := range tt.pgTypes {
				got, err := GoSimpleType(typ)
				if (err != nil) != tt.wantErr {
					t.Errorf("GoSimpleType() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GoSimpleType() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestIsBasic(t *testing.T) {
	type args struct {
		pgTypes []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true for basic types",
			args: args{
				pgTypes: []string{
					TypePGInt2,
					TypePGInt4,
					TypePGInt8,
					TypePGNumeric,
					TypePGFloat4,
					TypePGFloat8,
					TypePGText,
					TypePGVarchar,
					TypePGUuid,
					TypePGBpchar,
					TypePGBool,
					TypePGTimestamp,
					TypePGTimestamptz,
					TypePGDate,
					TypePGTime,
					TypePGTimetz,
					TypePGInterval,
					TypePGInet,
					TypePGCidr,
					TypePGPoint,
				},
			},
			want: true,
		},
		{
			name: "Should return false for non basic types",
			args: args{
				pgTypes: []string{
					TypePGJSON,
					TypePGJSONB,
					TypePGHstore,
					TypePGBytea,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, typ := range tt.args.pgTypes {
				if got := IsBasic(typ); got != tt.want {
					t.Errorf("IsBasic() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
