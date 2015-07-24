/**
 * Copyright 2015 @ S1N1 Team.
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
	// 提现
	KindBalanceApplyCash = 3

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
	Charge(chargeType int, title string, tradeNo string, amount float32) error

	// 订单抵扣消费
	OrderDiscount(tradeNo string, amount float32) error

	// 退款
	RequestBackBalance(backType int, title string, amount float32) error

	// 完成退款
	FinishBackBalance(id int, tradeNo string) error

	// 请求提现
	RequestApplyCash(applyType int, title string, amount float32) error

	// 确认提现
	ConfirmApplyCash(id int) error

	// 完成提现
	FinishApplyCash(id int, tradeNo string) error
}
