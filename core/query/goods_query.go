/**
 * Copyright 2015 @ z3q.net.
 * name : goods_query.go
 * author : jarryliu
 * date : 2016-04-23 12:37
 * description :
 * history :
 */
package query

import (
	"fmt"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
)

type GoodsQuery struct {
	db.Connector
}

func NewGoodsQuery(c db.Connector) *GoodsQuery {
	return &GoodsQuery{
		Connector: c,
	}
}

//根据关键词搜索上架的商品
func (g GoodsQuery) GetPagedOnShelvesGoodsByKeyword(shopId int32, start, end int,
	keyword, where, orderBy string) (int, []*valueobject.Goods) {
	var sql string
	total := 0
	keyword = "%" + keyword + "%"
	if len(where) != 0 {
		where = " AND " + where
	}
	if len(orderBy) != 0 {
		orderBy += ","
	}

	g.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_goods
         INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN cat_category ON gs_item.category_id=cat_category.id
		 WHERE gs_item.review_state=? AND gs_item.shelve_state=?
         AND (?=0 OR gs_item.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
         AND gs_item.name LIKE ? %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn, shopId, shopId, keyword)

	e := []*valueobject.Goods{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN cat_category ON gs_item.category_id=cat_category.id
		 WHERE gs_item.review_state=? AND gs_item.shelve_state=?
         AND (?=0 OR gs_item.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
         AND gs_item.name LIKE ? %s ORDER BY %s update_time DESC LIMIT ?,?`,
			where, orderBy)
		g.Connector.GetOrm().SelectByQuery(&e, sql, enum.ReviewPass,
			item.ShelvesOn, shopId, shopId, keyword, start, (end - start))
	}

	return total, e
}

//func (g GoodsQuery) GetGoodsComplex(goodsId int) *dto.GoodsComplex {
//	e := dto.GoodsComplex{}
//	sql := `SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
//		 INNER JOIN cat_category ON gs_item.category_id=cat_category.id
//		 WHERE cat_category.mch_id=? AND gs_item.review_state=? AND
//		 gs_item.shelve_state=? AND gs_item.name LIKE ? %s
//		 ORDER BY %s update_time DESC LIMIT ?,?`
//
//	g.Connector.GetOrm().GetByQuery(&e, sql,item.ShelvesOn, goodsId)
//	return &e
//}
