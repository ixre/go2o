/**
 * Copyright 2015 @ z3q.net.
 * name : ad_gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import "sort"

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

var _ sort.Interface = ValueGallery{}

// 轮播广告图片集合
type ValueGallery []*ValueImage

func (this ValueGallery) Len() int {
	return len(this)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (this ValueGallery) Less(i, j int) bool {
	return this[i].SortNumber < this[j].SortNumber || (this[i].SortNumber == this[j].SortNumber &&
		this[i].Id < this[j].Id)
}

// Swap swaps the elements with indexes i and j.
func (this ValueGallery) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
