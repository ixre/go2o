/**
 * Copyright 2015 @ z3q.net.
 * name : value_gallery_ad.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

// 广告图片
type ValueImage struct {
	// 图片编号
	Id int `db:"id" auto:"yes" pk:"true"`

	// 广告编号
	AdvertisementId int `db:"ad_id"`

	// 图片标题
	Title string `db:"title"`

	// 链接
	LinkUrl string `db:"link_url"`

	// 图片地址
	ImageUrl string `db:"image_url"`

	// 排列序号
	SortNumber int `db:"sort_number"`

	// 是否启用
	Enabled int `db:"enabled"`
}
