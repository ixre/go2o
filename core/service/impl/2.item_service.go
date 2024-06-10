/**
 * Copyright 2014 @ 56x.net.
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

	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/parser"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/types"
)

var _ proto.ItemServiceServer = new(itemService)

type itemService struct {
	serviceUtil
	itemRepo  item.IItemRepo
	itemQuery *query.ItemQuery
	cateRepo  product.ICategoryRepo
	labelRepo item.ISaleLabelRepo
	promRepo  promodel.IProductModelRepo
	mchRepo   merchant.IMerchantRepo
	valueRepo valueobject.IValueRepo
	sto       storage.Interface
	proto.UnimplementedItemServiceServer
}

func NewItemService(sto storage.Interface, cateRepo product.ICategoryRepo,
	goodsRepo item.IItemRepo,
	goodsQuery *query.ItemQuery,
	labelRepo item.ISaleLabelRepo, promRepo promodel.IProductModelRepo,
	mchRepo merchant.IMerchantRepo, valueRepo valueobject.IValueRepo) proto.ItemServiceServer {
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

// GetItem 获取商品
func (i *itemService) GetItem(_ context.Context, req *proto.GetItemRequest) (*proto.SItemDataResponse, error) {
	it := i.itemRepo.GetItem(req.ItemId)
	if it != nil {
		ret := parser.ItemDataDto(it.GetValue())
		skuArr := it.SkuArray()
		specArr := it.SpecArray()
		ret.Images = it.Images()
		ret.AttrArray = parser.AttrArrayDto(it.Product().Attr())
		ret.SkuArray = parser.SkuArrayDto(skuArr)
		ret.LevelPrices = parser.PriceArrayDto(it.GetLevelPrices())
		ret.SpecOptions = parser.SpecOptionsDto(specArr)
		ret.FlagData = i.getFlagData(it.GetValue().ItemFlag)
		return ret, nil
	}
	return nil, item.ErrNoSuchItem
}

// 获取商品的标志数据
func (i *itemService) getFlagData(flag int) *proto.SItemFlagData {
	return &proto.SItemFlagData{
		IsNewOnShelve:  domain.TestFlag(flag, item.FlagNewOnShelve),
		IsHotSales:     domain.TestFlag(flag, item.FlagHotSales),
		IsRecommend:    domain.TestFlag(flag, item.FlagRecommend),
		IsExchange:     domain.TestFlag(flag, item.FlagExchange),
		IsGift:         domain.TestFlag(flag, item.FlagGift),
		IsAffiliate:    domain.TestFlag(flag, item.FlagAffiliate),
		IsSelfSales:    domain.TestFlag(flag, item.FlagSelfSales),
		IsFreeDelivery: domain.TestFlag(flag, item.FlagFreeDelivery),
		IsSelfDelivery: domain.TestFlag(flag, item.FlagSelfDelivery),
	}
}

// SaveItem 保存商品
func (i *itemService) SaveItem(_ context.Context, r *proto.SaveItemRequest) (*proto.SaveItemResponse, error) {
	var gi item.IGoodsItemAggregateRoot
	it := parser.ParseGoodsItem(r)
	var err error
	if it.Id > 0 {
		gi = i.itemRepo.GetItem(it.Id)
		if gi == nil || gi.GetValue().VendorId != r.VendorId {
			return &proto.SaveItemResponse{
				ErrCode: 1,
				ErrMsg:  item.ErrNoSuchItem.Error(),
			}, nil
		}
	} else {
		gi = i.itemRepo.CreateItem(it)
	}
	err = gi.SetValue(it)
	if err == nil {
		err = gi.SetSku(it.SkuArray)
		if err == nil {
			if r.Images != nil && len(r.Images) > 0 {
				_ = gi.SetImages(r.Images)
			}
			it.Id, err = gi.Save()
		}
	}
	if err != nil {
		return &proto.SaveItemResponse{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, nil
	}
	if err == nil && r.FlagData != nil {
		i.saveItemFlag(gi, r)
	}
	return &proto.SaveItemResponse{
		ErrCode:  0,
		ItemId:   it.Id,
		ItemFlag: int32(it.ItemFlag),
	}, nil
}

// 保存商品标志
func (i *itemService) saveItemFlag(gi item.IGoodsItemAggregateRoot, r *proto.SaveItemRequest) {
	f := func(flag int, b bool) {
		if b {
			gi.GrantFlag(flag)
		} else {
			gi.GrantFlag(-flag)
		}
	}
	f(item.FlagNewOnShelve, r.FlagData.IsNewOnShelve)
	f(item.FlagHotSales, r.FlagData.IsHotSales)
	f(item.FlagRecommend, r.FlagData.IsRecommend)
	f(item.FlagExchange, r.FlagData.IsExchange)
	f(item.FlagGift, r.FlagData.IsGift)
	f(item.FlagAffiliate, r.FlagData.IsAffiliate)
	f(item.FlagSelfDelivery, r.FlagData.IsSelfDelivery)
	gi.Save()
}

// 附加商品的信息
func (i *itemService) attachUnifiedItem(item item.IGoodsItemAggregateRoot, extra bool) *proto.SUnifiedViewItem {
	ret := parser.ItemDtoV2(item.GetValue())
	skuService := i.itemRepo.SkuService()
	skuArr := item.SkuArray()
	ret.SkuArray = parser.SkuArrayDto(skuArr)
	ret.LevelPrices = parser.PriceArrayDto(item.GetLevelPrices())
	specArr := item.SpecArray()
	ret.FlagData = i.getFlagData(item.GetValue().ItemFlag)
	if extra {
		ret.ViewData = &proto.SItemViewData{
			Details: "",  //todo:??
			Thumbs:  nil, //todo:??
			Images:  nil, //todo:??
			SkuHtml: skuService.GetSpecHtml(specArr),
			SkuJson: string(skuService.GetItemSkuJson(skuArr)),
		}
	}
	return ret
}

// RecycleItem 回收商品
func (i *itemService) RecycleItem(_ context.Context, req *proto.RecycleItemRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(req.ItemId)
	var err error
	if it == nil {
		err = item.ErrNoSuchItem
	} else {
		if req.IsDestory {
			err = it.Destroy()
		} else {
			if req.Recycle {
				err = it.Recycle()
			} else {
				err = it.RecycleRevert()
			}
		}
	}
	return i.error(err), nil
}

// 根据SKU获取商品
func (i *itemService) GetItemBySku(_ context.Context, r *proto.ItemBySkuRequest) (*proto.SUnifiedViewItem, error) {
	v := i.itemRepo.GetValueGoodsBySku(r.ProductId, r.SkuId)
	if v != nil {
		item := i.itemRepo.CreateItem(v)
		return i.attachUnifiedItem(item, r.Extra), nil
	}
	return nil, item.ErrNoSuchSku
}

// GetItemAndSnapshot 获取商品用于销售的快照和信息
func (i *itemService) GetItemAndSnapshot(_ context.Context, r *proto.GetItemAndSnapshotRequest) (*proto.ItemSnapshotResponse, error) {
	it := i.itemRepo.GetItem(r.GetItemId())
	if it != nil {
		skuService := i.itemRepo.SkuService()
		sn := it.Snapshot()
		// 基础数据及其销售数量
		ret := parser.ParseItemSnapshotDto(sn)
		ret.SaleNum = it.GetValue().SaleNum
		ret.StockNum = it.GetValue().StockNum
		// 图片
		ret.Images = it.Images()
		if len(ret.Images) == 0 && len(sn.Image) > 0 {
			ret.Images = []string{sn.Image}
		}
		// 获取属性
		attrArr := it.Product().Attr()
		ret.AttrArray = parser.AttrValueArrayDto(attrArr)
		// 获取SKU和详情等
		skuArr := it.SkuArray()
		specArr := it.SpecArray()
		ret.SkuArray = parser.SkuArrayDto(skuArr)
		ret.SpecOptions = parser.SpecOptionsDto(specArr)
		// 视频介绍
		ret.IntroVideo = it.GetValue().IntroVideo
		// 商品标志
		ret.FlagData = i.getFlagData(it.GetValue().ItemFlag)
		// 产品详情
		ret.Description = it.Product().GetValue().Description
		// 返回SKU的HTML选择器
		if r.ReturnSkuHtml {
			ret.SkuHtml = skuService.GetSpecHtml(specArr)
		}
		if r.ReturnSkuJson {
			ret.SkuJson = string(skuService.GetItemSkuJson(skuArr)) //todo: 是否可以去掉
		}
		return ret, nil
	}
	return nil, item.ErrNoSuchItem
}

// 获取商品交易快照
func (i *itemService) GetTradeSnapshot(_ context.Context, id *proto.Int64) (*proto.STradeSnapshot, error) {
	sn := i.itemRepo.GetSalesSnapshot(id.Value)
	if sn != nil {
		return parser.ParseTradeSnapshot(sn), nil
	}
	return nil, item.ErrNoSuchSnapshot
}

// 获取SKU
func (i *itemService) GetSku(_ context.Context, request *proto.SkuId) (*proto.SSku, error) {
	it := i.itemRepo.GetItem(request.ItemId)
	if it != nil {
		sku := it.GetSku(request.SkuId)
		if sku != nil {
			return i.parseSkuDto(sku), nil
		}
	}
	return nil, item.ErrNoSuchItem
}

// 获取商品详细数据
func (i *itemService) GetItemDetailData(_ context.Context, request *proto.ItemDetailRequest) (*proto.String, error) {
	it := i.itemRepo.CreateItem(&item.GoodsItem{Id: request.ItemId})
	switch request.ItemType {
	case item.ItemWholesale:
		data := it.Wholesale().GetJsonDetailData()
		return &proto.String{Value: string(data)}, nil
	}
	return &proto.String{Value: "不支持的商品类型"}, nil
}

// 获取上架商品数据
func (i *itemService) GetItems(_ context.Context, r *proto.GetItemsRequest) (*proto.PagingGoodsResponse, error) {
	c := i.cateRepo.GlobCatService().GetCategory(int(r.CategoryId))
	var idList []int = nil
	var list []*item.GoodsItem
	if c != nil {
		idList = c.GetChildes()
	}
	if r.Random {
		list = i.itemQuery.GetRandomItem(idList, int(r.Begin), int(r.End), r.Where)

	} else {
		list = i.itemQuery.GetOnShelvesItem(idList, int(r.Begin), int(r.End), r.Where)
	}
	arr := make([]*proto.SUnifiedViewItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ItemDtoV2(*v)
	}
	return &proto.PagingGoodsResponse{
		Total: 0,
		Data:  arr,
	}, nil
}

func (i *itemService) GetAllSaleLabels(_ context.Context, _ *proto.Empty) (*proto.ItemLabelListResponse, error) {
	tags := i.labelRepo.LabelService().GetAllSaleLabels()
	ret := &proto.ItemLabelListResponse{
		Value: make([]*proto.SItemLabel, len(tags)),
	}
	for i, v := range tags {
		ret.Value[i] = parser.ParseSaleLabelDto(v.GetValue())
	}
	return ret, nil
}

func (i *itemService) GetSaleLabel(_ context.Context, id *proto.IdOrName) (*proto.SItemLabel, error) {
	var tag item.ISaleLabel
	if id.Id <= 0 {
		tag = i.labelRepo.LabelService().GetSaleLabelByCode(id.Name)
	} else {
		tag = i.labelRepo.LabelService().GetSaleLabel(int32(id.Id))
	}
	if tag != nil {
		return parser.ParseSaleLabelDto(tag.GetValue()), nil
	}
	return nil, errors.New("no such sale label")
}

// SaveSaleLabel 保存销售标签
func (i *itemService) SaveSaleLabel(_ context.Context, v *proto.SItemLabel) (*proto.Result, error) {
	ls := i.labelRepo.LabelService()
	var value = parser.FromSaleLabelDto(v)
	var lb item.ISaleLabel
	if v.Id > 0 {
		if lb = ls.GetSaleLabel(v.Id); lb == nil {
			return i.error(errors.New("没有销售标签")), nil
		}
	} else {
		lb = ls.CreateSaleLabel(value)
	}
	err := lb.SetValue(value)
	if err == nil {
		_, err = lb.Save()
	}
	return i.error(err), nil
}

// 删除销售标签
func (i *itemService) DeleteSaleLabel(_ context.Context, s *proto.Int64) (*proto.Result, error) {
	err := i.labelRepo.LabelService().DeleteSaleLabel(int32(s.Value))
	return i.error(err), nil
}

// 保存商品的会员价
func (i *itemService) SaveLevelPrices(_ context.Context, r *proto.SaveLevelPriceRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(r.ItemId)
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
	return i.error(err), nil
}

// 设置商品货架状态
func (i *itemService) SetShelveState(_ context.Context, r *proto.ShelveStateRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(r.ItemId)
	var err error
	if it == nil {
		return i.result(item.ErrNoSuchItem), nil
	}
	state := int32(types.ElseInt(r.ShelveOn,
		int(item.ShelvesOn),
		int(item.ShelvesDown)))
	switch r.ItemType {
	case proto.EItemSalesType_IT_WHOLESALE:
		err = it.Wholesale().SetShelve(state, r.Remark)
	case proto.EItemSalesType_IT_NORMAL:
		err = it.SetShelve(state, r.Remark)
	}
	return i.result(err), nil
}

// 审核商品
func (i *itemService) ReviewItem(_ context.Context, r *proto.ItemReviewRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(r.ItemId)
	var err error
	if it == nil {
		err = item.ErrNoSuchItem
	} else {
		err = it.Review(r.Pass, r.Remark)
	}
	return i.result(err), nil
}

// 标记为违规
func (i *itemService) SignAsIllegal(_ context.Context, r *proto.ItemIllegalRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(r.ItemId)
	var err error
	if it == nil {
		err = item.ErrNoSuchItem
	} else {
		err = it.Incorrect(r.Remark)
	}
	return i.error(err), nil
}

// 获取批发价格数组
func (i *itemService) GetWholesalePriceArray(_ context.Context, id *proto.SkuId) (*proto.SWsSkuPriceListResponse, error) {
	it := i.itemRepo.GetItem(id.ItemId)
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
func (i *itemService) SaveWholesalePrice(_ context.Context, r *proto.SaveSkuPricesRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(r.ItemId)
	arr := make([]*item.WsSkuPrice, len(r.Value))
	for i, v := range r.Value {
		e := parser.WsSkuPrice(v)
		e.ItemId = r.ItemId
		e.SkuId = r.SkuId
		arr[i] = e
	}
	err := it.Wholesale().SaveSkuPrice(r.SkuId, arr)
	return i.error(err), nil
}

// 获取批发折扣数组
func (i *itemService) GetWholesaleDiscountArray(_ context.Context, id *proto.GetWsDiscountRequest) (*proto.SWsItemDiscountListResponse, error) {

	it := i.itemRepo.GetItem(id.ItemId)
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
func (i *itemService) SaveWholesaleDiscount(_ context.Context, r *proto.SaveItemDiscountRequest) (*proto.Result, error) {
	it := i.itemRepo.GetItem(r.ItemId)
	arr := make([]*item.WsItemDiscount, len(r.Value))
	for i, v := range r.Value {
		e := parser.WsItemDiscount(v)
		e.ItemId = r.ItemId
		e.BuyerGid = int32(r.GroupId)
		arr[i] = e
	}
	err := it.Wholesale().SaveItemDiscount(int32(r.GroupId), arr)
	return i.error(err), nil
}

func (i *itemService) parseSkuDto(sku *item.Sku) *proto.SSku {
	return &proto.SSku{
		SkuId:       sku.Id,
		ProductId:   sku.ProductId,
		ItemId:      sku.ItemId,
		Title:       sku.Title,
		Image:       sku.Image,
		SpecData:    sku.SpecData,
		SpecWord:    sku.SpecWord,
		Code:        sku.Code,
		OriginPrice: sku.OriginPrice,
		Price:       sku.Price,
		Cost:        sku.Cost,
		Weight:      sku.Weight,
		Bulk:        sku.Bulk,
		Stock:       sku.Stock,
		SaleNum:     sku.SaleNum,
	}
}
