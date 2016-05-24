/**
 * Copyright 2015 @ z3q.net.
 * name : value_page.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package content

type ValuePage struct {
	// 编号
	Id int `db:"id" pk:"yes" auto:"yes"`

	// 商户编号
	MerchantId int `db:"merchant_id"`

	// 标题
	Title string `db:"title"`

	// 字符标识
	StrIndent string `db:"str_indent"`

	// 关键词
	KeyWord string `db:"keyword"`

	// 描述
	Description string `db:"description"`

	// 样式表地址
	CssPath string `db:"css_path"`

	// 内容
	Body string `db:"body"`

	// 修改时间
	UpdateTime int64 `db:"update_time"`

	// 是否启用
	Enabled int `db:"enabled"`
}
