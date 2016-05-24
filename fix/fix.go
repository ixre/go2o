/**
 * Copyright 2015 @ z3q.net.
 * name : fix
 * author : jarryliu
 * date : 2016-05-14 21:31
 * description : 自定义调整
 * history :
 */
package fix

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CustomFix() {

}

func SignalNotify(c chan bool) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM: // 退出时
			log.Println("[ OS][ TERM] - program has exit !")
			close(c)
		}
	}
}
