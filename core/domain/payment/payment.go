package payment

/**
 * Copyright 2015 @ z3q.net.
 * name : payment
 * author : jarryliu
 * date : 2016-07-03 09:25
 * description :
 * history :
 */

import (
	"errors"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"math"
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
	repo               payment.IPaymentRepo
	value              *payment.Order
	memberRepo         member.IMemberRepo
	valueRepo          valueobject.IValueRepo
	coupons            []promotion.ICouponPromotion
	orderManager       order.IOrderManager
	firstFinishPayment bool //第一次完成支付
	paymentUser        member.IMember
	buyer              member.IMember
}

func (p *paymentOrderImpl) GetAggregateRootId() int {
	return p.value.ID
}

func (p *paymentOrderImpl) Get() payment.Order {
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

// 支付标志
func (p *paymentOrderImpl) Flag() int {
	return p.value.PayFlag
}

// 支付途径支付信息
func (p *paymentOrderImpl) Channels() []*payment.TradeChan {
	if p.value.TradeChannels == nil {
		if p.GetAggregateRootId() <= 0 {
			return make([]*payment.TradeChan, 0)
		}
		p.value.TradeChannels = p.repo.GetTradeChannelItems(p.TradeNo())
	}
	return p.value.TradeChannels
}

// 提交支付订单
func (p *paymentOrderImpl) Submit() error {
	if id := p.GetAggregateRootId(); id > 0 {
		return payment.ErrOrderCommitted
	}
	p.prepareSubmit() // 提交之前进行操作
	// 检查支付单单号是否匹配
	if b := p.repo.CheckTradeNoMatch(p.value.TradeNo, p.GetAggregateRootId()); !b {
		return payment.ErrExistsTradeNo
	}
	return p.saveOrder()
}

// 合并支付
func (p *paymentOrderImpl) MergePay(orders []payment.IPaymentOrder) (mergeTradeNo string, finalFee int, err error) {
	if err := p.CheckPaymentState(); err != nil { // 验证支付单是否可以支付
		return "", 0, err
	}
	if len(orders) == 0 {
		panic(errors.New("will be merge trade orders is nil"))
	}
	finalFee = p.value.FinalFee
	tradeOrders := []string{p.TradeNo()}
	for _, v := range orders {
		// 检查支付单状态
		if err := v.CheckPaymentState(); err != nil {
			return "", 0, err
		}
		// 统计支付总金额
		finalFee += v.Get().FinalFee
		tradeOrders = append(tradeOrders, v.TradeNo())
	}
	// 清除欲合并的支付单
	err = p.repo.ResetMergePaymentOrders(tradeOrders)
	// 合并支付
	if err == nil {
		mergeTradeNo = "MG" + p.TradeNo()
		err = p.repo.SaveMergePaymentOrders(mergeTradeNo, tradeOrders)
	}
	return mergeTradeNo, finalFee, err
}

// 准备提交支付单
func (p *paymentOrderImpl) prepareSubmit() {
	unix := time.Now().Unix()
	p.value.SubmitTime = unix
	p.value.UpdateTime = unix
	// 初始化状态
	p.value.State = payment.StateAwaitingPayment
	// 初始化支付用户编号
	if p.value.PayUid <= 0 {
		p.value.PayUid = p.value.BuyerId
	}
}

// 在支付之前检查订单状态
func (p *paymentOrderImpl) CheckPaymentState() error {
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
	case payment.StateCancelled, payment.StateAborted:
		return payment.ErrOrderCancelled
	}
	return nil
}

// 检查是否支付完成, 且返回是否为第一次支付成功,
func (p *paymentOrderImpl) checkOrderFinalFee() error {
	if p.value.State == payment.StateAwaitingPayment {
		if p.value.ItemAmount <= 0 { // 检查支付金额
			return payment.ErrItemAmount
		}
		// 修正支付单共计金额
		p.value.TotalAmount = p.value.ItemAmount - p.value.DiscountAmount + p.value.AdjustAmount
		// 修正支付单金额
		p.value.FinalFee = p.value.ItemAmount - p.value.DeductAmount + p.value.ProcedureFee
		unix := time.Now().Unix()
		// 如果支付完成,则更新订单状态
		if p.value.FinalFee == 0 {
			p.value.State = payment.StateFinished
			p.firstFinishPayment = true
		}
		p.value.PaidTime = unix
	}
	return nil
}

