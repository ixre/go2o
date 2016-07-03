/**
 * Copyright 2015 @ z3q.net.
 * name : payment_service.go
 * author : jarryliu
 * date : 2016-07-03 13:24
 * description :
 * history :
 */
package dps

import "go2o/core/domain/interface/payment"

type paymentService struct {
	_rep payment.IPaymentRep
}

func NewPaymentService(rep payment.IPaymentRep) *paymentService {
	return &paymentService{
		_rep: rep,
	}
}

// 根据编号获取支付单
func (this *paymentService) GetPaymentOrder(id int) *payment.PaymentOrderBean {
	v := this._rep.GetPaymentOrder(id).GetValue()
	return &v
}

// 根据支付单号获取支付单
func (this *paymentService) GetPaymentOrderByNo(paymentNo string) *payment.PaymentOrderBean {
	if v := this._rep.GetPaymentOrderByNo(paymentNo); v != nil {
		v2 := v.GetValue()
		return &v2
	}
	return nil
}

// 创建支付单
func (this *paymentService) CreatePaymentOrder(v *payment.PaymentOrderBean,
) (int, error) {
	p := this._rep.CreatePaymentOrder(v)
	return p.Save()
}

// 创建支付单
func (this *paymentService) FinishPayment(tradeNo string, spName string,
	outerNo string) error {
	o := this._rep.GetPaymentOrderByNo(tradeNo)
	if o == nil {
		return payment.ErrNoSuchPaymentOrder
	}
	err := o.PaymentFinish(spName, outerNo)
	if err == nil {
		_, err = o.Save()
	}
	return err
}
