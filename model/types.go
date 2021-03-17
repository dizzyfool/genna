package model

import (
	"fmt"
)

const (
	// TypePGInt2 is a postgres type
	TypePGInt2 = "int2"
	// TypePGInt4 is a postgres type
	TypePGInt4 = "int4"
	// TypePGInt8 is a postgres type
	TypePGInt8 = "int8"
	// TypePGNumeric is a postgres type
	TypePGNumeric = "numeric"
	// TypePGFloat4 is a postgres type
	TypePGFloat4 = "float4"
	// TypePGFloat8 is a postgres type
	TypePGFloat8 = "float8"
	// TypePGText is a postgres type
	TypePGText = "text"
	// TypePGVarchar is a postgres type
	TypePGVarchar = "varchar"
	// TypePGUuid is a postgres type
	TypePGUuid = "uuid"
	// TypePGBpchar is a postgres type
	TypePGBpchar = "bpchar"
	// TypePGBytea is a postgres type
	TypePGBytea = "bytea"
	// TypePGBool is a postgres type
	TypePGBool = "bool"
	// TypePGTimestamp is a postgres type
	TypePGTimestamp = "timestamp"
	// TypePGTimestamptz is a postgres type
	TypePGTimestamptz = "timestamptz"
	// TypePGDate is a postgres type
	TypePGDate = "date"
	// TypePGTime is a postgres type
	TypePGTime = "time"
	// TypePGTimetz is a postgres type
	TypePGTimetz = "timetz"
	// TypePGInterval is a postgres type
	TypePGInterval = "interval"
	// TypePGJSONB is a postgres type
	TypePGJSONB = "jsonb"
	// TypePGJSON is a postgres type
	TypePGJSON = "json"
	// TypePGHstore is a postgres type
	TypePGHstore = "hstore"
	// TypePGInet is a postgres type
	TypePGInet = "inet"
	// TypePGCidr is a postgres type
	TypePGCidr = "cidr"
	// TypePGPoint is a postgres type
	TypePGPoint = "point"

	// TypeInt is a go type
	TypeInt = "int"
	// TypeInt32 is a go type
	TypeInt32 = "int32"
	// TypeInt64 is a go type
	TypeInt64 = "int64"
	// TypeFloat32 is a go type
	TypeFloat32 = "float32"
	// TypeFloat64 is a go type
	TypeFloat64 = "float64"
	// TypeString is a go type
	TypeString = "string"
	// TypeByteSlice is a go type
	TypeByteSlice = "[]byte"
	// TypeBool is a go type
	TypeBool = "bool"
	// TypeTime is a go type
	TypeTime = "time.Time"
	// TypeDuration is a go type
	TypeDuration = "time.Duration"
	// TypeMapInterface is a go type
	TypeMapInterface = "map[string]interface{}"
	// TypeMapString is a go type
	TypeMapString = "map[string]string"
	// TypeIP is a go type
	TypeIP = "net.IP"
	// TypeIPNet is a go type
	TypeIPNet = "net.IPNet"

	// TypeInterface is a go type
	TypeInterface = "interface{}"
)

// GoType generates simple go type from pg type
func GoType(pgType string) (string, error) {
	switch pgType {
	case TypePGInt2, TypePGInt4:
		return TypeInt, nil
	case TypePGInt8:
		return TypeInt64, nil
	case TypePGFloat4:
		return TypeFloat32, nil
	case TypePGNumeric, TypePGFloat8:
		return TypeFloat64, nil
	case TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint:
		return TypeString, nil
	case TypePGBytea:
		return TypeByteSlice, nil
	case TypePGBool:
		return TypeBool, nil
	case TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz:
		return TypeTime, nil
	case TypePGInterval:
		return TypeDuration, nil
	case TypePGJSONB, TypePGJSON:
		return TypeMapInterface, nil
	case TypePGHstore:
		return TypeMapString, nil
	case TypePGInet:
		return TypeIP, nil
	case TypePGCidr:
		return TypeIPNet, nil
	}

	return "", fmt.Errorf("unsupported type: %s", pgType)
}

// GoSlice generates go slice type from pg array
func GoSlice(pgType string, dimensions int) (string, error) {
	switch pgType {
	case TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz,
		TypePGInterval, TypePGHstore, TypePGInet, TypePGCidr:
		return "", fmt.Errorf("unsupported array type: %s", pgType)
	}

	typ, err := GoType(pgType)
	if err != nil {
		return "", err
	}

	for i := 0; i < dimensions; i++ {
		typ = fmt.Sprintf("[]%s", typ)
	}

	return typ, nil
}

// GoNullable generates all go types from pg type with pointer
func GoNullable(pgType string, useSQLNull bool) (string, error) {
	// avoiding pointers with sql.Null... types
	if useSQLNull {
		switch pgType {
		case TypePGInt2, TypePGInt4, TypePGInt8:
			return "sql.NullInt64", nil
		case TypePGNumeric, TypePGFloat4, TypePGFloat8:
			return "sql.NullFloat64", nil
		case TypePGBool:
			return "sql.NullBool", nil
		case TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint:
			return "sql.NullString", nil
		case TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz:
			return "pg.NullTime", nil
		}
	}

	typ, err := GoType(pgType)
	if err != nil {
		return "", err
	}

	switch pgType {
	case TypePGHstore, TypePGJSON, TypePGJSONB, TypePGBytea:
		// hstore & json & bytea types without pointers
		return typ, nil
	default:
		return fmt.Sprintf("*%s", typ), nil
	}
}

// GoImport generates import from go type
func GoImport(pgType string, nullable, useSQLNull bool, ver int) string {
	if nullable && useSQLNull {
		switch pgType {
		case TypePGInt2, TypePGInt4, TypePGInt8,
			TypePGNumeric, TypePGFloat4, TypePGFloat8,
			TypePGBool,
			TypePGText, TypePGVarchar, TypePGUuid, TypePGBpchar, TypePGPoint:
			return "database/sql"
		case TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz:
			if ver >= 9 {
				return fmt.Sprintf("github.com/go-pg/pg/v%d", ver)
			} else {
				return "github.com/go-pg/pg"
			}
		}
	}

	switch pgType {
	case TypePGInet, TypePGCidr:
		return "net"
	case TypePGTimestamp, TypePGTimestamptz, TypePGDate, TypePGTime, TypePGTimetz, TypePGInterval:
		return "time"
	}

	return ""
}
