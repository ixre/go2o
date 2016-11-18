/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 11:09
 * description :
 * history :
 */

package repository

import (
	"fmt"
	"github.com/jsix/gof/db"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/infrastructure/format"
)

var _ item.IItemRep = new(itemRep)

type itemRep struct {
	db.Connector
}

func NewItemRep(c db.Connector) item.IItemRep {
	return &itemRep{
		Connector: c,
	}
}

func (i *itemRep) GetValueItem(itemId int64) *item.Item {
	var e *item.Item = new(item.Item)
	//todo: supplier_id  == -1
	if i.Connector.GetOrm().GetByQuery(e, `select * FROM gs_item
			INNER JOIN gs_category c ON c.id = gs_item.category_id
			 WHERE gs_item.id=?`, itemId) == nil {
		return e
	}
	return nil
}

func (i *itemRep) GetItemByIds(ids ...int64) ([]*item.Item, error) {
	//todo: mchId
	var items []*item.Item

	//todo:改成database/sql方式，不使用orm
	err := i.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_item WHERE id IN (`+format.IdArrJoinStr(ids)+`)`)

	return items, err
}

func (i *itemRep) SaveValueItem(v *item.Item) (int64, error) {
	orm := i.Connector.GetOrm()
	if v.Id <= 0 {
		_, id, err := orm.Save(nil, v)
		return int(id), err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

func (i *itemRep) GetPagedOnShelvesItem(mchId int64, catIds []int64,
	start, end int) (total int, e []*item.Item) {
	var sql string

	var catIdStr string = format.IdArrJoinStr(catIds)
	sql = fmt.Sprintf(`SELECT * FROM gs_item INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		WHERE merchant_id=%d AND gs_category.id IN (%s) AND on_shelves=1 LIMIT %d,%d`, mchId, catIdStr, start, (end - start))

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_item INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		WHERE merchant_id=%d AND gs_category.id IN (%s) AND on_shelves=1`, mchId, catIdStr), &total)

	e = []*item.Item{}
	i.Connector.GetOrm().SelectByQuery(&e, sql)

	return total, e
}

// 获取货品销售总数
func (i *itemRep) GetItemSaleNum(mchId int64, id int64) int {
	var num int
	i.Connector.ExecScalar(`SELECT SUM(sale_num) FROM gs_goods WHERE item_id=?`,
		&num, id)
	return num
}

func (i *itemRep) DeleteItem(mchId, itemId int64) error {
	_, _, err := i.Connector.Exec(`
		DELETE f FROM gs_item AS f
		INNER JOIN gs_category AS c ON f.category_id=c.id
		WHERE f.id=? AND c.merchant_id=?`, itemId, mchId)
	return err
}
