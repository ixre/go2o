/**
 * Copyright 2015 @ S1N1 Team.
 * name : goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package valueobject

// g.id,g.item_id,g.sku_id,g.is_present,g.prom_flag,g.stock_num,g.sale_num,
//i.category_id,i.name as name,i.goods_no,i.img,i.price,i.sale_price

// 完整的商品信息
type Goods struct {
	Item_Id    int    `db:"gs_goods.item_id"`
	CategoryId int    `db:"gs_item.category_id"`
	Name       string `db:"gs_item.name"`
	// 货号
	GoodsNo string `db:"gs_item.goods_no"`
	Image   string `db:"gs_item.img"`

	//定价
	Price float32 `db:"gs_item.price"`

	//销售价
	SalePrice float32 `db:"gs_item.sale_price"`

	// 促销价
	PromPrice float32 `db:"-"`

	GoodsId   int `db:"gs_goods.id"`
	SkuId     int `db:"gs_goods.sku_id"`
	IsPresent int `db:"gs_goods.is_present"`

	// 促销标志
	PromotionFlag int `db:"gs_goods.prom_flag"`

	// 库存
	StockNum int `db:"gs_goods.stock_num"`

	// 已售件数
	SaleNum int `db:"gs_goods.sale_num"`
}

