package generator

const templateModel = `//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

{{range .Models}}
	type Columns{{.StructName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{end}}
	}
{{end}}

type ColumnsSt struct { {{range .Models}}
	{{.StructName}} Columns{{.StructName}}{{end}}
}

var Columns = ColumnsSt{ {{range .Models}}
	{{.StructName}}: Columns{{.StructName}}{ {{range .Columns}}
		{{.FieldName}}: "{{.FieldDBName}}",{{end}}{{if .HasRelations}}
		{{range .Relations}}
		{{.FieldName}}: "{{.FieldName}}",{{end}}{{end}}
	},{{end}}
}

{{range .Models}}
type Table{{.StructName}} struct {
	Name{{if .WithAlias}}, Alias{{end}} string
}
{{end}}

type TablesSt struct { {{range .Models}}
		{{.StructName}} Table{{.StructName}}{{end}}
}

var Tables = TablesSt { {{range .Models}}
	{{.StructName}}: Table{{.StructName}}{ 
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
{{end}}
`