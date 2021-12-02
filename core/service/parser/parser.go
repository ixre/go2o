/**
 * Copyright 2015 @ at3.net.
 * name : parser.go
 * author : jarryliu
 * date : 2016-11-17 15:07
 * description :
 * history :
 */
package parser

import (
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/types"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/service/proto"
)

func MemberDto(src *member.Member) *proto.SMember {
	return &proto.SMember{
		Id:             src.Id,
		User:           src.User,
		Code:           src.Code,
		Pwd:            src.Pwd,
		TradePwd:       src.TradePwd,
		Exp:            int64(src.Exp),
		Level:          int32(src.Level),
		PremiumUser:    int32(src.PremiumUser),
		PremiumExpires: src.PremiumExpires,
		InviteCode:     src.InviteCode,
		RegIp:          src.RegIp,
		RegFrom:        src.RegFrom,
		State:          int32(src.State),
		Flag:           int32(src.Flag),
		Avatar:         src.Avatar,
		Phone:          src.Phone,
		Email:          src.Email,
		Name:           src.Name,
		RealName:       src.RealName,
		DynamicToken:   src.DynamicToken,
		RegTime:        src.RegTime,
		LastLoginTime:  src.LastLoginTime,
	}
}

func ComplexMemberDto(src *member.ComplexMember) *proto.SComplexMember {
	return &proto.SComplexMember{
		Name:           src.Name,
		Avatar:         src.Avatar,
		Exp:            int32(src.Exp),
		Level:          int32(src.Level),
		LevelName:      src.LevelName,
		PremiumUser:    int32(src.PremiumUser),
		InviteCode:     src.InviteCode,
		TrustAuthState: int32(src.TrustAuthState),
		TradePwdHasSet: src.TradePwdHasSet,
		UpdateTime:     src.UpdateTime,
	}
}

func round(f float32, n int) float64 {
	return math.Round(float64(f), n)
}

func Address(src *proto.SAddress) *member.ConsigneeAddress {
	return &member.ConsigneeAddress{
		Id:             src.Id,
		ConsigneeName:  src.ConsigneeName,
		ConsigneePhone: src.ConsigneePhone,
		Province:       src.Province,
		City:           src.City,
		District:       src.District,
		Area:           src.Area,
		DetailAddress:  src.DetailAddress,
		IsDefault:      int(src.IsDefault),
	}
}

func TradeMethodDataDto(src *payment.TradeMethodData) *proto.STradeMethodData {
	return &proto.STradeMethodData{
		Method:     int32(src.Method),
		Code:       src.Code,
		Amount:     int32(src.Amount),
		Internal:   int32(src.Internal),
		OutTradeNo: src.OutTradeNo,
		PayTime:    src.PayTime,
	}
}

func ItemDto(src *item.GoodsItem) *proto.SOldItem {
	it := &proto.SOldItem{
		ItemId:       src.Id,
		ProductId:    src.ProductId,
		PromFlag:     src.PromFlag,
		CatId:        src.CategoryId,
		VendorId:     src.VendorId,
		BrandId:      src.BrandId,
		ShopId:       src.ShopId,
		ShopCatId:    src.ShopCatId,
		ExpressTid:   src.ExpressTid,
		Title:        src.Title,
		ShortTitle:   src.ShortTitle,
		Code:         src.Code,
		Image:        src.Image,
		IsPresent:    src.IsPresent,
		PriceRange:   src.PriceRange,
		StockNum:     src.StockNum,
		SaleNum:      src.SaleNum,
		SkuNum:       src.SkuNum,
		SkuId:        src.SkuId,
		Cost:         float64(src.Cost),
		Price:        float64(src.Price),
		RetailPrice:  float64(src.RetailPrice),
		Weight:       src.Weight,
		Bulk:         src.Bulk,
		ShelveState:  src.ShelveState,
		ReviewState:  src.ReviewState,
		ReviewRemark: src.ReviewRemark,
		SortNum:      src.SortNum,
		CreateTime:   src.CreateTime,
		UpdateTime:   src.UpdateTime,
		PromPrice:    float64(src.PromPrice),
	}
	if src.SkuArray != nil {
		it.SkuArray = make([]*proto.SSku, len(src.SkuArray))
		for i, v := range src.SkuArray {
			it.SkuArray[i] = SkuDto(v)
		}
	}
	return it
}

