/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:49
 * description :
 * history :
 */

package member

type IMember interface {
	GetAggregateRootId() int

	GetValue() ValueMember

	SetValue(*ValueMember) error

	// 获取账户
	GetAccount() *Account

	// 保存账户
	SaveAccount() error

	// 获取提现银行信息
	GetBank() BankInfo

	// 保存提现银行信息
	SaveBank(*BankInfo) error

	// 保存返现记录
	SaveIncomeLog(*IncomeLog) error

	//　保存积分记录
	SaveIntegralLog(*IntegralLog) error

	// 增加经验值
	AddExp(exp int) error

	//　增加积分
	// todo:partnerId 不需要
	AddIntegral(partnerId int, backType int, integral int, log string) error

	// 获取关联的会员
	GetRelation() *MemberRelation

	// 更新会员绑定
	SaveRelation(r *MemberRelation) error

	// 保存
	Save() (int, error)

	// 修改密码,旧密码可为空
	ModifyPassword(newPwd, oldPwd string) error

	// 用户是否已经存在
	UsrIsExist() bool

	// 创建配送地址
	CreateDeliver(*DeliverAddress) IDeliver

	// 获取配送地址
	GetDeliverAddrs() []IDeliver

	// 获取配送地址
	GetDeliver(int) IDeliver

	// 删除配送地址
	DeleteDeliver(int) error
}
