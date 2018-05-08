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
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
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
	rep                payment.IPaymentRepo
	value              *payment.PaymentOrder
	mmRepo             member.IMemberRepo
	valRepo            valueobject.IValueRepo
	coupons            []promotion.ICouponPromotion
	orderManager       order.IOrderManager
	firstFinishPayment bool //第一次完成支付
	paymentUser        member.IMember
	buyer              member.IMember
}

func (p *paymentOrderImpl) GetAggregateRootId() int {
	return p.value.ID
}

func (p *paymentOrderImpl) Get() payment.PaymentOrder {
	return *p.value
}

// 获取交易号
func (p *paymentOrderImpl) TradeNo() string {
	return p.value.TradeNo
}

// 支付单状态
func (p *paymentOrderImpl) State() int {
	return int(p.value.State)
}

// 提交支付订单
func (p *paymentOrderImpl) Submit() error {
	if id := p.GetAggregateRootId(); id > 0 {
		return payment.ErrOrderCommitted
	}
	// 检查支付单单号是否匹配
	if b := p.rep.CheckTradeNoMatch(p.value.TradeNo, p.GetAggregateRootId()); !b {
		return payment.ErrExistsTradeNo
	}
	_, err := p.save()
	return err
}

// 取消支付
func (p *paymentOrderImpl) Cancel() error {
	oriState := p.value.State //支付单原始状态
	p.value.State = payment.StateCancelled
	_, err := p.save()
	if err == nil {
		// log.Println(fmt.Sprintf("-- 支付单详情：%#v",p.value))
		mm := p.getBuyer()
		if mm == nil {
			return member.ErrNoSuchMember
		}

		/* //todo: ！！！
		pv := p.Get()
		acc := mm.GetAccount()
		//退回到余额
		if pv.BalanceDiscount > 0 {
			err = acc.Refund(member.AccountBalance,
				member.KindBalanceRefund, "订单退款", pv.TradeNo,
				pv.BalanceDiscount, member.DefaultRelateUser)
		}
		//退积分
		if pv.IntegralDiscount > 0 {
			//todo : 退换积分,暂时积分抵扣的不退款
		}
		*/

		// 如果已经支付，则将支付的款项退回到账户
		if p.value.FinalFee > 0 && oriState == payment.StateFinished {
			/* //todo:!!!
			//退到钱包账户
			if pv.PaymentSign == payment.SignWalletAccount {
				return acc.Refund(member.AccountWallet,
					member.KindWalletPaymentRefund,
					"订单退款", pv.TradeNo, pv.FinalFee,
					member.DefaultRelateUser)
			}
			*/
		}
	}
	return err
}

// 线下现金/刷卡支付,cash:现金,bank:刷卡金额,finalZero:是否金额必须为零
func (p *paymentOrderImpl) OfflineDiscount(cash int, bank int, finalZero bool) error {
	panic("implement me")
}

// 交易完成
func (p *paymentOrderImpl) TradeFinish() error {
	if p.value.State == payment.StateFinished {
		return payment.ErrOrderPayed
	}
	if p.value.State == payment.StateCancelled {
		return payment.ErrOrderHasCancel
	}
	p.value.State = payment.StateFinished
	p.value.OutTradeNo = ""
	p.value.PaidTime = time.Now().Unix()
	p.firstFinishPayment = true
	_, err := p.save()
	return err
}

// 支付完成,传入第三名支付名称,以及外部的交易号
func (p *paymentOrderImpl) PaymentFinish(spName string, outerNo string) error {
	outerNo = strings.TrimSpace(outerNo)
	if len(outerNo) < 8 {
		return payment.ErrOuterNo
	}
	if p.value.State == payment.StateFinished {
		return payment.ErrOrderPayed
	}
	if p.value.State == payment.StateCancelled {
		return payment.ErrOrderHasCancel
	}
	p.value.State = payment.StateFinished
	p.value.OutTradeNo = outerNo
	p.value.PaidTime = time.Now().Unix()
	p.firstFinishPayment = true
	_, err := p.save()
	return err
}

