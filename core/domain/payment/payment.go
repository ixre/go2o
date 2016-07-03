/**
 * Copyright 2015 @ z3q.net.
 * name : payment
 * author : jarryliu
 * date : 2016-07-03 09:25
 * description :
 * history :
 */
package payment

import (
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"time"
)

var _ payment.IPaymentOrder = new(paymentOrderImpl)

type paymentOrderImpl struct {
	_rep     payment.IPaymentRep
	_value   *payment.PaymentOrderBean
	_mmRep   member.IMemberRep
	_valRep  valueobject.IValueRep
	_coupons []promotion.ICouponPromotion
}

func (this *paymentOrderImpl) GetAggregateRootId() int {
	return this._value.Id
}

// 重新修正金额
func (this *paymentOrderImpl) fixFee() {
	v := this._value
	v.FinalFee = v.TotalFee - v.CouponFee - v.BalanceDiscount -
		v.IntegralDiscount - v.SubFee - v.SystemDiscount
}

// 更新订单状态
func (this *paymentOrderImpl) updateOrderFinish() {
	panic("未实现")
	//todo:  更新订单状态

	//this._value.PaymentSign = buyerType
	//if this._value.Status == enum.ORDER_WAIT_PAYMENT {
	//    this._value.Status = enum.ORDER_WAIT_CONFIRM
	//}
}

/// <summary>
/// 优惠券抵扣
/// </summary>
func (this *paymentOrderImpl) CouponDiscount(coupon promotion.ICouponPromotion) (
	float32, error) {
	if this._value.PaymentOpt&payment.OptUseCoupon == 0 {
		return 0, payment.ErrCanNotUseCoupon
	}
	//todo: 如可以使用多张优惠券,那么初始化应该获取支付单的所有优惠券
	if this._coupons == nil {
		this._coupons = []promotion.ICouponPromotion{}
	}
	this._coupons = append(this._coupons, coupon)
	// 支付金额应减去立减和系统支付的部分
	fee := this._value.TotalFee - this._value.SubFee -
		this._value.SystemDiscount
	for _, v := range this._coupons {
		this._value.CouponFee += v.GetCouponFee(fee)
	}
	this.fixFee()
	return this._value.CouponFee, nil
}

// 在支付之前检查订单状态
func (this *paymentOrderImpl) checkPayment() error {
	if this.GetAggregateRootId() <= 0 {
		return payment.ErrPaymentNotSave
	}
	switch this._value.State {
	case payment.StateFinishPayment:
		return payment.ErrOrderPayed
	case payment.StateHasCancel:
		return payment.ErrOrderHasCancel
	}
	return nil
}

// 应用余额支付
func (this *paymentOrderImpl) getBalanceDiscountFee(acc member.IAccount) float32 {
	if this._value.FinalFee <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= this._value.FinalFee {
		return this._value.FinalFee
	} else {
		return acv.Balance
	}
	return 0
}

// 使用余额支付
func (this *paymentOrderImpl) paymentWithBalance(buyerType int, fee float32) error {
	if this._value.PaymentOpt&payment.OptBalanceDiscount == 0 {
		return payment.ErrCanNotUseBalance
	}
	err := this.checkPayment()
	if err == nil {
		// 判断扣减金额,是否大于0
		acc := this._mmRep.GetMember(this._value.BuyUser).GetAccount()
		if fee := this.getBalanceDiscountFee(acc); fee == 0 {
			return member.ErrAccountBalanceNotEnough
		}
		// 从会员账户扣减,并更新支付单
		err = acc.PaymentDiscount(this.GetValue().PaymentNo, fee)
		if err == nil {
			this._value.BalanceDiscount = fee
			this.fixFee()
		}
	}
	return err
}

func (this *paymentOrderImpl) checkPaymentOk() (bool, error) {
	err := this.checkPayment()
	b := false
	if err == nil {
		unix := time.Now().Unix()
		// 如果支付完成,则更新订单状态
		if b = this._value.FinalFee == 0; b {
			this._value.State = payment.StateFinishPayment
			this.updateOrderFinish()
		}
		this._value.PaidTime = unix
	}
	return b, err
}

// 使用余额支付
func (this *paymentOrderImpl) BalanceDiscount(fee float32) error {
	return this.paymentWithBalance(payment.PaymentByBuyer, fee)
}

// 计算积分折算后的金额
func (this *paymentOrderImpl) mathIntegralFee(integral int) float32 {
	if integral > 0 {
		conf := this._valRep.GetGlobNumberConf()
		if conf.IntegralExchangeRate > 0 {
			return float32(integral) / float32(conf.IntegralExchangeRate)
		}
	}
	return 0
}

// 积分抵扣
func (this *paymentOrderImpl) IntegralDiscount(integral int) error {
	if this._value.PaymentOpt&payment.OptIntegralDiscount == 0 {
		return payment.ErrCanNotUseIntegral
	}

	err := this.checkPayment()
	if err == nil {
		// 判断扣减金额,是否大于0
		acc := this._mmRep.GetMember(this._value.BuyUser).GetAccount()
		if acc.GetValue().Integral < integral {
			return member.ErrAccountBalanceNotEnough
		}
		fee := this.mathIntegralFee(integral)
		err = acc.DiscountIntegral(this.GetValue().PaymentNo, integral, fee)
		if err == nil {
			this._value.IntegralDiscount = fee
			this.fixFee()
		}
	}
	return err
}

/// <summary>
/// 系统支付金额
/// </summary>
func (this *paymentOrderImpl) SystemPayment(fee float32) error {
	if this._value.PaymentOpt&payment.OptSystemPayment == 0 {
		return payment.ErrCanNotSystemDiscount
	}
	err := this.checkPayment()
	if err == nil {
		this._value.SystemDiscount += fee
		this.fixFee()
	}
	return err
}
func (this *paymentOrderImpl) BindOrder(orderId int) error {
	//todo: check order exists
	this._value.OrderId = orderId
	return nil
}
func (this *paymentOrderImpl) Save() (int, error) {
	_, err := this.checkPaymentOk()
	if err == nil {
		unix := time.Now().Unix()
		if this._value.CreateTime == 0 {
			this._value.CreateTime = unix
		}
		this._value.Id, err = this._rep.SavePaymentOrder(this._value)
	}
	return this.GetAggregateRootId(), err
}

func (this *paymentOrderImpl) PaymentFinish(tradeNo string) error {
	err := this.checkPayment()
	if err == nil {
		this._value.TradeNo = tradeNo
		this._value.FinalFee = 0
	}
	return err
}
func (this *paymentOrderImpl) GetValue() payment.PaymentOrderBean {
	return *this._value
}

/// <summary>
/// 取消支付
/// </summary>
func (this *paymentOrderImpl) Cancel() error {
	this._value.State = payment.StateHasCancel
	return nil
}

type PaymentRepBase struct {
}

func (this *PaymentRepBase) CreatePaymentOrder(v *payment.
	PaymentOrderBean, rep payment.IPaymentRep, mmRep member.IMemberRep,
	valRep valueobject.IValueRep) payment.IPaymentOrder {
	return &paymentOrderImpl{
		_rep:    rep,
		_value:  v,
		_mmRep:  mmRep,
		_valRep: valRep,
	}
}
