/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:47
 * description :
 * history :
 */

package sale

// 商品值
type ValueItem struct {
	Id         int    `db:"id" auto:"yes" pk:"yes"`
	CategoryId int    `db:"category_id"`
	Name       string `db:"name"`

	// 货号
	GoodsNo    string `db:"goods_no"`
	SmallTitle string `db:"small_title"`
	Image      string `db:"img"`
	//成本价
	Cost float32 `db:"cost"`
	//定价
	Price float32 `db:"price"`
<<<<<<< HEAD
	//参考销售价
=======
	//销售价
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	SalePrice float32 `db:"sale_price"`
	ApplySubs string  `db:"apply_subs"`

	//简单备注,如:(限时促销)
	Note        string `db:"note"`
	Description string `db:"description"`

	// 是否上架,1为上架
	OnShelves int `db:"on_shelves"`

	State      int   `db:"state"`
	CreateTime int64 `db:"create_time"`
	UpdateTime int64 `db:"update_time"`
}
