package model

import (
	"fmt"
	"regexp"
	"strings"
)

// Column stores information about column
// it does not store relation info
type Column struct {
	Name       string
	Type       string
	IsArray    bool
	Dimensions int
	IsNullable bool
	IsPK       bool
	IsFK       bool
	MaxLen     int
	Enum       []string
}

// Import gets import for column
func (c Column) Import() string {
	return GoImport(c.Type, c.IsNullable, c.IsArray, c.Dimensions, false)
}

// StructFieldName generates field name for struct
func (c Column) StructFieldName(keepPK bool) string {
	if c.IsPK && !keepPK {
		return "ID"
	}
	return StructFieldName(c.Name)
}

// StructFieldType generates field type for struct
func (c Column) StructFieldType() string {
	typ, err := GoType(c.Type, c.IsNullable, c.IsArray, c.Dimensions, false)
	if err != nil {
		return "interface{}"
	}

	return typ.String()
}

// StructFieldTag generates field tag for struct
func (c Column) StructFieldTag(softDeleteColumn string) string {
	// Ignoring unknown types
	if c.StructFieldType() == "interface{}" {
		return `sql:"-"`
	}

	tags := NewAnnotation()
	tags.AddTag("sql", c.Name)

	if c.IsPK {
		tags.AddTag("sql", "pk")
	}

	if c.Type == TypeHstore {
		tags.AddTag("sql", "hstore")
	} else if c.IsArray {
		tags.AddTag("sql", "array")
	}

	if !c.IsNullable && !c.IsPK {
		tags.AddTag("sql", "notnull")
	}

	if c.isSoftDeletable(softDeleteColumn) {
		tags.AddTag("pg", ",soft_delete")
	}

	return tags.String()
}

func (c Column) isSoftDeletable(softDeleteColumn string) bool {
	return softDeleteColumn != "" &&
		softDeleteColumn == c.Name &&
		c.IsNullable &&
		!c.IsPK &&
		!c.IsFK &&
		!c.IsArray &&
		IsDateTimeType(c.Type)
}

// IsSearchable checks if field can be added to search filters
func (c Column) IsSearchable() bool {
	return !c.IsArray && IsBasic(c.Type)
}

// SearchFieldType generates field type for search filters struct
func (c Column) SearchFieldType(strict bool) string {
	if !strict {
		return "interface{}"
	}

	typ, err := GoPointerType(c.Type, c.IsArray, c.Dimensions)
	if err != nil {
		return "interface{}"
	}

	return typ.String()
}

// SearchImport gets import for search column
func (c Column) SearchImport() string {
	if !c.IsSearchable() {
		return ""
	}

	typ, err := GoPointerType(c.Type, c.IsArray, c.Dimensions)
	if err != nil {
		return ""
	}

	return GoImportFromType(typ)
}

// Comment generates commentary for column
func (c Column) Comment() string {
	if _, err := GoType(c.Type, c.IsNullable, c.IsArray, c.Dimensions, false); err != nil {
		return fmt.Sprintf("// type %s not supported", c.Type)
	}

	return ""
}

// IsValidatable checks if field can be validated
func (c Column) IsValidatable() bool {
	// validate FK
	if c.IsFK {
		return true
	}

	// validate strings len
	if IsStringType(c.Type) && c.MaxLen > 0 {
		return true
	}

	// validate complex types
	if !c.IsNullable && (c.IsArray || IsComplexType(c.Type)) {
		return true
	}

	// validate enum
	if len(c.Enum) > 0 {
		return true
	}

	return false
}

// ValidationCheck return check type for validation
func (c Column) ValidationCheck() string {
	if c.IsValidatable() {
		if c.IsArray || IsComplexType(c.Type) {
			return "nil"
		}

		if c.IsFK {
			if c.IsNullable {
				return "pzero"
			}
			return "zero"
		}

		if IsStringType(c.Type) && c.MaxLen > 0 {
			if c.IsNullable {
				return "plen"
			}
			return "len"
		}

		if len(c.Enum) > 0 {
			if c.IsNullable {
				return "penum"
			}
			return "enum"
		}
	}

	return ""
}

// Validate checks current column for problems
func (c Column) Validate() error {
	if strings.Trim(c.Name, " ") == "" {
		return fmt.Errorf("column name is empty")
	}

	if regexp.MustCompile(`[^\w\d_]+`).MatchString(c.Name) {
		return fmt.Errorf("column name '%s' contains illegal character(s)", c.Name)
	}

	if c.IsPK && c.IsNullable {
		return fmt.Errorf("column can not be pkey and be nullable")
	}

	if c.IsArray {
		if c.Type == TypeHstore {
			return fmt.Errorf("array of hstores is not supported")
		}

		if c.Dimensions <= 0 {
			return fmt.Errorf("array column has %d dimesions", c.Dimensions)
		}
	}

	if !IsValid(c.Type, c.IsArray) {
		return fmt.Errorf("unsupported type '%s' (array = %t)", c.Type, c.IsArray)
	}

	return nil
}
