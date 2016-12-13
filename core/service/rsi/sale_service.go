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
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
)

type saleService struct {
	_rep        sale.ISaleRepo
	_goodsRepo  item.IGoodsRepo
	_goodsQuery *query.GoodsQuery
	_cateRepo   product.ICategoryRepo
}

func NewSaleService(r sale.ISaleRepo, cateRepo product.ICategoryRepo,
	goodsRepo item.IGoodsRepo, goodsQuery *query.GoodsQuery) *saleService {
	return &saleService{
		_rep:        r,
		_goodsRepo:  goodsRepo,
		_goodsQuery: goodsQuery,
		_cateRepo:   cateRepo,
	}
}

// 获取商品值
func (s *saleService) GetValueGoods(mchId, goodsId int32) *valueobject.Goods {
	sl := s._rep.GetSale(mchId)
	var goods item.IGoods = sl.GoodsManager().GetGoods(goodsId)
	if goods != nil {
		return goods.GetPackedValue()
	}
	return nil
}

// 根据SKU获取商品
func (s *saleService) GetGoodsBySku(mchId int32, itemId int32, sku int32) *valueobject.Goods {
	sl := s._rep.GetSale(mchId)
	var goods item.IGoods = sl.GoodsManager().GetGoodsBySku(itemId, sku)
	return goods.GetPackedValue()
}

// 根据SKU获取商品
func (s *saleService) GetValueGoodsBySku(mchId int32, itemId int32, sku int32) *item.GoodsItem {
	sl := s._rep.GetSale(mchId)
	gs := sl.GoodsManager().GetGoodsBySku(itemId, sku)
	if gs != nil {
		return gs.GetValue()
	}
	return nil
}

// 根据快照编号获取商品
func (s *saleService) GetGoodsBySnapshotId(snapshotId int32) *item.GoodsItem {
	snap := s._goodsRepo.GetSaleSnapshot(snapshotId)
	if snap != nil {
		return s._goodsRepo.GetValueGoodsById(snap.SkuId)
	}
	return nil
}

// 根据快照编号获取商品
func (s *saleService) GetSaleSnapshotById(snapshotId int32) *item.SalesSnapshot {
	return s._goodsRepo.GetSaleSnapshot(snapshotId)
}

// 保存商品
func (s *saleService) SaveGoods(mchId int32, gs *item.GoodsItem) (int32, error) {
	sl := s._rep.GetSale(mchId)
	if gs.Id > 0 {
		g := sl.GoodsManager().GetGoods(gs.Id)
		g.SetValue(gs)
		return g.Save()
	}
	g := sl.GoodsManager().CreateGoods(gs)
	return g.Save()
}

