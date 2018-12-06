package model

const (
	HasOne int = iota
	HasMany
	ManyToMany
)

// Relation stores relation of a table with target table
type Relation struct {
	Type int

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
func (r Relation) StructFieldType() string {
	return "*" + ModelName(r.TargetTable)
}

// StructFieldTag generates field tag for struct
func (r Relation) StructFieldTag() string {
	tags := NewAnnotation()

	return tags.AddTag("sql", "fk:"+r.SourceColumn).String()
}
