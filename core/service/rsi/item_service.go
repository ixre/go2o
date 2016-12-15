/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package rsi

import (
	"fmt"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/core/service/thrift/idl/gen-go/define"
	"go2o/core/service/thrift/parser"
)

type itemService struct {
	itemRepo    item.IGoodsItemRepo
	_goodsQuery *query.GoodsQuery
	_cateRepo   product.ICategoryRepo
	labelRepo   item.ISaleLabelRepo
}

func NewSaleService(cateRepo product.ICategoryRepo,
	goodsRepo item.IGoodsItemRepo, goodsQuery *query.GoodsQuery,
	labelRepo item.ISaleLabelRepo) *itemService {
	return &itemService{
		itemRepo:    goodsRepo,
		_goodsQuery: goodsQuery,
		_cateRepo:   cateRepo,
		labelRepo:   labelRepo,
	}
}

// 获取商品值
func (s *itemService) GetItemValue(itemId int32) *item.GoodsItem {
	item := s.itemRepo.GetItem(itemId)
	if item != nil {
		return item.GetValue()
	}
	return nil
}

// 根据SKU获取商品
func (s *itemService) GetGoodsBySku(mchId int32, itemId int32, sku int32) *valueobject.Goods {
	v := s.itemRepo.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		return s.itemRepo.CreateItem(v).GetPackedValue()
	}
	return nil
}

// 根据SKU获取商品
func (s *itemService) GetValueGoodsBySku(mchId int32, itemId int32, sku int32) *item.GoodsItem {
	v := s.itemRepo.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		return s.itemRepo.CreateItem(v).GetValue()
	}
	return nil
}

// 根据快照编号获取商品
func (s *itemService) GetGoodsBySnapshotId(snapshotId int32) *item.GoodsItem {
	snap := s.itemRepo.GetSaleSnapshot(snapshotId)
	if snap != nil {
		return s.itemRepo.GetValueGoodsById(snap.SkuId)
	}
	return nil
}

// 根据快照编号获取商品
func (s *itemService) GetSaleSnapshotById(snapshotId int32) *item.SalesSnapshot {
	return s.itemRepo.GetSaleSnapshot(snapshotId)
}

// 保存商品
func (s *itemService) SaveItem(gs *item.GoodsItem, vendorId int32) (_ *define.Result_, err error) {
	var gi item.IGoodsItem
	if gs.Id > 0 {
		gi = s.itemRepo.GetItem(gs.Id)
		if gi == nil || gi.GetValue().VendorId != vendorId {
			err = item.ErrNoSuchGoods
			goto R
		}
	} else {
		gi = s.itemRepo.CreateItem(gs)
	}
	err = gi.SetValue(gs)
	if err == nil {
		if gs.SkuArray != nil {
			//err = gi.SetSkus(gs.SkuArray)
		}
		gs.Id, err = gi.Save()
	}
R:
	return parser.Result(gs.Id, err), nil
}

