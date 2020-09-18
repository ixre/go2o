/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package impl

import (
	"context"
	"fmt"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/domain/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/core/service/parser"
	"go2o/core/service/proto"
	"strconv"
)

var _ proto.ItemServiceServer = new(itemService)

type itemService struct {
	serviceUtil
	itemRepo  item.IGoodsItemRepo
	itemQuery *query.ItemQuery
	cateRepo  product.ICategoryRepo
	labelRepo item.ISaleLabelRepo
	promRepo  promodel.IProModelRepo
	mchRepo   merchant.IMerchantRepo
	valueRepo valueobject.IValueRepo
	sto       storage.Interface
}


func NewSaleService(sto storage.Interface, cateRepo product.ICategoryRepo,
	goodsRepo item.IGoodsItemRepo, goodsQuery *query.ItemQuery,
	labelRepo item.ISaleLabelRepo, promRepo promodel.IProModelRepo,
	mchRepo merchant.IMerchantRepo, valueRepo valueobject.IValueRepo) *itemService {
	return &itemService{
		sto:       sto,
		itemRepo:  goodsRepo,
		itemQuery: goodsQuery,
		cateRepo:  cateRepo,
		labelRepo: labelRepo,
		promRepo:  promRepo,
		mchRepo:   mchRepo,
		valueRepo: valueRepo,
	}
}

// 获取商品值
func (s *itemService) GetItemValue(itemId int64) *proto.SOldItem {
	item := s.itemRepo.GetItem(itemId)
	if item != nil {
		return parser.ItemDto(item.GetValue())
	}
	return nil
}

// 获取SKU
func (s *itemService) GetSku(_ context.Context, request *proto.SkuRequest) (*proto.SSku, error) {
	item := s.itemRepo.GetItem(request.ItemId)
	if item != nil {
		sku := item.GetSku(request.SkuId)
		if sku != nil {
			return s.parseSkuDto(sku), nil
		}
	}
	return nil, nil
}

// 获取SKU数组
func (s *itemService) GetSkuArray(itemId int64) []*item.Sku {
	it := s.itemRepo.GetItem(itemId)
	if it != nil {
		return it.SkuArray()
	}
	return []*item.Sku{}
}

// 获取商品规格HTML信息
func (s *itemService) GetSkuHtmOfItem(itemId int64) (specHtm string) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: itemId})
	specArr := it.SpecArray()
	return s.itemRepo.SkuService().GetSpecHtml(specArr)
}

// 获取商品详细数据
func (s *itemService) GetItemDetailData(_ context.Context, request *proto.ItemDetailRequest) (*proto.String, error) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: request.ItemId})
	switch request.IType {
	case item.ItemWholesale:
		data := it.Wholesale().GetJsonDetailData()
		return &proto.String{Value: string(data)}, nil
	}
	return &proto.String{Value: "不支持的商品类型"}, nil
}

// 获取商品的Sku-JSON格式
func (s *itemService) GetItemSkuJson(_ context.Context, i *proto.Int64) (*proto.String, error) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: i.Value})
	skuBytes := s.itemRepo.SkuService().GetSkuJson(it.SkuArray())
	return &proto.String{Value: string(skuBytes)}, nil
}

// 保存商品
func (s *itemService) SaveItem(di *proto.SOldItem, vendorId int64) (_ *proto.Result, err error) {
	var gi item.IGoodsItem
	it := parser.Item(di)
	if it.ID > 0 {
		gi = s.itemRepo.GetItem(it.ID)
		if gi == nil || gi.GetValue().VendorId != vendorId {
			err = item.ErrNoSuchItem
			goto R
		}
	} else {
		gi = s.itemRepo.CreateItem(it)
	}
	err = gi.SetValue(it)
	if err == nil {
		err = gi.SetSku(it.SkuArray)
		if err == nil {
			it.ID, err = gi.Save()
		}
	}
R:
	r := s.result(err)
	r.Data = map[string]string{
		"ItemId": strconv.Itoa(int(it.ID)),
	}
	return r, nil
}

