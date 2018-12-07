package generator

import (
	"html/template"
	"io"
	"os"
	"path"

	"github.com/dizzyfool/genna/model"
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

func (g Generator) Process(tables []model.Table) error {
	index := map[string]int{}
	for i, t := range tables {
		index[model.Join(t.Schema, t.Name)] = i
	}

	toGenerate := model.DiscloseSchemas(g.options.Tables, tables)

	if g.options.FollowFKs {
		toGenerate = model.FollowFKs(toGenerate, tables)
	}

	for _, t := range toGenerate {
		if i, ok := index[t]; ok {
			table := tables[i]

			file, err := g.File(table)
			if err != nil {
				return err
			}

			if err := g.Generate(table, file); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g Generator) Generate(table model.Table, wr io.Writer) error {
	tmpl, err := template.New("model").Parse(TemplateModel)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(wr, "model", newTemplateTable(table, g.options))
}

func (g Generator) File(table model.Table) (*os.File, error) {
	filename := path.Join(g.options.Output, table.PackageName(g.options.SchemaAsPackage, g.options.Package), table.FileName()+".go")
	directory := path.Dir(filename)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
