package logging

import (
	"fmt"
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
		log.Print(safeLogMessage(fmt.Sprint(args...)))
	}
}

func Debugf(format string, args ...any) {
	if debug.Load() {
		log.Print(safeLogMessage(fmt.Sprintf(format, args...)))
	}
}

func Info(args ...any)                  { log.Print(safeLogMessage(fmt.Sprint(args...))) }
func Infof(format string, args ...any)  { log.Print(safeLogMessage(fmt.Sprintf(format, args...))) }
func Warn(args ...any)                  { log.Print(safeLogMessage(fmt.Sprint(args...))) }
func Warnf(format string, args ...any)  { log.Print(safeLogMessage(fmt.Sprintf(format, args...))) }
func Errorf(format string, args ...any) { log.Print(safeLogMessage(fmt.Sprintf(format, args...))) }
func Fatal(args ...any)                 { log.Fatal(safeLogMessage(fmt.Sprint(args...))) }
func Fatalf(format string, args ...any) { log.Fatal(safeLogMessage(fmt.Sprintf(format, args...))) }

func safeLogMessage(message string) string {
	message = strings.ReplaceAll(message, "\r", "")
	return strings.ReplaceAll(message, "\n", " ")
}
