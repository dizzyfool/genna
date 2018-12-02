package model

// Table stores information about table
type Table struct {
	Schema string
	Name   string

	// All available columns including pks and fks
	Columns []Column

	// All available relations
	Relations []Relation
}
