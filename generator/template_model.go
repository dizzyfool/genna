package generator

const templateModel = `//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
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
{{range .Models}}
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
		errors["{{.FieldName}}"] = "empty"
	}	
	{{else if eq .ValidationCheck "zero"}}
	if m.{{.FieldName}} == 0 {
		errors["{{.FieldName}}"] = "empty"
	}
	{{else if eq .ValidationCheck "pzero"}}
	if m.{{.FieldName}} != nil && *m.{{.FieldName}} == 0 {
		errors["{{.FieldName}}"] = "empty"
	}
	{{else if eq .ValidationCheck "len"}}
	if isExceedsLen(m.{{.FieldName}}, {{.MaxLen}}) {
		errors["{{.FieldName}}"] = "len"
	}
	{{else if eq .ValidationCheck "plen"}}
	if m.{{.FieldName}} != nil && isExceedsLen(*m.{{.FieldName}}, {{.MaxLen}}) {
		errors["{{.FieldName}}"] = "len"
	}
	{{else if eq .ValidationCheck "enum"}}
	switch m.{{.FieldName}} {
		case {{.Enum}}:
		default:
			errors["{{.FieldName}}"] = "value"
	}
	{{else if eq .ValidationCheck "penum"}}
	if m.{{.FieldName}} != nil { 
		switch *m.{{.FieldName}} {
			case {{.Enum}}:
			default:
				errors["{{.FieldName}}"] = "value"
		}
	}
	{{end}}
	{{end}}
	{{end}}

	return errors, len(errors) == 0
}
{{end}}
{{end}}

{{if .HasValidation}}
func isExceedsLen(v string, len int) bool {
	return utf8.RuneCountInString(v) > len
}
{{end}}
`
