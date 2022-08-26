package util

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type ZStdLogger struct {
	z zerolog.Logger
}

func NewZStdLogger(fd ...io.Writer) ZStdLogger {
	if len(fd) == 0 {
		fd = append(fd, os.Stdout)
	}
	return ZStdLogger{zerolog.New(fd[0]).With().Timestamp().Logger()}
}

func (l ZStdLogger) Fatal(args ...interface{}) {
	l.Fatalf("%v", args...)
}

func (l ZStdLogger) Fatalln(args ...interface{}) {
	l.Fatalf("%v", args...)
}

func (l ZStdLogger) Fatalf(format string, args ...interface{}) {
	l.z.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (l ZStdLogger) Print(args ...interface{}) {
	l.Printf("%v", args...)
}

func (l ZStdLogger) Println(args ...interface{}) {
	l.Printf("%v", args...)
}

func (l ZStdLogger) Printf(format string, args ...interface{}) {
	l.z.Info().Msg(fmt.Sprintf(format, args...))
}
