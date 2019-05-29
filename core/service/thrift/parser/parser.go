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
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/service/auto_gen/rpc/foundation_service"
	"go2o/core/service/auto_gen/rpc/mch_service"
	"go2o/core/service/auto_gen/rpc/member_service"
	"go2o/core/service/auto_gen/rpc/order_service"
	"go2o/core/service/auto_gen/rpc/payment_service"
	"go2o/core/service/auto_gen/rpc/ttype"
)

func MerchantDto(src *merchant.ComplexMerchant) *mch_service.SComplexMerchant {
	return &mch_service.SComplexMerchant{
		ID:            src.Id,
		MemberId:      src.MemberId,
		Usr:           src.Usr,
		Pwd:           src.Pwd,
		Name:          src.Name,
		SelfSales:     src.SelfSales,
		Level:         src.Level,
		Logo:          src.Logo,
		CompanyName:   src.CompanyName,
		Province:      src.Province,
		City:          src.City,
		District:      src.District,
		Enabled:       src.Enabled,
		ExpiresTime:   src.ExpiresTime,
		JoinTime:      src.JoinTime,
		UpdateTime:    src.UpdateTime,
		LoginTime:     src.LoginTime,
		LastLoginTime: src.LastLoginTime,
	}
}

func LevelDto(src *member.Level) *member_service.SLevel {
	return &member_service.SLevel{
		ID:            int32(src.ID),
		Name:          src.Name,
		RequireExp:    int32(src.RequireExp),
		ProgramSignal: src.ProgramSignal,
		Enabled:       int32(src.Enabled),
		IsOfficial:    int32(src.IsOfficial),
	}
}

func MemberDto(src *member.Member) *member_service.SMember {
	return &member_service.SMember{
		ID:             src.Id,
		Usr:            src.Usr,
		Pwd:            src.Pwd,
		TradePwd:       src.TradePwd,
		Exp:            int32(src.Exp),
		Level:          int32(src.Level),
		InvitationCode: src.InvitationCode,
		PremiumUser:    int32(src.PremiumUser),
		PremiumExpires: src.PremiumExpires,
		RegFrom:        src.RegFrom,
		RegIp:          src.RegIp,
		RegTime:        src.RegTime,
		CheckCode:      src.CheckCode,
		CheckExpires:   src.CheckExpires,
		Flag:           int32(src.Flag),
		State:          int32(src.State),
		LoginTime:      src.LoginTime,
		LastLoginTime:  src.LastLoginTime,
		UpdateTime:     src.UpdateTime,
		DynamicToken:   src.DynamicToken,
		TimeoutTime:    src.TimeoutTime,
	}
}

func Member(src *member_service.SMember) *member.Member {
	return &member.Member{
		Id:             src.ID,
		Usr:            src.Usr,
		Pwd:            src.Pwd,
		TradePwd:       src.TradePwd,
		Exp:            int(src.Exp),
		Level:          int(src.Level),
		InvitationCode: src.InvitationCode,
		PremiumUser:    int(src.PremiumUser),
		PremiumExpires: src.PremiumExpires,
		RegFrom:        src.RegFrom,
		RegIp:          src.RegIp,
		RegTime:        src.RegTime,
		CheckCode:      src.CheckCode,
		CheckExpires:   src.CheckExpires,
		Flag:           int(src.Flag),
		State:          int(src.State),
		LoginTime:      src.LoginTime,
		LastLoginTime:  src.LastLoginTime,
		UpdateTime:     src.UpdateTime,
		DynamicToken:   src.DynamicToken,
		TimeoutTime:    src.TimeoutTime,
	}
}

func MemberProfile(src *member.Profile) *member_service.SProfile {
	return &member_service.SProfile{
		MemberId:   src.MemberId,
		Name:       src.Name,
		Avatar:     src.Avatar,
		Sex:        src.Sex,
		BirthDay:   src.BirthDay,
		Phone:      src.Phone,
		Address:    src.Address,
		Im:         src.Im,
		Email:      src.Email,
		Province:   src.Province,
		City:       src.City,
		District:   src.District,
		Remark:     src.Remark,
		Ext1:       src.Ext1,
		Ext2:       src.Ext2,
		Ext3:       src.Ext3,
		Ext4:       src.Ext4,
		Ext5:       src.Ext5,
		Ext6:       src.Ext6,
		UpdateTime: src.UpdateTime,
	}
}

