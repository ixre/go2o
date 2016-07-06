/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
    "github.com/jsix/gof/db"
    "go2o/core/domain/interface/content"
    "fmt"
)

type ContentQuery struct {
    db.Connector
}

func NewContentQuery(c db.Connector) *ContentQuery {
    return &ContentQuery{c}
}

func (this *ContentQuery) PagedArticleList(catId, begin, size int, where string) (total int,
rows []*content.Article) {
    if len(where) != 0 {
        where = " AND " + where
    }

    this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM
		con_article WHERE cat_id=? %s`, where), &total, catId)

    rows = []*content.Article{}
    if total > 0 {
        this.Connector.GetOrm().SelectByQuery(&rows, fmt.Sprintf(`SELECT * FROM
		con_article WHERE cat_id=? %s ORDER BY update_time DESC LIMIT ?,?`, where),
            catId, begin,size)
    }

    return total, rows
}

// 获取页面列表
//func (this *MemberQuery) QueryPageList(memberId, page, size int,
//	where, orderBy string) (num int, rows []map[string]interface{}) {
//
//	d := this.Connector
//
//	if where != "" {
//		where = "WHERE " + where
//	}
//	if orderBy != "" {
//		orderBy = "ORDER BY " + orderBy
//	}
//	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0)
//			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
//			WHERE member_id=? %s`, where), &num, memberId)
//
//	sqlLine := fmt.Sprintf(`SELECT l.*,
//			record_time,
//			convert(l.fee,CHAR(10)) as fee
//			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
//			WHERE member_id=? %s %s LIMIT ?,?`,
//		where, orderBy)
//
//	d.Query(sqlLine, func(_rows *sql.Rows) {
//		rows = db.RowsToMarshalMap(_rows)
//		_rows.Close()
//	}, memberId, (page-1)*size, size)
//
//	return num, rows
//}
