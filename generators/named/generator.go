package named

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/generators/model"

	"github.com/spf13/cobra"
)

// CreateCommand creates generator command
func CreateCommand() *cobra.Command {
	return base.CreateCommand("model-named", "Basic go-pg model generator with named structures", New())
}

// Generator represents basic named generator
type Generator struct {
	*model.Basic
}

// New creates basic generator
func New() *Generator {
	return &Generator{
		Basic: model.New(),
	}
}

// Generate runs whole generation process
func (g *Generator) Generate() error {
	options := g.Options()
	return base.NewGenerator(options.URL).
		Generate(
			options.Tables,
			options.FollowFKs,
			options.UseSQLNulls,
			options.Output,
			Template,
			g.Packer(),
			options.GoPgVer,
			options.CustomTypes,
		)
}
