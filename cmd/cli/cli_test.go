package cli

import (
	"bytes"
	"testing"

	"errors"

	"github.com/pjbgf/go-test/should"
	"github.com/pjbgf/gosystract/cmd/systract"
)

func TestParseInputValues(t *testing.T) {
	assertThat := func(assumption string, args []string, expected string) {
		should := should.New(t)

		_, actual, _, err := parseInputValues(args)

		should.NotError(err, assumption)
		should.BeEqual(expected, actual, assumption)
	}

	assertThat("should handle template flag", []string{"gosystract", "--template=\"test\"", ""}, "test")
}

func TestRun(t *testing.T) {
	assertThat := func(assumption string, args []string,
		stub func() ([]systract.SystemCall, error), expected string,
		expectedToErr bool, expectedErr string) {

		should := should.New(t)
		gitcommit = "TESTVERSION"
		var stdOut, stdErr bytes.Buffer
		var hasErrored bool

		Run(&stdOut, &stdErr, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			return stub()
		}, func(code int) {
			hasErrored = true
		})

		actual := stdOut.String()
		actualErr := stdErr.String()

		should.BeEqual(expectedToErr, hasErrored, assumption)
		should.BeEqual(expected, actual, assumption)
		should.BeEqual(expectedErr, actualErr, assumption)
	}

	assertThat("should show usage when no args", []string{},
		func() ([]systract.SystemCall, error) { return nil, errors.New("invalid systax") },
		"",
		true,
		`gosystract version TESTVERSION
Usage:
gosystrac [flags] filePath

Flags:
	--dumpfile, -d    Handles a dump file instead of go executable.
	--template	  Defines a go template for the results.

error: invalid syntax
`)

	assertThat("should show syscalls found", []string{"gosystract", "filename"},
		func() ([]systract.SystemCall, error) {
			return []systract.SystemCall{{ID: 1, Name: "abc"}, {ID: 2, Name: "def"}}, nil
		},
		"2 system calls found:\n    abc (1)\n    def (2)\n", false, "")

	assertThat("should support custom go template for results",
		[]string{"gosystract", "--template=\"{{- range . }}\"{{.Name}}\",{{- end}}\"", "filename"},
		func() ([]systract.SystemCall, error) {
			return []systract.SystemCall{{ID: 1, Name: "abc"}, {ID: 2, Name: "def"}}, nil
		},
		"\"abc\",\"def\",", false, "")

	assertThat("should show message when no syscalls are found",
		[]string{"gosystract", "filename"},
		func() ([]systract.SystemCall, error) {
			return []systract.SystemCall{}, nil
		},
		"no systems calls were found\n", false, "")

	assertThat("should be able to handle exec files",
		[]string{"gosystract", "filename"},
		func() ([]systract.SystemCall, error) {
			return []systract.SystemCall{}, nil
		},
		"no systems calls were found\n", false, "")

	assertThat("should not write results if extract failed",
		[]string{"gosystract", "filename"},
		func() ([]systract.SystemCall, error) {
			return nil, errors.New("could not extract syscalls")
		},
		"",
		true, "\nerror: invalid syntax\n")

	assertThat("should error for invalid go template syntax",
		[]string{"gosystract", "--template=\"{{$%£}\"", "filename"},
		func() ([]systract.SystemCall, error) {
			return []systract.SystemCall{{ID: 1, Name: "abc"}, {ID: 2, Name: "def"}}, nil
		},
		"",
		true, "\nerror: invalid syntax\n")

	assertThat("should error for invalid go template syntax",
		[]string{"gosystract", "--template=\"{{.Something}}\"", "filename"},
		func() ([]systract.SystemCall, error) {
			return []systract.SystemCall{{ID: 1, Name: "abc"}, {ID: 2, Name: "def"}}, nil
		},
		"",
		true, "\nerror: invalid syntax\n")
}

func TestRun_SourceReaders(t *testing.T) {
	assertThat := func(assumption string, args []string, expected interface{}) {
		should := should.New(t)
		var stdOut, stdErr bytes.Buffer
		var hasErrored bool

		Run(&stdOut, &stdErr, args, func(source systract.SourceReader) ([]systract.SystemCall, error) {
			should.HaveSameType(expected, source, "should be able to handle dump files")
			return []systract.SystemCall{}, nil
		}, func(code int) {
			hasErrored = true
		})

		should.BeFalse(hasErrored, assumption)
	}

	assertThat("should be able to handle exec files",
		[]string{"gosystract", "filename"},
		&systract.ExeReader{})
	assertThat("should be able to handle dump files",
		[]string{"gosystract", "--dumpfile", "filename"},
		&systract.DumpReader{})
}