// 重新修正金额
func (p *paymentOrderImpl) fixFee() {
	/* todo: !!!
	v := p.value
	v.FinalFee = v.TotalAmount - v.CouponDiscount - v.BalanceDiscount -
		v.IntegralDiscount - v.SubAmount - v.SystemDiscount
	*/
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

	/** todo:!!!
	// 通知订单支付完成
	if p.value.OrderId > 0 {
		err := p.orderManager.NotifyOrderTradeSuccess(int64(p.value.OrderId))
		domain.HandleError(err, "domain")
	}
	*/
}

// 优惠券抵扣

func (p *paymentOrderImpl) CouponDiscount(coupon promotion.ICouponPromotion) (
	float32, error) {
	return 0, nil
	/** todo:!!!
if p.value.PaymentSign&payment.OptUseCoupon == 0 {
	return 0, payment.ErrCanNotUseCoupon
}
//todo: 如可以使用多张优惠券,那么初始化应该获取支付单的所有优惠券
if p.coupons == nil {
	p.coupons = []promotion.ICouponPromotion{}
}
p.coupons = append(p.coupons, coupon)
// 支付金额应减去立减和系统支付的部分
fee := p.value.TotalAmount - p.value.SubAmount -
	p.value.SystemDiscount
for _, v := range p.coupons {
	p.value.CouponDiscount += v.GetCouponFee(fee)
}
p.fixFee()
return p.value.CouponDiscount, nil
	*/
}

// 在支付之前检查订单状态
func (p *paymentOrderImpl) checkPayment() error {
	if p.GetAggregateRootId() <= 0 {
		return payment.ErrPaymentNotSave
	}
	switch p.value.State {
	case payment.StateAwaitingPayment:
		if p.value.FinalFee == 0 {
			return payment.ErrFinalFee
		}
	case payment.StateFinished:
		return payment.ErrOrderPayed
	case payment.StateCancelled:
		return payment.ErrOrderHasCancel
	}
	return nil
}

// 应用余额支付
func (p *paymentOrderImpl) getBalanceDiscountAmount(acc member.IAccount) int {
	if p.value.FinalFee <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if int(acv.Balance) >= p.value.FinalFee {
		return p.value.FinalFee
	} else {
		return int(float64(acv.Balance) * enum.RATE_AMOUNT)
	}
	return 0
}

func (p *paymentOrderImpl) getPaymentUser() member.IMember {
	if p.paymentUser == nil && p.value.PayUid > 0 {
		p.paymentUser = p.mmRepo.GetMember(p.value.PayUid)
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
		/** todo: !!!
		// 从会员账户扣减,并更新支付单
		err = acc.PaymentDiscount(p.value.TradeNo, amount, remark)
		if err == nil {
			p.value.BalanceDiscount = amount
			p.fixFee()
			_, err = p.save()
		}
		*/
	}
	return err
}

// 检查是否支付完成, 且返回是否为第一次支付成功,
func (p *paymentOrderImpl) checkPaymentOk() (bool, error) {
	b := false
	if p.value.State == payment.StateAwaitingPayment {
		unix := time.Now().Unix()
		// 如果支付完成,则更新订单状态
		if b = p.value.FinalFee == 0; b {
			p.value.State = payment.StateFinished
			p.firstFinishPayment = true
		}
		p.value.PaidTime = unix
	}
	return b, nil
}

// 计算积分折算后的金额
func (p *paymentOrderImpl) mathIntegralFee(integral int64) int {
	if integral > 0 {
		conf := p.valRepo.GetGlobNumberConf()
		if conf.IntegralDiscountRate > 0 {
			return int(float32(float64(integral)*enum.RATE_AMOUNT) / conf.IntegralDiscountRate)
		}
	}
	return 0
}

