/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package dproxy

import (
	"com/domain/interface/sale"
)

type saleService struct {
	saleRep sale.ISaleRep
}

func (this *saleService) GetValueProduct(partnerId, productId int) *sale.ValueProduct {
	sl := this.saleRep.GetSale(partnerId)
	pro := sl.GetProduct(productId)
	v := pro.GetValue()
	return &v
}

func (this *saleService) SaveProduct(partnerId int, v *sale.ValueProduct) (int, error) {
	sl := this.saleRep.GetSale(partnerId)
	pro := sl.GetProduct(v.Id)
	pro.SetValue(v)
	return pro.Save()
}

func (this *saleService) GetProductsByCid(partnerId, cid, num int) []*sale.ValueProduct {
	return this.saleRep.GetProductsByCid(partnerId, cid, num)
}

func (this *saleService) DeleteProduct(partnerId, productId int) error {
	sl := this.saleRep.GetSale(partnerId)
	return sl.DeleteProduct(productId)
}
