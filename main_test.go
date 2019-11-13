package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_E2E(t *testing.T) {
	assert := assert.New(t)

	fileName := "test/single-syscall.dump"
	cmd := exec.Command("go",
		"run",
		"main.go",
		fileName)

	output, err := cmd.Output()

	actual := string(output)
	expected := "1 system calls found:\n    exit_group (231)"

	assert.Nil(err)
	assert.Contains(actual, expected)
}
