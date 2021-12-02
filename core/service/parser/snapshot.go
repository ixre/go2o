package parser

import (
	"go2o/core/domain/interface/item"
	"go2o/core/service/proto"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : snapshort
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-19 16:36
 * description :
 * history :
 */

func ParseItemSnapshotDto(src *item.Snapshot) *proto.SItemSnapshot {
	return &proto.SItemSnapshot{
		ItemId:      src.ItemId,
		ProductId:   src.ProductId,
		Key:         src.Key,
		CategoryId:  int64(src.CatId),
		VendorId:    src.VendorId,
		BrandId:     int64(src.BrandId),
		ShopId:      src.ShopId,
		ShopCatId:   int64(src.ShopCatId),
		ExpressTid:  int64(src.ExpressTid),
		Title:       src.Title,
		ShortTitle:  src.ShortTitle,
		ProductCode: src.Code,
		MainImage:   src.Image,
		IsPresent:   src.IsPresent == 1,
		PriceRange:  src.PriceRange,
		SkuId:       src.SkuId,
		Stock:       &proto.SItemStock{},
		Price:       src.Price,
		RetailPrice: src.RetailPrice,
		UpdateTime:  src.UpdateTime,
		SkuArray:    nil,
		ViewData:    nil,
	}
}

func ParseTradeSnapshot(src *item.TradeSnapshot) *proto.STradeSnapshot {
	return &proto.STradeSnapshot{
		Id:          src.Id,
		ItemId:      src.ItemId,
		SkuId:       src.SkuId,
		SnapshotKey: src.SnapshotKey,
		SellerId:    src.SellerId,
		SellerName:  "",
		Title:       src.GoodsTitle,
		ProductCode: src.GoodsNo,
		Sku:         src.Sku,
		Image:       src.Image,
		Price:       src.Price,
		CreateTime:  src.CreateTime,
	}
}
