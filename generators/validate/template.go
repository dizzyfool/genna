package validate

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

const (
	ErrEmptyValue = "empty"
	ErrMaxLength  = "len"
	ErrWrongValue = "value"
)

{{range $model := .Entities}}
func (m {{.GoName}}) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	{{range .Columns}}
	{{if eq .Check "nil" }}
	if m.{{.GoName}} == nil {
		errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrEmptyValue
	}	
	{{else if eq .Check "zero"}}
	if m.{{.GoName}} == 0 {
		errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrEmptyValue
	}
	{{else if eq .Check "pzero"}}
	if m.{{.GoName}} != nil && *m.{{.GoName}} == 0 {
		errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrEmptyValue
	}
	{{else if eq .Check "len"}}
	if utf8.RuneCountInString(m.{{.GoName}}) > {{.MaxLen}} {
		errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrMaxLength
	}
	{{else if eq .Check "plen"}}
	if m.{{.GoName}} != nil && utf8.RuneCountInString(*m.{{.GoName}}) > {{.MaxLen}} {
		errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrMaxLength
	}
	{{else if eq .Check "enum"}}
	switch m.{{.GoName}} {
		case {{.Enum}}:
		default:
			errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrWrongValue
	}
	{{else if eq .Check "penum"}}
	if m.{{.GoName}} != nil { 
		switch *m.{{.GoName}} {
			case {{.Enum}}:
			default:
				errors[Columns.{{$model.GoName}}.{{.GoName}}] = ErrWrongValue
		}
	}
	{{end}}
	{{end}}

	return errors, len(errors) == 0
}
{{end}}
`
