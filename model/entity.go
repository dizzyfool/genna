package model

import (
	"github.com/LdDl/bungen/util"
)

// Entity stores information about table
type Entity struct {
	GoName       string
	GoNamePlural string
	PGName       string
	PGSchema     string
	PGFullName   string

	ViewName string

	Columns   []Column
	Relations []Relation

	Imports []string

	// helper indexes
	colIndex util.Index
	impIndex map[string]struct{}
}

// NewEntity creates new Entity from Postgres info
func NewEntity(schema, pgName string, columns []Column, relations []Relation) Entity {
	goName := util.EntityName(pgName)
	if schema != util.PublicSchema {
		goName = util.CamelCased(schema) + goName
	}

	goNamePlural := util.CamelCased(util.Sanitize(pgName))
	if schema != util.PublicSchema {
		goNamePlural = util.CamelCased(schema) + goNamePlural
	}

	entity := Entity{
		GoName:       goName,
		GoNamePlural: goNamePlural,
		PGName:       pgName,
		PGSchema:     schema,
		PGFullName:   util.JoinF(schema, pgName),

		Columns:   []Column{},
		Relations: []Relation{},
		colIndex:  util.NewIndex(),

		Imports:  []string{},
		impIndex: map[string]struct{}{},
	}

	if columns != nil {
		for _, col := range columns {
			entity.AddColumn(col)
		}
	}

	if relations != nil {
		for _, rel := range relations {
			entity.AddRelation(rel)
		}
	}

	return entity
}

// AddColumn adds column to entity
func (e *Entity) AddColumn(column Column) {
	if !e.colIndex.Available(column.GoName) {
		column.GoName = e.colIndex.GetNext(column.GoName)
	}
	e.colIndex.Add(column.GoName)

	e.Columns = append(e.Columns, column)

	if imp := column.Import; imp != "" {
		if _, ok := e.impIndex[imp]; !ok {
			e.impIndex[imp] = struct{}{}
			e.Imports = append(e.Imports, imp)
		}
	}
}

// AddRelation adds relation to entity
func (e *Entity) AddRelation(relation Relation) {
	if !e.colIndex.Available(relation.GoName) {
		relation.GoName = e.colIndex.GetNext(relation.GoName + util.Rel)
	}
	e.colIndex.Add(relation.GoName)

	e.Relations = append(e.Relations, relation)

	// adding relation to column
	for idx, field := range relation.FKFields {
		correspondingPK := relation.PKFields[idx]
		for i, column := range e.Columns {
			if column.PGName == field {
				e.Columns[i].AddRelation(&relation, correspondingPK)
			}
		}
	}
}

// HasMultiplePKs checks if entity has many primary keys
func (e *Entity) HasMultiplePKs() bool {
	counter := 0
	for _, col := range e.Columns {
		if col.IsPK {
			counter++
		}

		if counter > 1 {
			return true
		}
	}

	return false
}

// GetPKs returns columns which are marked as primary keys
func (e *Entity) GetPKs() []Column {
	var res []Column
	for _, col := range e.Columns {
		if col.IsPK {
			res = append(res, col)
		}
	}
	return res
}
