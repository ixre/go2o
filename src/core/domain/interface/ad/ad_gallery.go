/**
 * Copyright 2015 @ S1N1 Team.
 * name : ad_gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

type IGalleryAd interface {
	// 获取广告值
	GetAdValue()*ValueGallery

	// 设置广告值
	SetAdValue(*ValueGallery)error

	// 保存广告
	Save()(int,error)
}