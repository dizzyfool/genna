package model

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
			name:   "Should generate jsonb type",
			pgType: TypeJSONB,
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
			name:   "Should generate point type",
			pgType: TypePoint,
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
			name: "Should generate json array",
			args: args{TypeJSON, 1},
			want: "[]map[string]interface{}",
		},
		{
			name: "Should generate jsonb array",
			args: args{TypeJSONB, 1},
			want: "[]map[string]interface{}",
		},
		{
			name: "Should generate point array",
			args: args{TypePoint, 1},
			want: "[]string",
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
			name:   "Should generate point type",
			pgType: TypePoint,
			want:   "*string",
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
				pgTypes: []string{TypeInt2, TypeInt4},
			},
			want: types.NewPointer(types.Typ[types.Int]),
		},
		{
			name: "Should get int64",
			args: args{
				pgTypes: []string{TypeInt8},
			},
			want: types.NewPointer(types.Typ[types.Int64]),
		},
		{
			name: "Should get float32",
			args: args{
				pgTypes: []string{TypeNumeric, TypeFloat4},
			},
			want: types.NewPointer(types.Typ[types.Float32]),
		},
		{
			name: "Should get float64",
			args: args{
				pgTypes: []string{TypeFloat8},
			},
			want: types.NewPointer(types.Typ[types.Float64]),
		},
		{
			name: "Should get string",
			args: args{
				pgTypes: []string{TypeText, TypeVarchar, TypeUuid, TypeBpchar, TypePoint},
			},
			want: types.NewPointer(types.Typ[types.String]),
		},
		{
			name: "Should get []byte",
			args: args{
				pgTypes: []string{TypeBytea},
			},
			want: bytea{},
		},
		{
			name: "Should get bool",
			args: args{
				pgTypes: []string{TypeBool},
			},
			want: types.NewPointer(types.Typ[types.Bool]),
		},
		{
			name: "Should get time.Time",
			args: args{
				pgTypes: []string{TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz},
			},
			want: types.NewPointer(dateTime{}),
		},
		{
			name: "Should get duration",
			args: args{
				pgTypes: []string{TypeInterval},
			},
			want: types.NewPointer(interval{}),
		},
		{
			name: "Should get map[string]interface{}",
			args: args{
				pgTypes: []string{TypeJSONB, TypeJSON},
			},
			want: types.NewMap(types.Typ[types.String], intrfce{}),
		},
		{
			name: "Should get map[string]string",
			args: args{
				pgTypes: []string{TypeHstore},
			},
			want: types.NewMap(types.Typ[types.String], types.Typ[types.String]),
		},
		{
			name: "Should get netIP",
			args: args{
				pgTypes: []string{TypeInet},
			},
			want: types.NewPointer(netIP{}),
		},
		{
			name: "Should get netIPNet",
			args: args{
				pgTypes: []string{TypeCidr},
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
				pgTypes:    []string{TypeInt2, TypeInt4},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Int]),
		},
		{
			name: "Should get int64 array",
			args: args{
				pgTypes:    []string{TypeInt8},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Int64]),
		},
		{
			name: "Should get float32 array",
			args: args{
				pgTypes:    []string{TypeNumeric, TypeFloat4},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Float32]),
		},
		{
			name: "Should get float64 array",
			args: args{
				pgTypes:    []string{TypeFloat8},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Float64]),
		},
		{
			name: "Should get string array",
			args: args{
				pgTypes:    []string{TypeText, TypeVarchar, TypeUuid, TypeBpchar, TypePoint},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.String]),
		},
		{
			name: "Should get []byte array",
			args: args{
				pgTypes:    []string{TypeBytea},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(bytea{}),
		},
		{
			name: "Should get bool array",
			args: args{
				pgTypes:    []string{TypeBool},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.Typ[types.Bool]),
		},
		{
			name: "Should not get time.Time array",
			args: args{
				pgTypes:    []string{TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should not get duration array",
			args: args{
				pgTypes:    []string{TypeInterval},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should get map[string]interface{} array",
			args: args{
				pgTypes:    []string{TypeJSONB, TypeJSON},
				array:      true,
				dimensions: 1,
			},
			want: types.NewSlice(types.NewMap(types.Typ[types.String], intrfce{})),
		},
		{
			name: "Should not get map[string]string array",
			args: args{
				pgTypes:    []string{TypeHstore},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should not get netIP array",
			args: args{
				pgTypes:    []string{TypeInet},
				array:      true,
				dimensions: 1,
			},
			wantErr: true,
		},
		{
			name: "Should not get netIPNet array",
			args: args{
				pgTypes:    []string{TypeCidr},
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
			pgTypes: []string{TypeInt2, TypeInt4},
			want:    types.Typ[types.Int],
		},
		{
			name:    "Should get int64",
			pgTypes: []string{TypeInt8},
			want:    types.Typ[types.Int64],
		},
		{
			name:    "Should get float32",
			pgTypes: []string{TypeNumeric, TypeFloat4},
			want:    types.Typ[types.Float32],
		},
		{
			name:    "Should get float64",
			pgTypes: []string{TypeFloat8},
			want:    types.Typ[types.Float64],
		},
		{
			name:    "Should get string",
			pgTypes: []string{TypeText, TypeVarchar, TypeUuid, TypeBpchar, TypePoint},
			want:    types.Typ[types.String],
		},
		{
			name:    "Should get []byte",
			pgTypes: []string{TypeBytea},
			want:    bytea{},
		},
		{
			name:    "Should get bool",
			pgTypes: []string{TypeBool},
			want:    types.Typ[types.Bool],
		},
		{
			name:    "Should get time.Time",
			pgTypes: []string{TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz},
			want:    dateTime{},
		},
		{
			name:    "Should get duration",
			pgTypes: []string{TypeInterval},
			want:    interval{},
		},
		{
			name:    "Should get map[string]interface{}",
			pgTypes: []string{TypeJSONB, TypeJSON},
			want:    types.NewMap(types.Typ[types.String], intrfce{}),
		},
		{
			name:    "Should get map[string]string",
			pgTypes: []string{TypeHstore},
			want:    types.NewMap(types.Typ[types.String], types.Typ[types.String]),
		},
		{
			name:    "Should get netIP",
			pgTypes: []string{TypeInet},
			want:    netIP{},
		},
		{
			name:    "Should get netIPNet",
			pgTypes: []string{TypeCidr},
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
					TypeInt2,
					TypeInt4,
					TypeInt8,
					TypeNumeric,
					TypeFloat4,
					TypeFloat8,
					TypeText,
					TypeVarchar,
					TypeUuid,
					TypeBpchar,
					TypeBool,
					TypeTimestamp,
					TypeTimestamptz,
					TypeDate,
					TypeTime,
					TypeTimetz,
					TypeInterval,
					TypeInet,
					TypeCidr,
					TypePoint,
				},
			},
			want: true,
		},
		{
			name: "Should return false for non basic types",
			args: args{
				pgTypes: []string{
					TypeJSON,
					TypeJSONB,
					TypeHstore,
					TypeBytea,
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
