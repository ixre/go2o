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
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/order"
	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types"
)

func ItemDtoV2(src item.GoodsItem) *proto.SUnifiedViewItem {
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
		Price:      src.Price,
		PriceRange: src.PriceRange,
		StockNum:   src.StockNum,
		//OriginPrice:  float64(src.OriginPrice),
		//Weight:       src.Weight,
		//Bulk:         src.Bulk,
		ShelveState: src.ReviewState,
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

func ItemDataDto(src item.GoodsItem) *proto.SItemDataResponse {
	it := &proto.SItemDataResponse{
		ItemId:     src.Id,
		ProductId:  src.ProductId,
		ItemFlag:   int32(src.ItemFlag),
		CategoryId: int64(src.CategoryId),
		SkuId:      src.SkuId,
		VendorId:   src.VendorId,
		BrandId:    int64(src.BrandId),
		Thumbnail:  src.Image,
		ShopId:     src.ShopId,
		ExpressTid: int64(src.ExpressTid),
		Title:      src.Title,
		ShortTitle: src.ShortTitle,
		IntroVideo: src.IntroVideo,
		Code:       src.Code,
		//SaleNum:      src.SaleNum,
		//SkuNum:       src.SkuNum,
		Cost:        src.Cost,
		Price:       src.Price,
		PriceRange:  src.PriceRange,
		StockNum:    src.StockNum,
		OriginPrice: src.OriginPrice,
		//Weight:       src.Weight,
		//Bulk:         src.Bulk,
		ShelveState:   src.ShelveState,
		ReviewState:   src.ReviewState,
		ReviewRemark:  src.ReviewRemark,
		SafeguardFlag: int32(src.SafeguardFlag),
		//CreateTime:   src.CreateTime,
	}
	if src.SkuArray != nil {
		it.SkuArray = make([]*proto.SSku, len(src.SkuArray))
		for i, v := range src.SkuArray {
			it.SkuArray[i] = SkuDto(v)
		}
	}
	return it
}

func AttrArrayDto(src []*product.AttrValue) []*proto.SAttr {
	var dst = make([]*proto.SAttr, len(src))
	for i, v := range src {
		dst[i] = &proto.SAttr{
			Id:       v.Id,
			AttrId:   v.AttrId,
			AttrData: v.AttrData,
			AttrWord: v.AttrWord,
		}
	}
	return dst
}

func SkuArrayDto(src []*item.Sku) []*proto.SSku {
	var dst = make([]*proto.SSku, len(src))
	for i, v := range src {
		dst[i] = SkuDto(v)
	}
	return dst
}

func AttrValueArrayDto(src []*product.AttrValue) []*proto.SAttrValue {
	var dst = make([]*proto.SAttrValue, len(src))
	for i, v := range src {
		dst[i] = &proto.SAttrValue{
			AttrId:    v.AttrId,
			AttrName:  v.AttrName,
			AttrValue: v.AttrWord,
		}
	}
	return dst
}

func SpecOptionsDto(list promodel.SpecList) []*proto.SSpecOption {
	arr := make([]*proto.SSpecOption, len(list))
	s := func(l *promodel.Spec) []*proto.SSpecOptionItem {
		arr := make([]*proto.SSpecOptionItem, len(l.Items))
		for i, v := range l.Items {
			arr[i] = &proto.SSpecOptionItem{
				Value: fmt.Sprintf("%d:%d", l.Id, v.Id),
				Label: v.Value,
				Color: v.Color,
			}
		}
		return arr
	}
	for i, v := range list {
		arr[i] = &proto.SSpecOption{
			Name:  v.Name,
			Items: s(v),
		}
	}
	return arr
}

func SkuDto(src *item.Sku) *proto.SSku {
	return &proto.SSku{
		SkuId:       src.Id,
		ProductId:   src.ProductId,
		ItemId:      src.ItemId,
		Title:       src.Title,
		Image:       src.Image,
		SpecData:    src.SpecData,
		SpecWord:    src.SpecWord,
		Code:        src.Code,
		OriginPrice: src.OriginPrice,
		Price:       src.Price,
		Cost:        src.Cost,
		Weight:      src.Weight,
		Bulk:        src.Bulk,
		Stock:       src.Stock,
		SaleNum:     src.SaleNum,
	}
}

func Sku(src *proto.SSku) *item.Sku {
	return &item.Sku{
		Id:          src.SkuId,
		ProductId:   src.ProductId,
		ItemId:      src.ItemId,
		Title:       src.Title,
		Image:       src.Image,
		SpecData:    src.SpecData,
		SpecWord:    src.SpecWord,
		Code:        src.Code,
		OriginPrice: src.OriginPrice,
		Price:       src.Price,
		Cost:        src.Cost,
		Weight:      src.Weight,
		Bulk:        src.Bulk,
		Stock:       src.Stock,
		SaleNum:     src.SaleNum,
	}
}

func OrderItem(src *proto.SOrderItem) *order.ComplexItem {
	return &order.ComplexItem{
		ID:             src.Id,
		ItemId:         src.ItemId,
		SkuId:          src.SkuId,
		SnapshotId:     src.SnapshotId,
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		Amount:         src.ItemAmount,
		FinalAmount:    src.FinalAmount,
		IsShipped:      int32(types.ElseInt(src.IsShipped, 1, 0)),
		Data:           src.Data,
	}
}

func SubOrderItemDto(src *order.SubOrderItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		Id:             src.ID,
		SnapshotId:     src.SnapshotId,
		SkuId:          src.SkuId,
		ItemId:         src.ItemId,
		ItemTitle:      "",
		Image:          "",
		Price:          0,
		FinalPrice:     0,
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		ItemAmount:     src.Amount,
		FinalAmount:    src.FinalAmount,
		IsShipped:      src.IsShipped == 1,
		Data:           nil,
	}
}

