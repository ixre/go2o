/**
 * Copyright 2015 @ z3q.net.
 * name : snapshot
 * author : jarryliu
 * date : 2016-06-28 23:52
 * description :
 * history :
 */
package item

import (
	"fmt"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"time"
)

var _ item.ISnapshotManager = new(snapshotManagerImpl)

type snapshotManagerImpl struct {
	rep            item.IGoodsRepo
	itemRepo       product.IProductRepo
	skuId          int32
	gs             *item.GoodsItem
	gi             *product.Product
	latestSnapshot *item.Snapshot
}

func NewSnapshotManagerImpl(skuId int32, rep item.IGoodsRepo,
	itemRepo product.IProductRepo, gs *item.GoodsItem, gi *product.Product) item.ISnapshotManager {
	return &snapshotManagerImpl{
		rep:      rep,
		skuId:    skuId,
		gs:       gs,
		gi:       gi,
		itemRepo: itemRepo,
	}
}

// 获取最新的快照
func (s *snapshotManagerImpl) GetLatestSnapshot() *item.Snapshot {
	if s.latestSnapshot == nil {
		s.latestSnapshot = s.rep.GetLatestSnapshot(s.skuId)
	}
	return s.latestSnapshot
}

// 是否为新快照,与旧有快照进行数据对比
func (s *snapshotManagerImpl) CompareSnapshot(snap *item.Snapshot,
	latest *item.Snapshot) bool {
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

func (s *snapshotManagerImpl) getGoodsAndItem() (*item.GoodsItem, *product.Product) {
	if s.gs == nil {
		s.gs = s.rep.GetValueGoodsById(s.skuId)
	}
	if s.gi == nil {
		s.gi = s.itemRepo.GetProductValue(s.gs.ProductId)
	}
	return s.gs, s.gi
}

//func (s *snapshotManagerImpl)

// 检查快照
func (s *snapshotManagerImpl) checkSnapshot(snap *item.Snapshot, i *product.Product) (err error) {
	// 检查是否更新了上架状态
	if snap != nil && snap.ShelveState != i.ShelveState {
		snap.ShelveState = i.ShelveState
		_, err = s.rep.SaveSnapshot(snap)
	}
	return err
}

// 更新快照, 通过审核后,才会更新快照
func (s *snapshotManagerImpl) GenerateSnapshot() (int32, error) {
	gs, gi := s.getGoodsAndItem()
	if s.skuId <= 0 || gi == nil || gs == nil {
		return -1, item.ErrNoSuchGoods
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
func (s *snapshotManagerImpl) updateSnapshot(ls *item.Snapshot,
	gi *product.Product, gs *item.GoodsItem) (int32, error) {
	LevelSales := 0
	if len(s.rep.GetGoodsLevelPrice(s.skuId)) > 0 {
		LevelSales = 1
	}
	unix := time.Now().Unix()
	var snap *item.Snapshot = &item.Snapshot{
		SkuId:        s.skuId,
		VendorId:     gi.VendorId,
		Key:          fmt.Sprintf("%d-g%d-%d", gi.VendorId, s.skuId, unix),
		ItemId:       gs.ProductId,
		GoodsTitle:   gi.Name,
		GoodsNo:      gi.Code,
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
		s.latestSnapshot = snap
		return s.rep.SaveSnapshot(snap)
	}
	return snap.SkuId, nil
	//return 0, goods.ErrLatestSnapshot
}

// 根据KEY获取已销售商品的快照
func (s *snapshotManagerImpl) GetSaleSnapshotByKey(key string) *item.SalesSnapshot {
	return s.rep.GetSaleSnapshotByKey(key)
}

// 根据ID获取已销售商品的快照
func (s *snapshotManagerImpl) GetSaleSnapshot(id int32) *item.SalesSnapshot {
	return s.rep.GetSaleSnapshot(id)
}

// 获取最新的商品销售快照,如果商品有更新,则更新销售快照
func (s *snapshotManagerImpl) GetLatestSaleSnapshot() *item.SalesSnapshot {
	snap := s.rep.GetLatestSaleSnapshot(s.skuId)
	snapBasis := s.GetLatestSnapshot()
	if snap == nil || snap.CreateTime != snapBasis.UpdateTime {
		// 生成交易快照
		snap = s.createNewSaleSnap(snapBasis)
		snap.Id, _ = s.rep.SaveSaleSnapshot(snap)
	}
	return snap
}

// 通过商品快照创建新的商品销售快照
func (s *snapshotManagerImpl) createNewSaleSnap(snap *item.Snapshot) *item.SalesSnapshot {
	sn := &item.SalesSnapshot{
		//快照编号
		Id: 0,
		//商品SKU编号
		SkuId: s.skuId,
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