func ItemDtoV2(src *item.GoodsItem) *proto.SUnifiedViewItem {
	it := &proto.SUnifiedViewItem{
		ItemId:    src.Id,
		ProductId: src.ProductId,
		//PromFlag:     src.PromFlag,
		CategoryId: int64(src.CategoryId),
		SkuId:      src.SkuId,
		VendorId:   src.VendorId,
		BrandId:    int64(src.BrandId),
		//ShopId:       src.ShopId,
		//ShopCatId:    src.ShopCatId,
		ExpressTid: int64(src.ExpressTid),
		Title:      src.Title,
		//ShortTitle:   src.ShortTitle,
		Code:  src.Code,
		Image: src.Image,
		//SaleNum:      src.SaleNum,
		//SkuNum:       src.SkuNum,
		//SkuId:        src.SkuId,
		//Cost:         float64(src.Cost),
		Price: float64(src.Price),
		//IsPresent:    src.IsPresent,
		PriceRange: src.PriceRange,
		StockNum:   src.StockNum,
		//RetailPrice:  float64(src.RetailPrice),
		//Weight:       src.Weight,
		//Bulk:         src.Bulk,
		ShelveState: src.ShelveState,
		ReviewState: src.ReviewState,
		//ReviewRemark: src.ReviewRemark,
		//SortNum:      src.SortNum,
		//CreateTime:   src.CreateTime,
		UpdateTime: src.UpdateTime,
		Data:       nil,
	}
	if src.SkuArray != nil {
		it.SkuArray = make([]*proto.SSku, len(src.SkuArray))
		for i, v := range src.SkuArray {
			it.SkuArray[i] = SkuDto(v)
		}
	}
	return it
}

func SkuArrayDto(src []*item.Sku) []*proto.SSku {
	var dst = make([]*proto.SSku, len(src))
	for i, v := range src {
		dst[i] = SkuDto(v)
	}
	return dst
}

func Item(src *proto.SOldItem) *item.GoodsItem {
	it := &item.GoodsItem{
		Id:           src.ItemId,
		ProductId:    src.ProductId,
		PromFlag:     src.PromFlag,
		CategoryId:   src.CatId,
		VendorId:     src.VendorId,
		BrandId:      src.BrandId,
		ShopId:       src.ShopId,
		ShopCatId:    src.ShopCatId,
		ExpressTid:   src.ExpressTid,
		Title:        src.Title,
		ShortTitle:   src.ShortTitle,
		Code:         src.Code,
		Image:        src.Image,
		IsPresent:    src.IsPresent,
		PriceRange:   src.PriceRange,
		StockNum:     src.StockNum,
		SaleNum:      src.SaleNum,
		SkuNum:       src.SkuNum,
		SkuId:        src.SkuId,
		Cost:         float32(src.Cost),
		Price:        float32(src.Price),
		RetailPrice:  float32(src.RetailPrice),
		Weight:       src.Weight,
		Bulk:         src.Bulk,
		ShelveState:  src.ShelveState,
		ReviewState:  src.ReviewState,
		ReviewRemark: src.ReviewRemark,
		SortNum:      src.SortNum,
		CreateTime:   src.CreateTime,
		UpdateTime:   src.UpdateTime,
		PromPrice:    float32(src.PromPrice),
	}
	if src.SkuArray != nil {
		it.SkuArray = make([]*item.Sku, len(src.SkuArray))
		for i, v := range src.SkuArray {
			it.SkuArray[i] = Sku(v)
		}
	}
	return it
}

func SkuDto(src *item.Sku) *proto.SSku {
	return &proto.SSku{
		SkuId:       src.ID,
		ProductId:   src.ProductId,
		ItemId:      src.ItemId,
		Title:       src.Title,
		Image:       src.Image,
		SpecData:    src.SpecData,
		SpecWord:    src.SpecWord,
		Code:        src.Code,
		RetailPrice: math.Round(float64(src.RetailPrice), 2),
		Price:       math.Round(float64(src.Price), 2),
		Cost:        math.Round(float64(src.Cost), 2),
		Weight:      src.Weight,
		Bulk:        src.Bulk,
		Stock:       src.Stock,
		SaleNum:     src.SaleNum,
	}
}

func Sku(src *proto.SSku) *item.Sku {
	return &item.Sku{
		ID:          src.SkuId,
		ProductId:   src.ProductId,
		ItemId:      src.ItemId,
		Title:       src.Title,
		Image:       src.Image,
		SpecData:    src.SpecData,
		SpecWord:    src.SpecWord,
		Code:        src.Code,
		RetailPrice: float32(src.RetailPrice),
		Price:       float32(src.Price),
		Cost:        float32(src.Cost),
		Weight:      src.Weight,
		Bulk:        src.Bulk,
		Stock:       src.Stock,
		SaleNum:     src.SaleNum,
	}
}

