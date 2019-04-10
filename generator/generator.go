package generator

import (
	"bytes"
	"go/format"
	"html/template"
	"os"
	"path"

	"github.com/dizzyfool/genna/model"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Generator is a model files generator
type Generator struct {
	logger  *zap.Logger
	options Options
}

// Result stores result of generator
type Result struct {
	TotalTables       int
	GeneratedModels   int
}

// NewGenerator creates generator
func NewGenerator(options Options, logger *zap.Logger) *Generator {
	options.def()
	return &Generator{
		logger:  logger,
		options: options,
	}
}

// Process processing all tables
func (g Generator) Process(tables []model.Table) (*Result, error) {

	// disclosing asterisks
	toGenerate := model.DiscloseSchemas(tables, g.options.Tables)

	if g.options.FollowFKs {
		// adding models for foreign keys that was not selected for generation by user
		toGenerate = model.FollowFKs(tables, toGenerate)
	} else {
		// filtering relations for models that not listed for generation by user
		tables = model.FilterFKs(tables, toGenerate)
	}

	tmpl, err := template.New(model.DefaultPackage).Parse(templateModel)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template error")
	}

	// making intermediate structs for templates
	pkg := g.Package(tables, toGenerate)

	g.logger.Debug("generating", zap.String("package", pkg.Package))
	var buffer bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buffer, model.DefaultPackage, pkg); err != nil {
		return nil, errors.Wrap(err, "processing model template error")
	}

	unformatted := buffer.Bytes()
	// formatting by go-fmt
	content, err := format.Source(unformatted)
	if err != nil {
		g.logger.Info("formatting file error", zap.Error(err), zap.String("file", g.options.Output))
		// saving file even if there is fmt errors
		content = unformatted
	}

	g.logger.Debug("saving", zap.String("file", g.options.Output))
	file, err := g.File(g.options.Output)
	if err != nil {
		return nil, errors.Wrap(err, "open model file error")
	}

	if _, err := file.Write(content); err != nil {
		return nil, errors.Wrap(err, "writing content to file error")
	}

	return &Result{
		TotalTables:       len(tables),
		GeneratedModels:   len(toGenerate),
	}, nil
}

// Packages makes intermediate structs for templates
// tables - all tables in database
// toGenerate - tables with schemas need to generate, e.g. public.users
func (g Generator) Package(tables []model.Table, toGenerate []string) templatePackage {
	// index for quick access to model
	index := map[string]int{}
	for i, t := range tables {
		index[model.Join(t.Schema, t.Name)] = i
	}

	// on one big file
	toTemplate := make([]model.Table, 0)
	// just go though all tables
	for _, t := range toGenerate {
		if i, ok := index[t]; ok {
			toTemplate = append(toTemplate, tables[i])
		}
	}

	return newTemplatePackage(toTemplate, g.options)
}

// File creates a file for model
func (g Generator) File(filename string) (*os.File, error) {
	directory := path.Dir(filename)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
