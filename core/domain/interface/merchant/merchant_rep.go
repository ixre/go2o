/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:52
 * description :
 * history :
 */

package merchant

type IMerchantRep interface {
	// 获取商户管理器
	GetManager() IMerchantManager

	CreateMerchant(*Merchant) IMerchant

	// 创建会员申请商户密钥
	CreateSignUpToken(memberId int64) string

	// 根据商户申请密钥获取会员编号
	GetMemberFromSignUpToken(token string) int64

	// 获取商户的编号
	GetMerchantsId() []int64

	// 获取商户
	GetMerchant(id int64) IMerchant

	// 获取合作商主要的域名主机
	GetMerchantMajorHost(mchId int64) string

	// 保存
	SaveMerchant(*Merchant) (int64, error)

	// 获取账户
	GetAccount(mchId int64) *Account

	// 获取销售配置
	GetMerchantSaleConf(mchId int64) *SaleConf

	// 保存销售配置
	SaveMerchantSaleConf(v *SaleConf) error

	// 保存商户账户信息
	UpdateAccount(v *Account) error

	// 保存API信息
	SaveApiInfo(d *ApiInfo) error

	// 获取API信息
	GetApiInfo(mchId int64) *ApiInfo

	// 根据API编号获取商户编号
	GetMerchantIdByApiId(apiId string) int

	// 获取键值
	GetKeyValue(mchId int64, indent string, k string) string

	// 设置键值
	SaveKeyValue(mchId int64, indent string, k, v string, updateTime int64) error

	// 获取多个键值
	GetKeyMap(mchId int64, indent string, k []string) map[string]string

	// 检查是否包含值的键数量,keyStr为键模糊匹配
	CheckKvContainValue(mchId int64, indent string, value string, keyStr string) int

	// 根据关键字获取字典
	GetKeyMapByChar(mchId int64, indent string, keyword string) map[string]string

	//获取等级
	GetLevel(mchId, levelValue int64) *MemberLevel

	// 获取下一个等级
	GetNextLevel(mchId, levelVal int64) *MemberLevel

	// 获取会员等级
	GetMemberLevels(mchId int64) []*MemberLevel

	// 删除会员等级
	DeleteMemberLevel(mchId, id int64) error

	// 保存等级
	SaveMemberLevel(mchId int64, v *MemberLevel) (int64, error)

	//
	//  //修改线下支付利润
	//UpdateMechOfflineRate(id int, rate float32, return_rate float32) error
	////商户提现日志
	//SaveMachBlanceLog(v *BalanceLog) error
	////个人提现日志
	//SavePresionBlanceLog(v *member.PresentLog) error
	//
	//GetOfflineRate(id int) (float32, float32, error)
}
