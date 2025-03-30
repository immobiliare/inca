package util

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

const (
	testingDirPath = "/tmp/.testUtilFilesystem"
)

func TestUtilFilesystemIsDir(t *testing.T) {
	t.Parallel()

	test := is.New(t)

	test.NoErr(os.MkdirAll(testingDirPath, os.ModePerm))
	defer func() {
		if err := os.RemoveAll(testingDirPath); err != nil {
			t.Logf("Failed to clean up temp directory: %v", err)
		}
	}()

	test.True(IsDir(testingDirPath))
}

func TestUtilFilesystemIsNotDir(t *testing.T) {
	t.Parallel()

	is.New(t).True(!IsDir(testingDirPath))
}
