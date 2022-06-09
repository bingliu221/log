package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type option func(*Logger)

// Level defines what logs should be print
type Level int

// LogLevels
const (
	LevelError   Level = -2
	LevelWarning Level = -1
	LevelInfo    Level = 0 // default log level
	LevelDebug   Level = 1
)

var levelString = map[Level]string{
	LevelError:   "[E]",
	LevelWarning: "[W]",
	LevelInfo:    "[I]",
	LevelDebug:   "[D]",
}

func (lv Level) String() string {
	return levelString[lv]
}

func WithLevel(lv Level) option {
	return func(l *Logger) {
		if l == nil {
			return
		}
		l.setLevel(lv)
	}
}

func WithOutput(w io.Writer) option {
	return func(logger *Logger) {
		if logger == nil {
			return
		}
		logger.out = w
	}
}

// Logger is a simple custom logger support log levels
type Logger struct {
	mu   sync.Mutex
	lv   Level
	name string
	out  io.Writer

	bufPool sync.Pool
}

// New creates a new Logger
func New(name string, opts ...option) *Logger {
	logger := &Logger{
		name: name,
		out:  os.Stderr,

		bufPool: sync.Pool{
			New: func() any {
				buf := make([]byte, 0, 1024)
				return &buf
			},
		},
	}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

func (l *Logger) setLevel(lv Level) {
	l.lv = lv
}

func (l *Logger) setLevelString(s string) {
	switch strings.ToLower(s) {
	case "error", "e":
		l.setLevel(LevelError)
	case "warning", "w":
		l.setLevel(LevelWarning)
	case "info", "i":
		l.setLevel(LevelInfo)
	case "debug", "d":
		l.setLevel(LevelDebug)
	}
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int, lv Level) {
	ts := t.Format("2006/01/02 15:04:05.000 ")
	*buf = append(*buf, ts...)

	ls := lv.String()
	*buf = append(*buf, ls...)

	if l.name != "" {
		*buf = append(*buf, '[')
		*buf = append(*buf, l.name...)
		*buf = append(*buf, ']', ' ')
	}

	if l.lv == LevelDebug {
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		nu := strconv.Itoa(line)
		*buf = append(*buf, nu...)
		*buf = append(*buf, ' ')
	}
}

func (l *Logger) output(lv Level, s string) {
	now := time.Now()
	if lv > l.lv {
		return
	}

	var file string
	var line int
	if l.lv == LevelDebug {
		var ok bool
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
	}

	buf := l.bufPool.New().(*[]byte)
	defer l.bufPool.Put(buf)

	*buf = (*buf)[:0]
	l.formatHeader(buf, now, file, line, lv)
	*buf = append(*buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		*buf = append(*buf, '\n')
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := l.out.Write(*buf)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) Error(v ...interface{}) {
	l.output(LevelError, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(LevelError, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.output(LevelWarning, fmt.Sprint(v...))
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.output(LevelInfo, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.output(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(LevelDebug, fmt.Sprintf(format, v...))
}
