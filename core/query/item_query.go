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
	"bytes"
	"fmt"
	"github.com/ixre/gof/db"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"strings"
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
	start, end int32, where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

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

//根据关键词搜索上架的商品
func (i ItemQuery) GetPagedOnShelvesItemForWholesale(catId int32,
	start, end int32, where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM ws_item
         INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE item_info.cat_id=? AND ws_item.review_state=?
		 AND ws_item.shelve_state=? %s`, where), &total,
		catId, enum.ReviewPass, item.ShelvesOn)
	list := []*item.GoodsItem{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM  ws_item
         INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE item_info.cat_id=? AND ws_item.review_state=?
		 AND ws_item.shelve_state=? %s
		 ORDER BY %s item_info.update_time DESC LIMIT ?,?`,
			where, orderBy)
		i.Connector.GetOrm().SelectByQuery(&list, sql,
			catId, enum.ReviewPass, item.ShelvesOn, start, (end - start))
	}
	return total, list
}

//根据关键词搜索上架的商品
func (i ItemQuery) SearchOnShelvesItem(word string, start, end int32,
	where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}

	buf := bytes.NewBuffer(nil)
	if word != "" {
		buf.WriteString(" AND (item_info.title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%' OR item_info.short_title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%')")
	}
	buf.WriteString(where)
	where = buf.String()

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE  item_info.review_state=?
		 AND item_info.shelve_state=? %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn)
	list := []*item.GoodsItem{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_info
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE item_info.review_state=?
		 AND item_info.shelve_state=? %s
		 ORDER BY %s item_info.update_time DESC LIMIT ?,?`,
			where, orderBy)
		i.Connector.GetOrm().SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, (end - start))
	}
	return total, list
}

//根据关键词搜索上架的商品
func (i ItemQuery) SearchOnShelvesItemForWholesale(word string, start, end int32,
	where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}

	buf := bytes.NewBuffer(nil)
	if word != "" {
		buf.WriteString(" AND (item_info.title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%' OR item_info.short_title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%')")
	}
	buf.WriteString(where)
	where = buf.String()

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM ws_item
         INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE ws_item.review_state=?
		 AND ws_item.shelve_state=?  %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn)
	list := []*item.GoodsItem{}

	if total > 0 {
		sql = fmt.Sprintf(`SELECT item_info.id,item_info.product_id,item_info.prom_flag,
		item_info.cat_id,item_info.vendor_id,item_info.brand_id,item_info.shop_id,
		item_info.shop_cat_id,item_info.express_tid,item_info.title,
		item_info.short_title,item_info.code,item_info.image,
		item_info.is_present,ws_item.price_range,item_info.stock_num,
		item_info.sale_num,item_info.sku_num,item_info.sku_id,item_info.cost,
		ws_item.price,item_info.retail_price,item_info.weight,item_info.bulk,
		item_info.shelve_state,item_info.review_state,item_info.review_remark,
		item_info.sort_num,item_info.create_time,item_info.update_time
		 FROM ws_item INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 WHERE ws_item.review_state=?
		 AND ws_item.shelve_state=? %s
		 ORDER BY %s item_info.update_time DESC LIMIT ?,?`,
			where, orderBy)
		i.Connector.GetOrm().SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, (end - start))
	}
	return total, list
}

//根据分类获取上架的商品
func (i ItemQuery) GetOnShelvesItem(catIdArr []int32, start, end int32,
	where string) []*item.GoodsItem {
	list := []*item.GoodsItem{}
	if len(catIdArr) > 0 {
		catIdStr := format.I32ArrStrJoin(catIdArr)
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

	s := []string{where}
	if catId > 0 {
		if where != "" {
			s = append(s, " AND")
		}
		s = append(s, fmt.Sprintf("item_info.cat_id=%d)", catId))
	}
	search := strings.Join(s, "")

	list := []*item.GoodsItem{}
	sql := fmt.Sprintf(`SELECT * FROM item_info
    JOIN (SELECT ROUND(RAND() * (
      SELECT MAX(id)-? FROM item_info WHERE  item_info.review_state=?
         AND item_info.shelve_state=? %s
         )) AS id) AS r2
		 WHERE item_info.ID > r2.id
		  AND item_info.review_state=?
		 AND item_info.shelve_state=? %s LIMIT ?`,
		search, search)
	i.Connector.GetOrm().SelectByQuery(&list, sql,
		quantity, enum.ReviewPass, item.ShelvesOn,
		enum.ReviewPass, item.ShelvesOn, quantity)
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
		 INNER JOIN pro_category ON pro_product.cat_id=pro_category.id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=?
         AND (?=0 OR pro_product.supplier_id IN (SELECT vendor_id FROM mch_shop WHERE id=?))
         AND pro_product.name LIKE ? %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn, shopId, shopId, keyword)

	e := []*valueobject.Goods{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_info INNER JOIN pro_product ON pro_product.id = item_info.product_id
		 INNER JOIN pro_category ON pro_product.cat_id=pro_category.id
		 WHERE pro_product.review_state=? AND pro_product.shelve_state=?
         AND (?=0 OR pro_product.supplier_id IN (SELECT vendor_id FROM mch_shop WHERE id=?))
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
//		 INNER JOIN pro_category ON pro_product.cat_id=pro_category.id
//		 WHERE pro_category.mch_id=? AND pro_product.review_state=? AND
//		 pro_product.shelve_state=? AND pro_product.name LIKE ? %s
//		 ORDER BY %s update_time DESC LIMIT ?,?`
//
//	g.Connector.GetOrm().GetByQuery(&e, sql,item.ShelvesOn, goodsId)
//	return &e
//}
