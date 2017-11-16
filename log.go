// Package log implements a simple logger. It is fully compatible with the standard log package
// and, it provides extra functionality like labels, log levels and colored output.
package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	golog "log"
)

// These flags define which text to prefix to each log entry generated by the Logger.
const (
	// Bits or'ed together to control what's printed.
	// There is no control over the order they appear (the order listed
	// here) or the format they present (as described in the comments).
	// The prefix is followed by a colon only when Llongfile or Lshortfile
	// is specified.
	// For example, flags Ldate | Ltime (or LstdFlags) produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Llabel                        // log entry label: [DEBUG], [ERROR], [PANIC], ...
	Lcolor                        // colored output
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// Log levels.
const (
	LevelFatal   = iota       // Fatal log level
	LevelPanic                // Panic log level
	LevelError                // Error log level
	LevelWarn                 // Warning log level
	LevelInfo                 // Information log level
	LevelDebug                // Debug log level
	LevelDefault = LevelError // Default log level (LevelError)
)

// For colored output.
const (
	// ANSI colors
	colorNone   = 0
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorBlue   = 36
	colorWhite  = 37

	// ANSI escape sequence format
	escSeq = "\033[%dm"
)

var (
	// labelMap contains labels mapped to log levels.
	labelMap = []string{
		"FATAL",
		"PANIC",
		"ERROR",
		"WARN ",
		"INFO ",
		"DEBUG",
	}

	// colorMap contains ANSI colors mapped to log levels.
	colorMap = []int{
		colorRed,
		colorRed,
		colorRed,
		colorYellow,
		colorBlue,
		colorGreen,
	}
)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	l     *golog.Logger
	mu    sync.Mutex
	flag  int
	level int
}

// New returns a new Logger.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		l:     golog.New(out, prefix, flag),
		flag:  flag,
		level: LevelDefault,
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.l.SetOutput(w)
}

func (l *Logger) format(level int, s string) {
	if l.flag&Llabel != 0 {
		label := labelMap[level]

		if l.flag&Lcolor != 0 {
			color := colorMap[level]
			s = fmt.Sprintf("["+escSeq+"%s"+escSeq+"] "+escSeq+"%s"+escSeq, color, label, colorNone, colorWhite, s, colorNone)
		} else {
			s = fmt.Sprintf("[%s] %s", label, s)
		}
	}

	l.l.Output(3, s)
}

func (l *Logger) Output(calldepth int, s string) error {
	return l.l.Output(calldepth+1, s)
}

func (l *Logger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelInfo {
		l.format(LevelInfo, fmt.Sprint(v...))
	}
}

func (l *Logger) Println(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelInfo {
		l.format(LevelInfo, fmt.Sprintln(v...))
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelInfo {
		l.format(LevelInfo, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelFatal {
		l.format(LevelFatal, fmt.Sprint(v...))
	}
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelFatal {
		l.format(LevelFatal, fmt.Sprintln(v...))
	}
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelFatal {
		l.format(LevelFatal, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	s := fmt.Sprint(v...)
	if l.level >= LevelPanic {
		l.format(LevelPanic, s)
	}
	panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	s := fmt.Sprintln(v...)
	if l.level >= LevelPanic {
		l.format(LevelPanic, s)
	}
	panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	s := fmt.Sprintf(format, v...)
	if l.level >= LevelPanic {
		l.format(LevelPanic, s)
	}
	panic(s)
}

func (l *Logger) Error(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelError {
		l.format(LevelError, fmt.Sprint(v...))
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelError {
		l.format(LevelError, fmt.Sprintln(v...))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelError {
		l.format(LevelError, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelWarn {
		l.format(LevelWarn, fmt.Sprint(v...))
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelWarn {
		l.format(LevelWarn, fmt.Sprintln(v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelWarn {
		l.format(LevelWarn, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelInfo {
		l.format(LevelInfo, fmt.Sprint(v...))
	}
}

func (l *Logger) Infoln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelInfo {
		l.format(LevelInfo, fmt.Sprintln(v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelInfo {
		l.format(LevelInfo, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelDebug {
		l.format(LevelDebug, fmt.Sprint(v...))
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelDebug {
		l.format(LevelDebug, fmt.Sprintln(v...))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelDebug {
		l.format(LevelDebug, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Flags() (v int) {
	l.mu.Lock()
	v = l.flag
	l.mu.Unlock()
	return
}

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
	l.l.SetFlags(flag)
}

func (l *Logger) Level() (v int) {
	l.mu.Lock()
	v = l.level
	l.mu.Unlock()
	return
}

func (l *Logger) SetLevel(level int) {
	if level > LevelDebug {
		panic("invalid log level")
	}
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
}

func (l *Logger) Prefix() string {
	return l.l.Prefix()
}

func (l *Logger) SetPrefix(prefix string) {
	l.l.SetPrefix(prefix)
}

// Standard logger
var std = New(os.Stderr, "", LstdFlags)

func StdLogger() *Logger {
	return std
}

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func Output(calldepth int, s string) error {
	return std.l.Output(calldepth+1, s)
}

func Print(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelInfo {
		std.format(LevelInfo, fmt.Sprint(v...))
	}
}

func Println(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelInfo {
		std.format(LevelInfo, fmt.Sprintln(v...))
	}
}

func Printf(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelInfo {
		std.format(LevelInfo, fmt.Sprintf(format, v...))
	}
}

func Fatal(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelFatal {
		std.format(LevelFatal, fmt.Sprint(v...))
	}
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelFatal {
		std.format(LevelFatal, fmt.Sprintln(v...))
	}
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelFatal {
		std.format(LevelFatal, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func Panic(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	s := fmt.Sprint(v...)
	if std.level >= LevelPanic {
		std.format(LevelPanic, s)
	}
	panic(s)
}

func Panicln(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	s := fmt.Sprintln(v...)
	if std.level >= LevelPanic {
		std.format(LevelPanic, s)
	}
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	s := fmt.Sprintf(format, v...)
	if std.level >= LevelPanic {
		std.format(LevelPanic, s)
	}
	panic(s)
}

func Error(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelError {
		std.format(LevelError, fmt.Sprint(v...))
	}
}

func Errorln(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelError {
		std.format(LevelError, fmt.Sprintln(v...))
	}
}

func Errorf(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelError {
		std.format(LevelError, fmt.Sprintf(format, v...))
	}
}

func Warn(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelWarn {
		std.format(LevelWarn, fmt.Sprint(v...))
	}
}

func Warnln(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelWarn {
		std.format(LevelWarn, fmt.Sprintln(v...))
	}
}

func Warnf(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelWarn {
		std.format(LevelWarn, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelInfo {
		std.format(LevelInfo, fmt.Sprint(v...))
	}
}

func Infoln(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelInfo {
		std.format(LevelInfo, fmt.Sprintln(v...))
	}
}

func Infof(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelInfo {
		std.format(LevelInfo, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelDebug {
		std.format(LevelDebug, fmt.Sprint(v...))
	}
}

func Debugln(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelDebug {
		std.format(LevelDebug, fmt.Sprintln(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()
	if std.level >= LevelDebug {
		std.format(LevelDebug, fmt.Sprintf(format, v...))
	}
}

func Flags() int {
	return std.Flags()
}

func SetFlags(flag int) {
	std.SetFlags(flag)
}

func Level() int {
	return std.Level()
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func Prefix() string {
	return std.Prefix()
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}
