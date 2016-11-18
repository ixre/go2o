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
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
)

var _ goods.IGoodsRep = new(goodsRep)

type goodsRep struct {
	db.Connector
	_saleRep sale.ISaleRep
}

// 商品仓储
func NewGoodsRep(c db.Connector) *goodsRep {
	return &goodsRep{
		Connector: c,
	}
}
func (g *goodsRep) SetSaleRep(saleRep sale.ISaleRep) {
	g._saleRep = saleRep
}

// 根据SKU-ID获取商品,SKU-ID为商品ID
func (g *goodsRep) GetGoodsBySKuId(skuId int32) interface{} {
	snap := g.GetLatestSnapshot(skuId)
	if snap != nil {
		return g._saleRep.GetSale(snap.VendorId).
			GoodsManager().GetGoods(skuId)
	}
	return nil
}

// 获取商品
func (g *goodsRep) GetValueGoods(itemId int32, skuId int32) *goods.ValueGoods {
	var e *goods.ValueGoods = new(goods.ValueGoods)
	if g.Connector.GetOrm().GetBy(e, "item_id=? AND sku_id=?", itemId, skuId) == nil {
		return e
	}
	return nil
}

// 获取商品
func (g *goodsRep) GetValueGoodsById(goodsId int32) *goods.ValueGoods {
	var e *goods.ValueGoods = new(goods.ValueGoods)
	if g.Connector.GetOrm().Get(goodsId, e) == nil {
		return e
	}
	return nil
}

// 根据SKU获取商品
func (g *goodsRep) GetValueGoodsBySku(itemId, sku int32) *goods.ValueGoods {
	var e *goods.ValueGoods = new(goods.ValueGoods)
	if g.Connector.GetOrm().GetBy(e, "item_id=? AND sku_id=?", itemId, sku) == nil {
		return e
	}
	return nil
}

// 根据编号获取商品
func (g *goodsRep) GetGoodsByIds(ids ...int32) ([]*valueobject.Goods, error) {
	var items []*valueobject.Goods
	err := g.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_goods.item_id=gs_item.id
     WHERE gs_goods.id IN (`+format.IdArrJoinStr32(ids)+`)`)

	return items, err
}

// 获取会员价
func (g *goodsRep) GetGoodsLevelPrice(goodsId int32) []*goods.MemberPrice {
	var items []*goods.MemberPrice
	if g.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_member_price WHERE goods_id = ?`, goodsId) == nil {
		return items
	}
	return nil
}

// 保存会员价
func (g *goodsRep) SaveGoodsLevelPrice(v *goods.MemberPrice) (int32, error) {
	return orm.I32(orm.Save(g.GetOrm(), v, int(v.Id)))
}

// 移除会员价
func (g *goodsRep) RemoveGoodsLevelPrice(id int32) error {
	return g.Connector.GetOrm().DeleteByPk(goods.MemberPrice{}, id)
}

// 保存商品
func (g *goodsRep) SaveValueGoods(v *goods.ValueGoods) (int32, error) {
	return orm.I32(orm.Save(g.GetOrm(), v, int(v.Id)))
}

// 获取已上架的商品
func (g *goodsRep) GetPagedOnShelvesGoods(shopId int32, catIds []int32,
	start, end int, where, orderBy string) (int, []*valueobject.Goods) {
	var sql string
	total := 0
	catIdStr := ""
	if catIds != nil && len(catIds) > 0 {
		catIdStr = fmt.Sprintf(" AND gs_category.id IN (%s)",
			format.IdArrJoinStr32(catIds))
	}

	if len(where) != 0 {
		where = " AND " + where
	}
	if len(orderBy) != 0 {
		orderBy += ","
	}

	list := []*valueobject.Goods{}
	g.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_goods
	 INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE (?<=0 OR gs_item.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
		  %s AND gs_item.review_state=? AND gs_item.shelve_state=? %s`,
		catIdStr, where), &total, shopId, shopId, enum.ReviewPass, item.ShelvesOn)

	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE (?<=0 OR gs_item.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
		  %s AND gs_item.review_state=? AND gs_item.shelve_state=?
		  %s ORDER BY %s update_time DESC LIMIT ?,?`, catIdStr, where, orderBy)
		g.Connector.GetOrm().SelectByQuery(&list, sql, shopId, shopId,
			enum.ReviewPass, item.ShelvesOn, start, (end - start))
	}

	return total, list
}

// 获取指定数量已上架的商品
func (g *goodsRep) GetOnShelvesGoods(mchId int32, start, end int, sortBy string) []*valueobject.Goods {
	e := []*valueobject.Goods{}
	sql := fmt.Sprintf(`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE supplier_id=? AND gs_item.review_state=? AND gs_item.shelve_state=?
		 ORDER BY %s,update_time DESC LIMIT ?,?`,
		sortBy)

	g.Connector.GetOrm().SelectByQuery(&e, sql, mchId, enum.ReviewPass,
		item.ShelvesOn, start, (end - start))
	return e
}

// 保存快照
func (g *goodsRep) SaveSnapshot(v *goods.Snapshot) (int32, error) {
	var i int64
	var err error
	i, _, err = g.Connector.GetOrm().Save(v.SkuId, v)
	if i == 0 {
		_, _, err = g.Connector.GetOrm().Save(nil, v)
	}
	return v.SkuId, err
}

// 获取最新的商品快照
func (g *goodsRep) GetLatestSnapshot(skuId int32) *goods.Snapshot {
	e := &goods.Snapshot{}
	if g.Connector.GetOrm().Get(skuId, e) == nil {
		return e
	}
	return nil
}

// 根据指定商品快照
func (g *goodsRep) GetSnapshots(skuIdArr []int32) []goods.Snapshot {
	list := []goods.Snapshot{}
	g.Connector.GetOrm().SelectByQuery(&list,
		`SELECT * FROM gs_snapshot WHERE sku_id IN (`+
			format.IdArrJoinStr32(skuIdArr)+`)`)
	return list
}

// 获取最新的商品销售快照
func (g *goodsRep) GetLatestSaleSnapshot(skuId int32) *goods.SalesSnapshot {
	e := new(goods.SalesSnapshot)
	if g.Connector.GetOrm().GetBy(e, "sku_id=? ORDER BY id DESC", skuId) == nil {
		return e
	}
	return nil
}

// 获取指定的商品销售快照
func (g *goodsRep) GetSaleSnapshot(id int32) *goods.SalesSnapshot {
	e := new(goods.SalesSnapshot)
	if g.Connector.GetOrm().Get(id, e) == nil {
		return e
	}
	return nil
}

// 根据Key获取商品销售快照
func (g *goodsRep) GetSaleSnapshotByKey(key string) *goods.SalesSnapshot {
	var e *goods.SalesSnapshot = new(goods.SalesSnapshot)
	if g.Connector.GetOrm().GetBy(e, "key=?", key) == nil {
		return e
	}
	return nil
}

// 保存商品销售快照
func (g *goodsRep) SaveSaleSnapshot(v *goods.SalesSnapshot) (int32, error) {
	return orm.I32(orm.Save(g.Connector.GetOrm(), v, int(v.Id)))
}
