/**
 * Copyright 2015 @ z3q.net.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:48
 * description :
 * history :
 */
package member

const (
	// 退款
	KindBalanceBack = 0
	// 消费
	KindBalanceShopping = 1
	// 充值
	KindBalanceCharge = 2
	// 赠送
	KindBalancePresent = 3
	// 账户流通
	KindBalanceFlow = 4
	// 提现
	KindBalanceApplyCash = 5
	// 转账
	KindBalanceTransfer = 6
	// 冻结
	KindBalanceFreezes = 7
	// 解冻
	KindBalanceUnfreezes = 8
	// 冻结赠款
	KindBalanceFreezesPresent = 9
	// 解冻赠款
	KindBalanceUnfreezesPresent = 10

	// 系统充值
	TypeBalanceSystemCharge = 1
	// 网银充值
	TypeBalanceNetPayCharge = 2
	// 客服充值
	TypeBalanceServiceCharge = 3

	// 提现并充值到余额
	TypeApplyCashToCharge = 1
	// 提现到银行卡
	TypeApplyCashToBank = 2
	// 提现到第三方服务提供商（如：Paypal,支付宝等)
	TypeApplyCashToServiceProvider = 3

	// 退款到银行卡
	TypeBackToBank = 1
	// 退款到第三方
	TypeBackToServiceProvider = 2

	// 提现请求已提交
	StateApplySubmitted = 0
	// 提现已经确认
	StateApplyConfirmed = 1
	// 提现未通过
	StateApplyNotPass = 2
	// 提现完成
	StateApplyOver = 3

	StatusNormal = 0
	StatusOK     = 1
)

type IAccount interface {
	// 获取领域对象编号
	GetDomainId() int

	// 获取账户值
	GetValue() *AccountValue

	// 保存
	Save() (int, error)

	// 根据编号获取余额变动信息
	GetBalanceInfo(id int) *BalanceInfoValue

	// 根据号码获取余额变动信息
	GetBalanceInfoByNo(no string) *BalanceInfoValue

	// 保存余额变动信息
	SaveBalanceInfo(*BalanceInfoValue) (int, error)

	// 充值
	// @title 充值标题说明
	// @no    充值订单编号
	// @amount 金额
	ChargeBalance(chargeType int, title string, tradeNo string, amount float32) error

	// 赠送金额
	PresentBalance(title string, tradeNo string, amount float32)error

	// 流通账户余额变动，如扣除,amount传入负数金额
	ChargeFlowBalance(title string,tradeNo string,amount float32)error

	// 订单抵扣消费
	OrderDiscount(tradeNo string, amount float32) error

	// 退款
	RequestBackBalance(backType int, title string, amount float32) error

	// 完成退款
	FinishBackBalance(id int, tradeNo string) error

	// 请求提现,applyType：提现方式
	RequestApplyCash(applyType int, title string, amount float32) error

	// 确认提现
	ConfirmApplyCash(id int, pass bool, remark string) error

	// 完成提现
	FinishApplyCash(id int, tradeNo string) error

	// 冻结余额
	Freezes(amount float32) error

	// 解冻金额
	Unfreezes(amount float32) error

	// 冻结赠送金额
	FreezesPresent(amount float32) error

	// 解冻赠送金额
	UnfreezesPresent(amount float32) error
}
