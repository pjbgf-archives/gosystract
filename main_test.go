package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsage_NoArgs(t *testing.T) {
	assert := assert.New(t)

	cmd := exec.Command("go",
		"run",
		"main.go")

	output, err := cmd.CombinedOutput()

	actual := string(output)

	assert.Error(err, "when no args provided should error")
	assert.Contains(
		actual,
		usageMessage,
		"when no args provided should show usage")
}

func TestMain_NoSyscalls(t *testing.T) {
	assert := assert.New(t)

	fileName := "test/no-syscalls.dump"
	cmd := exec.Command("go",
		"run",
		"main.go",
		fileName)

	output, err := cmd.Output()

	actual := string(output)
	expected := noSyscallsFoundMessage

	assert.Nil(err)
	assert.Contains(actual, expected)
}

func TestMain(t *testing.T) {
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
