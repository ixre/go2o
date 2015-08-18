/**
 * Copyright 2014 @ S1N1 Team.
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
	"github.com/jrsix/gof/db"
	"go2o/src/core/domain/interface/member"
)

type MemberQuery struct {
	db.Connector
}

func NewMemberQuery(c db.Connector) *MemberQuery {
	return &MemberQuery{c}
}

// 获取返现记录
func (this *MemberQuery) QueryIncomeLog(memberId, page, size int,
	where, orderBy string) (num int, rows []map[string]interface{}) {

	d := this.Connector

	if where != "" {
		where = "WHERE " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0)
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s`, where), &num, memberId)

	sqlLine := fmt.Sprintf(`SELECT l.*,
			record_time,
			convert(l.fee,CHAR(10)) as fee
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s %s LIMIT ?,?`,
		where, orderBy)

	d.Query(sqlLine, func(_rows *sql.Rows) {
		rows = db.RowsToMarshalMap(_rows)
		_rows.Close()
	}, memberId, (page-1)*size, size)

	return num, rows
}

// 查询分页订单
func (this *MemberQuery) QueryPagerOrder(memberId, page, size int,
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

	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM pt_order WHERE
	 		member_id=? %s`, where), &num, memberId)

	d.Query(fmt.Sprintf(` SELECT id,
			order_no,
			member_id,
			partner_id,
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
		where, orderBy),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}

// 获取最近的余额变动信息
func (this *MemberQuery) GetLatestBalanceInfoByKind(memberId int, kind int) *member.BalanceInfoValue {
	var info = new(member.BalanceInfoValue)
	if err := this.GetOrm().GetBy(info, "member_id=? AND kind=? ORDER BY create_time DESC",
		memberId, kind); err == nil {
		return info
	}
	return nil
}
