/**
 * Copyright 2015 @ z3q.net.
 * name : order_query
 * author : jarryliu
 * date : 2016-07-08 15:32
 * description :
 * history :
 */
package query

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
)

type OrderQuery struct {
	db.Connector
}

func NewOrderQuery(conn db.Connector) *OrderQuery {
	return &OrderQuery{Connector: conn}
}

// 查询分页订单
func (this *OrderQuery) QueryPagerOrder(memberId, page, size int, pagination bool,
	where, orderBy string) (num int, rows []map[string]interface{}) {
	d := this.Connector
	if where != "" {
		where = "AND " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	} else {
		orderBy = " ORDER BY update_time DESC,create_time desc "
	}

	if pagination {
		d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM sale_order WHERE buyer_id=? %s`,
			where), &num, memberId)
		if num == 0 {
			return 0, nil
		}
	}

	d.Query(fmt.Sprintf(`SELECT id,order_no,discount_fee,final_fee,express_fee,is_paid
FROM sale_order where buyer_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}
