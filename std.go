package log

import (
	"fmt"
	"os"
)

var std *Logger

func init() {
	std = New("")
}

func SetName(name string) {
	std.name = name
}

func SetOutput(out io.Writer) {
	std.out = out
}

func SetLevelString(s string) {
	std.setLevelString(s)
}

func SetLevel(lv Level) {
	std.setLevel(lv)
}

func Error(v ...interface{}) {
	std.output(LevelError, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	std.output(LevelError, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	std.output(LevelWarning, fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	std.output(LevelWarning, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	std.output(LevelInfo, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	std.output(LevelInfo, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	std.output(LevelDebug, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	std.output(LevelDebug, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	std.output(LevelError, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	std.output(LevelError, fmt.Sprintf(format, v...))
	os.Exit(1)
}
