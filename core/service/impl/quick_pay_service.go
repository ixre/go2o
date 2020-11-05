package impl

import (
	"context"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/registry"
	"go2o/core/infrastructure/qpay"
	"go2o/core/infrastructure/qpay/hfb"
	"go2o/core/service/proto"
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

func (q quickPayServiceImpl) CheckSign(context context.Context, request *proto.CheckQPaySignRequest) (*proto.CheckQPaySignResponse, error) {
	panic("implement me")
}

func (q quickPayServiceImpl) RequestBankSideAuth(context context.Context, request *proto.BankAuthRequest) (*proto.BankAuthResponse, error) {
	panic("implement me")
}

func (q quickPayServiceImpl) QueryBankAuthResult(context context.Context, request *proto.BankAuthQueryRequest) (*proto.BankAuthQueryResponse, error) {
	panic("implement me")
}

func (q quickPayServiceImpl) DirectPayment(context context.Context, request *proto.QPaymentRequest) (*proto.QPaymentResponse, error) {
	panic("implement me")
}

func (q quickPayServiceImpl) BatchTransfer(context context.Context, request *proto.BatchTransferRequest) (*proto.BatchTransferResponse, error) {
	panic("implement me")
}

func initQPayConfig(repo registry.IRegistryRepo) {
	// 初始化HFB
	if _,err := repo.GetValue("qp_hfb_agent_id");err != nil {
		repo.CreateUserKey("qp_hfb_agent_id", "0000000", "汇付宝(快捷支付)商户编号")
	}
	if _,err := repo.GetValue("qp_hfb_md5_key");err != nil{
		repo.CreateUserKey("qp_hfb_md5_key", "CC08C5E3E69F4E6B85F1DC0B", "汇付宝(快捷支付)签名KEY(md5)")
	}
}
