package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pjbgf/should"
)

var usageMessage string = `gosystract version [ not set ]
Usage:
gosystrac [flags] filePath

Flags:
	--dumpfile, -d    Handles a dump file instead of go executable.
	--template	  Defines a go template for the results.
			  Example: --template="{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"
`

func TestMain_E2E(t *testing.T) {
	assertThat := func(assumption string, args []string, expected string, expectedErr error) {
		should := should.New(t)
		tmpfile, err := ioutil.TempFile("", "fakestdout.*")
		if err != nil {
			t.Errorf("could not setup test properly, got error: %s", err)
		}
		defer os.Remove(tmpfile.Name())

		var actualErr error
		os.Args = args
		os.Stdout = tmpfile
		onError = func(err error) {
			actualErr = err
		}

		main()

		contents, err := ioutil.ReadFile(tmpfile.Name())
		actual := string(contents)

		should.BeEqual(expectedErr, actualErr, assumption)
		should.BeEqual(expected, actual, assumption)
	}

	assertThat("should return exit_group call for single-syscall.dump",
		strings.Split("gosystract --dumpfile ../test/single-syscall.dump", " "),
		"1 system calls found:\n    exit_group (231)\n", nil)

	assertThat("should error for invalid syntax",
		strings.Split("gosystract", " "),
		usageMessage,
		errors.New("invalid syntax"))
}

func TestMain_ExitCodes_E2E(t *testing.T) {
	assertThat := func(assumption string, args []string, expected string) {
		should := should.New(t)
		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Run()

		e, _ := err.(*exec.ExitError)

		should.BeEqual(expected, e.Error(), assumption)
	}

	assertThat("should exit with exit code 1 for invalid syntax",
		strings.Split("go run main.go -test.run=TestMain_ExitCodes_E2E", " "),
		"exit status 1")
}
