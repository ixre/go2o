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
//RefundStatAwaitingVendor = 1 + iota
//// 退款取消
//RefundStatCancelled
//// 运营商拒绝退款
//RefundStatVendorDecline
//// 调解状态
//RefundStatIntercede
//// 等待确认退款
//RefundStatAwaittingConfirm
//// 退款成功
//RefundStatCompleted
)

type (
	// 退款单,同退货单。只是不退货物。退款单需要付款后才能退款。
	IRefundOrder interface {
		//// 获取领域对象编号
		//GetDomainId() int32
		//
		//// 获取值
		//Value() RefundOrder
		//
		//// 提交退款申请
		//Submit() error
		//
		//// 取消申请退款
		//Cancel() error
		//
		//// 拒绝退款
		//Decline(remark string) error
		//
		//// 同意退款
		//Agree() error
		//
		//// 确认退款
		//Confirm() error
		//
		//// 申请调解
		//RequestIntercede() error
		//
		//// 调解后直接操作退款或退回
		//IntercedeHandle(pass bool, remark string) error
	}
	// 退款单
	RefundOrder struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"no"`
		// 金额
		Amount float32 `db:"amount"`
		// 退款方式：1.退回余额  2: 原路退回
		//RefundType int `db:"refund_type"`
		// 是否为全部退款
		//AllRefund int `db:"all_refund"`
		// 是否已退款
		IsRefund int `db:"is_refund"`
	}
)
