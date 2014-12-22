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
	"ops/cf/db"
	"strconv"
	"strings"
)

type SaleRep struct {
	db.Connector
}

func (this *SaleRep) GetSale(partnerId int) sale.ISale {
	return sl.NewSale(partnerId, this)
}

func (this *SaleRep) GetProductByIds(partnerId int, ids ...int) ([]sale.IProduct, error) {
	//todo: partnerId
	var items []sale.ValueProduct
	var products []sale.IProduct
	var strIds []string = make([]string, len(ids))
	for i, v := range ids {
		strIds[i] = strconv.Itoa(v)
	}

	//todo:改成database/sql方式，不使用orm
	err := this.Connector.GetOrm().SelectByQuery(&items, sale.ValueProduct{},
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
