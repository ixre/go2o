/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:53
 * description :
 * history :
 */

package sale

import (
	"fmt"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"strconv"
	"time"
)

var _ sale.IItem = new(Item)

type Item struct {
	_value      *sale.ValueItem
	_saleRep    sale.ISaleRep
	_saleTagRep sale.ISaleTagRep
	_goodsRep   sale.IGoodsRep
	_promRep    promotion.IPromotionRep
	_sale       *Sale
	_saleTags   []*sale.SaleLabel
}

func newItem(sale *Sale, v *sale.ValueItem, saleRep sale.ISaleRep,
	saleTagRep sale.ISaleTagRep, goodsRep sale.IGoodsRep, promRep promotion.IPromotionRep) sale.IItem {
	return &Item{
		_value:      v,
		_saleRep:    saleRep,
		_saleTagRep: saleTagRep,
		_sale:       sale,
		_goodsRep:   goodsRep,
	}
}

func (this *Item) GetDomainId() int {
	return this._value.Id
}

func (this *Item) GetValue() sale.ValueItem {
	return *this._value
}

func (this *Item) SetValue(v *sale.ValueItem) error {
	if v.Id == this._value.Id {
		v.CreateTime = this._value.CreateTime
		v.GoodsNo = this._value.GoodsNo
		this._value = v
	}
	this._value.UpdateTime = time.Now().Unix()
	return nil
}

// 是否上架
func (this *Item) IsOnShelves() bool {
	return this._value.OnShelves == 1
}

// 获取商品的销售标签
func (this *Item) GetSaleTags() []*sale.SaleLabel {
	if this._saleTags == nil {
		this._saleTags = this._saleTagRep.GetItemSaleTags(this.GetDomainId())
	}
	return this._saleTags
}

// 保存销售标签
func (this *Item) SaveSaleTags(tagIds []int) error {
	err := this._saleTagRep.CleanItemSaleTags(this.GetDomainId())
	if err == nil {
		err = this._saleTagRep.SaveItemSaleTags(this.GetDomainId(), tagIds)
		this._saleTags = nil
	}
	return err
}

// 保存
func (this *Item) Save() (int, error) {
	this._sale.clearCache(this._value.Id)

	unix := time.Now().Unix()
	this._value.UpdateTime = unix

	if this.GetDomainId() <= 0 {
		this._value.CreateTime = unix
	}

	if this._value.GoodsNo == "" {
		cs := strconv.Itoa(this._value.CategoryId)
		us := strconv.Itoa(int(unix))
		l := len(cs)
		this._value.GoodsNo = fmt.Sprintf("%s%s", cs, us[4+l:])
	}

	id, err := this._saleRep.SaveValueItem(this._value)
	if err == nil {
		this._value.Id = id
		//todo: 保存商品
		this.saveGoods()

		// 创建快照
		//_, err = this.GenerateSnapshot()
	}
	return id, err
}

func (this *Item) saveGoods() {
	val := this._goodsRep.GetValueGoods(this.GetDomainId(), 0)
	if val == nil {
		val = &sale.ValueGoods{
			Id:            0,
			ItemId:        this.GetDomainId(),
			IsPresent:     0,
			SkuId:         0,
			PromotionFlag: 0,
			StockNum:      100,
			SaleNum:       100,
		}
	}
	goods := NewSaleGoods(this._sale, this, val, this._saleRep, this._goodsRep, this._promRep)
	goods.Save()
}

//// 生成快照
//func (this *Goods) GenerateSnapshot() (int, error) {
//	v := this._value
//	if v.Id <= 0 {
//		return 0, sale.ErrNoSuchGoods
//	}
//
//	if v.OnShelves == 0 {
//		return 0, sale.ErrNotOnShelves
//	}
//
//	merchantId := this._sale.GetAggregateRootId()
//	unix := time.Now().Unix()
//	cate := this._saleRep.GetCategory(merchantId, v.CategoryId)
//	var gsn *sale.GoodsSnapshot = &sale.GoodsSnapshot{
//		Key:          fmt.Sprintf("%d-g%d-%d", merchantId, v.Id, unix),
//		GoodsId:      this.GetDomainId(),
//		GoodsName:    v.Name,
//		GoodsNo:      v.GoodsNo,
//		SmallTitle:   v.SmallTitle,
//		CategoryName: cate.Name,
//		Image:        v.Image,
//		Cost:         v.Cost,
//		Price:        v.Price,
//		SalePrice:    v.SalePrice,
//		CreateTime:   unix,
//	}
//
//	if this.isNewSnapshot(gsn) {
//		this._latestSnapshot = gsn
//		return this._saleRep.SaveSnapshot(gsn)
//	}
//	return 0, sale.ErrLatestSnapshot
//}
//
//// 是否为新快照,与旧有快照进行数据对比
//func (this *Goods) isNewSnapshot(gsn *sale.GoodsSnapshot) bool {
//	latestGsn := this.GetLatestSnapshot()
//	if latestGsn != nil {
//		return latestGsn.GoodsName != gsn.GoodsName ||
//			latestGsn.SmallTitle != gsn.SmallTitle ||
//			latestGsn.CategoryName != gsn.CategoryName ||
//			latestGsn.Image != gsn.Image ||
//			latestGsn.Cost != gsn.Cost ||
//			latestGsn.Price != gsn.Price ||
//			latestGsn.SalePrice != gsn.SalePrice
//	}
//	return true
//}
//
//// 获取最新的快照
//func (this *Goods) GetLatestSnapshot() *sale.GoodsSnapshot {
//	if this._latestSnapshot == nil {
//		this._latestSnapshot = this._saleRep.GetLatestGoodsSnapshot(this.GetDomainId())
//	}
//	return this._latestSnapshot
//}
