/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-09 09:50
 * description :
 * history :
 */

package member

import (
	"github.com/ixre/go2o/core/dto"
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
	GetMember(memberId int64) IMemberAggregateRoot

	// 创建会员
	CreateMember(*Member) IMemberAggregateRoot

	// 删除会员
	DeleteMember(memberId int64) error

	// 创建会员,仅作为某些操作使用,不保存
	CreateMemberById(memberId int64) IMemberAggregateRoot

	// SaveMember 保存
	SaveMember(v *Member) (int64, error)
	// ResetMemberIdCache 重置会员缓存
	ResetMemberIdCache(field string, value string) error

	// 获取会员最后更新时间
	GetMemberLatestUpdateTime(id int64) int64

	// 根据手机号获取会员编号
	GetMemberIdByPhone(phone string) int64

	// 根据邮箱地址获取会员编号
	GetMemberIdByEmail(email string) int64

	// 获取会员编号
	GetMemberIdByUser(user string) int64

	// 根据编码获取会员
	GetMemberIdByCode(code string) int64

	// CheckUserExist 用户名是否存在
	CheckUserExist(user string, memberId int64) bool
	// CheckNicknameIsUse 昵称是否使用
	CheckNicknameIsUse(nickname string, memberId int64) bool
	// CheckPhoneBind 手机号码是否使用
	CheckPhoneBind(phone string, memberId int64) bool

	// 保存绑定
	SaveRelation(*InviteRelation) error

	// 获取账户
	GetAccount(memberId int64) *Account

	// 保存账户，传入会员编号
	SaveAccount(*Account) (int64, error)
	// 获取银行卡
	BankCards(memberId int64) []BankCard
	// 保存银行卡信息
	SaveBankCard(card *BankCard) error
	// 移除银行卡
	RemoveBankCard(id int64, no string) error
	// 获取收款码
	ReceiptsCodes(memberId int64) []ReceiptsCode
	// 保存收款码
	SaveReceiptsCode(code *ReceiptsCode, memberId int64) (int, error)
	// 保存积分记录
	SaveIntegralLog(*IntegralLog) error
	// 获取积分记录
	GetIntegralLog(id int) *IntegralLog
	// 保存余额日志
	SaveBalanceLog(*BalanceLog) (int32, error)
	// 获取余额日志
	GetBalanceLog(id int) *BalanceLog
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
	// 保存实名信息
	SaveTrustedInfo(id int, v *CerticationInfo) (int, error)
	// 获取实名信息
	GetTrustedInfo(memberId int) *CerticationInfo
	// 获取经验值对应的等级
	GetLevelValueByExp(mchId int64, exp int64) int

	// 获取会员升级记录
	GetLevelUpLog(id int) *LevelUpLog

	// 保存会员升级记录
	SaveLevelUpLog(l *LevelUpLog) (int32, error)

	// 保存地址
	SaveDeliverAddress(*ConsigneeAddress) (int64, error)

	// 获取全部配送地址
	GetDeliverAddress(memberId int64) []*ConsigneeAddress

	// 获取配送地址
	GetSingleDeliverAddress(memberId, addressId int64) *ConsigneeAddress

	// 删除配送地址
	DeleteAddress(memberId, addressId int64) error

	// 邀请
	GetMyInvitationMembers(memberId int64, begin, end int) (total int, rows []*dto.InvitationMember)

	// 获取下级会员数量
	GetSubInvitationNum(memberId int64, memberIdArr []int32) map[int32]int

	// GetInvitationCount 获取邀请会员数量
	GetInvitationCount(memberId int, level int) int

	// 获取推荐我的人
	GetInvitationMeMember(memberId int64) *Member

	// 保存余额变动信息
	SaveFlowAccountInfo(v *FlowAccountLog) (int32, error)

	// 保存理财账户信息
	SaveGrowAccount(memberId int64, balance, totalAmount,
		growEarnings, totalGrowEarnings float32, updateTime int64) error

	//收藏,favType 为收藏类型, referId为关联的ID
	Favorite(memberId int64, favType int, referId int64) error
	//是否已收藏
	Favored(memberId int64, favType int, referId int64) bool
	//取消收藏
	CancelFavorite(memberId int64, favType int, referId int64) error
	// 获取会员分页的优惠券列表
	GetMemberPagedCoupon(memberId int64, start, end int, where string) (total int, rows []*dto.SimpleCoupon)
	// Select MmBuyerGroup
	SelectMmBuyerGroup(where string, v ...interface{}) []*BuyerGroup
	// Save MmBuyerGroup
	SaveMmBuyerGroup(v *BuyerGroup) (int, error)

	// Save 会员锁定历史
	SaveLockHistory(v *MmLockHistory) (int, error)
	// Save 会员锁定记录
	SaveLockInfo(v *MmLockInfo) (int, error)
	// Delete 会员锁定记录
	DeleteLockInfos(memberId int64) error
	// 注册解锁任务
	RegisterUnlockJob(info *MmLockInfo)
	// 获取会员邀请的会员编号列表
	GetInviteChildren(id int64) []int64

	// GetOAuthAccount 关联第三方应用账号
	GetOAuthAccount(memberId int, appCode string) *OAuthAccount
	// SaveOAuthAccount 关联第三方应用账号
	SaveOAuthAccount(v *OAuthAccount) (int, error)
	// DeleteOAuthAccount 关联第三方应用账号
	DeleteOAuthAccount(primary interface{}) error
}
