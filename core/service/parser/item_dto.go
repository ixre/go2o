package parser

import (
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/service/proto"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : item_dto.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-19 13:50
 * description :
 * history :
 */

func ParseSaleLabelDto(src *item.Label) *proto.SItemLabel {
	return &proto.SItemLabel{
		Id:                   src.Id,
		Name:                 src.TagName,
		TagCode:              src.TagCode,
		LabelImage:           src.LabelImage,
		Enabled:              src.Enabled == 1,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

func FromSaleLabelDto(src *proto.SItemLabel) *item.Label {
	return &item.Label{
		Id:         src.Id,
		MerchantId: -1,
		TagCode:    src.TagCode,
		TagName:    src.Name,
		LabelImage: src.LabelImage,
		Enabled:    types.IntCond(src.Enabled, 1, 0),
	}
}

func PriceArrayDto(src []*item.MemberPrice) []*proto.SLevelPrice {
	dst := make([]*proto.SLevelPrice, len(src))
	for i, v := range src {
		dst[i] = LevelPriceDto(v)
	}
	return dst
}

func LevelPriceDto(src *item.MemberPrice) *proto.SLevelPrice {
	return &proto.SLevelPrice{
		Id:        int64(src.Id),
		Level:     int32(src.Level),
		Price:     int64(src.Price * 100),
		MaxNumber: int32(src.MaxQuota),
		Enabled:   src.Enabled == 1,
	}
}

func ParseLevelPrice(src *proto.SLevelPrice) *item.MemberPrice {
	return &item.MemberPrice{
		Id:       int(src.Id),
		Level:    int(src.Level),
		Price:    float32(src.Price) / 100,
		MaxQuota: int(src.MaxNumber),
		Enabled:  types.IntCond(src.Enabled, 1, 0),
	}
}

func ParseGoodsDto_(src *valueobject.Goods) *proto.SUnifiedViewItem {
	return &proto.SUnifiedViewItem{
		ItemId:      src.ItemId,
		ProductId:   src.ProductId,
		CategoryId:  int64(src.CategoryId),
		VendorId:    int64(src.VendorId),
		BrandId:     0,
		Title:       src.Title,
		Code:        "",
		SkuId:       src.SkuId,
		Image:       src.Image,
		Price:       float64(src.Price),
		PriceRange:  src.PriceRange,
		StockNum:    src.StockNum,
		ShelveState: item.ShelvesOn,
		ReviewState: 0,
		UpdateTime:  0,
	}
}

func WsSkuPriceDto(src *item.WsSkuPrice) *proto.SWsSkuPrice {
	return &proto.SWsSkuPrice{
		Id:             int64(src.ID),
		RequireNumber:  src.RequireQuantity,
		WholesalePrice: src.WholesalePrice,
	}
}

func WsSkuPrice(src *proto.SWsSkuPrice) *item.WsSkuPrice {
	return &item.WsSkuPrice{
		ID:              int32(src.Id),
		ItemId:          0,
		SkuId:           0,
		RequireQuantity: src.RequireNumber,
		WholesalePrice:  src.WholesalePrice,
	}
}

func WsItemDiscountDto(src *item.WsItemDiscount) *proto.SWsItemDiscount {
	return &proto.SWsItemDiscount{
		Id:            int64(src.ID),
		BuyerGroupId:  int64(src.BuyerGid),
		RequireAmount: int64(src.RequireAmount),
		DiscountRate:  int64(src.DiscountRate * 1000),
	}
}
func WsItemDiscount(src *proto.SWsItemDiscount) *item.WsItemDiscount {
	return &item.WsItemDiscount{
		ID:            int32(src.Id),
		RequireAmount: int32(src.RequireAmount),
		DiscountRate:  float64(src.DiscountRate) / 1000,
	}
}

func ParseGoodsItem(src *proto.SUnifiedViewItem) *item.GoodsItem {
	dst := &item.GoodsItem{
		Id:         src.ItemId,
		ProductId:  src.ProductId,
		PromFlag:   -1, //todo:??
		CategoryId: int32(src.CategoryId),
		VendorId:   src.VendorId,
		BrandId:    int32(src.BrandId),
		ShopCatId:  0,                     //todo:??
		ExpressTid: int32(src.ExpressTid), //src.,
		Title:      src.Title,
		ShortTitle: "", //src.Sho,
		Code:       src.Code,
		Image:      src.Image,
		IsPresent:  0, //todo:???
		PriceRange: src.PriceRange,
		StockNum:    src.StockNum,
		SaleNum:     0,
		SkuId:       src.SkuId,
		Cost:        0,
		Price:       0,
		RetailPrice: 0,
		SkuArray:    make([]*item.Sku, len(src.SkuArray)),
	}
	for i, v := range src.SkuArray {
		dst.SkuArray[i] = Sku(v)
	}
	return dst
}

func ParseOrderItem(v *dto.OrderItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		Id:             int64(v.Id),
		SnapshotId:     int64(v.SnapshotId),
		SkuId:          int64(v.SkuId),
		ItemId:         int64(v.ItemId),
		ItemTitle:      v.GoodsTitle,
		Image:          v.Image,
		Price:          float64(v.Price),
		FinalPrice:     float64(v.FinalPrice),
		Quantity:       int32(v.Quantity),
		ReturnQuantity: int32(v.ReturnQuantity),
		Amount:         float64(v.Amount),
		FinalAmount:    float64(v.FinalAmount),
		IsShipped:      v.IsShipped == 1,
	}
}
