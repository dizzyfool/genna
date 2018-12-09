package generator

import (
	"html/template"
	"testing"
)

func Test_templateTable_uniqualizeFields(t *testing.T) {
	type fields struct {
		StructName   string
		StructTag    template.HTML
		Columns      []templateColumn
		HasRelations bool
		Relations    []templateRelation
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Should uniqualize field names",
			fields: fields{
				Columns: []templateColumn{
					{FieldName: "User"},
					{FieldName: "User"},
					{FieldName: "User1"},
					{FieldName: "User1"},
					{FieldName: "WordRel"},
				},
				Relations: []templateRelation{
					{FieldName: "User"},
					{FieldName: "Word"},
					{FieldName: "Word"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := templateTable{
				StructName:   tt.fields.StructName,
				StructTag:    tt.fields.StructTag,
				Columns:      tt.fields.Columns,
				HasRelations: tt.fields.HasRelations,
				Relations:    tt.fields.Relations,
			}
			tbl.uniqualizeFields()

			index := map[string]bool{}
			for _, col := range tbl.Columns {
				if _, ok := index[col.FieldName]; !ok {
					index[col.FieldName] = true
					continue
				}
				t.Errorf("Column has not unique name %s", col.FieldName)
			}

			for _, rel := range tbl.Relations {
				if _, ok := index[rel.FieldName]; !ok {
					index[rel.FieldName] = true
					continue
				}
				t.Errorf("Relation has not unique name %s", rel.FieldName)
			}
		})
	}
}
