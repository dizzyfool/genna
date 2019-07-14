package search

import (
	"strings"

	"github.com/dizzyfool/genna/util"
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

	// Do not generate alias tag
	NoAlias bool

	// Strict types in filters
	Relaxed bool
}

// def fills default values of an options
func (o *Options) def() {
	if strings.Trim(o.Package, " ") == "" {
		o.Package = util.DefaultPackage
	}

	if len(o.Tables) == 0 {
		o.Tables = []string{"public.*"}
	}
}