// 获取上架商品数据（分页）
func (s *itemService) GetPagedOnShelvesItem(_ context.Context, r *proto.PagingGoodsRequest) (*proto.PagingGoodsResponse, error) {
	ret := &proto.PagingGoodsResponse{
		Total:                0,
		Data:                 make([]*proto.SUnifiedViewItem,0),
	}
	var total int32
	var list []*proto.SUnifiedViewItem
	switch r.ItemType {
	case proto.EItemSalesType_IT_NORMAL:
		total,list = s.getPagedOnShelvesItem(
			int32(r.CategoryId),
			int32(r.Params.Begin),
			int32(r.Params.End),
			r.Params.Where,
			r.Params.SortBy)
	case proto.EItemSalesType_IT_WHOLESALE:
		total,list = s.getPagedOnShelvesItemForWholesale(
			int32(r.CategoryId),
			int32(r.Params.Begin),
			int32(r.Params.End),
			r.Params.Where,
			r.Params.SortBy)
	}
	ret.Total = int64(total)
	ret.Data = list
	return ret,nil
}
func (s *itemService) getPagedOnShelvesItem(catId int32, start,
	end int32, where, sortBy string) (int32, []*proto.SUnifiedViewItem) {

	total, list := s.itemQuery.GetPagedOnShelvesItem(catId,
		start, end, where, sortBy)
	arr := make([]*proto.SUnifiedViewItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ItemDtoV2(v)
	}
	return total, arr
}

func (s *itemService) getPagedOnShelvesItemForWholesale(catId int32, start,
	end int32, where, sortBy string) (int32, []*proto.SUnifiedViewItem) {

	total, list := s.itemQuery.GetPagedOnShelvesItemForWholesale(catId,
		start, end, where, sortBy)
	arr := make([]*proto.SUnifiedViewItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		dto := parser.ItemDtoV2(v)
		s.attachWholesaleItemDataV2(dto)
		arr[i] = dto
	}
	return total, arr
}

// 获取上架商品数据（分页）
func (s *itemService) SearchOnShelvesItem(itemType int32, word string, start,
	end int32, where, sortBy string) (int32, []*proto.SOldItem) {

	switch itemType {
	case item.ItemNormal:
		return s.searchOnShelveItem(word, start, end, where, sortBy)
	case item.ItemWholesale:
		return s.searchOnShelveItemForWholesale(word, start, end, where, sortBy)
	}
	return 0, []*proto.SOldItem{}
}

func (s itemService) searchOnShelveItem(word string, start,
	end int32, where, sortBy string) (int32, []*proto.SOldItem) {
	total, list := s.itemQuery.SearchOnShelvesItem(word,
		start, end, where, sortBy)
	arr := make([]*proto.SOldItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ItemDto(v)
	}
	return total, arr
}

func (s itemService) searchOnShelveItemForWholesale(word string, start,
	end int32, where, sortBy string) (int32, []*proto.SOldItem) {
	total, list := s.itemQuery.SearchOnShelvesItemForWholesale(word,
		start, end, where, sortBy)
	arr := make([]*proto.SOldItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		dto := parser.ItemDto(v)
		s.attachWholesaleItemData(dto)
		arr[i] = dto

	}
	return total, arr
}

