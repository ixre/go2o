package impl

import (
	"context"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/registry"
	"go2o/core/infrastructure/qpay"
	"go2o/core/infrastructure/qpay/hfb"
	"go2o/core/service/proto"
	"strings"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : quick_pay_service.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-05 18:10
 * description :
 * history :
 */

var _ proto.QuickPayServiceServer = new(quickPayServiceImpl)

type quickPayServiceImpl struct {
	s            storage.Interface
	qp           qpay.QuickPayProvider
	registryRepo registry.IRegistryRepo
	serviceUtil
}

func NewQuickPayService(s storage.Interface,
	registryRepo registry.IRegistryRepo) *quickPayServiceImpl {
	initQPayConfig(registryRepo)
	qp := hfb.NewHfb(s)
	return &quickPayServiceImpl{
		s:            s,
		qp:           qp,
		registryRepo: registryRepo,
	}
}

func (q quickPayServiceImpl) QueryCardBin(_ context.Context, no *proto.BankCardNo) (*proto.CardBinQueryResponse, error) {
	r := q.qp.QueryCardBin(no.CardNo)
	return &proto.CardBinQueryResponse{
		ErrMsg:              r.ErrMsg,
		BankName:            r.BankName,
		BankCardNo:          r.BankCardNo,
		BankCode:            r.BankCode,
		CardType:            r.CardType,
		RequireBankSideAuth: r.RequireBankSideAuth,
	}, nil
}

func (q quickPayServiceImpl) CheckSign(_ context.Context, r *proto.CheckQPaySignRequest) (*proto.CheckQPaySignResponse, error) {
	b := q.qp.CheckSign(r.Params, r.SignType, r.Sign)
	return &proto.CheckQPaySignResponse{
		SignOk: b,
	}, nil
}

func (q quickPayServiceImpl) RequestBankSideAuth(_ context.Context, r *proto.BankAuthRequest) (*proto.BankAuthResponse, error) {
	ret, err := q.qp.RequestBankSideAuth(r.NonceId, r.BankCardNo, r.BankAccountName,
		r.IdCardNo, r.Mobile)
	rsp := &proto.BankAuthResponse{
		NonceId: r.NonceId,
	}
	if err != nil {
		rsp.ErrMsg = err.Error()
	} else {
		rsp.AuthForm = ret.AuthForm
		rsp.AuthData = ret.AuthData
	}
	return rsp, nil
}

func (q quickPayServiceImpl) QueryBankAuthResult(_ context.Context, r *proto.BankAuthQueryRequest) (*proto.BankAuthQueryResponse, error) {
	var result *qpay.BankAuthQueryResponse
	var err error
	if len(r.NonceId) > 0 {
		result, err = q.qp.QueryBankAuthByNonceId(r.NonceId)
	} else {
		result, err = q.qp.QueryBankAuth(r.BankCardNo)
	}
	if err != nil {
		return &proto.BankAuthQueryResponse{
			Code:          result.Code,
			ErrMsg:        err.Error(),
			BankAuthToken: "",
		}, nil
	}
	return &proto.BankAuthQueryResponse{
		Code:          result.Code,
		BankAuthToken: result.BankAuthToken,
	}, nil
}

func (q quickPayServiceImpl) DirectPayment(_ context.Context, r *proto.QPaymentRequest) (*proto.QPaymentResponse, error) {
	if len(strings.TrimSpace(r.Nonce)) == 0 {
		return &proto.QPaymentResponse{
			Code:    "2",
			ErrMsg:  "缺少参数:nonce",
			NonceId: r.Nonce,
		}, nil
	}
	if !q.checkSafeSign(r.Nonce, r.SafeSign) {
		return &proto.QPaymentResponse{
			Code:    "3",
			ErrMsg:  "请求签名不正确",
			NonceId: r.Nonce,
		}, nil
	}
	if r.TradeUserIp == "" {
		return &proto.QPaymentResponse{
			Code:    "4",
			ErrMsg:  "交易IP为空",
			NonceId: r.Nonce,
		}, nil
	}
	rsp, err := q.qp.DirectPayment(r.OrderNo, r.TradeFee, r.Subject,
		r.BankAuthToken, r.TradeUserIp, r.ReturnUrl, r.NotifyUrl)
	if err != nil {
		return &proto.QPaymentResponse{
			Code:    "1",
			ErrMsg:  err.Error(),
			NonceId: r.Nonce,
		}, nil
	}
	return &proto.QPaymentResponse{
		Code:    "0",
		BillNo:  rsp.BillNo,
		NonceId: r.Nonce,
	}, nil
}

func (q quickPayServiceImpl) BatchTransfer(_ context.Context, r *proto.BatchTransferRequest) (*proto.BatchTransferResponse, error) {
	if len(strings.TrimSpace(r.Nonce)) == 0 {
		return &proto.BatchTransferResponse{
			Code:    "2",
			ErrMsg:  "缺少参数:nonce",
			NonceId: r.Nonce,
		}, nil
	}
	if !q.checkSafeSign(r.Nonce, r.SafeSign) {
		return &proto.BatchTransferResponse{
			Code:    "3",
			ErrMsg:  "请求签名不正确",
			NonceId: r.Nonce,
		}, nil
	}
	if r.TradeUserIp == "" {
		return &proto.BatchTransferResponse{
			Code:    "4",
			ErrMsg:  "交易IP为空",
			NonceId: r.Nonce,
		}, nil
	}
	batchList := q.parseBatchList(r.BatchList)
	ret, err := q.qp.BatchTransfer(r.BatchTradeNo, r.BatchTradeFee,
		batchList, r.Nonce, r.TradeUserIp, r.NotifyUrl)
	if err != nil {
		return &proto.BatchTransferResponse{
			Code:    "1",
			ErrMsg:  err.Error(),
			NonceId: r.Nonce,
		}, nil
	}
	return &proto.BatchTransferResponse{
		Code:    ret.Code,
		ErrMsg:  err.Error(),
		NonceId: r.Nonce,
	}, nil
}

// 检查安全请求签名,lowercase(md5(nonce+secret))
func (q quickPayServiceImpl) checkSafeSign(nonce string, sign string) bool {
	secret, err := q.registryRepo.GetValue("qp_safe_secret")
	if err != nil {
		println("[ Go2o][ Warning]: ", err.Error(), " key: qp_safe_secret")
		return false
	}
	if len(secret) == 0 {
		println("[ Go2o][ Warning]: quick payment safe secret not set")
	}
	md5 := crypto.Md5([]byte(nonce + secret))
	return strings.ToLower(md5) == sign
}

func (q quickPayServiceImpl) parseBatchList(list []*proto.CardTransferRequest) []*qpay.CardTransferReq {
	dst := make([]*qpay.CardTransferReq, len(list))
	for i, v := range list {
		dst[i] = &qpay.CardTransferReq{
			OrderNo:           v.OrderNo,
			BankCode:          v.BankCode,
			TransferToCompany: v.TransferToCompany,
			TradeFee:          v.TradeFee,
			Subject:           v.Subject,
			Province:          v.Province,
			City:              v.City,
			StoreName:         v.StoreName,
		}
	}
	return dst
}

func initQPayConfig(repo registry.IRegistryRepo) {
	// 快捷支付安全签名密钥
	if _, err := repo.GetValue("qp_safe_secret"); err != nil {
		l := util.RandString(10)
		token := strings.ToLower(crypto.Md5([]byte(l)))
		repo.CreateUserKey("qp_safe_secret", token[8:26], "快捷支付安全签名密钥")
	}
	// 初始化HFB
	if _, err := repo.GetValue("qp_hfb_agent_id"); err != nil {
		repo.CreateUserKey("qp_hfb_agent_id", "0000000", "汇付宝(快捷支付)商户编号")
	}
	if _, err := repo.GetValue("qp_hfb_md5_key"); err != nil {
		repo.CreateUserKey("qp_hfb_md5_key", "CC08C5E3E69F4E6B85F1DC0B", "汇付宝(快捷支付)签名KEY(md5)")
	}

}
