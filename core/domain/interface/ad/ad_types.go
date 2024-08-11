/**
 * Copyright 2015 @ 56x.net.
 * name : ad_gallery
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

import "sort"

type (
	// 文本广告
	IHyperLinkAd interface {
		SetData(v *Data) error
	}
	// 图片广告
	IImageAd interface {
		SetData(v *Data) error
	}

	// 广告图片
	Data struct {
		// 图片编号
		Id int64 `db:"id" auto:"yes" pk:"yes"`

		// 广告编号
		AdId int64 `db:"ad_id"`

		// 图片标题
		Title string `db:"title"`

		// 链接
		LinkUrl string `db:"link_url"`

		// 图片地址
		ImageUrl string `db:"image_url"`

		// 排列序号
		SortNum int `db:"sort_num"`

		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 轮播广告
	IGalleryAd interface {
		// 获取广告数据
		GetAdValue() SwiperAd

		// 获取可用的广告数据
		GetEnabledAdValue() SwiperAd

		// 保存广告图片
		SaveImage(v []*Data) error

		// 获取图片项
		GetImage(id int64) *Data
	}
)

var _ sort.Interface = SwiperAd{}

// 轮播广告图片集合
type SwiperAd []*Data

func (v SwiperAd) Len() int {
	return len(v)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (v SwiperAd) Less(i, j int) bool {
	return v[i].SortNum < v[j].SortNum || (v[i].SortNum == v[j].SortNum &&
		v[i].Id < v[j].Id)
}

// Swap swaps the elements with indexes i and j.
func (v SwiperAd) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
