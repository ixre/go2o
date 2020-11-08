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
		CardType:            int32(r.CardType),
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
			Code:          2,
			ErrMsg:        err.Error(),
			BankAuthToken: "",
		}, nil
	}
	return &proto.BankAuthQueryResponse{
		Code:          int32(result.Code),
		ErrMsg:        result.Message,
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

func (q quickPayServiceImpl) QueryPaymentStatus(_ context.Context, r *proto.QPaymentQueryRequest) (*proto.QPaymentQueryResponse, error) {
	ret,err := q.qp.QueryPaymentStatus(r.OrderNo,r.Options)
	if err != nil{
		return &proto.QPaymentQueryResponse{
			Code:                 1,
			ErrMsg:               err.Error(),
		},nil
	}
	return &proto.QPaymentQueryResponse{
		Code:                 int32(ret.Code),
		ErrMsg:               ret.ErrMsg,
		BillNo:               ret.BillNo,
	},nil
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
	ret, err := q.qp.BatchTransfer(r.BatchTradeNo, batchList, r.Nonce, r.NotifyUrl)
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
			OrderNo:        v.OrderNo,
			BankCode:       v.BankCode,
			PersonTransfer: v.PersonTransfer,
			TradeFee:       v.TradeFee,
			BankCardNo:     v.BankCardNo,
			BankAccountName:v.BankAccountName,
			Subject:        v.Subject,
			Province:       v.Province,
			City:           v.City,
			StoreName:      v.StoreName,
		}
	}
	return dst
}

