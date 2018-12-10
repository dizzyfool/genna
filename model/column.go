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
func (c Column) StructFieldTag() string {
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

	return tags.String()
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
