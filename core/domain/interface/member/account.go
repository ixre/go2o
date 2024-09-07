/**
 * Copyright 2015 @ 56x.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:48
 * description :
 * history :
 */
package member

import (
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/domain"
)

type AccountType int

const (
	// AccountBalance 余额账户
	AccountBalance AccountType = 1
	// AccountIntegral 积分账户
	AccountIntegral AccountType = 2
	// AccountWallet 钱包账户
	AccountWallet AccountType = 3
	// AccountFlow 流通金账户
	AccountFlow AccountType = 4
	// AccountGrow 增长金账户
	AccountGrow AccountType = 7
)

const (
	// KindCustom 自定义的业务类型
	KindCustom int = 30
	// KindCharge 会员充值
	KindCharge = wallet.KCharge
	// KindCarry 入账
	KindCarry = wallet.KCarry
	// KindConsume 消耗
	KindConsume = wallet.KConsume
	// KindAdjust 客服调整
	KindAdjust = wallet.KAdjust
	// KindDiscount 支付抵扣
	KindDiscount = wallet.KDiscount
	// KindRefund 退款
	KindRefund int = wallet.KRefund
	// KindExchange 兑换充值, 比如将钱包充值到余额
	KindExchange int = wallet.KExchange
	// KindTransferIn 转入
	KindTransferIn int = wallet.KTransfer
	// KindTransferOut 转出
	KindTransferOut int = wallet.KTransfer
	// KindExpired 失效
	KindExpired int = wallet.KExpired
	// KindFreeze 冻结
	KindFreeze int = wallet.KFreeze
	// KindUnfreeze 解冻
	KindUnfreeze int = wallet.KUnfreeze
)

const (
	StatusOK = 1
)

const (
	// TypeIntegralPresent 赠送
	TypeIntegralPresent = 1
	// TypeIntegralFreeze 积分冻结
	TypeIntegralFreeze = 3
	// TypeIntegralUnfreeze 积分解冻
	TypeIntegralUnfreeze = 4
	// TypeIntegralShoppingPresent 购物赠送
	TypeIntegralShoppingPresent = 5
	// TypeIntegralPaymentDiscount 支付抵扣
	TypeIntegralPaymentDiscount = 6
)

