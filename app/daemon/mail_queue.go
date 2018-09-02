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
	"github.com/gomodule/redigo/redis"
	"go2o/core"
	"go2o/core/domain/interface/mss"
	mssIns "go2o/core/infrastructure/mss"
	"go2o/core/variable"
	"strconv"
	"time"
)

var (
	mailChan chan int
)

// 邮件队列
func startMailQueue(ss []Service) {
	conn := core.GetRedisConn()
	defer conn.Close()
	//var id int
	for {
		arr, err := redis.Values(conn.Do("BLPOP", variable.KvNewMailTask, 0))
		if err == nil {
			_, err = strconv.Atoi(string(arr[1].([]byte)))
			if err == nil {
				//todo: 此处获取所有需发送的邮件,应去掉从数据库批量查询操作
				sendForWaitingQueue(ss)
			}
		} else {
			appCtx.Log().Println("[ Daemon][ MailQueue][ Error] - ", err.Error())
			break
		}
	}
}

func sendForWaitingQueue(ss []Service) {
	var list = []*mss.MailTask{}
	err := appCtx.Db().GetOrm().Select(&list, "is_send = 0 OR is_failed = 1")
	if err == nil && len(list) > 0 {
		for _, s := range ss {
			if !s.HandleMailQueue(list) {
				break
			}
		}
	}
}

func handleMailQueue(list []*mss.MailTask) {
	mailChan = make(chan int, len(list))
	for _, v := range list {
		go func(ch chan int, t *mss.MailTask) {
			err := mssIns.SendMailWithDefaultConfig(t.Subject, []string{t.SendTo}, []byte(t.Body))
			if err != nil {
				appCtx.Log().Error(err)
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
