package validate

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	pkg    = "pkg"
	keepPK = "keep-pk"
)

// CreateCommand creates generator command
func CreateCommand(logger *zap.Logger) *cobra.Command {
	return base.CreateCommand("validation", "Validation generator for go-pg models", New(logger))
}

// Validate represents validate generator
type Validate struct {
	logger  *zap.Logger
	options Options
}

// New creates generator
func New(logger *zap.Logger) *Validate {
	return &Validate{
		logger: logger,
	}
}

// Logger gets logger
func (g *Validate) Logger() *zap.Logger {
	return g.logger
}

// Options gets options
func (g *Validate) Options() *Options {
	return &g.options
}

// AddFlags adds flags to command
func (g *Validate) AddFlags(command *cobra.Command) {
	base.AddFlags(command)

	flags := command.Flags()
	flags.SortFlags = false

	flags.StringP(pkg, "p", util.DefaultPackage, "package for model files")

	flags.BoolP(keepPK, "k", false, "keep primary key name as is (by default it should be converted to 'ID')")
}

// ReadFlags read flags from command
func (g *Validate) ReadFlags(command *cobra.Command) error {
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

	// setting defaults
	g.options.Def()

	return nil
}

// Generate runs whole generation process
func (g *Validate) Generate() error {
	return base.NewGenerator(g.options.URL, g.logger).
		Generate(
			g.options.Tables,
			g.options.FollowFKs,
			false,
			g.options.Output,
			templateValidate,
			g.Packer(),
		)
}

// Packer returns packer function for compile entities into package
func (g *Validate) Packer() base.Packer {
	return func(entities []model.Entity) (interface{}, error) {
		return NewTemplatePackage(entities, g.options), nil
	}
}
