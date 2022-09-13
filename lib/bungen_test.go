package bungen

import (
	"log"
	"os"
	"testing"
)

func prepareReq() (url string, logger *log.Logger) {
	logger = log.New(os.Stderr, "", log.LstdFlags)
	url = `postgres://some_user:some_password@localhost:5432/some_db?sslmode=disable`

	return
}

func TestBungen_Read(t *testing.T) {
	bungenCLI := New(prepareReq())

	t.Run("Should read DB", func(t *testing.T) {
		entities, err := bungenCLI.Read([]string{"public.*"}, true, false, nil)
		if err != nil {
			t.Errorf("Bungen.Read error %v", err)
			return
		}

		if ln := len(entities); ln != 3 {
			t.Errorf("len(entities) = %v, want %v", ln, 3)
			return
		}
	})
}
