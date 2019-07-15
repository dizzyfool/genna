package validate

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
}

// def fills default values of an options
func (o *Options) Def() {
	o.Options.Def()

	if strings.Trim(o.Package, " ") == "" {
		o.Package = util.DefaultPackage
	}
}