// 获取分页上架的商品
func (s *saleService) GetShopPagedOnShelvesGoods(shopId, categoryId int32, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cat := s._cateRepo.GlobCatService().GetCategory(categoryId)
		ids := cat.GetChildes()
		ids = append(ids, categoryId)
		total, list = s._goodsRepo.GetPagedOnShelvesGoods(shopId, ids, start, end, "", sortBy)
	} else {
		total, list = s._goodsRepo.GetPagedOnShelvesGoods(shopId, nil, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取上架商品数据（分页）
func (s *saleService) GetPagedOnShelvesGoods(shopId int32, categoryId int32, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cate := s._cateRepo.GlobCatService().GetCategory(categoryId)
		var ids []int32 = cate.GetChildes()
		ids = append(ids, categoryId)
		total, list = s._goodsRepo.GetPagedOnShelvesGoods(shopId,
			ids, start, end, "", sortBy)
	} else {
		total, list = s._goodsRepo.GetPagedOnShelvesGoods(shopId,
			[]int32{}, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取分页上架的商品
func (s *saleService) GetPagedOnShelvesGoodsByKeyword(shopId int32, start, end int,
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
		orderBy = "gs_goods.sale_num ASC"
	case "sale_1":
		where = ""
		orderBy = "gs_goods.sale_num DESC"
	case "rate_0":
	//todo:
	case "rate_1":
		//todo:
	}
	return s._goodsQuery.GetPagedOnShelvesGoodsByKeyword(shopId,
		start, end, word, where, orderBy)
}

// 删除产品
func (s *saleService) DeleteGoods(mchId, goodsId int32) error {
	sl := s._rep.GetSale(mchId)
	return sl.GoodsManager().DeleteGoods(goodsId)
}

func (s *saleService) GetAllSaleLabels(mchId int32) []*item.Label {
	sl := s._rep.GetSale(mchId)
	tags := sl.LabelManager().GetAllSaleLabels()

	lbs := make([]*item.Label, len(tags))
	for i, v := range tags {
		lbs[i] = v.GetValue()
	}
	return lbs
}

// 获取销售标签
func (s *saleService) GetSaleLabel(mchId, id int32) *item.Label {
	sl := s._rep.GetSale(mchId)
	if tag := sl.LabelManager().GetSaleLabel(id); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 获取销售标签
func (s *saleService) GetSaleLabelByCode(mchId int32, code string) *item.Label {
	sl := s._rep.GetSale(mchId)
	if tag := sl.LabelManager().GetSaleLabelByCode(code); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 保存销售标签
func (s *saleService) SaveSaleLabel(mchId int32, v *item.Label) (int32, error) {
	sl := s._rep.GetSale(mchId)
	if v.Id > 0 {
		tag := sl.LabelManager().GetSaleLabel(v.Id)
		tag.SetValue(v)
		return tag.Save()
	}
	return sl.LabelManager().CreateSaleLabel(v).Save()
}

// 根据销售标签获取指定数目的商品
func (s *saleService) GetValueGoodsBySaleLabel(mchId int32,
	code, sortBy string, begin int, end int) []*valueobject.Goods {
	sl := s._rep.GetSale(mchId)
	if tag := sl.LabelManager().GetSaleLabelByCode(code); tag != nil {
		list := tag.GetValueGoods(sortBy, begin, end)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
		}
		return list
	}
	return make([]*valueobject.Goods, 0)
}

// 根据分页销售标签获取指定数目的商品
func (s *saleService) GetPagedValueGoodsBySaleLabel(mchId int32,
	tagId int32, sortBy string, begin int, end int) (int, []*valueobject.Goods) {
	sl := s._rep.GetSale(mchId)
	tag := sl.LabelManager().CreateSaleLabel(&item.Label{
		Id: tagId,
	})
	return tag.GetPagedValueGoods(sortBy, begin, end)
}

// 删除销售标签
func (s *saleService) DeleteSaleLabel(mchId int32, id int32) error {
	return s._rep.GetSale(mchId).LabelManager().DeleteSaleLabel(id)
}

// 获取商品的会员价
func (s *saleService) GetGoodsLevelPrices(mchId, goodsId int32) []*item.MemberPrice {
	sl := s._rep.GetSale(mchId)
	if goods := sl.GoodsManager().GetGoods(goodsId); goods != nil {
		return goods.GetLevelPrices()
	}
	return make([]*item.MemberPrice, 0)
}

// 保存商品的会员价
func (s *saleService) SaveMemberPrices(mchId int32, goodsId int32,
	priceSet []*item.MemberPrice) error {
	sl := s._rep.GetSale(mchId)
	var err error
	if goods := sl.GoodsManager().GetGoods(goodsId); goods != nil {
		for _, v := range priceSet {
			if _, err = goods.SaveLevelPrice(v); err != nil {
				return err
			}
		}
	}
	return nil
}

//func (s *saleService) GetGoodsComplexInfo(goodsId int32) *dto.GoodsComplex {
//	return s._goodsQuery.GetGoodsComplex(goodsId)
//}

// 获取商品详情
func (s *saleService) GetGoodsDetails(mchId, goodsId, mLevel int32) (*valueobject.Goods, map[string]string) {
	sl := s._rep.GetSale(mchId)
	var goods item.IGoods = sl.GoodsManager().GetGoods(goodsId)
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
func (s *saleService) GetItemDescriptionByGoodsId(mchId, goodsId int32) string {
	sl := s._rep.GetSale(mchId)
	var goods item.IGoods = sl.GoodsManager().GetGoods(goodsId)
	return goods.GetItem().GetValue().Description
}

// 获取商品快照
func (s *saleService) GetSnapshot(skuId int32) *item.Snapshot {
	return s._goodsRepo.GetLatestSnapshot(skuId)
}

// 设置商品货架状态
func (s *saleService) SetShelveState(mchId int32, itemId int32, state int32, remark string) error {
	panic("应上架商品")
	//sl := s._rep.GetSale(mchId)
	//gi := sl.ItemManager().GetItem(itemId)
	//if gi == nil {
	//	return item.ErrNoSuchGoods
	//}
	//return gi.SetShelve(state, remark)
}

// 设置商品货架状态
func (s *saleService) ReviewItem(mchId int32, itemId int32, pass bool, remark string) error {
	panic("应申请商品")

	//sl := s._rep.GetSale(mchId)
	//gi := sl.ItemManager().GetItem(itemId)
	//if gi == nil {
	//	return item.ErrNoSuchGoods
	//}
	//return gi.Review(pass, remark)
}

// 标记为违规
func (s *saleService) SignIncorrect(supplierId int32, itemId int32, remark string) error {
	panic("应标记商品")
	//sl := s._rep.GetSale(supplierId)
	//gi := sl.ItemManager().GetItem(itemId)
	//if gi == nil {
	//	return item.ErrNoSuchGoods
	//}
	//return gi.Incorrect(remark)
}
