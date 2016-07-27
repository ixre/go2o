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

//todo: ??? save 改为submit
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

func (p *paymentOrderImpl) GetAggregateRootId() int {
	return p._value.Id
}

// 获取交易号
func (p *paymentOrderImpl) GetTradeNo() string {
	return p._value.TradeNo
}

// 重新修正金额
func (p *paymentOrderImpl) fixFee() {
	v := p._value
	v.FinalFee = v.TotalFee - v.CouponDiscount - v.BalanceDiscount -
		v.IntegralDiscount - v.SubAmount - v.SystemDiscount
}

// 更新订单状态, 需要注意,防止多次订单更新
func (p *paymentOrderImpl) notifyPaymentFinish() {
	err := p._rep.NotifyPaymentFinish(p.GetAggregateRootId())
	if err != nil {
		err = errors.New("Notify payment finish error :" + err.Error())
		domain.HandleError(err, "domain")
	}
	if p._value.OrderId > 0 {
		err = p._orderManager.PaymentForOnlineTrade(p._value.OrderId)
		domain.HandleError(err, "domain")
	}
}

// 优惠券抵扣

func (p *paymentOrderImpl) CouponDiscount(coupon promotion.ICouponPromotion) (
	float32, error) {
	if p._value.PaymentOpt&payment.OptUseCoupon == 0 {
		return 0, payment.ErrCanNotUseCoupon
	}
	//todo: 如可以使用多张优惠券,那么初始化应该获取支付单的所有优惠券
	if p._coupons == nil {
		p._coupons = []promotion.ICouponPromotion{}
	}
	p._coupons = append(p._coupons, coupon)
	// 支付金额应减去立减和系统支付的部分
	fee := p._value.TotalFee - p._value.SubAmount -
		p._value.SystemDiscount
	for _, v := range p._coupons {
		p._value.CouponDiscount += v.GetCouponFee(fee)
	}
	p.fixFee()
	return p._value.CouponDiscount, nil
}

// 在支付之前检查订单状态
func (p *paymentOrderImpl) checkPayment() error {
	if p.GetAggregateRootId() <= 0 {
		return payment.ErrPaymentNotSave
	}
	switch p._value.State {
	case payment.StateFinishPayment:
		return payment.ErrOrderPayed
	case payment.StateHasCancel:
		return payment.ErrOrderHasCancel
	}
	return nil
}

// 应用余额支付
func (p *paymentOrderImpl) getBalanceDiscountFee(acc member.IAccount) float32 {
	if p._value.FinalFee <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= p._value.FinalFee {
		return p._value.FinalFee
	} else {
		return acv.Balance
	}
	return 0
}

// 使用余额支付
func (p *paymentOrderImpl) paymentWithBalance(buyerType int) error {
	if p._value.PaymentOpt&payment.OptBalanceDiscount == 0 {
		return payment.ErrCanNotUseBalance
	}
	err := p.checkPayment()
	if err == nil {
		// 判断扣减金额,是否大于0
		acc := p._mmRep.GetMember(p._value.BuyUser).GetAccount()
		fee := p.getBalanceDiscountFee(acc)
		if fee == 0 {
			return member.ErrAccountBalanceNotEnough
		}
		// 从会员账户扣减,并更新支付单
		err = acc.PaymentDiscount(p.GetValue().TradeNo, fee)
		if err == nil {
			p._value.BalanceDiscount = fee
			p.fixFee()
		}
	}
	return err
}

// 检查是否支付完成, 且返回是否为第一次支付成功,
func (p *paymentOrderImpl) checkPaymentOk() (bool, error) {
	b := false
	if p._value.State == payment.StateAwaitingPayment {
		unix := time.Now().Unix()
		// 如果支付完成,则更新订单状态
		if b = p._value.FinalFee == 0; b {
			p._value.State = payment.StateFinishPayment
			p._firstFinishPayment = true
		}
		p._value.PaidTime = unix
	}
	return b, nil
}

// 使用会员的余额抵扣
func (p *paymentOrderImpl) BalanceDiscount() error {
	return p.paymentWithBalance(payment.PaymentByBuyer)
}

// 计算积分折算后的金额
func (p *paymentOrderImpl) mathIntegralFee(integral int) float32 {
	if integral > 0 {
		conf := p._valRep.GetGlobNumberConf()
		if conf.IntegralDiscountRate > 0 {
			return float32(integral) / conf.IntegralDiscountRate
		}
	}
	return 0
}

