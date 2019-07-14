package search

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"go.uber.org/zap"
)

func prepareReq() (url string, logger *zap.Logger) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"

	logger, _ = config.Build()
	url = `postgres://genna:genna@localhost:5432/genna?sslmode=disable`

	return
}

func TestGenerator_Generate(t *testing.T) {
	generator := New(prepareReq())
	output := path.Join(os.TempDir(), "base_test.go")
	fmt.Println(output)

	err := generator.Generate(Options{
		Output:    output,
		FollowFKs: true,
	})

	if err != nil {
		t.Errorf("generate error = %v", err)
		return
	}

	generated, err := ioutil.ReadFile(output)
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
