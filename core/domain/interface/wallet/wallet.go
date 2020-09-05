// create for src 29/11/2017 ( jarrysix@gmail.com )
package wallet

import "go2o/core/infrastructure/domain"

var (
	// 别名
	Alias = "钱包"
	// 金额放大比例
	AmountRateSize = 100
	// 提现暂停
	TakeOutPause = false
	// 最低提现金额
	MinTakeOutAmount = 100
	// 最高提现金额
	MaxTakeOutAmount = 10000000
)

const (
	// 个人钱包
	TPerson = 1
	// 企业钱包
	TMerchant = 2
)

const (
	// 抵扣
	FlagDiscount = 1 << iota
	// 充值
	FlagCharge
)

const (
	// 正常
	StatNormal = 1
	// 已禁用
	StatDisabled = 2
	// 已封停
	StatClosed = 3
)

const (
	// 未设置
	ReviewNotSet = 0
	// 等待审核
	ReviewAwaiting = 1
	// 审核失败
	ReviewReject = 2
	// 审核成功
	ReviewPass = 3
	// 已确认
	ReviewConfirm = 4
	// 审核终止
	ReviewAbort = 5
)

const (
	// 用户充值
	CUserCharge = 1
	// 系统自动充值
	CSystemCharge = 2
	// 客服充值
	CServiceAgentCharge = 3
	// 退款充值
	CRefundCharge = 4
)

const (
	// 赠送金额
	KCharge = 1
	// 客服赠送
	KServiceAgentCharge = 2
	// 钱包收入
	KIncome = 3
	// 失效
	KExpired = 4
	// 客服调整
	KAdjust = 5
	// 扣除
	KDiscount = 6
	// 转入
	KTransferIn = 7
	// 转出
	KTransferOut = 8

	// 冻结
	KFreeze = 9
	// 解冻
	KUnfreeze = 10

	// 转账退款
	KTransferRefund = 11
	// 提现退还到银行卡
	KTakeOutRefund = 12
	// 支付单退款
	KPaymentOrderRefund = 13

	// 提现到银行卡(人工提现)
	KTakeOutToBankCard = 14
	// 提现到第三方
	KTakeOutToThirdPart = 15
)

var (
	ErrSingletonWallet               = domain.NewError("err_wallet_singleton_wallet", "用户已存在相同类型的"+Alias)
	ErrRemarkLength                  = domain.NewError("err_wallet_remark_length", "备注不能超过40字")
	ErrMissingOperator               = domain.NewError("err_wallet_missing_operator", "缺少操作人员")
	ErrAmountZero                    = domain.NewError("err_wallet_amount_zero", "金额不能为零")
	ErrOutOfAmount                   = domain.NewError("err_wallet_not_enough_amount", Alias+"余额不足")
	ErrNoSuchTargetWalletAccount     = domain.NewError("err_wallet_no_such_target_wallet_account", "对方账户不存在")
	ErrNoSuchWalletAccount           = domain.NewError("err_wallet_no_such_wallet_account", "账户不存在")
	ErrTargetWalletAccountNotService = domain.NewError("err_target_wallet_account_not_service", "对方账户不可用")
	ErrWalletDisabled                = domain.NewError("err_wallet_disabled", "账户已被暂停")
	ErrWalletClosed                  = domain.NewError("err_wallet_closed", "账户已被关闭")
	ErrNotSupportTakeOutBusinessKind = domain.NewError("err_not_support_take_out_business_kind", "不支持的提现业务类型")
	ErrTakeOutPause                  = domain.NewError("err_wallet_take_out_pause", "当前"+Alias+"暂停提现")
	ErrLessThanMinTakeAmount         = domain.NewError("err_wallet_less_than_min_take_amount", "低于最低提现金额")
	ErrMoreThanMinTakeAmount         = domain.NewError("err_wallet_more_than_min_take_amount", "超过最大提现金额")
	ErrNoSuchTakeOutLog              = domain.NewError("err_wallet_no_such_take_out_log", "提现记录不存在")
	ErrTakeOutState                  = domain.NewError("err_wallet_member_take_out_state", "提现申请状态错误")
)

