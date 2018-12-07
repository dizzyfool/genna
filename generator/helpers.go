package generator

import (
	"fmt"
	"html/template"
	"path"

	"github.com/dizzyfool/genna/model"
)

// Stores package info
type templatePackage struct {
	FileName string

	Package    string
	HasImports bool
	Imports    []string

	Models []templateTable
}

// newMultiPackage creates a package with multiple models
func newMultiPackage(packageName string, tables []model.Table, options Options) *templatePackage {
	imports := make([]string, 0)

	models := make([]templateTable, len(tables))
	for i, table := range tables {
		imports = append(imports, table.Imports(options.SchemaPackage, options.ImportPath, options.Package)...)

		models[i] = newTemplateTable(table, options)
	}

	imports = model.UniqStrings(imports)

	return &templatePackage{
		FileName: path.Join(
			options.Output,
			tables[0].PackageName(options.SchemaPackage, options.Package),
			packageName+".go",
		),

		Package:    packageName,
		HasImports: len(imports) > 0,
		Imports:    imports,

		Models: models,
	}
}

// newSinglePackage creates a package with simple model
func newSinglePackage(table model.Table, options Options) *templatePackage {
	imports := table.Imports(options.SchemaPackage, options.ImportPath, options.Package)

	return &templatePackage{
		FileName: path.Join(
			options.Output,
			table.PackageName(options.SchemaPackage, options.Package),
			table.FileName()+".go",
		),

		Package:    table.PackageName(options.SchemaPackage, options.Package),
		HasImports: len(imports) > 0,
		Imports:    imports,

		Models: []templateTable{newTemplateTable(table, options)},
	}
}

// stores struct info
type templateTable struct {
	StructName string
	StructTag  template.HTML

	Columns []templateColumn

	HasRelations bool
	Relations    []templateRelation
}

func newTemplateTable(table model.Table, options Options) templateTable {
	columns := make([]templateColumn, len(table.Columns))
	for i, column := range table.Columns {
		columns[i] = newTemplateColumn(column, options)
	}

	relations := make([]templateRelation, len(table.Relations))
	for i, relation := range table.Relations {
		relations[i] = newTemplateRelation(relation, options)
	}

	return templateTable{
		StructName: table.ModelName(!options.SchemaPackage),
		StructTag:  template.HTML(fmt.Sprintf("`%s`", table.TableNameTag(options.NoDiscard, options.View))),

		Columns: columns,

		HasRelations: len(relations) > 0,
		Relations:    relations,
	}
}

// stores column info
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

// stores relation info
type templateRelation struct {
	FieldName string
	FieldType string
	FieldTag  template.HTML
}

func newTemplateRelation(relation model.Relation, options Options) templateRelation {
	return templateRelation{
		FieldName: relation.StructFieldName(),
		FieldType: relation.StructFieldType(!options.SchemaPackage, options.Package),
		FieldTag:  template.HTML(fmt.Sprintf("`%s`", relation.StructFieldTag())),
	}
}
