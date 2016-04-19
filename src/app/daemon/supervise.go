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
	"github.com/jsix/gof"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"strconv"
	"strings"
	"time"
)

// 监视新订单
func superviseOrder(ss []Service) {
	conn := core.GetRedisConn()
	defer conn.Close()
	var id int
	for {
		arr, err := redis.Values(conn.Do("BLPOP", variable.KvOrderBusinessQueue,
			0)) //取出队列的一个元素
		if err == nil {
			id, err = strconv.Atoi(string(arr[1].([]byte)))
			if err == nil { //通知订单更新
				order := dps.ShoppingService.GetOrderById(id)
				for _, v := range ss {
					if !v.OrderObs(order) {
						break
					}
				}
			}
		} else {
			appCtx.Log().Println("[ Daemon][ OrderQueue][ Error] - ", err.Error())
			break
		}
	}
}

// 监视新会员
func superviseMemberUpdate(ss []Service) {
	conn := core.GetRedisConn()
	defer conn.Close()
	var id int
	for {
		arr, err := redis.Values(conn.Do("BLPOP", variable.KvMemberUpdateQueue, 0))
		if err == nil { //通知会员修改,格式如: 1-[create|update]
			s := string(arr[1].([]byte))
			mArr := strings.Split(s, "-")
			id, err = strconv.Atoi(mArr[0])
			if err == nil {
				m := dps.MemberService.GetMember(id)
				for _, v := range ss {
					if !v.MemberObs(m, mArr[1] == "create") {
						break
					}
				}
			}
		} else {
			appCtx.Log().Println("[ Daemon][ MemberQueue][ Error] - ", err.Error())
			break
		}
	}
}

func detectOrderExpires(a gof.App) {
	conn := core.GetRedisConn()
	defer conn.Close()
	list, _ := redis.Strings(conn.Do("KEYS", variable.KvOrderExpiresTime+"*")) //获取标记为等待过期的订单
	ss := dps.ShoppingService
	for _, v := range list {
		if unix, err := redis.Int64(conn.Do("GET", v)); err == nil {
			if unix < time.Now().Unix() {
				arr := strings.Split(v[strings.LastIndex(v, ":")+1:], "_") //key: 商户号_订单号
				if partnerId, err := strconv.Atoi(arr[0]); err == nil {
					err = ss.CancelOrder(partnerId, arr[1], "订单超时,自动取消") //清除
					conn.Do("DEL", v)                                    //清除待取消记录
					//log.Println(unix,"--",time.Now().Unix(),v,err)
				}
			}
		}
	}
}

//
//func confirmNewOrder(app gof.App, dfs []Func) {
//
//	if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewCreatedOrder); i == enum.TRUE {
//		appCtx.Log().Printf("[ DAEMON][ ORDER][ CONFIRM] - begin invoke confirm handler.")
//		if dfs == nil || len(dfs) == 0 {
//			confirmOrderQueue(app)
//		} else {
//			for _, v := range dfs {
//				v(app)
//			}
//		}
//		appCtx.Storage().Set(variable.KvHaveNewCreatedOrder, enum.FALSE)
//	}
//}
//
//func completedOrderObs(app gof.App, dfs []Func) {
//	if len(dfs) < 0 {
//		return
//	}
//	if i, _ := appCtx.Storage().GetInt(variable.KvHaveNewCompletedOrder); i == enum.TRUE {
//		appCtx.Log().Printf("[ DAEMON][ ORDER][ FINISHED] - begin invoke finish handler.\n")
//		for _, v := range dfs {
//			v(app)
//		}
//		appCtx.Storage().Set(variable.KvHaveNewCompletedOrder, enum.FALSE)
//	}
//}
//
//type orderInfo struct {
//	PartnerId int
//	OrderNo   string
//}
//
//func confirmOrderQueue(app gof.App) {
//	var list []*orderInfo = []*orderInfo{}
//	appCtx.Db().GetOrm().SelectByQuery(&list, fmt.Sprintf("SELECT partner_id,order_no FROM pt_order WHERE status=%d",
//		enum.ORDER_WAIT_CONFIRM))
//	for _, v := range list {
//		err := dps.ShoppingService.ConfirmOrder(v.PartnerId, v.OrderNo)
//		if err != nil {
//			appCtx.Log().Printf("[ DAEMON][ ORDER][ ERROR] - %s\n", err.Error())
//		}
//	}
//}
