/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-08 10:53
 * description :
 * history :
 */

package sale

import (
	"com/domain/interface/sale"
)

type Product struct {
	val     *sale.ValueProduct
	saleRep sale.ISaleRep
}

func newProduct(val *sale.ValueProduct, saleRep sale.ISaleRep) sale.IProduct {
	return &Product{val: val, saleRep: saleRep}
}

func (this *Product) GetDomainId() int {
	return this.val.Id
}

func (this *Product) GetValue() sale.ValueProduct {
	return *this.val
}
