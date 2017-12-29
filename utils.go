package logs

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"sync"
)

var bufPool = sync.Pool{
	New: NewBuffer,
}

type buffer struct {
	err error
	buf *bytes.Buffer
}

func NewBuffer() interface{} {
	return &buffer{
		buf: bytes.NewBuffer(make([]byte, 1024)),
	}
}

func caller(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)

	return filename, line

}

func (b *buffer) formatTime() {
	if b.err != nil {
		return
	}

	_, b.err = b.buf.Write(nowData())
}

func (b *buffer) formatLevel(lv uint32) {
	if b.err != nil {
		return
	}

	_, b.err = b.buf.Write(levelStr[lv])
}

func (b *buffer) formatCaller(depth int) {
	if b.err != nil {
		return
	}

	filename, line := caller(depth)
	_, b.err = fmt.Fprintf(b.buf, "[%s:%d] ", filename, line)
}

func (b *buffer) formatMsg(format string, a ...interface{}) (err error) {
	if b.err != nil {
		return
	}

	_, b.err = fmt.Fprintf(b.buf, format, a...)
	return
}

func (b *buffer) formatMsgln(a ...interface{}) {
	if b.err != nil {
		return
	}

	_, b.err = fmt.Fprintln(b.buf, a...)
}

func (b *buffer) formatEOF() {
	if b.err != nil {
		return
	}

	_, b.err = b.buf.Write(eof)
}

func (b *buffer) bytes() []byte {
	return b.buf.Bytes()
}

func (b *buffer) reset() {
	b.buf.Reset()
}
