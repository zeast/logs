package logs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

const (
	// LevelDebug  level debug
	LevelDebug = iota

	// LevelInfo level info
	LevelInfo

	// LevelWarn level warn
	LevelWarn

	// LevelError level error
	LevelError

	// LevelFatal level fatal
	LevelFatal
)

var levelMap = map[string]uint32{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

var levelRvsMap = map[uint32]string{
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelFatal: "fatal",
}

var invalidLvErr = errors.New("invalid level")
var defaultTimeLayout = "2006-01-02 15:04:05"

var levelStr = map[uint32][]byte{
	LevelDebug: []byte(" [DEBUG] "),
	LevelInfo:  []byte(" [INFO]  "),
	LevelWarn:  []byte(" [WARN]  "),
	LevelError: []byte(" [ERROR] "),
	LevelFatal: []byte(" [FATAL] "),
}

// Flusher is the interface that warps then Flush method.
type Flusher interface {
	Flush() error
}

// Logger provides fast, leveled, unstructured logging.
// All methods are safe for concurrent use.
type Logger struct {
	sync.Mutex
	level      uint32
	timeLayout string

	writers [][]io.Writer

	funcCallDepth int
	funcCall      bool
}

// LogFuncCall set if log the call stack.
func (l *Logger) LogFuncCall(b bool) {
	l.Lock()
	l.funcCall = b
	if l.funcCallDepth == 0 {
		l.funcCallDepth = 5
	}
	l.Unlock()
}

// LogFuncCallDepth set the depth of call stack.
func (l *Logger) LogFuncCallDepth(depth int) {
	l.Lock()
	l.funcCallDepth = depth
	l.Unlock()
}

func (l *Logger) writef(lv uint32, format string, a ...interface{}) {
	buf := bufPool.Get().(*buffer)
	buf.reset()
	buf.formatTime(l.timeLayout)
	buf.formatLevel(lv)
	if l.funcCall {
		buf.formatCaller(l.funcCallDepth)
	}
	buf.formatMsg(format, a...)
	buf.formatEOF()

	if buf.err != nil {
		fmt.Println(buf.err)
		return
	}

	for _, w := range l.writers[lv] {
		if w == nil {
			continue
		}

		if _, err := w.Write(buf.bytes()); err != nil {
			fmt.Println(err)
		}

		if lv == LevelFatal {
			if f, ok := w.(Flusher); ok {
				f.Flush()
			}
		}
	}

	bufPool.Put(buf)

	if lv == LevelFatal {
		os.Exit(1)
	}
}

func (l *Logger) write(lv uint32, a ...interface{}) {
	buf := bufPool.Get().(*buffer)
	buf.reset()
	buf.formatTime(l.timeLayout)
	buf.formatLevel(lv)
	if l.funcCall {
		buf.formatCaller(l.funcCallDepth)
	}
	buf.formatMsgln(a...)

	if buf.err != nil {
		fmt.Println(buf.err)
		return
	}

	for _, w := range l.writers[lv] {
		if w == nil {
			continue
		}

		if _, err := w.Write(buf.bytes()); err != nil {
			fmt.Println(err)
		}
		if lv == LevelFatal {
			if f, ok := w.(Flusher); ok {
				f.Flush()
			}
		}
	}

	bufPool.Put(buf)

	if lv == LevelFatal {
		os.Exit(1)
	}
}

// LogDebug return if logger will log debug message.
func (l *Logger) LogDebug() bool {
	return atomic.LoadUint32(&l.level) <= LevelDebug
}

// LogInfo return if logger will log Info message.
func (l *Logger) LogInfo() bool {
	return atomic.LoadUint32(&l.level) <= LevelInfo
}

// LogWarn return if logger will log Warn message.
func (l *Logger) LogWarn() bool {
	return atomic.LoadUint32(&l.level) <= LevelWarn
}

// LogError return if logger will log Error message.
func (l *Logger) LogError() bool {
	return atomic.LoadUint32(&l.level) <= LevelError
}

// LogFatal return if logger will log Fatal message.
func (l *Logger) LogFatal() bool {
	return atomic.LoadUint32(&l.level) <= LevelFatal
}

// Debugf logs a message at level Debug.
func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.LogDebug() {
		l.writef(LevelDebug, format, a...)
	}
}

// Infof logs a message at level Info.
func (l *Logger) Infof(format string, a ...interface{}) {
	if l.LogInfo() {
		l.writef(LevelInfo, format, a...)
	}
}

// Warnf logs a message at level Warn.
func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.LogWarn() {
		l.writef(LevelWarn, format, a...)
	}
}

// Errorf logs a message at level Error.
func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.LogError() {
		l.writef(LevelError, format, a...)
	}
}

// Fatalf logs a message at level Fatal.
func (l *Logger) Fatalf(format string, a ...interface{}) {
	if l.LogFatal() {
		l.writef(LevelFatal, format, a...)
	}
}

// Debug logs a message at level Debug.
func (l *Logger) Debug(a ...interface{}) {
	if l.LogDebug() {
		l.write(LevelDebug, a...)
	}
}

// Info logs a message at level Info.
func (l *Logger) Info(a ...interface{}) {
	if l.LogInfo() {
		l.write(LevelInfo, a...)
	}
}

// Warn logs a message at level Warn.
func (l *Logger) Warn(a ...interface{}) {
	if l.LogWarn() {
		l.write(LevelWarn, a...)
	}
}

// Error logs a message at level Error.
func (l *Logger) Error(a ...interface{}) {
	if l.LogError() {
		l.write(LevelError, a...)
	}
}

// Fatal logs a message at level Fatal.
func (l *Logger) Fatal(a ...interface{}) {
	if l.LogFatal() {
		l.write(LevelFatal, a...)
	}
}

// SetLogLevel sets log level of default logger.
func (l *Logger) SetLogLevel(lv uint32) {
	if lv <= LevelFatal {
		atomic.StoreUint32(&l.level, lv)
	}
}

// LogLevel return log level of default logger.
func (l *Logger) LogLevel() uint32 {
	return atomic.LoadUint32(&l.level)
}

// SetLogLevelStr sets log level of logger.
func (l *Logger) SetLogLevelStr(lv string) {
	l.SetLogLevel(levelMap[lv])
}

// LogLevelStr sets log level of logger.
func (l *Logger) LogLevelStr() string {
	return levelRvsMap[l.LogLevel()]
}

// AddWriter add a writer to logger.
func (l *Logger) AddWriter(lv uint32, w io.Writer) error {
	if lv > LevelFatal {
		return invalidLvErr
	}

	l.writers[lv] = append(l.writers[lv], w)
	return nil
}

// SetBaseWriter set base writer of logger.
func (l *Logger) SetBaseWriter(w io.Writer) {
	for _, ws := range l.writers {
		ws[0] = w
	}
}

// SetTimeLayout set time layout of logger.
func (l *Logger) SetTimeLayout(layout string) {
	l.timeLayout = layout
}

// Flush writes buffered data to the underlying writer.
func (l *Logger) Flush() {
	for _, ws := range l.writers {
		for _, w := range ws {
			if f, ok := w.(Flusher); ok {
				f.Flush()
			}
		}
	}
}

// NewLogger return logger with default settings.
func NewLogger(w io.Writer) *Logger {

	l := Logger{
		level:      LevelDebug,
		timeLayout: defaultTimeLayout,
		writers: [][]io.Writer{
			{w},
			{w},
			{w},
			{w},
			{w},
		},
	}

	return &l
}
