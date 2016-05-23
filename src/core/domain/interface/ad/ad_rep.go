/**
 * Copyright 2015 @ z3q.net.
 * name : ad_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

// 广告仓储
type IAdRep interface {
	// 获取广告管理器
	GetAdManager() IAdManager

	// 获取广告分组
	GetAdGroups() []*AdGroup

	// 删除广告组
	DelAdGroup(id int) error

	// 获取广告位
	GetAdPositionsByGroupId(adGroupId int) []*AdPosition

	// 删除广告位
	DelAdPosition(id int) error

	// 保存广告位
	SaveAdPosition(a *AdPosition) (int, error)

	// 保存
	SaveAdGroup(value *AdGroup) (int, error)

	// 设置用户的广告
	SetUserAd(adUserId, posId, adId int) error

	// 根据名称获取广告编号
	GetIdByName(merchantId int, name string) int

	// 保存广告值
	SaveAdvertisementValue(*ValueAdvertisement) (int, error)

	// 保存广告图片
	SaveAdImageValue(*ValueImage) (int, error)

	// 获取广告
	GetValueAdvertisement(id int) *ValueAdvertisement

	// 根据名称获取广告
	GetValueAdvertisementByName(merchantId int, name string) *ValueAdvertisement

	// 获取轮播广告
	GetValueGallery(advertisementId int) ValueGallery

	// 获取图片项
	GetValueAdImage(advertisementId, id int) *ValueImage

	// 删除图片项
	DelAdImage(advertisementId, id int) error

	// 删除广告
	DelAdvertisement(merchantId, advertisementId int) error

	// 删除广告的图片数据
	DelImageDataForAdvertisement(advertisementId int) error

	// 删除广告的文字数据
	DelTextDataForAdvertisement(advertisementId int) error
}
