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
	"errors"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"strings"
	"time"
)

var _ payment.IPaymentOrder = new(paymentOrderImpl)

type paymentOrderImpl struct {
	_rep                payment.IPaymentRep
	_value              *payment.PaymentOrderBean
	_mmRep              member.IMemberRep
	_valRep             valueobject.IValueRep
	_coupons            []promotion.ICouponPromotion
	_orderManager       order.IOrderManager
	_firstFinishPayment bool //第一次完成支付
}

func (this *paymentOrderImpl) GetAggregateRootId() int {
	return this._value.Id
}

// 获取交易号
func (this *paymentOrderImpl) GetTradeNo() string {
	return this._value.TradeNo
}

// 重新修正金额
func (this *paymentOrderImpl) fixFee() {
	v := this._value
	v.FinalFee = v.TotalFee - v.CouponDiscount - v.BalanceDiscount -
		v.IntegralDiscount - v.SubFee - v.SystemDiscount
}

// 更新订单状态, 需要注意,防止多次订单更新
func (this *paymentOrderImpl) notifyPaymentFinish() {

	err := this._rep.NotifyPaymentFinish(this.GetAggregateRootId())
	if err != nil {
		err = errors.New("Notify payment finish error :" + err.Error())
		domain.HandleError(err, "domain")
	}
	if this._value.OrderId > 0 {
		err = this._orderManager.PaymentForOnlineTrade(this._value.OrderId)
		if err != nil {
			domain.HandleError(err, "domain")
		}
	}

	//todo:  更新订单状态

	//this._value.PaymentSign = buyerType
	//if this._value.Status == enum.ORDER_WAIT_PAYMENT {
	//    this._value.Status = enum.ORDER_WAIT_CONFIRM
	//}
}

// 优惠券抵扣

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
		this._value.CouponDiscount += v.GetCouponFee(fee)
	}
	this.fixFee()
	return this._value.CouponDiscount, nil
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
func (this *paymentOrderImpl) paymentWithBalance(buyerType int) error {
	if this._value.PaymentOpt&payment.OptBalanceDiscount == 0 {
		return payment.ErrCanNotUseBalance
	}
	err := this.checkPayment()
	if err == nil {
		// 判断扣减金额,是否大于0
		acc := this._mmRep.GetMember(this._value.BuyUser).GetAccount()
		fee := this.getBalanceDiscountFee(acc)
		if fee == 0 {
			return member.ErrAccountBalanceNotEnough
		}
		// 从会员账户扣减,并更新支付单
		err = acc.PaymentDiscount(this.GetValue().TradeNo, fee)
		if err == nil {
			this._value.BalanceDiscount = fee
			this.fixFee()
		}
	}
	return err
}

// 检查是否支付完成, 且返回是否为第一次支付成功,
func (this *paymentOrderImpl) checkPaymentOk() (bool, error) {
	b := false
	if this._value.State == payment.StateNotYetPayment {
		unix := time.Now().Unix()
		// 如果支付完成,则更新订单状态
		if b = this._value.FinalFee == 0; b {
			this._value.State = payment.StateFinishPayment
			this._firstFinishPayment = true
		}
		this._value.PaidTime = unix
	}
	return b, nil
}

// 使用会员的余额抵扣
func (this *paymentOrderImpl) BalanceDiscount() error {
	return this.paymentWithBalance(payment.PaymentByBuyer)
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
		err = acc.DiscountIntegral(this.GetValue().TradeNo, integral, fee)
		if err == nil {
			this._value.IntegralDiscount = fee
			this.fixFee()
		}
	}
	return err
}

// 系统支付金额
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

// 设置支付方式
func (this *paymentOrderImpl) SetPaymentSign(paymentSign int) error {
	//todo: 某个支付方式被暂停
	this._value.PaymentSign = paymentSign
	return nil
}

// 绑定订单号,如果交易号为空则绑定参数中传递的交易号
func (this *paymentOrderImpl) BindOrder(orderId int, tradeNo string) error {
	//todo: check order exists  and tradeNo exists
	this._value.OrderId = orderId
	if len(this._value.TradeNo) == 0 {
		this._value.TradeNo = tradeNo
	}
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

	//保存支付单后,通知支付成功。只通知一次
	if err == nil && this._firstFinishPayment {
		this._firstFinishPayment = false
		this.notifyPaymentFinish()
	}
	return this.GetAggregateRootId(), err
}

// 支付完成,传入第三名支付名称,以及外部的交易号
func (this *paymentOrderImpl) PaymentFinish(spName string, outerNo string) error {
	outerNo = strings.TrimSpace(outerNo)
	if len(outerNo) < 8 {
		return payment.ErrOuterNo
	}
	if this._value.State == payment.StateFinishPayment {
		return payment.ErrOrderPayed
	}
	if this._value.State == payment.StateHasCancel {
		return payment.ErrOrderHasCancel
	}
	this._value.State = payment.StateFinishPayment
	this._value.OuterNo = outerNo
	this._value.PaidTime = time.Now().Unix()
	this._firstFinishPayment = true

	return nil
}
func (this *paymentOrderImpl) GetValue() payment.PaymentOrderBean {
	return *this._value
}

// 取消支付
func (this *paymentOrderImpl) Cancel() error {
	this._value.State = payment.StateHasCancel
	return nil
}

type PaymentRepBase struct {
}

func (this *PaymentRepBase) CreatePaymentOrder(v *payment.
	PaymentOrderBean, rep payment.IPaymentRep, mmRep member.IMemberRep,
	orderManager order.IOrderManager, valRep valueobject.IValueRep) payment.IPaymentOrder {
	return &paymentOrderImpl{
		_rep:          rep,
		_value:        v,
		_mmRep:        mmRep,
		_valRep:       valRep,
		_orderManager: orderManager,
	}
}
