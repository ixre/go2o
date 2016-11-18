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
func (g *goodsRep) GetGoodsBySKuId(skuId int64) interface{} {
	snap := g.GetLatestSnapshot(skuId)
	if snap != nil {
		return g._saleRep.GetSale(snap.VendorId).
			GoodsManager().GetGoods(skuId)
	}
	return nil
}

// 获取商品
func (g *goodsRep) GetValueGoods(itemId int, skuId int) *goods.ValueGoods {
	var e *goods.ValueGoods = new(goods.ValueGoods)
	if g.Connector.GetOrm().GetBy(e, "item_id=? AND sku_id=?", itemId, skuId) == nil {
		return e
	}
	return nil
}

// 获取商品
func (g *goodsRep) GetValueGoodsById(goodsId int64) *goods.ValueGoods {
	var e *goods.ValueGoods = new(goods.ValueGoods)
	if g.Connector.GetOrm().Get(goodsId, e) == nil {
		return e
	}
	return nil
}

// 根据SKU获取商品
func (g *goodsRep) GetValueGoodsBySku(itemId, sku int64) *goods.ValueGoods {
	var e *goods.ValueGoods = new(goods.ValueGoods)
	if g.Connector.GetOrm().GetBy(e, "item_id=? AND sku_id=?", itemId, sku) == nil {
		return e
	}
	return nil
}

// 根据编号获取商品
func (g *goodsRep) GetGoodsByIds(ids ...int64) ([]*valueobject.Goods, error) {
	var items []*valueobject.Goods
	err := g.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_goods.item_id=gs_item.id
     WHERE gs_goods.id IN (`+format.IdArrJoinStr(ids)+`)`)

	return items, err
}

// 获取会员价
func (g *goodsRep) GetGoodsLevelPrice(goodsId int64) []*goods.MemberPrice {
	var items []*goods.MemberPrice
	if g.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_member_price WHERE goods_id = ?`, goodsId) == nil {
		return items
	}
	return nil
}

// 保存会员价
func (g *goodsRep) SaveGoodsLevelPrice(v *goods.MemberPrice) (id int, err error) {
	//if v.Id <= 0 {
	//	g.Connector.ExecScalar(`SELECT MAX(id) FROM gs_member_price where goods_id=? and level=?`,
	//		&v.Id, v.GoodsId, v.Level)
	//}
	if v.Id > 0 {
		_, _, err = g.Connector.GetOrm().Save(v.Id, v)
		id = v.Id
	} else {
		var id64 int64
		_, id64, err = g.Connector.GetOrm().Save(nil, v)
		v.Id = int(id64)
	}
	return id, err
}

// 移除会员价
func (g *goodsRep) RemoveGoodsLevelPrice(id int64) error {
	return g.Connector.GetOrm().DeleteByPk(goods.MemberPrice{}, id)
}

// 保存商品
func (g *goodsRep) SaveValueGoods(v *goods.ValueGoods) (id int, err error) {
	if v.Id > 0 {
		_, _, err = g.Connector.GetOrm().Save(v.Id, v)
	} else {
		var id64 int64
		_, id64, err = g.Connector.GetOrm().Save(nil, v)
		v.Id = int(id64)
	}
	return v.Id, err
}

// 获取已上架的商品
func (g *goodsRep) GetPagedOnShelvesGoods(shopId int, catIds []int, start, end int,
	where, orderBy string) (int, []*valueobject.Goods) {
	var sql string
	total := 0
	catIdStr := ""
	if catIds != nil && len(catIds) > 0 {
		catIdStr = fmt.Sprintf(" AND gs_category.id IN (%s)",
			format.IdArrJoinStr(catIds))
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
func (g *goodsRep) GetOnShelvesGoods(mchId int, start, end int, sortBy string) []*valueobject.Goods {
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
func (g *goodsRep) SaveSnapshot(v *goods.Snapshot) (int64, error) {
	var i int64
	var err error
	i, _, err = g.Connector.GetOrm().Save(v.SkuId, v)
	if i == 0 {
		_, _, err = g.Connector.GetOrm().Save(nil, v)
	}
	return v.SkuId, err
}

// 获取最新的商品快照
func (g *goodsRep) GetLatestSnapshot(skuId int64) *goods.Snapshot {
	e := &goods.Snapshot{}
	if g.Connector.GetOrm().Get(skuId, e) == nil {
		return e
	}
	return nil
}

// 根据指定商品快照
func (g *goodsRep) GetSnapshots(skuIdArr []int64) []goods.Snapshot {
	list := []goods.Snapshot{}
	g.Connector.GetOrm().SelectByQuery(&list,
		`SELECT * FROM gs_snapshot WHERE sku_id IN (`+
			format.IdArrJoinStr(skuIdArr)+`)`)
	return list
}

// 获取最新的商品销售快照
func (g *goodsRep) GetLatestSaleSnapshot(skuId int64) *goods.SalesSnapshot {
	e := new(goods.SalesSnapshot)
	if g.Connector.GetOrm().GetBy(e, "sku_id=? ORDER BY id DESC", skuId) == nil {
		return e
	}
	return nil
}

// 获取指定的商品销售快照
func (g *goodsRep) GetSaleSnapshot(id int64) *goods.SalesSnapshot {
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
func (g *goodsRep) SaveSaleSnapshot(v *goods.SalesSnapshot) (int64, error) {
	return orm.Save(g.Connector.GetOrm(), v, v.Id)
}
