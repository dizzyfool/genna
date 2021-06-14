package model

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

	// Soft delete column
	SoftDelete string

	// use sql.Null... instead of pointers
	UseSQLNulls bool

	// Do not generate alias tag
	NoAlias bool

	// Do not generate discard_unknown_columns tag
	NoDiscard bool

	// Override type for json/jsonb
	JSONTypes map[string]string

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
