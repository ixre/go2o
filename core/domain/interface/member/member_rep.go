/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:50
 * description :
 * history :
 */

package member

import (
	"go2o/core/domain/interface/merchant"
)

type IMemberRep interface {

	// 获取管理服务
	GetManager() IMemberManager

	// 获取会员等级
	GetMemberLevels_New() []*Level

	// 获取等级对应的会员数
	GetMemberNumByLevel_New(id int) int

	// 删除会员等级
	DeleteMemberLevel_New(id int) error

	// 保存会员等级
	SaveMemberLevel_New(v *Level) (int, error)

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
	CheckUsrExist(usr string, memberId int) bool

	// 手机号码是否使用
	CheckPhoneBind(phone string, memberId int) bool

	// 保存绑定
	SaveRelation(*MemberRelation) error

	// 获取账户
	GetAccount(memberId int) *AccountValue

	// 保存账户，传入会员编号
	SaveAccount(*AccountValue) (int, error)

	// 获取银行信息
	GetBankInfo(int) *BankInfo

	// 保存银行信息
	SaveBankInfo(*BankInfo) error

	// 保存积分记录
	SaveIntegralLog(*IntegralLog) error

	// 获取会员关联
	GetRelation(memberId int) *MemberRelation

	// 获取经验值对应的等级
	GetLevelValueByExp(merchantId int, exp int) int

	// 保存地址
	SaveDeliver(*DeliverAddress) (int, error)

	// 获取全部配送地址
	GetDeliverAddress(memberId int) []*DeliverAddress

	// 获取配送地址
	GetSingleDeliverAddress(memberId, deliverId int) *DeliverAddress

	// 删除配送地址
	DeleteDeliver(memberId, deliverId int) error

	// 邀请
	GetMyInvitationMembers(memberId, begin, end int) (total int, rows []*ValueMember)

	// 获取下级会员数量
	GetSubInvitationNum(memberId int, memberIdArr []int) map[int]int

	// 获取推荐我的人
	GetInvitationMeMember(memberId int) *ValueMember

	// 根据编号获取余额变动信息
	GetBalanceInfo(id int) *BalanceInfoValue

	// 根据号码获取余额变动信息
	GetBalanceInfoByNo(tradeNo string) *BalanceInfoValue

	// 保存余额变动信息
	SaveBalanceInfo(v *BalanceInfoValue) (int, error)

	// 保存理财账户信息
	SaveGrowAccount(memberId int, balance, totalAmount,
		growEarnings, totalGrowEarnings float32, updateTime int64) error

	//todo:商户需重构的等级方法
	/************  商户需重构的等级方法  *************/

	//获取等级
	GetLevel(merchantId, levelValue int) *merchant.MemberLevel

	// 获取下一个等级
	GetNextLevel(merchantId, levelVal int) *merchant.MemberLevel

	// 获取会员等级
	GetMemberLevels(merchantId int) []*merchant.MemberLevel

	// 删除会员等级
	DeleteMemberLevel(merchantId, id int) error

	// 保存等级
	SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error)
}
