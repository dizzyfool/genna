package validate

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"

	"github.com/spf13/cobra"
)

const (
	pkg    = "pkg"
	keepPK = "keep-pk"
)

// CreateCommand creates generator command
func CreateCommand() *cobra.Command {
	return base.CreateCommand("validation", "Validation generator for go-pg models", New())
}

// Validate represents validate generator
type Validate struct {
	options Options
}

// New creates generator
func New() *Validate {
	return &Validate{}
}

// Options gets options
func (g *Validate) Options() *Options {
	return &g.options
}

// SetOptions sets options
func (g *Validate) SetOptions(options Options) {
	g.options = options
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

	g.options.URL, g.options.Output, g.options.Tables, g.options.FollowFKs, g.options.GoPgVer, err = base.ReadFlags(command)
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
	return base.NewGenerator(g.options.URL).
		Generate(
			g.options.Tables,
			g.options.FollowFKs,
			false,
			g.options.Output,
			Template,
			g.Packer(),
			g.options.GoPgVer,
		)
}

// Packer returns packer function for compile entities into package
func (g *Validate) Packer() base.Packer {
	return func(entities []model.Entity) (interface{}, error) {
		return NewTemplatePackage(entities, g.options), nil
	}
}
