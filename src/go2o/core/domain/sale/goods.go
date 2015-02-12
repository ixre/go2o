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
	"fmt"
	"go2o/core/domain/interface/sale"
	"strconv"
	"time"
)

var _ sale.IGoods = new(Goods)

type Goods struct {
	value          *sale.ValueGoods
	saleRep        sale.ISaleRep
	sale           *Sale
	latestSnapshot *sale.GoodsSnapshot
}

func newGoods(sale *Sale, v *sale.ValueGoods, saleRep sale.ISaleRep) sale.IGoods {
	return &Goods{value: v,
		saleRep: saleRep,
		sale:    sale,
	}
}

func (this *Goods) GetDomainId() int {
	return this.value.Id
}

func (this *Goods) GetValue() sale.ValueGoods {
	return *this.value
}

func (this *Goods) SetValue(v *sale.ValueGoods) error {
	if v.Id == this.value.Id {
		v.CreateTime = this.value.CreateTime
		this.value = v
	}
	this.value.UpdateTime = time.Now().Unix()
	return nil
}

// 是否上架
func (this *Goods) IsOnShelves() bool {
	return this.value.OnShelves == 1
}

func (this *Goods) Save() (int, error) {
	this.sale.clearCache(this.value.Id)

	unix := time.Now().Unix()
	this.value.UpdateTime = unix

	if this.GetDomainId() <= 0 {
		this.value.CreateTime = unix
	}

	if this.value.GoodsNo == "" {
		cs := strconv.Itoa(this.value.CategoryId)
		us := strconv.Itoa(int(unix))
		l := len(cs)
		this.value.GoodsNo = fmt.Sprintf("%s%s", cs, us[4+l:])
	}

	id, err := this.saleRep.SaveGoods(this.value)
	if err == nil {
		// 创建快照
		_, err = this.GenerateSnapshot()
	}
	return id, err
}

// 生成快照
func (this *Goods) GenerateSnapshot() (int, error) {
	v := this.value
	if v.Id <= 0 {
		return 0, sale.ErrNoSuchGoods
	}

	if v.OnShelves == 0 {
		return 0, sale.ErrNotOnShelves
	}

	partnerId := this.sale.GetAggregateRootId()
	unix := time.Now().Unix()
	cate := this.saleRep.GetCategory(partnerId, v.CategoryId)
	var gsn *sale.GoodsSnapshot = &sale.GoodsSnapshot{
		Key:          fmt.Sprintf("%d-g%d-%d", partnerId, v.Id, unix),
		GoodsId:      this.GetDomainId(),
		GoodsName:    v.Name,
		GoodsNo:      v.GoodsNo,
		SmallTitle:   v.SmallTitle,
		CategoryName: cate.Name,
		Image:        v.Image,
		Cost:         v.Cost,
		Price:        v.Price,
		SalePrice:    v.SalePrice,
		CreateTime:   unix,
	}

	if this.isNewSnapshot(gsn) {
		this.latestSnapshot = gsn
		return this.saleRep.SaveSnapshot(gsn)
	}
	return 0, sale.ErrLatestSnapshot
}

// 是否为新快照,与旧有快照进行数据对比
func (this *Goods) isNewSnapshot(gsn *sale.GoodsSnapshot) bool {
	latestGsn := this.GetLatestSnapshot()
	if latestGsn != nil {
		return latestGsn.GoodsName != gsn.GoodsName ||
			latestGsn.SmallTitle != gsn.SmallTitle ||
			latestGsn.CategoryName != gsn.CategoryName ||
			latestGsn.Image != gsn.Image ||
			latestGsn.Cost != gsn.Cost ||
			latestGsn.Price != gsn.Price ||
			latestGsn.SalePrice != gsn.SalePrice
	}
	return true
}

// 获取最新的快照
func (this *Goods) GetLatestSnapshot() *sale.GoodsSnapshot {
	if this.latestSnapshot == nil {
		this.latestSnapshot = this.saleRep.GetLatestGoodsSnapshot(this.GetDomainId())
	}
	return this.latestSnapshot
}
