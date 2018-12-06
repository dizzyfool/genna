package model

import "go/types"

// Columns stores information about column
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
	return GoImport(c.Type, c.IsNullable, false)
}

// StructFieldName generates field name for struct
func (c Column) StructFieldName() string {
	if c.IsPK {
		return "ID"
	}
	return StructFieldName(c.Name)
}

// StructFieldType generates field type for struct
func (c Column) StructFieldType() string {
	var (
		typ types.Type
		err error
	)

	switch {
	case c.IsArray:
		typ, err = GoSliceType(c.Type, c.Dimensions, c.IsNullable)
	case c.IsNullable:
		typ, err = GoNullType(c.Type, false)
	default:
		typ, err = GoType(c.Type)
	}

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
		return tags.AddTag("sql", "pk").String()
	}

	if c.Type == TypeHstore {
		tags.AddTag("sql", "hstore")
	} else if c.IsArray {
		tags.AddTag("sql", "array")
	}

	if !c.IsNullable {
		tags.AddTag("sql", "notnull")
	}

	return tags.String()
}
