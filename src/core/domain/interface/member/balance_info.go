/**
 * Copyright 2015 @ z3q.net.
 * name : balance_info
 * author : jarryliu
 * date : 2015-07-24 09:42
 * description :
 * history :
 */
package member

// 账户余额变动信息
type IBalanceInfo interface {
	// 获取领域对象编号
	GetDomainId() int

	// 获取值
	GetValue() *BalanceInfoValue

	// 保存
	Save() (int, error)
}
