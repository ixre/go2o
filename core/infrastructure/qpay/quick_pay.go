package qpay

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : quick_pay.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 19:51
 * description :
 * history :
 */

type QuickPayProvider interface {
	// 查询银行卡信息
	QueryCardBin(bankCardNo string) *CardBinQueryResult
	// 检查签名是否正确
	CheckSign(params map[string]string, signType string, sign string) bool
	// 申请银行侧认证授权(某些银行需跳转到银行页面进行授权)
	RequestBankSideAuth(nonce string, bankCardNo string, accountName string,
		idCardNo string, mobile string) (*BankAuthResult, error)
	// 根据随机ID查询银行认证状态(提供给外部查询)
	QueryBankAuthByNonceId(id string) (*BankAuthQueryResponse, error)
	// 根据银行卡查询银行认证状态(提供给内部查询)
	QueryBankAuth(bankCardNo string) (*BankAuthQueryResponse, error)
	// 直接支付
	DirectPayment(orderNo string, fee int32, subject string, bankToken string,
		tradeIp string, notifyUrl string, returnUrl string) (*QPaymentResponse, error)
	// 批量付款
	BatchTransfer(batchTradeNo string, batchTradeFee int32,
		list []*CardTransferReq, nonce string,
		tradeIp string, notifyUrl string) (*BatchTransferResponse,error)
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
	CardType int `protobuf:"varint,4,opt,name=CardType,proto3" json:"CardType,omitempty"`
	// 是否需要银行端授权,如果否,则直接使用短信既可授权
	RequireBankSideAuth  bool     `protobuf:"varint,6,opt,name=RequireBankSideAuth,proto3" json:"RequireBankSideAuth,omitempty"`
}

// 申请认证返回结果,通常直接使用返回的表单提交获取银行的授权, 并再查询授权
type BankAuthResult struct {
	// 随机Id
	NonceId string `protobuf:"bytes,1,opt,name=NonceId,proto3" json:"NonceId,omitempty"`
	// 认证的Form表单
	AuthForm string `protobuf:"bytes,3,opt,name=AuthForm,proto3" json:"AuthForm,omitempty"`
	// 认证需要的其他数据
	AuthData map[string]string `protobuf:"bytes,4,rep,name=AuthData,proto3" json:"AuthData,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

// 银行授权查询响应
type BankAuthQueryResponse struct {
	// 状态码，0表示成功
	Code string `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty"`
	// 银行授权认证返回的Token,未处理授权时为空
	BankAuthToken string `protobuf:"bytes,3,opt,name=BankAuthToken,proto3" json:"BankAuthToken,omitempty"`
}

// 支付申请响应
type QPaymentResponse struct {
	// 状态码，0表示成功
	Code string `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty"`
	// 第三方支单据号
	BillNo string `protobuf:"bytes,3,opt,name=BillNo,proto3" json:"BillNo,omitempty"`
}

// 付款请求
type CardTransferReq struct {
	// 商户订单号
	OrderNo string `protobuf:"bytes,1,opt,name=OrderNo,proto3" json:"OrderNo,omitempty"`
	// 银行编号
	BankCode string `protobuf:"bytes,2,opt,name=BankCode,proto3" json:"BankCode,omitempty"`
	// 是否为对公转账
	TransferToCompany bool `protobuf:"varint,3,opt,name=TransferToCompany,proto3" json:"TransferToCompany,omitempty"`
	// 付款金额,实际金额*100,无小数
	TradeFee int32 `protobuf:"varint,4,opt,name=TradeFee,proto3" json:"TradeFee,omitempty"`
	// 付款事由
	Subject string `protobuf:"bytes,5,opt,name=Subject,proto3" json:"Subject,omitempty"`
	// 省份
	Province string `protobuf:"bytes,6,opt,name=Province,proto3" json:"Province,omitempty"`
	// 城市
	City string `protobuf:"bytes,7,opt,name=City,proto3" json:"City,omitempty"`
	// 支行信息
	StoreName string `protobuf:"bytes,8,opt,name=StoreName,proto3" json:"StoreName,omitempty"`
}

// 批量付款响应
type BatchTransferResponse struct {
	// 状态码，0表示成功
	Code string `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty"`
	// 随机ID
	NonceId string `protobuf:"bytes,3,opt,name=NonceId,proto3" json:"NonceId,omitempty"`
}

// 银行认证存储中间数据
type BankAuthSwapData struct{
	// 银行卡号
	BankCardNo string
	// 账户名
	AccountName string
	// 身份证号码
	IdCardNo string
	// 手机号码
	Mobile string
	// 银行名称
	BankName string
	// 返回的银行代号
	BankCode string
	// 银行卡类型（0=储蓄卡,1=信用卡）
	CardType int
}

// 转换私钥
func ParseRSAPrivateKey(s string)(*rsa.PrivateKey,error){
	block, _ := pem.Decode([]byte(s))
	if block == nil{
		return nil,errors.New("私钥不正确")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}