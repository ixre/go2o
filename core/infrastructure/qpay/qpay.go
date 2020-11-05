package qpay

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : qpay.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 19:51
 * description :
 * history :
 */

type QuickPayProvider interface {
	// 查询银行卡信息
	QueryCardBin(bankCardNo string) *CardBinQueryResult
}

// 银行卡查询结果
type CardBinQueryResult struct {
	// 错误信息
	ErrMsg string `protobuf:"bytes,1,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	// 银行名称
	BankName string `protobuf:"bytes,1,opt,name=BankName,proto3" json:"BankName,omitempty"`
	// 用户银行卡号
	BankCardNo string `protobuf:"bytes,2,opt,name=BankCardNo,proto3" json:"BankCardNo,omitempty"`
	// 返回的银行代号
	BankCode string `protobuf:"bytes,3,opt,name=BankCode,proto3" json:"BankCode,omitempty"`
	// 银行卡类型（0=储蓄卡,1=信用卡）
	CardType int32 `protobuf:"varint,4,opt,name=CardType,proto3" json:"CardType,omitempty"`
	// 是否需要银行端授权,如果否,则直接使用短信既可授权
	RequireBankSideAuth  bool     `protobuf:"varint,6,opt,name=RequireBankSideAuth,proto3" json:"RequireBankSideAuth,omitempty"`
}