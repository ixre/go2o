/**
 * Copyright 2015 @ to2.net.
 * name : MemberSummary
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

// 会员概览信息
type MemberSummary struct {
	MemberId int32 `db:"id" auto:"yes" pk:"yes"`
	// 用户名
	Usr string `db:"user"`
	// 昵称
	Name string `db:"name"`
	// 头像
	Avatar string `db:"avatar"`
	// 经验值
	Exp int32 `db:"exp"`
	// 等级
	Level int32 `db:"level"`
	// 等级名称
	LevelName string `db:"level_name"`
	// 等级标识
	LevelSign string `db:"program_sign"`
	// 等级是否为正式会员
	LevelOfficial int `db:"is_official"`
	// 邀请码
	InvitationCode string `db:"invitation_code"`
	// 积分
	Integral int64 `db:"integral"`
	// 账户余额
	Balance           float32 `db:"balance"`
	WalletBalance     float32 `db:"wallet_balance"`
	GrowBalance       float32 `db:"grow_balance"`
	GrowAmount        float32 `db:"grow_amount"`         // 理财总投资金额,不含收益
	GrowEarnings      float32 `db:"grow_earnings"`       // 当前收益金额
	GrowTotalEarnings float32 `db:"grow_total_earnings"` // 累积收益金额
	UpdateTime        int64   `db:"update_time"`
}
