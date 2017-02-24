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
	"go2o/core/domain/interface/enum"
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
		_gi:      gi,
		_itemRep: itemRep,
	}
}

// 获取最新的快照
func (s *snapshotManagerImpl) GetLatestSnapshot() *goods.Snapshot {
	if s._latestSnapshot == nil {
		s._latestSnapshot = s._rep.GetLatestSnapshot(s._skuId)
	}
	return s._latestSnapshot
}

// 是否为新快照,与旧有快照进行数据对比
func (s *snapshotManagerImpl) CompareSnapshot(snap *goods.Snapshot,
	latest *goods.Snapshot) bool {
	if latest != nil {
		return latest.GoodsTitle != snap.GoodsTitle ||
			latest.SmallTitle != snap.SmallTitle ||
			latest.CategoryId != snap.CategoryId ||
			latest.Image != snap.Image ||
			latest.Cost != snap.Cost ||
			latest.Price != snap.Price ||
			latest.SalePrice != snap.SalePrice ||
			latest.ShelveState != snap.ShelveState ||
			latest.LevelSales != snap.LevelSales ||
			latest.SaleNum != snap.SaleNum ||
			latest.StockNum != snap.StockNum ||
			latest.ExpressTplId != snap.ExpressTplId ||
			latest.Weight != snap.Weight
	}
	return true
}

func (s *snapshotManagerImpl) getGoodsAndItem() (*goods.ValueGoods, *item.Item) {
	if s._gs == nil {
		s._gs = s._rep.GetValueGoodsById(s._skuId)
	}
	if s._gi == nil {
		s._gi = s._itemRep.GetValueItem(s._gs.ItemId)
	}
	return s._gs, s._gi
}

//func (s *snapshotManagerImpl)

// 检查快照
func (s *snapshotManagerImpl) checkSnapshot(snap *goods.Snapshot, i *item.Item) (err error) {
	// 检查是否更新了上架状态
	if snap != nil && snap.ShelveState != i.ShelveState {
		snap.ShelveState = i.ShelveState
		_, err = s._rep.SaveSnapshot(snap)
	}
	return err
}

// 更新快照, 通过审核后,才会更新快照
func (s *snapshotManagerImpl) GenerateSnapshot() (int, error) {
	gs, gi := s.getGoodsAndItem()
	if s._skuId <= 0 || gi == nil || gs == nil {
		return -1, goods.ErrNoSuchGoods
	}
	ls := s.GetLatestSnapshot()
	// 检查快照
	err := s.checkSnapshot(ls, gi)
	// 审核通过后更新快照
	if err == nil && gi.ReviewState == enum.ReviewPass {
		return s.updateSnapshot(ls, gi, gs)
	}
	return 0, err
}

// 更新快照
func (s *snapshotManagerImpl) updateSnapshot(ls *goods.Snapshot,
	gi *item.Item, gs *goods.ValueGoods) (int, error) {
	LevelSales := 0
	if len(s._rep.GetGoodsLevelPrice(s._skuId)) > 0 {
		LevelSales = 1
	}
	unix := time.Now().Unix()
	var snap *goods.Snapshot = &goods.Snapshot{
		SkuId:        s._skuId,
		VendorId:     gi.VendorId,
		Key:          fmt.Sprintf("%d-g%d-%d", gi.VendorId, s._skuId, unix),
		ItemId:       gs.ItemId,
		GoodsTitle:   gi.Name,
		GoodsNo:      gi.GoodsNo,
		SmallTitle:   gi.SmallTitle,
		CategoryId:   gi.CategoryId,
		Image:        gi.Image,
		Weight:       gi.Weight,
		SalePrice:    gs.SalePrice,
		Cost:         gi.Cost,
		Price:        gi.Price,
		SaleNum:      gs.SaleNum,
		StockNum:     gs.StockNum,
		LevelSales:   LevelSales,
		ShelveState:  gi.ShelveState,
		ExpressTplId: gi.ExpressTplId,
		UpdateTime:   unix,
	}
	// 比较快照
	if s.CompareSnapshot(snap, ls) {
		s._latestSnapshot = snap
		return s._rep.SaveSnapshot(snap)
	}
	return snap.SkuId, nil
	//return 0, goods.ErrLatestSnapshot
}

// 根据KEY获取已销售商品的快照
func (s *snapshotManagerImpl) GetSaleSnapshotByKey(key string) *goods.SalesSnapshot {
	return s._rep.GetSaleSnapshotByKey(key)
}

// 根据ID获取已销售商品的快照
func (s *snapshotManagerImpl) GetSaleSnapshot(id int) *goods.SalesSnapshot {
	return s._rep.GetSaleSnapshot(id)
}

// 获取最新的商品销售快照,如果商品有更新,则更新销售快照
func (s *snapshotManagerImpl) GetLatestSaleSnapshot() *goods.SalesSnapshot {
	snap := s._rep.GetLatestSaleSnapshot(s._skuId)
	snapBasis := s.GetLatestSnapshot()
	if snap == nil || snap.CreateTime != snapBasis.UpdateTime {
		// 生成交易快照
		snap = s.createNewSaleSnap(snapBasis)
		snap.Id, _ = s._rep.SaveSaleSnapshot(snap)
	}
	return snap
}

// 通过商品快照创建新的商品销售快照
func (s *snapshotManagerImpl) createNewSaleSnap(snap *goods.Snapshot) *goods.SalesSnapshot {
	sn := &goods.SalesSnapshot{
		//快照编号
		Id: 0,
		//商品SKU编号
		SkuId: s._skuId,
		// 卖家编号
		SellerId: snap.VendorId,
		//商品标题
		GoodsTitle: snap.GoodsTitle,
		//货号
		GoodsNo: snap.GoodsNo,
		//货品编号
		ItemId: snap.ItemId,
		//分类编号
		CategoryId: snap.CategoryId,
		//SKU
		Sku: snap.Sku,
		//图片
		Image: snap.Image,
		// 供货价
		Cost: snap.Cost,
		//销售价
		Price: snap.SalePrice,
		// 快照时间
		CreateTime: snap.UpdateTime,
	}
	//快照编码: 商户编号+g商品编号+快照时间戳
	sn.SnapshotKey = fmt.Sprintf("%d-g%d-%d", sn.SellerId, sn.SkuId, sn.CreateTime)
	return sn
}
