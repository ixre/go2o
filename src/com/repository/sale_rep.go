/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 11:09
 * description :
 * history :
 */

package repository

import (
	"com/domain/interface/sale"
	sl "com/domain/sale"
	"com/infrastructure/log"
	"database/sql"
	"fmt"
	"github.com/newmin/gof/db"
	"strconv"
	"strings"
)

var _ sale.ISaleRep = new(saleRep)

type saleRep struct {
	db.Connector
	cache map[int]sale.ISale
}

func NewSaleRep(c db.Connector) sale.ISaleRep {
	return (&saleRep{
		Connector: c,
	}).init()
}

func (this *saleRep) init() sale.ISaleRep {
	this.cache = make(map[int]sale.ISale)
	return this
}

func (this *saleRep) GetSale(partnerId int) sale.ISale {
	v, ok := this.cache[partnerId]
	if !ok {
		v = sl.NewSale(partnerId, this)
		this.cache[partnerId] = v
	}
	return v
}

func (this *saleRep) GetValueProduct(partnerId, productId int) *sale.ValueProduct {
	var e *sale.ValueProduct = new(sale.ValueProduct)
	err := this.Connector.GetOrm().GetByQuery(e, `select * FROM it_item
			INNER JOIN it_category c ON c.id = it_item.cid WHERE it_item.id=?
			AND c.ptid=?`, productId, partnerId)
	if err != nil {
		return nil
	}
	return e
}

func (this *saleRep) GetProductByIds(partnerId int, ids ...int) ([]sale.IProduct, error) {
	//todo: partnerId
	var items []sale.ValueProduct
	var products []sale.IProduct
	var strIds []string = make([]string, len(ids))
	for i, v := range ids {
		strIds[i] = strconv.Itoa(v)
	}

	//todo:改成database/sql方式，不使用orm
	err := this.Connector.GetOrm().SelectByQuery(&items,
		`SELECT * FROM it_item WHERE id IN (`+strings.Join(strIds, ",")+`)`)

	if err != nil {
		return nil, err
	}

	s := this.GetSale(partnerId)

	products = make([]sale.IProduct, len(items))
	for i, v := range items {
		products[i] = s.CreateProduct(&v)
	}
	return products, err
}

func (this *saleRep) SaveProduct(v *sale.ValueProduct) (int, error) {
	orm := this.Connector.GetOrm()
	if v.Id <= 0 {
		_, id, err := orm.Save(nil, v)
		return int(id), err
	} else {
		_, _, err := orm.Save(v.Id, v)
		return v.Id, err
	}
}

func (this *saleRep) GetProductsByCid(partnerId, categoryId, num int) (e []*sale.ValueProduct) {
	var sql string
	if num <= 0 {
		sql = fmt.Sprintf(`SELECT * FROM it_item INNER JOIN it_category ON it_item.cid=it_category.id
		WHERE ptid=%d AND it_category.id=%d`, partnerId, categoryId)
	} else {
		sql = fmt.Sprintf(`SELECT * FROM it_item INNER JOIN it_category ON it_item.cid=it_category.id
		WHERE ptid=%d AND it_category.id=%d LIMIT 0,%d`, partnerId, categoryId, num)
	}

	e = []*sale.ValueProduct{}
	err := this.Connector.GetOrm().SelectByQuery(&e, sql)
	if err != nil {
		log.PrintErr(err)
		return nil
	}
	return e
}

func (this *saleRep) DeleteProduct(partnerId, productId int) error {
	_, _, err := this.Connector.Exec(`
		DELETE f,f2 FROM it_item AS f
		INNER JOIN it_category AS c ON f.cid=c.id
		INNER JOIN it_itemprop as f2 ON f2.id=f.id
		WHERE f.id=? AND c.ptid=?`, productId, partnerId)
	return err
}

//获取食物数量
//todo: 还未使用
func (this *saleRep) FoodItemsCount(partnerId, cid int) (count int) {
	this.Connector.QueryRow(`
		SELECT COUNT(0) FROM it_item f
	INNER JOIN it_category c ON f.cid = c.id
	 where c.ptid = ?
	AND (cid == -1 OR cid = ?)
	`, func(r *sql.Row) {
		r.Scan(count)
	}, partnerId, partnerId)
	return count
}
