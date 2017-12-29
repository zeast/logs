package logs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fileSyncFuncCall = "/tmp/test_sync_func_call.log"
var fileAsync = "/tmp/test_async.log"

func init() {
	_, err := os.Stat(fileSyncFuncCall)
	if err == nil {
		os.Remove(fileSyncFuncCall)
	}

	_, err = os.Stat(fileAsync)
	if err == nil {
		os.Remove(fileSyncFuncCall)
	}

}

func TestSyncFuncCall(t *testing.T) {
	w, err := NewFileWriter(
		FileConfig{
			Name: fileSyncFuncCall,
			Rotate: &Rotate{
				MaxSize: 10,
				Perm:    "0666",
			},
		},
	)

	assert.Nil(t, err)

	SetBaseWriter(w)
	LogFuncCall(true)

	if LogDebug() {
		Debug("debug message")
	}

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	Debugf("%s", "debug message")
	Infof("%s", "info message")
	Warnf("%s", "warn message")
	Errorf("%s", "error message")
}

func BenchmarkAsync(b *testing.B) {
	w, err := NewFileWriter(
		FileConfig{
			Name:  fileAsync,
			Async: true,
		},
	)

	if err != nil {
		panic(err)
	}

	SetBaseWriter(w)

	for i := 0; i < b.N; i++ {
		Debugf("%s", "debug message")
	}

	Flush()
}
