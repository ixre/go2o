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
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.PaymentServiceServer = new(paymentService)

type paymentService struct {
	repo       payment.IPaymentRepo
	orderRepo  order.IOrderRepo
	memberRepo member.IMemberRepo
	serviceUtil
	proto.UnimplementedPaymentServiceServer
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo) *paymentService {
	return &paymentService{
		repo:       rep,
		orderRepo:  orderRepo,
		memberRepo: memberRepo,
	}
}

// GetPaymentOrder 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(_ context.Context, req *proto.PaymentOrderRequest) (*proto.SPaymentOrder, error) {
	if po := p.repo.GetPaymentOrder(req.TradeNo); po != nil {
		v := po.Get()
		sp := p.parsePaymentOrderDto(&v)
		for _, t := range po.TradeMethods() {
			sp.TradeData = append(sp.TradeData, p.parseTradeMethodDataDto(t))
		}
		return sp, nil
	}
	return nil, payment.ErrNoSuchPaymentOrder
}

// SubmitPaymentOrder 创建支付单
func (p *paymentService) SubmitPaymentOrder(_ context.Context, order *proto.SPaymentOrder) (*proto.Result, error) {
	v := p.parsePaymentOrder(order)
	o := p.repo.CreatePaymentOrder(v)
	err := o.Submit()
	return p.result(err), nil
}

// AdjustOrder 调整支付单金额
func (p *paymentService) AdjustOrder(_ context.Context, r *proto.AdjustOrderRequest) (*proto.Result, error) {
	var err error
	o := p.repo.GetPaymentOrder(r.PaymentNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.Adjust(int(r.Amount * 100))
	}
	return p.result(err), nil
}

// DiscountByIntegral 积分抵扣支付单
func (p *paymentService) DiscountByIntegral(_ context.Context, r *proto.DiscountIntegralRequest) (*proto.Result, error) {
	var amount int
	var err error
	o := p.repo.GetPaymentOrderById(int(r.OrderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		amount, err = o.IntegralDiscount(int(r.Integral), r.IgnoreOut)
	}
	rs := p.result(err)
	rs.Data = map[string]string{"Amount": strconv.Itoa(amount)}
	return rs, nil
}

// DiscountByBalance 余额抵扣
func (p *paymentService) DiscountByBalance(_ context.Context, r *proto.DiscountBalanceRequest) (*proto.Result, error) {
	var err error
	o := p.repo.GetPaymentOrderById(int(r.OrderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDeduct(r.Remark)
	}
	return p.result(err), nil
}

// PaymentByWallet 钱包账户支付
func (p *paymentService) PaymentByWallet(_ context.Context, r *proto.WalletPaymentRequest) (rs *proto.Result, err error) {
	arr := p.repo.GetMergePayOrders(r.TradeNo)
	if len(arr) == 0 {
		// 单个订单支付
		ip := p.repo.GetPaymentOrder(r.TradeNo)
		if ip == nil {
			err = payment.ErrNoSuchPaymentOrder
		} else {
			err = ip.PaymentByWallet(r.Remark)
		}
		return p.result(err), nil
	}
	// 合并支付单支付
	payUid := arr[0].Get().PayerId
	var finalAmount int64 = 0
	for _, v := range arr {
		finalAmount += v.Get().FinalAmount
	}
	acc := p.memberRepo.GetAccount(payUid)
	if acc.Balance*100 < finalAmount {
		err = member.ErrAccountBalanceNotEnough
	} else {
		for _, v := range arr {
			if err = v.PaymentByWallet(r.Remark); err != nil {
				break
			}
		}
	}
	return p.result(err), nil
}

// HybridPayment 余额钱包混合支付，优先扣除余额。
func (p *paymentService) HybridPayment(_ context.Context, r *proto.HyperPaymentRequest) (rs *proto.Result, err error) {
	o := p.repo.GetPaymentOrder(r.TradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.HybridPayment(r.Remark)
	}
	return p.result(err), nil
}

// FinishPayment 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(_ context.Context, r *proto.FinishPaymentRequest) (rs *proto.Result, err error) {
	o := p.repo.GetPaymentOrder(r.PaymentOrderNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(r.SpName, r.SpTradeNo)
	}
	return p.result(err), nil
}

// GatewayV1 支付网关
func (p *paymentService) GatewayV1(_ context.Context, r *proto.PayGatewayRequest) (rs *proto.Result, err error) {
	mod := module.Get(module.PAY).(*module.PaymentModule)
	// 获取令牌
	if r.Action == "get_token" {
		token := mod.CreateToken(r.UserId)
		return p.success(map[string]string{"token": token}), nil
	}
	// 提交支付请求
	if r.Action == "submit" {
		err = mod.Submit(r.UserId, r.Data)
	}
	// 验证支付
	if r.Action == "payment" {
		err = mod.CheckAndPayment(r.UserId, r.Data)
	}
	return p.result(err), nil
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
	d.TradeState = payment.StateAwaitingPayment // 待支付
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
			OrderNo:      iv.OutOrderNo,
			Subject:      iv.Subject,
			TradeType:    iv.TradeType,
			State:        int32(iv.State),
			ProcedureFee: iv.ProcedureFee,
			FinalAmount:  iv.FinalAmount,
		}
		// 更新支付状态
		if so.State != payment.StateAwaitingPayment {
			d.TradeState = so.State
		}
		// 更新支付标志
		if i := int32(iv.PayFlag); d.PayFlag != i {
			d.PayFlag = i
		}
		// 更新支付金额
		d.TradeOrders = append(d.TradeOrders, so)
		d.ProcedureFee += so.ProcedureFee // 手续费
		d.FinalAmount += so.FinalAmount   // 最终金额
		d.TotalAmount += iv.TotalAmount   // 累计金额
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
		ret.ProcedureFee += iv.ProcedureFee // 手续费
		ret.FinalAmount += iv.FinalAmount   // 最终金额
		ret.TotalAmount += iv.TotalAmount   // 累计金额
	}
	return &ret, nil
}

