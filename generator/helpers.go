package generator

import (
	"fmt"
	"html/template"

	"github.com/dizzyfool/genna/model"
)

type templateTable struct {
	Package    string
	HasImports bool
	Imports    []string

	StructName string
	StructTag  template.HTML

	Columns []templateColumn

	HasRelations bool
	Relations    []templateRelation
}

func newTemplateTable(table model.Table, options Options) *templateTable {
	imports := table.Imports()

	columns := make([]templateColumn, len(table.Columns))
	for i, column := range table.Columns {
		columns[i] = newTemplateColumn(column, options)
	}

	relations := make([]templateRelation, len(table.Relations))
	for i, relation := range table.Relations {
		relations[i] = newTemplateRelation(relation, options)
	}

	return &templateTable{
		Package:    table.PackageName(options.SchemaAsPackage, options.Package),
		HasImports: len(imports) > 0,
		Imports:    imports,

		StructName: table.ModelName(!options.SchemaAsPackage),
		StructTag:  template.HTML(fmt.Sprintf("`%s`", table.TableNameTag(options.NoDiscard, options.View))),

		Columns: columns,

		HasRelations: len(relations) > 0,
		Relations:    relations,
	}
}

type templateColumn struct {
	FieldName string
	FieldType string
	FieldTag  template.HTML
}

func newTemplateColumn(column model.Column, options Options) templateColumn {
	return templateColumn{
		FieldName: column.StructFieldName(options.KeepPK),
		FieldType: column.StructFieldType(),
		FieldTag:  template.HTML(fmt.Sprintf("`%v`", column.StructFieldTag())),
	}
}

type templateRelation struct {
	FieldName string
	FieldType string
	FieldTag  template.HTML
}

func newTemplateRelation(relation model.Relation, options Options) templateRelation {
	return templateRelation{
		FieldName: relation.StructFieldName(),
		FieldType: relation.StructFieldType(!options.SchemaAsPackage),
		FieldTag:  template.HTML(fmt.Sprintf("`%s`", relation.StructFieldTag())),
	}
}
