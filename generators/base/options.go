package base

import (
	"strings"

	"github.com/dizzyfool/genna/model_old"
)

// Options for generator
type Options struct {
	// Output file path
	Output string

	// List of Tables to generate
	// Default []string{"public.*"}
	Tables []string

	// Package sets package name for model
	// Works only with SchemaPackage = false
	Package string

	// Generate model for foreign keys,
	// even if Tables not listed in Tables param
	// will not generate fks if schema not listed
	FollowFKs bool

	// Do not replace primary key name to ID
	KeepPK bool

	// Soft delete column
	SoftDelete string

	// use sql.Null... instead of pointers
	UseSqlNulls bool

	// Do not generate alias tag
	NoAlias bool

	// Do not generate discard_unknown_columns tag
	NoDiscard bool
}

// def fills default values of an options
func (o *Options) def() {
	if strings.Trim(o.Package, " ") == "" {
		o.Package = model_old.DefaultPackage
	}

	if len(o.Tables) == 0 {
		o.Tables = []string{"public.*"}
	}
}
