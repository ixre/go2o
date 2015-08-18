/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:50
 * description :
 * history :
 */

package member

import "go2o/src/core/domain/interface/valueobject"

type IMemberRep interface {

	// 根据用户名获取会员
	GetMemberValueByUsr(usr string) *ValueMember

	// 根据手机号码获取会员
	GetMemberValueByPhone(phone string) *ValueMember

	// 获取会员
	GetMember(memberId int) IMember

	// 创建会员
	CreateMember(*ValueMember) IMember

	// 保存
	SaveMember(v *ValueMember) (int, error)

	// 获取会员最后更新时间
	GetMemberLatestUpdateTime(int) int64

	// 锁定会员
	LockMember(id int, state int) error

	// 根据邀请码获取会员编号
	GetMemberIdByInvitationCode(string) int

	// 获取会员编号
	GetMemberIdByUser(string string) int

	// 用户名是否存在
	CheckUsrExist(string) bool

	// 保存绑定
	SaveRelation(*MemberRelation) error

	// 获取等级
	GetLevel(partnerId, levelValue int) *valueobject.MemberLevel

	// 获取下一个等级
	GetNextLevel(partnerId, levelVal int) *valueobject.MemberLevel

	// 获取会员等级
	GetMemberLevels(partnerId int) []*valueobject.MemberLevel

	// 删除会员等级
	DeleteMemberLevel(partnerId, id int) error

	// 保存等级
	SaveMemberLevel(partnerId int, v *valueobject.MemberLevel) (int, error)

	// 获取账户
	GetAccount(memberId int) *AccountValue

	// 保存账户，传入会员编号
	SaveAccount(*AccountValue) (int, error)

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
	GetLevelValueByExp(partnerId int, exp int) int

	// 保存地址
	SaveDeliver(*DeliverAddress) (int, error)

	// 获取全部配送地址
	GetDeliverAddress(memberId int) []*DeliverAddress

	// 获取配送地址
	GetSingleDeliverAddress(memberId, deliverId int) *DeliverAddress

	// 删除配送地址
	DeleteDeliver(memberId, deliverId int) error

	// 邀请
	GetMyInvitationMembers(memberId int) []*ValueMember

	// 获取下级会员数量
	GetSubInvitationNum(memberIds string) map[int]int

	// 获取推荐我的人
	GetInvitationMeMember(memberId int) *ValueMember

	// 根据编号获取余额变动信息
	GetBalanceInfo(id int) *BalanceInfoValue

	// 根据号码获取余额变动信息
	GetBalanceInfoByNo(tradeNo string) *BalanceInfoValue

	// 保存余额变动信息
	SaveBalanceInfo(v *BalanceInfoValue) (int, error)
}
