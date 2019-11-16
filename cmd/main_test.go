package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/pjbgf/should"
)

func TestMain_E2E(t *testing.T) {
	assertThat := func(assumption string, args []string, expected string) {
		should := should.New(t)
		cmd := exec.Command(args[0], args[1:]...)

		output, err := cmd.CombinedOutput()
		actual := string(output)

		should.NotError(err, assumption)
		should.BeEqual(expected, actual, assumption)
	}

	assertThat("should return exit_group call for single-syscall.dump",
		strings.Split("gosystract -test.run=TestMain_E2E --dumpfile ../test/single-syscall.dump", " "),
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
		strings.Split("gosystract -test.run=TestMain_ExitCodes_E2E", " "),
		"exit status 1")
}
