package ording

//
////todo: departed
//const (
//	//已作废
//	ORDER_CANCEL = -1
//	//订单已创建
//	ORDER_CREATED = 0
//	//订单处理中
//	ORDER_PROCESSING = 1
//	//订单发货中
//	ORDER_SENDING = 2
//	//订单完成
//	ORDER_FINISH = 3
//
//	//线上付款
//	PAY_ONLINE = 1
//	//线下付款
//	PAY_OFFLINE = 2
//
//	/****** 站点状态 *********/
//	//合作商网站关闭
//	PTSITE_CLOSED = 0
//	//合作商网站正常
//	PTSITE_NORMAL = 1
//
//	/****** 积分返回类型 *********/
//	INTEGRAL_TYPE_SYSTEMSEND = 1
//	INTEGRAL_TYPE_LOGINSEND  = 2
//	INTEGRAL_TYPE_ORDER      = 3
//	INTEGRAL_TYPE_BACK       = 4
//	INTEGRAL_TYPE_EXCHANGE   = 12
//)
//
//var (
//	ORDER_STATE_TEXTS      = [5]string{"已作废", "待确认", "等待配送", "配送中", "已完成"}
//	FRONT_SHOP_STATE_TEXTS = [3]string{"停用", "营业中", "暂停营业"}
//)
//
//func GetPaymentName(i int) string {
//	switch i {
//	case 1:
//		return "网上支付"
//	default:
//	case 2:
//		return "线下支付"
//	}
//	return ""
//}
//
//func GetOrderStatusName(status int) string {
//
//	//	<option value="100">一所有状态一</option>
//	//	<option value="0">待确认</option>
//	//	<option value="1">等待配送</option>
//	//	<option value="2">配送中</option>
//	//	<option value="3" style="color:green">已完成</option>
//	//	<option value="-1" style="color:#ff0000">已作废</option>
//
//	//	switch status{
//	//	case ORDER_CREATED:
//	//		return "待确认"
//	//	case ORDER_CANCEL:
//	//		return "已作废"
//	//	case ORDER_PROCESSING:
//	//		return "等待配送"
//	//	case ORDER_SENDING:
//	//		return "配送中"
//	//	case ORDER_FINISH:
//	//		return "已完成"
//	//	}
//	//	return ""
//
//	if l := status + 1; l < len(ORDER_STATE_TEXTS) {
//		return ORDER_STATE_TEXTS[status+1]
//	}
//	return "未知状态"
//}
//
//func GetFrontShopStateName(state int) string {
//	return FRONT_SHOP_STATE_TEXTS[state]
//}
