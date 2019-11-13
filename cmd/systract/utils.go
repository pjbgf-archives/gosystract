package systract

import (
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	return false
}

func sanitiseFileName(input string) (string, error) {
	if !path.IsAbs(input) {
		base, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "error getting current folder")
		}

		return filepath.Join(base, filepath.Clean(input)), nil
	}
	return input, nil
}
