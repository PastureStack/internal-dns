package logging

import (
	"io"
	"log"
	"strings"
	"sync/atomic"
)

var debug atomic.Bool

func SetLevelString(level string) {
	debug.Store(strings.EqualFold(level, "debug"))
}

func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

func Debug(args ...any) {
	if debug.Load() {
		log.Print(args...)
	}
}

func Debugf(format string, args ...any) {
	if debug.Load() {
		log.Printf(format, args...)
	}
}

func Info(args ...any)                  { log.Print(args...) }
func Infof(format string, args ...any)  { log.Printf(format, args...) }
func Warn(args ...any)                  { log.Print(args...) }
func Warnf(format string, args ...any)  { log.Printf(format, args...) }
func Errorf(format string, args ...any) { log.Printf(format, args...) }
func Fatal(args ...any)                 { log.Fatal(args...) }
func Fatalf(format string, args ...any) { log.Fatalf(format, args...) }
