package model

import (
	"path"
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
	SourceColumn string

	TargetSchema string
	TargetTable  string

	// Only for HasMany relation
	TargetColumn string
}

// StructFieldName generates field name for struct
func (r Relation) StructFieldName() string {
	return ReplaceSuffix(StructFieldName(r.SourceColumn), "ID", "")
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
	tags := NewAnnotation()

	return tags.AddTag("sql", "fk:"+r.SourceColumn).String()
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