// 附加批发商品的信息
func (s *itemService) attachWholesaleItemData(dto *proto.SOldItem) {
	dto.Data = make(map[string]string)
	vendor := s.mchRepo.GetMerchant(int(dto.VendorId))
	if vendor != nil {
		vv := vendor.GetValue()
		pStr := s.valueRepo.GetAreaName(int32(vv.Province))
		cStr := s.valueRepo.GetAreaName(int32(vv.City))
		dto.Data["VendorName"] = vv.CompanyName
		dto.Data["ShipArea"] = pStr + cStr
		// 认证信息
		ei := vendor.ProfileManager().GetEnterpriseInfo()
		if ei != nil && ei.Reviewed == enum.ReviewPass {
			dto.Data["Authorized"] = "true"
		} else {
			dto.Data["Authorized"] = "false"
		}
		// 品牌
		b := s.promRepo.BrandService().Get(dto.BrandId)
		if b != nil {
			dto.Data["BrandName"] = b.Name
			dto.Data["BrandImage"] = b.Image
			dto.Data["BrandId"] = strconv.Itoa(int(b.ID))
		}
	}
}

// 附加批发商品的信息
func (s *itemService) attachWholesaleItemDataV2(dto *proto.SUnifiedViewItem) {
	dto.Data = make(map[string]string)
	vendor := s.mchRepo.GetMerchant(int(dto.VendorId))
	if vendor != nil {
		vv := vendor.GetValue()
		pStr := s.valueRepo.GetAreaName(int32(vv.Province))
		cStr := s.valueRepo.GetAreaName(int32(vv.City))
		dto.Data["VendorName"] = vv.CompanyName
		dto.Data["ShipArea"] = pStr + cStr
		// 认证信息
		ei := vendor.ProfileManager().GetEnterpriseInfo()
		if ei != nil && ei.Reviewed == enum.ReviewPass {
			dto.Data["Authorized"] = "true"
		} else {
			dto.Data["Authorized"] = "false"
		}
		// 品牌
		b := s.promRepo.BrandService().Get(dto.BrandId)
		if b != nil {
			dto.Data["BrandName"] = b.Name
			dto.Data["BrandImage"] = b.Image
			dto.Data["BrandId"] = strconv.Itoa(int(b.ID))
		}
	}
}


// 获取上架商品数据（分页）
func (s *itemService) GetRandomItem(catId int32, quantity int32, where string) []*proto.SOldItem {
	hash := fmt.Sprintf("%d-%d-%s", catId, quantity, where)
	hash = crypto.Md5([]byte(hash))
	key := "go2o:shopQuery:cache:rd-item:" + hash
	var arr []*proto.SOldItem

	fn := func() interface{} {
		list := s.itemQuery.GetRandomItem(catId, quantity, where)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
			arr = append(arr, parser.ItemDto(v))
		}
		return arr
	}
	s.sto.RWJson(key, &arr, fn, 600)
	return arr
}

// 获取上架商品数据（分页）
func (s *itemService) GetBigCatItems(catId, quantity int32, where string) []*proto.SOldItem {
	c := s.cateRepo.GlobCatService().GetCategory(int(catId))
	if c != nil {
		ids := c.GetChildes()
		list := s.itemQuery.GetOnShelvesItem(ids, 0, quantity, where)
		arr := make([]*proto.SOldItem, len(list))
		for i, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
			arr[i] = parser.ItemDto(v)
		}
		return arr
	}
	return []*proto.SOldItem{}
}

// 根据SKU获取商品
func (s *itemService) GetGoodsBySku(mchId, itemId, sku int64) *valueobject.Goods {
	v := s.itemRepo.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		return s.itemRepo.CreateItem(v).GetPackedValue()
	}
	return nil
}

// 根据SKU获取商品
func (s *itemService) GetValueGoodsBySku(mchId int64, itemId, sku int64) *item.GoodsItem {
	v := s.itemRepo.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		return s.itemRepo.CreateItem(v).GetValue()
	}
	return nil
}

// 根据快照编号获取商品
func (s *itemService) GetGoodsBySnapshotId(snapshotId int64) *item.GoodsItem {
	snap := s.itemRepo.GetSalesSnapshot(snapshotId)
	if snap != nil {
		return s.itemRepo.GetValueGoodsById(snap.SkuId)
	}
	return nil
}

