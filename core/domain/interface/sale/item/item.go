/**
 * Copyright 2015 @ z3q.net.
 * name : item
 * author : jarryliu
 * date : 2016-06-29 09:31
 * description :
 * history :
 */
package item

import "go2o/core/domain/interface/valueobject"

type (
	IItemRep interface {
		// 获取货品
		GetValueItem(itemId int) *Item

		// 根据id获取货品
		GetItemByIds(ids ...int) ([]*Item, error)

		SaveValueItem(*Item) (int, error)

		// 获取在货架上的商品
		GetPagedOnShelvesItem(supplierId int, catIds []int, start, end int) (total int, goods []*Item)

		// 获取货品销售总数
		GetItemSaleNum(supplierId int, id int) int

		// 删除货品
		DeleteItem(supplierId, goodsId int) error
	}

	// 商品值
	Item struct {
		Id         int    `db:"id" auto:"yes" pk:"yes"`
		CategoryId int    `db:"category_id"`
		Name       string `db:"name"`
		//供应商编号(暂时同mch_id)
		VendorId int `db:"supplier_id"`
		// 货号
		GoodsNo    string `db:"goods_no"`
		SmallTitle string `db:"small_title"`
		Image      string `db:"img"`
		//成本价
		Cost float32 `db:"cost"`
		//定价
		Price float32 `db:"price"`
		//参考销售价
		SalePrice float32 `db:"sale_price"`
		ApplySubs string  `db:"apply_subs"`

		//简单备注,如:(限时促销)
		Remark      string `db:"remark"`
		Description string `db:"description"`

		// 是否上架,1为上架
		OnShelves int `db:"on_shelves"`

		State      int   `db:"state"`
		CreateTime int64 `db:"create_time"`
		UpdateTime int64 `db:"update_time"`
	}
)

// 转换包含部分数据的产品值对象
func ParseToPartialValueItem(v *valueobject.Goods) *Item {
	return &Item{
		Id:         v.Item_Id,
		CategoryId: v.CategoryId,
		Name:       v.Name,
		GoodsNo:    v.GoodsNo,
		Image:      v.Image,
		Price:      v.Price,
		SalePrice:  v.SalePrice,
	}
}
