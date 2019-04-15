package generator

import (
	"strings"

	"github.com/dizzyfool/genna/model"
)

// Options for generator
type Options struct {
	// Output file path
	Output string

	// List of tables to generate
	// Default []string{"public.*"}
	Tables []string

	// Package sets package name for model
	// Works only with SchemaPackage = false
	Package string

	// Generate model for foreign keys,
	// even if tables not listed in Tables param
	// will not generate fks if schema not listed
	FollowFKs bool

	// Do not replace primary key name to ID
	KeepPK bool

	// Soft delete column
	SoftDelete string

	// Generate model with views e.g. getUsers for users table
	View bool

	// Do not generate discard_unknown_columns tag
	NoDiscard bool

	// Do not generate alias tag
	NoAlias bool

	// Generate search filters
	WithSearch bool

	// Strict types in filters
	StrictSearch bool

	// Stores json field names as in db and target types for them
	// TODO implement
	JSONTypes map[string]string

	// Generate Hooks
	// TODO implement
	UseHooks bool
}

// def fills default values of an options
func (o *Options) def() {
	if strings.Trim(o.Package, " ") == "" {
		o.Package = model.DefaultPackage
	}

	if len(o.Tables) == 0 {
		o.Tables = []string{"public.*"}
	}
}
