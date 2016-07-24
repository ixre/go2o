/**
 * Copyright 2015 @ z3q.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
package dps

import (
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
)

type paymentService struct {
	_rep      payment.IPaymentRep
	_orderRep order.IOrderRep
}

func NewPaymentService(rep payment.IPaymentRep, orderRep order.IOrderRep) *paymentService {
	return &paymentService{
		_rep:      rep,
		_orderRep: orderRep,
	}
}

// 根据编号获取支付单
func (p *paymentService) GetPaymentOrder(id int) *payment.PaymentOrderBean {
	v := p._rep.GetPaymentOrder(id).GetValue()
	return &v
}

// 根据支付单号获取支付单
func (p *paymentService) GetPaymentOrderByNo(paymentNo string) *payment.PaymentOrderBean {
	if v := p._rep.GetPaymentOrderByNo(paymentNo); v != nil {
		v2 := v.GetValue()
		return &v2
	}
	return nil
}

// 创建支付单
func (p *paymentService) CreatePaymentOrder(v *payment.PaymentOrderBean,
) (int, error) {
	o := p._rep.CreatePaymentOrder(v)
	return o.Save()
}

// 创建支付单
func (p *paymentService) FinishPayment(tradeNo string, spName string,
	outerNo string) error {
	o := p._rep.GetPaymentOrderByNo(tradeNo)
	if o == nil {
		return payment.ErrNoSuchPaymentOrder
	}
	err := o.PaymentFinish(spName, outerNo)
	if err == nil {
		_, err = o.Save()
		//更改订单支付完成
		if err == nil {
			err = p._orderRep.Manager().PaymentForOnlineTrade(o.GetValue().OrderId)
		}
	}
	return err
}
