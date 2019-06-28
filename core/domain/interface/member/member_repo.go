/**
 * Copyright 2014 @ to2.net.
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
	GetProfile(memberId int64) *Profile

	// 保存资料
	SaveProfile(v *Profile) error

	// 获取会员等级
	GetMemberLevels_New() []*Level

	// 获取等级对应的会员数
	GetMemberNumByLevel_New(id int) int

	// 删除会员等级
	DeleteMemberLevel_New(id int) error

	// 保存会员等级
	SaveMemberLevel_New(v *Level) (int, error)

	// 根据用户名获取会员
	GetMemberByUser(user string) *Member

	// 根据手机号码获取会员
	GetMemberValueByPhone(phone string) *Member

	// 获取会员
	GetMember(memberId int64) IMember

	// 创建会员
	CreateMember(*Member) IMember

	// 删除会员
	DeleteMember(memberId int64) error

	// 创建会员,仅作为某些操作使用,不保存
	CreateMemberById(memberId int64) IMember

	// 保存
	SaveMember(v *Member) (int64, error)

	// 获取会员最后更新时间
	GetMemberLatestUpdateTime(id int64) int64

	// 根据邀请码获取会员编号
	GetMemberIdByInvitationCode(code string) int64

	// 根据手机号获取会员编号
	GetMemberIdByPhone(phone string) int64

	// 根据邮箱地址获取会员编号
	GetMemberIdByEmail(email string) int64

	// 获取会员编号
	GetMemberIdByUser(user string) int64

	// 根据编码获取会员
	GetMemberIdByCode(code string) int

	// 用户名是否存在
	CheckUsrExist(user string, memberId int64) bool

	// 手机号码是否使用
	CheckPhoneBind(phone string, memberId int64) bool

	// 保存绑定
	SaveRelation(*InviteRelation) error

	// 获取账户
	GetAccount(memberId int64) *Account

	// 保存账户，传入会员编号
	SaveAccount(*Account) (int64, error)

	// 获取银行信息
	GetBankInfo(memberId int64) *BankInfo

	// 保存银行信息
	SaveBankInfo(*BankInfo) error
	// 获取收款码
	GetCollectsCodes(memberId int64) []CollectsCode
	// 保存收款码
	SaveCollectsCode(code *CollectsCode, memberId int64) (int, error)
	// 保存积分记录
	SaveIntegralLog(*IntegralLog) error

	// 保存余额日志
	SaveBalanceLog(*BalanceLog) (int32, error)

	// 保存钱包账户日志
	SaveWalletAccountLog(*WalletAccountLog) (int32, error)

	// 获取钱包账户日志信息
	GetWalletLog(id int32) *WalletAccountLog

	// 增加会员当天提现次数
	AddTodayTakeOutTimes(memberId int64) error

	// 获取会员每日提现次数
	GetTodayTakeOutTimes(memberId int64) int

	// 获取会员关联
	GetRelation(memberId int64) *InviteRelation

	// 获取经验值对应的等级
	GetLevelValueByExp(mchId int32, exp int64) int

	// 获取会员升级记录
	GetLevelUpLog(id int32) *LevelUpLog

	// 保存会员升级记录
	SaveLevelUpLog(l *LevelUpLog) (int32, error)

	// 保存地址
	SaveDeliver(*Address) (int64, error)

	// 获取全部配送地址
	GetDeliverAddress(memberId int64) []*Address

	// 获取配送地址
	GetSingleDeliverAddress(memberId, addressId int64) *Address

	// 删除配送地址
	DeleteAddress(memberId, addressId int64) error

	// 邀请
	GetMyInvitationMembers(memberId int64, begin, end int) (total int, rows []*dto.InvitationMember)

	// 获取下级会员数量
	GetSubInvitationNum(memberId int64, memberIdArr []int32) map[int32]int

	// 获取推荐我的人
	GetInvitationMeMember(memberId int64) *Member

	// 保存余额变动信息
	SaveFlowAccountInfo(v *FlowAccountLog) (int32, error)

	// 保存理财账户信息
	SaveGrowAccount(memberId int64, balance, totalAmount,
		growEarnings, totalGrowEarnings float32, updateTime int64) error

	//收藏,favType 为收藏类型, referId为关联的ID
	Favorite(memberId int64, favType int, referId int32) error

	//是否已收藏
	Favored(memberId int64, favType int, referId int32) bool

	//取消收藏
	CancelFavorite(memberId int64, favType int, referId int32) error

	// 获取会员分页的优惠券列表
	GetMemberPagedCoupon(memberId int64, start, end int, where string) (total int, rows []*dto.SimpleCoupon)
	// Select MmBuyerGroup
	SelectMmBuyerGroup(where string, v ...interface{}) []*BuyerGroup
	// Save MmBuyerGroup
	SaveMmBuyerGroup(v *BuyerGroup) (int, error)
}
