/**
 * Copyright 2015 @ to2.net.
 * name : fix
 * author : jarryliu
 * date : 2016-05-14 21:31
 * description : 自定义调整
 * history :
 */
package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 监听进程信号,并执行操作。比如退出时应释放资源
func SignalNotify(c chan bool, fn func()) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)
	for {
		switch <-ch {
		case syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM: // 退出时
			log.Println("[ OS][ TERM] - program has exit !")
			fn()
			close(c)
		}
	}
}
