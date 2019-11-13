package internal

import (
	"bytes"
	"reflect"
	"testing"

	"errors"

	"github.com/pjbgf/gosystract/cmd/systract"
)

func TestParseInputValues(t *testing.T) {

	t.Run("should handle template flag", func(t *testing.T) {

		input := []string{"gosystract", "--template=\"test\"", ""}
		_, template, _, err := parseInputValues(input)

		got := template
		want := "test"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if err != nil {
			t.Error("should not error")
		}
	})
}

func TestRun(t *testing.T) {

	t.Run("should show usage when no args", func(t *testing.T) {
		args := []string{}
		var output bytes.Buffer

		err := Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			return nil, errors.New("invalid systax")
		})

		got := output.String()
		want := `Usage:
	gosystrac [flags] filePath

Flags:
	--dumpfile, -d  	Handles a dump file instead of go executables.
						To generate a dump file use: go tool objdump exeFilePath > file.dump

	--template			Define a go template for the results. 
						Example: {{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}
`

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
		args := []string{"gosystract", "filename", "--template=\"{{- range . }}\"{{.Name}}\",{{- end}}\""}
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

	t.Run("should be able to handle exec files", func(t *testing.T) {
		args := []string{"gosystract", "filename"}
		var output bytes.Buffer

		Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {

			// As interface types are only used for static typing, a
			// common idiom to find the reflection Type for an interface
			// type Foo is to use a *Foo value.
			exeReaderType := reflect.TypeOf((*systract.ExeReader)(nil))

			sourceReaderType := reflect.TypeOf(source)

			if sourceReaderType != exeReaderType {
				t.FailNow()
			}
			return []systract.SystemCall{}, nil
		})
	})

	t.Run("should be able to handle dump files", func(t *testing.T) {
		args := []string{"gosystract", "--dumpfile", "filename"}
		var output bytes.Buffer

		Run(&output, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {

			// As interface types are only used for static typing, a
			// common idiom to find the reflection Type for an interface
			// type Foo is to use a *Foo value.
			dumpReaderType := reflect.TypeOf((*systract.DumpReader)(nil))

			sourceReaderType := reflect.TypeOf(source)

			if sourceReaderType != dumpReaderType {
				t.FailNow()
			}
			return []systract.SystemCall{}, nil
		})
	})
}
