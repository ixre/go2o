/**
 * Copyright 2015 @ S1N1 Team.
 * name : value_text
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

type ValueText struct {
	// 编号
	Id int `db:"id" auto:"yes" pk:"true"`

	// 广告编号
	AdvertisementId int `db:"ad_id"`

	// 标题
	Title string `db:"title"`

	// 链接
	LinkUrl string `db:"link_url"`

	// 是否启用
	Enabled int `db:"enabled"`
}
