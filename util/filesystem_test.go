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
	defer os.RemoveAll(testingDirPath)

	test.True(IsDir(testingDirPath))
}

func TestUtilFilesystemIsNotDir(t *testing.T) {
	t.Parallel()

	is.New(t).True(!IsDir(testingDirPath))
}
