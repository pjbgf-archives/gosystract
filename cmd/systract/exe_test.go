package systract

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

func TestExeReader_GetReader_Integration(t *testing.T) {

	t.Run("should error when file not found", func(t *testing.T) {

		reader := NewExeReader("file-that-dont-exist")
		_, err := reader.GetReader()

		if err == nil {
			t.Error("should error")
		}
	})

	t.Run("should be able disassemble current executable", func(t *testing.T) {
		path, _ := os.Executable()

		reader := NewExeReader(path)
		r, err := reader.GetReader()

		if err != nil {
			t.Error("should not error")
		}

		if r == nil {
			t.Error("reader should not be nil")
		}
	})
}

func TestGetFileDumpReader(t *testing.T) {
	reader, err := getFileDumpReader("/bin/echo", "123456")

	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	got := scanner.Text()
	want := "123456"

	if err != nil {
		t.Error("should not error")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
