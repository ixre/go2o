/**
 * Copyright 2015 @ 56x.net.
 * name : defined.go
 * author : jarryliu
 * date : 2015-08-13 08:55
 * description :
 * history :
 */
package enum

const (
	Unknown int32 = -1
	False   int32 = 0
	True    int32 = 1
)

const (
	// 未提交审核
	ReviewNone int32 = 0
	// 待审核
	ReviewPending int32 = 1
	// 审核未通过
	ReviewRejected int32 = 2
	// 审核成功
	ReviewApproved int32 = 3
	// 已复核
	ReviewCompleted int32 = 4
	// 审核作废
	ReviewAbort int32 = 5
)

// 审核文本字典
var ReviewTextMap = map[int32]string{
	ReviewNone:      "未提交",
	ReviewPending:   "待审核",
	ReviewRejected:  "审核未通过",
	ReviewApproved:  "审核通过",
	ReviewCompleted: "交易完成",
	ReviewAbort:     "已取消",
}

// 审核状态名称
func ReviewString(r int32) string {
	return ReviewTextMap[r]
}

// 商户结算模式
type MchSettleMode int

const (
	// 按供货价销售额比例结算
	MchModeSettleByCost MchSettleMode = 1
	// 按销售额比例结算
	MchModeSettleByRate MchSettleMode = 2
	// 按单数结算
	MchModeSettleByOrderQuantity MchSettleMode = 3
)

// 金额依据
const (
	/** 未设置 */
	AmountBasisNotSet = 1
	/** 按金额 */
	AmountBasisByAmount = 2
	/** 按百分比 */
	AmountBasisByPercent = 3
)

const (
	// 百分比比例放大倍数，保留3位小数;0.56 * 10000 = 560
	RATE_PERCENT float64 = 10000
	// 金额比例放大倍数;0.95 * 100 = 95
	RATE_AMOUNT float64 = 100
	// 折扣比例放大倍数; 0.9 * 1000 = 900
	RATE_DISCOUNT float64 = 1000
)
