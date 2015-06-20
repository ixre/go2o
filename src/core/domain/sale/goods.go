/**
 * Copyright 2015 @ S1N1 Team.
 * name : sale_goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import (
	"fmt"
	"go2o/src/core/domain"
	"go2o/src/core/domain/interface/sale"
	"time"
	"go2o/src/core/domain/interface/valueobject"
)

var _ sale.IGoods = new(SaleGoods)
var _ domain.IDomain = new(SaleGoods)

type SaleGoods struct {
	_goods   sale.IItem
	_value   *sale.ValueGoods
	_saleRep sale.ISaleRep
	_sale    sale.ISale

	_latestSnapshot *sale.GoodsSnapshot
}

func NewSaleGoods(s sale.ISale, goods sale.IItem, value *sale.ValueGoods, rep sale.ISaleRep) sale.IGoods {
	v := &SaleGoods{
		_goods:          goods,
		_value:          value,
		_saleRep:        rep,
		_sale:           s,
		_latestSnapshot: nil,
	}
	return v.init()
}

func (this *SaleGoods) init() sale.IGoods {
	this._value.SalePrice = this._goods.GetValue().SalePrice
	this._value.Price = this._value.Price
	return this
}

//获取领域对象编号
func (this *SaleGoods) GetDomainId() int {
	return this._value.Id
}

// 获取货品
func (this *SaleGoods) GetItem() sale.IItem {
	return this._goods
}

// 设置值
func (this *SaleGoods) GetValue() *sale.ValueGoods {
	return this._value
}


// 获取包装过的商品信息
func (this *SaleGoods) GetPackedValue()*valueobject.Goods{
	item := this.GetItem().GetValue()
	gv := this.GetValue()
	goods := &valueobject.Goods{
		Item_Id :item.Id,
		CategoryId:item.CategoryId,
		Name :item.Name,
		GoodsNo:item.GoodsNo,
		Image:item.Image,
		Price:item.Price,
		SalePrice:item.SalePrice,
		GoodsId :this.GetDomainId(),
		SkuId:gv.SkuId,
		IsPresent:gv.IsPresent,
		PromotionFlag:gv.PromotionFlag,
		StockNum:gv.StockNum,
		SaleNum:gv.SaleNum,
	}
	return goods
}

// 设置值
func (this *SaleGoods) SetValue(v *sale.ValueGoods) error {
	this._value.IsPresent = v.IsPresent
	this._value.SaleNum = v.SaleNum
	this._value.StockNum = v.StockNum
	this._value.SkuId = v.SkuId
	//this._value.PromotionFlag = v.PromotionFlag
	return nil
}

// 保存
func (this *SaleGoods) Save() (int, error) {
	id, err := this._saleRep.SaveValueGoods(this._value)
	if err == nil{
		_,err = this.GenerateSnapshot()
	}
	this._value.Id = id
	return id, err

	//todo: save promotion
	// return id,err
}

// 生成快照
func (this *SaleGoods) GenerateSnapshot() (int, error) {
	v := this._value
	gi := this.GetItem()
	gv := gi.GetValue()

	if v.Id <= 0 {
		return -1, sale.ErrNoSuchGoods
	}

	if !gi.IsOnShelves() {
		return -1, sale.ErrNotOnShelves
	}

	partnerId := this._sale.GetAggregateRootId()
	unix := time.Now().Unix()
	cate := this._saleRep.GetCategory(partnerId, gv.CategoryId)
	var gsn *sale.GoodsSnapshot = &sale.GoodsSnapshot{
		Key:          fmt.Sprintf("%d-g%d-%d", partnerId, v.Id, unix),
		ItemId : gv.Id,
		GoodsId:      this.GetDomainId(),
		GoodsName:    gv.Name,
		GoodsNo:      gv.GoodsNo,
		SmallTitle:   gv.SmallTitle,
		CategoryName: cate.Name,
		Image:        gv.Image,
		Cost:         gv.Cost,
		SalePrice:    gv.SalePrice,
		Price:        this._value.Price,
		CreateTime:   unix,
	}

	if this.isNewSnapshot(gsn) {
		this._latestSnapshot = gsn
		return this._saleRep.SaveSnapshot(gsn)
	}
	return 0, sale.ErrLatestSnapshot
}

// 是否为新快照,与旧有快照进行数据对比
func (this *SaleGoods) isNewSnapshot(gsn *sale.GoodsSnapshot) bool {
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
func (this *SaleGoods) GetLatestSnapshot() *sale.GoodsSnapshot {
	if this._latestSnapshot == nil {
		this._latestSnapshot = this._saleRep.GetLatestGoodsSnapshot(this.GetDomainId())
	}
	return this._latestSnapshot
}
