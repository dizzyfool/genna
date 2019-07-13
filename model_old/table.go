package model_old

import (
	"fmt"
	"github.com/dizzyfool/genna/util"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	// SearchSuffix added to search filters struct
	SearchSuffix = "Search"

	defaultAlias = "t"
)

// Entity stores information about table
type Entity struct {
	Schema string
	Name   string

	// All available columns including pks and fks
	Columns []Column

	// All available relations
	Relations []Relation
}

// Imports returns all imports required by model
func (t Entity) Imports() []string {
	imports := make([]string, 0)
	index := make(map[string]struct{})

	for _, column := range t.Columns {
		if imp := column.Import(); imp != "" {
			if _, ok := index[imp]; !ok {
				imports = append(imports, imp)
				index[imp] = struct{}{}
			}
		}

		if validate := column.ValidationCheck(); validate == ValidateLen || validate == ValidatePLen {
			imp := "unicode/utf8"
			if _, ok := index[imp]; !ok {
				imports = append(imports, imp)
				index[imp] = struct{}{}
			}
		}
	}

	return imports
}

// EntityName returns model name in camel case and in singular form
func (t Entity) ModelName() string {
	name := util.EntityName(t.Name)
	if t.Schema != PublicSchema {
		name = util.CamelCased(t.Schema) + name
	}

	return name
}

// TableName returns valid table name with schema and quoted if needed
func (t Entity) TableName(quoted bool) string {
	table := t.Name
	if util.HasUpper(table) && quoted {
		table = fmt.Sprintf(`\"%s\"`, table)
	}

	if t.Schema == PublicSchema {
		return table
	}

	schema := t.Schema
	if util.HasUpper(schema) && quoted {
		schema = fmt.Sprintf(`\"%s\"`, schema)
	}

	return fmt.Sprintf("%s.%s", schema, table)
}

// ViewName returns view name for table starting with "get"
func (t Entity) ViewName() string {
	if t.Schema == PublicSchema {
		return fmt.Sprintf(`\"get%s\"`, util.CamelCased(t.Name))
	}

	schema := t.Schema
	if util.HasUpper(schema) {
		schema = fmt.Sprintf(`\"%s\"`, schema)
	}

	return fmt.Sprintf(`%s.\"get%s\"`, schema, util.CamelCased(t.Name))
}

// Alias generates alias name for table
func (t Entity) Alias() string {
	return defaultAlias
}

// TableNameTag returns tag for tableName property
func (t Entity) TableNameTag(withView, noDiscard, noAlias bool) string {
	annotation := util.NewAnnotation()

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
func (t Entity) HasMultiplePKs() bool {
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
func (t Entity) JoinAlias() string {
	return util.Underscore(t.ModelName())
}

// SearchModelName returns model name for search filters
func (t Entity) SearchModelName() string {
	return fmt.Sprintf("%s%s", t.ModelName(), SearchSuffix)
}

// SearchImports returns all imports required by search filters
func (t Entity) SearchImports() []string {
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

// HasValidation return true if any column of table can be validated
func (t Entity) HasValidation() bool {
	for _, column := range t.Columns {
		if column.IsValidatable() {
			return true
		}
	}

	return false
}

// Validate checks current table for problems
func (t Entity) Validate() error {
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
