/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-15 13:58
 * description :
 * history :
 */

package member

type Account struct {
	MemberId int `db:"member_id" pk:"yes" json:"memberId"`

	// 积分
	Integral int `db:"integral"`

	// 余额
	Balance float32 `db:"balance" json:"balance"`

	// 赠送余额
	PresentBalance float32 `db:"present_balance" json:"presentBalance"`

	// 总赠送金额
	TotalPresentFee float32 `db:"total_present_fee" json:"totalPresentFee"`

	// 总消费额
	TotalFee float32 `db:"total_fee" json:"totalFee"`

	// 总充值额
	TotalCharge float32 `db:"total_charge" json:"totalCharge"`

	// 总支付额
	TotalPay float32 `db:"total_pay" json:"totalPay"`

	// 更新时间
	UpdateTime int64 `db:"update_time" json:"updateTime"`
}
