// create for src 29/11/2017 ( jarrysix@gmail.com )
package wallet

import "github.com/ixre/go2o/core/infrastructure/domain"

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
	// FlagDiscount 抵扣
	FlagDiscount = 1 << iota
	// FlagCharge 充值
	FlagCharge
)

const (
	// StatNormal 正常
	StatNormal = 1
	// StatDisabled 已禁用
	StatDisabled = 2
	// StatClosed 已封停
	StatClosed = 3
)

const (
	// 未设置
	ReviewNotSet = 0
	// ReviewPending 等待审核
	ReviewPending = 1
	// ReviewReject 审核失败
	ReviewReject = 2
	// ReviewPass 审核成功
	ReviewPass = 3
	// ReviewConfirm 已确认
	ReviewConfirm = 4
	// ReviewAbort 审核终止
	ReviewAbort = 5
)

const (
	// CUserCharge 用户充值
	CUserCharge = 1
	// CSystemCharge 系统自动充值
	CSystemCharge = 2
	// CServiceAgentCharge 客服充值
	CServiceAgentCharge = 3
	// CRefundCharge 退款充值
	CRefundCharge = 4
)

const (
	// KCharge 赠送金额
	KCharge = 1
	// KCarry 钱包收入
	KCarry = 2
	// KExpired 失效
	KExpired = 3
	// KAdjust 客服调整
	KAdjust = 4
	// KConsume 消费
	KConsume = 5
	// KDiscount 扣除
	KDiscount = 6
	// KTransferIn 转入
	KTransferIn = 7
	// KTransferOut 转出
	KTransferOut = 8

	// KFreeze 冻结
	KFreeze = 9
	// KUnfreeze 解冻
	KUnfreeze = 10

	// KTransferRefund 转账退款
	KTransferRefund = 11
	// KWithdrawRefund 提现退还到银行卡
	KWithdrawRefund = 12
	// KPaymentOrderRefund 支付单退款
	KPaymentOrderRefund = 13

	// KWithdrawExchange 提现并兑换到余额
	KWithdrawExchange int = 21
	// KWithdrawToBankCard 提现到银行卡(人工提现)
	KWithdrawToBankCard = 22
	// KWithdrawToThirdPart 提现到第三方
	KWithdrawToThirdPart = 23
)

var (
	ErrSingletonWallet               = domain.NewError("err_wallet_singleton_wallet", "用户已存在相同类型的"+Alias)
	ErrWalletName                    = domain.NewError("err_wallet_name", "钱包名称为空或超出长度")
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
	ErrNoSuchAccountLog              = domain.NewError("err_wallet_no_such_take_out_log", "钱包记录不存在")
	ErrWithdrawState                 = domain.NewError("err_wallet_member_take_out_state", "提现申请状态错误")
	ErrNotSupport                    = domain.NewError("err_wallet_not_support", "不支持该操作")
)

