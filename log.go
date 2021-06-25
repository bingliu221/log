package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var defaultLevel level
var defaultOutputs []io.Writer = []io.Writer{os.Stderr}

type Config func(*Logger)

// Level defines what logs should be print
type level int

// LogLevels
const (
	levelError level = -2 + iota
	levelWarning
	levelInfo
	levelDebug
)

func (lv level) String() string {
	switch lv {
	case levelError:
		return "[E]"
	case levelWarning:
		return "[W]"
	case levelInfo:
		return "[I]"
	case levelDebug:
		return "[D]"
	}
	return ""
}

func setLevel(lv *level, s string) {
	if lv == nil {
		return
	}
	switch s {
	case "debug":
		*lv = levelDebug
	case "info":
		*lv = levelInfo
	case "warning":
		*lv = levelWarning
	case "error":
		*lv = levelError
	}
}

func SetLevel(s string) {
	setLevel(&std.lv, s)
}

func WithLevel(lv string) Config {
	return func(logger *Logger) {
		if logger == nil {
			return
		}
		setLevel(&logger.lv, lv)
	}
}

func SetOutputs(ws ...io.Writer) Config {
	return func(logger *Logger) {
		if logger == nil {
			return
		}
		logger.mu.Lock()
		defer logger.mu.Unlock()
		logger.outputs = ws
	}
}

func SetDefaultOutput(ws ...io.Writer) {
	defaultOutputs = ws
}

// Logger is a simple custom logger support log levels
type Logger struct {
	mu      sync.Mutex
	lv      level
	name    string
	buf     []byte
	outputs []io.Writer
}

// New creates a new Logger
func New(name string, configs ...Config) *Logger {
	logger := &Logger{
		name:    name,
		lv:      defaultLevel,
		outputs: defaultOutputs,
	}
	for _, config := range configs {
		config(logger)
	}
	return logger
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int, lv level) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	milliSec := t.Nanosecond() / 1e6
	ts := fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d.%03d", year, month, day, hour, min, sec, milliSec)
	*buf = append(*buf, ts...)
	*buf = append(*buf, ' ')
	ls := lv.String()
	*buf = append(*buf, ls...)
	if l.name != "" {
		*buf = append(*buf, '[')
		*buf = append(*buf, l.name...)
		*buf = append(*buf, ']', ' ')
	}
	if l.lv == levelDebug {
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		nu := strconv.Itoa(line)
		*buf = append(*buf, nu...)
		*buf = append(*buf, ' ')
	}
}

func (l *Logger) output(lv level, s string) error {
	now := time.Now()
	if lv > l.lv {
		return nil
	}
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.lv == levelDebug {
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line, lv)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	for _, w := range l.outputs {
		_, err := w.Write(l.buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) Error(v ...interface{}) {
	l.output(levelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(levelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.output(levelWarning, fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(levelWarning, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.output(levelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(levelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.output(levelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(levelDebug, fmt.Sprintf(format, v...))
}

var std *Logger = New("")

func Error(v ...interface{}) {
	std.output(levelError, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	std.output(levelError, fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	std.output(levelWarning, fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
	std.output(levelWarning, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	std.output(levelInfo, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	std.output(levelInfo, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	std.output(levelDebug, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	std.output(levelDebug, fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	std.output(levelError, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	std.output(levelError, fmt.Sprintf(format, v...))
	os.Exit(1)
}
