package model

import (
	"strings"
)

const (
	// HasOne is has one relation
	HasOne int = iota
	// HasMany is has many relation
	HasMany
	// ManyToMany is many to many relation
	ManyToMany
)

// Relation stores relation of a table with target table
type Relation struct {
	Type int

	SourceSchema string
	SourceTable  string

	// Only for HasOne relation
	SourceColumns []string

	TargetSchema string
	TargetTable  string

	// Only for HasMany relation
	TargetColumns []string
}

// StructFieldName generates field name for struct
func (r Relation) StructFieldName() string {
	names := make([]string, len(r.SourceColumns))
	for i, name := range r.SourceColumns {
		names[i] = ReplaceSuffix(StructFieldName(name), "ID", "")
	}

	return strings.Join(names, "")
}

// StructFieldType generates field type for struct
func (r Relation) StructFieldType() string {
	name := ModelName(r.TargetTable)
	if r.TargetSchema != PublicSchema {
		name = CamelCased(r.TargetSchema) + name
	}

	return "*" + name
}

// StructFieldTag generates field tag for struct
func (r Relation) StructFieldTag() string {
	tags := NewAnnotation().AddTag("pg", "fk:"+strings.Join(r.SourceColumns, ","))
	if len(r.SourceColumns) > 1 {
		tags.AddTag("sql", "-")
	}

	return tags.String()
}

// Comment generates commentary for relation
func (r Relation) Comment() string {
	if len(r.SourceColumns) > 1 {
		return "// multiple fields relations not supported by go-pg"
	}
	return ""
}