func MemberProfile2(src *member_service.SProfile) *member.Profile {
	return &member.Profile{
		MemberId:   src.MemberId,
		Name:       src.Name,
		Avatar:     src.Avatar,
		Sex:        src.Sex,
		BirthDay:   src.BirthDay,
		Phone:      src.Phone,
		Address:    src.Address,
		Im:         src.Im,
		Email:      src.Email,
		Province:   src.Province,
		City:       src.City,
		District:   src.District,
		Remark:     src.Remark,
		Ext1:       src.Ext1,
		Ext2:       src.Ext2,
		Ext3:       src.Ext3,
		Ext4:       src.Ext4,
		Ext5:       src.Ext5,
		Ext6:       src.Ext6,
		UpdateTime: src.UpdateTime,
	}
}

func ComplexMemberDto(src *member.ComplexMember) *member_service.SComplexMember {
	return &member_service.SComplexMember{
		MemberId:          src.MemberId,
		Usr:               src.Usr,
		Name:              src.Name,
		Avatar:            src.Avatar,
		Exp:               int32(src.Exp),
		Level:             int32(src.Level),
		LevelName:         src.LevelName,
		LevelSign:         src.LevelSign,
		LevelOfficial:     int32(src.LevelOfficial),
		PremiumUser:       int32(src.PremiumUser),
		PremiumExpires:    src.PremiumExpires,
		InvitationCode:    src.InvitationCode,
		TrustAuthState:    int32(src.TrustAuthState),
		State:             int32(src.State),
		Integral:          int64(src.Integral),
		Balance:           src.Balance,
		WalletBalance:     src.WalletBalance,
		GrowBalance:       src.GrowBalance,
		GrowAmount:        src.GrowAmount,
		GrowEarnings:      src.GrowEarnings,
		GrowTotalEarnings: src.GrowTotalEarnings,
		UpdateTime:        src.UpdateTime,
	}
}

func round(f float32, n int) float64 {
	return math.Round(float64(f), n)
}

func AccountDto(src *member.Account) *member_service.SAccount {
	return &member_service.SAccount{
		MemberId:          src.MemberId,
		Integral:          int64(src.Integral),
		FreezeIntegral:    int64(src.FreezeIntegral),
		Balance:           round(src.Balance, 2),
		FreezeBalance:     round(src.FreezeBalance, 2),
		ExpiredBalance:    round(src.ExpiredBalance, 2),
		WalletBalance:     round(src.WalletBalance, 2),
		FreezeWallet:      round(src.FreezeWallet, 2),
		ExpiredPresent:    round(src.ExpiredPresent, 2),
		TotalPresentFee:   round(src.TotalPresentFee, 2),
		FlowBalance:       round(src.FlowBalance, 2),
		GrowBalance:       round(src.GrowBalance, 2),
		GrowAmount:        round(src.GrowAmount, 2),
		GrowEarnings:      round(src.GrowEarnings, 2),
		GrowTotalEarnings: round(src.GrowTotalEarnings, 2),
		TotalExpense:      round(src.TotalExpense, 2),
		TotalCharge:       round(src.TotalCharge, 2),
		TotalPay:          round(src.TotalPay, 2),
		PriorityPay:       int64(src.PriorityPay),
		UpdateTime:        src.UpdateTime,
	}
}

func Account(src *member_service.SAccount) *member.Account {
	return &member.Account{
		MemberId:          src.MemberId,
		Integral:          int(src.Integral),
		FreezeIntegral:    int(src.FreezeIntegral),
		Balance:           float32(src.Balance),
		FreezeBalance:     float32(src.FreezeBalance),
		ExpiredBalance:    float32(src.ExpiredBalance),
		WalletBalance:     float32(src.WalletBalance),
		FreezeWallet:      float32(src.FreezeWallet),
		ExpiredPresent:    float32(src.ExpiredPresent),
		TotalPresentFee:   float32(src.TotalPresentFee),
		FlowBalance:       float32(src.FlowBalance),
		GrowBalance:       float32(src.GrowBalance),
		GrowAmount:        float32(src.GrowAmount),
		GrowEarnings:      float32(src.GrowEarnings),
		GrowTotalEarnings: float32(src.GrowTotalEarnings),
		TotalExpense:      float32(src.TotalExpense),
		TotalCharge:       float32(src.TotalCharge),
		TotalPay:          float32(src.TotalPay),
		PriorityPay:       int(src.PriorityPay),
		UpdateTime:        src.UpdateTime,
	}
}

