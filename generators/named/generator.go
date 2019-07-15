package named

import (
	"github.com/dizzyfool/genna/generators/base"
	"github.com/dizzyfool/genna/generators/model"

	"go.uber.org/zap"
)

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
			templateModel,
			g.Packer(),
		)
}
