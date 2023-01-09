/**
 * Copyright 2015 @ at3.net.
 * name : app.go
 * author : jarryliu
 * date : 2016-10-14 00:21
 * description :
 * history :
 */
package app

import (
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/shell"
	"time"
)


// 自动安装包
func AutoInstall() {
	execInstall()
	d := time.Second * 15
	t := time.NewTimer(d)
	for {
		select {
		case <-t.C:
			if err := execInstall(); err == nil {
				t.Reset(d)
			} else {
				break
			}
		}
	}
}

func execInstall() error {
	_, _, err := shell.Run("go install .", false)
	if err != nil {
		log.Println("[ Go2o][ Install]:", err)
	}
	return err
}