// 根据快照编号获取商品
func (s *itemService) GetSaleSnapshotById(snapshotId int64) *item.TradeSnapshot {
	return s.itemRepo.GetSalesSnapshot(snapshotId)
}

// 获取分页上架的商品
func (s *itemService) GetShopPagedOnShelvesGoods(_ context.Context, r *proto.PagingShopGoodsRequest) (*proto.PagingShopGoodsResponse, error) {
	ret := &proto.PagingShopGoodsResponse{
		Total: 0,
		Data:  make([]*proto.SGoods, 0),
	}
	var list []*valueobject.Goods
	var total int
	var ids []int
	if r.CategoryId > 0 {
		cat := s.cateRepo.GlobCatService().GetCategory(int(r.CategoryId))
		if cat == nil {
			return ret, nil
		}
		ids = cat.GetChildes()
		ids = append(ids, int(r.CategoryId))

	}
	total, list = s.itemRepo.GetPagedOnShelvesGoods(
		r.ShopId, ids,
		int(r.Params.Begin),
		int(r.Params.End),
		r.Params.Where,
		r.Params.SortBy)
	ret.Total = int64(total)
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		ret.Data = append(ret.Data, s.parseGoods(v))
	}
	return ret, nil
}

// 获取分页上架的商品
func (s *itemService) GetPagedOnShelvesGoodsByKeyword(shopId int64, start, end int,
	word, sortQuery string) (int, []*valueobject.Goods) {
	var where string
	var orderBy string
	switch sortQuery {
	case "price_0":
		where = ""
		orderBy = "it.price ASC"
	case "price_1":
		where = ""
		orderBy = "it.price DESC"
	case "sale_0":
		where = ""
		orderBy = "it.sale_num ASC"
	case "sale_1":
		where = ""
		orderBy = "it.sale_num DESC"
	case "rate_0":
		//todo:
	case "rate_1":
		//todo:
	}
	return s.itemQuery.GetPagedOnShelvesGoodsByKeyword(shopId,
		start, end, word, where, orderBy)
}

// 删除产品
func (s *itemService) DeleteGoods(mchId, itemId int64) error {
	gi := s.itemRepo.GetItem(itemId)
	if gi == nil || gi.GetValue().VendorId != mchId {
		return item.ErrNoSuchItem
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
func (s *itemService) GetValueGoodsBySaleLabel(_ context.Context, r *proto.GetItemsByLabelRequest) (*proto.PagingShopGoodsResponse, error) {
	tag := s.labelRepo.LabelService().GetSaleLabelByCode(r.Label)
	ret := &proto.PagingShopGoodsResponse{
		Data: make([]*proto.SGoods, 0),
	}
	if tag != nil {
		list := tag.GetValueGoods(r.SortBy, int(r.Begin), int(r.End))
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
			ret.Data = append(ret.Data, s.parseGoods(v))
		}
	}
	return ret, nil
}

// 根据分页销售标签获取指定数目的商品
func (s *itemService) GetPagedValueGoodsBySaleLabel(shopId int64, tagId int32, sortBy string, begin int, end int) (int,
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
func (s *itemService) GetGoodSMemberLevelPrices(itemId int64) []*item.MemberPrice {
	gi := s.itemRepo.GetItem(itemId)
	if gi != nil {
		return gi.GetLevelPrices()
	}
	return make([]*item.MemberPrice, 0)
}

// 保存商品的会员价
func (s *itemService) SaveMemberPrices(mchId, itemId int64,
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
func (s *itemService) GetGoodsDetails(itemId int64, mLevel int32) (
	*valueobject.Goods, map[string]string) {
	goods := s.itemRepo.GetItem(itemId)
	gv := goods.GetPackedValue()
	proMap := goods.GetPromotionDescribe()
	if b, price := goods.GetLevelPrice(int(mLevel)); b {
		gv.PromPrice = price
		proMap["会员专享"] = fmt.Sprintf("会员优惠,仅需<b>￥%s</b>",
			format.FormatFloat(price))
	}
	return gv, proMap
}

// 获取货品描述
func (s *itemService) GetItemDescriptionByGoodsId(itemId int64) string {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: itemId})
	pro := it.Product()
	if pro != nil {
		return pro.GetValue().Description
	}
	return ""
}

// 获取商品快照
func (s *itemService) GetSnapshot(skuId int64) *item.Snapshot {
	return s.itemRepo.GetLatestSnapshot(skuId)
}

// 设置商品货架状态
func (s *itemService) SetShelveState(vendorId, itemId int64,
	itemType int32, state int32, remark string) (_ *proto.Result, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchItem
	} else {
		switch itemType {
		case item.ItemWholesale:
			err = it.Wholesale().SetShelve(state, remark)
		default:
			err = it.SetShelve(state, remark)
		}
	}
	return s.result(err), nil
}