func initQPayConfig(repo registry.IRegistryRepo) {
	// 快捷支付安全签名密钥
	if _, err := repo.GetValue("qp_safe_secret"); err != nil {
		l := util.RandString(10)
		token := strings.ToLower(crypto.Md5([]byte(l)))
		_ = repo.CreateUserKey("qp_safe_secret", token[8:26], "快捷支付安全签名密钥")
	}
	// 初始化HFB
	if _, err := repo.GetValue("qp_hfb_agent_id"); err != nil {
		_ = repo.CreateUserKey("qp_hfb_agent_id", "0000000", "汇付宝(快捷支付)商户编号")
	}
	if _, err := repo.GetValue("qp_hfb_md5_key"); err != nil {
		_ = repo.CreateUserKey("qp_hfb_md5_key", "CC08C5E3E69F4E6B85F1DC0B", "汇付宝(快捷支付)签名KEY(md5)")
	}
	if _, err := repo.GetValue("qp_hfb_query_md5_key"); err != nil {
		_ = repo.CreateUserKey("qp_hfb_query_md5_key", "123456", "汇付宝(快捷支付)自定义签名KEY(md5)")
	}
	if _, err := repo.GetValue("qp_hfb_batch_3des_key"); err != nil {
		_ = repo.CreateUserKey("qp_hfb_batch_3des_key", "4865534416254C0F8837DFB3", "汇付宝(快捷支付)批付3DES(KEY)")
	}

	if _, err := repo.GetValue("qp_hfb_public_key"); err != nil {
		publicKey := "MIIBIjASBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsVR6LGVO7kbIBKKuAljjPS+V46Ij8+GVCIhIdx5Nj4kJsByM+wo7Nu8QiZczZsR/Yl9n0hYdb1blAO+O0sA4Dg2ALMJeYamxDe5acC+N5W1aVSiOsqiMmKIX7nOSYL2bPLx6uMG/VZjogZBoqHY5qTQH5AX4nQeqW3rAQACKljuqFTl580+TSZqv+QHcCKQqNDmmFW31a1icELoPWhZF7f+Ry1wr7Q4W1ScpLX3uZZadqsZtH7rvvk+SjxV3y5iCD8ZKFqRdxbuuXXcw+GEth6t0kp5EALkdmJFtIq4uI3lgyqCB+PJq4tyBDZOsU4tY/PqZJ+EbbrPRacRf7ecX0wIDAQAB"
		_ = repo.CreateUserKey("qp_hfb_public_key", publicKey, "汇付宝(快捷支付)公钥")
	}
	if _, err := repo.GetValue("qp_hfb_private_key"); err != nil {
		privateKey := "MIIEvgABADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDF/hqHZZb7r0S5KuuQ1zE4v6BT+irjybOR0mIBbRqUnlBlIK8eayxs7eTazEn7FIFjepvGMxgH/2tC6R7s45KaoQo5Yq9l/rvziyYI50U4SZor1mV24nlCNLbx5BqDBFcGwxOJqwZGVTelBVjDtOsper10rUjhtwDFcLSe82VoPQUt8k9H4zw8+0lC4DsK0JlNtRJNAi380Fmz5JV19+12D2N8Tn9+pqFXzjyvp2EyJ/hS8uHUXZGy3lh7cbeEkFu5sFcKB2RDSs++8Y5vyeXQ6RLqMlEbJIRcRRAeMaCZ2Vn5OATYQKCvTPmITTzKB7NoOvEOC9FO4V6HMjidZzBTAgMBAAECggEADufwi10EnvI1FFO85GyvEfyrT2c4L2oSENpr8nuKUsIQf2yUgo/DCnhmkGps73A9xYWHkMZr+r4qDyGJ6H/Bm86f/G4HkoA5Gj7RoD35IiG4b7B2dxrZ0jgxxchMjqyW+LVbFTRBBq6Hv+7FHgbS5Y6OEOiy4ftrHXI8xvLAIbbEa9k1EVmH2ZvA5iVTBuZGWsEAQMRrIBNpmyB3Lnmo7iK28vpEPLvxADtlr3/1vpwfIPMb2fUYkuMXsCPuxjGxtkiCNhahUyzzwGG8rvszx/JcP/vWwRC7IQQff+YONdGKrJT5VqchJV1oaKbLg9CbU1/xsuLOn2RZP1A3/ssdsQKBgQDrlYhZ8BYSa2l5euKX7r4NFGETD8UGnyJmCGPy22VstJ77vAvffVLkKSzWrZgOlmW8MdRfFUsLfPaolLx56rCtdgS6mwSh4kqz9nKMuQjQbpECJAJtZL4FuMjVKSL/71Kew3/Bc/MNo6uKGxiK54KjxFu4TXWplKHFAI1MPuhdvQKBgQDXJpkFta6XwWbtrBCrgN5+eROA9qP+xC0WF/Ar8jbNJAntoUYXFLkIMt1HJFKAPND71x54G0ZHHpL7LJCP/NiGhY19/4S1oBP79d67HPku9Kbrm1NXKUzafOv2rPXSK7uGR+XSgnnKbs5GicipcqZP3+OGOajb9xxjer0IpU//TwKBgQCfHy8r4FhoNJjXbsMicCV6XCt9XodsA4yOclhgLwSAujcwPUGfwNx+M7mPf01XfQpWZSnW12EK72sDTwNHLdgMMczb5dzpIxnmGC4jEs/7SNM1KPFixkr7PmaYY+K6EAI0LkRafGDM86Hn9IlNOTYqO3TgNaGl2zixAcBuoYb92QKBgFsS7aerFrMKnWVydsQCkyx6WDU5MoZ/yI4XqAUSTPxdiw5aPG88yG6eCWk6COpb1CMnFrDE6uTkHlfQr4kkAQxAsHprlWPE1XDMzXHre9fSnG4TnB3DT9MVGlWbNZu4A3N+L90CekekzBCz9os0Cw64uXlyIvaqDgxWQnrMb6alAoGBAJ44E3SOo9DD5UOk+6swf/YplhqG2sayJruVib+1D2dlWu/+LxJqQZJGI/jtLVO24q7XGdnlA1YXA85DRI9/VUPPOEaLpUI91KWHUaN0Cgcin/O02UR+UWWvtbNEhI8Huk4BDGOPrBxz1tI2Bw1IvkD6u/mKmiExhzCUX/oAAesT"
		_ = repo.CreateUserKey("qp_hfb_private_key", privateKey, "汇付宝(快捷支付)私钥")
	}
}
