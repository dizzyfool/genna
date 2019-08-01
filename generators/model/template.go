package model

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

var Columns = struct { {{range .Entities}}
	{{.GoName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{end}}
	}{{end}}
}{ {{range .Entities}}
	{{.GoName}}: struct { 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{end}}
	}{ {{range .Columns}}
		{{.GoName}}: "{{.PGName}}",{{end}}{{if .HasRelations}}
		{{range .Relations}}
		{{.GoName}}: "{{.GoName}}",{{end}}{{end}}
	},{{end}}
}

var Tables = struct { {{range .Entities}}
	{{.GoName}} struct {
		Name{{if not .NoAlias }}, Alias{{end}} string
	}{{end}}
}{ {{range .Entities}}
	{{.GoName}}: struct {
		Name{{if not .NoAlias}}, Alias{{end}} string
	}{ 
		Name: "{{.PGFullName}}"{{if not .NoAlias}},
		Alias: "{{.Alias}}",{{end}}
	},{{end}}
}
{{range $model := .Entities}}
type {{.GoName}} struct {
	tableName struct{} {{.Tag}}
	{{range .Columns}}
	{{.GoName}} {{.Type}} {{.Tag}} {{.Comment}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.GoName}} *{{.GoType}} {{.Tag}} {{.Comment}}{{end}}{{end}}
}
{{end}}
`