func PlatformConfDto(src *valueobject.PlatformConf) *foundation_service.PlatformConf {
	return &foundation_service.PlatformConf{
		Suspend:          src.Suspend,
		SuspendMessage:   src.SuspendMessage,
		MchGoodsCategory: src.MchGoodsCategory,
		MchPageCategory:  src.MchPageCategory,
	}
}

func PlatFromConf(src *foundation_service.PlatformConf) *valueobject.PlatformConf {
	return &valueobject.PlatformConf{
		Suspend:          src.Suspend,
		SuspendMessage:   src.SuspendMessage,
		MchGoodsCategory: src.MchGoodsCategory,
		MchPageCategory:  src.MchPageCategory,
	}
}

func AddressDto(src *member.Address) *member_service.SAddress {
	return &member_service.SAddress{
		ID:        src.ID,
		MemberId:  src.MemberId,
		RealName:  src.RealName,
		Phone:     src.Phone,
		Province:  src.Province,
		City:      src.City,
		District:  src.District,
		Area:      src.Area,
		Address:   src.Address,
		IsDefault: int32(src.IsDefault),
	}
}

func Address(src *member_service.SAddress) *member.Address {
	return &member.Address{
		ID:        src.ID,
		MemberId:  src.MemberId,
		RealName:  src.RealName,
		Phone:     src.Phone,
		Province:  src.Province,
		City:      src.City,
		District:  src.District,
		Area:      src.Area,
		Address:   src.Address,
		IsDefault: int(src.IsDefault),
	}
}

func PaymentOrder(src *payment_service.SPaymentOrder) *payment.Order {
	dst := &payment.Order{
		ID:             int(src.ID),
		SellerId:       int(src.SellerId),
		TradeType:      src.TradeType,
		TradeNo:        src.TradeNo,
		OrderType:      int(src.OrderType),
		OutOrderNo:     src.OutOrderNo,
		Subject:        src.Subject,
		BuyerId:        int64(src.BuyerId),
		PayUid:         int64(src.PayUid),
		TotalAmount:    int(src.TotalAmount),
		DiscountAmount: int(src.DiscountAmount),
		DeductAmount:   int(src.DeductAmount),
		AdjustAmount:   int(src.AdjustAmount),
		ItemAmount:     int(src.ItemAmount),
		ProcedureFee:   int(src.ProcedureFee),
		FinalFee:       int(src.FinalFee),
		PaidFee:        int(src.PaidFee),
		PayFlag:        int(src.PayFlag),
		FinalFlag:      int(src.FinalFlag),
		ExtraData:      src.ExtraData,
		State:          int(src.State),
		SubmitTime:     src.SubmitTime,
		ExpiresTime:    src.ExpiresTime,
		PaidTime:       src.PaidTime,
		TradeMethods:   make([]*payment.TradeMethodData, 0),
	}
	if src.SubOrder {
		dst.SubOrder = 1
	}
	return dst
}

