package impl

/**
 * Copyright 2015 @ 56x.net.
 * name : 2.payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
import (
	"context"
	"errors"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/module"
	"github.com/ixre/go2o/core/query"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.PaymentServiceServer = new(paymentService)

type paymentService struct {
	repo       payment.IPaymentRepo
	orderRepo  order.IOrderRepo
	memberRepo member.IMemberRepo
	query      *query.PaymentQuery
	serviceUtil
	proto.UnimplementedPaymentServiceServer
}

// QueryDivideOrders implements proto.PaymentServiceServer.
func (p *paymentService) QueryDivideOrders(ctx context.Context, req *proto.DivideOrdersRequest) (*proto.DivideOrdersResponse, error) {
	arr := p.query.QueryDivideOrders(int(req.MemberId), int(req.OrderType))
	ret := &proto.DivideOrdersResponse{
		Orders: make([]*proto.SDivideOrderInfo, 0),
	}
	for _, v := range arr {
		ret.Orders = append(ret.Orders, &proto.SDivideOrderInfo{
			TradeNo:       v.TradeNo,
			Amount:        int64(v.Amount),
			DividedAmount: int64(v.DividedAmount),
			CreateTime:    int64(v.CreateTime),
			DivideStatus:  int32(v.DivideStatus),
		})
	}
	return ret, nil
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo,
	query *query.PaymentQuery) proto.PaymentServiceServer {
	return &paymentService{
		repo:       rep,
		orderRepo:  orderRepo,
		memberRepo: memberRepo,
		query:      query,
	}
}

// GetPaymentOrder 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(_ context.Context, req *proto.PaymentOrderRequest) (*proto.SPaymentOrder, error) {
	if po := p.repo.GetPaymentOrder(req.TradeNo); po != nil {
		v := po.Get()
		sp := p.parsePaymentOrderDto(&v)
		for _, t := range po.TradeMethods() {
			pm := p.parseTradeMethodDataDto(t)
			pm.ChanName = po.ChanName(t.PayMethod)
			if len(pm.ChanName) == 0 {
				pm.ChanName = v.OutTradeSp
			}
			sp.TradeData = append(sp.TradeData, pm)
		}
		return sp, nil
	}
	return nil, payment.ErrNoSuchPaymentOrder
}

// SubmitPaymentOrder 创建支付单
func (p *paymentService) SubmitPaymentOrder(_ context.Context, order *proto.SPaymentOrder) (*proto.TxResult, error) {
	v := p.parsePaymentOrder(order)
	o := p.repo.CreatePaymentOrder(v)
	err := o.Submit()
	return p.errorV2(err), nil
}

// AdjustOrder 调整支付单金额
func (p *paymentService) AdjustOrder(_ context.Context, r *proto.AdjustOrderRequest) (*proto.TxResult, error) {
	var err error
	o := p.repo.GetPaymentOrder(r.PaymentNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.Adjust(int(r.Amount * 100))
	}
	return p.errorV2(err), nil
}

// DiscountByIntegral 积分抵扣支付单
func (p *paymentService) DiscountByIntegral(_ context.Context, r *proto.DiscountIntegralRequest) (*proto.TxResult, error) {
	var amount int
	var err error
	o := p.repo.GetPaymentOrderById(int(r.OrderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		amount, err = o.IntegralDiscount(int(r.Integral), r.IgnoreOut)
	}
	rs := p.errorV2(err)
	rs.Data = map[string]string{"Amount": strconv.Itoa(amount)}
	return rs, nil
}

// DiscountByBalance 余额抵扣
func (p *paymentService) DiscountByBalance(_ context.Context, r *proto.DiscountBalanceRequest) (*proto.TxResult, error) {
	var err error
	o := p.repo.GetPaymentOrderById(int(r.OrderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDeduct(r.Remark)
	}
	return p.errorV2(err), nil
}

// PaymentByWallet 钱包账户支付
func (p *paymentService) PaymentByWallet(_ context.Context, r *proto.WalletPaymentRequest) (rs *proto.TxResult, err error) {
	arr := p.repo.GetMergePayOrders(r.TradeNo)
	if len(arr) == 0 {
		// 单个订单支付
		ip := p.repo.GetPaymentOrder(r.TradeNo)
		if ip == nil {
			err = payment.ErrNoSuchPaymentOrder
		} else {
			err = ip.PaymentByWallet(r.Remark)
		}
		return p.errorV2(err), nil
	}
	// 合并支付单支付
	payUid := arr[0].Get().PayerId
	var finalAmount int64 = 0
	for _, v := range arr {
		finalAmount += int64(v.Get().FinalAmount)
	}
	acc := p.memberRepo.GetAccount(int(payUid))
	if acc.Balance < int(finalAmount) {
		err = member.ErrAccountBalanceNotEnough
	} else {
		for _, v := range arr {
			if err = v.PaymentByWallet(r.Remark); err != nil {
				break
			}
		}
	}
	return p.errorV2(err), nil
}

// HybridPayment 余额钱包混合支付，优先扣除余额。
func (p *paymentService) HybridPayment(_ context.Context, r *proto.HyperPaymentRequest) (rs *proto.TxResult, err error) {
	o := p.repo.GetPaymentOrder(r.TradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.HybridPayment(r.Remark)
	}
	return p.errorV2(err), nil
}

// FinishPayment 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(_ context.Context, r *proto.FinishPaymentRequest) (rs *proto.TxResult, err error) {
	o := p.repo.GetPaymentOrder(r.PaymentOrderNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(r.SpName, r.SpTradeNo)
	}
	return p.errorV2(err), nil
}

// GatewayV1 支付网关
func (p *paymentService) GatewayV1(_ context.Context, r *proto.PayGatewayRequest) (rs *proto.TxResult, err error) {
	mod := module.Get(module.PAY).(*module.PaymentModule)
	// 获取令牌
	if r.Action == "get_token" {
		token := mod.CreateToken(r.UserId)
		return p.successV2(map[string]string{"token": token}), nil
	}
	// 提交支付请求
	if r.Action == "submit" {
		err = mod.Submit(r.UserId, r.Data)
	}
	// 验证支付
	if r.Action == "payment" {
		err = mod.CheckAndPayment(r.UserId, r.Data)
	}
	return p.errorV2(err), nil
}

// GetPreparePaymentInfo 获取支付预交易数据
func (p *paymentService) GetPreparePaymentInfo(_ context.Context, r *proto.OrderInfoRequest) (*proto.SPrepareTradeData, error) {
	var arr []payment.IPaymentOrder
	if r.MergePay {
		arr = p.repo.GetMergePayOrders(r.TradeNo)
	} else {
		ip := p.repo.GetPaymentOrder(r.TradeNo)
		if ip != nil {
			arr = []payment.IPaymentOrder{ip}
		}
	}
	return p.getMergePaymentOrdersInfo(r.TradeNo, arr, false)
}

// 获取合并支付的支付单的支付数据
func (p *paymentService) getMergePaymentOrdersInfo(tradeNo string,
	tradeOrders []payment.IPaymentOrder, checkPay bool) (*proto.SPrepareTradeData, error) {
	d := &proto.SPrepareTradeData{ErrCode: 1, TradeOrders: []*proto.SPaymentOrderData{}}
	if len(tradeOrders) == 0 {
		d.ErrMsg = "无效的支付订单"
		return d, nil
	}
	d.TradeStatus = payment.StateAwaitingPayment // 待支付
	for _, ip := range tradeOrders {
		// 检查支付状态
		if checkPay {
			if err := ip.CheckPaymentState(); err != nil {
				d.ErrMsg = err.Error()
				return d, nil
			}
		}
		iv := ip.Get()
		so := &proto.SPaymentOrderData{
			OrderNo:        iv.OutOrderNo,
			Subject:        iv.Subject,
			TradeType:      iv.TradeType,
			Status:         int32(iv.Status),
			TransactionFee: int64(iv.TransactionFee),
			FinalAmount:    int64(iv.FinalAmount),
		}
		// 更新支付状态
		if so.Status != payment.StateAwaitingPayment {
			d.TradeStatus = so.Status
		}
		// 更新支付标志
		if i := int32(iv.PayFlag); d.PayFlag != i {
			d.PayFlag = i
		}
		// 更新支付金额
		d.TradeOrders = append(d.TradeOrders, so)
		d.TransactionFee += so.TransactionFee  // 手续费
		d.FinalAmount += so.FinalAmount        // 最终金额
		d.TotalAmount += int64(iv.TotalAmount) // 累计金额
	}
	d.ErrCode = 0
	d.TradeNo = tradeNo // 交易单号
	return d, nil
}

// GatewayV2 支付网关V2
func (p *paymentService) GatewayV2(_ context.Context, r *proto.PayGatewayV2Request) (*proto.PayGatewayResponse, error) {
	var arr []payment.IPaymentOrder
	if r.MergePay {
		arr = p.repo.GetMergePayOrders(r.TradeNo)
	} else {
		ip := p.repo.GetPaymentOrder(r.TradeNo)
		if ip != nil {
			arr = []payment.IPaymentOrder{ip}
		}
	}
	if len(arr) == 0 {
		return &proto.PayGatewayResponse{ErrCode: 1,
			ErrMsg: "支付单不存在"}, nil
	}
	for _, ip := range arr {
		if err := ip.CheckPaymentState(); err != nil {
			return &proto.PayGatewayResponse{ErrCode: 2,
				ErrMsg: err.Error()}, nil
		}
	}
	ret := proto.PayGatewayResponse{
		TradeNo: r.TradeNo,
	}
	for _, ip := range arr {
		iv := ip.Get()
		ret.TransactionFee += int64(iv.TransactionFee) // 手续费
		ret.FinalAmount += int64(iv.FinalAmount)       // 最终金额
		ret.TotalAmount += int64(iv.TotalAmount)       // 累计金额
	}
	return &ret, nil
}

// MixedPayment 混合支付
func (p *paymentService) MixedPayment(_ context.Context, _ *proto.MixedPaymentRequest) (*proto.TxResult, error) {
	return nil, errors.New("not support MixedPayment")
}

// parsePaymentOrder 转换为支付单
func (p *paymentService) parsePaymentOrder(src *proto.SPaymentOrder) *payment.Order {
	dst := &payment.Order{
		Id:             int(src.Id),
		SellerId:       int(src.SellerId),
		TradeType:      src.TradeType,
		TradeNo:        src.TradeNo,
		OrderType:      int(src.OrderType),
		OutOrderNo:     src.OutOrderNo,
		Subject:        src.Subject,
		BuyerId:        int(src.BuyerId),
		PayerId:        int(src.PayerId),
		TotalAmount:    int(src.TotalAmount),
		DeductAmount:   int(src.DeductAmount),
		AdjustAmount:   int(src.AdjustAmount),
		TransactionFee: int(src.TransactionFee),
		FinalAmount:    int(src.FinalAmount),
		PaidAmount:     int(src.PaidAmount),
		PayFlag:        int(src.PayFlag),
		FinalFlag:      int(src.FinalFlag),
		ExtraData:      src.ExtraData,
		Status:         int(src.Status),
		SubmitTime:     int(src.SubmitTime),
		ExpiresTime:    int(src.ExpiresTime),
		PaidTime:       int(src.PaidTime),
		TradeMethods:   make([]*payment.PayTradeData, 0),
	}
	if src.SubOrder {
		dst.SubOrder = 1
	}
	return dst
}

// parsePaymentOrderDto 转换为支付单数据
func (p *paymentService) parsePaymentOrderDto(src *payment.Order) *proto.SPaymentOrder {
	return &proto.SPaymentOrder{
		Id:             int32(src.Id),
		SellerId:       int32(src.SellerId),
		TradeType:      src.TradeType,
		TradeNo:        src.TradeNo,
		Subject:        src.Subject,
		BuyerId:        int64(src.BuyerId),
		PayerId:        int64(src.PayerId),
		TotalAmount:    int64(src.TotalAmount),
		DeductAmount:   int64(src.DeductAmount),
		AdjustAmount:   int64(src.AdjustAmount),
		TransactionFee: int64(src.TransactionFee),
		FinalAmount:    int64(src.FinalAmount),
		PaidAmount:     int64(src.PaidAmount),
		PayFlag:        int32(src.PayFlag),
		FinalFlag:      int32(src.FinalFlag),
		ExtraData:      src.ExtraData,
		Status:         int32(src.Status),
		SubmitTime:     int64(src.SubmitTime),
		ExpiresTime:    int64(src.ExpiresTime),
		PaidTime:       int64(src.PaidTime),
		SubOrder:       src.SubOrder == 1,
		OrderType:      int32(src.OrderType),
		OutOrderNo:     src.OutOrderNo,
		TradeData:      make([]*proto.STradeChanData, 0),
	}
}

// parseTradeMethodDataDto 转换为交易渠道数据
func (p *paymentService) parseTradeMethodDataDto(src *payment.PayTradeData) *proto.STradeChanData {
	return &proto.STradeChanData{
		ChanId:     int32(src.PayMethod),
		Amount:     int64(src.PayAmount),
		ChanCode:   src.OutTradeCode,
		OutTradeNo: src.OutTradeNo,
	}
}

// SaveIntegrateApp 保存集成支付应用
func (p *paymentService) SaveIntegrateApp(_ context.Context, app *proto.SIntegrateApp) (*proto.TxResult, error) {
	_, err := p.repo.SaveIntegrateApp(&payment.IntegrateApp{
		Id:            int(app.Id),
		AppName:       app.AppName,
		AppUrl:        app.AppUrl,
		Enabled:       int(app.Enabled),
		IntegrateType: int(app.IntegrateType),
		SortNumber:    int(app.SortNumber),
		Hint:          app.Hint,
		Highlight:     int(app.Highlight),
	})
	return p.errorV2(err), nil
}

// QueryIntegrateAppList 查询集成支付应用列表
func (p *paymentService) QueryIntegrateAppList(_ context.Context, _ *proto.Empty) (*proto.QueryIntegrateAppResponse, error) {
	arr := p.repo.FindAllIntegrateApp()
	ret := &proto.QueryIntegrateAppResponse{
		Value: make([]*proto.SIntegrateApp, 0),
	}
	for _, v := range arr {
		if v.Enabled == 1 {
			ret.Value = append(ret.Value, p.parseIntegrateApp(v))
		}
	}
	return ret, nil
}

// PrepareIntegrateParams 准备集成支付的参数
func (p *paymentService) PrepareIntegrateParams(_ context.Context, req *proto.IntegrateParamsRequest) (*proto.IntegrateParamsResponse, error) {
	arr := p.repo.FindAllIntegrateApp()
	var ret *payment.IntegrateApp
	for _, v := range arr {
		if v.Id == int(req.AppId) {
			ret = v
		}
	}
	if ret == nil {
		return &proto.IntegrateParamsResponse{
			ErrCode: 1,
			ErrMsg:  "no such payment app",
		}, nil
	}
	if ret.Enabled != 1 {
		return &proto.IntegrateParamsResponse{
			ErrCode: 2,
			ErrMsg:  ret.AppName + "暂不可用",
		}, nil
	}
	ord := p.repo.GetPaymentOrder(req.PayOrderNo)
	if ord == nil {
		return &proto.IntegrateParamsResponse{
			ErrCode: 3,
			ErrMsg:  "支付单无效",
		}, nil
	}
	ov := ord.Get()
	return &proto.IntegrateParamsResponse{
		AppId:       int32(ret.Id),
		AppName:     ret.AppName,
		AppUrl:      ret.AppUrl,
		Service:     "pay",
		OrderNo:     ov.TradeNo,
		OrderAmount: int32(ov.FinalAmount),
		Subject:     ov.Subject,
	}, nil
}

// parseIntegrateApp 转换为集成支付应用
func (p *paymentService) parseIntegrateApp(v *payment.IntegrateApp) *proto.SIntegrateApp {
	return &proto.SIntegrateApp{
		Id:            int32(v.Id),
		AppName:       v.AppName,
		AppUrl:        v.AppUrl,
		Enabled:       int32(v.Enabled),
		IntegrateType: int32(v.IntegrateType),
		SortNumber:    int32(v.SortNumber),
		Hint:          v.Hint,
		Highlight:     int32(v.Highlight),
	}
}
func (p *paymentService) DeleteIntegrateApp(_ context.Context, id *proto.PayIntegrateAppId) (*proto.TxResult, error) {
	err := p.repo.DeleteIntegrateApp(id.Value)
	return p.errorV2(err), nil
}

// Divide 分账
func (p *paymentService) Divide(_ context.Context, req *proto.PaymentDivideRequest) (*proto.TxResult, error) {
	if len(req.SubDivides) == 0 {
		return p.errorV2(errors.New("分账明细不正确")), nil
	}
	for _, v := range req.SubDivides {
		ip := p.repo.GetPaymentOrder(v.TradeNo)
		if ip == nil {
			return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
		}
		divides := make([]*payment.DivideData, len(v.Divides))
		for i, v := range v.Divides {
			divides[i] = &payment.DivideData{
				DivideType:   int(v.DivideType),
				UserId:       int(v.UserId),
				DivideAmount: int(v.DivideAmount),
			}
		}
		err := ip.Divide(req.OutTxNo, divides)
		if err != nil {
			return p.errorV2(err), nil
		}
	}
	return p.errorV2(nil), nil
}

// CompleteDivide 完成分账
func (p *paymentService) CompleteDivide(_ context.Context, req *proto.PaymentOrderRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrder(req.TradeNo)
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	err := ip.CompleteDivide()
	return p.errorV2(err), nil
}

// UpdateDivideStatus 更新分账状态
func (p *paymentService) UpdateDivideStatus(_ context.Context, req *proto.UpdateDivideStatusRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrderById(int(req.PayId))
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	err := ip.UpdateSubDivideStatus(int(req.DivideId), req.Success, req.DivideNo, req.Remark)
	return p.errorV2(err), nil
}

// RevertSubDivide 撤销分账
func (p *paymentService) RevertSubDivide(_ context.Context, req *proto.PaymentSubDivideRevertRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrderById(int(req.PayId))
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	err := ip.RevertSubDivide(int(req.DivideId), req.Reason)
	return p.errorV2(err), nil
}

// DivideSuccess 分账成功
func (p *paymentService) DivideSuccess(_ context.Context, req *proto.PaymentDivideSuccessRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrderById(int(req.PayId))
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	err := ip.DivideSuccess(req.OutTxNo)
	return p.errorV2(err), nil
}

// QueryRefundableOrders 查询可退款支付单
func (p *paymentService) QueryRefundableOrders(_ context.Context, req *proto.QueryRefundablePaymentOrdersRequest) (*proto.RefundablePaymentOrdersResponse, error) {
	arr := p.query.QueryRefundableRechargeOrders(int(req.MemberId))
	orders := make([]*proto.RefundablePaymentOrder, len(arr))
	for i, v := range arr {
		orders[i] = &proto.RefundablePaymentOrder{
			TradeNo:          v.TradeNo,
			RefundableAmount: int64(v.RefundableAmount),
		}
	}
	return &proto.RefundablePaymentOrdersResponse{
		Orders: orders,
	}, nil
}

// RequestRefund 请求退款
func (p *paymentService) RequestRefund(_ context.Context, req *proto.PaymentRefundRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrder(req.TradeNo)
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	if ip.Get().OrderType == payment.TypeOrder {
		return p.errorV2(errors.New("当前支付单不支持退款操作[TYPE_NOT_MATCH]")), nil
	}
	err := ip.Refund(map[int]int{
		payment.MPaySP: int(req.RefundAmount),
	}, req.Reason)
	return p.errorV2(err), nil
}

// RequestRefundAvail 请求全额退款
func (p *paymentService) RequestRefundAvail(_ context.Context, req *proto.PaymentRefundAvailRequest) (*proto.PaymentRefundAvailResponse, error) {
	ip := p.repo.GetPaymentOrder(req.TradeNo)
	if ip == nil {
		return &proto.PaymentRefundAvailResponse{
			Code:    1,
			Message: "支付单不存在",
		}, nil
	}
	if ip.Get().OrderType == payment.TypeOrder {
		return &proto.PaymentRefundAvailResponse{
			Code:    1,
			Message: "当前支付单不支持退款操作[TYPE_NOT_MATCH]",
		}, nil
	}
	amount, err := ip.RefundAvail(req.Remark)
	if err != nil {
		return &proto.PaymentRefundAvailResponse{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &proto.PaymentRefundAvailResponse{
		FinalAmount:  int64(ip.Get().FinalAmount),
		RefundAmount: int64(amount),
	}, nil
}

// GetSubMerchant implements proto.PaymentServiceServer.
func (p *paymentService) GetSubMerchant(_ context.Context, req *proto.SubMerchantCodeRequest) (*proto.SSubMerchant, error) {
	v := p.repo.SubMerchantManager().GetMerchant(req.Code)
	if v == nil {
		return &proto.SSubMerchant{}, nil
	}
	return &proto.SSubMerchant{
		Id:                    int64(v.Id),
		Code:                  v.Code,
		UserType:              int32(v.UserType),
		UserId:                int64(v.UserId),
		MchType:               int32(v.MchType),
		MchRole:               int32(v.MchRole),
		LicencePic:            v.LicencePic,
		SignName:              v.SignName,
		SignType:              int32(v.SignType),
		LicenceNo:             v.LicenceNo,
		ShortName:             v.ShortName,
		AccountLicencePic:     v.AccountLicencePic,
		LegalName:             v.LegalName,
		LegalLicenceType:      int32(v.LegalLicenceType),
		LegalLicenceNo:        v.LegalLicenceNo,
		LegalFrontPic:         v.LegalFrontPic,
		LegalBackPic:          v.LegalBackPic,
		ContactName:           v.ContactName,
		ContactPhone:          v.ContactPhone,
		ContactEmail:          v.ContactEmail,
		ContactLicenceNo:      v.ContactLicenceNo,
		AccountEmail:          v.AccountEmail,
		AccountPhone:          v.AccountPhone,
		PrimaryIndustryCode:   v.PrimaryIndustryCode,
		SecondaryIndustryCode: v.SecondaryIndustryCode,
		ProvinceCode:          int32(v.ProvinceCode),
		CityCode:              int32(v.CityCode),
		DistrictCode:          int32(v.DistrictCode),
		Address:               v.Address,
		SettleDirection:       int32(v.SettleDirection),
		SettleBankCode:        v.SettleBankCode,
		SettleAccountType:     int32(v.SettleAccountType),
		SettleBankAccount:     v.SettleBankAccount,
		IssueMchNo:            v.IssueMchNo,
		AgreementSignUrl:      v.AgreementSignUrl,
		IssueStatus:           int32(v.IssueStatus),
		IssueMessage:          v.IssueMessage,
		CreateTime:            int64(v.CreateTime),
		UpdateTime:            int64(v.UpdateTime),
	}, nil
}

// InitialSubMerchant implements proto.PaymentServiceServer.
func (p *paymentService) InitialSubMerchant(_ context.Context, req *proto.SubMerchantInitialRequest) (*proto.TxResult, error) {
	mgr := p.repo.SubMerchantManager()
	mch, err := mgr.InitialMerchant(int(req.UserType), int(req.UserId))
	if err == nil {
		return p.txResult(mch.Id, map[string]string{
			"code": mch.Code,
		}), nil
	}
	return p.errorV2(err), nil
}

// StageSubMerchant implements proto.PaymentServiceServer.
func (p *paymentService) StageSubMerchant(_ context.Context, req *proto.SSubMerchant) (*proto.TxResult, error) {
	mgr := p.repo.SubMerchantManager()
	dst := &payment.PayMerchant{
		Id:                    int(req.Id),
		Code:                  req.Code,
		UserType:              int(req.UserType),
		UserId:                int(req.UserId),
		MchType:               int(req.MchType),
		MchRole:               int(req.MchRole),
		LicencePic:            req.LicencePic,
		SignName:              req.SignName,
		SignType:              int(req.SignType),
		LicenceNo:             req.LicenceNo,
		ShortName:             req.ShortName,
		AccountLicencePic:     req.AccountLicencePic,
		LegalName:             req.LegalName,
		LegalLicenceType:      int(req.LegalLicenceType),
		LegalLicenceNo:        req.LegalLicenceNo,
		LegalFrontPic:         req.LegalFrontPic,
		LegalBackPic:          req.LegalBackPic,
		ContactName:           req.ContactName,
		ContactPhone:          req.ContactPhone,
		ContactEmail:          req.ContactEmail,
		ContactLicenceNo:      req.ContactLicenceNo,
		AccountEmail:          req.AccountEmail,
		AccountPhone:          req.AccountPhone,
		PrimaryIndustryCode:   req.PrimaryIndustryCode,
		SecondaryIndustryCode: req.SecondaryIndustryCode,
		ProvinceCode:          int(req.ProvinceCode),
		CityCode:              int(req.CityCode),
		DistrictCode:          int(req.DistrictCode),
		Address:               req.Address,
		SettleDirection:       int(req.SettleDirection),
		SettleBankCode:        req.SettleBankCode,
		SettleAccountType:     int(req.SettleAccountType),
		SettleBankAccount:     req.SettleBankAccount,
		IssueMchNo:            req.IssueMchNo,
		AgreementSignUrl:      req.AgreementSignUrl,
		IssueStatus:           int(req.IssueStatus),
		IssueMessage:          req.IssueMessage,
		CreateTime:            int(req.CreateTime),
		UpdateTime:            int(req.UpdateTime),
	}
	err := mgr.StageMerchant(dst)
	return p.errorV2(err), nil
}

// SubmitSubMerchant implements proto.PaymentServiceServer.
func (p *paymentService) SubmitSubMerchant(_ context.Context, req *proto.SubMerchantCodeRequest) (*proto.TxResult, error) {
	mgr := p.repo.SubMerchantManager()
	err := mgr.Submit(req.Code)
	return p.errorV2(err), nil
}

// UpdateSubMerchant implements proto.PaymentServiceServer.
func (p *paymentService) UpdateSubMerchant(_ context.Context, req *proto.SubMerchantUpdateRequest) (*proto.TxResult, error) {
	mgr := p.repo.SubMerchantManager()
	err := mgr.Update(req.Code, &payment.SubMerchantUpdateParams{
		Status:           int(req.Status),
		Remark:           req.Remark,
		MerchantCode:     req.MerchantCode,
		AgreementSignUrl: req.AgreementSignUrl,
	})
	return p.errorV2(err), nil
}
