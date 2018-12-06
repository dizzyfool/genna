package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go/types"
	"net"
	"time"

	"github.com/go-pg/pg"
)

const (
	Unknown = "unknown"

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

var typeMapping = map[string]bool{
	TypeInt2:        true,
	TypeInt4:        true,
	TypeInt8:        true,
	TypeNumeric:     true,
	TypeFloat4:      true,
	TypeFloat8:      true,
	TypeText:        true,
	TypeVarchar:     true,
	TypeBpchar:      true,
	TypeBytea:       true,
	TypeBool:        true,
	TypeTimestamp:   true,
	TypeTimestamptz: true,
	TypeDate:        true,
	TypeTime:        true,
	TypeTimetz:      true,
	TypeInterval:    true,
	TypeJsonb:       true,
	TypeJson:        true,
	TypeHstore:      true,
	TypeInet:        true,
	TypeCidr:        true,
}

// IsValid checks type
func IsValid(pgType string, array bool) bool {
	// checking for supported types
	if _, ok := typeMapping[pgType]; !ok {
		return false
	}

	// checking for supported array types
	if array {
		switch pgType {
		case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz,
			TypeInterval, TypeJsonb, TypeJson, TypeHstore, TypeInet, TypeCidr:
			return false
		}
	}

	return true
}

// GoImport generates import from pg type
func GoImport(pgType string, nullable, avoidPointers bool) string {
	// not valid types should not generate import
	if !IsValid(pgType, false) {
		return ""
	}

	switch pgType {
	case TypeInet, TypeCidr:
		return "net"
	case TypeInterval:
		return "time"
	case TypeJsonb, TypeJson:
		return "encoding/json"
	case TypeInt2, TypeInt4, TypeInt8, TypeNumeric, TypeFloat4, TypeFloat8, TypeBool, TypeText, TypeVarchar, TypeBpchar:
		if nullable && avoidPointers {
			return "database/sql"
		}
	case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz:
		if nullable {
			return "github.com/go-pg/pg"
		}
		return "time"
	}

	return ""
}

// GoType generates go type from pg type
func GoType(pgType string) (types.Type, error) {
	if !IsValid(pgType, false) {
		return nil, fmt.Errorf("type %s not supported", pgType)
	}

	switch pgType {
	case TypeInt2, TypeInt4:
		return types.Typ[types.Int], nil
	case TypeInt8:
		return types.Typ[types.Int64], nil
	case TypeNumeric, TypeFloat4:
		return types.Typ[types.Float32], nil
	case TypeFloat8:
		return types.Typ[types.Float64], nil
	case TypeText, TypeVarchar, TypeBpchar:
		return types.Typ[types.String], nil
	case TypeBytea:
		return Bytea{}, nil
	case TypeBool:
		return types.Typ[types.Bool], nil
	case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz:
		return DateTime{}, nil
	case TypeInterval:
		return Interval{}, nil
	case TypeJsonb, TypeJson:
		return JsonType{}, nil
	case TypeHstore:
		return types.NewMap(types.Typ[types.String], types.Typ[types.String]), nil
	case TypeInet:
		return NetIp{}, nil
	case TypeCidr:
		return NetIpNet{}, nil
	}

	return nil, fmt.Errorf("type %s not supported", pgType)
}

// GoSliceType generates go slice type from pg array
func GoSliceType(pgType string, dimensions int, nullable bool) (types.Type, error) {
	if !IsValid(pgType, true) {
		return nil, fmt.Errorf("type %s not supported for arrays", pgType)
	}

	typ, err := GoType(pgType)
	if err != nil {
		return nil, err
	}

	// making multidimensional array
	root := types.NewSlice(typ)
	if dimensions < 1 {
		return root, nil
	}

	for i := 1; i < dimensions; i++ {
		root = types.NewSlice(root)
	}

	// adding pointer if nullable
	if nullable {
		return types.NewPointer(root), nil
	}

	return root, nil
}

// GoNullType generates go pointer type from pg nullable type
func GoNullType(pgType string, avoidPointers bool) (types.Type, error) {
	// avoiding pointers with sql.Null... types
	if avoidPointers {
		switch pgType {
		case TypeInt2, TypeInt4, TypeInt8:
			return SqlNullInt64{}, nil
		case TypeNumeric, TypeFloat4, TypeFloat8:
			return SqlNullFloat64{}, nil
		case TypeBool:
			return SqlNullBool{}, nil
		case TypeText, TypeVarchar, TypeBpchar:
			return SqlNullString{}, nil
		}
	}

	switch pgType {
	case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz:
		return PgNullTime{}, nil
	default:
		// adding pointers for simple types
		if typ, err := GoType(pgType); err != nil {
			return nil, err
		} else {
			return types.NewPointer(typ), nil
		}
	}
}

// Custom types goes here

type DateTime time.Time

func (DateTime) Underlying() types.Type {
	return nil
}

func (DateTime) String() string {
	return "time.Time"
}

type Interval struct{}

func (Interval) Underlying() types.Type {
	return types.Typ[types.Int64]
}

func (Interval) String() string {
	return "time.Duration"
}

type JsonType json.RawMessage

func (JsonType) Underlying() types.Type {
	return types.NewSlice(types.Typ[types.Int8])
}

func (JsonType) String() string {
	return "json.RawMessage"
}

type NetIp net.IP

func (NetIp) Underlying() types.Type {
	return types.NewSlice(types.Typ[types.Int8])
}

func (NetIp) String() string {
	return "net.IP"
}

type NetIpNet net.IPNet

func (NetIpNet) Underlying() types.Type {
	return nil
}

func (NetIpNet) String() string {
	return "net.IPNet"
}

type PgNullTime pg.NullTime

func (PgNullTime) Underlying() types.Type {
	return nil
}

func (PgNullTime) String() string {
	return "pg.NullTime"
}

type Bytea []byte

func (Bytea) Underlying() types.Type {
	return types.NewSlice(types.Typ[types.Byte])
}

func (Bytea) String() string {
	return "[]byte"
}

type SqlNullInt64 sql.NullInt64

func (SqlNullInt64) String() string {
	return "sql.NullInt64"
}

func (SqlNullInt64) Underlying() types.Type {
	return nil
}

type SqlNullFloat64 sql.NullFloat64

func (SqlNullFloat64) String() string {
	return "sql.NullFloat64"
}

func (SqlNullFloat64) Underlying() types.Type {
	return nil
}

type SqlNullString struct{}

func (SqlNullString) String() string {
	return "sql.NullString"
}

func (SqlNullString) Underlying() types.Type {
	return nil
}

type SqlNullBool sql.NullBool

func (SqlNullBool) String() string {
	return "sql.NullBool"
}

func (SqlNullBool) Underlying() types.Type {
	return nil
}
