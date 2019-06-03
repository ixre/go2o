/**
 * Copyright 2015 @ to2.net.
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
	ItemId     int64 `db:"id"`
	ProductId  int64 `db:"product_id"`
	VendorId   int32 `db:"-"`
	ShopId     int32 `db:"-"`
	CategoryId int32 `db:"cat_id"`

	// 过渡的标题
	Name       string `db:"title"`
	Title      string `db:"title"`
	ShortTitle string `db:"-"`
	// 货号
	GoodsNo string `db:"code"`
	Image   string `db:"image"`

	//定价
	RetailPrice float32 `db:"retail_price"`

	//销售价
	Price float32 `db:"price"`

	// 促销价
	PromPrice float32 `db:"-"`
	// 价格区间
	PriceRange string `db:"price_range"`

	GoodsId   int64 `db:"it.id"`
	SkuId     int64 `db:"sku_id"`
	IsPresent int32 `db:"is_present"`

	// 促销标志
	PromotionFlag int32 `db:"prom_flag"`

	// 库存
	StockNum int32 `db:"stock_num"`

	// 已售件数
	SaleNum int32 `db:"sale_num"`
}
