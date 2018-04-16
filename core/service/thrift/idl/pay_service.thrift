namespace go define

include "ttype.thrift"


//支付单
struct SPaymentOrder {
    1: i32 ID
    2: string TradeNo
    3: i32 VendorId
    4: i32 Type
    5: i32 OrderId
    6: string Subject
    7: i64 BuyUser
    8: i64 PaymentUser
    9: double TotalAmount
    10: double BalanceDiscount
    11: double  IntegralDiscount
    12: double SystemDiscount
    13: double CouponDiscount
    14: double SubAmount
    15: double AdjustmentAmount
    16: double FinalFee
    17: i32 PaymentOptFlag
    18: i32 PaymentSign
    19: string OuterNo
    20: i64 CreateTime
    21: i64 PaidTime
    22: i32 State
    /** 交易类型 */
    23:string TradeType
}

// 支付服务
service PaymentService{
    // 创建支付单并提交
    ttype.Result SubmitPaymentOrder(1:SPaymentOrder o)
    // 根据支付单号获取支付单
    SPaymentOrder GetPaymentOrder(1:string orderNo)
    // 根据交易号获取支付单编号
    i32 GetPaymentOrderId(1:string tradeNo)
    // 根据编号获取支付单
    SPaymentOrder GetPaymentOrderById(1:i32 id)
    // 调整支付单金额
    ttype.Result AdjustOrder(1:string paymentNo, 2:double amount)
    // 余额抵扣
    ttype.Result DiscountByBalance(1:i32 orderId,2:string remark )
   // 积分抵扣支付单
    ttype.DResult DiscountByIntegral(1:i32 orderId,2:i64 integral,3:bool ignoreOut)
    // 钱包账户支付
    ttype.Result PaymentByWallet(1:i32 orderId,2:string remark)
    // 余额钱包混合支付，优先扣除余额。
    ttype.Result HybridPayment(1:i32 orderId,2:string remark)
    // 完成支付单支付，并传入支付方式及外部订单号
    ttype.Result FinishPayment(1:string tradeNo ,2:string spName,3:string outerNo)
    // 支付网关
    ttype.Result GatewayV1(1:string action,2:i64 userId,3:map<string,string> data)
}

