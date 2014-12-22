package log

import (
	"com/infrastructure"
	"runtime"
)

func PrintErr(err error) {
	_, f, line, _ := runtime.Caller(1)
	if err != nil && infrastructure.DebugMode {
		infrastructure.Context.Log().Println(
			"[ERROR]:%s , File:%s line:%d", err.Error(), f, line)
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
