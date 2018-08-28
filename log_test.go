package logs

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"time"
)

var TestFile = "/tmp/test.log"
var BenchAsyncFile = "/tmp/test_async.log"
var BenchSyncFile = "/tmp/test_sync.log"

func init() {
	for _, f := range []string{TestFile, BenchAsyncFile, BenchSyncFile} {
		_, err := os.Stat(f)
		if err == nil {
			os.Remove(f)
		}
	}

	t, err := time.Parse(defaultTimeLayout, "2018-08-28 11:09:00")
	if err != nil {
		panic(err)
	}

	now = func() time.Time {
		return t
	}
}

func TestSyncFuncCall(t *testing.T) {
	w, err := NewFileWriter(
		FileConfig{
			Name: TestFile,
		},
	)

	assert.Nil(t, err)

	SetBaseWriter(w)

	LogFuncCall(true)

	if LogDebug() {
		Debug("debug message")
	}

	if LogInfo() {
		Info("info message")
	}

	if LogWarn() {
		Warn("warn message")
	}

	if LogError() {
		Error("error message")
	}

	Debugf("%s", "debug message")
	Infof("%s", "info message")
	Warnf("%s", "warn message")
	Errorf("%s", "error message")

	file, err := os.Open(TestFile)
	assert.Nil(t, err)

	b, err := ioutil.ReadAll(file)
	assert.Nil(t, err)

	s := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	assert.Equal(t, len(s), 8)
	assert.Equal(t, s[0], "2018-08-28 11:09:00 [DEBUG] [log_test.go:50] debug message")
	assert.Equal(t, s[1], "2018-08-28 11:09:00 [INFO]  [log_test.go:54] info message")
	assert.Equal(t, s[2], "2018-08-28 11:09:00 [WARN]  [log_test.go:58] warn message")
	assert.Equal(t, s[3], "2018-08-28 11:09:00 [ERROR] [log_test.go:62] error message")
	assert.Equal(t, s[4], "2018-08-28 11:09:00 [DEBUG] [log_test.go:65] debug message")
	assert.Equal(t, s[5], "2018-08-28 11:09:00 [INFO]  [log_test.go:66] info message")
	assert.Equal(t, s[6], "2018-08-28 11:09:00 [WARN]  [log_test.go:67] warn message")
	assert.Equal(t, s[7], "2018-08-28 11:09:00 [ERROR] [log_test.go:68] error message")

}

func BenchmarkAsync(b *testing.B) {
	w, err := NewFileWriter(
		FileConfig{
			Name:  BenchAsyncFile,
			Async: true,
		},
	)

	if err != nil {
		panic(err)
	}

	SetBaseWriter(w)

	for i := 0; i < b.N; i++ {
		do()
	}

	Flush()
}

func BenchmarkSync(b *testing.B) {
	w, err := NewFileWriter(
		FileConfig{
			Name:  BenchSyncFile,
			Async: false,
		},
	)

	if err != nil {
		panic(err)
	}

	SetBaseWriter(w)

	for i := 0; i < b.N; i++ {
		Debugf("%s", "debug message")
	}

}

func do() {
	var wg = new(sync.WaitGroup)
	var num = 1000
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			Debugf("%s", "debug message")
		}()
	}

	wg.Wait()
}
