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
				Name: "locationId",
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
				SourceColumn: "locationId",
				TargetSchema: "geo",
				TargetTable:  "locations",
				TargetColumn: "locationId",
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

	_, filename, _, _ := runtime.Caller(0)

	generator := NewGenerator(Options{
		Package:         "test", // try model
		Tables:          []string{"*"},
		FollowFKs:       true, // TODO false is not working because need to ignore some relations
		Output:          path.Dir(filename) + "/../",
		SchemaAsPackage: false, // TODO true is not working because of invalid imports
		KeepPK:          false, // try true
		NoDiscard:       false, // try true
	})

	err := generator.Process([]model.Table{user, location})
	fmt.Print(err)
}
