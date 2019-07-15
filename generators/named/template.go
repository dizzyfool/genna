package named

const templateModel = `//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

{{range .Entities}}
	type Columns{{.GoName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{if .HasRelations}}
		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{end}}
	}
{{end}}
type ColumnsSt struct { {{range .Entities}}
	{{.GoName}} Columns{{.GoName}}{{end}}
}
var Columns = ColumnsSt{ {{range .Entities}}
	{{.GoName}}: Columns{{.GoName}}{ {{range .Columns}}
		{{.GoName}}: "{{.PGName}}",{{end}}{{if .HasRelations}}
		{{range .Relations}}
		{{.GoName}}: "{{.GoName}}",{{end}}{{end}}
	},{{end}}
}
{{range .Entities}}
type Table{{.GoName}} struct {
	Name{{if not .NoAlias}}, Alias{{end}} string
}
{{end}}
type TablesSt struct { {{range .Entities}}
		{{.GoName}} Table{{.GoName}}{{end}}
}
var Tables = TablesSt { {{range .Entities}}
	{{.GoName}}: Table{{.GoName}}{ 
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
