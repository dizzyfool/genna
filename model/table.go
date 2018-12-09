package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Table stores information about table
type Table struct {
	Schema string
	Name   string

	// All available columns including pks and fks
	Columns []Column

	// All available relations
	Relations []Relation
}

// FileName get valid file name for model
func (t Table) FileName() string {
	name := Underscore(t.Name)
	if t.Schema != PublicSchema {
		name = Underscore(t.Schema) + "_" + name
	}
	return name
}

// PackageName get valid package name for model from schema
func (t Table) PackageName(withSchema bool, publicAlias string) string {
	if publicAlias == "" {
		publicAlias = DefaultPackage
	}

	if t.Schema == PublicSchema || !withSchema {
		return publicAlias
	}

	return PackageName(t.Schema)
}

// Model returns all imports required by model
func (t Table) Imports(withRelations bool, importPath, publicAlias string) []string {
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

	if withRelations {
		for _, relation := range t.Relations {
			if relation.TargetSchema == t.Schema {
				continue
			}

			if imp := relation.Import(importPath, publicAlias); imp != "" {
				if _, ok := index[imp]; !ok {
					imports = append(imports, imp)
					index[imp] = struct{}{}
				}
			}
		}
	}

	return imports
}

// Model returns model name in camel case and in singular form
func (t Table) ModelName(withSchema bool) string {
	name := ModelName(t.Name)
	if withSchema && t.Schema != PublicSchema {
		name = CamelCased(t.Schema) + name
	}

	return name
}

// TableName returns valid table name with schema and quoted if needed
func (t Table) TableName() string {
	table := t.Name
	if HasUpper(table) {
		table = fmt.Sprintf(`\"%s\"`, table)
	}

	if t.Schema == PublicSchema {
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
	if t.Schema == PublicSchema {
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

func (t Table) HasMultiplePKs() bool {
	count := 0
	for _, column := range t.Columns {
		if column.IsPK {
			count++
			if count >= 2 {
				return true
			}
		}
	}

	return false
}

func (t Table) Validate() error {
	if strings.Trim(t.Schema, " ") == "" {
		return fmt.Errorf("shema name is empty")
	}

	if strings.Trim(t.Name, " ") == "" {
		return fmt.Errorf("table name is empty")
	}

	rgxp := regexp.MustCompile(`[^\w\d_]+`)
	if rgxp.Match([]byte(t.Schema)) {
		return fmt.Errorf("shema name '%s' contains illegal character(s)", t.Schema)
	}

	if rgxp.Match([]byte(t.Name)) {
		return fmt.Errorf("table name '%s' contains illegal character(s)", t.Name)
	}

	if len(t.Columns) == 0 {
		return fmt.Errorf("table has no columns")
	}

	for _, column := range t.Columns {
		if err := column.Validate(); err != nil {
			return errors.Wrap(err, "column '%s' is not valid")
		}
		if column.IsFK && len(t.Relations) == 0 {
			return fmt.Errorf("table has fkey(s) but no relations")
		}
	}

	return nil
}
