package generator

import (
	"bytes"
	"go/format"
	"html/template"
	"os"
	"path"
	"sort"
	"strings"

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
	TotalTables     int
	GeneratedModels int
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

	// making intermediate structs for templates
	pkg := g.Package(tables, sortTables(toGenerate))

	// generating model
	if err := g.generateAndSave(model.DefaultPackage, templateModel, g.options.Output, pkg); err != nil {
		return nil, err
	}

	// generating search filters
	if g.options.WithSearch {
		output := addSuffix(g.options.Output, "_search")
		if err := g.generateAndSave(model.SearchSuffix, templateSearch, output, pkg); err != nil {
			return nil, err
		}
	}

	return &Result{
		TotalTables:     len(tables),
		GeneratedModels: len(toGenerate),
	}, nil
}

// Package makes intermediate struct for templates
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

func (g Generator) generateAndSave(name, tmpl string, output string, pkg templatePackage) error {
	parsed, err := template.New(name).Parse(tmpl)
	if err != nil {
		return errors.Wrap(err, "parsing template error")
	}

	g.logger.Debug("generating", zap.String("package", pkg.Package))
	var buffer bytes.Buffer
	if err := parsed.ExecuteTemplate(&buffer, name, pkg); err != nil {
		return errors.Wrap(err, "processing model template error")
	}

	unformatted := buffer.Bytes()
	// formatting by go-fmt
	content, err := format.Source(unformatted)
	if err != nil {
		g.logger.Info("formatting file error", zap.Error(err), zap.String("file", output))
		// saving file even if there is fmt errors
		content = unformatted
	}

	g.logger.Debug("saving", zap.String("file", output))
	file, err := g.File(output)
	if err != nil {
		return errors.Wrap(err, "open model file error")
	}

	if _, err := file.Write(content); err != nil {
		return errors.Wrap(err, "writing content to file error")
	}

	return nil
}

func addSuffix(filename, suffix string) string {
	dir, file := path.Split(filename)
	last := strings.LastIndex(file, ".")
	if last == -1 {
		return filename + suffix
	}

	return path.Join(dir, file[:last]+suffix+file[last:])
}

func sortTables(slice []string) []string {
	sort.Slice(slice, func(i, j int) bool {
		si, ti := model.Split(slice[i])
		sj, tj := model.Split(slice[j])

		if si == sj {
			return ti < tj
		}

		if si == model.PublicSchema {
			return true
		}
		if sj == model.PublicSchema {
			return false
		}

		return si < sj
	})

	return slice
}
