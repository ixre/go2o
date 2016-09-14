/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:53
 * description :
 * history :
 */

package sale

import (
	"errors"
	"fmt"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"strconv"
	"time"
)

var _ sale.IItem = new(itemImpl)

type itemImpl struct {
	_manager      *itemManagerImpl
	_value        *item.Item
	_saleRep      sale.ISaleRep
	_itemRep      item.IItemRep
	_saleLabelRep sale.ISaleLabelRep
	_goodsRep     goods.IGoodsRep
	_expressRep   express.IExpressRep
	_promRep      promotion.IPromotionRep
	_sale         *saleImpl
	_saleLabels   []*sale.Label
	_valueRep     valueobject.IValueRep
}

func newItemImpl(mgr *itemManagerImpl, sale *saleImpl, v *item.Item,
	itemRep item.IItemRep, saleRep sale.ISaleRep,
	saleLabelRep sale.ISaleLabelRep, goodsRep goods.IGoodsRep,
	valRep valueobject.IValueRep, expressRep express.IExpressRep,
	promRep promotion.IPromotionRep) sale.IItem {
	return &itemImpl{
		_manager:      mgr,
		_value:        v,
		_itemRep:      itemRep,
		_saleRep:      saleRep,
		_saleLabelRep: saleLabelRep,
		_sale:         sale,
		_expressRep:   expressRep,
		_goodsRep:     goodsRep,
		_valueRep:     valRep,
	}
}

func (i *itemImpl) GetDomainId() int {
	return i._value.Id
}

func (i *itemImpl) GetValue() item.Item {
	return *i._value
}

func (i *itemImpl) checkValue(v *item.Item) error {
	registry := i._valueRep.GetRegistry()

	// 检测供应商
	if v.VendorId <= 0 || v.VendorId != i._value.VendorId {
		return item.ErrVendor
	}

	// 检测是否上传图片
	if v.Image == registry.GoodsDefaultImage {
		return item.ErrNotUploadImage
	}

	// 检测运费模板
	if v.ExpressTplId <= 0 {
		return shipment.ErrNotSetExpressTemplate
	}
	tpl := i._expressRep.GetUserExpress(v.VendorId).GetTemplate(v.ExpressTplId)
	if tpl == nil {
		return express.ErrNoSuchTemplate
	}
	if !tpl.Enabled() {
		return express.ErrTemplateNotEnabled
	}

	// 检测价格
	return i.checkPrice(v)
}

// 设置值
func (i *itemImpl) SetValue(v *item.Item) error {
	if err := i.checkValue(v); err != nil {
		return err
	}
	if v.Id == i._value.Id {
		v.CreateTime = i._value.CreateTime
		v.GoodsNo = i._value.GoodsNo
		i._value = v
	}
	i._value.UpdateTime = time.Now().Unix()
	return nil
}

// 是否上架
func (i *itemImpl) IsOnShelves() bool {
	return i._value.ShelveState == item.ShelvesOn
}

// 获取商品的销售标签
func (i *itemImpl) GetSaleLabels() []*sale.Label {
	if i._saleLabels == nil {
		i._saleLabels = i._saleLabelRep.GetItemSaleLabels(i.GetDomainId())
	}
	return i._saleLabels
}

// 保存销售标签
func (i *itemImpl) SaveSaleLabels(tagIds []int) error {
	err := i._saleLabelRep.CleanItemSaleLabels(i.GetDomainId())
	if err == nil {
		err = i._saleLabelRep.SaveItemSaleLabels(i.GetDomainId(), tagIds)
		i._saleLabels = nil
	}
	return err
}

// 重置审核状态
func (i *itemImpl) resetReview() {
	i._value.ReviewState = enum.ReviewAwaiting
}

// 判断价格是否正确
func (i *itemImpl) checkPrice(v *item.Item) error {
	rate := (v.SalePrice - v.Cost) / v.SalePrice
	if rate <= 0 {
		return goods.ErrSalePriceLessThanCost
	}
	conf := i._valueRep.GetRegistry()
	minRate := conf.GoodsMinProfitRate
	if rate < minRate {
		return errors.New(fmt.Sprintf(goods.ErrGoodsMinProfitRate.Error(),
			strconv.Itoa(int(minRate*100))+"%"))
	}
	return nil
}

// 设置上架
func (i *itemImpl) SetShelve(state int, remark string) error {
	if state == item.ShelvesIncorrect && len(remark) == 0 {
		return item.ErrNilRejectRemark
	}
	i._value.ShelveState = state
	i._value.Remark = remark
	_, err := i.Save()
	return err
}

// 保存
func (i *itemImpl) Save() (int, error) {
	i._sale.clearCache(i._value.Id)
	unix := time.Now().Unix()
	i._value.UpdateTime = unix
	if i.GetDomainId() <= 0 {
		i._value.CreateTime = unix
	}
	if i._value.GoodsNo == "" {
		cs := strconv.Itoa(i._value.CategoryId)
		us := strconv.Itoa(int(unix))
		l := len(cs)
		i._value.GoodsNo = fmt.Sprintf("%s%s", cs, us[4+l:])
	}
	i.resetReview()

	id, err := i._itemRep.SaveValueItem(i._value)
	if err == nil {
		i._value.Id = id

		//todo: 保存商品
		i.saveGoods()

		// 创建快照
		//_, err = i.GenerateSnapshot()
	}
	return id, err
}

