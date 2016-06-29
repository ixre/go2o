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
	_manager        *goodsManagerImpl
	_goods          sale.IItem
	_value          *goods.ValueGoods
	_saleRep        sale.ISaleRep
	_goodsRep       goods.IGoodsRep
	_itemRep        item.IItemRep
	_promRep        promotion.IPromotionRep
	_sale           sale.ISale
	_latestSnapshot *goods.GoodsSnapshot
	_levelPrices    []*goods.MemberPrice
	_promDescribes  map[string]string
	_snapManager    goods.ISnapshotManager
}

func NewSaleGoods(m *goodsManagerImpl, s sale.ISale, goods sale.IItem, value *goods.ValueGoods, rep sale.ISaleRep,
	goodsRep goods.IGoodsRep, promRep promotion.IPromotionRep) sale.IGoods {
	v := &tmpGoodsImpl{
		_manager:        m,
		_goods:          goods,
		_value:          value,
		_saleRep:        rep,
		_goodsRep:       goodsRep,
		_promRep:        promRep,
		_sale:           s,
		_latestSnapshot: nil,
	}
	return v.init()
}

func (this *tmpGoodsImpl) init() sale.IGoods {
	this._value.Price = this._value.Price
	this._value.SalePrice = this._goods.GetValue().SalePrice
	this._value.PromPrice = this._goods.GetValue().SalePrice
	return this
}

//获取领域对象编号
func (this *tmpGoodsImpl) GetDomainId() int {
	return this._value.Id
}

// 商品快照
func (this *tmpGoodsImpl) SnapshotManager() goods.ISnapshotManager {
	if this._snapManager == nil {
		gi := this.GetItem().GetValue()
		this._snapManager = goodsImpl.NewSnapshotManagerImpl(this.GetDomainId(),
			this._goodsRep, this._itemRep, this.GetValue(), &gi)
	}
	return this._snapManager
}

// 获取货品
func (this *tmpGoodsImpl) GetItem() sale.IItem {
	return this._goods
}

// 设置值
func (this *tmpGoodsImpl) GetValue() *goods.ValueGoods {
	return this._value
}

// 获取包装过的商品信息
func (this *tmpGoodsImpl) GetPackedValue() *valueobject.Goods {
	item := this.GetItem().GetValue()
	gv := this.GetValue()
	goods := &valueobject.Goods{
		Item_Id:       item.Id,
		CategoryId:    item.CategoryId,
		Name:          item.Name,
		GoodsNo:       item.GoodsNo,
		Image:         item.Image,
		Price:         item.Price,
		SalePrice:     item.SalePrice,
		PromPrice:     item.SalePrice,
		GoodsId:       this.GetDomainId(),
		SkuId:         gv.SkuId,
		IsPresent:     gv.IsPresent,
		PromotionFlag: gv.PromotionFlag,
		StockNum:      gv.StockNum,
		SaleNum:       gv.SaleNum,
	}
	return goods
}

