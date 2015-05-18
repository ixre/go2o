/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-12 16:53
 * description :
 * history :
 */

package partner

import (
	"go2o/src/core/domain/interface/partner/user"
)

type IPartner interface {
	GetAggregateRootId() int

	GetValue() ValuePartner

	SetValue(*ValuePartner) error

	// 保存
	Save() (int, error)

	// 获取商户的域名
	GetMajorHost() string

	// 获取销售配置
	GetSaleConf() SaleConf

	// 保存销售配置
	SaveSaleConf(*SaleConf) error

	// 获取站点配置
	GetSiteConf() SiteConf

	// 保存站点配置
	SaveSiteConf(*SiteConf) error

	// 获取API信息
	GetApiInfo() ApiInfo

	// 保存API信息
	SaveApiInfo(*ApiInfo) error

	// 新建商店
	CreateShop(*ValueShop) IShop

	// 获取所有商店
	GetShops() []IShop

	// 获取商店
	GetShop(int) IShop

	// 删除门店
	DeleteShop(shopId int) error

	// 返回用户服务
	UserManager() user.IUserManager

	// 返回设置服务
	ConfManager()IConfManager

	// 获取会员等级服务
	LevelManager() ILevelManager
}