func Order(src *proto.SSingleOrder) *order.ComplexOrder {
	o := &order.ComplexOrder{
		OrderId:        src.OrderId,
		OrderType:      src.OrderType,
		OrderNo:        src.OrderNo,
		BuyerId:        src.BuyerId,
		VendorId:       src.SellerId,
		ShopId:         src.ShopId,
		Subject:        src.Subject,
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		FinalAmount:    src.FinalAmount,

		Consignee: &order.ComplexConsignee{
			ConsigneePerson: src.Consignee.ConsigneePerson,
			ConsigneePhone:  src.Consignee.ConsigneePhone,
			ShippingAddress: src.Consignee.ShippingAddress,
		},
		BuyerComment: src.BuyerComment,
		CreateTime:   src.SubmitTime,
		State:        src.State,
		Items:        make([]*order.ComplexItem, len(src.Items)),
		Data:         src.Data,
	}
	if src.Items != nil {
		for i, v := range src.Items {
			o.Items[i] = OrderItem(v)
		}
	}
	return o
}

func OrderItem(src *proto.SOrderItem) *order.ComplexItem {
	return &order.ComplexItem{
		ID:             src.Id,
		ItemId:         src.ItemId,
		SkuId:          src.SkuId,
		SnapshotId:     src.SnapshotId,
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		Amount:         src.Amount,
		FinalAmount:    src.FinalAmount,
		IsShipped:      int32(types.ElseInt(src.IsShipped, 1, 0)),
		Data:           src.Data,
	}
}

func SubOrderItemDto(src *order.SubOrderItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		Id:                   src.ID,
		SnapshotId:           src.SnapshotId,
		SkuId:                src.SkuId,
		ItemId:               src.ItemId,
		ItemTitle:            "", //todo: 没有显示图片
		Image:                "",
		Price:                0,
		FinalPrice:           0,
		Quantity:             src.Quantity,
		ReturnQuantity:       src.ReturnQuantity,
		Amount:               float64(src.Amount),
		FinalAmount:          float64(src.FinalAmount),
		IsShipped:            src.IsShipped == 1,
		Data:                 nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

func SubOrderDto(src *order.NormalSubOrder) *proto.SSingleOrder {
	o := &proto.SSingleOrder{
		OrderId:       src.ID,
		ParentOrderId: src.OrderId,
		OrderNo:       src.OrderNo,
		//OrderType:            src.,
		BuyerId:        src.BuyerId,
		SellerId:       src.VendorId,
		ShopId:         src.ShopId,
		Subject:        src.Subject,
		ItemAmount:     float64(src.ItemAmount),
		DiscountAmount: float64(src.DiscountAmount),
		DeductAmount:   0,
		AdjustAmount:   0,
		ExpressFee:     float64(src.ExpressFee),
		PackageFee:     float64(src.PackageFee),
		ProcedureFee:   0,
		//TotalAmount:          src.ItemAmount,
		FinalAmount: float64(src.FinalAmount),
		Consignee: &proto.SConsigneeInfo{
			ConsigneePerson: "",
			ConsigneePhone:  "",
			ShippingAddress: "",
		},
		BuyerComment: src.BuyerComment,
		State:        src.State,
		SubmitTime:   src.CreateTime,
		Items:        make([]*proto.SOrderItem, len(src.Items)),
	}
	for i, v := range src.Items {
		o.Items[i] = SubOrderItemDto(v)
	}
	return o
}

func OrderItemDto(src *order.ComplexItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		SnapshotId: src.SnapshotId,
		SkuId:      src.SkuId,
		ItemId:     src.ItemId,
		ItemTitle:  src.ItemTitle,
		Image:      src.MainImage,
		Price:      float64(src.Price),
		//FinalPrice:           0,
		Quantity:             src.Quantity,
		ReturnQuantity:       src.ReturnQuantity,
		Amount:               src.Amount,
		FinalAmount:          src.FinalAmount,
		IsShipped:            src.IsShipped == 1,
		Data:                 src.Data,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

func OrderDto(src *order.ComplexOrder) *proto.SSingleOrder {
	o := &proto.SSingleOrder{
		OrderId: src.OrderId,
		//SubOrderId:      src.SubOrderId,
		OrderType:      src.OrderType,
		OrderNo:        src.OrderNo,
		BuyerId:        src.BuyerId,
		SellerId:       src.VendorId,
		ShopId:         src.ShopId,
		Subject:        src.Subject,
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		FinalAmount:    src.FinalAmount,
		Consignee: &proto.SConsigneeInfo{
			ConsigneePerson: src.Consignee.ConsigneePerson,
			ConsigneePhone:  src.Consignee.ConsigneePhone,
			ShippingAddress: src.Consignee.ShippingAddress,
		},
		BuyerComment: src.BuyerComment,
		SubmitTime:   src.CreateTime,
		State:        src.State,
		Items:        make([]*proto.SOrderItem, len(src.Items)),
		Data:         src.Data,
	}
	if src.Items != nil {
		for i, v := range src.Items {
			o.Items[i] = OrderItemDto(v)
		}
	}
	return o
}
