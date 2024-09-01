/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:52
 * description :
 * history :
 */

package merchant

import "github.com/ixre/go2o/core/infrastructure/fw"

type IMerchantRepo interface {
	// 仓储实现
	fw.Repository[Merchant]
	// 账单仓储
	BillRepo() fw.Repository[MerchantBill]
	// 商户结算仓储
	SettleRepo() fw.Repository[SettleConf]
	// 获取商户管理器
	GetManager() IMerchantManager
	// 创建商户
	CreateMerchant(*Merchant) IMerchantAggregateRoot

	// 获取商户的编号
	GetMerchantsId() []int32

	// 获取商户
	GetMerchant(id int) IMerchantAggregateRoot
	// 根据登录用户名获取商户
	GetMerchantByUsername(user string) IMerchantAggregateRoot
	// 获取合作商主要的域名主机
	GetMerchantMajorHost(mchId int) string

	// 保存
	SaveMerchant(*Merchant) (int, error)

	// 获取账户
	GetAccount(mchId int) *Account
	// 保存会员账户
	SaveAccount(a *Account) (int, error)

	// 获取销售配置
	GetMerchantSaleConf(mchId int64) *SaleConf

	// 保存销售配置
	SaveMerchantSaleConf(v *SaleConf) error
	// 保存API信息
	SaveApiInfo(d *ApiInfo) error

	// 获取API信息
	GetApiInfo(mchId int) *ApiInfo

	// 根据API编号获取商户编号
	GetMerchantIdByApiId(apiId string) int64

	// 获取键值
	GetKeyValue(mchId int, indent string, k string) string

	// 设置键值
	SaveKeyValue(mchId int, indent string, k, v string, updateTime int64) error

	// 获取多个键值
	GetKeyMap(mchId int, indent string, k []string) map[string]string

	// 检查是否包含值的键数量,keyStr为键模糊匹配
	CheckKvContainValue(mchId int, indent string, value string, keyStr string) int

	// 根据关键字获取字典
	GetKeyMapByChar(mchId int, indent string, keyword string) map[string]string

	//获取等级
	GetLevel(mchId, levelValue int32) *MemberLevel

	// 获取下一个等级
	GetNextLevel(mchId, levelVal int32) *MemberLevel

	// 获取会员等级
	GetMemberLevels(mchId int64) []*MemberLevel

	// 删除会员等级
	DeleteMemberLevel(mchId, id int32) error

	// 保存等级
	SaveMemberLevel(mchId int64, v *MemberLevel) (int32, error)

	// Get MchBuyerGroupSetting
	GetMchBuyerGroupByGroupId(mchId, groupId int32) *MchBuyerGroupSetting
	// Select MchBuyerGroupSetting
	SelectMchBuyerGroup(mchId int64) []*MchBuyerGroupSetting
	// Save MchBuyerGroupSetting
	SaveMchBuyerGroup(v *MchBuyerGroupSetting) (int, error)

	// auto generate by gof
	// Get MchTradeConf
	GetMchTradeConf(primary interface{}) *TradeConf
	// GetBy MchTradeConf
	GetMchTradeConfBy(where string, v ...interface{}) *TradeConf
	// Select MchTradeConf
	SelectMchTradeConf(where string, v ...interface{}) []*TradeConf
	// Save MchTradeConf
	SaveMchTradeConf(v *TradeConf) (int, error)
	// Delete MchTradeConf
	DeleteMchTradeConf(primary interface{}) error
	// 验证商户用户名是否存在
	CheckUserExists(user string, id int) bool
	// CheckMemberBind 验证会员是否绑定商户
	CheckMemberBind(memberId int, mchId int) bool

	//  //修改线下支付利润
	//UpdateMechOfflineRate(id int, rate float32, return_rate float32) error
	////商户提现日志
	//SaveMachBlanceLog(v *BalanceLog) error
	////个人提现日志
	//SavePresionBlanceLog(v *member.PresentLog) error
	//
	//GetOfflineRate(id int32) (float32, float32, error)
	// // 根据外部订单号查找账户日志
	// GetBalanceLogByOuterNo(outerTradeNo string) *BalanceLog

	// 根据会员编号获取商户
	GetMerchantByMemberId(memberId int) IMerchantAggregateRoot
	// 查找账户流水
	GetBalanceAccountLog(id int) *BalanceLog
	// 保存账户日志
	SaveBalanceAccountLog(*BalanceLog) (int, error)

	// SaveAuthenticate 保存商户认证信息
	SaveAuthenticate(v *Authenticate) (int, error)
	// GetAuthenticateBy 获取商户认证信息
	GetMerchantAuthenticate(mchId int, version int) *Authenticate
	// DeleteOthersAuthenticate 删除其他认证信息
	DeleteOthersAuthenticate(mchId int, id int) error
}
