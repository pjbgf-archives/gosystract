package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pjbgf/should"
)

func TestMain_E2E(t *testing.T) {
	assertThat := func(assumption string, args []string, expected string) {
		should := should.New(t)
		tmpfile, err := ioutil.TempFile("", "fakestdout.*")
		if err != nil {
			t.Errorf("could not setup test properly, got error: %s", err)
		}
		defer os.Remove(tmpfile.Name())

		os.Args = []string{"gosystract", "--dumpfile", "../test/single-syscall.dump"}
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
