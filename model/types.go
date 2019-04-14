package model

import (
	"database/sql"
	"fmt"
	"go/types"
	"net"
	"time"

	"github.com/go-pg/pg"
)

const (
	// Unknown represents unknown type
	Unknown = "unknown"
	// TypeInt2 is a postgres type
	TypeInt2 = "int2"
	// TypeInt4 is a postgres type
	TypeInt4 = "int4"
	// TypeInt8 is a postgres type
	TypeInt8 = "int8"
	// TypeNumeric is a postgres type
	TypeNumeric = "numeric"
	// TypeFloat4 is a postgres type
	TypeFloat4 = "float4"
	// TypeFloat8 is a postgres type
	TypeFloat8 = "float8"
	// TypeText is a postgres type
	TypeText = "text"
	// TypeVarchar is a postgres type
	TypeVarchar = "varchar"
	// TypeUuid is a postgres type
	TypeUuid = "uuid"
	// TypeBpchar is a postgres type
	TypeBpchar = "bpchar"
	// TypeBytea is a postgres type
	TypeBytea = "bytea"
	// TypeBool is a postgres type
	TypeBool = "bool"
	// TypeTimestamp is a postgres type
	TypeTimestamp = "timestamp"
	// TypeTimestamptz is a postgres type
	TypeTimestamptz = "timestamptz"
	// TypeDate is a postgres type
	TypeDate = "date"
	// TypeTime is a postgres type
	TypeTime = "time"
	// TypeTimetz is a postgres type
	TypeTimetz = "timetz"
	// TypeInterval is a postgres type
	TypeInterval = "interval"
	// TypeJSONB is a postgres type
	TypeJSONB = "jsonb"
	// TypeJSON is a postgres type
	TypeJSON = "json"
	// TypeHstore is a postgres type
	TypeHstore = "hstore"
	// TypeInet is a postgres type
	TypeInet = "inet"
	// TypeCidr is a postgres type
	TypeCidr = "cidr"
	// TypeCidr is a postgres type
	TypePoint = "point"
)

var (
	typeMapping = map[string]bool{
		TypeInt2:        true,
		TypeInt4:        true,
		TypeInt8:        true,
		TypeNumeric:     true,
		TypeFloat4:      true,
		TypeFloat8:      true,
		TypeText:        true,
		TypeVarchar:     true,
		TypeUuid:        true,
		TypeBpchar:      true,
		TypeBytea:       true,
		TypeBool:        true,
		TypeTimestamp:   true,
		TypeTimestamptz: true,
		TypeDate:        true,
		TypeTime:        true,
		TypeTimetz:      true,
		TypeInterval:    true,
		TypeJSONB:       true,
		TypeJSON:        true,
		TypeHstore:      true,
		TypeInet:        true,
		TypeCidr:        true,
		TypePoint:       true,
	}
	basicTypes = map[string]bool{
		TypeInt2:        true,
		TypeInt4:        true,
		TypeInt8:        true,
		TypeNumeric:     true,
		TypeFloat4:      true,
		TypeFloat8:      true,
		TypeText:        true,
		TypeVarchar:     true,
		TypeUuid:        true,
		TypeBpchar:      true,
		TypeBool:        true,
		TypeTimestamp:   true,
		TypeTimestamptz: true,
		TypeDate:        true,
		TypeTime:        true,
		TypeTimetz:      true,
		TypeInterval:    true,
		TypeInet:        true,
		TypeCidr:        true,
		TypePoint:       true,
	}
)

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
			TypeInterval, TypeHstore, TypeInet, TypeCidr:
			return false
		}
	}

	return true
}

// GoImport generates import from pg type
func GoImport(pgType string, nullable, array bool, dimensions int, avoidPointers bool) string {
	// not valid types should not generate import
	if !IsValid(pgType, array) {
		return ""
	}

	typ, err := GoType(pgType, nullable, array, dimensions, avoidPointers)
	if err != nil {
		return ""
	}

	return GoImportFromType(typ)
}

// GoImportFromType generates import from go type
func GoImportFromType(typ types.Type) string {
	switch v := typ.(type) {
	case *types.Pointer:
		return GoImportFromType(v.Elem())
	case dateTime, interval, *dateTime, *interval:
		return "time"
	case netIP, netIPNet, *netIP, *netIPNet:
		return "net"
	//case jsonType, *jsonType:
	//	return "encoding/json"
	case pgNullTime, *pgNullTime:
		return "github.com/go-pg/pg"
	case sqlNullInt64, sqlNullFloat64, sqlNullBool, sqlNullString:
		return "database/sql"
	case *sqlNullInt64, *sqlNullFloat64, *sqlNullBool, *sqlNullString:
		return "database/sql"
	}

	return ""
}

// GoType generates all go types from pg type
func GoType(pgType string, nullable, array bool, dimensions int, avoidPointers bool) (types.Type, error) {
	switch {
	case array:
		return GoSliceType(pgType, dimensions)
	case nullable:
		return GoNullType(pgType, avoidPointers)
	default:
		return GoSimpleType(pgType)
	}
}

