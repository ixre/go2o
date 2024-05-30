package log

import (
	"fmt"
	"runtime"

	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/gof/log"
)

func isDebug() bool {
	return false
}
func getLogger() log.ILogger {
	return provide.GetApp().Log()
}
func Error(err error) {
	_, f, line, _ := runtime.Caller(1)
	if err != nil && isDebug() {
		getLogger().Println(
			fmt.Sprintf("[ ERROR]:%s , File:%s line:%d", err.Error(), f, line))
	}
}

func Println(v ...interface{}) {
	if isDebug() {
		getLogger().Println(v...)
	}
}

func Printf(s string, v ...interface{}) {
	if isDebug() {
		getLogger().Printf(s, v...)
	}
}