type (
	// 钱包
	IWallet interface {
		// 获取聚合根编号
		GetAggregateRootId() int64

		// 哈希值
		Hash() string

		// 节点编号
		NodeId() int

		// 获取
		Get() Wallet

		// 状态
		State() int

		//  获取日志
		GetLog(logId int64) WalletLog

		// 保存
		Save() (int64, error)

		// 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
		Adjust(value int, title, outerNo string, opuId int, opuName string) error

		// 支付抵扣,must是否必须大于0
		Discount(value int, title, outerNo string, must bool) error

		// 冻结余额
		Freeze(value int, title, outerNo string, opuId int, opuName string) error

		// 解冻金额
		Unfreeze(value int, title, outerNo string, opuId int, opuName string) error

		// 将冻结金额标记为失效
		FreezeExpired(value int, remark string) error

		// 收入
		Income(value int, tradeFee int, title, outerNo string) error

		// 充值,kind: 业务类型
		Charge(value int, by int, title, outerNo string, opuId int, opuName string) error

		// 退款,kind: 业务类型
		Refund(value int, kind int, title, outerNo string, opuId int, opuName string) error

		// 转账,title如:转账给xxx， toTitle: 转账收款xxx
		Transfer(toWalletId int64, value int, tradeFee int, title, toTitle, remark string) error

		// 接收转账
		ReceiveTransfer(fromWalletId int64, value int, tradeNo, title, remark string) error

		// 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,tradeFee为手续费
		RequestTakeOut(value int, tradeFee int, kind int, title string) (int64, string, error)

		// 确认提现
		ReviewTakeOut(takeId int64, pass bool, remark string, opuId int, opuName string) error

		// 完成提现
		FinishTakeOut(takeId int64, outerNo string) error

		// 分页钱包日志
		PagingLog(begin int, over int, opt map[string]string, sort string) (int, []*WalletLog)
	}

	// 钱包仓储
	IWalletRepo interface {
		// 创建钱包
		CreateWallet(v *Wallet) IWallet
		// 获取钱包账户
		GetWallet(walletId int64) IWallet
		// 根据用户编号获取钱包账户
		GetWalletByUserId(userId int64, walletType int) IWallet
		// 获取日志
		GetLog(walletId int64, logId int64) *WalletLog
		// 检查钱包是否匹配/是否存在
		CheckWalletUserMatch(userId int64, walletType int, walletId int64) bool
		// 获取分页钱包日志
		PagingWalletLog(walletId int64, nodeId int, begin int, over int, where string, sort string) (int, []*WalletLog)

		// auto generate by gof
		// Get WalletLog
		GetWalletLog_(primary interface{}) *WalletLog
		// GetBy WalletLog
		GetWalletLogBy_(where string, v ...interface{}) *WalletLog
		// Select WalletLog
		SelectWalletLog_(where string, v ...interface{}) []*WalletLog
		// Save WalletLog
		SaveWalletLog_(v *WalletLog) (int, error)
		// Delete WalletLog
		DeleteWalletLog_(primary interface{}) error
		// Batch Delete WalletLog
		BatchDeleteWalletLog_(where string, v ...interface{}) (int64, error)

		// auto generate by gof
		// Get Wallet
		GetWallet_(primary interface{}) *Wallet
		// GetBy Wallet
		GetWalletBy_(where string, v ...interface{}) *Wallet
		// Select Wallet
		SelectWallet_(where string, v ...interface{}) []*Wallet
		// Save Wallet
		SaveWallet_(v *Wallet) (int, error)
		// Delete Wallet
		DeleteWallet_(primary interface{}) error
		// Batch Delete Wallet
		BatchDeleteWallet_(where string, v ...interface{}) (int64, error)
	}

	Wallet struct {
		// 钱包编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 哈希值
		HashCode string `db:"hash_code"`
		// 节点编号
		NodeId int `db:"node_id"`
		// 用户编号
		UserId int64 `db:"user_id"`
		// 钱包类型
		WalletType int `db:"wallet_type"`
		// 钱包标志
		WalletFlag int `db:"wallet_flag"`
		// 余额
		Balance int `db:"balance"`
		// 赠送余额
		PresentBalance int `db:"present_balance"`
		// 调整金额
		AdjustAmount int `db:"adjust_amount"`
		// 冻结余额
		FreezeAmount int `db:"freeze_amount"`
		// 结余金额
		LatestAmount int `db:"latest_amount"`
		// 失效账户余额
		ExpiredAmount int `db:"expired_amount"`
		// 总充值金额
		TotalCharge int `db:"total_charge"`
		// 累计赠送金额
		TotalPresent int `db:"total_present"`
		// 总支付额
		TotalPay int `db:"total_pay"`
		// 状态
		State int `db:"state"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 钱包日志
	WalletLog struct {
		// 编号
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 钱包编号
		WalletId int64 `db:"wallet_id"`
		// 业务类型
		Kind int `db:"kind"`
		// 标题
		Title string `db:"title"`
		// 外部通道
		OuterChan string `db:"outer_chan"`
		// 外部订单号
		OuterNo string `db:"outer_no"`
		// 变动金额
		Value int `db:"value"`
		// 余额
		Balance int `db:"balance"`
		// 交易手续费
		TradeFee int `db:"trade_fee"`
		// 操作人员用户编号
		OperatorId int `db:"op_uid"`
		// 操作人员名称
		OperatorName string `db:"op_name"`
		// 备注
		Remark string `db:"remark"`
		// 审核状态
		ReviewState int `db:"review_state"`
		// 审核备注
		ReviewRemark string `db:"review_remark"`
		// 审核时间
		ReviewTime int64 `db:"review_time"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
