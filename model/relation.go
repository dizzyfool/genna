package model

import (
	"strings"

	"github.com/dizzyfool/genna/util"
)

// Relation stores relation
type Relation struct {
	FKField []string
	GoName  string

	GoType string
}

// NewRelation creates relation from pg info
func NewRelation(sourceColumns []string, targetSchema, targetTable string) Relation {
	names := make([]string, len(sourceColumns))
	for i, name := range sourceColumns {
		names[i] = util.ReplaceSuffix(util.ColumnName(name), util.ID, "")
	}

	typ := util.EntityName(targetTable)
	if targetSchema != util.PublicSchema {
		typ = util.CamelCased(targetSchema) + typ
	}

	return Relation{
		FKField: sourceColumns,
		GoName:  strings.Join(names, ""),
		GoType:  typ,
	}
}
