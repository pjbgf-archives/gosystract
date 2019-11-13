package internal

import (
	"bytes"
	"testing"

	"errors"

	"github.com/pjbgf/gosystract/cmd/systract"
)

func TestRun(t *testing.T) {

	t.Run("should show usage when no args", func(t *testing.T) {
		args := []string{}
		var output bytes.Buffer

		Run(&output, args, func(dumpFileName string) ([]systract.SystemCall, error) {
			return nil, errors.New("invalid systax")
		})

		got := output.String()
		want := `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