func PaymentOrderDto(src *payment.Order) *payment_service.SPaymentOrder {
	return &payment_service.SPaymentOrder{
		ID:             int32(src.ID),
		SellerId:       int32(src.SellerId),
		TradeType:      src.TradeType,
		TradeNo:        src.TradeNo,
		Subject:        src.Subject,
		BuyerId:        int32(src.BuyerId),
		PayUid:         int32(src.PayUid),
		TotalAmount:    int32(src.TotalAmount),
		DiscountAmount: int32(src.DiscountAmount),
		DeductAmount:   int32(src.DeductAmount),
		AdjustAmount:   int32(src.AdjustAmount),
		ItemAmount:     int32(src.ItemAmount),
		ProcedureFee:   int32(src.ProcedureFee),
		FinalFee:       int32(src.FinalFee),
		PaidFee:        int32(src.PaidFee),
		PayFlag:        int32(src.PayFlag),
		FinalFlag:      int32(src.FinalFlag),
		ExtraData:      src.ExtraData,
		State:          int32(src.State),
		SubmitTime:     int64(src.SubmitTime),
		ExpiresTime:    int64(src.ExpiresTime),
		PaidTime:       int64(src.PaidTime),
		SubOrder:       src.SubOrder == 1,
		OrderType:      int32(src.OrderType),
		OutOrderNo:     src.OutOrderNo,
		TradeData:      make([]*payment_service.STradeMethodData, 0),
	}
}

func TradeMethodDataDto(src *payment.TradeMethodData) *payment_service.STradeMethodData {
	return &payment_service.STradeMethodData{
		Method:     int32(src.Method),
		Code:       src.Code,
		Amount:     int32(src.Amount),
		Internal:   int32(src.Internal),
		OutTradeNo: src.OutTradeNo,
		PayTime:    src.PayTime,
	}
}

func MemberRelationDto(src *member.Relation) *member_service.SMemberRelation {
	return &member_service.SMemberRelation{
		MemberId:      src.MemberId,
		CardId:        src.CardCard,
		InviterId:     src.InviterId,
		InviterStr:    src.InviterStr,
		RegisterMchId: src.RegMchId,
	}
}

func TrustedInfoDto(src *member.TrustedInfo) *member_service.STrustedInfo {
	return &member_service.STrustedInfo{
		MemberId:    src.MemberId,
		RealName:    src.RealName,
		CardId:      src.CardId,
		TrustImage:  src.TrustImage,
		ReviewState: int32(src.ReviewState),
		ReviewTime:  src.ReviewTime,
		Remark:      src.Remark,
		UpdateTime:  src.UpdateTime,
	}
}

func ItemDto(src *item.GoodsItem) *ttype.SOldItem {
	it := &ttype.SOldItem{
		ItemId:       src.ID,
		ProductId:    src.ProductId,
		PromFlag:     src.PromFlag,
		CatId:        src.CatId,
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
		it.SkuArray = make([]*ttype.SSku, len(src.SkuArray))
		for i, v := range src.SkuArray {
			it.SkuArray[i] = SkuDto(v)
		}
	}
	return it
}

