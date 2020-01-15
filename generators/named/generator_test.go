package named

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	generator := New()
	options := generator.Options()

	options.Def()
	options.URL = `postgres://genna:genna@localhost:5432/genna?sslmode=disable`
	options.Output = path.Join(os.TempDir(), "model_test.go")
	options.FollowFKs = true

	generator.SetOptions(options)

	if err := generator.Generate(); err != nil {
		t.Errorf("generate error = %v", err)
		return
	}

	generated, err := ioutil.ReadFile(options.Output)
	if err != nil {
		t.Errorf("file not generated = %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	check, err := ioutil.ReadFile(path.Join(path.Dir(filename), "generator_test.output"))
	if err != nil {
		t.Errorf("check file not found = %v", err)
	}

	if string(generated) != string(check) {
		t.Errorf("generated not mathed with check")
		return
	}
}
