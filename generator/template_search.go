package generator

const templateSearch = `//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import ({{if .HasSearchImports}}{{range .SearchImports}}
	"{{.}}"{{end}}
	{{end}}
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// base filters

type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	custom map[string][]interface{}
}

func (s *search) apply(table string, values map[string]interface{}, query *orm.Query) *orm.Query {
	for field, value := range values {
		if value != nil {
			query.Where("?.? = ?", pg.F(table), pg.F(field), value)
		}
	}

	if s.custom != nil {
		for condition, params := range s.custom {
			query.Where(condition, params...)
		}
	}

	return query
}

func (s *search) with(condition string, params ... interface{}) {
	if s.custom == nil {
		s.custom = map[string][]interface{}{}
	}
	s.custom[condition] = params
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier
}

{{range $model := .Entities}}
type {{.SearchStructName}} struct {
	search

	{{range .Columns}}{{if .IsSearchable}}
	{{.FieldName}} {{.SearchFieldType}}{{end}}{{end}}
}

func (s *{{.SearchStructName}}) Apply(query *orm.Query) *orm.Query {
	return s.apply(Tables.{{.StructName}}.{{if .WithAlias}}Alias{{else}}Name{{end}}, map[string]interface{}{ {{range .Columns}}{{if .IsSearchable}}
		Columns.{{$model.StructName}}.{{.FieldName}}: s.{{.FieldName}},{{end}}{{end}}
	}, query)
}

func (s *{{.SearchStructName}}) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}
{{end}}
`
