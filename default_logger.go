package logs

import (
	"io"
	"os"
)

var _logger = NewLogger(os.Stdout)

func LogDebug() bool {
	return _logger.LogDebug()
}

func LogInfo() bool {
	return _logger.LogInfo()
}

func LogWarn() bool {
	return _logger.LogWarn()
}

func LogError() bool {
	return _logger.LogError()
}

func LogFatal() bool {
	return _logger.LogFatal()
}

func Debugf(format string, a ...interface{}) {
	_logger.Debugf(format, a...)
}

func Infof(format string, a ...interface{}) {
	_logger.Infof(format, a...)
}

func Warnf(format string, a ...interface{}) {
	_logger.Warnf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	_logger.Errorf(format, a...)
}

func Fatalf(format string, a ...interface{}) {
	_logger.Fatalf(format, a...)
}

func Debug(a ...interface{}) {
	_logger.Debug(a...)
}

func Info(a ...interface{}) {
	_logger.Info(a...)
}

func Warn(a ...interface{}) {
	_logger.Warn(a...)
}

func Error(a ...interface{}) {
	_logger.Error(a...)
}

func Fatal(a ...interface{}) {
	_logger.Fatal(a...)
}

func SetLogLevel(lv uint32) {
	_logger.SetLogLevel(lv)
}

func LogLevel() uint32 {
	return _logger.LogLevel()
}

func SetLogLevelStr(lv string) {
	_logger.SetLogLevelStr(lv)
}

func LogLevelStr() string {
	return _logger.LogLevelStr()
}

func LogFuncCall(b bool) {
	_logger.LogFuncCall(b)
}

func LogFuncCallDepth(depth int) {
	_logger.LogFuncCallDepth(depth)
}

func AddWriter(lv uint32, w io.Writer) error {
	return _logger.AddWriter(lv, w)
}

func SetBaseWriter(w io.Writer) {
	_logger.SetBaseWriter(w)
}

func Flush() {
	_logger.Flush()
}
