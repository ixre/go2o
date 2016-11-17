/**
 * Copyright 2015 @ at3.net.
 * name : model.go
 * author : jarryliu
 * date : 2016-11-15 19:37
 * description :
 * history :
 */
package dao

type (
	// 二维码模板
	CommQrTemplate struct {
		// 编号
		Id int64 `db:"id" pk:"yes" auto:"yes"`
		// 模板标题
		Title string `db:"title"`
		// 背景图片
		BgImage string `db:"bg_image"`
		// 垂直偏离量
		OffsetX int `db:"offset_x"`
		// 垂直偏移量
		OffsetY int `db:"offset_y"`
		// 二维码模板文本
		Comment string `db:"comment"`
		// 回调地址
		CallbackUrl string `db:"callback_url"`
		// 是否启用
		Enabled int `db:"enabled"`
	}
)
