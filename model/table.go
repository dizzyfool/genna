package model

import "fmt"

// Table stores information about table
type Table struct {
	Schema string
	Name   string

	// All available columns including pks and fks
	Columns []Column

	// All available relations
	Relations []Relation
}

// Model returns all imports required by model
func (t Table) Imports() []string {
	imports := make([]string, 0)
	index := make(map[string]struct{})

	for _, column := range t.Columns {
		if imp := column.Import(); imp != "" {
			if _, ok := index[imp]; !ok {
				imports = append(imports, imp)
				index[imp] = struct{}{}
			}
		}
	}

	return imports
}

// Model returns model name in camel case and in singular form
func (t Table) ModelName() string {
	return ModelName(t.Name)
}

// TableName returns valid table name with schema and quoted if needed
func (t Table) TableName() string {
	table := t.Name
	if HasUpper(table) {
		table = fmt.Sprintf(`\"%s\"`, table)
	}

	if t.Schema == "public" {
		return table
	}

	schema := t.Schema
	if HasUpper(schema) {
		schema = fmt.Sprintf(`\"%s\"`, schema)
	}

	return fmt.Sprintf("%s.%s", schema, table)
}

// ViewName returns view name for table starting with "get"
func (t Table) ViewName() string {
	if t.Schema == "public" {
		return fmt.Sprintf(`\"get%s\"`, CamelCased(t.Name))
	}

	schema := t.Schema
	if HasUpper(schema) {
		schema = fmt.Sprintf(`\"%s\"`, schema)
	}

	return fmt.Sprintf(`%s.\"get%s\"`, schema, CamelCased(t.Name))
}

// TableNameTag returns tag for tableName property
func (t Table) TableNameTag(noDiscard, withView bool) string {
	annotation := NewAnnotation()

	annotation.AddTag("sql", t.TableName())
	if withView {
		annotation.AddTag("sql", fmt.Sprintf("select:%s", t.ViewName()))
	}

	if !noDiscard {
		// leading comma is required
		annotation.AddTag("pg", ",discard_unknown_columns")
	}

	return annotation.String()
}
