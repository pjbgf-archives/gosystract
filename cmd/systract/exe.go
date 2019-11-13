package systract

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
)

// ExeReader represents a go executables reader.
// Internally it will call go tool objdump in order to get a disassembled dump of the file.
type ExeReader struct {
	filePath string
}

// NewExeReader initialises a new ExeReader
func NewExeReader(exeFilePath string) *ExeReader {
	return &ExeReader{exeFilePath}
}

// GetReader returns a io.Reader based of the filePath
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
	objDumpFilePath := fmt.Sprintf("/usr/local/go/pkg/tool/%s_%s/objdump", runtime.GOOS, runtime.GOARCH)

	/* #nosec filePath is pre-processed by sanitiseFileName */
	cmd := exec.Command(objDumpFilePath, filePath)
	output, err := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return output, err
}
