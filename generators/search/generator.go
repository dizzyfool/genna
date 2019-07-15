package search

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	pkg     = "pkg"
	keepPK  = "keep-pk"
	noAlias = "no-alias"
	relaxed = "relaxed"
)

// Search represents search generator
type Search struct {
	logger  *zap.Logger
	options Options
}

// New creates generator
func New(logger *zap.Logger) *Search {
	return &Search{
		logger: logger,
	}
}

// Logger gets logger
func (g *Search) Logger() *zap.Logger {
	return g.logger
}

// AddFlags adds flags to command
func (g *Search) AddFlags(command *cobra.Command) {
	base.AddFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.StringP(pkg, "p", util.DefaultPackage, "package for model files")

	flags.BoolP(keepPK, "k", false, "keep primary key name as is (by default it should be converted to 'ID')")

	flags.BoolP(noAlias, "w", false, `do not set 'alias' tag to "t"`)

	flags.BoolP(relaxed, "r", false, "use interface{} type in search filters\n")
}

// ReadFlags read flags from command
func (g *Search) ReadFlags(command *cobra.Command) error {
	var err error

	g.options.URL, g.options.Output, g.options.Tables, g.options.FollowFKs, err = base.ReadFlags(command)
	if err != nil {
		return err
	}

	flags := command.Flags()

	if g.options.Package, err = flags.GetString(pkg); err != nil {
		return err
	}

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
	return base.NewGenerator(g.options.URL, g.logger).
		Generate(
			g.options.Tables,
			g.options.FollowFKs,
			false,
			g.options.Output,
			templateSearch,
			g.Packer(),
		)
}

// Packer returns packer function for compile entities into package
func (g *Search) Packer() base.Packer {
	return func(entities []model.Entity) interface{} {
		return NewTemplatePackage(entities, g.options)
	}
}
