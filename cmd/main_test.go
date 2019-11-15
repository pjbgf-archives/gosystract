package main

import (
	"bytes"
	"testing"

	"github.com/pjbgf/should"
)

func TestMain_Integration(t *testing.T) {
	assertThat := func(assumption, filePath, expected string) {
		should := should.New(t)
		var output bytes.Buffer
		run(&output, []string{"gosystract", "--dumpfile", "../test/single-syscall.dump"})

		actual := output.String()

		should.BeEqual(actual, expected, assumption)
	}

	assertThat("should return exit_group call for single-syscall.dump", "../test/single-syscall.dump", "1 system calls found:\n    exit_group (231)\n")
}
