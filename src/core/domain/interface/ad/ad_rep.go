/**
 * Copyright 2015 @ S1N1 Team.
 * name : ad_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

// 广告仓储
type IAdvertisementRep interface {
	// 获取商户的广告管理
	GetPartnerAdvertisement(partnerId int) IPartnerAdvertisement

	// 根据名称获取广告编号
	GetIdByName(partnerId int, name string) int

	// 保存广告值
	SaveAdvertisementValue(*ValueAdvertisement) (int, error)

	// 保存广告图片
	SaveAdImageValue(*ValueImage) (int, error)

	// 获取广告
	GetValueAdvertisement(partnerId, id int) *ValueAdvertisement

	// 根据名称获取广告
	GetValueAdvertisementByName(partnerId int, name string) *ValueAdvertisement

	// 获取轮播广告
	GetValueGallery(advertisementId int) ValueGallery

	// 获取图片项
	GetValueAdImage(advertisementId, id int) *ValueImage

	// 删除图片项
	DelAdImage(advertisementId, id int) error

	// 删除广告
	DelAdvertisement(partnerId, advertisementId int) error

	// 删除广告的图片数据
	DelImageDataForAdvertisement(advertisementId int) error

	// 删除广告的文字数据
	DelTextDataForAdvertisement(advertisementId int) error
}
