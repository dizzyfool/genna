package validate

import (
	"fmt"
	"strings"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

const (
	// Nil is nil check types
	Nil = "nil"
	// Zero is 0 check types
	Zero = "zero"
	// PZero is 0 check types for pointers
	PZero = "pzero"
	// Len is length check types
	Len = "len"
	// PLen is length check types for pointers
	PLen = "plen"
	// Enum is allowed values check types
	Enum = "enum"
	// PEnum is allowed values check types for pointers
	PEnum = "penum"
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
	for _, entity := range entities {
		mdl := NewTemplateEntity(entity, options)
		if len(mdl.Columns) == 0 {
			continue
		}

		for _, imp := range mdl.Imports {
			imports.Add(imp)
		}

		models = append(models, mdl)
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

	Columns []TemplateColumn
	Imports []string
}

// NewTemplateEntity creates an entity for template
func NewTemplateEntity(entity model.Entity, options Options) TemplateEntity {
	imports := util.NewSet()

	var columns []TemplateColumn
	for _, column := range entity.Columns {
		if !isValidatable(column) {
			continue
		}

		tmpl := NewTemplateColumn(column, options)

		columns = append(columns, tmpl)
		if tmpl.Import != "" {
			imports.Add(tmpl.Import)
		}
	}

	return TemplateEntity{
		Entity: entity,

		Columns: columns,
		Imports: imports.Elements(),
	}
}

// TemplateColumn stores column info
type TemplateColumn struct {
	model.Column

	Check string
	Enum  string

	Import string
}

// NewTemplateColumn creates a column for template
func NewTemplateColumn(column model.Column, options Options) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	tmpl := TemplateColumn{
		Column: column,

		Check: check(column),
	}

	if len(column.Values) > 0 {
		tmpl.Enum = fmt.Sprintf(`"%s"`, strings.Join(column.Values, `", "`))
	}

	if tmpl.Check == PLen || tmpl.Check == Len {
		tmpl.Import = "unicode/utf8"
	}

	return tmpl
}

// isValidatable checks if field can be validated
func isValidatable(c model.Column) bool {
	// validate FK
	if c.IsFK {
		return true
	}

	// validate complex types
	if !c.Nullable && (c.IsArray || c.GoType == model.TypeMapInterface || c.GoType == model.TypeMapString) {
		return true
	}

	// validate strings len
	if c.GoType == model.TypeString && c.MaxLen > 0 {
		return true
	}

	// validate enum
	if len(c.Values) > 0 {
		return true
	}

	return false
}

// check return check type for validation
func check(c model.Column) string {
	if !isValidatable(c) {
		return ""
	}

	if c.IsArray || c.GoType == model.TypeMapInterface || c.GoType == model.TypeMapString {
		return Nil
	}

	if c.IsFK {
		if c.Nullable {
			return PZero
		}
		return Zero
	}

	if c.GoType == model.TypeString && c.MaxLen > 0 {
		if c.Nullable {
			return PLen
		}
		return Len
	}

	if len(c.Values) > 0 {
		if c.Nullable {
			return PEnum
		}
		return Enum
	}

	return ""
}
