/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-09 09:50
 * description :
 * history :
 */

package member

type IMemberRep interface {

	// 根据用户名获取会员
	GetMemberValueByUsr(usr string) *ValueMember

	GetMember(memberId int) (IMember, error)

	// 创建会员
	CreateMember(*ValueMember) IMember

	// 获取等级
	GetLevel(levelVal int) *MemberLevel

	// 获取下一个等级
	GetNextLevel(levelVal int) *MemberLevel

	// 获取配送地址
	GetDeliverAddr(memberId, deliverAddId int) *DeliverAddress

	// 获取账户
	GetAccount(memberId int) *Account

	// 保存账户，传入会员编号
	SaveAccount(*Account) error

	// 获取银行信息
	GetBankInfo(int) *BankInfo

	// 保存银行信息
	SaveBankInfo(*BankInfo) error

	// 保存返现记录
	SaveIncomeLog(*IncomeLog) error

	// 保存积分记录
	SaveIntegralLog(*IntegralLog) error

	// 获取会员关联
	GetRelation(memberId int) *MemberRelation

	// 获取经验值对应的等级
	GetLevelByExp(exp int) int

	// 保存
	SaveMember(v *ValueMember) (int, error)

	// 用户名是否存在
	CheckUsrExist(string) bool

	// 保存绑定
	SaveRelation(*MemberRelation) error
}
