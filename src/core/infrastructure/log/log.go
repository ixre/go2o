package log

import (
	"fmt"
	"go2o/src/core/infrastructure"
	"runtime"
)

func Error(err error) {
	_, f, line, _ := runtime.Caller(1)
	if err != nil && infrastructure.DebugMode {
		infrastructure.GetApp().Log().Println(
			fmt.Sprintf("[ ERROR]:%s , File:%s line:%d", err.Error(), f, line))
	}
}

func Println(v ...interface{}) {
	if infrastructure.DebugMode {
		infrastructure.GetApp().Log().Println(v...)
	}
}

func Printf(s string, v ...interface{}) {
	if infrastructure.DebugMode {
		infrastructure.GetApp().Log().Printf(s, v...)
	}
}