// 获取促销信息
func (this *tmpGoodsImpl) GetPromotions() []promotion.IPromotion {
	var vp []*promotion.PromotionInfo = this._promRep.GetPromotionOfGoods(this.GetDomainId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = this._promRep.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (this *tmpGoodsImpl) GetLevelPrice(level int) (bool, float32) {
	lvp := this.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < this._value.SalePrice {
			return true, v.Price
		}
	}
	return false, this._value.SalePrice
}

// 获取促销价
func (this *tmpGoodsImpl) GetPromotionPrice(level int) float32 {
	b, price := this.GetLevelPrice(level)
	if b {
		return price
	}
	return this._value.SalePrice
}

// 获取促销描述
func (this *tmpGoodsImpl) GetPromotionDescribe() map[string]string {
	if this._promDescribes == nil {
		proms := this.GetPromotions()
		this._promDescribes = make(map[string]string, len(proms))
		for _, v := range proms {
			key := v.TypeName()
			if txt, ok := this._promDescribes[key]; !ok {
				this._promDescribes[key] = v.GetValue().ShortName
			} else {
				this._promDescribes[key] = txt + "；" + v.GetValue().ShortName
			}

			//			if v.Type() == promotion.TypeFlagCashBack {
			//				if txt, ok := this._promDescribes[key]; !ok {
			//					this._promDescribes[key] = v.GetValue().ShortName
			//				} else {
			//					this._promDescribes[key] = txt + "；" + v.GetValue().ShortName
			//				}
			//			} else if v.Type() == promotion.TypeFlagCoupon {
			//				if txt, ok := this._promDescribes[key]; !ok {
			//					this._promDescribes[key] = v.GetValue().ShortName
			//				} else {
			//					this._promDescribes[key] = txt + "；" + v.GetValue().ShortName
			//				}
			//			}

			//todo: other promotion implement
		}
	}
	return this._promDescribes
}

// 获取会员价
func (this *tmpGoodsImpl) GetLevelPrices() []*goods.MemberPrice {
	if this._levelPrices == nil {
		this._levelPrices = this._goodsRep.GetGoodsLevelPrice(this.GetDomainId())
	}
	return this._levelPrices
}

// 保存会员价
func (this *tmpGoodsImpl) SaveLevelPrice(v *goods.MemberPrice) (int, error) {
	v.GoodsId = this.GetDomainId()
	if this._value.SalePrice == v.Price {
		if v.Id > 0 {
			this._goodsRep.RemoveGoodsLevelPrice(v.Id)
		}
		return -1, nil
	}
	return this._goodsRep.SaveGoodsLevelPrice(v)
}

// 设置值
func (this *tmpGoodsImpl) SetValue(v *goods.ValueGoods) error {
	this._value.IsPresent = v.IsPresent
	this._value.SaleNum = v.SaleNum
	this._value.StockNum = v.StockNum
	this._value.SkuId = v.SkuId
	//this._value.PromotionFlag = v.PromotionFlag
	return nil
}

// 保存
func (this *tmpGoodsImpl) Save() (int, error) {
	id, err := this._goodsRep.SaveValueGoods(this._value)
	if err == nil {
		_, err = this.SnapshotManager().GenerateSnapshot()
	}
	this._value.Id = id
	return id, err
	//todo: save promotion
	// return id,err
}

// 更新销售数量
func (this *tmpGoodsImpl) AddSaleNum(quantity int) error {
	// 减去库存
	if quantity > 0 {
		if quantity > this._value.StockNum {
			return sale.ErrOutOfStock
		}
		this._value.StockNum -= quantity
		this._value.SaleNum += quantity
		_, err := this.Save()
		return err
	}
	return sale.ErrGoodsNum
}

// 取消销售
func (this *tmpGoodsImpl) CancelSale(quantity int, orderNo string) error {
	// 减去库存
	if quantity > 0 {
		this._value.StockNum += quantity
		this._value.SaleNum -= quantity
		_, err := this.Save()
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

func (this *goodsManagerImpl) init() sale.IGoodsManager {
	return this
}

// 创建商品
func (this *goodsManagerImpl) CreateGoods(s *goods.ValueGoods) sale.IGoods {
	return NewSaleGoods(this, this._sale, nil, s, this._sale._saleRep,
		this._sale._goodsRep, this._sale._promRep)
}

// 创建商品
func (this *goodsManagerImpl) CreateGoodsByItem(item sale.IItem, v *goods.ValueGoods) sale.IGoods {
	return NewSaleGoods(this, this._sale, item, v, this._sale._saleRep,
		this._sale._goodsRep, this._sale._promRep)
}

// 根据产品编号获取商品
func (this *goodsManagerImpl) GetGoods(goodsId int) sale.IGoods {
	var v *goods.ValueGoods = this._sale._goodsRep.GetValueGoodsById(goodsId)
	if v != nil {
		pv := this._sale._itemRep.GetValueItem(v.ItemId)
		if pv != nil {
			return this.CreateGoodsByItem(this._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 根据产品SKU获取商品
func (this *goodsManagerImpl) GetGoodsBySku(itemId, sku int) sale.IGoods {
	var v *goods.ValueGoods = this._sale._goodsRep.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		pv := this._sale._itemRep.GetValueItem(v.ItemId)
		if pv != nil {
			return this.CreateGoodsByItem(this._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 删除商品
func (this *goodsManagerImpl) DeleteGoods(goodsId int) error {
	gs := this.GetGoods(goodsId)
	if gs.GetValue().SaleNum > 0 {
		return goods.ErrNoSuchSnapshot
	}

	//todo: delete goods
	err := this._sale._itemRep.DeleteItem(this._mchId, goodsId)
	if err != nil {
		this._sale.clearCache(goodsId)
	}
	return err
}

// 获取指定数量已上架的商品
func (this *goodsManagerImpl) GetOnShelvesGoods(start, end int,
	sortBy string) []*valueobject.Goods {
	return this._sale._goodsRep.GetOnShelvesGoods(this._mchId,
		start, end, sortBy)
}
