package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// SearchSuffix added to search filters struct
const SearchSuffix = "Search"

// Table stores information about table
type Table struct {
	Schema string
	Name   string

	// All available columns including pks and fks
	Columns []Column

	// All available relations
	Relations []Relation
}

// Imports returns all imports required by model
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

// ModelName returns model name in camel case and in singular form
func (t Table) ModelName() string {
	name := ModelName(t.Name)
	if t.Schema != PublicSchema {
		name = CamelCased(t.Schema) + name
	}

	return name
}

// TableName returns valid table name with schema and quoted if needed
func (t Table) TableName(quoted bool) string {
	table := t.Name
	if HasUpper(table) && quoted {
		table = fmt.Sprintf(`\"%s\"`, table)
	}

	if t.Schema == PublicSchema {
		return table
	}

	schema := t.Schema
	if HasUpper(schema) && quoted {
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

// Alias generates alias name for table
func (t Table) Alias() string {
	// default alias is "t" for filters
	return "t"

	//alias := strings.ToLower(t.Name)
	//if t.Schema != PublicSchema {
	//	alias = fmt.Sprintf(`%s_%s`, strings.ToLower(t.Schema), alias)
	//}
	//
	//return alias
}

// TableNameTag returns tag for tableName property
func (t Table) TableNameTag(withView, noDiscard, noAlias bool) string {
	annotation := NewAnnotation()

	annotation.AddTag("sql", t.TableName(true))
	if withView {
		annotation.AddTag("sql", fmt.Sprintf("select:%s", t.ViewName()))
	}

	if !noAlias {
		annotation.AddTag("sql", fmt.Sprintf("alias:%s", t.Alias()))
	}

	if !noDiscard {
		// leading comma is required
		annotation.AddTag("pg", ",discard_unknown_columns")
	}

	return annotation.String()
}

// HasMultiplePKs returns true if model have several PKs
// can disable converting PK's name to ID
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

// JoinAlias used in raws relations
func (t Table) JoinAlias() string {
	return Underscore(t.ModelName())
}

// SearchModelName returns model name for search filters
func (t Table) SearchModelName() string {
	return fmt.Sprintf("%s%s", t.ModelName(), SearchSuffix)
}

// SearchImports returns all imports required by search filters
func (t Table) SearchImports() []string {
	imports := make([]string, 0)
	index := make(map[string]struct{})

	for _, column := range t.Columns {
		if imp := column.SearchImport(); imp != "" {
			if _, ok := index[imp]; !ok {
				imports = append(imports, imp)
				index[imp] = struct{}{}
			}
		}
	}

	return imports
}

// Validate checks current table for problems
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
