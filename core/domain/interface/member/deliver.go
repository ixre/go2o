/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 22:06
 * description :
 * history :
 */

package member

type (
	// 收货地址
	IDeliverAddress interface {
		GetDomainId() int
		GetValue() DeliverAddress
		SetValue(*DeliverAddress) error
		Save() (int, error)
	}

	// 收货地址
	DeliverAddress struct {
		//编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		//会员编号
		MemberId int `db:"member_id"`
		//收货人
		RealName string `db:"real_name"`
		//电话
		Phone string `db:"phone"`
		//省
		Province int `db:"province"`
		//市
		City int `db:"city"`
		//区
		District int `db:"district"`
		//地区(省市区连接)
		Area string `db:"area"`
		//地址
		Address string `db:"address"`
		//是否默认
		IsDefault int `db:"is_default"`
	}
)
