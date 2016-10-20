/**
 * Copyright 2015 @ z3q.net.
 * name : sale_goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import (
	"go2o/core/domain"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
	goodsImpl "go2o/core/domain/sale/goods"
)

var _ sale.IGoods = new(tmpGoodsImpl)
var _ domain.IDomain = new(tmpGoodsImpl)

// 临时的商品实现  todo: 要与item分开
type tmpGoodsImpl struct {
	manager       *goodsManagerImpl
	goods         sale.IItem
	value         *goods.ValueGoods
	saleRep       sale.ISaleRep
	goodsRep      goods.IGoodsRep
	itemRep       item.IItemRep
	promRep       promotion.IPromotionRep
	sale          sale.ISale
	levelPrices   []*goods.MemberPrice
	promDescribes map[string]string
	snapManager   goods.ISnapshotManager
}

func NewSaleGoods(m *goodsManagerImpl, s sale.ISale,
	itemRep item.IItemRep, goods sale.IItem,
	value *goods.ValueGoods, rep sale.ISaleRep,
	goodsRep goods.IGoodsRep, promRep promotion.IPromotionRep) sale.IGoods {
	v := &tmpGoodsImpl{
		manager:  m,
		goods:    goods,
		value:    value,
		saleRep:  rep,
		itemRep:  itemRep,
		goodsRep: goodsRep,
		promRep:  promRep,
		sale:     s,
	}
	return v.init()
}

func (g *tmpGoodsImpl) init() sale.IGoods {
	g.value.Price = g.value.Price
	if g.goods != nil {
		g.value.SalePrice = g.goods.GetValue().SalePrice
		g.value.PromPrice = g.goods.GetValue().SalePrice
	}
	return g
}

//获取领域对象编号
func (g *tmpGoodsImpl) GetDomainId() int {
	return g.value.Id
}

// 商品快照
func (g *tmpGoodsImpl) SnapshotManager() goods.ISnapshotManager {
	if g.snapManager == nil {
		var item *item.Item
		gi := g.GetItem()
		if gi != nil {
			v := gi.GetValue()
			item = &v
		}
		g.snapManager = goodsImpl.NewSnapshotManagerImpl(g.GetDomainId(),
			g.goodsRep, g.itemRep, g.GetValue(), item)
	}
	return g.snapManager
}

// 获取货品
func (g *tmpGoodsImpl) GetItem() sale.IItem {
	return g.goods
}

// 设置值
func (g *tmpGoodsImpl) GetValue() *goods.ValueGoods {
	return g.value
}

// 获取包装过的商品信息
func (g *tmpGoodsImpl) GetPackedValue() *valueobject.Goods {
	item := g.GetItem().GetValue()
	gv := g.GetValue()
	goods := &valueobject.Goods{
		Item_Id:       item.Id,
		CategoryId:    item.CategoryId,
		Name:          item.Name,
		GoodsNo:       item.GoodsNo,
		Image:         item.Image,
		Price:         item.Price,
		SalePrice:     item.SalePrice,
		PromPrice:     item.SalePrice,
		GoodsId:       g.GetDomainId(),
		SkuId:         gv.SkuId,
		IsPresent:     gv.IsPresent,
		PromotionFlag: gv.PromotionFlag,
		StockNum:      gv.StockNum,
		SaleNum:       gv.SaleNum,
	}
	return goods
}

// 获取促销信息
func (g *tmpGoodsImpl) GetPromotions() []promotion.IPromotion {
	var vp []*promotion.PromotionInfo = g.promRep.GetPromotionOfGoods(
		g.GetDomainId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = g.promRep.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (g *tmpGoodsImpl) GetLevelPrice(level int) (bool, float32) {
	lvp := g.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < g.value.SalePrice {
			return true, v.Price
		}
	}
	return false, g.value.SalePrice
}

// 获取促销价
func (g *tmpGoodsImpl) GetPromotionPrice(level int) float32 {
	b, price := g.GetLevelPrice(level)
	if b {
		return price
	}
	return g.value.SalePrice
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
func (g *tmpGoodsImpl) GetLevelPrices() []*goods.MemberPrice {
	if g.levelPrices == nil {
		g.levelPrices = g.goodsRep.GetGoodsLevelPrice(g.GetDomainId())
	}
	return g.levelPrices
}

// 保存会员价
func (g *tmpGoodsImpl) SaveLevelPrice(v *goods.MemberPrice) (int, error) {
	v.GoodsId = g.GetDomainId()
	if g.value.SalePrice == v.Price {
		if v.Id > 0 {
			g.goodsRep.RemoveGoodsLevelPrice(v.Id)
		}
		return -1, nil
	}
	return g.goodsRep.SaveGoodsLevelPrice(v)
}

// 设置值
func (g *tmpGoodsImpl) SetValue(v *goods.ValueGoods) error {
	g.value.IsPresent = v.IsPresent
	g.value.SaleNum = v.SaleNum
	g.value.StockNum = v.StockNum
	g.value.SkuId = v.SkuId
	//g._value.PromotionFlag = v.PromotionFlag
	return nil
}

// 保存
func (g *tmpGoodsImpl) Save() (int, error) {
	id, err := g.goodsRep.SaveValueGoods(g.value)
	if err == nil {
		g.value.Id = id
		_, err = g.SnapshotManager().GenerateSnapshot()
	}
	//todo: save promotion
	return id, err
}

// 更新销售数量
func (g *tmpGoodsImpl) AddSalesNum(quantity int) error {
	if quantity <= 0 {
		return sale.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return sale.ErrOutOfStock
	}
	g.value.SaleNum += quantity
	_, err := g.Save()
	return err
}

// 取消销售
func (g *tmpGoodsImpl) CancelSale(quantity int, orderNo string) error {
	if quantity <= 0 {
		return sale.ErrGoodsNum
	}
	g.value.SaleNum -= quantity
	_, err := g.Save()
	return err
}

// 占用库存
func (g *tmpGoodsImpl) TakeStock(quantity int) error {
	if quantity <= 0 {
		return sale.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return sale.ErrOutOfStock
	}
	g.value.StockNum -= quantity
	_, err := g.Save()
	return err
}

// 释放库存
func (g *tmpGoodsImpl) FreeStock(quantity int) error {
	if quantity <= 0 {
		return sale.ErrGoodsNum
	}
	g.value.StockNum += quantity
	_, err := g.Save()
	return err
}

var _ sale.IGoodsManager = new(goodsManagerImpl)

type goodsManagerImpl struct {
	_sale   *saleImpl
	_valRep valueobject.IValueRep
	_mchId  int
}

func NewGoodsManager(mchId int, s *saleImpl,
	valRep valueobject.IValueRep) sale.IGoodsManager {
	c := &goodsManagerImpl{
		_sale:   s,
		_mchId:  mchId,
		_valRep: valRep,
	}
	return c.init()
}

func (g *goodsManagerImpl) init() sale.IGoodsManager {
	return g
}

// 创建商品
func (g *goodsManagerImpl) CreateGoods(s *goods.ValueGoods) sale.IGoods {
	return NewSaleGoods(g, g._sale, g._sale.itemRep,
		nil, s, g._sale.saleRep,
		g._sale.goodsRep, g._sale.promRep)
}

// 创建商品
func (g *goodsManagerImpl) CreateGoodsByItem(item sale.IItem, v *goods.ValueGoods) sale.IGoods {
	return NewSaleGoods(g, g._sale, g._sale.itemRep,
		item, v, g._sale.saleRep, g._sale.goodsRep,
		g._sale.promRep)
}

// 根据产品编号获取商品
func (g *goodsManagerImpl) GetGoods(goodsId int) sale.IGoods {
	var v *goods.ValueGoods = g._sale.goodsRep.GetValueGoodsById(goodsId)
	if v != nil {
		pv := g._sale.itemRep.GetValueItem(v.ItemId)
		if pv != nil {
			return g.CreateGoodsByItem(g._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 根据产品SKU获取商品
func (g *goodsManagerImpl) GetGoodsBySku(itemId, sku int) sale.IGoods {
	var v *goods.ValueGoods = g._sale.goodsRep.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		pv := g._sale.itemRep.GetValueItem(v.ItemId)
		if pv != nil {
			return g.CreateGoodsByItem(g._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 删除商品
func (g *goodsManagerImpl) DeleteGoods(goodsId int) error {
	gs := g.GetGoods(goodsId)
	if gs.GetValue().SaleNum > 0 {
		return goods.ErrNoSuchSnapshot
	}

	//todo: delete goods
	err := g._sale.itemRep.DeleteItem(g._mchId, goodsId)
	if err != nil {
		g._sale.clearCache(goodsId)
	}
	return err
}

// 获取指定数量已上架的商品
func (g *goodsManagerImpl) GetOnShelvesGoods(start, end int,
	sortBy string) []*valueobject.Goods {
	return g._sale.goodsRep.GetOnShelvesGoods(g._mchId,
		start, end, sortBy)
}
