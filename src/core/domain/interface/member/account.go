/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-15 13:58
 * description :
 * history :
 */

package member

type Account struct {
	MemberId       int     `db:"member_id" pk:"yes" json:"memberId"`
	Integral       int     `db:"integral"`
	Balance        float32 `db:"balance" json:"balance"`
	PresentBalance float32 `db:"present_balance" json:"presentBalance"`
	TotalFee       float32 `db:"total_fee" json:"totalFee"`
	TotalCharge    float32 `db:"total_charge" json:"totalCharge"`
	TotalPay       float32 `db:"total_pay" json:"totalPay"`
	UpdateTime     int64   `db:"update_time" json:"updateTime"`
}
