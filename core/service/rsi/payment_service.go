/**
 * Copyright 2015 @ z3q.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
package rsi

import (
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/module"
	"go2o/gen-code/thrift/define"
	"go2o/core/service/thrift/parser"
)

type paymentService struct {
	_rep       payment.IPaymentRepo
	_orderRepo order.IOrderRepo
}

func NewPaymentService(rep payment.IPaymentRepo, orderRepo order.IOrderRepo) *paymentService {
	return &paymentService{
		_rep:       rep,
		_orderRepo: orderRepo,
	}
}

// 根据编号获取支付单
func (p *paymentService) GetPaymentOrderById(id int32) (*define.PaymentOrder, error) {
	po := p._rep.GetPaymentOrderById(id)
	if po != nil {
		v := po.GetValue()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 根据交易号获取支付单编号
func (p *paymentService) GetPaymentOrderId(tradeNo string) (int32, error) {
	po := p._rep.GetPaymentOrder(tradeNo)
	if po != nil {
		return po.GetAggregateRootId(), nil
	}
	return 0, nil
}

// 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrder(paymentNo string) (*define.PaymentOrder, error) {
	if po := p._rep.GetPaymentOrder(paymentNo); po != nil {
		v := po.GetValue()
		return parser.PaymentOrderDto(&v), nil
	}
	return nil, nil
}

// 创建支付单
func (p *paymentService) SubmitPaymentOrder(s *define.PaymentOrder) (*define.Result_, error) {
	v := parser.PaymentOrder(s)
	o := p._rep.CreatePaymentOrder(v)
	return parser.Result(o.Commit()), nil
}

// 调整支付单金额
func (p *paymentService) AdjustOrder(paymentNo string, amount float64) (*define.Result_, error) {
	var err error
	o := p._rep.GetPaymentOrder(paymentNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.Adjust(float32(amount))
	}
	return parser.Result(0, err), nil
}

func (p *paymentService) SetPrefixOfTradeNo(id int32, prefix string) error {
	o := p._rep.GetPaymentOrderById(id)
	if o == nil {
		return payment.ErrNoSuchPaymentOrder
	}
	return o.TradeNoPrefix(prefix)
}

// 积分抵扣支付单
func (p *paymentService) DiscountByIntegral(orderId int32,
	integral int64, ignoreOut bool) (r *define.DResult_, err error) {
	var amount float32
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		amount, err = o.IntegralDiscount(integral, ignoreOut)
	}
	return parser.DResult(float64(amount), err), nil
}

// 余额抵扣
func (p *paymentService) DiscountByBalance(orderId int32, remark string) (*define.Result_, error) {
	var err error
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.BalanceDiscount(remark)
	}
	return parser.Result(0, err), nil
}

// 钱包账户支付
func (p *paymentService) PaymentByWallet(orderId int32, remark string) (r *define.Result_, err error) {
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentByWallet(remark)
	}
	return parser.Result(0, err), nil
}

// 余额钱包混合支付，优先扣除余额。
func (p *paymentService) HybridPayment(orderId int32, remark string) (r *define.Result_, err error) {
	o := p._rep.GetPaymentOrderById(orderId)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.HybridPayment(remark)
	}
	return parser.Result(0, err), nil
}

// 完成支付单支付，并传入支付方式及外部订单号
func (p *paymentService) FinishPayment(tradeNo string, spName string,
	outerNo string) (r *define.Result_, err error) {
	o := p._rep.GetPaymentOrder(tradeNo)
	if o == nil {
		err = payment.ErrNoSuchPaymentOrder
	} else {
		err = o.PaymentFinish(spName, outerNo)
	}
	return parser.Result(0, err), nil
}

// 支付网关
func (p *paymentService) GatewayV1(action string, userId int64, data map[string]string) (r *define.Result_, err error) {
	mod := module.Get(module.M_PAY).(*module.PaymentModule)
	rlt := &define.Result_{}
	// 提交支付请求
	if action == "submit" {
		err = mod.Submit(userId, data)
	}
	// 获取令牌
	if action == "get_token" {
		rlt.Message = mod.CreateToken(userId)
	}
	// 验证支付
	if action == "payment" {
		err = mod.CheckAndPayment(userId, data)
	}
	if err != nil {
		rlt.Result_ = false
		rlt.Message = err.Error()
	} else {
		rlt.Result_ = true
	}
	return rlt, nil
}