func OrderItemDto(src *order.ComplexItem) *proto.SOrderItem {
	return &proto.SOrderItem{
		SnapshotId:     src.SnapshotId,
		SkuId:          src.SkuId,
		ItemId:         src.ItemId,
		ItemTitle:      src.ItemTitle,
		Image:          src.MainImage,
		Price:          src.Price,
		FinalPrice:     src.FinalPrice,
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		ItemAmount:     src.Amount,
		FinalAmount:    src.FinalAmount,
		IsShipped:      src.IsShipped == 1,
		Data:           src.Data,
	}
}

func OrderDto(src *order.ComplexOrder) *proto.SSingleOrder {
	d := src.Details[0]
	o := &proto.SSingleOrder{
		OrderId:          d.Id,
		OrderType:        src.OrderType,
		OrderNo:          d.OrderNo,
		BuyerId:          src.BuyerId,
		SellerId:         src.SellerId,
		SellerName:       "",
		Items:            make([]*proto.SOrderItem, len(d.Items)),
		ItemAmount:       int32(src.ItemAmount),
		DiscountAmount:   int32(src.DiscountAmount),
		ExpressFee:       int32(src.ExpressFee),
		PackageFee:       int32(src.PackageFee),
		FinalAmount:      int32(src.FinalAmount),
		BuyerComment:     d.BuyerComment,
		Consignee:        &proto.SConsigneeInfo{ConsigneeName: src.Consignee.ConsigneeName, ConsigneePhone: src.Consignee.ConsigneePhone, ShippingAddress: src.Consignee.ShippingAddress, IsModified: src.ConsigneeModified == 1},
		SubmitTime:       src.CreateTime,
		Status:           int32(d.Status),
		StateDesc:        "",
		ExpiresTime:      0,
		PaymentTime:      0,
		TradeData:        []*proto.SOrderPayChanData{},
		ShipTime:         0,
		ShipLogisticCode: "",
		ShipExpressName:  "",
		Data:             src.Data,
	}
	if d.Items != nil {
		for i, v := range d.Items {
			o.Items[i] = OrderItemDto(v)
		}
	}
	return o
}

func ParentOrderDto(src *order.ComplexOrder) *proto.SParentOrder {
	o := &proto.SParentOrder{
		BuyerId:        src.BuyerId,
		BuyerUser:      src.BuyerUser,
		OrderNo:        src.OrderNo,
		ItemCount:      int64(src.ItemCount),
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		DeductAmount:   0,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		FinalAmount:    src.FinalAmount,
		Consignee: &proto.SConsigneeInfo{
			ConsigneeName:   src.Consignee.ConsigneeName,
			ConsigneePhone:  src.Consignee.ConsigneePhone,
			ShippingAddress: src.Consignee.ShippingAddress,
		},
		SubOrders:  []*proto.SSubOrder{},
		Status:     int32(src.Status),
		IsPaid:     src.IsPaid == 1,
		StatusText: order.OrderStatus(src.Status).String(),
		CreateTime: src.CreateTime,
	}
	for _, v := range src.Details {
		d := &proto.SSubOrder{
			OrderNo:        v.OrderNo,
			ShopId:         v.ShopId,
			ShopName:       v.ShopName,
			ItemAmount:     v.ItemAmount,
			DiscountAmount: v.DiscountAmount,
			DeductAmount:   0, //v.DeductAmount,
			AdjustAmount:   0, //v.AdjustAmount,
			ExpressFee:     v.ExpressFee,
			PackageFee:     v.PackageFee,
			ProcedureFee:   0, //v.ProcedureFee,
			TotalAmount:    0,
			FinalAmount:    v.FinalAmount,
			BuyerComment:   v.BuyerComment,
			PaymentTime:    v.PaymentTime,
			CloseTime:      v.CloseTime,
			Status:         int32(v.Status),
			Items:          []*proto.SOrderItem{},
			StatusText:     order.OrderStatus(v.Status).String(),
		}
		for _, it := range v.Items {
			d.Items = append(d.Items, OrderItemDto(it))
		}
		o.SubOrders = append(o.SubOrders, d)
	}
	return o
}

// PrepareOrderDto 转换为预生成订单
func PrepareOrderDto(src *order.ComplexOrder) *proto.PrepareOrderResponse {
	o := &proto.PrepareOrderResponse{
		OrderType:      src.OrderType,
		ItemAmount:     src.ItemAmount,
		DiscountAmount: src.DiscountAmount,
		ExpressFee:     src.ExpressFee,
		PackageFee:     src.PackageFee,
		FinalAmount:    src.FinalAmount,
		Consignee: &proto.SConsigneeInfo{
			ConsigneeName:   src.Consignee.ConsigneeName,
			ConsigneePhone:  src.Consignee.ConsigneePhone,
			ShippingAddress: src.Consignee.ShippingAddress,
		},
	}
	return o
}
