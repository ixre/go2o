/**
 * Copyright 2015 @ z3q.net.
 * name : after_sales
 * author : jarryliu
 * date : 2016-07-16 14:41
 * description :
 * history :
 */
package after_sales

type (

	// 售后单
	IAfterSalesOrder interface {
		// 获取领域编号
		GetDomainId() int

		// 提交售后申请
		Submit() error

		// 取消申请
		Cancel() error

		// 拒绝售后服务
		Decline(reason string) error

		// 同意售后服务
		Agree() error

		// 系统确认
		Confirm() error

		// 申请调解
		RequestIntercede() error

		// 退回商品
		ReturnShip(spName string, spOrder string, image string) error

		// 已收货
		ReturnReceive() error

		// 设置要退回货物信息
		SetItem(itemId int, quantity int) error
	}

	// 售后单
	AfterSalesOrder struct {
		// 编号
		Id int `db:"id"`
		// 订单编号
		OrderId int `db:"order_id"`
		// 类型，退货、换货、维修
		Type int `db:"type"`
		// 退货的商品项编号
		ItemId int `db:"item_id"`
		// 商品数量
		Quantity int `db:"quantity"`
		// 售后原因
		Reason string `db:"reason"`
		// 联系人
		PersonName string `db:"person_name"`
		// 联系电话
		PersonPhone string `db:"person_phone"`
		// 退货的快递服务商编号
		RefundSpName string `db:"rsp_name"`
		// 退货的快递单号
		RefundSpOrder string `db:"rsp_order"`
		// 退货凭证
		RefundSpImage string `db:"rsp_image"`
		// 备注(系统)
		Remark string `db:"remark"`
		// 运营商备注
		VendorRemark string `db:"vendor_remark"`
		// 售后单状态
		State int `db:"state"`
		// 提交时间
		CreateTime int `db:"create_time"`
		// 更新时间
		UpdateTime int `db:"update_time"`
	}
)
