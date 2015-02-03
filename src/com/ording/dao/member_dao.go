package dao

import (
	"com/ording"
	"database/sql"
	"fmt"
	"github.com/atnet/gof/db"
)

type memberDao struct {
	db.Connector
}

func (this *memberDao) GetIncomeLog(memberId, page, size int,
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
			date_format(l.record_time,%s) as record_time,
			convert(l.fee,CHAR(10)) as fee
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s %s LIMIT ?,?`, "'%Y-%m-%d %T'",
		where, orderby),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}

//@deparend
//验证用户密码
func (this *memberDao) Verify(usr, pwd string) bool {
	var id int
	encPwd := ording.EncodeMemberPwd(usr, pwd)
	if err := this.Connector.ExecScalar("SELECT id FROM mm_member WHERE usr=? AND pwd=?", &id, usr, encPwd); err != nil {
		return false
	}
	return id != 0
}
