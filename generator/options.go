package generator

import (
	"strings"

	"github.com/dizzyfool/genna/model"
)

// Options for generator
type Options struct {
	// Directory path where files should be saved
	Output string

	// List of tables to generate
	// Default []string{"public.*"}
	Tables []string

	// Prefix for imports
	ImportPath string

	// Package sets package name for model
	// Works only with SchemaPackage = false
	Package string

	// Generate every schema as separate package
	SchemaPackage bool

	// Generates one file for package
	// SchemaPackage | MultiFile | Result
	// true          | true      | each generated package will contain one file
	// true          | false     | each generated package will contain several files, one per model
	// false         | false     | one package for all models separated to different files
	// false         | true      | one big file for all models
	// TODO Make this param as MODE ?
	MultiFile bool

	// Generate model for foreign keys,
	// even if tables not listed in Tables param
	// will not generate fks if schema not listed
	FollowFKs bool

	// Generate model with views e.g. getUsers for users table
	View bool

	// Do not replace primary key name to ID
	KeepPK bool

	// Do not generate discard_unknown_columns tag
	NoDiscard bool

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

	if o.SchemaPackage {
		o.Package = model.DefaultPackage
	}

	if len(o.Tables) == 0 {
		o.Tables = []string{"public.*"}
	}
}
