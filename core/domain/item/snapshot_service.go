/**
 * Copyright 2015 @ to2.net.
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

var _ item.ISnapshotService = new(snapshotServiceImpl)

type snapshotServiceImpl struct {
	itemRepo item.IGoodsItemRepo
}

func NewSnapshotServiceImpl(repo item.IGoodsItemRepo) item.ISnapshotService {
	return &snapshotServiceImpl{
		itemRepo: repo,
	}
}

// 获取最新的快照
func (s *snapshotServiceImpl) GetLatestSnapshot(itemId int64) *item.Snapshot {
	return s.itemRepo.GetLatestSnapshot(itemId)
}

// 是否为新快照,与旧有快照进行数据对比
func (s *snapshotServiceImpl) CompareSnapshot(snap *item.Snapshot,
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

//func (s *snapshotManagerImpl)

// 检查快照
func (s *snapshotServiceImpl) checkSnapshot(snap *item.Snapshot, it *item.GoodsItem) (err error) {
	// 检查是否更新了上架状态
	if snap != nil && snap.ShelveState != it.ShelveState {
		snap.ShelveState = it.ShelveState
		_, err = s.itemRepo.SaveSnapshot(snap)
	}
	return err
}

// 更新快照, 通过审核后,才会更新快照
func (s *snapshotServiceImpl) GenerateSnapshot(it *item.GoodsItem) (int64, error) {
	if it.ID <= 0 || it == nil {
		return -1, item.ErrNoSuchItem
	}
	ls := s.GetLatestSnapshot(it.ID)
	// 检查快照
	err := s.checkSnapshot(ls, it)
	// 审核通过后更新快照
	if err == nil && it.ReviewState == enum.ReviewPass {
		return s.updateSnapshot(ls, it)
	}
	return 0, err
}

// 更新快照
func (s *snapshotServiceImpl) updateSnapshot(ls *item.Snapshot,
	it *item.GoodsItem) (int64, error) {
	//todo: ???  SKU的会员价
	levelSales := 0
	if len(s.itemRepo.GetGoodsLevelPrice(it.ID)) > 0 {
		levelSales = 1
	}
	unix := time.Now().Unix()
	var snap *item.Snapshot = &item.Snapshot{
		ItemId:      it.ID,
		Key:         fmt.Sprintf("%d-g%d-%d", it.VendorId, it.ID, unix),
		CatId:       it.CatId,
		VendorId:    it.VendorId,
		BrandId:     it.BrandId,
		ProductId:   it.ProductId,
		ShopId:      it.ShopId,
		ShopCatId:   it.ShopCatId,
		ExpressTid:  it.ExpressTid,
		SkuId:       it.SkuId,
		Title:       it.Title,
		ShortTitle:  it.ShortTitle,
		Code:        it.Code,
		Image:       it.Image,
		IsPresent:   it.IsPresent,
		PriceRange:  it.PriceRange,
		Cost:        it.Cost,
		Price:       it.Price,
		RetailPrice: it.RetailPrice,
		Weight:      it.Weight,
		Bulk:        it.Bulk,
		LevelSales:  int32(levelSales),
		ShelveState: it.ShelveState,
		UpdateTime:  it.UpdateTime,
	}
	// 比较快照,如果为最新则更新快照
	if s.CompareSnapshot(snap, ls) {
		return s.itemRepo.SaveSnapshot(snap)
	}
	return snap.ItemId, nil
}

// 根据KEY获取已销售商品的快照
func (s *snapshotServiceImpl) GetSaleSnapshotByKey(key string) *item.TradeSnapshot {
	return s.itemRepo.GetSaleSnapshotByKey(key)
}

// 根据ID获取已销售商品的快照
func (s *snapshotServiceImpl) GetSalesSnapshot(id int64) *item.TradeSnapshot {
	return s.itemRepo.GetSalesSnapshot(id)
}

// 获取最新的商品销售快照,如果商品有更新,则更新销售快照
func (s *snapshotServiceImpl) GetLatestSalesSnapshot(itemId, skuId int64) *item.TradeSnapshot {
	snap := s.itemRepo.GetLatestSalesSnapshot(skuId)
	snapBasis := s.GetLatestSnapshot(itemId)
	if snap == nil || snap.CreateTime != snapBasis.UpdateTime {
		// 生成交易快照
		snap = s.createNewSaleSnap(skuId, snapBasis)
		snap.Id, _ = s.itemRepo.SaveSalesSnapshot(snap)
	}
	return snap
}

// 通过商品快照创建新的商品销售快照
func (s *snapshotServiceImpl) createNewSaleSnap(skuId int64, snap *item.Snapshot) *item.TradeSnapshot {
	sn := &item.TradeSnapshot{
		//快照编号
		Id: 0,
		//商品编号
		ItemId: snap.ItemId,
		//商品SKU编号
		SkuId: skuId,
		// 卖家编号
		SellerId: snap.VendorId,
		//商品标题
		GoodsTitle: snap.Title,
		//货号
		GoodsNo: snap.Code,
		//分类编号
		CategoryId: snap.CatId,
		//SKU
		Sku: "",
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
	sn.SnapshotKey = fmt.Sprintf("%d-g%d-%d", sn.SellerId, skuId, sn.CreateTime)
	// 绑定SKU的信息
	if skuId > 0 {
		if sku := s.itemRepo.GetItemSku(skuId); sku != nil {
			sn.Sku = sku.SpecWord
			sn.Price = sku.Price
			sn.Cost = sku.Cost
			if sku.Image != "" {
				sn.Image = sku.Image
			}
			if sku.Title != "" {
				sn.GoodsTitle = sku.Title
			}
			if sku.Code != "" {
				sn.GoodsNo = sku.Code
			}
		}
	}
	return sn
}
