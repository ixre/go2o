/**
 * Copyright 2015 @ z3q.net.
 * name : goods_rep
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package repository

import (
	"fmt"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
)

var _ sale.IGoodsRep = new(goodsRep)

type goodsRep struct {
	db.Connector
}

// 商品仓储
func NewGoodsRep(c db.Connector) sale.IGoodsRep {
	return &goodsRep{
		Connector: c,
	}
}

// 获取商品
func (this *goodsRep) GetValueGoods(itemId int, skuId int) *sale.ValueGoods {
	var e *sale.ValueGoods = new(sale.ValueGoods)
	if this.Connector.GetOrm().GetBy(e, "item_id=? AND sku_id=?", itemId, skuId) == nil {
		return e
	}
	return nil
}

// 获取商品
func (this *goodsRep) GetValueGoodsById(goodsId int) *sale.ValueGoods {
	var e *sale.ValueGoods = new(sale.ValueGoods)
	if this.Connector.GetOrm().Get(goodsId, e) == nil {
		return e
	}
	return nil
}

// 根据SKU获取商品
func (this *goodsRep) GetValueGoodsBySku(itemId, sku int) *sale.ValueGoods {
	var e *sale.ValueGoods = new(sale.ValueGoods)
	if this.Connector.GetOrm().GetBy(e, "item_id=? AND sku_id=?", itemId, sku) == nil {
		return e
	}
	return nil
}

// 根据编号获取商品
func (this *goodsRep) GetGoodsByIds(ids ...int) ([]*valueobject.Goods, error) {
	var items []*valueobject.Goods
	err := this.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_goods.item_id=gs_item.id
	 WHERE gs_goods.id IN (`+format.GetCategoryIdStr(ids)+`)`)

	return items, err
}

// 获取会员价
func (this *goodsRep) GetGoodsLevelPrice(goodsId int) []*sale.MemberPrice {
	var items []*sale.MemberPrice
	if this.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_member_price WHERE goods_id = ?`, goodsId) == nil {
		return items
	}
	return nil
}

// 保存会员价
func (this *goodsRep) SaveGoodsLevelPrice(v *sale.MemberPrice) (id int, err error) {

	if v.Id <= 0 {
		this.Connector.ExecScalar(`SELECT MAX(id) FROM gs_member_price where goods_id=? and level=?`,
			&v.Id, v.GoodsId, v.Level)
	}

	if v.Id > 0 {
		_, _, err = this.Connector.GetOrm().Save(v.Id, v)
		id = v.Id
	} else {
		_, _, err = this.Connector.GetOrm().Save(nil, v)
		if err == nil {
			err = this.Connector.ExecScalar(`SELECT MAX(id) FROM gs_member_price where goods_id=?`, &id, v.GoodsId)
		}
	}
	return id, err
}

// 移除会员价
func (this *goodsRep) RemoveGoodsLevelPrice(id int) error {
	return this.Connector.GetOrm().DeleteByPk(sale.MemberPrice{}, id)
}

// 保存商品
func (this *goodsRep) SaveValueGoods(v *sale.ValueGoods) (id int, err error) {
	if v.Id > 0 {
		_, _, err = this.Connector.GetOrm().Save(v.Id, v)
		id = v.Id
	} else {
		_, _, err = this.Connector.GetOrm().Save(nil, v)
		if err == nil {
			err = this.Connector.ExecScalar(`SELECT MAX(id) FROM gs_goods where items_id=?`, &id, v.ItemId)
		}
	}
	return id, err

}

// 获取已上架的商品
func (this *goodsRep) GetPagedOnShelvesGoods(mchId int, catIds []int, start, end int, where, orderBy string) (total int, e []*valueobject.Goods) {
	var sql string

	 catIdStr := ""
	if catIds != nil && len(catIds)> 0 {
		catIdStr = fmt.Sprintf(" AND gs_category.id IN (%s)",
			format.GetCategoryIdStr(catIds))
	}

	if len(where) != 0 {
		where = " AND " + where
	}
	if len(orderBy) != 0 {
		orderBy += ","
	}

	this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_goods
	 INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE (?<=0 OR supplier_id=?) %s AND gs_item.state=1
		 AND gs_item.on_shelves=1 %s`, catIdStr, where), &total,mchId, mchId)

	e = []*valueobject.Goods{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE (?<=0 OR supplier_id=?) %s AND gs_item.state=1
		 AND gs_item.on_shelves=1 %s ORDER BY %s update_time DESC LIMIT %d,%d`, catIdStr,
			where, orderBy, start, (end - start))

		this.Connector.GetOrm().SelectByQuery(&e, sql,mchId,mchId)
	}

	return total, e
}

// 获取指定数量已上架的商品
func (this *goodsRep) GetOnShelvesGoods(merchantId int, start, end int, sortBy string) []*valueobject.Goods {
	e := []*valueobject.Goods{}
	sql := fmt.Sprintf(`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE supplier_id=? AND gs_item.state=1
		 AND gs_item.on_shelves=1 ORDER BY %s,update_time DESC LIMIT ?,?`,
		sortBy)

	this.Connector.GetOrm().SelectByQuery(&e, sql, merchantId, start, (end - start))
	return e
}