// 积分抵扣,返回抵扣的金额及错误,ignoreAmount:是否忽略超出订单金额的积分
func (p *paymentOrderImpl) IntegralDiscount(integral int64, ignoreAmount bool) (float32, error) {
	if p.value.PaymentSign&payment.OptIntegralDiscount != payment.OptIntegralDiscount {
		return 0, payment.ErrCanNotUseIntegral
	}
	err := p.checkPayment()
	if err != nil {
		return 0, err
	}
	return 0, err
	/* todo:!!! 积分兑换
	// 判断扣减金额是否大于0
	amount := p.mathIntegralFee(integral)
	// 如果不忽略超出订单支付金额的积分,那么按实际来抵扣
	if !ignoreAmount && amount > p.value.FinalFee {
		conf := p.valRepo.GetGlobNumberConf()
		amount = p.value.FinalFee
		integral = int64(amount * conf.IntegralDiscountRate)
	}

	if amount > 0 {
		acc := p.mmRepo.GetMember(p.value.BuyUser).GetAccount()
		// 抵扣积分

		//log.Println("----", p._value.BuyUser, acc.GetValue().Integral, "discount:", integral)
		//log.Printf("-----%#v\n", acc.GetValue())
		err = acc.IntegralDiscount(member.TypeIntegralPaymentDiscount,
			p.Get().TradeNo, integral, "")
		if err == nil {
			p.value.IntegralDiscount = amount
			p.fixFee()
			_, err = p.save()
		}
	}
	return amount, err

	*/
}

// 系统支付金额
func (p *paymentOrderImpl) SystemPayment(fee float32) error {
	return nil
	/* todo:!!!
	if p.value.PaymentSign&payment.OptSystemPayment == 0 {
		return payment.ErrCanNotSystemDiscount
	}
	err := p.checkPayment()
	if err == nil {
		p.value.SystemDiscount += fee
		p.fixFee()
	}
	return err
	*/
}

func (p *paymentOrderImpl) getBuyer() member.IMember {
	if p.buyer == nil {
		p.buyer = p.mmRepo.GetMember(p.value.BuyerId)
	}
	return p.buyer
}

func (p *paymentOrderImpl) intAmount(a float32) int {
	return int(a * float32(enum.RATE_AMOUNT))
}

// 余额钱包混合支付，优先扣除余额。
func (p *paymentOrderImpl) HybridPayment(remark string) error {
	buyer := p.getBuyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	v := p.Get()
	acc := buyer.GetAccount().GetValue()
	optFlag := int(v.PayFlag)
	// 判断是否能余额支付
	if p := payment.OptBalanceDiscount; optFlag&p != p {
		return payment.ErrNotSupportPaymentOpt
	}
	// 如果余额够支付，则优先余额支付
	if p.intAmount(acc.Balance) >= v.FinalFee {
		return p.BalanceDiscount(remark)
	}
	// 判断是否能钱包支付
	if p := payment.OptWalletPayment; optFlag&p != p {
		return payment.ErrNotSupportPaymentOpt
	}
	// 判断是否余额不足
	if p.intAmount(acc.Balance+acc.WalletBalance) < v.FinalFee {
		return payment.ErrNotEnoughAmount
	}
	err := p.BalanceDiscount(remark)
	if err == nil {
		err = p.PaymentByWallet(remark)
	}
	return err
}

// 使用会员的余额抵扣
func (p *paymentOrderImpl) BalanceDiscount(remark string) error {
	return p.paymentWithBalance(payment.PaymentByBuyer, remark)
}