// 取消支付,并退款
func (p *paymentOrderImpl) Cancel() (err error) {
	if p.value.State == payment.StateCancelled {
		return payment.ErrOrderCancelled
	}
	p.value.State = payment.StateCancelled
	if err = p.saveOrder(); err != nil {
		return err
	}
	mm := p.getBuyer()
	if mm == nil {
		return member.ErrNoSuchMember
	}
	pv := p.Get()
	acc := mm.GetAccount()
	chanMap := p.getPaymentChannelMap()
	//退回到余额
	if v := chanMap[payment.MBalance]; v > 0 {
		err = acc.Refund(member.AccountBalance,
			member.KindBalanceRefund, "订单退款", pv.TradeNo,
			float32(v/100), member.DefaultRelateUser)
	}
	//退积分
	if v := chanMap[payment.MIntegral]; v > 0 {
		//todo : 退换积分,暂时积分抵扣的不退款
	}
	// 如果已经支付，则将支付的款项退回到账户
	if v := chanMap[payment.MWallet]; v > 0 {
		return acc.Refund(member.AccountWallet,
			member.KindWalletPaymentRefund,
			"订单退款", pv.TradeNo, float32(v/100),
			member.DefaultRelateUser)
	}
	return err
}

// 线下现金/刷卡支付,cash:现金,bank:刷卡金额,finalZero:是否金额必须为零
func (p *paymentOrderImpl) OfflineDiscount(cash int, bank int, finalZero bool) error {
	if err := p.CheckPaymentState(); err != nil {
		return err
	}
	if !p.andMethod(p.value.PayFlag, payment.MCash) {
		return payment.ErrNotSupportPaymentChannel
	}
	if !p.andMethod(p.value.PayFlag, payment.MBankCard) {
		return payment.ErrNotSupportPaymentChannel
	}
	if cash+bank > p.value.FinalFee {
		return payment.ErrOutOfFinalFee
	}
	if finalZero && p.value.FinalFee > cash+bank {
		return payment.ErrNotMatchFinalFee
	}
	p.value.DeductAmount += cash + bank
	err := p.saveOrder()
	if err == nil {
		if cash > 0 {
			err = p.saveTradeChan(cash, payment.MCash, "")
		}
		if bank > 0 {
			err = p.saveTradeChan(bank, payment.MBankCard, "")
		}
	}
	return err
}

// 交易完成
func (p *paymentOrderImpl) TradeFinish() error {
	if err := p.CheckPaymentState(); err != nil {
		return err
	}
	p.value.State = payment.StateFinished
	p.value.OutTradeNo = ""
	p.value.PaidTime = time.Now().Unix()
	p.firstFinishPayment = true
	return p.saveOrder()
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
		return payment.ErrOrderCancelled
	}
	p.value.State = payment.StateFinished
	p.value.OutTradeNo = outerNo
	p.value.PaidTime = time.Now().Unix()
	p.firstFinishPayment = true
	return p.saveOrder()
}

// 更新订单状态, 需要注意,防止多次订单更新
func (p *paymentOrderImpl) notifyPaymentFinish() {
	if p.GetAggregateRootId() <= 0 {
		panic(payment.ErrNoSuchPaymentOrder)
	}
	// 通知订单支付完成
	if p.value.OutOrderNo != "" {
		subOrder := p.value.SubOrder == 1
		err := p.orderManager.NotifyOrderTradeSuccess(p.value.OutOrderNo, subOrder)
		domain.HandleError(err, "domain")
	}
}

// 优惠券抵扣

func (p *paymentOrderImpl) CouponDiscount(coupon promotion.ICouponPromotion) (
	int, error) {
	//** todo:!!! 应该在订单除使用优惠券
	//todo: 如可以使用多张优惠券,那么初始化应该获取支付单的所有优惠券
	if p.coupons == nil {
		p.coupons = []promotion.ICouponPromotion{}
	}
	p.coupons = append(p.coupons, coupon)
	// 支付金额应减去立减和系统支付的部分
	fee := p.value.FinalFee
	for _, v := range p.coupons {
		p.value.DiscountAmount += int(v.GetCouponFee(float32(fee/100)) * 100)
	}
	return p.value.DiscountAmount, nil
}

// 应用余额支付
func (p *paymentOrderImpl) getBalanceDiscountAmount(acc member.IAccount) int {
	if p.value.FinalFee <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if int(acv.Balance*100) >= p.value.FinalFee {
		return p.value.FinalFee
	} else {
		return int(float64(acv.Balance * 100))
	}
	return 0
}

func (p *paymentOrderImpl) getPaymentUser() member.IMember {
	if p.paymentUser == nil && p.value.PayUid > 0 {
		p.paymentUser = p.memberRepo.GetMember(p.value.PayUid)
	}
	return p.paymentUser
}

