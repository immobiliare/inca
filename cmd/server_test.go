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
	t.Parallel()

	test := is.New(t)
	test.NoErr(os.WriteFile(testingConfigPath, []byte(testingConfig), 0644))
	defer func() {
		if err := os.Remove(testingConfigPath); err != nil {
			t.Logf("Failed to remove test config file: %v", err)
		}
	}()

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
