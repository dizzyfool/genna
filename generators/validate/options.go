package validate

import (
	"strings"

	"github.com/LdDl/bungen/generators/base"
	"github.com/LdDl/bungen/util"
)

// Options for generator
type Options struct {
	base.Options

	// Package sets package name for model
	// Works only with SchemaPackage = false
	Package string

	// Do not replace primary key name to ID
	KeepPK bool
}

// Def fills default values of an options
func (o *Options) Def() {
	o.Options.Def()

	if strings.Trim(o.Package, " ") == "" {
		o.Package = util.DefaultPackage
	}
}
