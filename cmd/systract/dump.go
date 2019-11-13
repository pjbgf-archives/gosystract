package systract

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type DumpReader struct {
	filePath string
}

func NewDumpReader(dumpFilePath string) *DumpReader {
	return &DumpReader{dumpFilePath}
}

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
