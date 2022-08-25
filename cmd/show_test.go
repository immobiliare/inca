package cmd

import (
	"testing"

	"github.com/matryer/is"
)

const (
	tempDir = "/tmp"
	tempCrt = tempDir + "/crt.pem"
	tempKey = tempDir + "/key.pem"
)

func TestCmdShow(t *testing.T) {
	test := is.New(t)
	_, err := mockExecute(cmdGen, []string{
		"gen",
		"--ca",
		"--name", name,
		"--output", tempDir,
	}...)
	test.NoErr(err)

	_, err = mockExecute(cmdShow, []string{"show", tempCrt, tempKey}...)
	test.NoErr(err)
}