// 验证是否支持支付方式
func (p *paymentOrderImpl) andMethod(flag, method int) bool {
	v := 1 << uint(method-1)
	return flag&v == v
}

// 使用余额抵扣
func (p *paymentOrderImpl) BalanceDiscount(remark string) error {
	if b := p.andMethod(p.value.PayFlag, payment.MBalance); !b { // 检查支付方式
		return payment.ErrNotSupportPaymentChannel
	}
	if err := p.CheckPaymentState(); err != nil { // 检查支付单状态
		return err
	}
	pu := p.getPaymentUser()
	if pu == nil {
		return member.ErrNoSuchMember
	}
	acc := pu.GetAccount()
	amount := p.getBalanceDiscountAmount(acc)
	if amount == 0 {
		return member.ErrAccountBalanceNotEnough
	}
	err := acc.PaymentDiscount(p.value.TradeNo, float32(amount/100), remark)
	if err == nil {
		p.value.DeductAmount += amount // 修改抵扣金额
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(amount, payment.MBalance, "")
		}
	}
	return err
}

// 计算积分折算后的金额
func (p *paymentOrderImpl) getIntegralExchangeAmount(integral int) int {
	if integral > 0 {
		conf := p.valueRepo.GetGlobNumberConf()
		if conf.IntegralDiscountRate > 0 {
			return int(float32(integral)/conf.IntegralDiscountRate) * 100
		}
	}
	return 0
}

// 积分抵扣,返回抵扣的金额及错误,ignoreAmount:是否忽略超出订单金额的积分
func (p *paymentOrderImpl) IntegralDiscount(integral int,
	ignoreAmount bool) (amount int, err error) {
	if !p.andMethod(p.value.PayFlag, payment.MIntegral) {
		return 0, payment.ErrNotSupportPaymentChannel
	}
	if err = p.CheckPaymentState(); err != nil {
		return 0, err
	}
	// 判断扣减金额是否大于0
	amount = p.getIntegralExchangeAmount(integral)
	// 如果不忽略超出订单支付金额的积分,那么按实际来抵扣
	if !ignoreAmount && amount > p.value.FinalFee {
		amount = p.value.FinalFee
		conf := p.valueRepo.GetGlobNumberConf()
		integral = int(float32(amount) * conf.IntegralDiscountRate)
	}
	if amount <= 0 {
		return 0, nil
	}
	acc := p.memberRepo.GetMember(p.value.BuyerId).GetAccount()
	//log.Println("----", p.value.BuyerId, acc.GetValue().Integral, "discount:", integral)
	//log.Printf("-----%#v\n", acc.GetValue())
	err = acc.IntegralDiscount(member.TypeIntegralPaymentDiscount,
		p.Get().TradeNo, int64(integral), "")
	// 抵扣积分
	if err == nil {
		p.value.DeductAmount += amount
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(amount, payment.MIntegral, "")
		}
	}
	return amount, err
}

// 系统支付金额
func (p *paymentOrderImpl) SystemPayment(fee int) error {
	if !p.andMethod(p.value.PayFlag, payment.MSystemPay) {
		return payment.ErrNotSupportPaymentChannel
	}
	err := p.CheckPaymentState()
	if err == nil {
		p.value.DeductAmount += fee
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(fee, payment.MSystemPay, "")
		}
	}
	return err
}

// 保存支付信息
func (p *paymentOrderImpl) saveTradeChan(amount int, payChan int, chanData string) error {
	c := &payment.TradeChan{
		TradeNo:      p.TradeNo(),
		PayChan:      payChan,
		InternalChan: 1,
		PayAmount:    amount,
		ChanData:     chanData,
	}
	_, err := p.repo.SavePaymentTradeChan(p.TradeNo(), c)
	return err
}

