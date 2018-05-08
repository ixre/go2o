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
	"fmt"
	"github.com/garyburd/redigo/redis"
	"go2o/core"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"go2o/core/variable"
	"log"
	"strconv"
	"strings"
	"time"
)

// 监视新订单
func superviseOrder(ss []Service) {
	sv := rsi.ShoppingService
	notify := func(orderNo string, sub bool, ss []Service) {
		o, _ := sv.GetOrder(thrift.Context, orderNo, sub)
		if o != nil {
			for _, v := range ss {
				if !v.OrderObs(o) {
					break
				}
			}
		}
	}
	// 监听队列
	conn := core.GetRedisConn()
	defer conn.Close()
	for {
		arr, err := redis.ByteSlices(conn.Do("BLPOP",
			variable.KvOrderBusinessQueue, 0)) //取出队列的一个元素
		if err == nil {
			//通知订单更新
			orderNo := string(arr[1])
			sub := strings.HasPrefix(orderNo, "sub!")
			if sub {
				orderNo = orderNo[4:]
			}
			//log.Println("----- 订单号：",orderNo, "; 是否子订单：", isSub)
			go notify(orderNo, sub, ss)

		} else {
			appCtx.Log().Println("[ Daemon][ OrderQueue][ Error]:",
				err.Error(), "; retry after 10 seconds.")
			time.Sleep(time.Second * 10)
		}
	}
}

// 监视新会员
func superviseMemberUpdate(ss []Service) {
	sv := rsi.MemberService
	notify := func(id int64, action string, ss []Service) {
		m, _ := sv.GetMember(thrift.Context, id)
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
				go notify(int64(id), mArr[1], ss)
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
	sv := rsi.PaymentService
	notify := func(id int, ss []Service) {
		order, _ := sv.GetPaymentOrderById(thrift.Context, int32(id))
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

// 从RDS键中找到订单号，如：go2o:queue:sub!1234345435 , go2o:queue:2
func testIdFromRdsKey(key string) (orderNo string, sub bool, err error) {
	arr := strings.Split(key, ":")
	orderNo = arr[len(arr)-1]
	sub = strings.HasPrefix(orderNo, "sub!")
	if sub {
		orderNo = orderNo[4:]
	}
	return orderNo, sub, err
}

// 检测已过期的订单并标记
func detectOrderExpires() {
	if appCtx.Debug() {
		log.Println("[ Order]: detect order time out ...")
	}
	conn := core.GetRedisConn()
	defer conn.Close()
	tick := getTick(time.Now())
	key := fmt.Sprintf("%s:*:%s", variable.KvOrderExpiresTime, tick)
	//key = "go2o:order:timeout:11-0-2:*"
	//获取标记为等待过期的订单
	ss := rsi.ShoppingService
	list, err := redis.Strings(conn.Do("KEYS", key))
	if err == nil {
		for _, oKey := range list {
			orderNo, isSub, err := testIdFromRdsKey(oKey)
			if err == nil && orderNo != "" {
				err = ss.CancelOrder(orderNo, isSub, "订单超时,自动取消")
				//清除待取消记录
				conn.Do("DEL", oKey)
				//log.Println("---",orderId,"---",unix, "--", time.Now().Unix(), v, err)
			}
		}
	} else {
		log.Println("[ Daemon][ Order][ Cancel][ Error]:",
			err.Error(), "; retry after 10 seconds.")
		time.Sleep(time.Second * 10)
	}
}

// 订单自动收货
func orderAutoReceive() {
	if appCtx.Debug() {
		log.Println("[ Order]: order auto receive ...")
	}
	conn := core.GetRedisConn()
	defer conn.Close()
	tick := getTick(time.Now())
	key := fmt.Sprintf("%s:*:%s", variable.KvOrderAutoReceive, tick)
	//key = "go2o:order:autoreceive:11-0-2:*"
	//获取标记为自动收货的订单
	ss := rsi.ShoppingService
	list, err := redis.Strings(conn.Do("KEYS", key))
	if err == nil {
		for _, oKey := range list {
			orderNo, isSub, err := testIdFromRdsKey(oKey)
			//log.Println("----",oKey,orderId,isSub,err)
			if err == nil && orderNo != "" {
				err = ss.BuyerReceived(orderNo, isSub)
			}
		}
	} else {
		log.Println("[ Daemon][ Order][ Receive][ Error]:",
			err.Error(), "; retry after 10 seconds.")
		time.Sleep(time.Second * 10)
	}
}
