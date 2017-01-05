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
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
)

type ItemQuery struct {
	db.Connector
}

func NewItemQuery(c db.Connector) *ItemQuery {
	return &ItemQuery{
		Connector: c,
	}
}

//根据关键词搜索上架的商品
func (i ItemQuery) GetPagedOnShelvesItem(catId int32,
	start, end int, where, orderBy string) (int, []*item.GoodsItem) {
	var sql string
	total := 0

	if len(orderBy) != 0 {
		orderBy += ","
	}

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE item_info.cat_id=? AND item_info.review_state=?
		 AND item_info.shelve_state=? %s`, where), &total,
		catId, enum.ReviewPass, item.ShelvesOn)
	list := []*item.GoodsItem{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE item_info.cat_id=? AND item_info.review_state=?
		 AND item_info.shelve_state=? %s
		 ORDER BY %s item_info.update_time DESC LIMIT ?,?`,
			where, orderBy)
		i.Connector.GetOrm().SelectByQuery(&list, sql,
			catId, enum.ReviewPass, item.ShelvesOn, start, (end - start))
	}
	return total, list
}

//根据分类获取上架的商品
func (i ItemQuery) GetOnShelvesItem(catIdArr []int32, start, end int32,
	where string) []*item.GoodsItem {
	list := []*item.GoodsItem{}
	if len(catIdArr) > 0 {
		catIdStr := format.IdArrJoinStr32(catIdArr)
		sql := fmt.Sprintf(`SELECT * FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE item_info.cat_id IN(%s) AND item_info.review_state=?
		 AND item_info.shelve_state=? %s
		 ORDER BY item_info.update_time DESC LIMIT ?,?`, catIdStr, where)
		i.Connector.GetOrm().SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, (end - start))
	}
	return list
}

// 搜索随机的商品列表
func (i ItemQuery) GetRandomItem(catId, quantity int32, where string) []*item.GoodsItem {
	/*
	       随机查询： 要减去获取的条数，以确保至少有2条数据
	   SELECT * FROM item_info

	   JOIN (SELECT ROUND(RAND() * (
	     SELECT MAX(id)-2 FROM item_info
	        )) AS id) AS r2

	        WHERE item_info.id > r2.id LIMIT 2

	*/

	list := []*item.GoodsItem{}
	sql := fmt.Sprintf(`SELECT * FROM item_info

    JOIN (SELECT ROUND(RAND() * (
      SELECT MAX(id)-? FROM item_info WHERE item_info.cat_id=?
        AND item_info.review_state=? AND item_info.shelve_state=? %s
         )) AS id) AS r2
		 WHERE item_info.Id > r2.id AND item_info.cat_id=?
		  AND item_info.review_state=?
		 AND item_info.shelve_state=? %s LIMIT ?`,
		where, where)
	i.Connector.GetOrm().SelectByQuery(&list, sql,
		quantity, catId, enum.ReviewPass, item.ShelvesOn,
		catId, enum.ReviewPass, item.ShelvesOn, quantity)
	return list
}

//根据关键词搜索上架的商品
func (i ItemQuery) GetPagedOnShelvesGoodsByKeyword(shopId int32, start, end int,
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

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=?
         AND (?=0 OR pro_product.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
         AND pro_product.name LIKE ? %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn, shopId, shopId, keyword)

	e := []*valueobject.Goods{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_info INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 INNER JOIN cat_category ON pro_product.cat_id=cat_category.id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=?
         AND (?=0 OR pro_product.supplier_id IN (SELECT mch_id FROM mch_shop WHERE id=?))
         AND pro_product.name LIKE ? %s ORDER BY %s update_time DESC LIMIT ?,?`,
			where, orderBy)
		i.Connector.GetOrm().SelectByQuery(&e, sql, enum.ReviewPass,
			item.ShelvesOn, shopId, shopId, keyword, start, (end - start))
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
//	g.Connector.GetOrm().GetByQuery(&e, sql,item.ShelvesOn, goodsId)
//	return &e
//}
