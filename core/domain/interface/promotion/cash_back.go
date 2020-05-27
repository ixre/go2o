/**
 * Copyright 2015 @ to2.net.
 * name : cash_back
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

// 返现促销
type ICashBackPromotion interface {
	// 获取领域编号
	GetDomainId() int32

	// 设置详细的促销信息
	SetDetailsValue(*ValueCashBack) error

	// 获取自定义数据
	GetDataTag() map[string]string
}

type ValueCashBack struct {
	// 编号
	Id int32 `db:"id" pk:"yes"`

	// 最低金额要求
	MinFee int `db:"min_fee"`

	// 返还金额
	BackFee int `db:"back_fee"`

	// 返还方式,1:充值到余额 2:直接抵扣订单
	BackType int `db:"back_type"`

	// 自定义数据,用于分析处理自定义的数据
	DataTag string `db:"data_tag"`
}
