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
	"github.com/jsix/gof/math"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/service/thrift/idl/gen-go/define"
)

func Member(src *member.Member) *define.Member {
	return &define.Member{
		ID:             src.Id,
		Usr:            src.Usr,
		Pwd:            src.Pwd,
		TradePwd:       src.TradePwd,
		Exp:            src.Exp,
		Level:          src.Level,
		InvitationCode: src.InvitationCode,
		RegFrom:        src.RegFrom,
		RegIp:          src.RegIp,
		RegTime:        src.RegTime,
		CheckCode:      src.CheckCode,
		CheckExpires:   src.CheckExpires,
		State:          src.State,
		LoginTime:      src.LoginTime,
		LastLoginTime:  src.LastLoginTime,
		UpdateTime:     src.UpdateTime,
		DynamicToken:   src.DynamicToken,
		TimeoutTime:    src.TimeoutTime,
	}
}

func Member2(src *define.Member) *member.Member {
	return &member.Member{
		Id:             src.ID,
		Usr:            src.Usr,
		Pwd:            src.Pwd,
		TradePwd:       src.TradePwd,
		Exp:            src.Exp,
		Level:          src.Level,
		InvitationCode: src.InvitationCode,
		RegFrom:        src.RegFrom,
		RegIp:          src.RegIp,
		RegTime:        src.RegTime,
		CheckCode:      src.CheckCode,
		CheckExpires:   src.CheckExpires,
		State:          src.State,
		LoginTime:      src.LoginTime,
		LastLoginTime:  src.LastLoginTime,
		UpdateTime:     src.UpdateTime,
		DynamicToken:   src.DynamicToken,
		TimeoutTime:    src.TimeoutTime,
	}
}

