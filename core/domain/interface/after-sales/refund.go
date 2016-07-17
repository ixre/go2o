/**
 * Copyright 2015 @ z3q.net.
 * name : refund
 * author : jarryliu
 * date : 2016-07-16 14:44
 * description :
 * history :
 */
package afterSales

const (
	// 等待运营商确认
	RefundStatAwaitingVendor = 1 + iota
	// 退款取消
	RefundStatCancelled
	// 运营商拒绝退款
	RefundStatVendorDecline
	// 调解状态
	RefundStatIntercede
	// 等待确认退款
	RefundStatAwaittingConfirm
	// 退款成功
	RefundStatCompleted
)

type (
	// 退款单
	IRefundOrder interface {
		// 获取领域对象编号
		GetDomainId() int

		// 获取值
		GetValue() RefundOrder

		// 提交退款申请
		Submit() error

		// 取消申请退款
		Cancel() error

		// 拒绝退款
		Decline(remark string) error

		// 同意退款
		Agree() error

		// 确认退款
		Confirm() error

		// 申请调解
		RequestIntercede() error

		// 调解后直接操作退款或退回
		IntercedeHandle(pass bool, remark string) error
	}
	// 退款单
	RefundOrder struct {
		// 编号
		Id int `db:"id"`
		// 订单编号
		OrderId int `db:"order_id"`
		// 金额
		Amount float32 `db:"amount"`
		// 退款方式：1.退回余额  2: 原路退回
		RefundType int `db:"refund_type"`
		// 是否为全部退款
		AllRefund int `db:"all_refund"`
		// 退款的商品项编号
		ItemId int `db:"item_id"`
		// 联系人
		PersonName string `db:"person_name"`
		// 联系电话
		PersonPhone string `db:"person_phone"`
		// 退款原因
		Reason string `db:"reason"`
		// 退款单备注(系统)
		Remark string `db:"remark"`
		// 运营商备注
		VendorRemark string `db:"vendor_remark"`
		// 退款状态
		State int `db:"state"`
		// 提交时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
