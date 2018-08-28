package logs

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var MinBufSize = 4096
var eof = []byte{10}

type FileWriter struct {
	sync.Mutex
	writer        *bufio.Writer
	file          *os.File
	openTime      time.Time
	curSize       int
	flushInterval time.Duration
	cfg           FileConfig
}

type FileConfig struct {
	Name    string
	Async   bool
	Perm    string
	perm    os.FileMode
	BufSize int
	Rotate  *Rotate
}

type Rotate struct {
	MaxSize int
	MaxDays int
	Perm    string
	Daily   bool
	perm    os.FileMode
}

func NewFileWriter(cfg FileConfig) (*FileWriter, error) {
	var w = new(FileWriter)
	w.cfg = cfg

	if w.cfg.BufSize < MinBufSize {
		w.cfg.BufSize = MinBufSize
	}

	if w.cfg.Perm == "" {
		w.cfg.Perm = "0644"
	}

	perm, err := strconv.ParseUint(w.cfg.Perm, 8, 32)
	if err != nil {
		return nil, fmt.Errorf("perm error. %s", err)
	}
	w.cfg.perm = os.FileMode(perm)

	if w.cfg.Rotate != nil {
		if w.cfg.Rotate.Perm == "" {
			w.cfg.Rotate.Perm = w.cfg.Perm
		}

		perm, err := strconv.ParseUint(w.cfg.Rotate.Perm, 8, 32)
		if err != nil {
			return nil, fmt.Errorf("rotate perm error. %s", err)
		}
		w.cfg.Rotate.perm = os.FileMode(perm)
	}

	if err := w.initFile(); err != nil {
		return nil, err
	}

	if w.cfg.Async {
		w.flushInterval = time.Second
		go w.autoFlush()
	}

	return w, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	if w.needRotate() {
		if w.cfg.Async {
			w.writer.Flush()
		}

		w.doRotate()
	}

	if w.cfg.Async {
		n, err = w.writer.Write(p)
	} else {
		n, err = w.file.Write(p)
	}

	w.curSize += n

	w.Unlock()

	return
}

func (w *FileWriter) Flush() (err error) {
	if !w.cfg.Async {
		return nil
	}

	w.Lock()
	err = w.writer.Flush()
	w.Unlock()
	return
}

func (w *FileWriter) needRotate() bool {
	if w.cfg.Rotate == nil {
		return false
	}

	if (w.cfg.Rotate.Daily && now().Day() != w.openTime.Day()) ||
		(w.cfg.Rotate.MaxSize > 0 && w.curSize >= w.cfg.Rotate.MaxSize) {
		return true
	}

	return false
}

func (w *FileWriter) doRotate() {
	var err error
	var fName string

	if w.cfg.Rotate.MaxSize > 0 {
		for i := 0; ; i++ {
			fName = w.cfg.Name + fmt.Sprintf(".%s.%03d", w.openTime.Format("2006-01-02"), i)
			_, err = os.Lstat(fName)
			if err != nil {
				break
			}
		}
	} else {
		fName = w.cfg.Name + fmt.Sprintf(".%s", w.openTime.Format("2006-01-02"))
		_, err = os.Lstat(fName)
		if err == nil {
			for i := 0; ; i++ {
				fName = w.cfg.Name + fmt.Sprintf(".%s.%03d", w.openTime.Format("2006-01-02"), i)
				_, err = os.Lstat(fName)
				if err != nil {
					break
				}
			}
		}
	}

	w.file.Close()

	if err = os.Rename(w.cfg.Name, fName); err == nil {
		os.Chmod(fName, w.cfg.Rotate.perm)
	}

	w.initFile()

	if w.cfg.Rotate.MaxDays > 0 {
		go w.deleteOldFile()
	}
}

func (w *FileWriter) deleteOldFile() {
	dir := filepath.Dir(w.cfg.Name)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil {
			return
		}

		if !info.IsDir() && info.ModTime().Add(24*time.Hour*time.Duration(w.cfg.Rotate.MaxDays)).Before(time.Now()) && strings.HasPrefix(path, w.cfg.Name) {
			os.Remove(path)
		}
		return
	})
}

func (w *FileWriter) initFile() error {
	dir, _ := path.Split(w.cfg.Name)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(w.cfg.Name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, w.cfg.perm)
	if err != nil {
		return err
	}

	fInfo, err := file.Stat()
	if err != nil {
		return err
	}

	w.curSize = int(fInfo.Size())

	if w.cfg.Async {
		if w.writer != nil {
			w.writer.Reset(file)
		} else {
			w.writer = bufio.NewWriterSize(file, w.cfg.BufSize)
		}
	}

	w.openTime = time.Now()
	w.file = file

	return nil
}

func (w *FileWriter) autoFlush() {
	ticker := time.NewTicker(w.flushInterval)
	for range ticker.C {
		w.Flush()
	}
}
