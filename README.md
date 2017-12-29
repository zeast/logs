# Logs
package logs a simple unstructured log implements for Go.

## Features
* Rotate. You can deploy without other rotation system.
* Async write. Save IO.
* Change log level atomic.In general, you can just log error message, if you want to see some details, you can change it without stopping server.

## Example
```go
package main

import (
	"encoding/json"

	"github.com/zeast/logs"
)

func main() {
	w, err := logs.NewFileWriter(
		logs.FileConfig{
			Name: "/tmp/xxx.log",
			Perm: "0644",
			Rotate: &logs.Rotate{
				MaxSize: 1024 * 1024, //1G
				MaxDays: 10,          //10 day
				Perm:    "0444",
			},
		},
	)
	
	//check err
	_ = err
	
	logs.SetBaseWriter(w)
	
	logs.LogFuncCall(true)
	
	logs.Debugf("debug message")
	
	var b1 = []byte{123, 34, 107, 34, 58, 34, 115, 111, 109, 101, 32, 109, 97, 114, 115, 104, 97, 108, 32, 98, 121, 116, 101, 115, 34, 125}
	if logs.LogDebug() {
		//unmarhsal for human readable
		var v map[string]string
		json.Unmarshal(b1, &v)
		logs.Debug(v)
	}
	
	var b2 = []byte("big byte slice")
	if logs.LogInfo() {
		//if you use logs.Info() immediate, one memory copy will happen
		logs.Info(string(b2))
	}

}
```