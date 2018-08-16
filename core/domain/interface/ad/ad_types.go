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

type (
	// 文本广告
	IHyperLinkAd interface {
		SetData(v *HyperLink) error
	}

	//  超链接
	HyperLink struct {
		Id      int32  `db:"id"`
		AdId    int32  `db:"ad_id"`
		Title   string `db:"title"`
		LinkUrl string `db:"link_url"`
	}

	// 图片广告
	IImageAd interface {
		SetData(v *Image) error
	}

	// 广告图片
	Image struct {
		// 图片编号
		Id int32 `db:"id" auto:"yes" pk:"true"`

		// 广告编号
		AdId int32 `db:"ad_id"`

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
		GetAdValue() ValueGallery

		// 获取可用的广告数据
		GetEnabledAdValue() ValueGallery

		// 保存广告图片
		SaveImage(v *Image) (int32, error)

		// 获取图片项
		GetImage(id int32) *Image

		// 删除图片项
		DelImage(id int32) error
	}
)

var _ sort.Interface = ValueGallery{}

// 轮播广告图片集合
type ValueGallery []*Image

func (v ValueGallery) Len() int {
	return len(v)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (v ValueGallery) Less(i, j int) bool {
	return v[i].SortNum < v[j].SortNum || (v[i].SortNum == v[j].SortNum &&
		v[i].Id < v[j].Id)
}

// Swap swaps the elements with indexes i and j.
func (v ValueGallery) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
