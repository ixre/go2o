/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"github.com/garyburd/redigo/redis"
	"go2o/core"
	"go2o/core/service/dps"
	"go2o/core/variable"
	"strconv"
	"strings"
	"time"
)

// 监视新订单
func superviseOrder(ss []Service) {
	sv := dps.ShoppingService
	notify := func(id int, ss []Service) {
		o := sv.GetSubOrder(int32(id))
		if o != nil {
			for _, v := range ss {
				if !v.OrderObs(o) {
					break
				}
			}
		}
	}
	// 监听队列
	id := 0
	conn := core.GetRedisConn()
	defer conn.Close()
	for {
		arr, err := redis.ByteSlices(conn.Do("BLPOP",
			variable.KvOrderBusinessQueue, 0)) //取出队列的一个元素
		if err == nil {
			//通知订单更新
			id, err = strconv.Atoi(string(arr[1]))
			if err == nil {
				go notify(id, ss)
			}
		} else {
			appCtx.Log().Println("[ Daemon][ OrderQueue][ Error]:",
				err.Error(), "; retry after 10 seconds.")
			time.Sleep(time.Second * 10)
		}

	}
}

// 监视新会员
func superviseMemberUpdate(ss []Service) {
	sv := dps.MemberService
	notify := func(id int32, action string, ss []Service) {
		m, _ := sv.GetMember(id)
		if m != nil {
			for _, v := range ss {
				if !v.MemberObs(m, action == "create") {
					break
				}
			}
		}
	}
	id := 0
	conn := core.GetRedisConn()
	defer conn.Close()
	for {
		arr, err := redis.ByteSlices(conn.Do("BLPOP",
			variable.KvMemberUpdateQueue, 0))
		if err == nil {
			//通知会员修改,格式如: 1-[create|update]
			s := string(arr[1])
			mArr := strings.Split(s, "-")
			id, err = strconv.Atoi(mArr[0])
			if err == nil {
				go notify(int32(id), mArr[1], ss)
			}
		} else {
			appCtx.Log().Println("[ Daemon][ MemberQueue][ Error]:",
				err.Error(), "; retry after 10 seconds.")
			time.Sleep(time.Second * 10)
		}
	}
}

// 监视支付单完成
func supervisePaymentOrderFinish(ss []Service) {
	sv := dps.PaymentService
	notify := func(id int, ss []Service) {
		order := sv.GetPaymentOrder(int32(id))
		if order != nil {
			for _, v := range ss {
				if !v.PaymentOrderObs(order) {
					break
				}
			}
		}
	}
	id := 0
	conn := core.GetRedisConn()
	defer conn.Close()
	for {
		arr, err := redis.ByteSlices(conn.Do("BLPOP",
			variable.KvPaymentOrderFinishQueue, 0))
		if err == nil {
			//通知服务
			s := string(arr[1])
			id, err = strconv.Atoi(s)
			if err == nil {
				go notify(id, ss)
			}
		} else {
			appCtx.Log().Println("[ Daemon][ PaymentOrderQueue][ Error]:",
				err.Error(), "; retry after 10 seconds.")
			time.Sleep(time.Second * 10)
		}
	}
}

// 检测已过期的订单并标记
func detectOrderExpires() {
	conn := core.GetRedisConn()
	defer conn.Close()
	//获取标记为等待过期的订单
	orderId := 0
	list, _ := redis.Strings(conn.Do("KEYS", variable.KvOrderExpiresTime+"*"))
	ss := dps.ShoppingService
	for _, v := range list {
		unix, err := redis.Int64(conn.Do("GET", v))
		if err == nil {
			//获取过期时间
			if unix < time.Now().Unix() {
				//订单号
				orderId, err = strconv.Atoi(v[len(variable.KvOrderExpiresTime):])
				err = ss.CancelOrder(int32(orderId), "订单超时,自动取消")
				//清除待取消记录
				conn.Do("DEL", v)
				//log.Println("---",orderId,"---",unix, "--", time.Now().Unix(), v, err)
			}
		} else {
			appCtx.Log().Println("[ Daemon][ Order][ Cancel][ Error]:",
				err.Error(), "; retry after 10 seconds.")
			time.Sleep(time.Second * 10)
		}
	}
}