func MemberProfile(src *member.Profile) *define.Profile {
	return &define.Profile{
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

func MemberProfile2(src *define.Profile) *member.Profile {
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

func SummaryDto(src *dto.MemberSummary) *define.MemberSummary {
	return &define.MemberSummary{
		MemberId:          src.MemberId,
		Usr:               src.Usr,
		Name:              src.Name,
		Avatar:            src.Avatar,
		Exp:               src.Exp,
		Level:             src.Level,
		LevelName:         src.LevelName,
		LevelSign:         src.LevelSign,
		LevelOfficial:     int64(src.LevelOfficial),
		InvitationCode:    src.InvitationCode,
		Integral:          int64(src.Integral),
		Balance:           round(src.Balance, 2),
		PresentBalance:    round(src.PresentBalance, 2),
		GrowBalance:       round(src.GrowBalance, 2),
		GrowAmount:        round(src.GrowAmount, 2),
		GrowEarnings:      round(src.GrowEarnings, 2),
		GrowTotalEarnings: round(src.GrowTotalEarnings, 2),
		UpdateTime:        src.UpdateTime,
	}
}

func round(f float32, n int) float64 {
	return math.Round(float64(f), n)
}

func AccountDto(src *member.Account) *define.Account {
	return &define.Account{
		MemberId:          src.MemberId,
		Integral:          int64(src.Integral),
		FreezeIntegral:    int64(src.FreezeIntegral),
		Balance:           round(src.Balance, 2),
		FreezeBalance:     round(src.FreezeBalance, 2),
		ExpiredBalance:    round(src.ExpiredBalance, 2),
		PresentBalance:    round(src.PresentBalance, 2),
		FreezePresent:     round(src.FreezePresent, 2),
		ExpiredPresent:    round(src.ExpiredPresent, 2),
		TotalPresentFee:   round(src.TotalPresentFee, 2),
		FlowBalance:       round(src.FlowBalance, 2),
		GrowBalance:       round(src.GrowBalance, 2),
		GrowAmount:        round(src.GrowAmount, 2),
		GrowEarnings:      round(src.GrowEarnings, 2),
		GrowTotalEarnings: round(src.GrowTotalEarnings, 2),
		TotalConsumption:  round(src.TotalConsumption, 2),
		TotalCharge:       round(src.TotalCharge, 2),
		TotalPay:          round(src.TotalPay, 2),
		PriorityPay:       int64(src.PriorityPay),
		UpdateTime:        src.UpdateTime,
	}
}

func Account(src *define.Account) *member.Account {
	return &member.Account{
		MemberId:          src.MemberId,
		Integral:          int(src.Integral),
		FreezeIntegral:    int(src.FreezeIntegral),
		Balance:           float32(src.Balance),
		FreezeBalance:     float32(src.FreezeBalance),
		ExpiredBalance:    float32(src.ExpiredBalance),
		PresentBalance:    float32(src.PresentBalance),
		FreezePresent:     float32(src.FreezePresent),
		ExpiredPresent:    float32(src.ExpiredPresent),
		TotalPresentFee:   float32(src.TotalPresentFee),
		FlowBalance:       float32(src.FlowBalance),
		GrowBalance:       float32(src.GrowBalance),
		GrowAmount:        float32(src.GrowAmount),
		GrowEarnings:      float32(src.GrowEarnings),
		GrowTotalEarnings: float32(src.GrowTotalEarnings),
		TotalConsumption:  float32(src.TotalConsumption),
		TotalCharge:       float32(src.TotalCharge),
		TotalPay:          float32(src.TotalPay),
		PriorityPay:       int(src.PriorityPay),
		UpdateTime:        src.UpdateTime,
	}
}

func PlatformConfDto(src *valueobject.PlatformConf) *define.PlatformConf {
	return &define.PlatformConf{
		Name:             src.Name,
		Logo:             src.Logo,
		Suspend:          src.Suspend,
		SuspendMessage:   src.SuspendMessage,
		MchGoodsCategory: src.MchGoodsCategory,
		MchPageCategory:  src.MchPageCategory,
	}
}

func PlatFromConf(src *define.PlatformConf) *valueobject.PlatformConf {
	return &valueobject.PlatformConf{
		Name:             src.Name,
		Logo:             src.Logo,
		Suspend:          src.Suspend,
		SuspendMessage:   src.SuspendMessage,
		MchGoodsCategory: src.MchGoodsCategory,
		MchPageCategory:  src.MchPageCategory,
	}
}

func AddressDto(src *member.Address) *define.Address {
	return &define.Address{
		ID:        src.Id,
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

func PaymentOrder(src *define.PaymentOrder) *payment.PaymentOrder {
	return &payment.PaymentOrder{
		Id:               src.ID,
		TradeNo:          src.TradeNo,
		VendorId:         src.VendorId,
		Type:             src.Type,
		OrderId:          src.OrderId,
		Subject:          src.Subject,
		BuyUser:          src.BuyUser,
		PaymentUser:      src.PaymentUser,
		TotalFee:         float32(src.TotalFee),
		BalanceDiscount:  float32(src.BalanceDiscount),
		IntegralDiscount: float32(src.IntegralDiscount),
		SystemDiscount:   float32(src.SystemDiscount),
		CouponDiscount:   float32(src.CouponDiscount),
		SubAmount:        float32(src.SubAmount),
		AdjustmentAmount: float32(src.AdjustmentAmount),
		FinalAmount:      float32(src.FinalAmount),
		PaymentOptFlag:   src.PaymentOptFlag,
		PaymentSign:      src.PaymentSign,
		OuterNo:          src.OuterNo,
		CreateTime:       src.CreateTime,
		PaidTime:         src.PaidTime,
		State:            src.State,
	}
}

func PaymentOrderDto(src *payment.PaymentOrder) *define.PaymentOrder {
	return &define.PaymentOrder{
		ID:               src.Id,
		TradeNo:          src.TradeNo,
		VendorId:         src.VendorId,
		Type:             src.Type,
		OrderId:          src.OrderId,
		Subject:          src.Subject,
		BuyUser:          src.BuyUser,
		PaymentUser:      src.PaymentUser,
		TotalFee:         round(src.TotalFee, 2),
		BalanceDiscount:  round(src.BalanceDiscount, 2),
		IntegralDiscount: round(src.IntegralDiscount, 2),
		SystemDiscount:   round(src.SystemDiscount, 2),
		CouponDiscount:   round(src.CouponDiscount, 2),
		SubAmount:        round(src.SubAmount, 2),
		AdjustmentAmount: round(src.AdjustmentAmount, 2),
		FinalAmount:      round(src.FinalAmount, 2),
		PaymentOptFlag:   src.PaymentOptFlag,
		PaymentSign:      src.PaymentSign,
		OuterNo:          src.OuterNo,
		CreateTime:       src.CreateTime,
		PaidTime:         src.PaidTime,
		State:            src.State,
	}
}

func MemberRelationDto(src *member.Relation) *define.MemberRelation {
	return &define.MemberRelation{
		MemberId:      src.MemberId,
		CardId:        src.CardId,
		InviterId:     src.InviterId,
		InviterStr:    src.InviterStr,
		RegisterMchId: src.RegisterMchId,
	}
}

func TrustedInfoDto(src *member.TrustedInfo) *define.TrustedInfo {
	return &define.TrustedInfo{
		MemberId:   src.MemberId,
		RealName:   src.RealName,
		CardId:     src.CardId,
		TrustImage: src.TrustImage,
		Reviewed:   src.Reviewed,
		ReviewTime: src.ReviewTime,
		Remark:     src.Remark,
		UpdateTime: src.UpdateTime,
	}
}
