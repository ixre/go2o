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
	"fmt"
	"github.com/ixre/gof/db"
	"go2o/core/domain/interface/content"
)

type ContentQuery struct {
	db.Connector
}

func NewContentQuery(c db.Connector) *ContentQuery {
	return &ContentQuery{c}
}

func (cq *ContentQuery) PagedArticleList(catId int32, begin, size int, where string) (total int,
	rows []*content.Article) {
	if len(where) != 0 {
		where = " AND " + where
	}
	cq.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM
		article_list WHERE cat_id= $1 %s`, where), &total, catId)
	rows = []*content.Article{}
	if total > 0 {
		cq.Connector.GetOrm().SelectByQuery(&rows, fmt.Sprintf(`SELECT * FROM
		article_list WHERE cat_id= $1 %s ORDER BY update_time DESC LIMIT $3 OFFSET $2`, where),
			catId, begin, size)
		for i := 0; i < len(rows); i++ {
			rows[i].Content = ""
		}
	}

	return total, rows
}
