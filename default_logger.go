package logs

import (
	"io"
	"os"
)

var _logger = NewLogger(os.Stdout)

// LogDebug return if default logger will log debug message.
func LogDebug() bool {
	return _logger.LogDebug()
}

// LogInfo return if default logger will log Info message.
func LogInfo() bool {
	return _logger.LogInfo()
}

// LogWarn return if default logger will log warn message.
func LogWarn() bool {
	return _logger.LogWarn()
}

// LogError return if default logger will log error message.
func LogError() bool {
	return _logger.LogError()
}

// LogFatal return if default logger will log fatal message.
func LogFatal() bool {
	return _logger.LogFatal()
}

// Debugf logs a message at level Debug on the default logger.
func Debugf(format string, a ...interface{}) {
	_logger.Debugf(format, a...)
}

// Infof logs a message at level Info on the default logger.
func Infof(format string, a ...interface{}) {
	_logger.Infof(format, a...)
}

// Warnf logs a message at level Warn on the default logger.
func Warnf(format string, a ...interface{}) {
	_logger.Warnf(format, a...)
}

// Errorf logs a message at level Error on the default logger.
func Errorf(format string, a ...interface{}) {
	_logger.Errorf(format, a...)
}

// Fatalf logs a message at level Fatal on the default logger.
func Fatalf(format string, a ...interface{}) {
	_logger.Fatalf(format, a...)
}

// Debug logs a message at level Debug on the default logger.
func Debug(a ...interface{}) {
	_logger.Debug(a...)
}

// Info logs a message at level Info on the default logger.
func Info(a ...interface{}) {
	_logger.Info(a...)
}

// Warn logs a message at level Warn on the default logger.
func Warn(a ...interface{}) {
	_logger.Warn(a...)
}

// Error logs a message at level Error on the default logger.
func Error(a ...interface{}) {
	_logger.Error(a...)
}

// Fatal logs a message at level Fatal on the default logger.
func Fatal(a ...interface{}) {
	_logger.Fatal(a...)
}

// SetLogLevel sets log level of default logger.
func SetLogLevel(lv uint32) {
	_logger.SetLogLevel(lv)
}

// LogLevel return log level of default logger.
func LogLevel() uint32 {
	return _logger.LogLevel()
}

// SetLogLevelStr sets log level of default logger.
func SetLogLevelStr(lv string) {
	_logger.SetLogLevelStr(lv)
}

// LogLevelStr sets log level of default logger.
func LogLevelStr() string {
	return _logger.LogLevelStr()
}

// LogFuncCall set if log the call stack.
func LogFuncCall(b bool) {
	_logger.LogFuncCall(b)
}

// LogFuncCallDepth set the depth of call stack.
func LogFuncCallDepth(depth int) {
	_logger.LogFuncCallDepth(depth)
}

// AddWriter add a writer to default logger.
func AddWriter(lv uint32, w io.Writer) error {
	return _logger.AddWriter(lv, w)
}

// SetBaseWriter set base writer of default logger.
func SetBaseWriter(w io.Writer) {
	_logger.SetBaseWriter(w)
}

// SetTimeLayout set time layout of default logger.
func SetTimeLayout(layout string) {
	_logger.SetTimeLayout(layout)
}

// Flush writes buffered data to the underlying writer.
func Flush() {
	_logger.Flush()
}

// DefaultLogger return the default logger.
func DefaultLogger() *Logger {
	return _logger;
}