/**
 * Copyright 2015 @ z3q.net.
 * name : sale_goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package item

import (
	"go2o/core/domain"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
)

var _ item.IGoods = new(tmpGoodsImpl)
var _ domain.IDomain = new(tmpGoodsImpl)

// 临时的商品实现  todo: 要与item分开
type tmpGoodsImpl struct {
	manager       *goodsManagerImpl
	pro           product.IProduct
	value         *item.GoodsItem
	goodsRepo     item.IGoodsRepo
	productRepo   product.IProductRepo
	promRepo      promotion.IPromotionRepo
	levelPrices   []*item.MemberPrice
	promDescribes map[string]string
	snapManager   item.ISnapshotManager
}

func NewSaleGoods(m *goodsManagerImpl,
	itemRepo product.IProductRepo, pro product.IProduct,
	value *item.GoodsItem,
	goodsRepo item.IGoodsRepo, promRepo promotion.IPromotionRepo) item.IGoods {
	v := &tmpGoodsImpl{
		manager:     m,
		pro:         pro,
		value:       value,
		productRepo: itemRepo,
		goodsRepo:   goodsRepo,
		promRepo:    promRepo,
	}
	return v.init()
}

func (g *tmpGoodsImpl) init() item.IGoods {
	g.value.Price = g.value.Price
	if g.pro != nil {
		g.value.Price = g.pro.GetValue().SalePrice
		g.value.PromPrice = g.pro.GetValue().SalePrice
	}
	return g
}

//获取领域对象编号
func (g *tmpGoodsImpl) GetDomainId() int32 {
	return g.value.Id
}

// 商品快照
func (g *tmpGoodsImpl) SnapshotManager() item.ISnapshotManager {
	if g.snapManager == nil {
		var item *product.Product
		gi := g.GetItem()
		if gi != nil {
			v := gi.GetValue()
			item = &v
		}
		g.snapManager = NewSnapshotManagerImpl(g.GetDomainId(),
			g.goodsRepo, g.productRepo, g.GetValue(), item)
	}
	return g.snapManager
}

// 获取货品
func (g *tmpGoodsImpl) GetItem() product.IProduct {
	return g.pro
}

// 设置值
func (g *tmpGoodsImpl) GetValue() *item.GoodsItem {
	return g.value
}

// 获取包装过的商品信息
func (g *tmpGoodsImpl) GetPackedValue() *valueobject.Goods {
	item := g.GetItem().GetValue()
	gv := g.GetValue()
	goods := &valueobject.Goods{
		ProductId:     item.Id,
		CategoryId:    item.CategoryId,
		Name:          item.Name,
		GoodsNo:       item.Code,
		Image:         item.Image,
		Price:         item.Price,
		SalePrice:     item.SalePrice,
		PromPrice:     item.SalePrice,
		GoodsId:       g.GetDomainId(),
		SkuId:         gv.SkuId,
		IsPresent:     gv.IsPresent,
		PromotionFlag: gv.PromFlag,
		StockNum:      gv.StockNum,
		SaleNum:       gv.SaleNum,
	}
	return goods
}

// 获取促销信息
func (g *tmpGoodsImpl) GetPromotions() []promotion.IPromotion {
	var vp []*promotion.PromotionInfo = g.promRepo.GetPromotionOfGoods(
		g.GetDomainId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = g.promRepo.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (g *tmpGoodsImpl) GetLevelPrice(level int32) (bool, float32) {
	lvp := g.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < g.value.Price {
			return true, v.Price
		}
	}
	return false, g.value.Price
}

// 获取促销价
func (g *tmpGoodsImpl) GetPromotionPrice(level int32) float32 {
	b, price := g.GetLevelPrice(level)
	if b {
		return price
	}
	return g.value.Price
}

// 获取促销描述
func (g *tmpGoodsImpl) GetPromotionDescribe() map[string]string {
	if g.promDescribes == nil {
		proms := g.GetPromotions()
		g.promDescribes = make(map[string]string, len(proms))
		for _, v := range proms {
			key := v.TypeName()
			if txt, ok := g.promDescribes[key]; !ok {
				g.promDescribes[key] = v.GetValue().ShortName
			} else {
				g.promDescribes[key] = txt + "；" + v.GetValue().ShortName
			}

			//			if v.Type() == promotion.TypeFlagCashBack {
			//				if txt, ok := g._promDescribes[key]; !ok {
			//					g._promDescribes[key] = v.GetValue().ShortName
			//				} else {
			//					g._promDescribes[key] = txt + "；" + v.GetValue().ShortName
			//				}
			//			} else if v.Type() == promotion.TypeFlagCoupon {
			//				if txt, ok := g._promDescribes[key]; !ok {
			//					g._promDescribes[key] = v.GetValue().ShortName
			//				} else {
			//					g._promDescribes[key] = txt + "；" + v.GetValue().ShortName
			//				}
			//			}

			//todo: other promotion implement
		}
	}
	return g.promDescribes
}

// 获取会员价
func (g *tmpGoodsImpl) GetLevelPrices() []*item.MemberPrice {
	if g.levelPrices == nil {
		g.levelPrices = g.goodsRepo.GetGoodsLevelPrice(g.GetDomainId())
	}
	return g.levelPrices
}

// 保存会员价
func (g *tmpGoodsImpl) SaveLevelPrice(v *item.MemberPrice) (int32, error) {
	v.GoodsId = g.GetDomainId()
	if g.value.Price == v.Price {
		if v.Id > 0 {
			g.goodsRepo.RemoveGoodsLevelPrice(v.Id)
		}
		return -1, nil
	}
	return g.goodsRepo.SaveGoodsLevelPrice(v)
}

// 设置值
func (g *tmpGoodsImpl) SetValue(v *item.GoodsItem) error {
	g.value.IsPresent = v.IsPresent
	g.value.SaleNum = v.SaleNum
	g.value.StockNum = v.StockNum
	g.value.SkuId = v.SkuId
	//g._value.PromotionFlag = v.PromotionFlag
	return nil
}

// 保存
func (g *tmpGoodsImpl) Save() (int32, error) {
	id, err := g.goodsRepo.SaveValueGoods(g.value)
	if err == nil {
		g.value.Id = id
		_, err = g.SnapshotManager().GenerateSnapshot()
	}
	//todo: save promotion
	return id, err
}

// 更新销售数量
func (g *tmpGoodsImpl) AddSalesNum(quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return item.ErrOutOfStock
	}
	g.value.SaleNum += quantity
	_, err := g.Save()
	return err
}

// 取消销售
func (g *tmpGoodsImpl) CancelSale(quantity int32, orderNo string) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	g.value.SaleNum -= quantity
	_, err := g.Save()
	return err
}

// 占用库存
func (g *tmpGoodsImpl) TakeStock(quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return item.ErrOutOfStock
	}
	g.value.StockNum -= quantity
	_, err := g.Save()
	return err
}

// 释放库存
func (g *tmpGoodsImpl) FreeStock(quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	g.value.StockNum += quantity
	_, err := g.Save()
	return err
}

var _ item.IGoodsManager = new(goodsManagerImpl)

type goodsManagerImpl struct {
	_productRepo product.IProductRepo
	_itemRepo    item.IGoodsRepo
	_promRepo    promotion.IPromotionRepo
	_valRepo     valueobject.IValueRepo
	_mchId       int32
}

func NewGoodsManager(mchId int32,
	itemRepo item.IGoodsRepo,
	productRepo product.IProductRepo,
	promRepo promotion.IPromotionRepo,
	valRepo valueobject.IValueRepo) item.IGoodsManager {
	c := &goodsManagerImpl{
		_mchId:       mchId,
		_productRepo: productRepo,
		_itemRepo:    itemRepo,
		_promRepo:    promRepo,
		_valRepo:     valRepo,
	}
	return c.init()
}

func (g *goodsManagerImpl) init() item.IGoodsManager {
	return g
}

// 创建商品
func (g *goodsManagerImpl) CreateGoods(v *item.GoodsItem) item.IGoods {
	return NewSaleGoods(g, g._productRepo,
		nil, v, g._itemRepo, g._promRepo)
}

// 创建商品
func (g *goodsManagerImpl) CreateGoodsByItem(pro product.IProduct, v *item.GoodsItem) item.IGoods {
	return NewSaleGoods(g, g._productRepo,
		pro, v, g._itemRepo, g._promRepo)
}

// 根据产品编号获取商品
func (g *goodsManagerImpl) GetGoods(goodsId int32) item.IGoods {
	var v *item.GoodsItem = g._itemRepo.GetValueGoodsById(goodsId)
	if v != nil {
		pro := g._productRepo.GetProduct(v.ProductId)
		if pro != nil {
			return g.CreateGoodsByItem(pro, v)
		}
	}
	return nil
}

// 根据产品SKU获取商品
func (g *goodsManagerImpl) GetGoodsBySku(itemId, skuId int32) item.IGoods {
	var v *item.GoodsItem = g._itemRepo.GetValueGoodsBySku(itemId, skuId)
	if v != nil {
		pro := g._productRepo.GetProduct(v.ProductId)
		if pro != nil {
			return g.CreateGoodsByItem(pro, v)
		}
	}
	return nil
}

// 删除商品
func (g *goodsManagerImpl) DeleteGoods(goodsId int32) error {
	gs := g.GetGoods(goodsId)
	if gs.GetValue().SaleNum > 0 {
		return item.ErrNoSuchSnapshot
	}
	//todo: delete goods
	return g._productRepo.DeleteProduct(gs.GetValue().ProductId)
}

// 获取指定数量已上架的商品
func (g *goodsManagerImpl) GetOnShelvesGoods(start, end int,
	sortBy string) []*valueobject.Goods {
	return g._itemRepo.GetOnShelvesGoods(g._mchId,
		start, end, sortBy)
}
