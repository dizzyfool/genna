package search

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/model"
	"github.com/spf13/cobra"
)

const (
	keepPK  = "keep-pk"
	noAlias = "no-alias"
	relaxed = "relaxed"
)

// CreateCommand creates generator command
func CreateCommand() *cobra.Command {
	return base.CreateCommand("search", "Search generator for go-pg models", New())
}

// Search represents search generator
type Search struct {
	options Options
}

// New creates generator
func New() *Search {
	return &Search{}
}

// Options gets options
func (g *Search) Options() *Options {
	return &g.options
}

// SetOptions sets options
func (g *Search) SetOptions(options Options) {
	g.options = options
}

// AddFlags adds flags to command
func (g *Search) AddFlags(command *cobra.Command) {
	base.AddFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.BoolP(keepPK, "k", false, "keep primary key name as is (by default it should be converted to 'ID')")

	flags.BoolP(noAlias, "w", false, `do not set 'alias' tag to "t"`)

	flags.BoolP(relaxed, "r", false, "use interface{} type in search filters\n")
}

// ReadFlags read flags from command
func (g *Search) ReadFlags(command *cobra.Command) error {
	var err error

	g.options.URL, g.options.Output, g.options.Package, g.options.Tables, g.options.FollowFKs, g.options.GoPgVer, g.options.CustomTypes, err = base.ReadFlags(command)
	if err != nil {
		return err
	}

	flags := command.Flags()

	if g.options.KeepPK, err = flags.GetBool(keepPK); err != nil {
		return err
	}

	if g.options.NoAlias, err = flags.GetBool(noAlias); err != nil {
		return err
	}

	if g.options.Relaxed, err = flags.GetBool(relaxed); err != nil {
		return err
	}

	// setting defaults
	g.options.Def()

	return nil
}

// Generate runs whole generation process
func (g *Search) Generate() error {
	return base.NewGenerator(g.options.URL).
		Generate(
			g.options.Tables,
			g.options.FollowFKs,
			false,
			g.options.Output,
			Template,
			g.Packer(),
			g.options.GoPgVer,
			g.options.CustomTypes,
		)
}

// Repack runs generator with custom packer
func (g *Search) Repack(packer base.Packer) error {
	return base.NewGenerator(g.options.URL).
		Generate(
			g.options.Tables,
			g.options.FollowFKs,
			false,
			g.options.Output,
			Template,
			packer,
			g.options.GoPgVer,
			g.options.CustomTypes,
		)
}

// Packer returns packer function for compile entities into package
func (g *Search) Packer() base.Packer {
	return func(entities []model.Entity) (interface{}, error) {
		return NewTemplatePackage(entities, g.options), nil
	}
}
