/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package dao

import (
	"database/sql"
	"github.com/jrsix/gof/db"
	"github.com/jrsix/gof/web/ui/tree"
	"go2o/src/core/ording/dao/entity"
)

type commDao struct {
	db.Connector
}

func (this *commDao) GetChinaTree() *tree.TreeNode {
	var nodes []tree.TreeNode = []tree.TreeNode{}
	this.Connector.Query(`SELECT
				Id,
				Pid,
				Name,
				Id as Value
				FROM china`,
		func(rows *sql.Rows) {
			for rows.Next() {
				node := tree.TreeNode{}
				rows.Scan(&node.Id, &node.Pid, &node.Text, &node.Value)
				nodes = append(nodes, node)
			}
			rows.Close()
		})

	return tree.List2Tree(nodes)
}

//获取父级位置
func (this *commDao) GetParentPlace(id int) *entity.Area {
	var e *entity.Area
	this.Connector.QueryRow(`SELECT Id, Pid, Name FROM china
		WHERE Id = (SELECT china2.Pid
		FROM china china2 WHERE china2.Id=?)`,
		func(row *sql.Row) {
			e = &entity.Area{}
			row.Scan(&e.Id, &e.Pid, &e.Name)
		}, id)
	return e
}
