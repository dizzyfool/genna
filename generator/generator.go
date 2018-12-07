package generator

import (
	"bytes"
	"github.com/dizzyfool/genna/model"
	"go/format"
	"html/template"
	"os"
	"path"
)

type Generator struct {
	options Options
}

func NewGenerator(options Options) *Generator {
	options.def()
	return &Generator{
		options: options,
	}
}

// Processing all tables
func (g Generator) Process(tables []model.Table) error {
	// disclosing asterisks
	toGenerate := model.DiscloseSchemas(tables, g.options.Tables)

	if g.options.FollowFKs {
		// adding models for foreign keys that was not selected for generation by user
		toGenerate = model.FollowFKs(tables, toGenerate)
	} else {
		// filtering relations for models that not listed for generation by user
		tables = model.FilterFKs(tables, toGenerate)
	}

	tmpl, err := template.New(model.DefaultPackage).Parse(TemplateModel)
	if err != nil {
		return err
	}

	// making intermediate structs for templates
	packages := g.Packages(tables, toGenerate)

	for _, pkg := range packages {
		var buffer bytes.Buffer

		err := tmpl.ExecuteTemplate(&buffer, model.DefaultPackage, pkg)
		if err != nil {
			return err
		}

		// formatting by go-fmt
		content, err := format.Source(buffer.Bytes())
		if err != nil {
			return err
		}

		file, err := g.File(pkg.FileName)
		if err != nil {
			return err
		}

		if _, err := file.Write(content); err != nil {
			return err
		}
	}

	return nil
}

// Schemas get schemas from table names
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
func (g Generator) Packages(tables []model.Table, toGenerate []string) []*templatePackage {
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

		return []*templatePackage{newMultiPackage(g.options.Package, toTemplate, g.options)}
	}

	result := make([]*templatePackage, 0)

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

	// many files for each model
	if g.options.MultiFile {
		for _, t := range toGenerate {
			if i, ok := index[t]; ok {
				result = append(result, newSinglePackage(tables[i], g.options))
			}
		}
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
