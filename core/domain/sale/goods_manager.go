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
	_manager       *goodsManagerImpl
	_goods         sale.IItem
	_value         *goods.ValueGoods
	_saleRep       sale.ISaleRep
	_goodsRep      goods.IGoodsRep
	_itemRep       item.IItemRep
	_promRep       promotion.IPromotionRep
	_sale          sale.ISale
	_levelPrices   []*goods.MemberPrice
	_promDescribes map[string]string
	_snapManager   goods.ISnapshotManager
}

func NewSaleGoods(m *goodsManagerImpl, s sale.ISale,
	itemRep item.IItemRep, goods sale.IItem,
	value *goods.ValueGoods, rep sale.ISaleRep,
	goodsRep goods.IGoodsRep, promRep promotion.IPromotionRep) sale.IGoods {
	v := &tmpGoodsImpl{
		_manager:  m,
		_goods:    goods,
		_value:    value,
		_saleRep:  rep,
		_itemRep:  itemRep,
		_goodsRep: goodsRep,
		_promRep:  promRep,
		_sale:     s,
	}
	return v.init()
}

func (g *tmpGoodsImpl) init() sale.IGoods {
	g._value.Price = g._value.Price
	if g._goods != nil {
		g._value.SalePrice = g._goods.GetValue().SalePrice
		g._value.PromPrice = g._goods.GetValue().SalePrice
	}
	return g
}

//获取领域对象编号
func (g *tmpGoodsImpl) GetDomainId() int {
	return g._value.Id
}

// 商品快照
func (g *tmpGoodsImpl) SnapshotManager() goods.ISnapshotManager {
	if g._snapManager == nil {
		var item *item.Item
		gi := g.GetItem()
		if gi != nil {
			v := gi.GetValue()
			item = &v
		}
		g._snapManager = goodsImpl.NewSnapshotManagerImpl(g.GetDomainId(),
			g._goodsRep, g._itemRep, g.GetValue(), item)
	}
	return g._snapManager
}

// 获取货品
func (g *tmpGoodsImpl) GetItem() sale.IItem {
	return g._goods
}

// 设置值
func (g *tmpGoodsImpl) GetValue() *goods.ValueGoods {
	return g._value
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
	var vp []*promotion.PromotionInfo = g._promRep.GetPromotionOfGoods(
		g.GetDomainId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = g._promRep.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (g *tmpGoodsImpl) GetLevelPrice(level int) (bool, float32) {
	lvp := g.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < g._value.SalePrice {
			return true, v.Price
		}
	}
	return false, g._value.SalePrice
}

// 获取促销价
func (g *tmpGoodsImpl) GetPromotionPrice(level int) float32 {
	b, price := g.GetLevelPrice(level)
	if b {
		return price
	}
	return g._value.SalePrice
}

// 获取促销描述
func (g *tmpGoodsImpl) GetPromotionDescribe() map[string]string {
	if g._promDescribes == nil {
		proms := g.GetPromotions()
		g._promDescribes = make(map[string]string, len(proms))
		for _, v := range proms {
			key := v.TypeName()
			if txt, ok := g._promDescribes[key]; !ok {
				g._promDescribes[key] = v.GetValue().ShortName
			} else {
				g._promDescribes[key] = txt + "；" + v.GetValue().ShortName
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
	return g._promDescribes
}

// 获取会员价
func (g *tmpGoodsImpl) GetLevelPrices() []*goods.MemberPrice {
	if g._levelPrices == nil {
		g._levelPrices = g._goodsRep.GetGoodsLevelPrice(g.GetDomainId())
	}
	return g._levelPrices
}

// 保存会员价
func (g *tmpGoodsImpl) SaveLevelPrice(v *goods.MemberPrice) (int, error) {
	v.GoodsId = g.GetDomainId()
	if g._value.SalePrice == v.Price {
		if v.Id > 0 {
			g._goodsRep.RemoveGoodsLevelPrice(v.Id)
		}
		return -1, nil
	}
	return g._goodsRep.SaveGoodsLevelPrice(v)
}

// 设置值
func (g *tmpGoodsImpl) SetValue(v *goods.ValueGoods) error {
	g._value.IsPresent = v.IsPresent
	g._value.SaleNum = v.SaleNum
	g._value.StockNum = v.StockNum
	g._value.SkuId = v.SkuId
	//g._value.PromotionFlag = v.PromotionFlag
	return nil
}

// 保存
func (g *tmpGoodsImpl) Save() (int, error) {
	id, err := g._goodsRep.SaveValueGoods(g._value)
	if err == nil {
		g._value.Id = id
		_, err = g.SnapshotManager().GenerateSnapshot()
	}
	//todo: save promotion
	return id, err
}

// 更新销售数量
func (g *tmpGoodsImpl) AddSaleNum(quantity int) error {
	// 减去库存
	if quantity > 0 {
		if quantity > g._value.StockNum {
			return sale.ErrOutOfStock
		}
		g._value.StockNum -= quantity
		g._value.SaleNum += quantity
		_, err := g.Save()
		return err
	}
	return sale.ErrGoodsNum
}

// 取消销售
func (g *tmpGoodsImpl) CancelSale(quantity int, orderNo string) error {
	// 减去库存
	if quantity > 0 {
		g._value.StockNum += quantity
		g._value.SaleNum -= quantity
		_, err := g.Save()
		return err
	}
	return sale.ErrGoodsNum
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
	return NewSaleGoods(g, g._sale, g._sale._itemRep,
		nil, s, g._sale._saleRep,
		g._sale._goodsRep, g._sale._promRep)
}

// 创建商品
func (g *goodsManagerImpl) CreateGoodsByItem(item sale.IItem, v *goods.ValueGoods) sale.IGoods {
	return NewSaleGoods(g, g._sale, g._sale._itemRep,
		item, v, g._sale._saleRep, g._sale._goodsRep,
		g._sale._promRep)
}

// 根据产品编号获取商品
func (g *goodsManagerImpl) GetGoods(goodsId int) sale.IGoods {
	var v *goods.ValueGoods = g._sale._goodsRep.GetValueGoodsById(goodsId)
	if v != nil {
		pv := g._sale._itemRep.GetValueItem(v.ItemId)
		if pv != nil {
			return g.CreateGoodsByItem(g._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 根据产品SKU获取商品
func (g *goodsManagerImpl) GetGoodsBySku(itemId, sku int) sale.IGoods {
	var v *goods.ValueGoods = g._sale._goodsRep.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		pv := g._sale._itemRep.GetValueItem(v.ItemId)
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
	err := g._sale._itemRep.DeleteItem(g._mchId, goodsId)
	if err != nil {
		g._sale.clearCache(goodsId)
	}
	return err
}

// 获取指定数量已上架的商品
func (g *goodsManagerImpl) GetOnShelvesGoods(start, end int,
	sortBy string) []*valueobject.Goods {
	return g._sale._goodsRep.GetOnShelvesGoods(g._mchId,
		start, end, sortBy)
}
