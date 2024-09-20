// create for src 29/11/2017 ( jarrysix@gmail.com )
package wallet

import "github.com/ixre/go2o/core/infrastructure/domain"

var (
	// 别名
	Alias = "钱包"
	// 金额放大比例
	AmountRateSize = 100
	// 提现暂停
	WithdrawIsPaused = false
	// 最低提现金额
	MinWithdrawAmount = 100
	// 最高提现金额
	MaxWithdrawAmount = 10000000
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
	// 待提交审核
	ReviewStaging = 0
	// ReviewPending 等待审核
	ReviewPending = 1
	// ReviewReject 审核失败
	ReviewReject = 2
	// ReviewApproved 审核成功
	ReviewApproved = 3
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
	// KCharge 充值金额
	KCharge = 1
	// KCarry 钱包收入
	KCarry = 2
	// KConsume 消费
	KConsume = 3
	// KAdjust 客服调整
	KAdjust = 4
	// KDiscount 扣除
	KDiscount = 5
	// KindRefund 退款
	KRefund = 6
	// KFreeze 冻结
	KFreeze = 7
	// KUnfreeze 解冻
	KUnfreeze = 8
	// KTransferIn 转入
	KTransfer = 9
	// 提现(预留，应增加提现方式)
	KWithdraw = 10
	// KExpired 失效
	KExpired = 11
	// KindExchange 兑换充值, 比如将钱包充值到余额
	KExchange = 12

	// KTransferRefund 转账退款
	KTransferRefund = 13
	// KWithdrawRefund 提现退还到银行卡
	KWithdrawRefund = 14
	// KPaymentOrderRefund 支付单退款
	KPaymentOrderRefund = 15

	// todo: 充值用2开头， 提现用3开头

	// KWithdrawExchange 提现并兑换到余额
	KWithdrawExchange int = 30
	// KWithdrawToBankCard 提现到银行卡(人工提现)
	KWithdrawToBankCard = 31
	// KWithdrawToPayWallet 提现到第三方支付钱包
	KWithdrawToPayWallet = 32
	// KWithdrawCustom 自定义提现
	KWithdrawCustom = 33
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
	ErrNotSupportWithdrawKind        = domain.NewError("err_not_support_take_out_business_kind", "不支持的提现业务类型")
	ErrTakeOutPause                  = domain.NewError("err_wallet_take_out_pause", "当前"+Alias+"暂停提现")
	ErrLessThanMinWithdrawAmount     = domain.NewError("err_wallet_less_than_min_take_amount", "低于最低提现金额")
	ErrMoreThanMinWithdrawAmount     = domain.NewError("err_wallet_more_than_min_take_amount", "超过最大提现金额")
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
		OuterTxNo string
		// 备注
		TransactionRemark string
		// 交易流水编号,对冻结流水进行更新时,传递该参数
		TransactionId int
		// 外部交易用户编号,可为空
		OuterTxUid int
	}

	// TakeOutTransaction 提现交易
	WithdrawTransaction struct {
		// 提现金额
		Amount int
		// 提现手续费
		TransactionFee int
		// 提现类型
		Kind int
		// 提现标题
		TransactionTitle string
		// 银行名称
		BankName string
		// 账号
		AccountNo string
		// 账户名称
		AccountName string
	}
	Operator struct {
		OperatorUid  int
		OperatorName string
	}

	// IWallet 钱包
	IWallet interface {
		// GetAggregateRootId 获取聚合根编号
		domain.IAggregateRoot
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
		Save() (int, error)

		// Adjust 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
		Adjust(value int, title, outerNo string, remark string, operatorUid int, operatorName string) error

		// Consume 消费
		Consume(amount int, title string, outerNo string, remark string) error

		// 预扣消费,将冻结转为消费,扣款后不自动退回余额
		PrefreezeConsume(data TransactionData) error

		// Discount 抵扣,must是否必须大于0
		Discount(amount int, title, outerNo string, must bool) error

		// Freeze 冻结余额,返回LogId
		Freeze(data TransactionData, operator Operator) (int, error)

		// Unfreeze 解冻金额, 传入isRefundBalance是否退回余额,如果为true，则解冻金额后，将自动退回余额
		Unfreeze(amount int, title, outerNo string, isRefundBalance bool, operatorUid int, operatorName string) error

		// FreezeExpired 将冻结金额标记为失效
		FreezeExpired(amount int, remark string) error

		// CarryTo 收入/入账, freeze是否先冻结; 返回交易流水ID
		CarryTo(tx TransactionData, freeze bool) (transactionId int, err error)

		// ReviewCarryTo 审核入账
		ReviewCarryTo(transactionId int, pass bool, reason string) error

		// Charge 充值,kind: 业务类型
		Charge(value int, kind int, title, outerNo string, remark string, operatorUid int, operatorName string) error

		// Refund 退款,kind: 业务类型
		Refund(value int, kind int, title, outerNo string, operatorUid int, operatorName string) error

		// Transfer 转账,title如:转账给xxx， toTitle: 转账收款xxx
		Transfer(toWalletId int64, value int, transactionFee int, title, toTitle, remark string) error

		// ReceiveTransfer 接收转账
		ReceiveTransfer(fromWalletId int64, value int, tradeNo, title, remark string) error

		// RequestWithdrawal 申请提现,kind：提现方式,返回info_id,交易号 及错误,amount为提现金额,transactionFee为手续费
		RequestWithdrawal(tx WithdrawTransaction) (int, string, error)

		// ReviewWithdrawal 确认提现
		ReviewWithdrawal(transactionId int, pass bool, remark string, operatorUid int, operatorName string) error

		// FinishWithdrawal 完成提现
		FinishWithdrawal(transactionId int, outerTxNo string) error

		// PagingLog 分页钱包日志
		PagingLog(begin int, over int, opt map[string]string, sort string) (int, []*WalletLog)
	}

	// 钱包仓储
	IWalletRepo interface {
		// 创建钱包
		CreateWallet(userId int, username string, walletType int, walletName string, flag int) IWallet
		// 获取钱包账户
		GetWallet(walletId int) IWallet
		// 根据用户编号获取钱包账户
		GetWalletByUserId(userId int64, walletType int) IWallet
		// 获取日志
		GetLog(walletId int, logId int64) *WalletLog
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
)