//todo: 过渡方法,应有SKU,不根据Item生成Goods
func (i *itemImpl) saveGoods() {
	val := i._goodsRep.GetValueGoods(i.GetDomainId(), 0)
	if val == nil {
		val = &goods.ValueGoods{
			Id:            0,
			ItemId:        i.GetDomainId(),
			IsPresent:     0,
			SkuId:         0,
			PromotionFlag: 0,
			StockNum:      100,
			SaleNum:       100,
		}
	}
	goods := NewSaleGoods(nil, i._sale, i._itemRep, i, val,
		i._saleRep, i._goodsRep, i._promRep)
	goods.Save()
}

//// 生成快照
//func (i *Goods) GenerateSnapshot() (int, error) {
//	v := i._value
//	if v.Id <= 0 {
//		return 0, sale.ErrNoSuchGoods
//	}
//
//	if v.OnShelves == 0 {
//		return 0, sale.ErrNotOnShelves
//	}
//
//	merchantId := i._sale.GetAggregateRootId()
//	unix := time.Now().Unix()
//	cate := i._saleRep.GetCategory(merchantId, v.CategoryId)
//	var gsn *goods.GoodsSnapshot = &goods.GoodsSnapshot{
//		Key:          fmt.Sprintf("%d-g%d-%d", merchantId, v.Id, unix),
//		GoodsId:      i.GetDomainId(),
//		GoodsName:    v.Name,
//		GoodsNo:      v.GoodsNo,
//		SmallTitle:   v.SmallTitle,
//		CategoryName: cate.Name,
//		Image:        v.Image,
//		Cost:         v.Cost,
//		Price:        v.Price,
//		SalePrice:    v.SalePrice,
//		CreateTime:   unix,
//	}
//
//	if i.isNewSnapshot(gsn) {
//		i._latestSnapshot = gsn
//		return i._saleRep.SaveSnapshot(gsn)
//	}
//	return 0, sale.ErrLatestSnapshot
//}
//
//// 是否为新快照,与旧有快照进行数据对比
//func (i *Goods) isNewSnapshot(gsn *goods.GoodsSnapshot) bool {
//	latestGsn := i.GetLatestSnapshot()
//	if latestGsn != nil {
//		return latestGsn.GoodsName != gsn.GoodsName ||
//			latestGsn.SmallTitle != gsn.SmallTitle ||
//			latestGsn.CategoryName != gsn.CategoryName ||
//			latestGsn.Image != gsn.Image ||
//			latestGsn.Cost != gsn.Cost ||
//			latestGsn.Price != gsn.Price ||
//			latestGsn.SalePrice != gsn.SalePrice
//	}
//	return true
//}
//
//// 获取最新的快照
//func (i *Goods) GetLatestSnapshot() *goods.GoodsSnapshot {
//	if i._latestSnapshot == nil {
//		i._latestSnapshot = i._saleRep.GetLatestGoodsSnapshot(i.GetDomainId())
//	}
//	return i._latestSnapshot
//}

var _ sale.IItemManager = new(itemManagerImpl)

type itemManagerImpl struct {
	_sale       *saleImpl
	_itemRep    item.IItemRep
	_valRep     valueobject.IValueRep
	_expressRep express.IExpressRep
	_vendorId   int
}

func NewItemManager(vendorId int, s *saleImpl,
	itemRep item.IItemRep, expressRep express.IExpressRep,
	valRep valueobject.IValueRep) sale.IItemManager {
	c := &itemManagerImpl{
		_sale:       s,
		_vendorId:   vendorId,
		_valRep:     valRep,
		_itemRep:    itemRep,
		_expressRep: expressRep,
	}
	return c.init()
}

func (i *itemManagerImpl) init() sale.IItemManager {
	return i
}

func (i *itemManagerImpl) CreateItem(v *item.Item) sale.IItem {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}
	if v.UpdateTime == 0 {
		v.UpdateTime = v.CreateTime
	} //todo: 判断category
	return newItemImpl(i, i._sale, v, i._itemRep,
		i._sale._saleRep, i._sale._labelRep,
		i._sale._goodsRep, i._valRep, i._expressRep,
		i._sale._promRep)
}

// 删除货品
func (i *itemManagerImpl) DeleteItem(id int) error {
	var err error
	num := i._itemRep.GetItemSaleNum(i._vendorId, id)

	if num == 0 {
		err = i._itemRep.DeleteItem(i._vendorId, id)
		if err != nil {
			i._sale.clearCache(id)
		}
	} else {
		err = sale.ErrCanNotDeleteItem
	}
	return err
}

// 根据产品编号获取产品
func (i *itemManagerImpl) GetItem(itemId int) sale.IItem {
	pv := i._itemRep.GetValueItem(itemId)
	if pv != nil && pv.VendorId == i._vendorId {
		return i.CreateItem(pv)
	}
	return nil
}
