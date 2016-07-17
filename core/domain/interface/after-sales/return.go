/**
 * Copyright 2015 @ z3q.net.
 * name : return
 * author : jarryliu
 * date : 2016-07-16 14:51
 * description :
 * history :
 */
package afterSales

type (
	// 退款货接口
	IReturnOrder interface {
		// 同意退货
		Return() error
	}
	// 退款单
	ReturnOrder struct {
		// 编号
		Id int `db:"Id"`
		// 金额
		Amount float32 `db:"Amount"`
		// 是否已退款
		IsRefund int `db:"IsRefund"`
	}
)
