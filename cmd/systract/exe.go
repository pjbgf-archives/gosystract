package systract

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type ExeReader struct {
	filePath string
}

func NewExeReader(exeFilePath string) *ExeReader {
	return &ExeReader{exeFilePath}
}

func (e *ExeReader) GetReader() (io.Reader, error) {
	filePath, err := sanitiseFileName(e.filePath)
	if err != nil {
		return nil, err
	}
	if !fileExists(filePath) {
		return nil, errors.New("file does not exist or permission denied")
	}

	return getFileDump(filePath)
}

func getFileDump(filePath string) (io.Reader, error) {
	/* #nosec filePath is pre-processed by sanitiseFileName */
	cmd := exec.Command("go", "tool", "objdump", filePath)
	output, err := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return output, err
}
