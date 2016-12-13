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
	"go2o/core/domain/interface/product"
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

	g.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=?
         AND (?=0 OR pro_product.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
         AND pro_product.name LIKE ? %s`, where), &total,
		enum.ReviewPass, product.ShelvesOn, shopId, shopId, keyword)

	e := []*valueobject.Goods{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_info INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=?
         AND (?=0 OR pro_product.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
         AND pro_product.name LIKE ? %s ORDER BY %s update_time DESC LIMIT ?,?`,
			where, orderBy)
		g.Connector.GetOrm().SelectByQuery(&e, sql, enum.ReviewPass,
			product.ShelvesOn, shopId, shopId, keyword, start, (end - start))
	}

	return total, e
}

//func (g GoodsQuery) GetGoodsComplex(goodsId int) *dto.GoodsComplex {
//	e := dto.GoodsComplex{}
//	sql := `SELECT * FROM item_info INNER JOIN pro_product ON pro_product.id = item_info.product_id
//		 INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
//		 WHERE cat_category.mch_id=? AND pro_product.review_state=? AND
//		 pro_product.shelve_state=? AND pro_product.name LIKE ? %s
//		 ORDER BY %s update_time DESC LIMIT ?,?`
//
//	g.Connector.GetOrm().GetByQuery(&e, sql,product.ShelvesOn, goodsId)
//	return &e
//}
