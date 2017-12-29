package logs

import (
	"sync"
	"time"
)

type timeFormat struct {
	sync.RWMutex
	tk      *time.Ticker
	day     int
	nowData []byte
}

var layout = "2006-01-02 15:04:05"

var _time = new(timeFormat)

func init() {
	_time.tk = time.NewTicker(time.Second)
	refresh(time.Now())

	go func() {
		for now := range _time.tk.C {
			refresh(now)
		}
	}()

}

func refresh(now time.Time) {
	_time.Lock()
	_time.nowData = []byte(now.Format(layout))
	_time.day = now.Day()
	_time.Unlock()
}

func nowData() (data []byte) {
	_time.RLock()
	data = _time.nowData
	_time.RUnlock()
	return
}

func day() (d int) {
	_time.RLock()
	d = _time.day
	_time.RUnlock()
	return
}
