/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:49
 * description :
 * history :
 */

package member

import "go2o/src/core/domain/interface/valueobject"

type IMember interface {
	// 获取聚合根编号
	GetAggregateRootId() int

	// 获取值
	GetValue() ValueMember

	// 邀请管理
	Invitation() IInvitationManager

	// 设置值
	SetValue(*ValueMember) error

	// 获取账户
	GetAccount() IAccount

	// 锁定会员
	Lock() error

	// 解锁会员
	Unlock() error

	// 获取提现银行信息
	GetBank() BankInfo

	// 保存提现银行信息
	SaveBank(*BankInfo) error

	//　保存积分记录
	SaveIntegralLog(*IntegralLog) error

	// 增加经验值
	AddExp(exp int) error

	// 获取等级
	GetLevel() *valueobject.MemberLevel

	//　增加积分
	// todo:partnerId 不需要
	AddIntegral(partnerId int, backType int, integral int, log string) error

	// 获取关联的会员
	GetRelation() *MemberRelation

	// 更新会员绑定
	SaveRelation(r *MemberRelation) error

	// 更换用户名
	ChangeUsr(string) error

	// 保存
	Save() (int, error)

	// 修改密码,旧密码可为空
	ModifyPassword(newPwd, oldPwd string) error

	// 修改交易密码，旧密码可为空
	ModifyTradePassword(newPwd, oldPwd string) error

	// 创建配送地址
	CreateDeliver(*DeliverAddress) (IDeliver, error)

	// 获取配送地址
	GetDeliverAddress() []IDeliver

	// 获取配送地址
	GetDeliver(int) IDeliver

	// 删除配送地址
	DeleteDeliver(int) error
}
