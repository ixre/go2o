namespace go define
include "ttype.thrift"


/** 钱包服务 */
service WalletService{
    // 创建钱包，并返回钱包编号
    ttype.Result CreateWallet(1:i64 userId,2:i32 walletType,3:i32 flag,4:string remark)

    // 获取钱包编号，如果钱包不存在，则返回0
    i32 GetWalletId(1:i64 userId,2:i32 walletType)

    // 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
	ttype.Result Adjust(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,  5:i32 opuId, 6:string opuName)

    // 支付抵扣,must是否必须大于0
    ttype.Result Discount(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,5:i32 opuId, 6:string opuName,7:bool must)

    // 冻结余额
    ttype.Result Freeze(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,  5:i32 opuId, 6:string opuName)

    // 解冻金额
    ttype.Result Unfreeze(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,  5:i32 opuId, 6:string opuName)

    // 充值,kind: 业务类型
    ttype.Result Charge(1:i64 walletId,2:i32 value,3:i32 by, 4:string title,5:string outerNo,  6:i32 opuId, 7:string opuName)

    // 转账,title如:转账给xxx， toTitle: 转账收款xxx
    ttype.Result Transfer(1:i64 walletId,2:i64 toWalletId, 3:i32 value,4:i32 tradeFee,5:string remark)

    // 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,tradeFee为手续费
    ttype.Result RequestTakeOut(1:i64 walletId,2:i32 value,3:i32 tradeFee,4:i32 kind,5:string title)

    // 确认提现
    ttype.Result ReviewTakeOut(1:i64 walletId,2:i64 takeId,3:bool reviewPass,4:string remark, 5:i32 opuId, 6:string opuName)

    // 完成提现
    ttype.Result FinishTakeOut(1:i64 walletId,2:i64 takeId,3:string outerNo)
}

// 钱包类型
enum WalletType{
	// 个人钱包
	TPerson = 1
	// 商家钱包
	TMerchant = 2
}

// 钱包标志
enum WalletFlag{
	// 抵扣
	FlagDiscount = 1
	// 充值
	FlagCharge = 2
}

// 充值方式
enum ChargeKind{
	// 用户充值
	CUserCharge = 1
	// 系统自动充值
	CSystemCharge = 2
	// 客服充值
	CServiceAgentCharge = 3
	// 退款充值
	CRefundCharge = 4
}

// 提现方式
enum TakeOutKind{
    // 提现到银行卡(人工提现)
	KTakeOutToBankCard = 16
	// 提现到第三方
	KTakeOutToThirdPart = 17
}
