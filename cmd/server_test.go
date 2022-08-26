package cmd

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/matryer/is"
)

const (
	testingConfigPath = "/tmp/.testCmdServer.yml"
	tesTingConfig     = `storage:
    type: fs
    path: ./
`
)

func TestCmdServer(t *testing.T) {
	test := is.New(t)
	test.NoErr(os.WriteFile(testingConfigPath, []byte(tesTingConfig), 0644))
	defer os.Remove(testingConfigPath)

	go func() {
		time.Sleep(3 * time.Second)
		test.NoErr(syscall.Kill(syscall.Getpid(), syscall.SIGINT))
	}()
	_, err := testExecute(cmdServer,
		"server",
		"--config", testingConfigPath,
	)
	test.NoErr(err)
}
