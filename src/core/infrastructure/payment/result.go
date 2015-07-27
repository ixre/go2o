/**
 * Copyright 2015 @ S1N1 Team.
 * name : result.go
 * author : jarryliu
 * date : 2015-07-27 21:52
 * description :
 * history :
 */
package payment

type Result struct{
	// 状态
	Status int
	// 订单号
	OrderNo string
	// 交易号
	TradeNo string
	// 金额
	Fee float32
}