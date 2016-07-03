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
func (this *paymentService) GetPaymentOrder(id int) payment.PaymentOrderBean {
	return this._rep.GetPaymentOrder(id).GetValue()
}

// 根据支付单号获取支付单
func (this *paymentService) GetPaymentOrderByNo(paymentNo string,
) payment.PaymentOrderBean {
	return this._rep.GetPaymentOrderByNo(paymentNo).GetValue()
}

// 创建支付单
func (this *paymentService) CreatePaymentOrder(v *payment.PaymentOrderBean,
) (int, error) {
	p := this._rep.CreatePaymentOrder(v)
	return p.Save()
}
