package systract

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/pjbgf/go-test/should"
)

func TestExeReader_GetReader_Integration(t *testing.T) {
	assertThat := func(assumption, filePath string, expectedErr bool) {
		should := should.New(t)
		reader := NewExeReader(filePath)

		_, err := reader.GetReader()

		hasErrored := err != nil
		should.BeEqual(expectedErr, hasErrored, assumption)
	}

	assertThat("should error when file not found",
		"file-that-dont-exist",
		true)

	assertThat("should be able disassemble current executable",
		"../../test/simple-app",
		false)

	// test the handling of current chdir being deleted
	wdSnapshot, _ := os.Getwd()
	tmpFolder, err := ioutil.TempDir("", "zaz-test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	os.Chdir(tmpFolder)
	os.Remove(tmpFolder)

	assertThat("should error if current directly disappears",
		"any-file", true)

	// returns snapshotted working directory to ensure other tests' repeatability
	os.Chdir(wdSnapshot)
}

func TestGetFileDumpReader(t *testing.T) {
	assertThat := func(assumption, objDumpPath, input, expectedPrefix string,
		expectedErr error) {
		should := should.New(t)

		reader, err := getFileDumpReader(objDumpPath, input)
		scanner := bufio.NewScanner(reader)
		scanner.Scan()
		actual := scanner.Text()
		reader.Close()

		hasPrefix := strings.HasPrefix(actual, expectedPrefix)

		should.BeTrue(hasPrefix, assumption)
		should.BeEqual(expectedErr, err, assumption)
	}

	assertThat("should support custom objDump path", "/bin/echo", "123456", "123456", nil)
	assertThat("should fallback to default if path does not exist", "/bin/echo1", "../../test/simple-app", "TEXT internal/cpu.Initialize(SB)", nil)
}
