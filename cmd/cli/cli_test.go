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

		err := Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			return nil, errors.New("invalid systax")
		})

		got := output.String()
		want := `gosystract returns the names and IDs of all system calls being called inside a go application.
Usage: 
	gosystrac goapp.dump
	gosystrac goapp.dump "{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"

To generate a dump file from a go application use: 
	go tool objdump goapp > goapp.dump`

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if err == nil {
			t.Error("should have returned error")
		}
	})

	t.Run("should show syscalls found", func(t *testing.T) {
		args := []string{"gosystract", "filename"}
		var output bytes.Buffer

		err := Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			return []systract.SystemCall{{ID: 1, Name: "abc"}, {ID: 2, Name: "def"}}, nil
		})

		got := output.String()
		want := "2 system calls found:\n    abc (1)\n    def (2)\n"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("should support custom go template for results", func(t *testing.T) {
		args := []string{"gosystract", "filename", "{{- range . }}\"{{.Name}}\",{{- end}}"}
		var output bytes.Buffer

		err := Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			return []systract.SystemCall{{ID: 1, Name: "abc"}, {ID: 2, Name: "def"}}, nil
		})

		got := output.String()
		want := "\"abc\",\"def\","

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if err != nil {
			t.Error(err)
		}
	})

	t.Run("should show message when no syscalls are found", func(t *testing.T) {
		args := []string{"gosystract", "filename"}
		var output bytes.Buffer

		err := Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			return []systract.SystemCall{}, nil
		})

		got := output.String()
		want := "no systems calls were found\n"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if err != nil {
			t.Error(err)
		}
	})
}
