package model

import "go/types"

const (
	TypeInt2        = "int2"
	TypeInt4        = "int4"
	TypeInt8        = "int8"
	TypeNumeric     = "numeric"
	TypeFloat4      = "float4"
	TypeFloat8      = "float8"
	TypeText        = "text"
	TypeVarchar     = "varchar"
	TypeBpchar      = "bpchar"
	TypeBytea       = "bytea"
	TypeBool        = "bool"
	TypeTimestamp   = "timestamp"
	TypeTimestamptz = "timestamptz"
	TypeDate        = "date"
	TypeTime        = "time"
	TypeTimetz      = "timetz"
	TypeInterval    = "interval"
	TypeJsonb       = "jsonb"
	TypeJson        = "json"
	TypeHstore      = "hstore"
	TypeInet        = "inet"
	TypeCidr        = "cidr"
)

func GoType(pgType string) types.Type {
	switch pgType {
	case TypeInt2, TypeInt4:
		return types.Typ[types.Int]
	case TypeInt8:
		return types.Typ[types.Int64]
	case TypeNumeric, TypeFloat4:
		return types.Typ[types.Float32]
	case TypeFloat8:
		return types.Typ[types.Float64]
	case TypeText, TypeVarchar, TypeBpchar:
		return types.Typ[types.String]
	case TypeBytea:
		return types.NewSlice(types.Typ[types.Int8])
	case TypeBool:
		return types.Typ[types.Bool]
	case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz:
		//TODO make type for it
		return nil
	case TypeInterval:
		//TODO make type for it
		return nil
	case TypeJsonb, TypeJson:
		//TODO make type for it
		return nil
	case TypeHstore:
		//TODO make type for it
		return types.NewMap(types.Typ[types.String], types.Typ[types.String])
	case TypeInet:
		//TODO make type for it
		return nil
	case TypeCidr:
		//TODO make type for it
		return nil
	}

	return nil
}

func GoSliceType(pgType string, dimensions int) types.Type {
	root := types.NewSlice(GoType(pgType))
	if dimensions < 1 {
		return root
	}

	for i := 0; i < dimensions; i++ {
		root = types.NewSlice(root)
	}

	return root
}