// Wallet 钱包
type Wallet struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 哈希值
	HashCode string `json:"hashCode" db:"hash_code" gorm:"column:hash_code" bson:"hashCode"`
	// 节点编号
	NodeId int `json:"nodeId" db:"node_id" gorm:"column:node_id" bson:"nodeId"`
	// 用户编号
	UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
	// 钱包类型
	WalletType int `json:"walletType" db:"wallet_type" gorm:"column:wallet_type" bson:"walletType"`
	// 钱包标志
	WalletFlag int `json:"walletFlag" db:"wallet_flag" gorm:"column:wallet_flag" bson:"walletFlag"`
	// 钱包名称
	WalletName string `json:"walletName" db:"wallet_name" gorm:"column:wallet_name" bson:"walletName"`
	// 余额
	Balance int `json:"balance" db:"balance" gorm:"column:balance" bson:"balance"`
	// 赠送余额
	PresentBalance int `json:"presentBalance" db:"present_balance" gorm:"column:present_balance" bson:"presentBalance"`
	// 调整禁遏
	AdjustAmount int `json:"adjustAmount" db:"adjust_amount" gorm:"column:adjust_amount" bson:"adjustAmount"`
	// 冻结金额
	FreezeAmount int `json:"freezeAmount" db:"freeze_amount" gorm:"column:freeze_amount" bson:"freezeAmount"`
	// 结余金额
	LatestAmount int `json:"latestAmount" db:"latest_amount" gorm:"column:latest_amount" bson:"latestAmount"`
	// 失效账户余额
	ExpiredAmount int `json:"expiredAmount" db:"expired_amount" gorm:"column:expired_amount" bson:"expiredAmount"`
	// 总充值金额
	TotalCharge int `json:"totalCharge" db:"total_charge" gorm:"column:total_charge" bson:"totalCharge"`
	// 累计赠送金额
	TotalPresent int `json:"totalPresent" db:"total_present" gorm:"column:total_present" bson:"totalPresent"`
	// 总支付额
	TotalPay int `json:"totalPay" db:"total_pay" gorm:"column:total_pay" bson:"totalPay"`
	// 状态
	State int `json:"state" db:"state" gorm:"column:state" bson:"state"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	// 用户名
	Username string `json:"userName" db:"user_name" gorm:"column:user_name" bson:"userName"`
}

func (w Wallet) TableName() string {
	return "wal_wallet"
}

// WalletLog 钱包流水明细
type WalletLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 钱包编号
	WalletId int `json:"walletId" db:"wallet_id" gorm:"column:wallet_id" bson:"walletId"`
	// 业务类型
	Kind int `json:"kind" db:"kind" gorm:"column:kind" bson:"kind"`
	// 标题
	Subject string `json:"subject" db:"subject" gorm:"column:subject" bson:"subject"`
	// 外部通道
	OuterChan string `json:"outerChan" db:"outer_chan" gorm:"column:outer_chan" bson:"outerChan"`
	// 外部订单号
	OuterTxNo string `json:"outerTxNo" db:"outer_tx_no" gorm:"column:outer_tx_no" bson:"outerTxNo"`
	// 变动金额
	ChangeValue int `json:"changeValue" db:"change_value" gorm:"column:change_value" bson:"changeValue"`
	// 余额
	Balance int `json:"balance" db:"balance" gorm:"column:balance" bson:"balance"`
	// 交易手续费
	TransactionFee int `json:"transactionFee" db:"transaction_fee" gorm:"column:transaction_fee" bson:"transactionFee"`
	// 操作人员用户编号
	OprUid int `json:"oprUid" db:"opr_uid" gorm:"column:opr_uid" bson:"oprUid"`
	// 操作人员名称
	OprName string `json:"oprName" db:"opr_name" gorm:"column:opr_name" bson:"oprName"`
	// 提现账号
	AccountNo string `json:"accountNo" db:"account_no" gorm:"column:account_no" bson:"accountNo"`
	// 提现银行账户名称
	AccountName string `json:"accountName" db:"account_name" gorm:"column:account_name" bson:"accountName"`
	// 提现银行
	BankName string `json:"bankName" db:"bank_name" gorm:"column:bank_name" bson:"bankName"`
	// 审核状态
	ReviewStatus int `json:"reviewStatus" db:"review_status" gorm:"column:review_status" bson:"reviewStatus"`
	// 审核备注
	ReviewRemark string `json:"reviewRemark" db:"review_remark" gorm:"column:review_remark" bson:"reviewRemark"`
	// 审核时间
	ReviewTime int `json:"reviewTime" db:"review_time" gorm:"column:review_time" bson:"reviewTime"`
	// 备注
	Remark string `json:"remark" db:"remark" gorm:"column:remark" bson:"remark"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	// 钱包用户
	WalletUser string `json:"walletUser" db:"wallet_user" gorm:"column:wallet_user" bson:"walletUser"`
	// 交易外部用户
	OuterTxUid int `json:"outerTxUid" db:"outer_tx_uid" gorm:"column:outer_tx_uid" bson:"outerTxUid"`
}

func (w WalletLog) TableName() string {
	return "wal_wallet_log"
}
