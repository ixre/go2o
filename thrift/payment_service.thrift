namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.payment_service
include "ttype.thrift"

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
    ttype.Result DiscountByIntegral(1:i32 orderId,2:i64 integral,3:bool ignoreOut)
    // 钱包账户支付
    ttype.Result PaymentByWallet(1:string tradeNo,2:bool mergePay,3:string remark)
    // 余额钱包混合支付，优先扣除余额。
    ttype.Result HybridPayment(1:i32 orderId,2:string remark)
    // 完成支付单支付，并传入支付方式及外部订单号
    ttype.Result FinishPayment(1:string tradeNo ,2:string spName,3:string outerNo)


    // 支付网关
    ttype.Result GatewayV1(1:string action,2:i64 userId,3:map<string,string> data)
    // 获取支付预交易数据
    SPrepareTradeData GetPaymentOrderInfo(1:string tradeNo,2:bool mergePay)

   /**
    * 支付单混合支付
    *
    * @param storeCode 店铺编号
    * @param tradeNo   交易号
    * @param data  支付数据
    * @return 支付结果,返回:order_state
    */
   ttype.Result MixedPayment(1:string tradeNo,3:list<SRequestPayData> data)

}


/** 支付方式的位值 */
enum EMethodFlag{
    /** 余额抵扣 */
    Balance = 1
    /** 钱包支付 */
    Wallet = 2
    /** 积分兑换 */
    Integral = 4
    /** 用户卡 */
    UserCard = 8
    /** 用户券 */
    UserCoupon = 16
    /** 现金支付 */
    Cash = 32
    /** 银行卡支付 */
    BankCard = 64
    /** 第三方支付,如支付宝等 */
    PaySP = 128
    /** 卖家支付通道 */
    SellerPay = 256
    /** 系统支付通道 */
    SystemPay = 512
}

// 支付方式
enum EPaymentMethod{
    /** 余额抵扣通道 */
    Balance = 1
    /** 钱包支付通道 */
    Wallet = 2
    /** 积分兑换通道 */
    Integral = 3
    /** 用户卡通道 */
    UserCard = 4
    /** 用户券通道 */
    UserCoupon = 5
    /** 现金支付通道 */
    Cash = 6
    /** 银行卡支付通道 */
    BankCard = 7
    /** 第三方支付 */
    PaySP = 8
    /** 卖家支付通道 */
    SellerPay = 9
    /** 系统支付通道 */
    SystemPay = 10
}

/** 支付单 */
struct SPaymentOrder{
    /** 编号 */
    1:i32 ID
    /** 卖家编号 */
    2:i32 SellerId
    /** 交易类型 */
    3:string TradeType
    /** 交易号 */
    4:string TradeNo
    /** 合并支付交单单号 */
    5:string MergeTradeNo
    /** 支付单详情 */
    6:string Subject
    /** 是否为子订单 */
	7:bool SubOrder
    /** 支付单的类型，如购物或其他 */
    8:i32 OrderType
    /** 外部订单号 */
    9:string OutOrderNo
    /** 买家编号 */
    10:i32 BuyerId
    /** 支付用户编号 */
    11:i32 PayUid
    /** 商品金额 */
    12:i32 ItemAmount
    /** 优惠金额  */
    13:i32 DiscountAmount
    /** 调整金额 */
    14:i32 AdjustAmount
    /** 抵扣金额  */
    15:i32 DeductAmount
    /** 共计金额 */
    16:i32 TotalAmount
    /** 手续费 */
    17:i32 ProcedureFee
    /** 实付金额 */
    18:i32 PaidFee
    /** 最终应付金额 */
    19:i32 FinalFee
    /** 可⽤支付方式  */
    20:i32 PayFlag
    /** 实际使用支付方式 */
    21:i32 FinalFlag
    /** 其他支付信息 */
    22:string ExtraData
    /** 交易支付渠道 */
    23:i32 TradeChannel
    /** 外部交易提供商 */
    24:string OutTradeSp
    /** 外部交易订单号 */
    25:string OutTradeNo
    /** 订单状态 */
    26:i32 State
    /** 提交时间 */
    27:i64 SubmitTime
    /** 过期时间 */
    28:i64 ExpiresTime
    /** 支付时间 */
    29:i64 PaidTime
    /** 更新时间 */
    30:i64 UpdateTime
    /** 交易途径交易信息 */
    31:list<SPayTradeChan> TradeChannels
}

/** 支付单预交易数据 */
struct SPrepareTradeData{
    /** 错误码 */
    1:i32 ErrCode
    /** 错误消息 */
    2:string ErrMsg
    /** 交易号 */
    3:string TradeNo
    /** 支付标志 */
    4:i32 PayFlag
    /** 交易状态 */
    5:i32 TradeState
    /** 交易订单 */
    6:list<SPaymentOrderData> TradeOrders
    /** 累计金额 */
    7:i32 TotalAmount
    /** 手续费 */
    8:i32 ProcedureFee
    /** 最终支付金额 */
    9:i32 FinalFee
}

/** 支付单数据 */
struct SPaymentOrderData{
    /** 交易订单号 */
    1:string OrderNo
    /** 标题 */
    2:string Subject
    /** 交易类型 */
    3:string TradeType
    /** 状态 */
    4:i32 State
    /** 手续费 */
    5:i32 ProcedureFee
    /** 最终支付金额 */
    6:i32 FinalFee
}

/** 请求支付数据 */
struct SRequestPayData{
    /** 支付通道标志 */
    1:i32 Method
    /** 通道标签 */
    2:string Tag
    /** 支付金额 */
    3:i32 Amount
}

/** 支付单项 */
struct SPayTradeChan{
    /** 编号 */
    1:i32 ID
    /** 交易单号 */
    2:i32 TradeNo
    /** 支付途径 */
    3:i32 PayChan
    /** 是否为内置支付途径 */
    4:i32 InternalChan
    /** 支付金额 */
    5:i32 PayAmount
    /** 通道数据 */
    6:string ChanData
}
