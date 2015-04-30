/**
 * Copyright 2014 @ ops Inc.
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

	// 初始化商户
	InitPartner(partnerId int) error

	// 获取销售配置
	GetSaleConf(int) *SaleConf

	SaveSaleConf(*SaleConf) error

	// 获取站点配置
	GetSiteConf(int) *SiteConf

	SaveSiteConf(*SiteConf) error

	SaveShop(*ValueShop) (int, error)

	GetValueShop(partnerId, shopId int) *ValueShop

	GetShopsOfPartner(partnerId int) []*ValueShop

	DeleteShop(partnerId, shopId int) error
}
