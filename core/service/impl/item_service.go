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
	"errors"
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
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
	promRepo  promodel.IProductModelRepo
	mchRepo   merchant.IMerchantRepo
	valueRepo valueobject.IValueRepo
	sto       storage.Interface
}

func NewSaleService(sto storage.Interface, cateRepo product.ICategoryRepo,
	goodsRepo item.IGoodsItemRepo, goodsQuery *query.ItemQuery,
	labelRepo item.ISaleLabelRepo, promRepo promodel.IProductModelRepo,
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

// 获取商品
func (s *itemService) GetItem(_ context.Context, id *proto.Int64) (*proto.SUnifiedViewItem, error) {
	item := s.itemRepo.GetItem(id.Value)
	if item != nil {
		return s.attachUnifiedItem(item), nil
	}
	return nil, nil
}

// 保存商品
func (s *itemService) SaveItem(_ context.Context, r *proto.SUnifiedViewItem) (*proto.Result, error) {
	var gi item.IGoodsItem
	it := parser.ParseGoodsItem(r)
	var err error
	if it.Id > 0 {
		gi = s.itemRepo.GetItem(it.Id)
		if gi == nil || gi.GetValue().VendorId != r.VendorId {
			return s.error(item.ErrNoSuchItem), nil
		}
	} else {
		gi = s.itemRepo.CreateItem(it)
	}
	err = gi.SetValue(it)
	if err == nil {
		err = gi.SetSku(it.SkuArray)
		if err == nil {
			it.Id, err = gi.Save()
		}
	}
	ret := s.error(err)
	if err == nil {
		r.Data = map[string]string{
			"ItemId": strconv.Itoa(int(it.Id)),
		}
	}
	return ret, nil
}

// 附加商品的信息
func (s *itemService) attachUnifiedItem(item item.IGoodsItem) *proto.SUnifiedViewItem {
	ret := parser.ItemDtoV2(item.GetValue())
	skuService := s.itemRepo.SkuService()
	skuArr := item.SkuArray()
	ret.SkuArray = parser.SkuArrayDto(skuArr)
	ret.LevelPrices = parser.PriceArrayDto(item.GetLevelPrices())
	specArr := item.SpecArray()
	ret.ViewData = &proto.SItemViewData{
		Details: "",  //todo:??
		Thumbs:  nil, //todo:??
		Images:  nil, //todo:??
		SkuHtml: skuService.GetSpecHtml(specArr),
		SkuJson: string(skuService.GetSkuJson(skuArr)),
	}
	return ret
}

// 根据SKU获取商品
func (s *itemService) GetItemBySku(_ context.Context, r *proto.ItemBySkuRequest) (*proto.SUnifiedViewItem, error) {
	v := s.itemRepo.GetValueGoodsBySku(r.ProductId, r.SkuId)
	if v != nil {
		item := s.itemRepo.CreateItem(v)
		return s.attachUnifiedItem(item), nil
	}
	return nil, nil
}

// 获取商品用于销售的快照和信息
func (s *itemService) GetItemSnapshot(_ context.Context, id *proto.Int64) (*proto.SItemSnapshot, error) {
	item := s.itemRepo.GetItem(id.Value)
	if item != nil {
		skuService := s.itemRepo.SkuService()
		sn := item.Snapshot()
		// 基础数据及其销售数量
		ret := parser.ParseItemSnapshotDto(sn)
		ret.Stock.SaleNum = item.GetValue().SaleNum
		ret.Stock.StockNum = item.GetValue().StockNum
		// 获取SKU和详情等
		skuArr := item.SkuArray()
		ret.SkuArray = parser.SkuArrayDto(skuArr)
		specArr := item.SpecArray()
		prod := item.Product()
		ret.ViewData = &proto.SItemViewData{
			Details: prod.GetValue().Description,
			Thumbs:  nil, //todo:??
			Images:  nil, //todo:??
			SkuHtml: skuService.GetSpecHtml(specArr),
			SkuJson: string(skuService.GetSkuJson(skuArr)),
		}
		return ret, nil
	}
	return nil, nil
}

// 获取商品交易快照
func (s *itemService) GetTradeSnapshot(_ context.Context, id *proto.Int64) (*proto.STradeSnapshot, error) {
	sn := s.itemRepo.GetSalesSnapshot(id.Value)
	if sn != nil {
		return parser.ParseTradeSnapshot(sn), nil
	}
	return nil, nil
}

// 获取SKU
func (s *itemService) GetSku(_ context.Context, request *proto.SkuId) (*proto.SSku, error) {
	item := s.itemRepo.GetItem(request.ItemId)
	if item != nil {
		sku := item.GetSku(request.SkuId)
		if sku != nil {
			return s.parseSkuDto(sku), nil
		}
	}
	return nil, nil
}

// 获取商品详细数据
func (s *itemService) GetItemDetailData(_ context.Context, request *proto.ItemDetailRequest) (*proto.String, error) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{Id: request.ItemId})
	switch request.IType {
	case item.ItemWholesale:
		data := it.Wholesale().GetJsonDetailData()
		return &proto.String{Value: string(data)}, nil
	}
	return &proto.String{Value: "不支持的商品类型"}, nil
}

