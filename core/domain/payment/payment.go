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
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"regexp"
	"strings"
	"time"
)

var _ payment.IPaymentOrder = new(paymentOrderImpl)
var (
	letterRegexp        = regexp.MustCompile("^[A-Z]+$")
	tradeNoPrefixRegexp = regexp.MustCompile("^[A-Za-z]+\\d+$")
)

type paymentOrderImpl struct {
	rep                payment.IPaymentRep
	value              *payment.PaymentOrder
	mmRep              member.IMemberRep
	valRep             valueobject.IValueRep
	coupons            []promotion.ICouponPromotion
	orderManager       order.IOrderManager
	firstFinishPayment bool //第一次完成支付
	paymentUser        member.IMember
	buyer              member.IMember
}

func (p *paymentOrderImpl) GetAggregateRootId() int64 {
	return p.value.Id
}

// 获取交易号
func (p *paymentOrderImpl) GetTradeNo() string {
	return p.value.TradeNo
}

// 为交易号增加一个2位的前缀
func (p *paymentOrderImpl) TradeNoPrefix(prefix string) error {
	if tradeNoPrefixRegexp.MatchString(p.value.TradeNo) {
		return payment.ErrTradeNoExistsPrefix
	}
	if !letterRegexp.MatchString(prefix) {
		return payment.ErrTradeNoPrefix
	}
	p.value.TradeNo = prefix + p.value.TradeNo
	_, err := p.save()
	return err
}

// 重新修正金额
func (p *paymentOrderImpl) fixFee() {
	v := p.value
	v.FinalAmount = v.TotalFee - v.CouponDiscount - v.BalanceDiscount -
		v.IntegralDiscount - v.SubAmount - v.SystemDiscount
}

// 更新订单状态, 需要注意,防止多次订单更新
func (p *paymentOrderImpl) notifyPaymentFinish() {
	if p.GetAggregateRootId() <= 0 {
		panic(payment.ErrNoSuchPaymentOrder)
	}
	//err := p._rep.NotifyPaymentFinish(p.GetAggregateRootId())
	//if err != nil {
	//	err = errors.New("Notify payment finish error :" + err.Error())
	//	domain.HandleError(err, "domain")
	//}

	// 通知订单支付完成
	if p.value.OrderId > 0 {
		err := p.orderManager.PaymentForOnlineTrade(p.value.OrderId)
		domain.HandleError(err, "domain")
	}
}

// 优惠券抵扣

func (p *paymentOrderImpl) CouponDiscount(coupon promotion.ICouponPromotion) (
	float32, error) {
	if p.value.PaymentSign&payment.OptUseCoupon == 0 {
		return 0, payment.ErrCanNotUseCoupon
	}
	//todo: 如可以使用多张优惠券,那么初始化应该获取支付单的所有优惠券
	if p.coupons == nil {
		p.coupons = []promotion.ICouponPromotion{}
	}
	p.coupons = append(p.coupons, coupon)
	// 支付金额应减去立减和系统支付的部分
	fee := p.value.TotalFee - p.value.SubAmount -
		p.value.SystemDiscount
	for _, v := range p.coupons {
		p.value.CouponDiscount += v.GetCouponFee(fee)
	}
	p.fixFee()
	return p.value.CouponDiscount, nil
}

// 在支付之前检查订单状态
func (p *paymentOrderImpl) checkPayment() error {
	if p.GetAggregateRootId() <= 0 {
		return payment.ErrPaymentNotSave
	}
	switch p.value.State {
	case payment.StateAwaitingPayment:
		if p.value.FinalAmount == 0 {
			return payment.ErrFinalFee
		}
	case payment.StateFinishPayment:
		return payment.ErrOrderPayed
	case payment.StateHasCancel:
		return payment.ErrOrderHasCancel
	}
	return nil
}

