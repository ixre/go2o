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
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/infrastructure/format"
)

var _ item.IItemRepo = new(itemRepo)

type itemRepo struct {
	db.Connector
}

func NewItemRepo(c db.Connector) item.IItemRepo {
	return &itemRepo{
		Connector: c,
	}
}

func (i *itemRepo) GetValueItem(itemId int32) *item.Item {
	var e *item.Item = new(item.Item)
	//todo: supplier_id  == -1
	if i.Connector.GetOrm().GetByQuery(e, `select * FROM gs_item
			INNER JOIN cat_category c ON c.id = gs_item.category_id
			 WHERE gs_item.id=?`, itemId) == nil {
		return e
	}
	return nil
}

func (i *itemRepo) GetItemByIds(ids ...int32) ([]*item.Item, error) {
	//todo: mchId
	var items []*item.Item

	//todo:改成database/sql方式，不使用orm
	err := i.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_item WHERE id IN (`+format.IdArrJoinStr32(ids)+`)`)

	return items, err
}

func (i *itemRepo) SaveValueItem(v *item.Item) (int32, error) {
	return orm.I32(orm.Save(i.GetOrm(), v, int(v.Id)))
}

func (i *itemRepo) GetPagedOnShelvesItem(mchId int32, catIds []int32,
	start, end int) (total int, e []*item.Item) {
	var sql string

	var catIdStr string = format.IdArrJoinStr32(catIds)
	sql = fmt.Sprintf(`SELECT * FROM gs_item INNER JOIN cat_category ON gs_item.category_id=cat_category.id
		WHERE merchant_id=%d AND cat_category.id IN (%s) AND on_shelves=1 LIMIT %d,%d`, mchId, catIdStr, start, (end - start))

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_item INNER JOIN cat_category ON gs_item.category_id=cat_category.id
		WHERE merchant_id=%d AND cat_category.id IN (%s) AND on_shelves=1`, mchId, catIdStr), &total)

	e = []*item.Item{}
	i.Connector.GetOrm().SelectByQuery(&e, sql)

	return total, e
}

// 获取货品销售总数
func (i *itemRepo) GetItemSaleNum(mchId int32, id int32) int {
	var num int
	i.Connector.ExecScalar(`SELECT SUM(sale_num) FROM gs_goods WHERE item_id=?`,
		&num, id)
	return num
}

func (i *itemRepo) DeleteItem(mchId, itemId int32) error {
	_, _, err := i.Connector.Exec(`
		DELETE f FROM gs_item AS f
		INNER JOIN cat_category AS c ON f.category_id=c.id
		WHERE f.id=? AND c.merchant_id=?`, itemId, mchId)
	return err
}
