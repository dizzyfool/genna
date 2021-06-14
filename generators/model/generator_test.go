package model

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/dizzyfool/genna/model"
)

func TestGenerator_Generate(t *testing.T) {
	generator := New()

	generator.options.Def()
	generator.options.URL = `postgres://genna:genna@localhost:5432/genna?sslmode=disable`
	generator.options.Output = path.Join(os.TempDir(), "model_test.go")
	generator.options.FollowFKs = true
	generator.options.CustomTypes.Add(model.TypePGUuid, "uuid.UUID", "github.com/google/uuid")
	generator.options.GoPgVer = 10
	//generator.options.AddJSONTag = true

	if err := generator.Generate(); err != nil {
		t.Errorf("generate error = %v", err)
		return
	}

	generated, err := ioutil.ReadFile(generator.options.Output)
	if err != nil {
		t.Errorf("file not generated = %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	check, err := ioutil.ReadFile(path.Join(path.Dir(filename), "generator_test.output"))
	if err != nil {
		t.Errorf("check file not found = %v", err)
	}

	if string(generated) != string(check) {
		t.Errorf("generated does not match with check")
		return
	}
}
