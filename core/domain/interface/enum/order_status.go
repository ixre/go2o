/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:25
 * description :
 * history :
 */

package enum

type OrderState int

const (
	// 已取消
	//ORDER_CANCEL = 0
	// 订单待支付
	//ORDER_WAIT_PAYMENT = 1
	// 订单待确认
	//ORDER_WAIT_CONFIRM = 2
	// 订单待发货
	//ORDER_WAIT_DELIVERY = 3
	// 订单配送中,等待收货
	//ORDER_WAIT_RECEIVE = 4
	// 已收货
	//ORDER_RECEIVED = 5
	// 订单完成
	//ORDER_COMPLETED = 6

	// 可进行流程的状态
	ORDER_SETUP_STATE = "2,4"
)

//
//func (t OrderState) String() string {
//	switch t {
//	case ORDER_CANCEL:
//		return "已取消"
//	//case ORDER_WAIT_PAYMENT:
//	//	return "待付款"
//	case ORDER_WAIT_CONFIRM:
//		return "待确认"
//	case ORDER_WAIT_DELIVERY:
//		return "待发货"
//	case ORDER_WAIT_RECEIVE:
//		return "配送中"
//	case ORDER_RECEIVED:
//		return "已收货"
//	case ORDER_COMPLETED:
//		return "已完成"
//	}
//	return "Error State"
//}

const (
	/****** 站点状态 *********/
	//合作商网站关闭
	PARTNER_SITE_CLOSED = 0
	//合作商网站正常
	PARTNER_SITE_NORMAL = 1
	/****** 积分返回类型 *********/
	INTEGRAL_TYPE_SYSTEM_PRESENT = 1
	INTEGRAL_TYPE_LOGIN_PRESENT  = 2
	INTEGRAL_TYPE_ORDER          = 3
	INTEGRAL_TYPE_BACK           = 4
	INTEGRAL_TYPE_EXCHANGE       = 12
)

var (
	FRONT_SHOP_STATE_TEXTS = [3]string{"停用", "营业中", "暂停营业"}
)

// 获取支付方式名称
func GetPaymentName(i int32) string {
	switch i {
	case 1:
		return "在线支付"
	case 2:
		return "货到付款"
	default:
	case 3:
		return "转账汇款"
	}
	return ""
}

func GetFrontShopStateName(state int32) string {
	return FRONT_SHOP_STATE_TEXTS[state]
}