// 积分抵扣,返回抵扣的金额及错误,ignoreAmount:是否忽略超出订单金额的积分
func (p *paymentOrderImpl) IntegralDiscount(integral int, ignoreAmount bool) (float32, error) {
	var amount float32 = 0
	if p._value.PaymentOpt&payment.OptIntegralDiscount == 0 {
		return 0, payment.ErrCanNotUseIntegral
	}
	err := p.checkPayment()
	if err != nil {
		return 0, err
	}
	// 判断扣减金额是否大于0
	amount = p.mathIntegralFee(integral)
	// 如果不忽略超出订单支付金额的积分,那么按实际来抵扣
	if !ignoreAmount && amount > p._value.FinalFee {
		conf := p._valRep.GetGlobNumberConf()
		amount = p._value.FinalFee
		integral = int(amount * conf.IntegralDiscountRate)
	}
	if amount > 0 {
		acc := p._mmRep.GetMember(p._value.BuyUser).GetAccount()
		// 抵扣积分

		//log.Println("----", p._value.BuyUser, acc.GetValue().Integral, "discount:", integral)
		//log.Printf("-----%#v\n", acc.GetValue())
		err = acc.IntegralDiscount(member.TypeIntegralPaymentDiscount,
			p.GetValue().TradeNo, integral, "")
		if err == nil {
			p._value.IntegralDiscount = amount
			p.fixFee()
			_, err = p.Save()
		}
	}
	return amount, err
}

// 系统支付金额
func (p *paymentOrderImpl) SystemPayment(fee float32) error {
	if p._value.PaymentOpt&payment.OptSystemPayment == 0 {
		return payment.ErrCanNotSystemDiscount
	}
	err := p.checkPayment()
	if err == nil {
		p._value.SystemDiscount += fee
		p.fixFee()
	}
	return err
}

// 设置支付方式
func (p *paymentOrderImpl) SetPaymentSign(paymentSign int) error {
	//todo: 某个支付方式被暂停
	p._value.PaymentSign = paymentSign
	return nil
}

// 绑定订单号,如果交易号为空则绑定参数中传递的交易号
func (p *paymentOrderImpl) BindOrder(orderId int, tradeNo string) error {
	//todo: check order exists  and tradeNo exists
	p._value.OrderId = orderId
	if len(p._value.TradeNo) == 0 {
		p._value.TradeNo = tradeNo
	}
	return nil
}

func (p *paymentOrderImpl) Save() (int, error) {
	return p.save()
}

func (p *paymentOrderImpl) save() (int, error) {
	_, err := p.checkPaymentOk()
	if err == nil {
		unix := time.Now().Unix()
		if p._value.CreateTime == 0 {
			p._value.CreateTime = unix
		}
		p._value.Id, err = p._rep.SavePaymentOrder(p._value)
	}

	//保存支付单后,通知支付成功。只通知一次
	if err == nil && p._firstFinishPayment {
		p._firstFinishPayment = false
		go p.notifyPaymentFinish()
	}
	return p.GetAggregateRootId(), err
}

// 支付完成,传入第三名支付名称,以及外部的交易号
func (p *paymentOrderImpl) PaymentFinish(spName string, outerNo string) error {
	outerNo = strings.TrimSpace(outerNo)
	if len(outerNo) < 8 {
		return payment.ErrOuterNo
	}
	if p._value.State == payment.StateFinishPayment {
		return payment.ErrOrderPayed
	}
	if p._value.State == payment.StateHasCancel {
		return payment.ErrOrderHasCancel
	}
	p._value.State = payment.StateFinishPayment
	p._value.OuterNo = outerNo
	p._value.PaidTime = time.Now().Unix()
	p._firstFinishPayment = true

	return nil
}
func (p *paymentOrderImpl) GetValue() payment.PaymentOrderBean {
	return *p._value
}

// 取消支付
func (p *paymentOrderImpl) Cancel() error {
	p._value.State = payment.StateHasCancel
	return nil
}

// 调整金额,如果调整的金额与实付金额一致,则取消支付单
func (p *paymentOrderImpl) Adjust(amount float32) error {
	p._value.AdjustmentAmount += amount
	p._value.FinalFee += amount
	if p._value.FinalFee <= 0 {
		_, err := p.checkPaymentOk()
		return err
	}
	_, err := p.Save()
	return err
}

type PaymentRepBase struct {
}

func (p *PaymentRepBase) CreatePaymentOrder(v *payment.
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
