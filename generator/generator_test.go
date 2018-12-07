package generator

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/dizzyfool/genna/model"
)

func TestDo(t *testing.T) {
	user := model.Table{
		Schema: model.PublicSchema,
		Name:   "users",
		Columns: []model.Column{
			{
				Name:       "userId",
				Type:       model.TypeInt8,
				IsPK:       true,
				IsNullable: false,
			},
			{
				Name:       "email",
				Type:       model.TypeVarchar,
				IsNullable: true,
			},
			{
				Name:       "locationId",
				Type:       model.TypeInt8,
				IsNullable: false,
				IsFK:       true,
			},
			{
				Name: "companyId",
				Type: model.TypeInt8,
				IsFK: true,
			},
			{
				Name: "createdAt",
				Type: model.TypeTimestamp,
			},
		},
		Relations: []model.Relation{
			{
				Type:         model.HasOne,
				SourceSchema: model.PublicSchema,
				SourceTable:  "users",
				SourceColumn: "locationId",
				TargetSchema: "geo",
				TargetTable:  "locations",
				TargetColumn: "locationId",
			},
			{
				Type:         model.HasOne,
				SourceSchema: model.PublicSchema,
				SourceTable:  "users",
				SourceColumn: "companyId",
				TargetSchema: model.PublicSchema,
				TargetTable:  "companies",
				TargetColumn: "companyId",
			},
		},
	}

	company := model.Table{
		Schema: model.PublicSchema,
		Name:   "companies",
		Columns: []model.Column{
			{
				Name:       "companyId",
				Type:       model.TypeInt8,
				IsPK:       true,
				IsNullable: false,
			},
			{
				Name:       "title",
				Type:       model.TypeVarchar,
				IsNullable: true,
			},
		},
	}

	location := model.Table{
		Schema: "geo",
		Name:   "locations",
		Columns: []model.Column{
			{
				Name:       "locationId",
				Type:       model.TypeInt8,
				IsPK:       true,
				IsNullable: false,
			},
			{
				Name:       "title",
				Type:       model.TypeVarchar,
				IsNullable: true,
			},
		},
	}

	lang := model.Table{
		Schema: "geo",
		Name:   "languages",
		Columns: []model.Column{
			{
				Name:       "languageId",
				Type:       model.TypeInt8,
				IsPK:       true,
				IsNullable: false,
			},
			{
				Name:       "title",
				Type:       model.TypeVarchar,
				IsNullable: true,
			},
		},
	}

	unused := model.Table{
		Schema: model.PublicSchema,
		Name:   "unused",
		Columns: []model.Column{
			{
				Name:       "unusedId",
				Type:       model.TypeInt8,
				IsPK:       true,
				IsNullable: false,
			},
			{
				Name:       "title",
				Type:       model.TypeVarchar,
				IsNullable: true,
			},
		},
	}

	_, filename, _, _ := runtime.Caller(0)

	// just for test
	generator := NewGenerator(Options{
		Package:       "test", // model.DefaultPackage,
		Tables:        []string{"public.users", "geo.*"},
		FollowFKs:     true,
		Output:        path.Dir(filename) + "/../test/",
		SchemaPackage: true,
		MultiFile:     true,
		ImportPath:    "github.com/dizzyfool/genna/test",
		KeepPK:        false, // try true
		NoDiscard:     false, // try true
	})

	err := generator.Process([]model.Table{unused, user, company, location, lang})
	fmt.Print(err)
}
