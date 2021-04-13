package genna

import (
	"log"
	"os"
	"testing"
)

func prepareReq() (url string, logger *log.Logger) {
	logger = log.New(os.Stderr, "", log.LstdFlags)
	url = `postgres://genna:genna@localhost:5432/genna?sslmode=disable`

	return
}

func TestGenna_Read(t *testing.T) {
	genna := New(prepareReq())

	t.Run("Should read DB", func(t *testing.T) {
		entities, err := genna.Read([]string{"public.*"}, true, false, 9, nil)
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
