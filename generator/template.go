package generator

const templateModel = `package {{.Package}}{{if and .HasImports .WithModel}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}{{if .WithColumns}}

var Columns = struct { {{range .Models}}
	{{.StructName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string 
	}{{end}}
}{ {{range .Models}}
	{{.StructName}}: struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string 
	}{ {{range .Columns}}
		{{.FieldName}}: "{{.FieldDBName}}",{{end}}
	},{{end}}
}{{end}}{{if .WithModel}}
{{range .Models}}
type {{.StructName}} struct {
	tableName struct{} {{.StructTag}}
	{{range .Columns}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}}{{end}}{{end}}
}
{{end}}{{end}}
`
