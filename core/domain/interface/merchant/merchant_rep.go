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
	CreateMerchant(*Merchant) (IMerchant, error)

	// 获取商户的编号
	GetMerchantsId() []int

	// 获取商户,并槛车商户状态,返回错误
	GetMerchant(int) (IMerchant, error)

	// 获取合作商主要的域名主机
	GetMerchantMajorHost(int) string

	// 保存
	SaveMerchant(*Merchant) (int, error)

	// 获取账户
	GetAccount(mchId int) *Account

	// 获取销售配置
	GetMerchantSaleConf(int) *SaleConf

	// 保存销售配置
	SaveMerchantSaleConf(v *SaleConf) error

	// 保存API信息
	SaveApiInfo(d *ApiInfo) error

	// 获取API信息
	GetApiInfo(merchantId int) *ApiInfo

	// 根据API编号获取商户编号
	GetMerchantIdByApiId(apiId string) int

	// 获取键值
	GetKeyValue(merchantId int, indent string, k string) string

	// 设置键值
	SaveKeyValue(merchantId int, indent string, k, v string, updateTime int64) error

	// 获取多个键值
	GetKeyMap(merchantId int, indent string, k []string) map[string]string

	// 检查是否包含值的键数量,keyStr为键模糊匹配
	CheckKvContainValue(merchantId int, indent string, value string, keyStr string) int

	// 根据关键字获取字典
	GetKeyMapByChar(merchantId int, indent string, keyword string) map[string]string

	//获取等级
	GetLevel(merchantId, levelValue int) *MemberLevel

	// 获取下一个等级
	GetNextLevel(merchantId, levelVal int) *MemberLevel

	// 获取会员等级
	GetMemberLevels(merchantId int) []*MemberLevel

	// 删除会员等级
	DeleteMemberLevel(merchantId, id int) error

	// 保存等级
	SaveMemberLevel(merchantId int, v *MemberLevel) (int, error)
}