func Item(src *ttype.SOldItem) *item.GoodsItem {
	it := &item.GoodsItem{
		ID:           src.ItemId,
		ProductId:    src.ProductId,
		PromFlag:     src.PromFlag,
		CatId:        src.CatId,
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

func SkuDto(src *item.Sku) *ttype.SSku {
	return &ttype.SSku{
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

func Sku(src *ttype.SSku) *item.Sku {
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

func ShoppingCartItem(src *ttype.SShoppingCartItem) *cart.NormalCartItem {
	i := &cart.NormalCartItem{
		ItemId:   src.ItemId,
		SkuId:    src.SkuId,
		Quantity: src.Quantity,
		Checked:  util.BoolExt.TInt32(src.Checked, 1, 0),
		ShopId:   src.ShopId,
	}
	return i
}

func SubOrderItemDto(src *order.SubOrderItem) *order_service.SComplexItem {
	return &order_service.SComplexItem{
		ID:             int64(src.ID),
		OrderId:        src.OrderId,
		ItemId:         int64(src.ItemId),
		SkuId:          int64(src.SkuId),
		SnapshotId:     int64(src.SnapshotId),
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		Amount:         float64(src.Amount),
		FinalAmount:    float64(src.FinalAmount),
		IsShipped:      int32(src.IsShipped),
	}
}

func SubOrderDto(src *order.NormalSubOrder) *order_service.SComplexOrder {
	o := &order_service.SComplexOrder{
		OrderId:        src.OrderId,
		SubOrderId:     src.OrderId,
		OrderNo:        src.OrderNo,
		BuyerId:        int64(src.BuyerId),
		VendorId:       src.VendorId,
		ShopId:         src.ShopId,
		Subject:        src.Subject,
		ItemAmount:     float64(src.ItemAmount),
		DiscountAmount: float64(src.DiscountAmount),
		ExpressFee:     float64(src.ExpressFee),
		PackageFee:     float64(src.PackageFee),
		FinalAmount:    float64(src.FinalAmount),
		CreateTime:     src.CreateTime,
		UpdateTime:     src.UpdateTime,
		State:          int32(src.State),
		Items:          make([]*order_service.SComplexItem, len(src.Items)),
	}
	for i, v := range src.Items {
		o.Items[i] = SubOrderItemDto(v)
	}
	return o
}

func OrderItemDto(src *order.ComplexItem) *order_service.SComplexItem {
	return &order_service.SComplexItem{
		ID:             src.ID,
		OrderId:        src.OrderId,
		ItemId:         src.ItemId,
		SkuId:          src.SkuId,
		SnapshotId:     src.SnapshotId,
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		Amount:         src.Amount,
		FinalAmount:    src.FinalAmount,
		IsShipped:      src.IsShipped,
		Data:           src.Data,
	}
}

func OrderDto(src *order.ComplexOrder) *order_service.SComplexOrder {
	o := &order_service.SComplexOrder{
		OrderId:         src.OrderId,
		SubOrderId:      src.SubOrderId,
		OrderType:       src.OrderType,
		OrderNo:         src.OrderNo,
		BuyerId:         src.BuyerId,
		VendorId:        src.VendorId,
		ShopId:          src.ShopId,
		Subject:         src.Subject,
		ItemAmount:      src.ItemAmount,
		DiscountAmount:  src.DiscountAmount,
		ExpressFee:      src.ExpressFee,
		PackageFee:      src.PackageFee,
		FinalAmount:     src.FinalAmount,
		ConsigneePerson: src.ConsigneePerson,
		ConsigneePhone:  src.ConsigneePhone,
		ShippingAddress: src.ShippingAddress,
		BuyerComment:    src.BuyerComment,
		CreateTime:      src.CreateTime,
		UpdateTime:      src.UpdateTime,
		State:           src.State,
		Items:           make([]*order_service.SComplexItem, len(src.Items)),
		Data:            src.Data,
	}
	if src.Items != nil {
		for i, v := range src.Items {
			o.Items[i] = OrderItemDto(v)
		}
	}
	return o
}

func Order(src *order_service.SComplexOrder) *order.ComplexOrder {
	o := &order.ComplexOrder{
		OrderId:         src.OrderId,
		SubOrderId:      src.SubOrderId,
		OrderType:       src.OrderType,
		OrderNo:         src.OrderNo,
		BuyerId:         src.BuyerId,
		VendorId:        src.VendorId,
		ShopId:          src.ShopId,
		Subject:         src.Subject,
		ItemAmount:      src.ItemAmount,
		DiscountAmount:  src.DiscountAmount,
		ExpressFee:      src.ExpressFee,
		PackageFee:      src.PackageFee,
		FinalAmount:     src.FinalAmount,
		ConsigneePerson: src.ConsigneePerson,
		ConsigneePhone:  src.ConsigneePhone,
		ShippingAddress: src.ShippingAddress,
		BuyerComment:    src.BuyerComment,
		CreateTime:      src.CreateTime,
		UpdateTime:      src.UpdateTime,
		State:           src.State,
		Items:           make([]*order.ComplexItem, len(src.Items)),
		Data:            src.Data,
	}
	if src.Items != nil {
		for i, v := range src.Items {
			o.Items[i] = OrderItem(v)
		}
	}
	return o
}

func OrderItem(src *order_service.SComplexItem) *order.ComplexItem {
	return &order.ComplexItem{
		ID:             src.ID,
		OrderId:        src.OrderId,
		ItemId:         src.ItemId,
		SkuId:          src.SkuId,
		SnapshotId:     src.SnapshotId,
		Quantity:       src.Quantity,
		ReturnQuantity: src.ReturnQuantity,
		Amount:         src.Amount,
		FinalAmount:    src.FinalAmount,
		IsShipped:      src.IsShipped,
		Data:           src.Data,
	}
}
