/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
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
	//推荐人（会员）
	RefereesId int `db:"invi_member_id"`
	//注册关联商户编号
	RegisterMerchantId int `db:"reg_merchant_id"`
}
