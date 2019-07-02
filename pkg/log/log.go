package log

import "os"

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
