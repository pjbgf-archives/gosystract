package systract

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// DumpReader represents a go disassembled files reader
type DumpReader struct {
	filePath string
}

// NewDumpReader initialises a new DumpReader
func NewDumpReader(dumpFilePath string) *DumpReader {
	return &DumpReader{dumpFilePath}
}

// GetReader returns a io.Reader based of the filePath
func (d *DumpReader) GetReader() (io.Reader, error) {
	filePath, err := sanitiseFileName(d.filePath)
	if err != nil {
		return nil, err
	}
	if !fileExists(filePath) {
		return nil, errors.New("file does not exist or permission denied")
	}

	/* #nosec filePath is pre-processed by sanitiseFileName */
	return os.Open(filePath)
}