func (p *paymentOrderImpl) getBuyer() member.IMember {
	if p.buyer == nil {
		p.buyer = p.memberRepo.GetMember(p.value.BuyerId)
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
	// 判断是否能余额支付
	if !p.andMethod(p.value.PayFlag, payment.MBalance) {
		return payment.ErrNotSupportPaymentChannel
	}
	// 如果余额够支付，则优先余额支付
	if p.intAmount(acc.Balance) >= v.FinalFee {
		return p.BalanceDiscount(remark)
	}
	// 判断是否能钱包支付
	if !p.andMethod(p.value.PayFlag, payment.MWallet) {
		return payment.ErrNotSupportPaymentChannel
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

// 钱包账户支付
func (p *paymentOrderImpl) PaymentByWallet(remark string) error {
	if !p.andMethod(p.value.PayFlag, payment.MWallet) {
		return payment.ErrNotSupportPaymentChannel
	}
	buyer := p.getBuyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	amount := p.value.FinalFee
	// 判断并从钱包里扣款
	acc := buyer.GetAccount()
	if p.intAmount(acc.GetValue().WalletBalance) < amount {
		return payment.ErrNotEnoughAmount
	}
	err := acc.DiscountWallet(remark, p.TradeNo(), float32(float32(amount)/100),
		member.DefaultRelateUser, true)
	if err == nil {
		p.value.DeductAmount += amount
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(amount, payment.MWallet, "")
		}
	}
	return err
}

// 使用会员卡支付,cardCode:会员卡编码,amount:支付金额
func (p *paymentOrderImpl) PaymentWithCard(cardCode string, amount int) error {
	panic("not implement")
}

// 设置支付方式
func (p *paymentOrderImpl) SetTradeSP(spName string) error {
	err := p.CheckPaymentState()
	if err == nil {
		p.value.OutTradeSp = spName
	}
	return nil
}

func (p *paymentOrderImpl) saveOrder() error {
	// 检查支付单
	err := p.checkOrderFinalFee()
	if err == nil {
		p.value.UpdateTime = time.Now().Unix()
		p.value.ID, err = p.repo.SavePaymentOrder(p.value)
	}
	//保存支付单后,通知支付成功。只通知一次
	if err == nil && p.firstFinishPayment {
		p.firstFinishPayment = false
		go p.notifyPaymentFinish()
	}
	return err
}

// 退款
func (p *paymentOrderImpl) Refund(amount int) (err error) {
	mm := p.getBuyer()
	if mm == nil {
		return member.ErrNoSuchMember
	}
	pv := p.Get()
	acc := mm.GetAccount()
	chanMap := p.getPaymentChannelMap()
	// 先退回到余额
	if v := chanMap[payment.MBalance]; v > 0 {
		final := int(math.Min(float64(v), float64(amount)))
		if amount > v {
			amount = amount - final
		}
		err = acc.Refund(member.AccountBalance,
			member.KindBalanceRefund, "订单退款", pv.TradeNo,
			float32(final)/100, member.DefaultRelateUser)
		if err == nil {
			p.value.DeductAmount -= final
		}
	}
	// 如果已经支付，则将支付的款项退回到账户
	if v := chanMap[payment.MWallet]; v > 0 {
		final := int(math.Min(float64(v), float64(amount)))
		if amount > v {
			amount = amount - final
		}
		err = acc.Refund(member.AccountWallet,
			member.KindWalletPaymentRefund,
			"订单退款", pv.TradeNo, float32(final)/100,
			member.DefaultRelateUser)
		if err == nil {
			p.value.DeductAmount -= amount
		}
	}
	//todo: 原路退回，目前全部退回钱包
	if amount > 0 {
		err = acc.Refund(member.AccountWallet,
			member.KindWalletPaymentRefund,
			"订单退款", pv.TradeNo, float32(amount)/100,
			member.DefaultRelateUser)
		if err == nil {
			p.value.DeductAmount -= amount
		}
	}
	return p.saveOrder()
}

// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
func (p *paymentOrderImpl) Adjust(amount int) error {
	p.value.AdjustAmount += amount
	p.value.FinalFee += amount
	if p.value.FinalFee <= 0 {
		return p.checkOrderFinalFee()
	}
	return p.saveOrder()
}

// 获取支付途径支付信息字典
func (p *paymentOrderImpl) getPaymentChannelMap() map[int]int {
	mp := make(map[int]int, 0)
	arr := p.Channels()
	for _, v := range arr {
		if v.PayAmount > 0 {
			c, ok := mp[v.PayChan]
			if ok {
				mp[v.PayChan] = c + v.PayAmount
			} else {
				mp[v.PayChan] = v.PayAmount
			}
		}
	}
	return mp
}

type RepoBase struct {
}

func (p *RepoBase) CreatePaymentOrder(v *payment.
	Order, repo payment.IPaymentRepo, mmRepo member.IMemberRepo,
	orderManager order.IOrderManager, valRepo valueobject.IValueRepo) payment.IPaymentOrder {
	return &paymentOrderImpl{
		repo:         repo,
		value:        v,
		memberRepo:   mmRepo,
		valueRepo:    valRepo,
		orderManager: orderManager,
	}
}
