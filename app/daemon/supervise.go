/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ixre/go2o/core/initial"
	"github.com/ixre/go2o/core/service"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof/util"
)

// 监视新订单
func superviseOrder(ss []Service) {
	notify := func(orderNo string, sub bool, ss []Service) {
		trans, cli, _ := service.OrderServiceClient()
		// 这里应处理子订单和父订单

		//o, _ := cli.GetOrder(context.TODO(), &proto.GetOrderRequest{
		//	OrderNo:  orderNo,
		//	SubOrder: sub,
		//})
		o, _ := cli.GetOrder(context.TODO(), &proto.OrderRequest{
			OrderNo: orderNo,
		})
		trans.Close()
		if o != nil {
			for _, v := range ss {
				if !v.OrderObs(o) {
					break
				}
			}
		}
	}
	// 监听队列
	conn := initial.GetRedisConn()
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
	trans, cli, _ := service.MemberServiceClient()
	defer trans.Close()
	notify := func(id int64, action string, ss []Service) {
		m, _ := cli.GetMember(context.TODO(), &proto.MemberIdRequest{MemberId: id})
		if m != nil {
			for _, v := range ss {
				if !v.MemberObs(m, action == "create") {
					break
				}
			}
		}
	}
	id := 0
	conn := initial.GetRedisConn()
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
	notify := func(id int, ss []Service) {
		//trans, cli, _ := service.PaymentServiceClient()
		//defer trans.Close()
		// order, _ := cli.GetPaymentOrderById(context.TODO(), &proto.Int32{Value: int32(id)})
		// if order != nil {
		// 	for _, v := range ss {
		// 		if !v.PaymentOrderObs(order) {
		// 			break
		// 		}
		// 	}
		// }
	}
	id := 0
	conn := initial.GetRedisConn()
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
func memberAutoUnlock() {
	if appCtx.Debug() {
		log.Println("[ Order]: execute member unlock job ...")
	}
	conn := initial.GetRedisConn()
	defer conn.Close()
	tick := util.GetMinuteSlice(time.Now(), 1)
	key := fmt.Sprintf("%s:%s:*", variable.KvMemberAutoUnlock, tick)
	//获取标记为等待过期的订单
	list, err := redis.Strings(conn.Do("KEYS", key))
	if err == nil {
		trans, cli, err := service.MemberServiceClient()
		if err == nil {
			for _, oKey := range list {
				memberId, _ := redis.Int64(conn.Do("GET", oKey))
				if memberId > 0 {
					cli.Unlock(context.TODO(), &proto.MemberIdRequest{MemberId: memberId})
				}
				conn.Do("DEL", oKey)
			}
			trans.Close()
		}
	} else {
		log.Println("[ Daemon][ Member][ Unlock][ Error]:",
			err.Error(), "; retry after 10 seconds.")
		time.Sleep(time.Second * 10)
	}
}

// 订单自动收货
func orderAutoReceive() {
	if appCtx.Debug() {
		log.Println("[ Order]: order auto receive ...")
	}
	conn := initial.GetRedisConn()
	defer conn.Close()
	tick := getTick(time.Now())
	key := fmt.Sprintf("%s:*:%s", variable.KvOrderAutoReceive, tick)
	//key = "go2o:order:autoreceive:11-0-2:*"
	//获取标记为自动收货的订单
	list, err := redis.Strings(conn.Do("KEYS", key))
	if err == nil {
		for _, oKey := range list {
			orderNo, isSub, err := testIdFromRdsKey(oKey)
			//log.Println("----",oKey,orderId,isSub,err)
			if err == nil && orderNo != "" {
				trans, cli, _ := service.OrderServiceClient()
				ret, _ := cli.BuyerReceived(context.TODO(), &proto.OrderNo{
					OrderNo: orderNo,
					Sub:     isSub,
				})
				trans.Close()
				if ret.ErrCode > 0 {
					err = errors.New(ret.ErrMsg)
				}

			}
		}
	} else {
		log.Println("[ Daemon][ Order][ Receive][ Error]:",
			err.Error(), "; retry after 10 seconds.")
		time.Sleep(time.Second * 10)
	}
}
