package search

import (
	"strings"

	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/util"
)

// Options for generator
type Options struct {
	base.Options

	// Package sets package name for model
	// Works only with SchemaPackage = false
	Package string

	// Do not replace primary key name to ID
	KeepPK bool

	// Do not generate alias tag
	NoAlias bool

	// Strict types in filters
	Relaxed bool

	// Add json tag to models
	AddJSONTag bool
}

// Def fills default values of an options
func (o *Options) Def() {
	o.Options.Def()

	if strings.Trim(o.Package, " ") == "" {
		o.Package = util.DefaultPackage
	}
}
