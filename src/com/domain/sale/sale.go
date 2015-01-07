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

var _ sale.ISale = new(Sale)

const MAX_CACHE_SIZE int = 1000

type Sale struct {
	partnerId int
	saleRep   sale.ISaleRep
	proCache  map[int]sale.IProduct
}

func NewSale(partnerId int, saleRep sale.ISaleRep) sale.ISale {
	return (&Sale{
		partnerId: partnerId,
		saleRep:   saleRep,
	}).init()
}

func (this *Sale) init() sale.ISale {
	this.proCache = make(map[int]sale.IProduct)
	return this
}

func (this *Sale) clearCache(productId int){
	delete(this.proCache, productId)
}

func (this *Sale) chkCache() {
	if len(this.proCache) >= MAX_CACHE_SIZE {
		this.proCache = make(map[int]sale.IProduct)
	}
}


func (this *Sale) GetAggregateRootId() int {
	return this.partnerId
}

func (this *Sale) CreateProduct(val *sale.ValueProduct) sale.IProduct {
	return newProduct(this,val, this.saleRep)
}

// 根据产品编号获取产品
func (this *Sale) GetProduct(productId int) sale.IProduct {
	p, ok := this.proCache[productId]
	if !ok {
		this.chkCache()
		pv := this.saleRep.GetValueProduct(this.GetAggregateRootId(), productId)

		if pv != nil {
			p = this.CreateProduct(pv)
			this.proCache[productId] = p
		}
	}
	return p
}

// 删除商品
func (this *Sale) DeleteProduct(productId int) error {
	err := this.saleRep.DeleteProduct(this.GetAggregateRootId(), productId)
	if err != nil {
		this.clearCache(productId)
	}
	return err
}
