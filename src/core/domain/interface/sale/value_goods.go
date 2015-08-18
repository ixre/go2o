/**
 * Copyright 2015 @ z3q.net.
 * name : value_sale_goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

type ValueGoods struct {
	Id int `db:"id" pk:"yes" auto:"yes"`

	// 货品编号
	ItemId int `db:"item_id"`

	// 是否为赠品
	IsPresent int `db:"is_present"`

	// 规格
	SkuId int `db:"sku_id"`

	// 促销标志
	PromotionFlag int `db:"prom_flag"`

	// 库存
	StockNum int `db:"stock_num"`

	// 已售件数
	SaleNum int `db:"sale_num"`

	// 销售价
	SalePrice float32 `db:"-"`

	// 促销价
	PromPrice float32 `db:"-"`

	// 实际价
	Price float32 `db:"-"`
}
