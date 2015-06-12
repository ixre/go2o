/**
 * Copyright 2015 @ S1N1 Team.
 * name : payment_opt
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package enum

var(
    PaymentOptionNames []string  =[]string{"货到付款","在线付款","转账汇款"}
)
const(
    PaymentOfflineCashPay = 1  // 线下现金付款
    PaymentOnlinePay int = 2   // 线上付款
    PaymentRemit int = 3       // 转账汇款
)
