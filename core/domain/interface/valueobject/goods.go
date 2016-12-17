/**
 * Copyright 2015 @ z3q.net.
 * name : goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package valueobject

// g.id,g.product_id,g.sku_id,g.is_present,g.prom_flag,g.stock_num,g.sale_num,
//i.cat_id,i.name as name,i.code,i.img,i.price,i.sale_price

// 完整的商品信息
type Goods struct {
	ProductId  int32  `db:"item_info.product_id"`
	VendorId   int32  `db:"-"`
	ShopId     int32  `db:"-"`
	CategoryId int32  `db:"pro_product.cat_id"`
	Name       string `db:"pro_product.name"`
	// 货号
	GoodsNo string `db:"pro_product.code"`
	Image   string `db:"pro_product.img"`

	//定价
	Price float32 `db:"pro_product.price"`

	//销售价
	SalePrice float32 `db:"pro_product.sale_price"`

	// 促销价
	PromPrice float32 `db:"-"`

	GoodsId   int32 `db:"item_info.id"`
	SkuId     int32 `db:"item_info.sku_id"`
	IsPresent int32 `db:"item_info.is_present"`

	// 促销标志
	PromotionFlag int32 `db:"item_info.prom_flag"`

	// 库存
	StockNum int32 `db:"item_info.stock_num"`

	// 已售件数
	SaleNum int32 `db:"item_info.sale_num"`
}