// GoType generates all go types from pg type with pointer
func GoPointerType(pgType string, array bool, dimensions int) (types.Type, error) {
	if array {
		return GoSliceType(pgType, dimensions)
	}

	switch pgType {
	case TypeJSONB, TypeJSON, TypeHstore, TypeBytea:
		return GoSimpleType(pgType)
	default:
		typ, err := GoSimpleType(pgType)
		if err != nil {
			return nil, err
		}

		return types.NewPointer(typ), nil
	}
}

// GoSimpleType generates simple go type from pg type
func GoSimpleType(pgType string) (types.Type, error) {
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
	case TypeText, TypeVarchar, TypeUuid, TypeBpchar, TypePoint:
		return types.Typ[types.String], nil
	case TypeBytea:
		return bytea{}, nil
	case TypeBool:
		return types.Typ[types.Bool], nil
	case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz:
		return dateTime{}, nil
	case TypeInterval:
		return interval{}, nil
	case TypeJSONB, TypeJSON:
		return types.NewMap(types.Typ[types.String], intrfce{}), nil
	case TypeHstore:
		return types.NewMap(types.Typ[types.String], types.Typ[types.String]), nil
	case TypeInet:
		return netIP{}, nil
	case TypeCidr:
		return netIPNet{}, nil
	}

	return nil, fmt.Errorf("type %s not supported", pgType)
}

// GoSliceType generates go slice type from pg array
func GoSliceType(pgType string, dimensions int) (types.Type, error) {
	if !IsValid(pgType, true) {
		return nil, fmt.Errorf("type %s not supported for arrays", pgType)
	}

	typ, err := GoSimpleType(pgType)
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

	return root, nil
}

// GoNullType generates go pointer type from pg nullable type
func GoNullType(pgType string, avoidPointers bool) (types.Type, error) {
	// avoiding pointers with sql.Null... types
	if avoidPointers {
		switch pgType {
		case TypeInt2, TypeInt4, TypeInt8:
			return sqlNullInt64{}, nil
		case TypeNumeric, TypeFloat4, TypeFloat8:
			return sqlNullFloat64{}, nil
		case TypeBool:
			return sqlNullBool{}, nil
		case TypeText, TypeVarchar, TypeUuid, TypeBpchar, TypePoint:
			return sqlNullString{}, nil
		}
	}

	switch pgType {
	case TypeTimestamp, TypeTimestamptz, TypeDate, TypeTime, TypeTimetz:
		return pgNullTime{}, nil
	case TypeHstore, TypeJSON, TypeJSONB:
		// hstore & jason types without pointers
		return GoSimpleType(pgType)
	default:
		// adding pointers for simple types
		typ, err := GoSimpleType(pgType)
		if err != nil {
			return nil, err
		}
		return types.NewPointer(typ), nil
	}
}

// IsBasic returns true if type is number/string/bool
func IsBasic(pgType string) bool {
	return basicTypes[pgType]
}

// Custom types goes here

type dateTime time.Time

func (dateTime) Underlying() types.Type {
	return nil
}

func (dateTime) String() string {
	return "time.Time"
}

type interval struct{}

func (interval) Underlying() types.Type {
	return types.Typ[types.Int64]
}

func (interval) String() string {
	return "time.Duration"
}

/*
type jsonType json.RawMessage

func (jsonType) Underlying() types.Type {
	return types.NewSlice(types.Typ[types.Int8])
}

func (jsonType) String() string {
	return "json.RawMessage"
}
*/

type netIP net.IP

func (netIP) Underlying() types.Type {
	return types.NewSlice(types.Typ[types.Int8])
}

func (netIP) String() string {
	return "net.IP"
}

type netIPNet net.IPNet

func (netIPNet) Underlying() types.Type {
	return nil
}

func (netIPNet) String() string {
	return "net.IPNet"
}

type pgNullTime pg.NullTime

func (pgNullTime) Underlying() types.Type {
	return nil
}

func (pgNullTime) String() string {
	return "pg.NullTime"
}

type bytea []byte

func (bytea) Underlying() types.Type {
	return types.NewSlice(types.Typ[types.Byte])
}

func (bytea) String() string {
	return "[]byte"
}

type sqlNullInt64 sql.NullInt64

func (sqlNullInt64) String() string {
	return "sql.NullInt64"
}

func (sqlNullInt64) Underlying() types.Type {
	return nil
}

type sqlNullFloat64 sql.NullFloat64

func (sqlNullFloat64) String() string {
	return "sql.NullFloat64"
}

func (sqlNullFloat64) Underlying() types.Type {
	return nil
}

type sqlNullString struct{}

func (sqlNullString) String() string {
	return "sql.NullString"
}

func (sqlNullString) Underlying() types.Type {
	return nil
}

type sqlNullBool sql.NullBool

func (sqlNullBool) String() string {
	return "sql.NullBool"
}

func (sqlNullBool) Underlying() types.Type {
	return nil
}

type intrfce struct{}

func (intrfce) String() string {
	return "interface{}"
}

func (intrfce) Underlying() types.Type {
	return &types.Interface{}
}
