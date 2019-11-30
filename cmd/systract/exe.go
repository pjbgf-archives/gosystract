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

// GetReader returns a io.ReadCloser based of the filePath
func (e *ExeReader) GetReader() (io.ReadCloser, error) {
	filePath, err := sanitiseFileName(e.filePath)
	if err != nil {
		return nil, err
	}
	if !fileExists(filePath) {
		return nil, errors.New("file does not exist or permission denied")
	}

	objDumpFilePath := getObjDumpFilePath()
	return getFileDumpReader(objDumpFilePath, filePath)
}

func getObjDumpFilePath() string {
	return fmt.Sprintf("/usr/local/go/pkg/tool/%s_%s/objdump", runtime.GOOS, runtime.GOARCH)
}

func getFileDumpReader(objDumpFilePath, filePath string) (io.ReadCloser, error) {
	/* #nosec filePath is pre-processed by sanitiseFileName */
	cmd := exec.Command(objDumpFilePath, filePath)

	if !fileExists(objDumpFilePath) {
		/* #nosec filePath is pre-processed by sanitiseFileName */
		cmd = exec.Command("go", "tool", "objdump", filePath)
	}

	output, err := cmd.StdoutPipe()
	defer cmd.Start()

	return output, err
}
