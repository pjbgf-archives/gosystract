package systract

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pjbgf/go-test/should"
)

func TestDumpReader_GetReader_Integration(t *testing.T) {
	assertThat := func(assumption, filePath string, expectedErr bool) {
		should := should.New(t)
		reader := NewDumpReader(filePath)

		_, err := reader.GetReader()

		hasErrored := err != nil
		should.BeEqual(expectedErr, hasErrored, assumption)
	}

	assertThat("should error when file not found",
		"file-that-dont-exist",
		true)
	assertThat("should be able to process sample dump",
		"../../test/single-syscall.dump",
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
