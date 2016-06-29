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
	iRep item.IItemRep, gs *goods.ValueGoods, gi *item.Item) goods.ISnapshotManager {
	return &snapshotManagerImpl{
		_rep:   rep,
		_skuId: skuId,
		_gs:    gs,
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
			latest.SalePrice != snap.SalePrice
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

// 更新快照
func (this *snapshotManagerImpl) GenerateSnapshot() (int, error) {
	ls := this.GetLatestSnapshot()
	gs, gi := this.getGoodsAndItem()

	if this._skuId <= 0 || gi == nil || gs == nil {
		return -1, goods.ErrNoSuchGoods
	}

	if gi.OnShelves != 1 {
		//是否上架
		return -1, goods.ErrNotOnShelves
	}

	unix := time.Now().Unix()
	var gsn *goods.Snapshot = &goods.Snapshot{
		SkuId:      this._skuId,
		Key:        fmt.Sprintf("%d-g%d-%d", gi.VendorId, this._skuId, unix),
		ItemId:     gs.Id,
		GoodsId:    this._skuId,
		GoodsTitle: gi.Name,
		GoodsNo:    gi.GoodsNo,
		SmallTitle: gi.SmallTitle,
		CategoryId: gi.CategoryId,
		Image:      gi.Image,
		SalePrice:  gs.SalePrice,
		Price:      gs.Price,
		UpdateTime: unix,
	}

	if this.CompareSnapshot(gsn, ls) {
		this._latestSnapshot = gsn
		return this._rep.SaveSnapshot(gsn)
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