// 获取上架商品数据（分页）
func (s *itemService) GetPagedOnShelvesItem(_ context.Context, r *proto.PagingGoodsRequest) (*proto.PagingGoodsResponse, error) {
	ret := &proto.PagingGoodsResponse{
		Total: 0,
		Data:  make([]*proto.SUnifiedViewItem, 0),
	}
	var total int32
	var list []*proto.SUnifiedViewItem
	switch r.ItemType {
	case proto.EItemSalesType_IT_NORMAL:
		total, list = s.getPagedOnShelvesItem(
			int32(r.CategoryId),
			int32(r.Params.Begin),
			int32(r.Params.End),
			r.Params.Where,
			r.Params.SortBy)
	case proto.EItemSalesType_IT_WHOLESALE:
		total, list = s.getPagedOnShelvesItemForWholesale(
			int32(r.CategoryId),
			int32(r.Params.Begin),
			int32(r.Params.End),
			r.Params.Where,
			r.Params.SortBy)
	}
	ret.Total = int64(total)
	ret.Data = list
	return ret, nil
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
		b := s.promRepo.BrandService().Get(int32(dto.BrandId))
		if b != nil {
			dto.Data["BrandName"] = b.Name
			dto.Data["BrandImage"] = b.Image
			dto.Data["BrandId"] = strconv.Itoa(int(b.ID))
		}
	}
}

