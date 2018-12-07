package generator

import "strings"

type Options struct {
	// Directory path where files should be saved
	Output string

	// Package sets package name for model
	// With SchemaAsPackage = true this param works only for public schema
	// Default 'model'
	Package string

	// Generate every schema as separate package
	SchemaAsPackage bool

	// If SchemaAsPackage is true
	// PackagePrefix holds prefix for foreign keys
	PackagePrefix string

	// List of tables to generate
	// Default []string{"public.*"}
	Tables []string

	// Generate model with views e.g. getUsers for users table
	View bool

	// Generate model for foreign keys
	FollowFKs bool

	// Stores json field names as in db and target types for them
	// TODO implement
	JsonTypes map[string]string

	// Generate Hooks
	// TODO implement
	UseHooks bool

	// Do not replace primary key name to ID
	KeepPK bool

	// Do not generate discard_unknown_columns tag
	NoDiscard bool
}

func (o *Options) def() {
	if strings.Trim(o.Package, " ") == "" {
		o.Package = "model"
	}

	if len(o.Tables) == 0 {
		o.Tables = []string{"public.*"}
	}
}