// 获取分页上架的商品
func (s *itemService) GetShopPagedOnShelvesGoods(shopId, categoryId int32, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cat := s._cateRepo.GlobCatService().GetCategory(categoryId)
		ids := cat.GetChildes()
		ids = append(ids, categoryId)
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId, ids, start, end, "", sortBy)
	} else {
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId, nil, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取上架商品数据（分页）
func (s *itemService) GetPagedOnShelvesGoods(shopId int32, categoryId int32, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cate := s._cateRepo.GlobCatService().GetCategory(categoryId)
		var ids []int32 = cate.GetChildes()
		ids = append(ids, categoryId)
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId,
			ids, start, end, "", sortBy)
	} else {
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId,
			[]int32{}, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取分页上架的商品
func (s *itemService) GetPagedOnShelvesGoodsByKeyword(shopId int32, start, end int,
	word, sortQuery string) (int, []*valueobject.Goods) {
	var where string
	var orderBy string
	switch sortQuery {
	case "price_0":
		where = ""
		orderBy = "pro_product.sale_price ASC"
	case "price_1":
		where = ""
		orderBy = "pro_product.sale_price DESC"
	case "sale_0":
		where = ""
		orderBy = "item_info.sale_num ASC"
	case "sale_1":
		where = ""
		orderBy = "item_info.sale_num DESC"
	case "rate_0":
	//todo:
	case "rate_1":
		//todo:
	}
	return s._goodsQuery.GetPagedOnShelvesGoodsByKeyword(shopId,
		start, end, word, where, orderBy)
}

// 删除产品
func (s *itemService) DeleteGoods(mchId, goodsId int32) error {
	gi := s.itemRepo.GetItem(goodsId)
	if gi == nil || gi.GetValue().VendorId != mchId {
		return item.ErrNoSuchGoods
	}
	return gi.Destroy()
}

func (s *itemService) GetAllSaleLabels() []*item.Label {
	tags := s.labelRepo.LabelService().GetAllSaleLabels()
	lbs := make([]*item.Label, len(tags))
	for i, v := range tags {
		lbs[i] = v.GetValue()
	}
	return lbs
}

// 获取销售标签
func (s *itemService) GetSaleLabel(labelId int32) *item.Label {
	tag := s.labelRepo.LabelService().GetSaleLabel(labelId)
	if tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 获取销售标签
func (s *itemService) GetSaleLabelByCode(code string) *item.Label {
	tag := s.labelRepo.LabelService().GetSaleLabelByCode(code)
	if tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 保存销售标签
func (s *itemService) SaveSaleLabel(v *item.Label) (int32, error) {
	ls := s.labelRepo.LabelService()
	var lb item.ISaleLabel

	if v.Id > 0 {
		lb = ls.GetSaleLabel(v.Id)
		if lb == nil {
			panic("没有销售标签")
		}
	} else {
		lb = ls.CreateSaleLabel(v)
	}
	err := lb.SetValue(v)
	if err == nil {
		return lb.Save()
	}
	return v.Id, err
}

// 根据销售标签获取指定数目的商品
func (s *itemService) GetValueGoodsBySaleLabel(code, sortBy string, begin int, end int) []*valueobject.Goods {
	tag := s.labelRepo.LabelService().GetSaleLabelByCode(code)
	if tag != nil {
		list := tag.GetValueGoods(sortBy, begin, end)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
		}
		return list
	}
	return make([]*valueobject.Goods, 0)
}

// 根据分页销售标签获取指定数目的商品
func (s *itemService) GetPagedValueGoodsBySaleLabel(shopId int32, tagId int32, sortBy string, begin int, end int) (int,
	[]*valueobject.Goods) {
	tag := s.labelRepo.LabelService().CreateSaleLabel(&item.Label{
		Id: tagId,
	})
	return tag.GetPagedValueGoods(sortBy, begin, end)
}

// 删除销售标签
func (s *itemService) DeleteSaleLabel(id int32) error {
	return s.labelRepo.LabelService().DeleteSaleLabel(id)
}

// 获取商品的会员价
func (s *itemService) GetGoodsLevelPrices(itemId int32) []*item.MemberPrice {
	gi := s.itemRepo.GetItem(itemId)
	if gi != nil {
		return gi.GetLevelPrices()
	}
	return make([]*item.MemberPrice, 0)
}

// 保存商品的会员价
func (s *itemService) SaveMemberPrices(mchId int32, itemId int32,
	priceSet []*item.MemberPrice) (err error) {
	gi := s.itemRepo.GetItem(itemId)
	if gi != nil {
		for _, v := range priceSet {
			if _, err = gi.SaveLevelPrice(v); err != nil {
				return err
			}
		}
	}
	return err
}

//func (s *saleService) GetGoodsComplexInfo(goodsId int32) *dto.GoodsComplex {
//	return s._goodsQuery.GetGoodsComplex(goodsId)
//}

// 获取商品详情
func (s *itemService) GetGoodsDetails(itemId, mLevel int32) (
	*valueobject.Goods, map[string]string) {
	goods := s.itemRepo.GetItem(itemId)
	gv := goods.GetPackedValue()
	proMap := goods.GetPromotionDescribe()
	if b, price := goods.GetLevelPrice(mLevel); b {
		gv.PromPrice = price
		proMap["会员专享"] = fmt.Sprintf("会员优惠,仅需<b>￥%s</b>",
			format.FormatFloat(price))
	}
	return gv, proMap
}

// 获取货品描述
func (s *itemService) GetItemDescriptionByGoodsId(itemId int32) string {
	goods := s.itemRepo.GetItem(itemId)
	return goods.GetItem().GetValue().Description
}

// 获取商品快照
func (s *itemService) GetSnapshot(skuId int32) *item.Snapshot {
	return s.itemRepo.GetLatestSnapshot(skuId)
}

// 设置商品货架状态
func (s *itemService) SetShelveState(vendorId int32, itemId int32,
	state int32, remark string) (_ *define.Result_, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchGoods
	} else {
		err = it.SetShelve(state, remark)
	}
	return parser.Result(0, err), nil
}

// 设置商品货架状态
func (s *itemService) ReviewItem(vendorId int32, itemId int32,
	pass bool, remark string) (_ *define.Result_, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchGoods
	} else {
		err = it.Review(pass, remark)
	}
	return parser.Result(0, err), nil
}

// 标记为违规
func (s *itemService) SignIncorrect(vendorId int32, itemId int32,
	remark string) (_ *define.Result_, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchGoods
	} else {
		err = it.Incorrect(remark)
	}
	return parser.Result(0, err), nil
}
