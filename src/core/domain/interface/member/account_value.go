/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-15 13:58
 * description :
 * history :
 */

package member

// 账户值对象
type AccountValue struct {
	MemberId        int     `db:"member_id" pk:"yes" json:"memberId"`
	Integral        int     `db:"integral"`                                 // 积分
	Balance         float32 `db:"balance" json:"balance"`                   // 账户余额
	FreezesFee      float32 `db:"freezes_fee" json:"freezesFee"`            // 冻结金额
	PresentBalance  float32 `db:"present_balance" json:"presentBalance"`    // 奖金账户余额
	FreezesPresent  float32 `db:"freezes_present" json:"freezesPresent"`    // 冻结赠送额
	TotalPresentFee float32 `db:"total_present_fee" json:"totalPresentFee"` // 总赠送金额
	FlowBalance     float32 `db:"flow_balance" json:"flowBalance"`          // 流动账户余额
	GrowBalance     float32 `db:"grow_balance" json:"growBalance"`          // 当前增利账户金额
	TotalGrowAmount float32 `db:"grow_amount" json:"growAmount"`            // 累积增利金额
	TotalFee        float32 `db:"total_fee" json:"totalFee"`                // 总消费额
	TotalCharge     float32 `db:"total_charge" json:"totalCharge"`          // 总充值额
	TotalPay        float32 `db:"total_pay" json:"totalPay"`                // 总支付额
	UpdateTime      int64   `db:"update_time" json:"updateTime"`            // 更新时间
}
