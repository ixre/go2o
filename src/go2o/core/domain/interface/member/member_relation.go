/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-09 10:28
 * description :
 * history :
 */
package member

//会员关联表
type MemberRelation struct {
	MemberId int `db:"member_id" pk:"yes"`
	//会员卡号
	CardId string `db:"card_id"`
	//推荐人
	TgId int `db:"tg_id"`
	//注册关联商家编号
	Reg_PtId int `db:"reg_ptid"`
}
