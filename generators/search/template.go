package search

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import ({{if .HasImports}}{{range .Imports}}
	"{{.}}"{{end}}
	{{end}}
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

const condition =  "?.? = ?"

// base filters
type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	appliers[] applier
}

func (s *search) apply(query *orm.Query) {
	for _, applier := range s.appliers {
		query.Apply(applier)
	}
}

func (s *search) where(query *orm.Query, table, field string, value interface{}) {
	query.Where(condition, pg.F(table), pg.F(field), value)
}

func (s *search) WithApply(a applier) {
	if s.appliers == nil {
		s.appliers = []applier{}
	}
	s.appliers = append(s.appliers, a)
}

func (s *search) With(condition string, params ...interface{}) {
	s.WithApply(func(query *orm.Query) (*orm.Query, error) {
		return query.Where(condition, params...), nil
	})
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier

	With(condition string, params ...interface{})
	WithApply(a applier)
}

{{range $model := .Entities}}
type {{.GoName}}Search struct {
	search 

	{{range .Columns}}
	{{.GoName}} {{.Type}}{{end}}
}

func (s *{{.GoName}}Search) Apply(query *orm.Query) *orm.Query { {{range .Columns}}{{if .Relaxed}}
	if !reflect.ValueOf(s.{{.GoName}}).IsNil(){ {{else}}
	if s.{{.GoName}} != nil { {{end}}{{if .UseCustomRender}}
		{{.CustomRender}}{{else}} 
		s.where(query, Tables.{{$model.GoName}}.{{if not $model.NoAlias}}Alias{{else}}Name{{end}}, Columns.{{$model.GoName}}.{{.GoName}}, s.{{.GoName}}){{end}}
	}{{end}}

	s.apply(query)
	
	return query
}

func (s *{{.GoName}}Search) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}
{{end}}
`
