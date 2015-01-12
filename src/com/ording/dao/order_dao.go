package dao

import (
	"database/sql"
	"fmt"
	"github.com/atnet/gof/db"
)

type orderDao struct {
	db.Connector
}

func (this *orderDao) GetMemberPagerOrder(memberId, page, size int,
	where, orderby string) (num int, rows []map[string]interface{}) {

	d := this.Connector

	if where != "" {
		where = "AND " + where
	}
	if orderby != "" {
		orderby = "ORDER BY " + orderby
	} else {
		orderby = " ORDER BY update_time DESC,create_time desc "
	}

	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM pt_order WHERE
	 		shop_id=? %s`, where), &num, memberId)

	dtStr := `%Y-%m-%d %T`
	d.Query(fmt.Sprintf(` SELECT id,
			order_no,
			member_id,
			pt_id,
			shop_id,
			items,
			replace(items_info,'\n','<br />') as items_info,
			total_fee,
			fee,
			pay_fee,
			pay_method,
			is_payed,
			note,
			status,
            date_format(create_time,'%s') as create_time,
            date_format(deliver_time,'%s') as deliver_time,
            date_format(update_time,'%s') as update_time
            FROM pt_order WHERE member_id=? %s %s LIMIT ?,?`,
		dtStr, dtStr, dtStr, where, orderby),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}
