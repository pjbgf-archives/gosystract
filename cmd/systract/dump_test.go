package systract

import (
	"os"
	"testing"
)

func TestDumpReader_GetReader_Integration(t *testing.T) {

	t.Run("should error when file not found", func(t *testing.T) {

		reader := NewDumpReader("file-that-dont-exist")
		_, err := reader.GetReader()

		if err == nil {
			t.Error("should error")
		}
	})

	t.Run("should be able open current executable", func(t *testing.T) {
		path, _ := os.Executable()

		reader := NewDumpReader(path)
		r, err := reader.GetReader()

		if err != nil {
			t.Error("should not error")
		}

		if r == nil {
			t.Error("reader should not be nil")
		}
	})
}