// 应用余额支付
func (p *paymentOrderImpl) getBalanceDiscountAmount(acc member.IAccount) float32 {
	if p.value.FinalAmount <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= p.value.FinalAmount {
		return p.value.FinalAmount
	} else {
		return acv.Balance
	}
	return 0
}

func (p *paymentOrderImpl) getPaymentUser() member.IMember {
	if p.paymentUser == nil && p.value.PaymentUser > 0 {
		p.paymentUser = p.mmRep.GetMember(p.value.PaymentUser)
	}
	return p.paymentUser
}

// 使用余额支付
func (p *paymentOrderImpl) paymentWithBalance(buyerType int, remark string) error {
	if p.value.PaymentSign&payment.OptBalanceDiscount == 0 {
		return payment.ErrCanNotUseBalance
	}
	err := p.checkPayment()
	if err == nil {
		// 判断扣减金额,是否大于0
		pu := p.getPaymentUser()
		if pu == nil {
			return member.ErrNoSuchMember
		}
		acc := pu.GetAccount()
		amount := p.getBalanceDiscountAmount(acc)
		if amount == 0 {
			return member.ErrAccountBalanceNotEnough
		}
		// 从会员账户扣减,并更新支付单
		err = acc.PaymentDiscount(p.value.TradeNo, amount, remark)
		if err == nil {
			p.value.BalanceDiscount = amount
			p.fixFee()
			_, err = p.save()
		}
	}
	return err
}

// 检查是否支付完成, 且返回是否为第一次支付成功,
func (p *paymentOrderImpl) checkPaymentOk() (bool, error) {
	b := false
	if p.value.State == payment.StateAwaitingPayment {
		unix := time.Now().Unix()
		// 如果支付完成,则更新订单状态
		if b = p.value.FinalAmount == 0; b {
			p.value.State = payment.StateFinishPayment
			p.firstFinishPayment = true
		}
		p.value.PaidTime = unix
	}
	return b, nil
}

// 使用会员的余额抵扣
func (p *paymentOrderImpl) BalanceDiscount(remark string) error {
	return p.paymentWithBalance(payment.PaymentByBuyer, remark)
}

// 计算积分折算后的金额
func (p *paymentOrderImpl) mathIntegralFee(integral int) float32 {
	if integral > 0 {
		conf := p.valRep.GetGlobNumberConf()
		if conf.IntegralDiscountRate > 0 {
			return float32(integral) / conf.IntegralDiscountRate
		}
	}
	return 0
}

// 积分抵扣,返回抵扣的金额及错误,ignoreAmount:是否忽略超出订单金额的积分
func (p *paymentOrderImpl) IntegralDiscount(integral int, ignoreAmount bool) (float32, error) {
	var amount float32 = 0
	if p.value.PaymentSign&payment.OptIntegralDiscount != payment.OptIntegralDiscount {
		return 0, payment.ErrCanNotUseIntegral
	}
	err := p.checkPayment()
	if err != nil {
		return 0, err
	}
	// 判断扣减金额是否大于0
	amount = p.mathIntegralFee(integral)
	// 如果不忽略超出订单支付金额的积分,那么按实际来抵扣
	if !ignoreAmount && amount > p.value.FinalAmount {
		conf := p.valRep.GetGlobNumberConf()
		amount = p.value.FinalAmount
		integral = int(amount * conf.IntegralDiscountRate)
	}

	if amount > 0 {
		acc := p.mmRep.GetMember(p.value.BuyUser).GetAccount()
		// 抵扣积分

		//log.Println("----", p._value.BuyUser, acc.GetValue().Integral, "discount:", integral)
		//log.Printf("-----%#v\n", acc.GetValue())
		err = acc.IntegralDiscount(member.TypeIntegralPaymentDiscount,
			p.GetValue().TradeNo, integral, "")
		if err == nil {
			p.value.IntegralDiscount = amount
			p.fixFee()
			_, err = p.save()
		}
	}
	return amount, err
}

