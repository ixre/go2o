namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.wallet_service

include "ttype.thrift"

/** 钱包服务 */
service WalletService{
    /** 创建钱包，并返回钱包编号 */
    ttype.Result CreateWallet(1:i64 userId,2:i32 walletType,3:i32 flag,4:string remark)

    /** 获取钱包编号，如果钱包不存在，则返回0 */
    i64 GetWalletId(1:i64 userId,2:i32 walletType)

    /** 获取钱包账户 */
    SWallet GetWallet(1:i64 walletId)

    /** 获取钱包日志 */
    SWalletLog GetWalletLog(1:i64 walletId,2:i64 id)

    /** 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称 */
	ttype.Result Adjust(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,  5:i32 opuId, 6:string opuName)

    /** 支付抵扣,must是否必须大于0 */
    ttype.Result Discount(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,5:bool must)

    /** 冻结余额 */
    ttype.Result Freeze(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,  5:i32 opuId, 6:string opuName)

    /** 解冻金额 */
    ttype.Result Unfreeze(1:i64 walletId,2:i32 value, 3:string title,4:string outerNo,  5:i32 opuId, 6:string opuName)

    /** 充值,kind: 业务类型 */
    ttype.Result Charge(1:i64 walletId,2:i32 value,3:i32 by, 4:string title,5:string outerNo,  6:i32 opuId, 7:string opuName)

    /** 转账,title如:转账给xxx， toTitle: 转账收款xxx */
    ttype.Result Transfer(1:i64 walletId,2:i64 toWalletId, 3:i32 value,4:i32 tradeFee,5:string remark)

    /** 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,tradeFee为手续费 */
    ttype.Result RequestTakeOut(1:i64 walletId,2:i32 value,3:i32 tradeFee,4:i32 kind,5:string title)

    /** 确认提现 */
    ttype.Result ReviewTakeOut(1:i64 walletId,2:i64 takeId,3:bool reviewPass,4:string remark, 5:i32 opuId, 6:string opuName)

    /** 完成提现 */
    ttype.Result FinishTakeOut(1:i64 walletId,2:i64 takeId,3:string outerNo)

    /** 获取分页钱包日志 */
    ttype.SPagingResult PagingWalletLog(1:i64 walletId,2:ttype.SPagingParams params)
}

/** 钱包类型 */
enum EWalletType{
	/** 个人钱包 */
	TPerson = 1
	/** 商家钱包 */
	TMerchant = 2
}

/** 钱包标志 */
enum EWalletFlag{
	/** 抵扣 */
	FlagDiscount = 1
	/** 充值 */
	FlagCharge = 2
}

/** 充值方式 */
enum EChargeKind{
	/** 用户充值 */
	CUserCharge = 1
	/** 系统自动充值 */
	CSystemCharge = 2
	/** 客服充值 */
	CServiceAgentCharge = 3
	/** 退款充值 */
	CRefundCharge = 4
}

/** 提现方式 */
enum ETakeOutKind{
    /** 提现到银行卡(人工提现) */
	KTakeOutToBankCard = 14
	/** 提现到第三方 */
	KTakeOutToThirdPart = 15
}

/** 钱包日志种类 */
enum EWalletLogKind{
	/** 赠送金额 */
	KCharge = 1
	/** 客服赠送 */
	KServiceAgentCharge = 2
	/** 钱包收入 */
	KIncome = 3
	/** 失效 */
	KExpired = 4
	/** 客服调整 */
	KAdjust = 5
	/** 扣除 */
	KDiscount = 6
	/** 转入 */
	KTransferIn = 7
	/** 转出 */
	KTransferOut = 8

	/** 冻结 */
	KFreeze = 9
	/** 解冻 */
	KUnfreeze = 10

	/** 转账退款 */
	KTransferRefund = 11
	/** 提现退还到银行卡 */
	KTakeOutRefund = 12
	/** 支付单退款 */
	KPaymentOrderRefund = 13

	/** 提现到银行卡(人工提现) */
	KTakeOutToBankCard = 14
	/** 提现到第三方 */
	KTakeOutToThirdPart = 15
}

/** 钱包 */
struct SWallet {
    /** 钱包编号 */
    1:i64 ID
    /** 哈希值 */
    2:string HashCode
    /** 节点编号 */
    3:i32 NodeId
    /** 用户编号 */
    4:i64 UserId
    /** 钱包类型 */
    5:i32 WalletType
    /** 钱包标志 */
    6:i32 WalletFlag
    /** 余额 */
    7:i32 Balance
    /** 赠送余额 */
    8:i32 PresentBalance
    /** 调整金额 */
    9:i32 AdjustAmount
    /** 冻结余额 */
    10:i32 FreezeAmount
    /** 结余金额 */
    11:i32 LatestAmount
    /** 失效账户余额 */
    12:i32 ExpiredAmount
    /** 总充值金额 */
    13:i32 TotalCharge
    /** 累计赠送金额 */
    14:i32 TotalPresent
    /** 总支付额 */
    15:i32 TotalPay
    /** 状态 */
    16:i32 State
    /** 备注 */
    17:string Remark
    /** 创建时间 */
    18:i64 CreateTime
    /** 更新时间 */
    19:i64 UpdateTime
}

/** 钱包日志 */
struct SWalletLog {
    /** 编号 */
    1:i64 ID
    /** 钱包编号 */
    2:i64 WalletId
    /** 业务类型 */
    3:i32 Kind
    /** 标题 */
    4:string Title
    /** 外部通道 */
    5:string OuterChan
    /** 外部订单号 */
    6:string OuterNo
    /** 变动金额 */
    7:i32 Value
    /** 余额 */
    8:i32 Balance
    /** 交易手续费 */
    9:i32 TradeFee
    /** 操作人员用户编号 */
    10:i32 OperatorId
    /** 操作人员名称 */
    11:string OperatorName
    /** 备注 */
    12:string Remark
    /** 审核状态 */
    13:i32 ReviewState
    /** 审核备注 */
    14:string ReviewRemark
    /** 审核时间 */
    15:i64 ReviewTime
    /** 创建时间 */
    16:i64 CreateTime
    /** 更新时间 */
    17:i64 UpdateTime
}