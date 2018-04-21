package rsi

/**
 * Copyright 2015 @ z3q.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
import (
	"context"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/module"
	"go2o/core/service/thrift/parser"
	"go2o/gen-code/thrift/define"
)

var _ define.PaymentService = new(paymentService)

type paymentService struct {
	_repo      payment.IPaymentRepo
	_orderRepo order.IOrderRepo
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo) *paymentService {
	return &paymentService{
		_repo:      rep,
		_orderRepo: orderRepo,
	}
}

// 根据编号获取支付单
func (p *paymentService) GetPaymentOrderById(ctx context.Context, id int32) (*define.SPaymentOrder, error) {
	po := p._repo.GetPaymentOrderById(id)
	if po != nil {
		v := po.GetValue()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 根据交易号获取支付单编号
func (p *paymentService) GetPaymentOrderId(ctx context.Context, tradeNo string) (int32, error) {
	po := p._repo.GetPaymentOrder(tradeNo)
	if po != nil {
		return po.GetAggregateRootId(), nil
	}
	return 0, nil
}

// 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(ctx context.Context, paymentNo string) (*define.SPaymentOrder, error) {
	if po := p._repo.GetPaymentOrder(paymentNo); po != nil {
		v := po.GetValue()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 创建支付单
func (p *paymentService) SubmitPaymentOrder(ctx context.Context, s *define.SPaymentOrder) (*define.Result_, error) {
	v := parser.PaymentOrder(s)
	o := p._repo.CreatePaymentOrder(v)
	err := o.Commit()
	return parser.Result(nil, err), nil
}

// 调整支付单金额
func (p *paymentService) AdjustOrder(ctx context.Context, paymentNo string, amount float64) (*define.Result_, error) {
	var err error
	o := p._repo.GetPaymentOrder(paymentNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.Adjust(float32(amount))
	}
	return parser.Result(0, err), nil
}

// 设置交易单号前缀
func (p *paymentService) SetPrefixOfTradeNo(id int32, prefix string) error {
	o := p._repo.GetPaymentOrderById(id)
	if o == nil {
		return payment.ErrNoSuchPaymentOrder
	}
	return o.TradeNoPrefix(prefix)
}

// 积分抵扣支付单
func (p *paymentService) DiscountByIntegral(ctx context.Context, orderId int32,
	integral int64, ignoreOut bool) (r *define.DResult_, err error) {
	var amount float32
	o := p._repo.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		amount, err = o.IntegralDiscount(integral, ignoreOut)
	}
	return parser.DResult(float64(amount), err), nil
}

// 余额抵扣
func (p *paymentService) DiscountByBalance(ctx context.Context, orderId int32, remark string) (*define.Result_, error) {
	var err error
	o := p._repo.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDiscount(remark)
	}
	return parser.Result(0, err), nil
}

// 钱包账户支付
func (p *paymentService) PaymentByWallet(ctx context.Context, orderId int32, remark string) (r *define.Result_, err error) {
	o := p._repo.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentByWallet(remark)
	}
	return parser.Result(0, err), nil
}

// 余额钱包混合支付，优先扣除余额。
func (p *paymentService) HybridPayment(ctx context.Context, orderId int32, remark string) (r *define.Result_, err error) {
	o := p._repo.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.HybridPayment(remark)
	}
	return parser.Result(0, err), nil
}

// 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(ctx context.Context, tradeNo string, spName string,
	outerNo string) (r *define.Result_, err error) {
	o := p._repo.GetPaymentOrder(tradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(spName, outerNo)
	}
	return parser.Result(nil, err), nil
}

// 支付网关
func (p *paymentService) GatewayV1(ctx context.Context, action string, userId int64, data map[string]string) (r *define.Result_, err error) {
	mod := module.Get(module.M_PAY).(*module.PaymentModule)
	// 获取令牌
	if action == "get_token" {
		token := mod.CreateToken(userId)
		return parser.Result(token, nil), nil
	}
	// 提交支付请求
	if action == "submit" {
		err = mod.Submit(userId, data)
	}
	// 验证支付
	if action == "payment" {
		err = mod.CheckAndPayment(userId, data)
	}
	return parser.Result(nil, err), nil
}