// 钱包账户支付
func (p *paymentOrderImpl) PaymentByWallet(remark string) error {
	amount := p.value.FinalFee
	buyer := p.getBuyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	acc := buyer.GetAccount()
	av := acc.GetValue()
	if p.intAmount(av.WalletBalance) < amount {
		return payment.ErrNotEnoughAmount
	}
	if remark == "" {
		remark = "支付订单"
	}
	err := acc.DiscountWallet(remark, p.TradeNo(), float32(float64(amount)/enum.RATE_AMOUNT),
		member.DefaultRelateUser, true)
	if err == nil {
		/** todo:!!!
		p.value.PaymentSign = payment.SignWalletAccount
		*/
		// 标记为第一次支付
		p.firstFinishPayment = true
		p.value.State = payment.StateFinished
		p.value.PaidTime = time.Now().Unix()
		_, err = p.save()
	}
	return err
}

// 设置支付方式
func (p *paymentOrderImpl) SetPaymentSign(paymentSign int32) error {
	/** todo:!!!
	//todo: 某个支付方式被暂停
	p.value.PaymentSign = paymentSign
	*/
	return nil
}

// 绑定订单号,如果交易号为空则绑定参数中传递的交易号
func (p *paymentOrderImpl) BindOrder(orderId int64, tradeNo string) error {
	/* todo:!!!
	//todo: check order exists  and tradeNo exists
	p.value.OrderId = int32(orderId)
	if len(p.value.TradeNo) == 0 {
		p.value.TradeNo = tradeNo
	}
	*/
	return nil
}

func (p *paymentOrderImpl) save() (int, error) {
	_, err := p.checkPaymentOk()
	if err == nil {
		unix := time.Now().Unix()
		if p.value.SubmitTime == 0 {
			p.value.SubmitTime = unix
		}
		p.value.ID, err = p.rep.SavePaymentOrder(p.value)
	}

	//保存支付单后,通知支付成功。只通知一次
	if err == nil && p.firstFinishPayment {
		p.firstFinishPayment = false
		go p.notifyPaymentFinish()
	}
	return p.GetAggregateRootId(), err
}

// 退款
func (p *paymentOrderImpl) Refund(amount float64) (err error) {
	mm := p.getBuyer()
	if mm == nil {
		return member.ErrNoSuchMember
	}
	return nil
	/* todo:!!!
	pv := p.Get()
	originState := pv.State
	acc := mm.GetAccount()
	//先通过退回到余额
	if pv.BalanceDiscount > 0 {
		final := math.Min(float64(pv.BalanceDiscount), amount)
		if amount > float64(pv.BalanceDiscount) {
			amount = amount - final
		}
		err = acc.Refund(member.AccountBalance,
			member.KindBalanceRefund, "订单退款", pv.TradeNo,
			float32(final), member.DefaultRelateUser)
		if err == nil {
			p.value.BalanceDiscount -= float32(final)
		}
	}
	if amount > 0 && p.value.FinalFee > 0 &&
		originState == payment.StateFinished {
		//退到钱包账户
		if pv.PaymentSign == payment.SignWalletAccount {
			err = acc.Refund(member.AccountWallet,
				member.KindWalletPaymentRefund,
				"订单退款", pv.TradeNo, float32(amount),
				member.DefaultRelateUser)
			if err == nil {
				p.value.FinalFee -= float32(amount)
			}
		}
	}
	if err == nil {
		_, err = p.save()
	}
	return err
	*/
}

// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
func (p *paymentOrderImpl) Adjust(amount int) error {
	p.value.AdjustAmount += amount
	p.value.FinalFee += amount
	if p.value.FinalFee <= 0 {
		_, err := p.checkPaymentOk()
		return err
	}
	_, err := p.save()
	return err
}

type PaymentRepBase struct {
}

func (p *PaymentRepBase) CreatePaymentOrder(v *payment.
PaymentOrder, rep payment.IPaymentRepo, mmRepo member.IMemberRepo,
	orderManager order.IOrderManager, valRepo valueobject.IValueRepo) payment.IPaymentOrder {
	return &paymentOrderImpl{
		rep:          rep,
		value:        v,
		mmRepo:       mmRepo,
		valRepo:      valRepo,
		orderManager: orderManager,
	}
}
