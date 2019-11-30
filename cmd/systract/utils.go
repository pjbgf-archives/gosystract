package systract

import (
	"errors"
	"os"
	"path"
	"path/filepath"
)

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	return false
}

func sanitiseFileName(input string) (string, error) {
	p := filepath.Clean(input)
	if !path.IsAbs(p) {
		base, err := os.Getwd()
		if err != nil {
			return "", errors.New("error getting current folder")
		}

		return filepath.Join(base, p), nil
	}
	return p, nil
}
