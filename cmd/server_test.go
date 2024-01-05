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
	testingConfig     = `templates_path: ../server/views
storage:
    type: fs
    path: /tmp
`
)

func TestCmdServer(t *testing.T) {
	test := is.New(t)
	test.NoErr(os.WriteFile(testingConfigPath, []byte(testingConfig), 0644))
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
