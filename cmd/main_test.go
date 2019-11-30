package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pjbgf/go-test/should"
)

var usageMessage string = `gosystract version [ not set ]
Usage:
gosystrac [flags] filePath

Flags:
	--dumpfile, -d    Handles a dump file instead of go executable.
	--template	  Defines a go template for the results.
			  Example: --template="{{- range . }}{{printf "%d - %s\n" .ID .Name}}{{- end}}"
`

func TestMain(t *testing.T) {
	assertThat := func(assumption string, args []string, expected string) {
		should := should.New(t)
		tmpfile, err := ioutil.TempFile("", "fakestdout.*")
		if err != nil {
			t.Errorf("could not setup test properly, got error: %s", err)
		}
		defer os.Remove(tmpfile.Name())

		os.Args = args
		os.Stdout = tmpfile

		main()

		contents, err := ioutil.ReadFile(tmpfile.Name())
		actual := string(contents)

		should.BeEqual(expected, actual, assumption)
	}

	assertThat("should return exit_group call for single-syscall.dump",
		strings.Split("gosystract --dumpfile ../test/single-syscall.dump", " "),
		"1 system calls found:\n    exit_group (231)\n")
}

func TestMain_ErrorCodes(t *testing.T) {
	assertThat := func(assumption, command, expectedErr, expectedOutput string) {
		should := should.New(t)
		exe, _ := os.Executable()

		cmd := exec.Command(exe, "-test.run", "^TestMain_ErrorCodes_Inception$")
		cmd.Env = append(cmd.Env, fmt.Sprintf("ErrorCodes_Args=%s", command))

		output, err := cmd.CombinedOutput()

		e, ok := err.(*exec.ExitError)

		if !ok {
			t.Log("was expecting exit code which did not happen")
			t.FailNow()
		}

		actualOutput := string(output)

		should.BeEqual(expectedErr, e.Error(), assumption)
		should.BeEqual(expectedOutput, actualOutput, assumption)
	}

	assertThat("should exit with code 1 if no args provided", "gosystract", "exit status 1",
		`gosystract version [ not set ]
Usage:
gosystrac [flags] filePath

Flags:
	--dumpfile, -d    Handles a dump file instead of a go executable.
	--template	  Defines a go template for the results.

error: invalid syntax
`)
}

func TestMain_ErrorCodes_Inception(t *testing.T) {
	args := os.Getenv("ErrorCodes_Args")
	if args != "" {
		os.Args = strings.Split(args, " ")

		main()
	}
}
