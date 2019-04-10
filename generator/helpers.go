package generator

import (
	"fmt"
	"github.com/dizzyfool/genna/model"
	"html/template"
)

// Stores package info
type templatePackage struct {
	Package     string
	HasImports  bool
	Imports     []string

	Models []templateTable
}

// newMultiPackage creates a package with multiple models
func newTemplatePackage(tables []model.Table, options Options) templatePackage {
	imports := make([]string, 0)

	models := make([]templateTable, len(tables))
	for i, table := range tables {
		imports = append(imports, table.Imports()...)

		models[i] = newTemplateTable(table, options)
		models[i].uniqualizeFields()
	}

	imports = model.UniqStrings(imports)

	return templatePackage{
		Package:     options.Package,
		HasImports:  len(imports) > 0,
		Imports:     imports,

		Models: models,
	}
}

// stores struct info
type templateTable struct {
	StructName string
	StructTag  template.HTML

	JoinAlias string
	TableName string

	Columns []templateColumn

	HasRelations bool
	Relations    []templateRelation
}

func newTemplateTable(table model.Table, options Options) templateTable {
	if table.HasMultiplePKs() {
		options.KeepPK = true
	}

	columns := make([]templateColumn, len(table.Columns))
	for i, column := range table.Columns {
		columns[i] = newTemplateColumn(column, options)
	}

	relations := make([]templateRelation, len(table.Relations))
	for i, relation := range table.Relations {
		relations[i] = newTemplateRelation(relation, options)
	}

	return templateTable{
		StructName: table.ModelName(),
		StructTag:  template.HTML(fmt.Sprintf("`%s`", table.TableNameTag(options.View, options.NoAlias, options.NoDiscard))),

		JoinAlias: table.JoinAlias(),
		TableName: table.Name,

		Columns: columns,

		HasRelations: len(relations) > 0,
		Relations:    relations,
	}
}

// stores column info
type templateColumn struct {
	FieldName    string
	FieldDBName  string
	FieldType    string
	FieldTag     template.HTML
	FieldComment template.HTML
}

func newTemplateColumn(column model.Column, options Options) templateColumn {
	return templateColumn{
		FieldName:    column.StructFieldName(options.KeepPK),
		FieldDBName:  column.Name,
		FieldType:    column.StructFieldType(),
		FieldTag:     template.HTML(fmt.Sprintf("`%s`", column.StructFieldTag())),
		FieldComment: template.HTML(column.Comment()),
	}
}

// stores relation info
type templateRelation struct {
	FieldName    string
	FieldType    string
	FieldTag     template.HTML
	FieldComment template.HTML
}

func newTemplateRelation(relation model.Relation, options Options) templateRelation {
	return templateRelation{
		FieldName:    relation.StructFieldName(),
		FieldType:    relation.StructFieldType(),
		FieldTag:     template.HTML(fmt.Sprintf("`%s`", relation.StructFieldTag())),
		FieldComment: template.HTML(relation.Comment()),
	}
}

func (t templateTable) uniqualizeFields() {
	index := map[string]bool{}

	for i, column := range t.Columns {
		fieldName := column.FieldName

		if _, ok := index[fieldName]; !ok {
			index[fieldName] = true
			continue
		}

		suffix := 1
	couter:
		for {
			fieldName = fmt.Sprintf("%s%d", column.FieldName, suffix)

			for _, col := range t.Columns {
				if col.FieldName == fieldName {
					suffix++
					continue couter
				}
			}
			t.Columns[i].FieldName = fieldName
			break
		}
	}

	for i, relation := range t.Relations {
		fieldName := relation.FieldName

		if _, ok := index[fieldName]; !ok {
			index[fieldName] = true
			continue
		}

		suffix := 0

	router:
		for {
			if suffix == 0 {
				fieldName = fmt.Sprintf("%sRel", relation.FieldName)
			} else {
				fieldName = fmt.Sprintf("%sRel%d", relation.FieldName, suffix)
			}

			for _, col := range t.Columns {
				if col.FieldName == fieldName {
					suffix++
					continue router
				}
			}

			for _, rel := range t.Relations {
				if rel.FieldName == fieldName {
					suffix++
					continue router
				}
			}

			t.Relations[i].FieldName = fieldName
			break
		}
	}
}