// MixedPayment 混合支付
func (p *paymentService) MixedPayment(_ context.Context, _ *proto.MixedPaymentRequest) (*proto.Result, error) {
	return nil, errors.New("not support MixedPayment")
}

func (p *paymentService) parsePaymentOrder(src *proto.SPaymentOrder) *payment.Order {
	dst := &payment.Order{
		Id:             int(src.Id),
		SellerId:       int(src.SellerId),
		TradeType:      src.TradeType,
		TradeNo:        src.TradeNo,
		OrderType:      int(src.OrderType),
		OutOrderNo:     src.OutOrderNo,
		Subject:        src.Subject,
		BuyerId:        src.BuyerId,
		PayerId:        src.PayerId,
		TotalAmount:    src.TotalAmount,
		DiscountAmount: src.DiscountAmount,
		DeductAmount:   src.DeductAmount,
		AdjustAmount:   src.AdjustAmount,
		ItemAmount:     src.ItemAmount,
		ProcedureFee:   src.ProcedureFee,
		FinalAmount:    src.FinalAmount,
		PaidAmount:     src.PaidAmount,
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

func (p *paymentService) parsePaymentOrderDto(src *payment.Order) *proto.SPaymentOrder {
	return &proto.SPaymentOrder{
		Id:             int32(src.Id),
		SellerId:       int32(src.SellerId),
		TradeType:      src.TradeType,
		TradeNo:        src.TradeNo,
		Subject:        src.Subject,
		BuyerId:        src.BuyerId,
		PayerId:        src.PayerId,
		TotalAmount:    src.TotalAmount,
		DiscountAmount: src.DiscountAmount,
		DeductAmount:   src.DeductAmount,
		AdjustAmount:   src.AdjustAmount,
		ItemAmount:     src.ItemAmount,
		ProcedureFee:   src.ProcedureFee,
		FinalAmount:    src.FinalAmount,
		PaidAmount:     src.PaidAmount,
		PayFlag:        int32(src.PayFlag),
		FinalFlag:      int32(src.FinalFlag),
		ExtraData:      src.ExtraData,
		State:          int32(src.State),
		SubmitTime:     src.SubmitTime,
		ExpiresTime:    src.ExpiresTime,
		PaidTime:       src.PaidTime,
		SubOrder:       src.SubOrder == 1,
		OrderType:      int32(src.OrderType),
		OutOrderNo:     src.OutOrderNo,

		OutTradeNo: src.OutTradeNo,
		OutTradeSp: src.OutTradeSp,
		TradeData:  make([]*proto.STradeMethodData, 0),
	}
}

func (p *paymentService) parseTradeMethodDataDto(src *payment.TradeMethodData) *proto.STradeMethodData {
	return &proto.STradeMethodData{
		Method:     int32(src.Method),
		Code:       src.Code,
		Amount:     src.Amount,
		Internal:   int32(src.Internal),
		OutTradeNo: src.OutTradeNo,
		PayTime:    src.PayTime,
	}
}

func (p *paymentService) SaveIntegrateApp(_ context.Context, app *proto.SIntegrateApp) (*proto.Result, error) {
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
	return p.error(err), nil
}

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
func (p *paymentService) DeleteIntegrateApp(_ context.Context, id *proto.PayIntegrateAppId) (*proto.Result, error) {
	err := p.repo.DeleteIntegrateApp(id.Value)
	return p.error(err), nil
}
