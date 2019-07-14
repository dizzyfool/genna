package genna

import (
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

func TestGenna_Read(t *testing.T) {
	genna := New(prepareReq())

	t.Run("Should read db", func(t *testing.T) {
		entities, err := genna.Read([]string{"public.*"}, true)
		if err != nil {
			t.Errorf("Genna.Read error %v", err)
			return
		}

		if ln := len(entities); ln != 3 {
			t.Errorf("len(entities) = %v, want %v", ln, 3)
			return
		}
	})
}