// 获取上架商品数据
func (s *itemService) GetItems(_ context.Context, r *proto.GetItemsRequest) (*proto.PagingGoodsResponse, error) {
	/*
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
	*/
	c := s.cateRepo.GlobCatService().GetCategory(int(r.CategoryId))
	var idList []int = nil
	var list []*item.GoodsItem
	if c != nil {
		idList = c.GetChildes()
	}
	if r.Random {
		list = s.itemQuery.GetRandomItem(idList, int(r.Begin), int(r.End), r.Where)

	} else {
		list = s.itemQuery.GetOnShelvesItem(idList, int(r.Begin), int(r.End), r.Where)
	}
	arr := make([]*proto.SUnifiedViewItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ItemDtoV2(v)
	}
	return &proto.PagingGoodsResponse{
		Total: 0,
		Data:  arr,
	}, nil
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

func (s *itemService) GetAllSaleLabels(_ context.Context, _ *proto.Empty) (*proto.ItemLabelListResponse, error) {
	tags := s.labelRepo.LabelService().GetAllSaleLabels()
	ret := &proto.ItemLabelListResponse{
		Value: make([]*proto.SItemLabel, len(tags)),
	}
	for i, v := range tags {
		ret.Value[i] = parser.ParseSaleLabelDto(v.GetValue())
	}
	return ret, nil
}

func (s *itemService) GetSaleLabel(_ context.Context, id *proto.IdOrName) (*proto.SItemLabel, error) {
	var tag item.ISaleLabel
	if id.Id <= 0 {
		tag = s.labelRepo.LabelService().GetSaleLabelByCode(id.Name)
	} else {
		tag = s.labelRepo.LabelService().GetSaleLabel(int32(id.Id))
	}
	if tag != nil {
		return parser.ParseSaleLabelDto(tag.GetValue()), nil
	}
	return nil, nil
}

// 保存销售标签
func (s *itemService) SaveSaleLabel(_ context.Context, v *proto.SItemLabel) (*proto.Result, error) {
	ls := s.labelRepo.LabelService()
	var value = parser.FromSaleLabelDto(v)
	var lb item.ISaleLabel
	if v.Id > 0 {
		if lb = ls.GetSaleLabel(v.Id); lb == nil {
			return s.error(errors.New("没有销售标签")), nil
		}
	} else {
		lb = ls.CreateSaleLabel(value)
	}
	err := lb.SetValue(value)
	if err == nil {
		_, err = lb.Save()
	}
	return s.error(err), nil
}

// 删除销售标签
func (s *itemService) DeleteSaleLabel(_ context.Context, i *proto.Int64) (*proto.Result, error) {
	err := s.labelRepo.LabelService().DeleteSaleLabel(int32(i.Value))
	return s.error(err), nil
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
func (s *itemService) GetPagedValueGoodsBySaleLabel_(_ context.Context, r *proto.SaleLabelItemsRequest_) (*proto.PagingGoodsResponse, error) {
	tag := s.labelRepo.LabelService().CreateSaleLabel(&item.Label{
		Id: r.LabelId,
	})
	total, list := tag.GetPagedValueGoods(r.Params.SortBy, int(r.Params.Begin), int(r.Params.End))
	arr := make([]*proto.SUnifiedViewItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ParseGoodsDto_(v)
	}
	return &proto.PagingGoodsResponse{
		Total: int64(total),
		Data:  arr,
	}, nil
}

// 保存商品的会员价
func (s *itemService) SaveLevelPrices(_ context.Context, r *proto.SaveLevelPriceRequest) (*proto.Result, error) {
	it := s.itemRepo.GetItem(r.ItemId)
	var err error
	if it != nil {
		for _, v := range r.Value {
			e := parser.ParseLevelPrice(v)
			e.GoodsId = r.ItemId
			if _, err := it.SaveLevelPrice(e); err != nil {
				break
			}
		}
	}
	return s.error(err), nil
}

// 设置商品货架状态
func (s *itemService) SetShelveState(_ context.Context, r *proto.ShelveStateRequest) (*proto.Result, error) {
	it := s.itemRepo.GetItem(r.ItemId)
	var err error
	if it == nil || it.GetValue().VendorId != r.SellerId {
		err = item.ErrNoSuchItem
	} else {
		state := int32(types.IntCond(r.ShelveOn,
			int(item.ShelvesOn),
			int(item.ShelvesDown)))
		switch r.ItemType {
		case proto.EItemSalesType_IT_WHOLESALE:
			err = it.Wholesale().SetShelve(state, r.Remark)
		case proto.EItemSalesType_IT_NORMAL:
			err = it.SetShelve(state, r.Remark)
		}
	}
	return s.result(err), nil
}

// 审核商品
func (s *itemService) ReviewItem(_ context.Context, r *proto.ItemReviewRequest) (*proto.Result, error) {
	it := s.itemRepo.GetItem(r.ItemId)
	var err error
	if it == nil {
		err = item.ErrNoSuchItem
	} else {
		err = it.Review(r.Pass, r.Remark)
	}
	return s.result(err), nil
}

// 标记为违规
func (s *itemService) SignAsIllegal(_ context.Context, r *proto.ItemIllegalRequest) (*proto.Result, error) {
	it := s.itemRepo.GetItem(r.ItemId)
	var err error
	if it == nil {
		err = item.ErrNoSuchItem
	} else {
		err = it.Incorrect(r.Remark)
	}
	return s.error(err), nil
}

// 获取批发价格数组
func (s *itemService) GetWholesalePriceArray(_ context.Context, id *proto.SkuId) (*proto.SWsSkuPriceListResponse, error) {
	it := s.itemRepo.GetItem(id.ItemId)
	arr := it.Wholesale().GetSkuPrice(id.SkuId)
	ret := make([]*proto.SWsSkuPrice, len(arr))
	for i, v := range arr {
		ret[i] = parser.WsSkuPriceDto(v)
	}
	return &proto.SWsSkuPriceListResponse{
		Value: ret,
	}, nil
}

// 保存批发价格
func (s *itemService) SaveWholesalePrice(_ context.Context, r *proto.SaveSkuPricesRequest) (*proto.Result, error) {
	it := s.itemRepo.GetItem(r.ItemId)
	arr := make([]*item.WsSkuPrice, len(r.Value))
	for i, v := range r.Value {
		e := parser.WsSkuPrice(v)
		e.ItemId = r.ItemId
		e.SkuId = r.SkuId
		arr[i] = e
	}
	err := it.Wholesale().SaveSkuPrice(r.SkuId, arr)
	return s.error(err), nil
}

// 获取批发折扣数组
func (s *itemService) GetWholesaleDiscountArray(_ context.Context, id *proto.GetWsDiscountRequest) (*proto.SWsItemDiscountListResponse, error) {

	it := s.itemRepo.GetItem(id.ItemId)
	arr := it.Wholesale().GetItemDiscount(int32(id.GroupId))
	ret := make([]*proto.SWsItemDiscount, len(arr))
	for i, v := range arr {
		ret[i] = parser.WsItemDiscountDto(v)
	}
	return &proto.SWsItemDiscountListResponse{
		Value: ret,
	}, nil
}

// 保存批发折扣
func (s *itemService) SaveWholesaleDiscount(_ context.Context, r *proto.SaveItemDiscountRequest) (*proto.Result, error) {
	it := s.itemRepo.GetItem(r.ItemId)
	arr := make([]*item.WsItemDiscount, len(r.Value))
	for i, v := range r.Value {
		e := parser.WsItemDiscount(v)
		e.ItemId = r.ItemId
		e.BuyerGid = int32(r.GroupId)
		arr[i] = e
	}
	err := it.Wholesale().SaveItemDiscount(int32(r.GroupId), arr)
	return s.error(err), nil
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
