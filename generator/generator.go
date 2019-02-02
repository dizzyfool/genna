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
	GeneratedFiles    int
	GeneratedPackages int
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
	packages := g.Packages(tables, toGenerate)
	index := map[string]struct{}{}

	for _, pkg := range packages {
		if _, ok := index[pkg.Package]; !ok {
			g.logger.Debug("generating", zap.String("package", pkg.Package))
			index[pkg.Package] = struct{}{}
		}

		var buffer bytes.Buffer

		err := tmpl.ExecuteTemplate(&buffer, model.DefaultPackage, pkg)
		if err != nil {
			return nil, errors.Wrap(err, "processing model template error")
		}

		unformatted := buffer.Bytes()
		// formatting by go-fmt
		content, err := format.Source(unformatted)
		if err != nil {
			g.logger.Info("formatting file error", zap.Error(err), zap.String("file", pkg.FileName))
			// saving file even if there is fmt errors
			content = unformatted
		}

		g.logger.Debug("saving", zap.String("file", pkg.FileName))
		file, err := g.File(pkg.FileName)
		if err != nil {
			return nil, errors.Wrap(err, "open model file error")
		}

		if _, err := file.Write(content); err != nil {
			return nil, errors.Wrap(err, "writing content to file error")
		}
	}

	return &Result{
		TotalTables:       len(tables),
		GeneratedModels:   len(toGenerate),
		GeneratedFiles:    len(packages),
		GeneratedPackages: len(index),
	}, nil
}

// SchemasWithTables gets schemas from table names
func SchemasWithTables(tables []string) map[string][]string {
	schemas := map[string][]string{}
	for _, t := range tables {
		schema, _ := model.Split(t)
		if _, ok := schemas[schema]; !ok {
			schemas[schema] = []string{}
		}

		schemas[schema] = append(schemas[schema], t)
	}

	return schemas
}

// Packages makes intermediate structs for templates
// tables - all tables in database
// toGenerate - tables with schemas need to generate, e.g. public.users
func (g Generator) Packages(tables []model.Table, toGenerate []string) []templatePackage {
	// index for quick access to model
	index := map[string]int{}
	for i, t := range tables {
		index[model.Join(t.Schema, t.Name)] = i
	}

	// on one big file
	if !g.options.SchemaPackage && !g.options.MultiFile {
		toTemplate := make([]model.Table, 0)
		// just go though all tables
		for _, t := range toGenerate {
			if i, ok := index[t]; ok {
				toTemplate = append(toTemplate, tables[i])
			}
		}

		return []templatePackage{newMultiPackage(g.options.Package, toTemplate, g.options)}
	}

	result := make([]templatePackage, 0)

	// single file for each package
	if g.options.SchemaPackage && !g.options.MultiFile {
		swt := SchemasWithTables(toGenerate)
		for _, tbls := range swt {
			toTemplate := make([]model.Table, 0)
			for _, t := range tbls {
				if i, ok := index[t]; ok {
					toTemplate = append(toTemplate, tables[i])
				}
			}
			result = append(result, newMultiPackage(
				toTemplate[0].PackageName(true, g.options.Package), toTemplate, g.options,
			))
		}

		return result
	}

	// many files for each model separated by packages
	if g.options.SchemaPackage && g.options.MultiFile {
		swt := SchemasWithTables(toGenerate)
		for _, tbls := range swt {
			toColumns := make([]model.Table, 0)
			for _, t := range tbls {
				if i, ok := index[t]; ok {
					toColumns = append(toColumns, tables[i])

					result = append(result, newSinglePackage(tables[i], g.options))
				}
			}

			result = append(result, newColumnsPackage(
				toColumns[0].PackageName(true, g.options.Package), toColumns, g.options,
			))
		}

		return result
	}

	// many files for each model in one package
	if !g.options.SchemaPackage && g.options.MultiFile {
		swt := SchemasWithTables(toGenerate)
		toColumns := make([]model.Table, 0)
		for _, tbls := range swt {
			for _, t := range tbls {
				if i, ok := index[t]; ok {
					toColumns = append(toColumns, tables[i])

					result = append(result, newSinglePackage(tables[i], g.options))
				}
			}
		}

		result = append(result, newColumnsPackage(
			g.options.Package, toColumns, g.options,
		))

		return result
	}

	return nil
}

// File creates a file for model
func (g Generator) File(filename string) (*os.File, error) {
	directory := path.Dir(filename)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
