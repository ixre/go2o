/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 11:44
 * description :
 * history :
 */

package sale

import (
	"com/domain/interface/sale"
)

type Sale struct {
	partnerId int
	saleRep   sale.ISaleRep
}

func NewSale(partnerId int, saleRep sale.ISaleRep) sale.ISale {
	return &Sale{
		partnerId: partnerId,
		saleRep:   saleRep,
	}
}

func (this *Sale) CreateProduct(val *sale.ValueProduct) sale.IProduct {
	return newProduct(val, this.saleRep)
}
