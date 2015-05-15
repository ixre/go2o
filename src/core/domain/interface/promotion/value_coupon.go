/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-02 21:50
 * description :
 * history :
 */

package promotion

type ValueCoupon struct {
	Id int `db:"id" auto:"yes" pk:"yes"`
	//合作商编号
	PartnerId int `db:"partner_id"`
	//优惠码
	Code string `db:"code"`

	// 优惠券描述
	// 如果未填写则自动产生
	Description string `db:"description"`

	// 优惠码可用数量
	Amount int `db:"amount"`

	// 优惠码数量
	TotalAmount int `db:"total_amount"`

	//优惠金额
	Fee int `db:"fee"`

	//赠送积分
	Integral int `db:"integral"`

	//订单折扣(不打折为100)
	Discount int `db:"discount"`

	//等级限制
	MinLevel int `db:"min_level"`

	//订单金额限制
	MinFee int `db:"min_fee"`

	BeginTime int64 `db:"begin_time"`
	OverTime  int64 `db:"over_time"`

	//是否允许使用
	AllowEnable int `db:"allow_enable"`

	//是否需要绑定。反之可以直接使用
	NeedBind int `db:"need_bind"`

	CreateTime int64 `db:"create_time"`
	UpdateTime int64 `db:"update_time"`
}
