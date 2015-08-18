/**
 * Copyright 2015 @ z3q.net.
 * name : ad_gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

type IGalleryAd interface {
	// 获取广告数据
	GetAdValue() ValueGallery

	// 获取可用的广告数据
	GetEnabledAdValue() ValueGallery

	// 保存广告图片
	SaveImage(v *ValueImage) (int, error)

	// 获取图片项
	GetImage(id int) *ValueImage

	// 删除图片项
	DelImage(id int) error
}
