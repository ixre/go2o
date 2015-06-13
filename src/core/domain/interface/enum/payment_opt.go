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
    PaymentOptionNames []string  =[]string{"在线付款","货到付款","转账汇款"}
)
const(
    PaymentOnlinePay int = 1   // 线上付款
    PaymentOfflineCashPay = 2  // 线下现金付款
    PaymentRemit int = 3       // 转账汇款
)
