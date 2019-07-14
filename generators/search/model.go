package search

import (
	"fmt"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// Stores package info
type TemplatePackage struct {
	Package string

	HasImports bool
	Imports    []string

	Entities []TemplateEntity
}

// newMultiPackage creates a package with multiple models
func NewTemplatePackage(entities []model.Entity, options Options) TemplatePackage {

	imports := util.NewSet()

	var models []TemplateEntity
	for i, entity := range entities {
		mdl := NewTemplateEntity(entity, options)
		if len(mdl.Columns) == 0 {
			continue
		}

		models = append(models, mdl)
		for _, imp := range models[i].Imports {
			imports.Add(imp)
		}
	}

	return TemplatePackage{
		Package: options.Package,

		HasImports: imports.Len() > 0,
		Imports:    imports.Elements(),

		Entities: models,
	}
}

// stores struct info
type TemplateEntity struct {
	model.Entity

	NoAlias bool
	Alias   string

	Columns []TemplateColumn

	Imports []string
}

func NewTemplateEntity(entity model.Entity, options Options) TemplateEntity {
	if entity.HasMultiplePKs() {
		options.KeepPK = true
	}

	imports := util.NewSet()

	var columns []TemplateColumn
	for _, column := range entity.Columns {
		if column.IsArray || column.GoType == model.TypeMapInterface || column.GoType == model.TypeMapString {
			continue
		}

		columns = append(columns, NewTemplateColumn(column, options))
		if column.Import != "" {
			imports.Add(column.Import)
		}
	}

	return TemplateEntity{
		Entity: entity,

		NoAlias: options.NoAlias,
		Alias:   util.DefaultAlias,

		Columns: columns,
		Imports: imports.Elements(),
	}
}

// stores column info
type TemplateColumn struct {
	model.Column
}

func NewTemplateColumn(column model.Column, options Options) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	if options.Relaxed {
		column.GoType = model.TypeInterface
	} else {
		column.GoType = fmt.Sprintf("*%s", column.GoType)
	}

	return TemplateColumn{
		Column: column,
	}
}
