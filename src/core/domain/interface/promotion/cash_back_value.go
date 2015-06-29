/**
 * Copyright 2015 @ S1N1 Team.
 * name : cash_back_value
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package promotion

type ValueCashBack struct {
	// 编号
	Id int `db:"id" pk:"yes"`

	// 最低金额要求
	MinFee int `db:"min_fee"`

	// 返还金额
	BackFee int `db:"back_fee"`

	// 返还方式,1:充值到余额 2:直接抵扣订单
	BackType int `db:"back_type"`

	// 自定义数据,用于分析处理自定义的数据
	DataTag string `db:"data_tag"`
}
