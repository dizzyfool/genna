package named

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/generators/model"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// CreateCommand creates generator command
func CreateCommand(logger *zap.Logger) *cobra.Command {
	return base.CreateCommand("model-named", "Basic go-pg model generator with named structures", New(logger))
}

// Generator represents basic named generator
type Generator struct {
	*model.Basic
}

// New creates basic generator
func New(logger *zap.Logger) *Generator {
	return &Generator{
		Basic: model.New(logger),
	}
}

// Generate runs whole generation process
func (g *Generator) Generate() error {
	options := g.Options()
	logger := g.Logger()

	return base.NewGenerator(options.URL, logger).
		Generate(
			options.Tables,
			options.FollowFKs,
			options.UseSQLNulls,
			options.Output,
			Template,
			g.Packer(),
		)
}
