/**
 * Copyright 2015 @ S1N1 Team.
 * name : account
 * author : jarryliu
 * date : 2015-07-24 08:48
 * description :
 * history :
 */
package member

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
}
