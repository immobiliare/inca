package cmd

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/matryer/is"
)

const (
	tempFile = "/tmp/.testCmdServer.yml"
	config   = `storage:
    type: fs
    path: ./
`
)

func TestCmdServer(t *testing.T) {
	test := is.New(t)
	test.NoErr(ioutil.WriteFile(tempFile, []byte(config), 0644))
	defer os.Remove(tempFile)

	go func() {
		time.Sleep(3 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	_, err := mockExecute(cmdServer, []string{
		"server",
		"--config", tempFile,
	}...)
	test.NoErr(err)
}
