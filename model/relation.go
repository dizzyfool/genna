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
