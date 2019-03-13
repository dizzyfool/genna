package model

import (
	"path"
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
	return StructFieldName(Singular(r.TargetTable))
}

// StructFieldType generates field type for struct
// withSchema adds schema to filed name
// publicAlias rewrites public schema name
func (r Relation) StructFieldType(withSchema bool, publicAlias string) string {
	if publicAlias == "" {
		publicAlias = DefaultPackage
	}

	name := ModelName(r.TargetTable)
	if withSchema && r.TargetSchema != PublicSchema {
		name = CamelCased(r.TargetSchema) + name
	}

	if !withSchema && r.TargetSchema != r.SourceSchema {
		if r.TargetSchema == PublicSchema {
			name = publicAlias + "." + name
		} else {
			name = PackageName(r.TargetSchema) + "." + name
		}
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

// Import gets import for relation
// importPath adds prefix to import path
func (r Relation) Import(importPath, publicAlias string) string {
	if publicAlias == "" {
		publicAlias = DefaultPackage
	}

	if r.TargetSchema == r.SourceSchema {
		return ""
	}

	if r.TargetSchema == PublicSchema {
		return path.Join(importPath, publicAlias)
	}

	return path.Join(importPath, PackageName(r.TargetSchema))
}
