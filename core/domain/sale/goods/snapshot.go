/**
 * Copyright 2015 @ z3q.net.
 * name : snapshot
 * author : jarryliu
 * date : 2016-06-28 23:52
 * description :
 * history :
 */
package goods

import (
	"fmt"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"time"
)

var _ goods.ISnapshotManager = new(snapshotManagerImpl)

type snapshotManagerImpl struct {
	_rep            goods.IGoodsRep
	_itemRep        item.IItemRep
	_skuId          int
	_gs             *goods.ValueGoods
	_gi             *item.Item
	_latestSnapshot *goods.Snapshot
}

func NewSnapshotManagerImpl(skuId int, rep goods.IGoodsRep,
	itemRep item.IItemRep, gs *goods.ValueGoods, gi *item.Item) goods.ISnapshotManager {
	return &snapshotManagerImpl{
		_rep:     rep,
		_skuId:   skuId,
		_gs:      gs,
		_itemRep: itemRep,
	}
}

// 获取最新的快照
func (this *snapshotManagerImpl) GetLatestSnapshot() *goods.Snapshot {
	if this._latestSnapshot == nil {
		this._latestSnapshot = this._rep.GetLatestSnapshot(this._skuId)
	}
	return this._latestSnapshot
}

// 是否为新快照,与旧有快照进行数据对比
func (this *snapshotManagerImpl) CompareSnapshot(snap *goods.Snapshot,
	latest *goods.Snapshot) bool {
	if latest != nil {
		return latest.GoodsTitle != snap.GoodsTitle ||
			latest.SmallTitle != snap.SmallTitle ||
			latest.CategoryId != snap.CategoryId ||
			latest.Image != snap.Image ||
			latest.Price != snap.Price ||
			latest.SalePrice != snap.SalePrice ||
			latest.OnShelves != snap.OnShelves ||
			latest.LevelSales != snap.LevelSales ||
			latest.SaleNum != snap.SaleNum ||
			latest.StockNum != snap.StockNum
	}
	return true
}

func (this *snapshotManagerImpl) getGoodsAndItem() (*goods.ValueGoods, *item.Item) {
	if this._gs == nil {
		this._gs = this._rep.GetValueGoodsById(this._skuId)
	}
	if this._gi == nil {
		this._gi = this._itemRep.GetValueItem(this._gs.ItemId)
	}
	return this._gs, this._gi
}

// 更新快照, 通过审核后,才会更新快照
func (this *snapshotManagerImpl) GenerateSnapshot() (int, error) {
	ls := this.GetLatestSnapshot()
	gs, gi := this.getGoodsAndItem()

	if this._skuId <= 0 || gi == nil || gs == nil {
		return -1, goods.ErrNoSuchGoods
	}

	// 是否审核通过
	if gi.ReviewPass == 0 {
		return -1, item.ErrNotBeReview
	}

	LevelSales := 0
	if len(this._rep.GetGoodsLevelPrice(this._skuId)) > 0 {
		LevelSales = 1
	}

	unix := time.Now().Unix()
	var snap *goods.Snapshot = &goods.Snapshot{
		SkuId:      this._skuId,
		VendorId:   gi.VendorId,
		Key:        fmt.Sprintf("%d-g%d-%d", gi.VendorId, this._skuId, unix),
		ItemId:     gs.ItemId,
		GoodsTitle: gi.Name,
		GoodsNo:    gi.GoodsNo,
		SmallTitle: gi.SmallTitle,
		CategoryId: gi.CategoryId,
		Image:      gi.Image,
		SalePrice:  gs.SalePrice,
		Price:      gi.Price,
		SaleNum:    gs.SaleNum,
		StockNum:   gs.StockNum,
		LevelSales: LevelSales,
		OnShelves:  gi.OnShelves,
		UpdateTime: unix,
	}

	if this.CompareSnapshot(snap, ls) {
		this._latestSnapshot = snap
		return this._rep.SaveSnapshot(snap)
	}

	return 0, goods.ErrLatestSnapshot
}

// 生成交易快照
func (this *snapshotManagerImpl) GenerateSaleSnapshot() (int, error) {
	return -1, nil
}

// 根据KEY获取已销售商品的快照
func (this *snapshotManagerImpl) GetSaleSnapshotByKey(key string) *goods.GoodsSnapshot {
	return this._rep.GetSaleSnapshotByKey(key)
}

// 根据ID获取已销售商品的快照
func (this *snapshotManagerImpl) GetSaleSnapshot(id int) *goods.GoodsSnapshot {
	return this._rep.GetSaleSnapshot(id)
}
