/**
 * Copyright 2015 @ S1N1 Team.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:48
 * description :
 * history :
 */
package member


const(
	// 退款
	TypeBalanceChargeBack = 0
	// 系统充值
	TypeBalanceSystemCharge = 1
	// 网银充值
	TypeBalanceNetPayCharge = 2
	// 客服充值
	TypeBalanceServiceCharge = 3
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
	SaveBalanceInfo(*BalanceInfoValue)(int,error)

	// 充值
	// @title 充值标题说明
	// @no    充值订单编号
	// @amount 金额
	Charge(chargeType int,title string,no string,amount float32)(error)

	// 退款
	ChargeBack(title string,no string,amount float32)(error)

	// 请求提现
	RequestApplyCash(amount float32)(error)

	// 确认提现
	ConfirmApplyCash(id int)(error)

	// 完成提现
	FinishApplyCash(id int)(error)
}
