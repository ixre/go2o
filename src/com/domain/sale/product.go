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

var _ sale.IProduct = new(Product)

type Product struct {
	value   *sale.ValueProduct
	saleRep sale.ISaleRep
	sale    *Sale
}

func newProduct(sale *Sale, v *sale.ValueProduct, saleRep sale.ISaleRep) sale.IProduct {
	return &Product{value: v,
		saleRep: saleRep,
		sale:    sale,
	}
}

func (this *Product) GetDomainId() int {
	return this.value.Id
}

func (this *Product) GetValue() sale.ValueProduct {
	return *this.value
}

func (this *Product) SetValue(v *sale.ValueProduct) error {
	if v.Id == this.value.Id {
		this.value = v
	}
	return nil
}

func (this *Product) Save() (int, error) {
	this.sale.clearCache(this.value.Id)
	return this.saleRep.SaveProduct(this.value)
}
