package search

const templateSearch = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

import ({{if .HasImports}}{{range .Imports}}
	"{{.}}"{{end}}
	{{end}}
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/types"
)

// base filters
type applier func(query *orm.Query) (*orm.Query, error)

type filterParams struct {
	Table types.ValueAppender
	Field types.ValueAppender
	Value interface{}
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier
}

{{range $model := .Entities}}
type {{.GoName}}Search struct {
	{{range .Columns}}
	{{.GoName}} {{.GoType}}{{end}}
}

func (s *{{.GoName}}Search) Apply(query *orm.Query) *orm.Query { {{range .Columns}}{{if .Relaxed}}
	if !reflect.ValueOf(s.{{.GoName}}).IsNil(){ {{else}}
	if s.{{.GoName}} != nil { {{end}}
		query.Where("{{.Condition}}", filterParams{ {{if .TableExpr}}
			Table: pg.F({{.TableExpr}}),{{else}}
			Table: pg.F("{{.TableName}}"),{{end}}{{if .FieldExpr}}
			Field: pg.F({{.FieldExpr}}),{{else}}
			Field: pg.F("{{.PGName}}"),{{end}}
			Value: s.{{.GoName}},
		})
	}{{end}}
	
	return query
}

func (s *{{.GoName}}Search) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}
{{end}}
`
