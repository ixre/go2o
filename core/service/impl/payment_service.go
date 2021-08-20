package impl

/**
 * Copyright 2015 @ to2.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
import (
	"context"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/module"
	"go2o/core/service/proto"
	"strconv"
)

var _ proto.PaymentServiceServer = new(paymentService)

type paymentService struct {
	repo       payment.IPaymentRepo
	orderRepo  order.IOrderRepo
	memberRepo member.IMemberRepo
	serviceUtil
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo,
	memberRepo member.IMemberRepo) *paymentService {
	return &paymentService{
		repo:       rep,
		orderRepo:  orderRepo,
		memberRepo: memberRepo,
	}
}

// 根据编号获取支付单
func (p *paymentService) GetPaymentOrderById(_ context.Context, id *proto.Int32) (*proto.SPaymentOrder, error) {
	po := p.repo.GetPaymentOrderById(int(id.Value))
	if po != nil {
		v := po.Get()
		return p.parsePaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 根据交易号获取支付单编号
func (p *paymentService) GetPaymentOrderId(_ context.Context, tradeNo *proto.String) (*proto.Int32, error) {
	po := p.repo.GetPaymentOrder(tradeNo.Value)
	if po != nil {
		return &proto.Int32{Value: int32(po.GetAggregateRootId())}, nil
	}
	return &proto.Int32{Value: 0}, nil
}

// 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(_ context.Context, paymentNo *proto.String) (*proto.SPaymentOrder, error) {
	if po := p.repo.GetPaymentOrder(paymentNo.Value); po != nil {
		v := po.Get()
		sp := p.parsePaymentOrderDto(&v)
		for _, t := range po.TradeMethods() {
			sp.TradeData = append(sp.TradeData, p.parseTradeMethodDataDto(t))
		}
		return sp, nil
	}
	return nil, nil
}

// 创建支付单
func (p *paymentService) SubmitPaymentOrder(_ context.Context, order *proto.SPaymentOrder) (*proto.Result, error) {
	v := p.parsePaymentOrder(order)
	o := p.repo.CreatePaymentOrder(v)
	err := o.Submit()
	return p.result(err), nil
}

// 调整支付单金额
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

// 积分抵扣支付单
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

// 余额抵扣
func (p *paymentService) DiscountByBalance(_ context.Context, r *proto.DiscountBalanceRequest) (*proto.Result, error) {
	var err error
	o := p.repo.GetPaymentOrderById(int(r.OrderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDiscount(r.Remark)
	}
	return p.result(err), nil
}

// 钱包账户支付
func (p *paymentService) PaymentByWallet(_ context.Context, r *proto.WalletPaymentRequest) (rs *proto.Result, err error) {
	// 单个支付订单支付
	if !r.MergePay {
		ip := p.repo.GetPaymentOrder(r.TradeNo)
		if ip == nil {
			err = payment.ErrNoSuchPaymentOrder
		} else {
			err = ip.PaymentByWallet(r.Remark)
		}
		return p.result(err), nil
	}
	// 合并支付单
	arr := p.repo.GetMergePayOrders(r.TradeNo)
	payUid := arr[0].Get().PayUid
	finalFee := 0
	for _, v := range arr {
		finalFee += v.Get().FinalFee
	}
	acc := p.memberRepo.GetAccount(payUid)
	if int(acc.Balance*100) < finalFee {
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

// 余额钱包混合支付，优先扣除余额。
func (p *paymentService) HybridPayment(_ context.Context, r *proto.HyperPaymentRequest) (rs *proto.Result, err error) {
	o := p.repo.GetPaymentOrderById(int(r.OrderId))
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.HybridPayment(r.Remark)
	}
	return p.result(err), nil
}

// 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(_ context.Context, r *proto.FinishPaymentRequest) (rs *proto.Result, err error) {
	o := p.repo.GetPaymentOrder(r.TradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(r.SpName, r.OuterNo)
	}
	return p.result(err), nil
}

// 支付网关
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

// 获取支付预交易数据
func (p *paymentService) GetPaymentOrderInfo(_ context.Context, r *proto.OrderInfoRequest) (*proto.SPrepareTradeData, error) {
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
			ProcedureFee: int32(iv.ProcedureFee),
			FinalFee:     int32(iv.FinalFee),
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
		d.ProcedureFee += so.ProcedureFee      // 手续费
		d.FinalFee += so.FinalFee              // 最终金额
		d.TotalAmount += int32(iv.TotalAmount) // 累计金额
	}
	d.ErrCode = 0
	d.TradeNo = tradeNo // 交易单号
	return d, nil
}

// 混合支付
func (p *paymentService) MixedPayment(_ context.Context, _ *proto.MixedPaymentRequest) (*proto.Result, error) {
	return nil, nil
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

func (p *paymentService) parsePaymentOrderDto(src *payment.Order) *proto.SPaymentOrder {
	return &proto.SPaymentOrder{
		Id:             int32(src.Id),
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
		SubmitTime:     src.SubmitTime,
		ExpiresTime:    src.ExpiresTime,
		PaidTime:       src.PaidTime,
		SubOrder:       src.SubOrder == 1,
		OrderType:      int32(src.OrderType),
		OutOrderNo:     src.OutOrderNo,
		TradeData:      make([]*proto.STradeMethodData, 0),
	}
}

func (p *paymentService) parseTradeMethodDataDto(src *payment.TradeMethodData) *proto.STradeMethodData {
	return &proto.STradeMethodData{
		Method:     int32(src.Method),
		Code:       src.Code,
		Amount:     int32(src.Amount),
		Internal:   int32(src.Internal),
		OutTradeNo: src.OutTradeNo,
		PayTime:    src.PayTime,
	}
}
