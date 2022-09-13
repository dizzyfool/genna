package search

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import ({{if .HasImports}}{{range .Imports}}
	"{{.}}"{{end}}
	{{end}}
	"github.com/uptrace/bun"
)

const condition =  "?.? = ?"

// base filters
type applier func(query bun.QueryBuilder) (bun.QueryBuilder, error)

type search struct {
	appliers[] applier
}

func (s *search) apply(query bun.QueryBuilder) {
	for _, applier := range s.appliers {
		applier(query)
	}
}

func (s *search) where(query bun.QueryBuilder, table, field string, value interface{}) {
	
	query.Where(condition, bun.Ident(table), bun.Ident(field), value)
	
}

func (s *search) WithApply(a applier) {
	if s.appliers == nil {
		s.appliers = []applier{}
	}
	s.appliers = append(s.appliers, a)
}

func (s *search) With(condition string, params ...interface{}) {
	s.WithApply(func(query bun.QueryBuilder) (bun.QueryBuilder, error) {
		return query.Where(condition, params...), nil
	})
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query bun.QueryBuilder) bun.QueryBuilder
	Q() applier

	With(condition string, params ...interface{})
	WithApply(a applier)
}

{{range $model := .Entities}}
type {{.GoName}}Search struct {
	search 

	{{range .Columns}}
	{{.GoName}} {{.Type}}{{if .HasTags}} {{.Tag}}{{end}}{{end}}
}

func (s *{{.GoName}}Search) Apply(query bun.QueryBuilder) bun.QueryBuilder { {{range .Columns}}{{if .Relaxed}}
	if !reflect.ValueOf(s.{{.GoName}}).IsNil(){ {{else}}
	if s.{{.GoName}} != nil { {{end}}{{if .UseCustomRender}}
		{{.CustomRender}}{{else}} 
		s.where(query, Tables.{{$model.GoName}}.{{if not $model.NoAlias}}Alias{{else}}Name{{end}}, Columns.{{$model.GoName}}.{{.GoName}}, s.{{.GoName}}){{end}}
	}{{end}}

	s.apply(query)
	
	return query
}

func (s *{{.GoName}}Search) Q() applier {
	return func(query bun.QueryBuilder) (bun.QueryBuilder, error) {
		return s.Apply(query), nil
	}
}
{{end}}
`
