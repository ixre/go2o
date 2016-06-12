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
	"fmt"
	"go2o/core/domain"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	"time"
)

var _ sale.IGoods = new(SaleGoods)
var _ domain.IDomain = new(SaleGoods)

type SaleGoods struct {
	_manager        *goodsManagerImpl
	_goods          sale.IItem
	_value          *sale.ValueGoods
	_saleRep        sale.ISaleRep
	_goodsRep       sale.IGoodsRep
	_promRep        promotion.IPromotionRep
	_sale           sale.ISale
	_latestSnapshot *sale.GoodsSnapshot
	_levelPrices    []*sale.MemberPrice
	_promDescribes  map[string]string
}

func NewSaleGoods(m *goodsManagerImpl, s sale.ISale, goods sale.IItem, value *sale.ValueGoods, rep sale.ISaleRep,
	goodsRep sale.IGoodsRep, promRep promotion.IPromotionRep) sale.IGoods {
	v := &SaleGoods{
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

func (this *SaleGoods) init() sale.IGoods {
	this._value.Price = this._value.Price
	this._value.SalePrice = this._goods.GetValue().SalePrice
	this._value.PromPrice = this._goods.GetValue().SalePrice
	return this
}

//获取领域对象编号
func (this *SaleGoods) GetDomainId() int {
	return this._value.Id
}

// 获取货品
func (this *SaleGoods) GetItem() sale.IItem {
	return this._goods
}

// 设置值
func (this *SaleGoods) GetValue() *sale.ValueGoods {
	return this._value
}

// 获取包装过的商品信息
func (this *SaleGoods) GetPackedValue() *valueobject.Goods {
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
func (this *SaleGoods) GetPromotions() []promotion.IPromotion {
	var vp []*promotion.PromotionInfo = this._promRep.GetPromotionOfGoods(this.GetDomainId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = this._promRep.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (this *SaleGoods) GetLevelPrice(level int) (bool, float32) {
	lvp := this.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < this._value.SalePrice {
			return true, v.Price
		}
	}
	return false, this._value.SalePrice
}

// 获取促销价
func (this *SaleGoods) GetPromotionPrice(level int) float32 {
	b, price := this.GetLevelPrice(level)
	if b {
		return price
	}
	return this._value.SalePrice
}

// 获取促销描述
func (this *SaleGoods) GetPromotionDescribe() map[string]string {
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
func (this *SaleGoods) GetLevelPrices() []*sale.MemberPrice {
	if this._levelPrices == nil {
		this._levelPrices = this._goodsRep.GetGoodsLevelPrice(this.GetDomainId())
	}
	return this._levelPrices
}

// 保存会员价
func (this *SaleGoods) SaveLevelPrice(v *sale.MemberPrice) (int, error) {
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
func (this *SaleGoods) SetValue(v *sale.ValueGoods) error {
	this._value.IsPresent = v.IsPresent
	this._value.SaleNum = v.SaleNum
	this._value.StockNum = v.StockNum
	this._value.SkuId = v.SkuId
	//this._value.PromotionFlag = v.PromotionFlag
	return nil
}

// 保存
func (this *SaleGoods) Save() (int, error) {
	id, err := this._goodsRep.SaveValueGoods(this._value)
	if err == nil {
		_, err = this.GenerateSnapshot()
	}
	this._value.Id = id
	return id, err

	//todo: save promotion
	// return id,err
}

// 生成快照
func (this *SaleGoods) GenerateSnapshot() (int, error) {
	v := this._value
	gi := this.GetItem()
	gv := gi.GetValue()

	if v.Id <= 0 {
		return -1, sale.ErrNoSuchGoods
	}

	if !gi.IsOnShelves() {
		return -1, sale.ErrNotOnShelves
	}

	merchantId := this._sale.GetAggregateRootId()
	unix := time.Now().Unix()
	cate := this._sale.CategoryManager().GetCategory(gv.CategoryId)
	var gsn *sale.GoodsSnapshot = &sale.GoodsSnapshot{
		Key:          fmt.Sprintf("%d-g%d-%d", merchantId, v.Id, unix),
		ItemId:       gv.Id,
		GoodsId:      this.GetDomainId(),
		GoodsName:    gv.Name,
		GoodsNo:      gv.GoodsNo,
		SmallTitle:   gv.SmallTitle,
		CategoryName: cate.GetValue().Name,
		Image:        gv.Image,
		Cost:         gv.Cost,
		SalePrice:    gv.SalePrice,
		Price:        this._value.Price,
		CreateTime:   unix,
	}

	if this.isNewSnapshot(gsn) {
		this._latestSnapshot = gsn
		return this._saleRep.SaveSnapshot(gsn)
	}

	return 0, sale.ErrLatestSnapshot
}

// 更新销售数量
func (this *SaleGoods) AddSaleNum(quantity int) error {
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
func (this *SaleGoods) CancelSale(quantity int, orderNo string) error {
	// 减去库存
	if quantity > 0 {
		this._value.StockNum += quantity
		this._value.SaleNum -= quantity
		_, err := this.Save()
		return err
	}
	return sale.ErrGoodsNum
}

// 是否为新快照,与旧有快照进行数据对比
func (this *SaleGoods) isNewSnapshot(gsn *sale.GoodsSnapshot) bool {
	latestGsn := this.GetLatestSnapshot()
	if latestGsn != nil {
		return latestGsn.GoodsName != gsn.GoodsName ||
			latestGsn.SmallTitle != gsn.SmallTitle ||
			latestGsn.CategoryName != gsn.CategoryName ||
			latestGsn.Image != gsn.Image ||
			latestGsn.Cost != gsn.Cost ||
			latestGsn.Price != gsn.Price ||
			latestGsn.SalePrice != gsn.SalePrice
	}
	return true
}

// 获取最新的快照
func (this *SaleGoods) GetLatestSnapshot() *sale.GoodsSnapshot {
	if this._latestSnapshot == nil {
		this._latestSnapshot = this._saleRep.GetLatestGoodsSnapshot(this.GetDomainId())
	}
	return this._latestSnapshot
}

var _ sale.IGoodsManager = new(goodsManagerImpl)

type goodsManagerImpl struct {
	_sale   *SaleImpl
	_valRep valueobject.IValueRep
	_mchId  int
}

func NewGoodsManager(mchId int, s *SaleImpl,
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
func (this *goodsManagerImpl) CreateGoods(s *sale.ValueGoods) sale.IGoods {
	return NewSaleGoods(this, this._sale, nil, s, this._sale._saleRep,
		this._sale._goodsRep, this._sale._promRep)
}

// 创建商品
func (this *goodsManagerImpl) CreateGoodsByItem(item sale.IItem, v *sale.ValueGoods) sale.IGoods {
	return NewSaleGoods(this, this._sale, item, v, this._sale._saleRep,
		this._sale._goodsRep, this._sale._promRep)
}

// 根据产品编号获取商品
func (this *goodsManagerImpl) GetGoods(goodsId int) sale.IGoods {
	var v *sale.ValueGoods = this._sale._goodsRep.GetValueGoodsById(goodsId)
	if v != nil {
		pv := this._sale._saleRep.GetValueItem(this._mchId, v.ItemId)
		if pv != nil {
			return this.CreateGoodsByItem(this._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 根据产品SKU获取商品
func (this *goodsManagerImpl) GetGoodsBySku(itemId, sku int) sale.IGoods {
	var v *sale.ValueGoods = this._sale._goodsRep.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		pv := this._sale._saleRep.GetValueItem(this._mchId, v.ItemId)
		if pv != nil {
			return this.CreateGoodsByItem(this._sale.ItemManager().CreateItem(pv), v)
		}
	}
	return nil
}

// 删除商品
func (this *goodsManagerImpl) DeleteGoods(goodsId int) error {
	goods := this.GetGoods(goodsId)
	if goods.GetValue().SaleNum > 0 {
		return sale.ErrNoSuchSnapshot
	}

	//todo: delete goods
	err := this._sale._saleRep.DeleteItem(this._mchId, goodsId)
	if err != nil {
		this._sale.clearCache(goodsId)
	}
	return err
}

// 获取指定的商品快照
func (this *goodsManagerImpl) GetGoodsSnapshot(id int) *sale.GoodsSnapshot {
	return this._sale._saleRep.GetGoodsSnapshot(id)
}

// 根据Key获取商品快照
func (this *goodsManagerImpl) GetGoodsSnapshotByKey(key string) *sale.GoodsSnapshot {
	return this._sale._saleRep.GetGoodsSnapshotByKey(key)
}

// 获取指定数量已上架的商品
func (this *goodsManagerImpl) GetOnShelvesGoods(start, end int,
	sortBy string) []*valueobject.Goods {
	return this._sale._goodsRep.GetOnShelvesGoods(this._mchId,
		start, end, sortBy)
}
