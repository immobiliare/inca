package cmd

import (
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
	test.NoErr(os.WriteFile(tempFile, []byte(config), 0644))
	defer os.Remove(tempFile)

	go func() {
		time.Sleep(3 * time.Second)
		test.NoErr(syscall.Kill(syscall.Getpid(), syscall.SIGINT))
	}()
	_, err := mockExecute(cmdServer,
		"server",
		"--config", tempFile,
	)
	test.NoErr(err)
}
