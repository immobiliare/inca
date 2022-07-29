package util

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type ZStdLogger struct {
	z zerolog.Logger
}

func NewZStdLogger() ZStdLogger {
	return ZStdLogger{zerolog.New(os.Stdout).With().Timestamp().Logger()}
}

func (l ZStdLogger) Fatal(args ...interface{}) {
	l.Fatalf("%v\n", args...)
}

func (l ZStdLogger) Fatalln(args ...interface{}) {
	l.Fatalf("%v\n", args...)
}

func (l ZStdLogger) Fatalf(format string, args ...interface{}) {
	l.z.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (l ZStdLogger) Print(args ...interface{}) {
	l.Printf("%v\n", args...)
}

func (l ZStdLogger) Println(args ...interface{}) {
	l.Printf("%v\n", args...)
}

func (l ZStdLogger) Printf(format string, args ...interface{}) {
	l.z.Info().Msg(fmt.Sprintf(format, args...))
}
