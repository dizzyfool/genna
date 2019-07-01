package generator

const templateModel = `//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}{{if .HasValidation}}
const (
	ErrEmptyValue = "empty"
	ErrMaxLength  = "len"
	ErrWrongValue = "value"
){{end}}

var Columns = struct { {{range .Models}}
	{{.StructName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{end}}
	}{{end}}
}{ {{range .Models}}
	{{.StructName}}: struct { 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{end}}
	}{ {{range .Columns}}
		{{.FieldName}}: "{{.FieldDBName}}",{{end}}{{if .HasRelations}}
		{{range .Relations}}
		{{.FieldName}}: "{{.FieldName}}",{{end}}{{end}}
	},{{end}}
}

var Tables = struct { {{range .Models}}
	{{.StructName}} struct {
		Name{{if .WithAlias}}, Alias{{end}} string
	}{{end}}
}{ {{range .Models}}
	{{.StructName}}: struct {
		Name{{if .WithAlias}}, Alias{{end}} string
	}{ 
		Name: "{{.TableName}}"{{if .WithAlias}},
		Alias: "{{.TableAlias}}",{{end}}
	},{{end}}
}
{{range $model := .Models}}
type {{.StructName}} struct {
	tableName struct{} {{.StructTag}}
	{{range .Columns}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}} {{.FieldComment}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}} {{.FieldComment}}{{end}}{{end}}
}
{{if .HasValidation}}
func (m {{.StructName}}) Validate() (errors map[string]string, valid bool) {
	errors = map[string]string{}

	{{range .Columns}}
	{{if .IsValidatable}}
	{{if eq .ValidationCheck "nil" }}
	if m.{{.FieldName}} == nil {
		errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrEmptyValue
	}	
	{{else if eq .ValidationCheck "zero"}}
	if m.{{.FieldName}} == 0 {
		errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrEmptyValue
	}
	{{else if eq .ValidationCheck "pzero"}}
	if m.{{.FieldName}} != nil && *m.{{.FieldName}} == 0 {
		errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrEmptyValue
	}
	{{else if eq .ValidationCheck "len"}}
	if utf8.RuneCountInString(m.{{.FieldName}}) > {{.MaxLen}} {
		errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrMaxLength
	}
	{{else if eq .ValidationCheck "plen"}}
	if m.{{.FieldName}} != nil && utf8.RuneCountInString(*m.{{.FieldName}}) > {{.MaxLen}} {
		errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrMaxLength
	}
	{{else if eq .ValidationCheck "enum"}}
	switch m.{{.FieldName}} {
		case {{.Enum}}:
		default:
			errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrWrongValue
	}
	{{else if eq .ValidationCheck "penum"}}
	if m.{{.FieldName}} != nil { 
		switch *m.{{.FieldName}} {
			case {{.Enum}}:
			default:
				errors[Columns.{{$model.StructName}}.{{.FieldName}}] = ErrWrongValue
		}
	}
	{{end}}
	{{end}}
	{{end}}

	return errors, len(errors) == 0
}
{{end}}
{{end}}
`
