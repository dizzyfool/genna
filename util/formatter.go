package util

import (
	"go/format"
	"golang.org/x/xerrors"
	"os"
	"path"
)

// FmtAndSave formats go code and saves file
// if formatting failed it still saves file but return error also
func FmtAndSave(unformatted []byte, filename string) (bool, error) {
	// formatting by go-fmt
	content, fmtErr := format.Source(unformatted)
	if fmtErr != nil {
		// saving file even if there is fmt errors
		content = unformatted
	}

	file, err := file(filename)
	if err != nil {
		return false, xerrors.Errorf("open model file error: %w", err)
	}

	if _, err := file.Write(content); err != nil {
		return false, xerrors.Errorf("writing content to file error: %w", err)
	}

	return true, fmtErr
}

func file(filename string) (*os.File, error) {
	directory := path.Dir(filename)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
