package base

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dizzyfool/genna/lib"
	"github.com/dizzyfool/genna/util"

	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type Generator struct {
	genna.Genna
}

func New(url string, logger *zap.Logger) Generator {
	return Generator{
		Genna: genna.New(url, logger),
	}
}

func (g Generator) Generate(options Options) error {
	entities, err := g.Read(options.Tables, options.FollowFKs, options.UseSqlNulls)
	if err != nil {
		return xerrors.Errorf("read database error: %w", err)
	}

	parsed, err := template.New("base").Parse(templateModel)
	if err != nil {
		return xerrors.Errorf("parsing template error: %w", err)
	}

	pack := NewTemplatePackage(entities, options)

	var buffer bytes.Buffer
	if err := parsed.ExecuteTemplate(&buffer, "base", pack); err != nil {
		return xerrors.Errorf("processing model template error: %w", err)
	}

	saved, err := util.FmtAndSave(buffer.Bytes(), options.Output)
	if err != nil {
		if !saved {
			return xerrors.Errorf("saving file error: %w", err)
		}
		g.Logger.Error("formatting file error", zap.Error(err), zap.String("file", options.Output))
	}

	g.Logger.Info(fmt.Sprintf("succesfully generated %d models\n", len(entities)))

	return nil
}