// 设置商品货架状态
func (s *itemService) ReviewItem(vendorId, itemId int64,
	pass bool, remark string) (_ *proto.Result, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchItem
	} else {
		err = it.Review(pass, remark)
	}
	return s.result(err), nil
}

// 标记为违规
func (s *itemService) SignGoodsIllegal(vendorId, itemId int64,
	remark string) (_ *proto.Result, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchItem
	} else {
		err = it.Incorrect(remark)
	}
	return s.result(err), nil
}

// 获取批发价格数组
func (s *itemService) GetWholesalePriceArray(itemId, skuId int64) []*item.WsSkuPrice {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().GetSkuPrice(skuId)
}

// 保存批发价格
func (s *itemService) SaveWholesalePrice(itemId, skuId int64, arr []*item.WsSkuPrice) error {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().SaveSkuPrice(skuId, arr)
}

// 获取批发折扣数组
func (s *itemService) GetWholesaleDiscountArray(itemId int64, groupId int32) []*item.WsItemDiscount {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().GetItemDiscount(groupId)
}

// 保存批发折扣
func (s *itemService) SaveWholesaleDiscount(itemId int64, groupId int32, arr []*item.WsItemDiscount) error {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().SaveItemDiscount(groupId, arr)
}

func (s *itemService) parseSkuDto(sku *item.Sku) *proto.SSku {
	return &proto.SSku{
		SkuId:       sku.ID,
		ProductId:   sku.ProductId,
		ItemId:      sku.ItemId,
		Title:       sku.Title,
		Image:       sku.Image,
		SpecData:    sku.SpecData,
		SpecWord:    sku.SpecWord,
		Code:        sku.Code,
		RetailPrice: math.Round(float64(sku.RetailPrice), 2),
		Price:       math.Round(float64(sku.Price), 2),
		Cost:        math.Round(float64(sku.Cost), 2),
		Weight:      sku.Weight,
		Bulk:        sku.Bulk,
		Stock:       sku.Stock,
		SaleNum:     sku.SaleNum,
	}
}

func (s *itemService) parseGoods(v *valueobject.Goods) *proto.SGoods {
	return &proto.SGoods{
		ItemId:        v.ItemId,
		ProductId:     v.ProductId,
		VendorId:      int64(v.VendorId),
		ShopId:        int64(v.ShopId),
		CategoryId:    v.CategoryId,
		Title:         v.Title,
		ShortTitle:    v.ShortTitle,
		GoodsNo:       v.GoodsNo,
		Image:         v.Image,
		RetailPrice:   float64(v.RetailPrice),
		Price:         float64(v.ProductId),
		PromPrice:     float64(v.PromPrice),
		PriceRange:    v.PriceRange,
		GoodsId:       v.GoodsId,
		SkuId:         v.SkuId,
		IsPresent:     v.IsPresent == 1,
		PromotionFlag: v.PromotionFlag,
		StockNum:      v.StockNum,
		SaleNum:       v.SaleNum,
	}
}
