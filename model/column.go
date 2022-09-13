package model

import (
	"github.com/LdDl/bungen/util"
)

type columnRelWrap struct {
	*Relation
	RelationPK string
}

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
	// Relation *Relation
	Relation *columnRelWrap

	Import string

	MaxLen int
	Values []string
}

// NewColumn creates Column from Postgres info
func NewColumn(pgName string, pgType string, nullable, sqlNulls, array bool, dims int, pk, fk bool, len int, values []string, customTypes CustomTypeMapping) Column {
	var (
		err error
		ok  bool
	)

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
		GoName:     util.ColumnName(pgName),
	}

	if customTypes == nil {
		customTypes = CustomTypeMapping{}
	}

	if column.GoType, ok = customTypes.GoType(pgType); !ok || column.GoType == "" {
		if column.GoType, err = GoType(pgType); err != nil {
			column.GoType = "interface{}"
		}
	}

	switch {
	case column.IsArray:
		column.Type, err = GoSlice(pgType, dims)
	case column.Nullable:
		column.Type, err = GoNullable(pgType, sqlNulls, customTypes)
	default:
		column.Type = column.GoType
	}

	if err != nil {
		column.Type = column.GoType
	}

	if column.Import, ok = customTypes.GoImport(pgType); !ok {
		column.Import = GoImport(pgType, nullable, sqlNulls)
	}

	return column
}

// AddRelation adds relation to column. Should be used if FK
func (c *Column) AddRelation(relation *Relation, relPK string) {
	c.Relation = &columnRelWrap{
		Relation:   relation,
		RelationPK: relPK,
	}
}
