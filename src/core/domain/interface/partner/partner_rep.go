/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:52
 * description :
 * history :
 */

package partner

type IPartnerRep interface {
	CreatePartner(*ValuePartner) (IPartner, error)

	// 获取商户的编号
	GetPartnersId() []int

	GetPartner(int) (IPartner, error)

	// 获取合作商主要的域名主机
	GetPartnerMajorHost(int) string

	// 保存
	SavePartner(*ValuePartner) (int, error)

	// 获取销售配置
	GetSaleConf(int) *SaleConf

	SaveSaleConf(partnerId int, v *SaleConf) error

	// 获取站点配置
	GetSiteConf(int) *SiteConf

	SaveSiteConf(partnerId int, v *SiteConf) error

	// 保存API信息
	SaveApiInfo(partnerId int, d *ApiInfo) error

	// 获取API信息
	GetApiInfo(partnerId int) *ApiInfo

	// 根据API编号获取商户编号
	GetPartnerIdByApiId(apiId string) int

	SaveShop(*ValueShop) (int, error)

	GetValueShop(partnerId, shopId int) *ValueShop

	GetShopsOfPartner(partnerId int) []*ValueShop

	DeleteShop(partnerId, shopId int) error

	// 获取键值
	GetKeyValue(partnerId int, k string) string
	// 设置键值
	SaveKeyValue(partnerId int, k, v string) error
	// 获取多个键值
	GetKeyMap(partnerId int, k []string) map[string]string
	// 检查是否包含值的键数量,keyStr为键模糊匹配
	CheckKvContainValue(partnerId string,value string,keyStr string)int
}
