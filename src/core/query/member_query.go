/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"database/sql"
	"fmt"
	"github.com/atnet/gof/db"
)

type MemberQuery struct {
	db.Connector
}

func NewMemberQuery(c db.Connector) *MemberQuery {
	return &MemberQuery{c}
}

// 获取返现记录
func (this *MemberQuery) QueryIncomeLog(memberId, page, size int,
	where, orderby string) (num int, rows []map[string]interface{}) {

	d := this.Connector

	if where != "" {
		where = "WHERE " + where
	}
	if orderby != "" {
		orderby = "ORDER BY " + orderby
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0)
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s`, where), &num, memberId)

	d.Query(fmt.Sprintf(`SELECT l.*,
			record_time,
			convert(l.fee,CHAR(10)) as fee
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s %s LIMIT ?,?`,
		where, orderby),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}

// 查询分页订单
func (this *MemberQuery) QueryPagerOrder(memberId, page, size int,
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

	d.Query(fmt.Sprintf(` SELECT id,
			order_no,
			member_id,
			pt_id,
			shop_id,
			replace(items_info,'\n','<br />') as items_info,
			total_fee,
			fee,
			pay_fee,
			payment_opt,
			is_paid,
			note,
			status,
			paid_time,
            create_time,
            deliver_time,
            update_time
            FROM pt_order WHERE member_id=? %s %s LIMIT ?,?`,
		where, orderby),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}