// 系统支付金额
func (p *paymentOrderImpl) SystemPayment(fee float32) error {
	if p.value.PaymentSign&payment.OptSystemPayment == 0 {
		return payment.ErrCanNotSystemDiscount
	}
	err := p.checkPayment()
	if err == nil {
		p.value.SystemDiscount += fee
		p.fixFee()
	}
	return err
}

func (p *paymentOrderImpl) getBuyer() member.IMember {
	if p.buyer == nil {
		p.buyer = p.mmRep.GetMember(p.value.BuyUser)
	}
	return p.buyer
}

// 赠送账户支付
func (p *paymentOrderImpl) PresentAccountPayment(remark string) error {
	amount := p.value.FinalAmount
	buyer := p.getBuyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	acc := buyer.GetAccount()
	av := acc.GetValue()
	if av.PresentBalance < amount {
		return payment.ErrNotEnoughAmount
	}
	if remark == "" {
		remark = "支付订单"
	}
	err := acc.DiscountPresent(remark, p.GetTradeNo(), amount,
		member.DefaultRelateUser, true)
	if err == nil {
		p.value.PaymentSign = payment.SignPresentAccount
		p.value.FinalAmount = 0
		p.value.PaidTime = time.Now().Unix()
		_, err = p.save()
	}
	return err
}

// 设置支付方式
func (p *paymentOrderImpl) SetPaymentSign(paymentSign int) error {
	//todo: 某个支付方式被暂停
	p.value.PaymentSign = paymentSign
	return nil
}

// 绑定订单号,如果交易号为空则绑定参数中传递的交易号
func (p *paymentOrderImpl) BindOrder(orderId int64, tradeNo string) error {
	//todo: check order exists  and tradeNo exists
	p.value.OrderId = orderId
	if len(p.value.TradeNo) == 0 {
		p.value.TradeNo = tradeNo
	}
	return nil
}

// 提交支付订单
func (p *paymentOrderImpl) Commit() (int64, error) {
	if id := p.GetAggregateRootId(); id > 0 {
		return id, payment.ErrOrderCommitted
	}
	if p.GetTradeNo() == "" {

	}
	return p.save()
}

func (p *paymentOrderImpl) save() (int64, error) {
	_, err := p.checkPaymentOk()
	if err == nil {
		unix := time.Now().Unix()
		if p.value.CreateTime == 0 {
			p.value.CreateTime = unix
		}
		p.value.Id, err = p.rep.SavePaymentOrder(p.value)
	}

	//保存支付单后,通知支付成功。只通知一次
	if err == nil && p.firstFinishPayment {
		p.firstFinishPayment = false
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
	if p.value.State == payment.StateFinishPayment {
		return payment.ErrOrderPayed
	}
	if p.value.State == payment.StateHasCancel {
		return payment.ErrOrderHasCancel
	}
	p.value.State = payment.StateFinishPayment
	p.value.OuterNo = outerNo
	p.value.PaidTime = time.Now().Unix()
	p.firstFinishPayment = true
	_, err := p.save()
	return err
}
func (p *paymentOrderImpl) GetValue() payment.PaymentOrder {
	return *p.value
}

// 取消支付
func (p *paymentOrderImpl) Cancel() error {
	p.value.State = payment.StateHasCancel
	return nil
}

// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
func (p *paymentOrderImpl) Adjust(amount float32) error {
	p.value.AdjustmentAmount += amount
	p.value.FinalAmount += amount
	if p.value.FinalAmount <= 0 {
		_, err := p.checkPaymentOk()
		return err
	}
	_, err := p.save()
	return err
}

type PaymentRepBase struct {
}

func (p *PaymentRepBase) CreatePaymentOrder(v *payment.
	PaymentOrder, rep payment.IPaymentRep, mmRep member.IMemberRep,
	orderManager order.IOrderManager, valRep valueobject.IValueRep) payment.IPaymentOrder {
	return &paymentOrderImpl{
		rep:          rep,
		value:        v,
		mmRep:        mmRep,
		valRep:       valRep,
		orderManager: orderManager,
	}
}
