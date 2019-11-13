package internal

import (
	"errors"
	"io"

	"github.com/pjbgf/gosystract/cmd/systract"
)

var (
	invalidSyntaxMessage string = "invalid systax"
	usageMessage         string = `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`
)

// Run does basic handling of user input
func Run(output io.Writer, args []string, extract func(dumpFileName string) ([]systract.SystemCall, error)) error {

	if len(args) < 2 || len(args) > 3 {
		output.Write([]byte(usageMessage))
		return errors.New(invalidSyntaxMessage)
	}

	return nil
}
