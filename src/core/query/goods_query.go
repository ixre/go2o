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
	"go2o/src/core/domain/interface/valueobject"
)

type GoodsQuery struct {
	db.Connector
}

func NewGoodsQuery(c db.Connector) *GoodsQuery {
	return &GoodsQuery{
		Connector: c,
	}
}

func (this GoodsQuery) GetPagedOnShelvesGoodsByKeyword(partnerId int, start, end int,
	keyword, where, orderBy string) (total int, goods []*valueobject.Goods) {
	var sql string

	keyword = "%" + keyword + "%"
	if len(where) != 0 {
		where = " AND " + where
	}
	if len(orderBy) != 0 {
		orderBy += ","
	}

	this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_goods
         INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE gs_category.merchant_id=? AND gs_item.state=1
		 AND gs_item.on_shelves=1 AND gs_item.name LIKE ? %s`, where), &total, partnerId, keyword)

	e := []*valueobject.Goods{}
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM gs_goods INNER JOIN gs_item ON gs_item.id = gs_goods.item_id
		 INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		 WHERE gs_category.merchant_id=? AND gs_item.state=1
		 AND gs_item.on_shelves=1 AND gs_item.name LIKE ? %s ORDER BY %s update_time DESC LIMIT ?,?`,
			where, orderBy)

		this.Connector.GetOrm().SelectByQuery(&e, sql, partnerId, keyword, start, (end - start))
	}

	return total, e
}