type (
	IAccount interface {
		domain.IDomain

		// GetValue 获取账户值
		GetValue() *Account

		// Save 保存
		Save() (int, error)

		// Wallet 电子钱包
		Wallet() wallet.IWallet

		// SetPriorityPay 设置优先(默认)支付方式, account 为账户类型
		SetPriorityPay(account AccountType, enabled bool) error

		// Charge 用户充值,金额放大100倍（应只充值钱包）
		Charge(account AccountType, title string, amount int, outerNo string, remark string) error

		// CarryTo 入账,review是否先冻结审核, transactionFee手续费; 返回日志ID
		CarryTo(account AccountType, d AccountOperateData, review bool, transactionFee int) (int, error)

		// ReviewCarryTo 审核入账
		ReviewCarryTo(account AccountType, transactionId int, pass bool, reason string) error

		// Consume 消耗
		Consume(account AccountType, title string, amount int, outerNo string, remark string) error

		// Adjust 客服调整
		Adjust(account AccountType, title string, amount int, remark string, relateUser int64) error

		// Discount 抵扣, 如果账户扣除后不存在为消耗,反之为抵扣(内部,购物时需要抵扣一部分)
		Discount(account AccountType, title string, amount int, outerNo string, remark string) error

		// Refund 退款
		Refund(account AccountType, title string, amount int, outerNo string, remark string) error

		// Freeze 账户冻结
		Freeze(account AccountType, p AccountOperateData, relateUser int64) (int, error)

		// Unfreeze 账户解冻, isRefundBalance 是否退回到余额
		Unfreeze(account AccountType, p AccountOperateData, isRefundBalance bool, relateUser int64) error

		// 预扣消费,将冻结转为消费,扣款后不自动退回余额
		PrefreezeConsume(transactionId int, transactionTitle string, transactionRemark string) error

		// FreezeExpired 将冻结金额标记为失效
		FreezeExpired(account AccountType, amount int, remark string) error

		// PaymentDiscount 支付单抵扣消费,tradeNo为支付单单号
		PaymentDiscount(tradeNo string, amount int, remark string) error

		// GetWalletLog 获取钱包账户日志
		GetWalletLog(id int64) wallet.WalletLog

		// RequestWithdrawal 申请提现(只支持钱包),drawType：提现方式,返回info_id,交易号 及错误
		RequestWithdrawal(w *wallet.WithdrawTransaction) (int, string, error)

		// ReviewWithdrawal 提现审核
		ReviewWithdrawal(transactionId int, pass bool, reason string) error

		// FinishWithdrawal 完成提现(打款),outerTransactionNo为外部交易号
		FinishWithdrawal(transactionId int, outerTransactionNo string) error

		// TransferAccount 转账
		TransferAccount(account AccountType, toMember int64, amount int,
			transactionFee int, remark string) error

		// ReceiveTransfer 接收转账
		ReceiveTransfer(account AccountType, fromMember int64, tradeNo string,
			amount int, remark string) error

		// TransferBalance 转账余额到其他账户
		TransferBalance(account AccountType, amount int, tradeNo string, toTitle, fromTitle string) error

		// TransferFlow 转账活动账户,kind为转账类型，如 KindBalanceTransfer等
		// commission手续费
		TransferFlow(kind int, amount int, commission float32, tradeNo string,
			toTitle string, fromTitle string) error

		// TransferFlowTo 将活动金转给其他人
		TransferFlowTo(memberId int64, kind int, amount int, commission float32,
			tradeNo string, toTitle string, fromTitle string) error
	}
	// 账户操作数据
	AccountOperateData struct {
		// 描述
		TransactionTitle string
		// 金额
		Amount int
		// 外部订单号
		OuterTransactionNo string
		// 备注
		TransactionRemark string
		// 交易流水编号,对冻结流水进行更新时,传递该参数
		TransactionId int
		// 关联的外部用户编号,可为空
		OuterTxUid int
	}

	// WalletAccountLog 钱包账户日志
	WalletAccountLog struct {
		Id int64 `db:"id" auto:"yes" pk:"yes"`
		// 会员编号
		MemberId int64 `db:"member_id"`
		// 外部单号
		OuterNo string `db:"outer_no"`
		// 业务类型
		Kind int `db:"kind"`
		// 标题
		Title string `db:"title"`
		// 金额
		Amount int64 `db:"amount"`
		// 手续费
		ProcedureFee int64 `db:"procedure_fee"`
		// 关联操作人,仅在客服操作时,记录操作人
		RelateUser int64 `db:"rel_user"`
		// 状态
		ReviewStatus int32 `db:"review_status"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// FlowAccountLog 活动账户日志信息(todo: 活动账户还在用,暂时不删除)
	FlowAccountLog struct {
		Id int64 `db:"id" auto:"yes" pk:"yes"`
		// 会员编号
		MemberId int64 `db:"member_id"`
		// 外部单号
		OuterNo string `db:"outer_no"`
		// 业务类型
		Kind int `db:"kind"`
		// 标题
		Title string `db:"subject"`
		// 金额
		Amount int64 `db:"change_value"`
		// 手续费
		CsnFee int64 `db:"procedure_fee"`
		// 引用编号
		RelateUser int64 `db:"rel_user"`
		// 审核状态
		ReviewStatus int `db:"review_status"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)

func (b *BalanceLog) TableName() string {
	return "mm_balance_log"
}

// MmBalanceLog 余额日志
type BalanceLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会员编号
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
	// 类型
	Kind int16 `json:"kind" db:"kind" gorm:"column:kind" bson:"kind"`
	// 标题
	Subject string `json:"subject" db:"subject" gorm:"column:subject" bson:"subject"`
	// 外部交易号
	OuterNo string `json:"outerNo" db:"outer_no" gorm:"column:outer_no" bson:"outerNo"`
	// 金额
	ChangeValue int `json:"changeValue" db:"change_value" gorm:"column:change_value" bson:"changeValue"`
	// 手续费
	ProcedureFee int `json:"procedureFee" db:"procedure_fee" gorm:"column:procedure_fee" bson:"procedureFee"`
	// 审核状态
	ReviewStatus int `json:"reviewStatus" db:"review_status" gorm:"column:review_status" bson:"reviewStatus"`
	// 关联用户
	RelateUser int `json:"relUser" db:"rel_user" gorm:"column:rel_user" bson:"relUser"`
	// 备注
	Remark string `json:"remark" db:"remark" gorm:"column:remark" bson:"remark"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	// 变动后的余额
	Balance int `json:"balance" db:"balance" gorm:"column:balance" bson:"balance"`
}

// MmIntegralLog 积分明细
type IntegralLog struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会员编号
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
	// 类型
	Kind int `json:"kind" db:"kind" gorm:"column:kind" bson:"kind"`
	// 标题
	Subject string `json:"subject" db:"subject" gorm:"column:subject" bson:"subject"`
	// 关联的编号
	OuterNo string `json:"outerNo" db:"outer_no" gorm:"column:outer_no" bson:"outerNo"`
	// 积分值
	ChangeValue int `json:"changeValue" db:"change_value" gorm:"column:change_value" bson:"changeValue"`
	// 备注
	Remark string `json:"remark" db:"remark" gorm:"column:remark" bson:"remark"`
	// 关联用户
	RelateUser int `json:"relateUser" db:"rel_user" gorm:"column:rel_user" bson:"relUser"`
	// 审核状态
	ReviewStatus int `json:"reviewStatus" db:"review_status" gorm:"column:review_status" bson:"reviewStatus"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	// 变动后的余额
	Balance int `json:"balance" db:"balance" gorm:"column:balance" bson:"balance"`
}

