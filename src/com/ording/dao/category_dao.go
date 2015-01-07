package dao

import (
	"com/ording/entity"
	"database/sql"
	"github.com/newmin/gof/db"
	"github.com/newmin/gof/web/ui/tree"
	"strconv"
	"time"
)

type categoryDao struct {
	db.Connector
}

func (this *categoryDao) GetCategoriesOfPartner(partnerId int) (categories []entity.Category) {
	categories = []entity.Category{}
	this.Connector.Query(`SELECT
				id,
				pid,
				ptid,
				name,
				order_index,
				create_time,
				enabled,
				descript
				FROM it_category
				WHERE ptid =? ORDER BY order_index,name`,
		func(rows *sql.Rows) {
			for rows.Next() {
				category := entity.Category{}
				var createTime string
				rows.Scan(&category.Id, &category.Pid, &category.Ptid, &category.Name,
					&category.Idx, &createTime, &category.IsEnabled, &category.Descript)
				category.CreateTime, _ = time.Parse("2006-01-02 15:04:05", createTime)

				categories = append(categories, category)
			}
			rows.Close()
		}, partnerId)
	return categories
}

func (this *categoryDao) GetCategoryById(partnerId int, id int) (category *entity.Category) {
	this.Connector.QueryRow(`SELECT
				Id,
				Pid,
				Ptid,
				Name,
				Idx,
				CreateTime,
				Descript,
				IsEnabled
				FROM it_category
				WHERE Ptid =? AND Id=?`,
		func(row *sql.Row) {
			category = &entity.Category{}
			var createTime string

			//IsEnabled bool放在最后，否则不能取到Descript的值
			row.Scan(&category.Id, &category.Pid, &category.Ptid, &category.Name,
				&category.Idx, &createTime, &category.Descript, &category.IsEnabled)
			category.CreateTime, _ = time.Parse("2006-01-02 15:04:05", createTime)

		}, partnerId, id)
	return category
}

func (this *categoryDao) SaveCategory(category *entity.Category) (id int, err error) {
	if category.Id <= 0 {
		//多行字符用``
		_, id, err := this.Connector.Exec(`INSERT INTO it_category
			(
			Pid,
			Ptid,
			Name,
			Idx,
			CreateTime,
			IsEnabled,
			Descript)
			VALUES
			(
			?,
			?,
			?,
			?,
			?,
			?,
			?)
`, category.Pid, category.Ptid, category.Name,
			category.Idx, category.CreateTime,
			category.IsEnabled, category.Descript)
		return id, err
	} else {
		_, _, err := this.Connector.Exec(`UPDATE it_category
					SET
					Pid = ?,
					Name =?,
					Idx =?,
					IsEnabled = ?,
					Descript =?
					WHERE Ptid=? AND Id =?
`, category.Pid, category.Name,
			category.Idx, category.IsEnabled,
			category.Descript, category.Ptid, category.Id)
		return category.Id, err
	}
}

func (this *categoryDao) DeleteCategoryAndRelation(partnerId int, categoryId int) error {

	//删除子类
	_, _, err := this.Connector.Exec("DELETE FROM it_category WHERE Ptid=? AND Pid=?",
		partnerId, categoryId)

	//删除分类
	this.Connector.Exec("DELETE FROM it_category WHERE Ptid=? AND Id=?",
		partnerId, categoryId)

	//清理项
	this.Connector.Exec(`DELETE FROM it_item WHERE Cid NOT IN
		(SELECT Id FROM it_category WHERE Ptid=?`, categoryId)

	return err
}

//删除分类
func DelCategory(partnerId int, categoryId int) error {
	return Category().DeleteCategoryAndRelation(partnerId, categoryId)
}

func GetCategoryTreeNode(partnerId int) *tree.TreeNode {
	var categories []entity.Category = Category().GetCategoriesOfPartner(partnerId)
	rootNode := &tree.TreeNode{
		Text:   "根节点",
		Value:  "",
		Url:    "",
		Icon:   "",
		Open:   true,
		Childs: nil}
	iterCategoryTree(rootNode, 0, categories)
	return rootNode
}

func iterCategoryTree(node *tree.TreeNode, pid int, categories []entity.Category) {
	node.Childs = []*tree.TreeNode{}
	for _, cate := range categories {
		if cate.Pid == pid {
			cNode := &tree.TreeNode{
				Text:   cate.Name,
				Value:  strconv.Itoa(cate.Id),
				Url:    "",
				Icon:   "",
				Open:   true,
				Childs: nil}
			node.Childs = append(node.Childs, cNode)
			iterCategoryTree(cNode, cate.Id, categories)
		}
	}
}
