/**
 * Copyright 2015 @ to2.net.
 * name : payment_opt
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package enum

var (
	PaymentOptionNames  []string = []string{"在线付款", "货到付款", "转账汇款"}
	PaymentHelpContents []string = []string{
		"<span style='color:red'>请完成最后一步，点击按钮在线支付订单。</span>",
		"请您在收到商品时候，现金支付，或使用POS机刷卡(不同地区可能不支持).",
		"请通过银行汇款至我们的银行账户后，并联系客服进行订单付款。详情点击<a target='_blank' href='/content/page?id=bank-transfer'>这里</a>查看。",
	}
)

const (
	PaymentOnlinePay      int32 = 1 // 线上付款
	PaymentOfflineCashPay int32 = 2 // 线下现金付款
	PaymentRemit          int32 = 3 // 转账汇款
)

// 获取支付帮助内容
//todo: 需要商户可以自定义设置
func GetPaymentHelpContent(opt int32) string {
	if opt-1 <= 0 {
		opt = 1
	}
	return PaymentHelpContents[opt-1]
}
