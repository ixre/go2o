/**
 * Copyright 2015 @ z3q.net.
 * name : value_sale_tag
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

// 销售标签
type ValueSaleTag struct {
	Id int `db:"id" auto:"yes" pk:"yes"`

	// 商户编号
	PartnerId int `db:"partner_id"`

	// 标签代码
	TagCode string `db:"tag_code"`

	// 标签名
	TagName string `db:"tag_name"`

	// 商品的遮盖图
	GoodsImage string `db:"goods_image"`

	// 是否内部
	IsInternal int `db:"is_internal"`

	// 是否启用
	Enabled int `db:"enabled"`
}
