package generator

const templateModel = `package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

{{range .Models}}
type {{.StructName}} struct {
	tableName struct{} {{.StructTag}}
	{{range .Columns}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}}{{end}}{{end}}
}
{{end}}
`
