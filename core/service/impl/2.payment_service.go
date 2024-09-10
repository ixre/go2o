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
			pm.ChanName = po.ChanName(t.Method)
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

func (p *paymentService) parseTradeMethodDataDto(src *payment.TradeMethodData) *proto.STradeChanData {
	return &proto.STradeChanData{
		ChanId:     int32(src.Method),
		Amount:     src.Amount,
		ChanCode:   src.Code,
		OutTradeNo: src.OutTradeNo,
	}
}

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
func (p *paymentService) DeleteIntegrateApp(_ context.Context, id *proto.PayIntegrateAppId) (*proto.TxResult, error) {
	err := p.repo.DeleteIntegrateApp(id.Value)
	return p.errorV2(err), nil
}

// Divide implements proto.PaymentServiceServer.
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

// FinishDivide implements proto.PaymentServiceServer.
func (p *paymentService) FinishDivide(_ context.Context, req *proto.PaymentOrderRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrder(req.TradeNo)
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	err := ip.FinishDivide()
	return p.errorV2(err), nil
}

func (p *paymentService) UpdateDivideStatus(_ context.Context, req *proto.UpdateDivideStatusRequest) (*proto.TxResult, error) {
	ip := p.repo.GetPaymentOrderById(int(req.PayId))
	if ip == nil {
		return p.errorV2(payment.ErrNoSuchPaymentOrder), nil
	}
	err := ip.UpdateDivideStatus(int(req.DivideId), req.Success, req.DivideNo, req.Remark)
	return p.errorV2(err), nil
}
