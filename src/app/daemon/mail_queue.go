/**
 * Copyright 2015 @ S1N1 Team.
 * name : mail_queue
 * author : jarryliu
 * date : 2015-07-27 17:06
 * description :
 * history :
 */
package daemon

import (
	"go2o/src/core/domain/interface/partner/mss"
	mssIns "go2o/src/core/infrastructure/mss"
	"go2o/src/core/variable"
	"time"
	"go2o/src/core/domain/interface/enum"
)

var (
	mailChan chan int
)

func startMailQueue() {
	for {
		if i, _ := appCtx.Storage().GetInt(variable.KvNewMailTask); i == enum.FALSE {
			sendQueue()
			appCtx.Storage().Set(variable.KvNewMailTask, enum.TRUE)
		}
		time.Sleep(time.Second * 5)
	}
}

func sendQueue() {
	var list = []*mss.MailTask{}
	appCtx.Db().GetOrm().Select(&list, "is_send = 0 OR is_failed = 1")
	mailChan = make(chan int, len(list))
	for _, v := range list {
		go func(ch chan int, t *mss.MailTask) {
			err := mssIns.SendMailWithDefaultConfig(t.Subject, []string{t.SendTo}, []byte(t.Body))
			if err != nil {
				appCtx.Log().PrintErr(err)
				t.IsFailed = 1
				t.IsSend = 1
			} else {
				t.IsSend = 1
				t.IsFailed = 0
			}
			t.SendTime = time.Now().Unix()
			appCtx.Db().GetOrm().Save(t.Id, t)
			mailChan <- 0
		}(mailChan, v)
		<-mailChan
	}
}
