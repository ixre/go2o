package dao

import (
	"com/ording/entity"
	"database/sql"
	"ops/cf/db"
	"ops/cf/web/ui/tree"
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
func (this *commDao) GetParentPlace(id int) *entity.China {
	var e *entity.China
	this.Connector.QueryRow(`SELECT Id, Pid, Name FROM china
		WHERE Id = (SELECT china2.Pid
		FROM china china2 WHERE china2.Id=?)`,
		func(row *sql.Row) {
			e = &entity.China{}
			row.Scan(&e.Id, &e.Pid, &e.Name)
		}, id)
	return e
}
