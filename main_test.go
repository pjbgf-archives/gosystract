package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_E2E(t *testing.T) {
	assert := assert.New(t)

	var output bytes.Buffer
	run(&output, []string{"gosystract", "test/single-syscall.dump"})

	actual := output.String()
	expected := "1 system calls found:\n    exit_group (231)"

	assert.Contains(actual, expected)
}
