/**
 * Copyright 2015 @ 56x.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:48
 * description :
 * history :
 */
package member

import "github.com/ixre/go2o/core/domain/interface/wallet"

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
	KindCharge = 1
	// KindCarry 入账
	KindCarry = 2
	// KindConsume 消耗
	KindConsume = 3
	// KindAdjust 客服调整
	KindAdjust = 4
	// KindDiscount 支付抵扣
	KindDiscount = 5
	// KindRefund 退款
	KindRefund int = 6
	// KindExchange 兑换充值, 比如将钱包充值到余额
	KindExchange int = 7
	// KindTransferIn 转入
	KindTransferIn int = 8
	// KindTransferOut 转出
	KindTransferOut int = 9
	// KindExpired 失效
	KindExpired int = 10
	// KindFreeze 冻结
	KindFreeze int = 11
	// KindUnfreeze 解冻
	KindUnfreeze int = 12
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
		// GetDomainId 获取领域对象编号
		GetDomainId() int64

		// GetValue 获取账户值
		GetValue() *Account

		// Save 保存
		Save() (int64, error)

		// Wallet 电子钱包
		Wallet() wallet.IWallet

		// SetPriorityPay 设置优先(默认)支付方式, account 为账户类型
		SetPriorityPay(account AccountType, enabled bool) error

		// Charge 用户充值,金额放大100倍（应只充值钱包）
		Charge(account AccountType, title string, amount int, outerNo string, remark string) error

		// CarryTo 入账,freeze是否先冻结, procedureFee手续费; 返回日志ID
		CarryTo(account AccountType, d AccountOperateData, freeze bool, procedureFee int) (int, error)

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

		// Unfreeze 账户解冻
		Unfreeze(account AccountType, p AccountOperateData, relateUser int64) error

		// PaymentDiscount 支付单抵扣消费,tradeNo为支付单单号
		PaymentDiscount(tradeNo string, amount int, remark string) error

		// FreezeExpired 将冻结金额标记为失效
		FreezeExpired(account AccountType, amount int, remark string) error

		// GetWalletLog 获取钱包账户日志
		GetWalletLog(id int64) wallet.WalletLog

		// RequestWithdrawal 申请提现,applyType：提现方式,返回info_id,交易号 及错误
		RequestWithdrawal(takeKind int, title string, amount int,
			tradeFee int, bankAccountNo string) (int64, string, error)

		// ReviewWithdrawal 确认提现
		ReviewWithdrawal(id int64, pass bool, remark string) error

		// FinishWithdrawal 完成提现
		FinishWithdrawal(id int64, tradeNo string) error

		// TransferAccount 转账
		TransferAccount(account AccountType, toMember int64, amount int,
			tradeFee int, remark string) error

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

	// Account 账户值对象
	Account struct {
		// 会员编号
		MemberId int64 `db:"member_id" pk:"yes"`
		// 积分
		Integral int `db:"integral"`
		// 不可用积分
		FreezeIntegral int `db:"freeze_integral"`
		// 余额
		Balance int64 `db:"balance"`
		// 不可用余额
		FreezeBalance int64 `db:"freeze_balance"`
		// 失效的账户余额
		ExpiredBalance int64 `db:"expired_balance"`
		// 钱包代码
		WalletCode string `db:"wallet_code"`
		//奖金账户余额
		WalletBalance int64 `db:"wallet_balance"`
		//冻结赠送金额
		FreezeWallet int64 `db:"freeze_wallet"`
		//失效的赠送金额
		ExpiredWallet int64 `db:"expired_wallet"`
		//总赠送金额
		TotalWalletAmount int64 `db:"total_wallet_amount"`
		//流动账户余额
		FlowBalance int64 `db:"flow_balance"`
		//当前理财账户余额
		GrowBalance int64 `db:"grow_balance"`
		//理财总投资金额,不含收益
		GrowAmount int64 `db:"grow_amount"`
		//当前收益金额
		GrowEarnings int64 `db:"grow_earnings"`
		//累积收益金额
		GrowTotalEarnings int64 `db:"grow_total_earnings"`
		//总消费金额
		TotalExpense int64 `db:"total_expense"`
		//总充值金额
		TotalCharge int64 `db:"total_charge"`
		//总支付额
		TotalPay int64 `db:"total_pay"`
		// 优先(默认)支付选项
		PriorityPay int `db:"priority_pay"`
		//更新时间
		UpdateTime int64 `db:"update_time"`
	}

	AccountOperateData struct {
		Title   string
		Amount  int
		OuterNo string
		Remark  string
	}

	// IntegralLog 积分记录
	IntegralLog struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 会员编号
		MemberId int `db:"member_id"`
		// 类型
		Kind int `db:"kind"`
		// 标题
		Title string `db:"subject"`
		// 关联的编号
		OuterNo string `db:"outer_no"`
		// 积分值
		Value int `db:"change_value"`
		// 余额
		Balance int `db:"balance"`
		// 备注
		Remark string `db:"remark"`
		// 关联用户
		RelateUser int `db:"rel_user"`
		// 审核状态
		ReviewState int16 `db:"audit_state"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// BalanceLog 余额日志
	BalanceLog struct {
		Id       int64  `db:"id" auto:"yes" pk:"yes"`
		MemberId int64  `db:"member_id"`
		OuterNo  string `db:"outer_no"`
		// 业务类型
		Kind int `db:"kind"`

		Title string `db:"subject"`
		// 金额
		Amount int64 `db:"change_value"`
		// 余额
		Balance int `db:"balance"`
		// 手续费
		ProcedureFee int64 `db:"procedure_fee"`
		// 关联操作人,仅在客服操作时,记录操作人
		RelateUser int64 `db:"rel_user"`
		// 状态
		ReviewState int32 `db:"audit_state"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
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
		ProcedureFee int64 `db:"csn_fee"`
		// 关联操作人,仅在客服操作时,记录操作人
		RelateUser int64 `db:"rel_user"`
		// 状态
		ReviewState int32 `db:"review_state"`
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
		ReviewState int `db:"audit_state"`
		// 备注
		Remark string `db:"remark"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 更新时间
		UpdateTime int64 `db:"update_time"`
	}
)
