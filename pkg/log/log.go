package log

import (
	"fmt"
	"io"
	"os"
)

// PlainLogger is a logger that doesn't do anything fancy.
type PlainLogger struct {
	w io.Writer
}

// NewPlainLogger creates a new plain logger.
func NewPlainLogger(w io.Writer) *PlainLogger {
	return &PlainLogger{
		w: w,
	}
}

// Printf prints with a format
func (l *PlainLogger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

// Println prints to the logger with a new line at the end.
func (l *PlainLogger) Println(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
}

// Warning prints to the logger as a warning.
func (l *PlainLogger) Warning(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
}

// Error prints to the logger as an error.
func (l *PlainLogger) Error(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
}

// Fatal will output an error and exit the program.
func (l *PlainLogger) Fatal(v ...interface{}) {
	l.Error(v...)
	os.Exit(1)
}

// Output writes the output to the logger.
func (l *PlainLogger) Output(out string) {
	_, err := l.w.Write([]byte(out))
	if err != nil {
		fmt.Println("Logging Error: ", err)
	}
}

// DefaultLogger is the default logging object.
var DefaultLogger PrintLogger = NewColorLogger(os.Stderr, Cyan)

// Println prints to the default logger with a new line at the end.
func Println(v ...interface{}) {
	DefaultLogger.Println(v...)
}

// Printf prints a formatted message to the default logger.
func Printf(format string, v ...interface{}) {
	DefaultLogger.Printf(format, v...)
}

// Warning prints to the default logger as a warning.
func Warning(v ...interface{}) {
	DefaultLogger.Warning(v...)
}

// Error prints an error to the default logger.
func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

// Fatal prints an error to the default logger and exits the program.
func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}
