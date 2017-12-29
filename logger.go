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
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelMap = map[string]uint32{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

var invalidLvErr = errors.New("invalid level")

var levelStr = map[uint32][]byte{
	LevelDebug: []byte(" [DEBUG] "),
	LevelInfo:  []byte(" [INFO]  "),
	LevelWarn:  []byte(" [WARN]  "),
	LevelError: []byte(" [ERROR] "),
	LevelFatal: []byte(" [FATAL] "),
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	sync.Mutex
	level uint32

	writers [][]io.Writer

	funcCallDepth int
	funcCall      bool
}

func (l *Logger) LogFuncCall(b bool) {
	l.Lock()
	l.funcCall = b
	if l.funcCallDepth == 0 {
		l.funcCallDepth = 5
	}
	l.Unlock()
}

func (l *Logger) LogFuncCallDepth(depth int) {
	l.Lock()
	l.funcCallDepth = depth
	l.Unlock()
}

func (l *Logger) writef(lv uint32, format string, a ...interface{}) {
	buf := bufPool.Get().(*buffer)
	buf.reset()
	buf.formatTime()
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
		os.Exit(0)
	}
}

func (l *Logger) write(lv uint32, a ...interface{}) {
	buf := bufPool.Get().(*buffer)
	buf.reset()
	buf.formatTime()
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
		os.Exit(0)
	}
}

func (l *Logger) LogDebug() bool {
	return atomic.LoadUint32(&l.level) <= LevelDebug
}

func (l *Logger) LogInfo() bool {
	return atomic.LoadUint32(&l.level) <= LevelInfo
}

func (l *Logger) LogWarn() bool {
	return atomic.LoadUint32(&l.level) <= LevelWarn
}

func (l *Logger) LogError() bool {
	return atomic.LoadUint32(&l.level) <= LevelError
}

func (l *Logger) LogFatal() bool {
	return atomic.LoadUint32(&l.level) <= LevelFatal
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.LogDebug() {
		l.writef(LevelDebug, format, a...)
	}
}

func (l *Logger) Infof(format string, a ...interface{}) {
	if l.LogInfo() {
		l.writef(LevelInfo, format, a...)
	}
}

func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.LogWarn() {
		l.writef(LevelWarn, format, a...)
	}
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.LogError() {
		l.writef(LevelError, format, a...)
	}
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	if l.LogFatal() {
		l.writef(LevelFatal, format, a...)
	}
}

func (l *Logger) Debug(a ...interface{}) {
	if l.LogDebug() {
		l.write(LevelDebug, a...)
	}
}

func (l *Logger) Info(a ...interface{}) {
	if l.LogInfo() {
		l.write(LevelInfo, a...)
	}
}

func (l *Logger) Warn(a ...interface{}) {
	if l.LogWarn() {
		l.write(LevelWarn, a...)
	}
}

func (l *Logger) Error(a ...interface{}) {
	if l.LogError() {
		l.write(LevelError, a...)
	}
}

func (l *Logger) Fatal(a ...interface{}) {
	if l.LogFatal() {
		l.write(LevelFatal, a...)
	}
}

func (l *Logger) SetLogLevel(lv uint32) {
	if lv <= LevelFatal {
		atomic.StoreUint32(&l.level, lv)
	}
}

func (l *Logger) SetLogLevelStr(lv string) {
	l.SetLogLevel(levelMap[lv])
}

func (l *Logger) AddWriter(lv uint32, w io.Writer) error {
	if lv > LevelFatal {
		return invalidLvErr
	}

	l.writers[LevelError] = append(l.writers[LevelError], w)
	return nil
}

func (l *Logger) SetBaseWriter(w io.Writer) {
	for _, ws := range l.writers {
		ws[0] = w
	}
}

func (l *Logger) Flush() {
	for _, ws := range l.writers {
		for _, w := range ws {
			if f, ok := w.(Flusher); ok {
				f.Flush()
			}
		}
	}
}

func NewLogger(w io.Writer) *Logger {

	l := Logger{
		level: LevelDebug,
		writers: [][]io.Writer{
			[]io.Writer{w},
			[]io.Writer{w},
			[]io.Writer{w},
			[]io.Writer{w},
			[]io.Writer{w},
		},
	}

	return &l
}
