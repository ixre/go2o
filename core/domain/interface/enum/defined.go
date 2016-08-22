/**
 * Copyright 2015 @ z3q.net.
 * name : defined.go
 * author : jarryliu
 * date : 2015-08-13 08:55
 * description :
 * history :
 */
package enum

const (
	TRUE  int = 1
	FALSE int = 0
)

const (
	// 未设置
	ReviewNotSet = 0
	// 等待审核
	ReviewAwaiting = 1
	// 审核成功
	ReviewPass = 2
	// 审核失败
	ReviewReject = 3
)

// 商户结算模式
type MchSettleMode int

const (
	// 结算供货价
	MchModeSettleByCost MchSettleMode = 1
	// 按比例结算
	MchModeSetttleByRate MchSettleMode = 2
)
