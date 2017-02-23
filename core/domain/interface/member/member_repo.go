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
	"go2o/core/dto"
)

type IMemberRepo interface {
	// 获取管理服务
	GetManager() IMemberManager

	// 获取资料或初始化
	GetProfile(memberId int32) *Profile

	// 保存资料
	SaveProfile(v *Profile) error

	// 获取会员等级
	GetMemberLevels_New() []*Level

	// 获取等级对应的会员数
	GetMemberNumByLevel_New(id int32) int

	// 删除会员等级
	DeleteMemberLevel_New(id int32) error

	// 保存会员等级
	SaveMemberLevel_New(v *Level) (int32, error)

	// 根据用户名获取会员
	GetMemberByUsr(usr string) *Member

	// 根据手机号码获取会员
	GetMemberValueByPhone(phone string) *Member

	// 获取会员
	GetMember(memberId int32) IMember

	// 创建会员
	CreateMember(*Member) IMember

	// 删除会员
	DeleteMember(id int32) error

	// 创建会员,仅作为某些操作使用,不保存
	CreateMemberById(memberId int32) IMember

	// 保存
	SaveMember(v *Member) (int32, error)

	// 获取会员最后更新时间
	GetMemberLatestUpdateTime(id int32) int64

	// 根据邀请码获取会员编号
	GetMemberIdByInvitationCode(code string) int32

	// 根据手机号获取会员编号
	GetMemberIdByPhone(phone string) int32

	// 根据邮箱地址获取会员编号
	GetMemberIdByEmail(email string) int32

	// 获取会员编号
	GetMemberIdByUser(user string) int32

	// 用户名是否存在
	CheckUsrExist(usr string, memberId int32) bool

	// 手机号码是否使用
	CheckPhoneBind(phone string, memberId int32) bool

	// 保存绑定
	SaveRelation(*Relation) error

	// 获取账户
	GetAccount(memberId int32) *Account

	// 保存账户，传入会员编号
	SaveAccount(*Account) (int32, error)

	// 获取银行信息
	GetBankInfo(id int32) *BankInfo

	// 保存银行信息
	SaveBankInfo(*BankInfo) error

	// 保存积分记录
	SaveIntegralLog(*IntegralLog) error

	// 保存余额日志
	SaveBalanceLog(*BalanceLog) (int32, error)

	// 保存钱包账户日志
	SavePresentLog(*PresentLog) (int32, error)

	// 获取钱包账户日志信息
	GetPresentLog(id int32) *PresentLog

	// 增加会员当天提现次数
	AddTodayTakeOutTimes(memberId int32) error

	// 获取会员每日提现次数
	GetTodayTakeOutTimes(memberId int32) int

	// 获取会员关联
	GetRelation(memberId int32) *Relation

	// 获取经验值对应的等级
	GetLevelValueByExp(mchId int32, exp int64) int

	// 获取会员升级记录
	GetLevelUpLog(id int32) *LevelUpLog

	// 保存会员升级记录
	SaveLevelUpLog(l *LevelUpLog) (int32, error)

	// 保存地址
	SaveDeliver(*Address) (int32, error)

	// 获取全部配送地址
	GetDeliverAddress(memberId int32) []*Address

	// 获取配送地址
	GetSingleDeliverAddress(memberId, adressId int32) *Address

	// 删除配送地址
	DeleteAddress(memberId, adressId int32) error

	// 邀请
	GetMyInvitationMembers(memberId int32, begin, end int) (total int, rows []*dto.InvitationMember)

	// 获取下级会员数量
	GetSubInvitationNum(memberId int32, memberIdArr []int32) map[int32]int

	// 获取推荐我的人
	GetInvitationMeMember(memberId int32) *Member

	// 根据编号获取余额变动信息
	GetBalanceInfo(id int32) *BalanceInfo

	// 根据号码获取余额变动信息
	GetBalanceInfoByNo(tradeNo string) *BalanceInfo

	// 保存余额变动信息
	SaveBalanceInfo(v *BalanceInfo) (int32, error)

	// 保存理财账户信息
	SaveGrowAccount(memberId int32, balance, totalAmount,
		growEarnings, totalGrowEarnings float32, updateTime int64) error

	//收藏,favType 为收藏类型, referId为关联的ID
	Favorite(memberId int32, favType int, referId int32) error

	//是否已收藏
	Favored(memberId int32, favType int, referId int32) bool

	//取消收藏
	CancelFavorite(memberId int32, favType int, referId int32) error

	// 获取会员分页的优惠券列表
	GetMemberPagedCoupon(memberId int32, start, end int, where string) (total int, rows []*dto.SimpleCoupon)
	// Select MmBuyerGroup
	SelectMmBuyerGroup(where string, v ...interface{}) []*BuyerGroup
	// Save MmBuyerGroup
	SaveMmBuyerGroup(v *BuyerGroup) (int, error)
}
