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

func (this *itemRep) GetValueItem(itemId int) *item.Item {
	var e *item.Item = new(item.Item)
	//todo: supplier_id  == -1
	if this.Connector.GetOrm().GetByQuery(e, `select * FROM gs_item
			INNER JOIN gs_category c ON c.id = gs_item.category_id
			 WHERE gs_item.id=?`, itemId) == nil {
		return e
	}
	return nil
}

func (this *itemRep) GetItemByIds(ids ...int) ([]*item.Item, error) {
	//todo: merchantId
	var items []*item.Item

	//todo:改成database/sql方式，不使用orm
	err := this.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM gs_item WHERE id IN (`+format.GetCategoryIdStr(ids)+`)`)

	return items, err
}

func (this *itemRep) SaveValueItem(v *item.Item) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, id, err := orm.Save(nil, v)
		return int(id), err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

func (this *itemRep) GetPagedOnShelvesItem(merchantId int, catIds []int,
	start, end int) (total int, e []*item.Item) {
	var sql string

	var catIdStr string = format.GetCategoryIdStr(catIds)
	sql = fmt.Sprintf(`SELECT * FROM gs_item INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		WHERE merchant_id=%d AND gs_category.id IN (%s) AND on_shelves=1 LIMIT %d,%d`, merchantId, catIdStr, start, (end - start))

	this.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM gs_item INNER JOIN gs_category ON gs_item.category_id=gs_category.id
		WHERE merchant_id=%d AND gs_category.id IN (%s) AND on_shelves=1`, merchantId, catIdStr), &total)

	e = []*item.Item{}
	this.Connector.GetOrm().SelectByQuery(&e, sql)

	return total, e
}

// 获取货品销售总数
func (this *itemRep) GetItemSaleNum(merchantId int, id int) int {
	var num int
	this.Connector.ExecScalar(`SELECT SUM(sale_num) FROM gs_goods WHERE item_id=?`,
		&num, id)
	return num
}

func (this *itemRep) DeleteItem(merchantId, itemId int) error {
	_, _, err := this.Connector.Exec(`
		DELETE f FROM gs_item AS f
		INNER JOIN gs_category AS c ON f.category_id=c.id
		WHERE f.id=? AND c.merchant_id=?`, itemId, merchantId)
	return err
}
