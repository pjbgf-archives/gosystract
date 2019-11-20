package systract

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pjbgf/go-test/should"
)

func TestFileExists(t *testing.T) {
	assertThat := func(assumption, filePath string, expected bool) {
		should := should.New(t)

		exists := fileExists(filePath)

		should.BeEqual(expected, exists, assumption)
	}

	currentFile, _ := os.Executable()
	assertThat("should return true for existing files", currentFile, true)
	assertThat("should return false when file cannot be found", "somefilethatdoesnotexist", false)
}

func TestSanitiseFileName(t *testing.T) {
	// creates tmp folder to make the expected filepath more predictable
	tmpFolder, err := ioutil.TempDir("", "zaz-test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.Chdir(tmpFolder)

	assertThat := func(assumption, filePath, expected string, expectedErr error) {
		should := should.New(t)

		fileName, err := sanitiseFileName(filePath)

		should.BeEqual(expectedErr, err, assumption)
		should.BeEqual(expected, fileName, assumption)
	}

	assertThat("should concate single dot relative paths with current dir",
		"./filename", fmt.Sprintf("%s/%s", tmpFolder, "filename"), nil)
	assertThat("should concate filenames with current dir", "filename",
		fmt.Sprintf("%s/%s", tmpFolder, "filename"), nil)
	assertThat("should do nothing for absolute paths", "/root/filename",
		"/root/filename", nil)
	assertThat("should remove .. from rooted paths", "/../etc/passwd",
		"/etc/passwd", nil)

	os.Remove(tmpFolder)

	assertThat("should error if can't get current folder", "", "",
		errors.New("error getting current folder"))
}
