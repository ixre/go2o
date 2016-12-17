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
	"time"
)

var _ item.ISnapshotManager = new(snapshotManagerImpl)

type snapshotManagerImpl struct {
	itemRepo       item.IGoodsItemRepo
	skuId          int32
	gsItem         *item.GoodsItem
	latestSnapshot *item.Snapshot
}

func NewSnapshotManagerImpl(skuId int32, repo item.IGoodsItemRepo,
	gs *item.GoodsItem) item.ISnapshotManager {
	return &snapshotManagerImpl{
		itemRepo: repo,
		skuId:    skuId,
		gsItem:   gs,
	}
}

// 获取最新的快照
func (s *snapshotManagerImpl) GetLatestSnapshot() *item.Snapshot {
	if s.latestSnapshot == nil {
		s.latestSnapshot = s.itemRepo.GetLatestSnapshot(s.skuId)
	}
	return s.latestSnapshot
}

// 是否为新快照,与旧有快照进行数据对比
func (s *snapshotManagerImpl) CompareSnapshot(snap *item.Snapshot,
	latest *item.Snapshot) bool {
	if latest != nil {
		return latest.Title != snap.Title ||
			latest.ShortTitle != snap.ShortTitle ||
			latest.CatId != snap.CatId ||
			latest.Image != snap.Image ||
			latest.Cost != snap.Cost ||
			latest.RetailPrice != snap.RetailPrice ||
			latest.Price != snap.Price ||
			latest.ExpressTid != snap.ExpressTid ||
			latest.Weight != snap.Weight ||
			latest.Bulk != snap.Bulk ||
			latest.PriceRange != snap.PriceRange ||
			latest.ShopCatId != snap.ShopCatId ||
			latest.ShortTitle != snap.ShortTitle ||
			latest.ShopId != snap.ShopId ||
			latest.ProductId != snap.ProductId
	}
	return true
}

func (s *snapshotManagerImpl) getGoodsAndItem() *item.GoodsItem {
	if s.gsItem == nil {
		s.gsItem = s.itemRepo.GetValueGoodsById(s.skuId)
	}
	return s.gsItem
}

//func (s *snapshotManagerImpl)

// 检查快照
func (s *snapshotManagerImpl) checkSnapshot(snap *item.Snapshot, it *item.GoodsItem) (err error) {
	// 检查是否更新了上架状态
	if snap != nil && snap.ShelveState != it.ShelveState {
		snap.ShelveState = it.ShelveState
		_, err = s.itemRepo.SaveSnapshot(snap)
	}
	return err
}

// 更新快照, 通过审核后,才会更新快照
func (s *snapshotManagerImpl) GenerateSnapshot() (int32, error) {
	it := s.getGoodsAndItem()
	if s.skuId <= 0 || it == nil {
		return -1, item.ErrNoSuchGoods
	}
	ls := s.GetLatestSnapshot()
	// 检查快照
	err := s.checkSnapshot(ls, it)
	// 审核通过后更新快照
	if err == nil && it.ReviewState == enum.ReviewPass {
		return s.updateSnapshot(ls, it)
	}
	return 0, err
}

// 更新快照
func (s *snapshotManagerImpl) updateSnapshot(ls *item.Snapshot, it *item.GoodsItem) (int32, error) {
	levelSales := 0
	if len(s.itemRepo.GetGoodsLevelPrice(s.skuId)) > 0 {
		levelSales = 1
	}
	unix := time.Now().Unix()
	var snap *item.Snapshot = &item.Snapshot{
		ItemId:      it.Id,
		Key:         fmt.Sprintf("%d-g%d-%d", it.VendorId, s.skuId, unix),
		CatId:       it.CatId,
		VendorId:    it.VendorId,
		BrandId:     it.BrandId,
		ProductId:   it.ProductId,
		ShopId:      it.ShopId,
		ShopCatId:   it.ShopCatId,
		ExpressTid:  it.ExpressTid,
		Title:       it.Title,
		ShortTitle:  it.ShortTitle,
		Code:        it.Code,
		Image:       it.Image,
		IsPresent:   it.IsPresent,
		PriceRange:  it.PriceRange,
		SkuId:       it.SkuId,
		Cost:        it.Cost,
		Price:       it.Price,
		RetailPrice: it.RetailPrice,
		Weight:      it.Weight,
		Bulk:        it.Bulk,
		LevelSales:  int32(levelSales),
		ShelveState: it.ShelveState,
		UpdateTime:  it.UpdateTime,
	}

	// 比较快照
	if s.CompareSnapshot(snap, ls) {
		s.latestSnapshot = snap
		return s.itemRepo.SaveSnapshot(snap)
	}
	return snap.SkuId, nil
	//return 0, goods.ErrLatestSnapshot
}

// 根据KEY获取已销售商品的快照
func (s *snapshotManagerImpl) GetSaleSnapshotByKey(key string) *item.SalesSnapshot {
	return s.itemRepo.GetSaleSnapshotByKey(key)
}

// 根据ID获取已销售商品的快照
func (s *snapshotManagerImpl) GetSaleSnapshot(id int32) *item.SalesSnapshot {
	return s.itemRepo.GetSaleSnapshot(id)
}

// 获取最新的商品销售快照,如果商品有更新,则更新销售快照
func (s *snapshotManagerImpl) GetLatestSaleSnapshot() *item.SalesSnapshot {
	snap := s.itemRepo.GetLatestSaleSnapshot(s.skuId)
	snapBasis := s.GetLatestSnapshot()
	if snap == nil || snap.CreateTime != snapBasis.UpdateTime {
		// 生成交易快照
		snap = s.createNewSaleSnap(snapBasis)
		snap.Id, _ = s.itemRepo.SaveSaleSnapshot(snap)
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
		GoodsTitle: snap.Title,
		//货号
		GoodsNo: snap.Code,
		//货品编号
		ItemId: snap.ItemId,
		//分类编号
		CategoryId: snap.CatId,
		//SKU
		Sku: "", // snap.SkuId,
		//图片
		Image: snap.Image,
		// 供货价
		Cost: snap.Cost,
		//销售价
		Price: snap.Price,
		// 快照时间
		CreateTime: snap.UpdateTime,
	}
	//快照编码: 商户编号+g商品编号+快照时间戳
	sn.SnapshotKey = fmt.Sprintf("%d-g%d-%d", sn.SellerId, sn.SkuId, sn.CreateTime)
	return sn
}
