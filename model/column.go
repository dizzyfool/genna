package model

import (
	"github.com/dizzyfool/genna/util"
)

// Column stores information about column
type Column struct {
	GoName string
	PGName string

	Type string

	GoType string
	PGType string

	Nullable bool

	IsArray    bool
	Dimensions int

	IsPK bool
	IsFK bool

	Import string

	MaxLen int
	Values []string
}

// NewColumn creates Column from pg info
func NewColumn(pgName string, pgType string, nullable, sqlNulls, array bool, dims int, pk, fk bool, len int, values []string) Column {
	var err error

	array, dims = fixIsArray(pgType, array, dims)

	column := Column{
		PGName:     pgName,
		PGType:     pgType,
		Nullable:   nullable,
		IsArray:    array,
		Dimensions: dims,
		IsPK:       pk,
		IsFK:       fk,
		MaxLen:     len,
		Values:     values,
	}

	column.GoName = util.ColumnName(pgName)

	column.GoType, err = goType(pgType)
	if err != nil {
		column.GoType = "interface{}"
	}

	switch {
	case column.IsArray:
		column.Type, err = goSlice(pgType, dims)
	case column.Nullable:
		column.Type, err = goNullable(pgType, sqlNulls)
	default:
		column.Type = column.GoType
	}

	if err != nil {
		column.Type = column.GoType
	}

	column.Import = goImport(pgType, nullable, sqlNulls)

	return column
}
