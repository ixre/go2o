package log

import (
	"fmt"
	"go2o/core/infrastructure"
	"runtime"
)

func PrintErr(err error) {
	_, f, line, _ := runtime.Caller(1)
	if err != nil && infrastructure.DebugMode {
		infrastructure.Context.Log().Println(
			fmt.Sprintf("[ERROR]:%s , File:%s line:%d", err.Error(), f, line))
	}
}

func Println(v ...interface{}) {
	if infrastructure.DebugMode {
		infrastructure.Context.Log().Println(v...)
	}
}

func Printf(s string, v ...interface{}) {
	if infrastructure.DebugMode {
		infrastructure.Context.Log().Printf(s, v...)
	}
}
