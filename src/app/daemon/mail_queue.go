/**
 * Copyright 2015 @ z3q.net.
 * name : mail_queue
 * author : jarryliu
 * date : 2015-07-27 17:06
 * description :
 * history :
 */
package daemon

import (
	"go2o/src/core/domain/interface/enum"
	"go2o/src/core/domain/interface/partner/mss"
	mssIns "go2o/src/core/infrastructure/mss"
	"go2o/src/core/variable"
	"time"
)

var (
	mailChan chan int
)

// 邮件队列
func startMailQueue(ss []Service) {
	for {
		if i, _ := appCtx.Storage().GetInt(variable.KvNewMailTask); i == enum.FALSE {
			defer func() {
				appCtx.Storage().Set(variable.KvNewMailTask, enum.TRUE)
			}()
			var list = []*mss.MailTask{}
			err := appCtx.Db().GetOrm().Select(&list, "is_send = 0 OR is_failed = 1")
			if err == nil && len(list) > 0 {
				for _, s := range ss {
					if !s.HandleMailQueue(list) {
						break
					}
				}
			}
			time.Sleep(tickerDuration)
		}
	}
}

func handleMailQueue(list []*mss.MailTask) {
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
