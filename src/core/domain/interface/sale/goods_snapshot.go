/**
 * Copyright 2013 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-02-04 20:13
 * description :
 * history :
 */
package sale

// 商品快照
type GoodsSnapshot struct {
	Id           int    `db:"id" auto:"yes" pk:"yes"`
	Key          string `db:"snapshot_key"`
	GoodsId      int    `db:"goods_id"`
	GoodsName    string `db:"goods_name"`
	GoodsNo      string `db:"goods_no"`
	SmallTitle   string `db:"small_title"`
	CategoryName string `db:"category_name"`
	Image        string `db:"img"`

	//成本价
	Cost float32 `db:"cost"`

	//定价
	Price float32 `db:"price"`

	//销售价
	SalePrice  float32 `db:"sale_price"`
	CreateTime int64   `db:"create_time"`
}