func (m IntegralLog) TableName() string {
	return "mm_integral_log"
}

// MmAccount 会员账户
type Account struct {
	// 会员编号
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" pk:"yes" bson:"memberId"`
	// 积分
	Integral int `json:"integral" db:"integral" gorm:"column:integral" bson:"integral"`
	// 冻结积分
	FreezeIntegral int `json:"freezeIntegral" db:"freeze_integral" gorm:"column:freeze_integral" bson:"freezeIntegral"`
	// 余额
	Balance int `json:"balance" db:"balance" gorm:"column:balance" bson:"balance"`
	// 冻结余额
	FreezeBalance int `json:"freezeBalance" db:"freeze_balance" gorm:"column:freeze_balance" bson:"freezeBalance"`
	// 失效的余额
	ExpiredBalance int `json:"expiredBalance" db:"expired_balance" gorm:"column:expired_balance" bson:"expiredBalance"`
	// 钱包余额
	WalletBalance int `json:"walletBalance" db:"wallet_balance" gorm:"column:wallet_balance" bson:"walletBalance"`
	// 冻结钱包余额,作废
	FreezeWallet int `json:"freezeWallet" db:"freeze_wallet" gorm:"column:freeze_wallet" bson:"freezeWallet"`
	// ,作废
	ExpiredWallet int `json:"expiredWallet" db:"expired_wallet" gorm:"column:expired_wallet" bson:"expiredWallet"`
	// TotalWalletAmount
	TotalWalletAmount int `json:"totalWalletAmount" db:"total_wallet_amount" gorm:"column:total_wallet_amount" bson:"totalWalletAmount"`
	// FlowBalance
	FlowBalance int `json:"flowBalance" db:"flow_balance" gorm:"column:flow_balance" bson:"flowBalance"`
	// GrowBalance
	GrowBalance int `json:"growBalance" db:"grow_balance" gorm:"column:grow_balance" bson:"growBalance"`
	// GrowAmount
	GrowAmount int `json:"growAmount" db:"grow_amount" gorm:"column:grow_amount" bson:"growAmount"`
	// GrowEarnings
	GrowEarnings int `json:"growEarnings" db:"grow_earnings" gorm:"column:grow_earnings" bson:"growEarnings"`
	// GrowTotalEarnings
	GrowTotalEarnings int `json:"growTotalEarnings" db:"grow_total_earnings" gorm:"column:grow_total_earnings" bson:"growTotalEarnings"`
	// 累计充值
	TotalCharge int `json:"totalCharge" db:"total_charge" gorm:"column:total_charge" bson:"totalCharge"`
	// 累计支付
	TotalPay int `json:"totalPay" db:"total_pay" gorm:"column:total_pay" bson:"totalPay"`
	// 累计消费
	TotalExpense int `json:"totalExpense" db:"total_expense" gorm:"column:total_expense" bson:"totalExpense"`
	// PriorityPay
	PriorityPay int `json:"priorityPay" db:"priority_pay" gorm:"column:priority_pay" bson:"priorityPay"`
	// UpdateTime
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	// 钱包代码
	WalletCode string `json:"walletCode" db:"wallet_code" gorm:"column:wallet_code" bson:"walletCode"`
	// 可开票金额
	InvoiceableAmount int `json:"invoiceableAmount" db:"invoiceable_amount" gorm:"column:invoiceable_amount" bson:"invoiceableAmount"`
}

func (m Account) TableName() string {
	return "mm_account"
}

// OauthAccount 关联第三方应用账号
type OAuthAccount struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会员ID
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
	// 应用代码,如wechat-mp
	AppCode string `json:"appCode" db:"app_code" gorm:"column:app_code" bson:"appCode"`
	// 第三方应用id
	OpenId string `json:"openId" db:"open_id" gorm:"column:open_id" bson:"openId"`
	// UnionId
	UnionId string `json:"unionId" db:"union_id" gorm:"column:union_id" bson:"unionId"`
	// AuthToken
	AuthToken string `json:"authToken" db:"auth_token" gorm:"column:auth_token" bson:"authToken"`
	// ProfilePhoto
	ProfilePhoto string `json:"profilePhoto" db:"profile_photo" gorm:"column:profile_photo" bson:"profilePhoto"`
	// UpdateTime
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (m OAuthAccount) TableName() string {
	return "mm_oauth_account"
}
