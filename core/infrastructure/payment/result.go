/**
 * Copyright 2015 @ z3q.net.
 * name : result.go
 * author : jarryliu
 * date : 2015-07-27 21:52
 * description :
 * history :
 */
package payment

type Result struct {
	// 状态
	Status int
	// 错误消息
	ErrMsg string
	// 外部交易号(系统订单号)
	OutTradeNo string
	// 交易号
	TradeNo string
	// 金额
	Fee float32
}
