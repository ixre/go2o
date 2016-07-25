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
	"go2o/core"
	"go2o/core/repository"
	"go2o/core/variable"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CustomFix() {
	variable.AliasMemberIM = "微信"
	variable.MemberImNote = "填写微信后才可领取红包"
	variable.MemberImRequired = true
	variable.AliasMemberExt1 = "QQ"
	variable.MemberExt1Note = ""
	variable.MemberExt1Show = false
	variable.AliasMemberExt2 = "支付宝"
	variable.MemberExt2Note = "便于支付宝打款"
	variable.MemberExt2Show = !true
	variable.AliasMemberExt3 = "扩展3"
	variable.MemberExt3Note = ""
	variable.MemberExt3Show = false
	variable.AliasMemberExt4 = "扩展4"
	variable.MemberExt4Note = ""
	variable.MemberExt4Show = false
	variable.AliasMemberExt5 = "扩展5"
	variable.MemberExt5Note = ""
	variable.MemberExt5Show = false
	variable.AliasMemberExt6 = "扩展6"
	variable.MemberExt6Note = ""
	variable.MemberExt6Show = false

	// 注册后赠送10w积分
	repository.DefaultRegistry.PresentIntegralNumOfRegister = 10000000
}

func SignalNotify(c chan bool) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM: // 退出时
			log.Println("[ OS][ TERM] - program has exit !")
			dispose()
			close(c)
		}
	}
}

func dispose() {
	core.GetRedisPool().Close()
}
