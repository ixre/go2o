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
	ReviewNotSet int = 0
	// 等待审核
	ReviewAwaiting int = 1
	// 审核失败
	ReviewReject int = 2
	// 审核成功
	ReviewPass int = 3
	// 已确认
	ReviewConfirm int = 4
)

// 审核文本字典
var ReviewTextMap = map[int]string{
	ReviewNotSet:   "未提交",
	ReviewAwaiting: "待审核",
	ReviewReject:   "审核不通过",
	ReviewPass:     "审核通过",
	ReviewConfirm:  "已确认",
}

// 审核状态名称
func ReviewString(r int) string {
	return ReviewTextMap[r]
}

// 商户结算模式
type MchSettleMode int

const (
	// 结算供货价
	MchModeSettleByCost MchSettleMode = 1
	// 按比例结算
	MchModeSettleByRate MchSettleMode = 2
)
