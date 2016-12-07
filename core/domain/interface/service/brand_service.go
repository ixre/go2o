package service

// 产品品牌
type ProBrand struct {
	// 编号
	Id int64 `db:"id" pk:"yes" auto:"yes"`
	// 品牌名称
	Name string `db:"name"`
	// 品牌图片
	Image string `db:"image"`
	// 品牌网址
	SiteUrl string `db:"site_url"`
	// 介绍
	Intro string `db:"intro"`
	// 是否审核
	Review bool `db:"review"`
	// 加入时间
	CreateTime int64 `db:"create_time"`
}
