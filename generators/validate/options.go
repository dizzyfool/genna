package validate

import (
	"strings"

	"github.com/dizzyfool/genna/util"
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
