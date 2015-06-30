/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:25
 * description :
 * history :
 */

package enum

type OrderState int

const (
	// 已作废
	ORDER_CANCEL = 0
	// 订单已创建
	ORDER_CREATED = 1
	// 订单已确认
	ORDER_CONFIRMED = 2
	// 订单处理中
	ORDER_PROCESSING = 3
	// 订单发货中
	ORDER_SENDING = 4
	// 已收货
	ORDER_RECEIVED = 5
	// 订单完成
	ORDER_COMPLETED = 6

	// 可进行流程的状态
	ORDER_SETUP_STATE = "1,2,3,4,5"
)

func (t OrderState) String() string {
	switch t {
	case ORDER_CANCEL:
		return "已取消"
	case ORDER_CREATED:
		return "待确认"
	case ORDER_CONFIRMED:
		return "已确认"
	case ORDER_PROCESSING:
		return "处理中"
	case ORDER_SENDING:
		return "配送中"
	case ORDER_RECEIVED:
		return "已收货"
	case ORDER_COMPLETED:
		return "已完成"
	}
	return "Error State"
}

const (
	//线下付款
	PAY_OFFLINE = 1

	//线上付款
	PAY_ONLINE = 2
)

const (
	/****** 站点状态 *********/
	//合作商网站关闭
	PARTNER_SITE_CLOSED = 0
	//合作商网站正常
	PARTNER_SITE_NORMAL = 1

	/****** 积分返回类型 *********/
	INTEGRAL_TYPE_SYSTEMSEND = 1
	INTEGRAL_TYPE_LOGINSEND  = 2
	INTEGRAL_TYPE_ORDER      = 3
	INTEGRAL_TYPE_BACK       = 4
	INTEGRAL_TYPE_EXCHANGE   = 12
)

var (
	FRONT_SHOP_STATE_TEXTS = [3]string{"停用", "营业中", "暂停营业"}
)

func GetPaymentName(i int) string {
	switch i {
	case 1:
		return "网上支付"
	case 2:
		return "现金支付"
	default:
	case 3:
		return "银行转账"
	}
	return ""
}

func GetFrontShopStateName(state int) string {
	return FRONT_SHOP_STATE_TEXTS[state]
}
