/**
 * Copyright 2015 @ z3q.net.
 * name : member_summary
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

// 会员概览信息
type MemberSummary struct {
	Id                int     `db:"id" auto:"yes" pk:"yes"`
	Usr               string  `db:"usr"`
	Name              string  `db:"name"`
	Avatar            string  `db:"avatar"`
	Exp               int     `db:"exp"`
	Level             int     `db:"level"`
	LevelName         string  `db:"level_name"`
	Integral          int     `db:"integral"`
	Balance           float32 `db:"balance"`
	PresentBalance    float32 `db:"present_balance"`
	GrowBalance       float32 `db:"grow_balance"`
	GrowAmount        float32 `db:"grow_amount"`         // 理财总投资金额,不含收益
	GrowEarnings      float32 `db:"grow_earnings"`       // 当前收益金额
	GrowTotalEarnings float32 `db:"grow_total_earnings"` // 累积收益金额
	UpdateTime        int64   `db:"update_time"`
}