type (
	// TransactionData 钱包交易数据
	TransactionData struct {
		// 描述
		TransactionTitle string
		// 金额(含手续费)
		Amount int
		// 交易手续费
		TransactionFee int
		// 外部单号,如果非系统订单，添加前缀，如：XT:20140109345
		OuterNo string
		// 备注
		TransactionRemark string
		// 交易流水编号,对冻结流水进行更新时,传递该参数
		TransactionId int
	}
	Operator struct {
		OperatorUid  int
		OperatorName string
	}

	// IWallet 钱包
	IWallet interface {
		// GetAggregateRootId 获取聚合根编号
		GetAggregateRootId() int64

		// Hash 哈希值
		Hash() string

		// NodeId 节点编号
		NodeId() int

		// Get 获取
		Get() Wallet

		// State 状态
		State() int

		// GetLog 获取日志
		GetLog(logId int64) WalletLog

		// Save 保存
		Save() (int64, error)

		// Adjust 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
		Adjust(value int, title, outerNo string, remark string, operatorUid int, operatorName string) error

		// Consume 消费
		Consume(amount int, title string, outerNo string, remark string) error

		// Discount 抵扣,must是否必须大于0
		Discount(amount int, title, outerNo string, must bool) error

		// Freeze 冻结余额,返回LogId
		Freeze(data TransactionData, operator Operator) (int, error)

		// Unfreeze 解冻金额
		Unfreeze(amount int, title, outerNo string, operatorUid int, operatorName string) error

		// FreezeExpired 将冻结金额标记为失效
		FreezeExpired(amount int, remark string) error

		// CarryTo 收入/入账, freeze是否先冻结; 返回交易流水ID
		CarryTo(tx TransactionData, freeze bool) (transactionId int, err error)

		// ReviewCarryTo 审核入账
		ReviewCarryTo(requestId int, pass bool, reason string) error

		// Charge 充值,kind: 业务类型
		Charge(value int, kind int, title, outerNo string, remark string, operatorUid int, operatorName string) error

		// Refund 退款,kind: 业务类型
		Refund(value int, kind int, title, outerNo string, operatorUid int, operatorName string) error

		// Transfer 转账,title如:转账给xxx， toTitle: 转账收款xxx
		Transfer(toWalletId int64, value int, transactionFee int, title, toTitle, remark string) error

		// ReceiveTransfer 接收转账
		ReceiveTransfer(fromWalletId int64, value int, tradeNo, title, remark string) error

		// RequestWithdrawal 申请提现,kind：提现方式,返回info_id,交易号 及错误,amount为提现金额,transactionFee为手续费
		RequestWithdrawal(amount int, transactionFee int, kind int, title string,
			accountNo string, accountName string, bankName string) (int64, string, error)

		// ReviewWithdrawal 确认提现
		ReviewWithdrawal(takeId int64, pass bool, remark string, operatorUid int, operatorName string) error

		// FinishWithdrawal 完成提现
		FinishWithdrawal(takeId int64, outerNo string) error

		// PagingLog 分页钱包日志
		PagingLog(begin int, over int, opt map[string]string, sort string) (int, []*WalletLog)
	}

	// 钱包仓储
	IWalletRepo interface {
		// 创建钱包
		CreateWallet(userId int64, username string, walletType int, walletName string, flag int) IWallet
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
		// 根据钱包代码获取钱包
		GetWalletByCode(code string) IWallet
	}

	Wallet struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 哈希值
		HashCode string `db:"hash_code"`
		// 节点编号
		NodeId int `db:"node_id"`
		// 用户编号
		UserId int64 `db:"user_id"`
		// 用户名,方便查询数据
		Username string `db:"user_name"`
		// 钱包类型
		WalletType int `db:"wallet_type"`
		// 钱包标志
		WalletFlag int `db:"wallet_flag"`
		// 钱包名称
		WalletName string `db:"wallet_name"`
		// 余额
		Balance int64 `db:"balance"`
		// 赠送余额
		PresentBalance int64 `db:"present_balance"`
		// 调整禁遏
		AdjustAmount int `db:"adjust_amount"`
		// 冻结金额
		FreezeAmount int `db:"freeze_amount"`
		// 结余金额
		LatestAmount int `db:"latest_amount"`
		// 失效账户余额
		ExpiredAmount int `db:"expired_amount"`
		// 总充值金额
		TotalCharge int64 `db:"total_charge"`
		// 累计赠送金额
		TotalPresent int `db:"total_present"`
		// 总支付额
		TotalPay int `db:"total_pay"`
		// 状态
		State int16 `db:"state"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// WalletLog 钱包日志
	WalletLog struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 钱包编号
		WalletId int64 `db:"wallet_id"`
		// 钱包用户名,冗余用于展示
		WalletUser string `db:"wallet_user"`
		// 业务类型
		Kind int `db:"kind"`
		// 标题
		Subject string `db:"subject"`
		// 外部通道
		OuterChan string `db:"outer_chan"`
		// 外部订单号
		OuterNo string `db:"outer_no"`
		// 变动金额
		ChangeValue int64 `db:"change_value"`
		// 余额
		Balance int64 `db:"balance"`
		// 交易手续费
		TransactionFee int `db:"procedure_fee"`
		// 操作人员用户编号
		OperatorUid int `db:"opr_uid"`
		// 操作人员名称
		OperatorName string `db:"opr_name"`
		// 提现账号
		AccountNo string `db:"account_no"`
		// 提现账户名称
		AccountName string `db:"account_name"`
		// 提现银行名称
		BankName string `db:"bank_name"`
		// 审核状态
		ReviewStatus int `db:"review_status"`
		// 审核备注
		ReviewRemark string `db:"review_remark"`
		// 审核时间
		ReviewTime int64 `db:"review_time"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
