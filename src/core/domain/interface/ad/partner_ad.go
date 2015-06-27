/**
 * Copyright 2015 @ S1N1 Team.
 * name : adversment
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

type IPartnerAdvertisement interface {
	// 获取聚合根标识
	GetAggregateRootId() int

	// 初始化内置默认的广告
	InitInternalAdvertisements()

	// 删除广告
	DeleteAdvertisement(advertisementId int)error

	// 根据编号获取广告
	GetById(int)IAdvertisement

	// 根据名称获取广告
	GetByName(string)IAdvertisement

	// 创建广告对象
	CreateAdvertisement(*ValueAdvertisement)IAdvertisement
}