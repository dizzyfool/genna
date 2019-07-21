package search

import (
	"fmt"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// TemplatePackage stores package info
type TemplatePackage struct {
	Package string

	HasImports bool
	Imports    []string

	Entities []TemplateEntity
}

// NewTemplatePackage creates a package for template
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

// TemplateEntity stores struct info
type TemplateEntity struct {
	model.Entity

	NoAlias bool
	Alias   string

	Columns []TemplateColumn

	Imports []string
}

// NewTemplateEntity creates an entity for template
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

		columns = append(columns, NewTemplateColumn(entity, column, options))
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

// TemplateColumn stores column info
type TemplateColumn struct {
	model.Column

	Relaxed   bool
	FieldExpr string
	TableExpr string
	Condition string
}

// NewTemplateColumn creates a column for template
func NewTemplateColumn(entity model.Entity, column model.Column, options Options) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	if options.Relaxed {
		column.GoType = model.TypeInterface
	} else {
		column.GoType = fmt.Sprintf("*%s", column.GoType)
	}

	tableExpr := fmt.Sprintf("Tables.%s.Alias", entity.GoName)
	if options.NoAlias {
		tableExpr = fmt.Sprintf("Tables.%s.Name", entity.GoName)
	}

	return TemplateColumn{
		Relaxed:   options.Relaxed,
		Column:    column,
		Condition: "?table.?field = ?value",
		TableExpr: tableExpr,
		FieldExpr: fmt.Sprintf("Columns.%s.%s", entity.GoName, column.GoName),
	}
}
