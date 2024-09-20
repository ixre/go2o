/**
 * Copyright 2015 @ 56x.net.
 * name : payment
 * author : jarryliu
 * date : 2016-07-03 09:25
 * description :
 * history :
 */
package payment

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/gof/domain/eventbus"
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
	registryRepo       registry.IRegistryRepo
	coupons            []promotion.ICouponPromotion
	orderManager       order.IOrderManager
	firstFinishPayment bool //第一次完成支付
	paymentUser        member.IMemberAggregateRoot
	buyer              member.IMemberAggregateRoot
}

func (p *paymentOrderImpl) GetAggregateRootId() int {
	return p.value.Id
}

func (p *paymentOrderImpl) Get() payment.Order {
	return *p.value
}

// TradeNo 获取交易号
func (p *paymentOrderImpl) TradeNo() string {
	return p.value.TradeNo
}

// State 支付单状态
func (p *paymentOrderImpl) State() int {
	return int(p.value.Status)
}

// Flag 支付标志
func (p *paymentOrderImpl) Flag() int {
	return p.value.PayFlag
}

// TradeMethods 支付途径支付信息
func (p *paymentOrderImpl) TradeMethods() []*payment.PayTradeData {
	if p.value.TradeMethods == nil {
		if p.GetAggregateRootId() <= 0 {
			return make([]*payment.PayTradeData, 0)
		}
		p.value.TradeMethods = p.repo.GetTradeChannelItems(p.TradeNo())
	}
	return p.value.TradeMethods
}

// Submit 提交支付订单
func (p *paymentOrderImpl) Submit() error {
	if id := p.GetAggregateRootId(); id > 0 {
		return payment.ErrOrderCommitted
	}
	err := p.prepareSubmit() // 提交之前进行操作
	if err != nil {
		return err
	}

	if len(p.value.TradeNo) == 0 {
		p.value.TradeNo = p.generateTradeNo()
	} else {
		// 检查外部指定的支付单单号是否重复
		if b := p.repo.CheckTradeNoMatch(p.value.TradeNo, p.GetAggregateRootId()); !b {
			return payment.ErrExistsTradeNo
		}
	}
	err = p.saveOrder()
	if err == nil {
		// 保存支付单的支付方式,主要用于拆分子订单提交
		for _, v := range p.value.TradeMethods {
			v.OrderId = p.GetAggregateRootId()
			v.Id, _ = p.repo.SavePaymentTradeChan(p.TradeNo(), v)
		}
	}
	return err
}

// 生成交易号
func (p *paymentOrderImpl) generateTradeNo() string {
	var orderNo string
	i := 0
	for {
		orderNo = domain.NewTradeNo(p.value.OrderType, int(p.value.BuyerId))
		if p.repo.CheckTradeNoMatch(orderNo, p.GetAggregateRootId()) {
			break
		}
		if i++; i > 10 {
			log.Println("生成交易号失败")
			return ""
		}
	}
	return orderNo
}

// MergePay 合并支付
func (p *paymentOrderImpl) MergePay(orders []payment.IPaymentOrder) (mergeTradeNo string, finalAmount int, err error) {
	if err := p.CheckPaymentState(); err != nil { // 验证支付单是否可以支付
		return "", 0, err
	}
	if len(orders) == 0 {
		return "", 0, errors.New("will be merge trade orders is nil")
	}
	finalAmount = int(p.value.FinalAmount)
	tradeOrders := []string{p.TradeNo()}
	for _, v := range orders {
		// 检查支付单状态
		if err := v.CheckPaymentState(); err != nil {
			return "", 0, err
		}
		// 统计支付总金额
		finalAmount += int(v.Get().FinalAmount)
		tradeOrders = append(tradeOrders, v.TradeNo())
	}
	// 清除欲合并的支付单
	err = p.repo.ResetMergePaymentOrders(tradeOrders)
	// 合并支付
	if err == nil {
		mergeTradeNo = "MG" + p.TradeNo()
		err = p.repo.SaveMergePaymentOrders(mergeTradeNo, tradeOrders)
	}
	return mergeTradeNo, finalAmount, err
}

// 准备提交支付单
func (p *paymentOrderImpl) prepareSubmit() error {
	unix := time.Now().Unix()
	p.value.SubmitTime = int(unix)
	p.value.UpdateTime = int(unix)
	// 初始化状态
	p.value.Status = payment.StateAwaitingPayment
	// 初始化支付用户编号
	if p.value.PayerId <= 0 {
		p.value.PayerId = p.value.BuyerId
	}
	if p.value.ExpiresTime == 0 {
		return errors.New("支付单过期时间未设置")
	}
	return nil
}

// 在支付之前检查订单状态
func (p *paymentOrderImpl) CheckPaymentState() error {
	if p.GetAggregateRootId() <= 0 {
		return payment.ErrPaymentNotSave
	}
	switch p.value.Status {
	case payment.StateAwaitingPayment:
		if p.value.FinalAmount == 0 {
			return payment.ErrFinalAmount
		}
	case payment.StateFinished:
		return payment.ErrOrderPayed
	case payment.StateRefunded:
		return payment.ErrOrderRefunded
	case payment.StateClosed:
		return payment.ErrOrderClosed
	}
	return nil
}

// 检查是否支付完成, 且返回是否为第一次支付成功,
func (p *paymentOrderImpl) checkOrderFinalAmount() error {
	if p.value.Status == payment.StateAwaitingPayment {
		if p.value.TotalAmount <= 0 { // 检查支付金额
			return payment.ErrItemAmount
		}
		// 修正支付单共计金额
		//p.value.TotalAmount = p.value.ItemAmount - p.value.DiscountAmount + p.value.AdjustAmount
		// 修正支付单金额
		p.value.FinalAmount = p.value.TotalAmount - p.value.DeductAmount + p.value.TransactionFee
		// 如果支付完成,则更新订单状态
		if p.value.FinalAmount == 0 {
			p.value.Status = payment.StateFinished
			p.firstFinishPayment = true
			p.value.PaidTime = int(time.Now().Unix())
		}
	}
	return nil
}

// 取消支付,并退款
func (p *paymentOrderImpl) Cancel() (err error) {
	// 如果已取消或订单再次回调到支付单取消, 不做任何处理
	if p.value.Status == payment.StateClosed {
		return nil
	}
	p.value.Status = payment.StateClosed
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
			"订单退款", int(v), pv.TradeNo, "")
		if err != nil {
			log.Println("[ GO2O][ ERROR]: payment refund balance failed,",
				err.Error(), "orderNo", p.value.TradeNo)
		}
	}
	//退积分
	if v := chanMap[payment.MIntegral]; v > 0 {
		err = acc.Refund(member.AccountIntegral,
			"订单退款", int(v), pv.TradeNo, "")
		if err != nil {
			log.Println("[ GO2O][ ERROR]: payment refund integral failed,",
				err.Error(), "orderNo", p.value.TradeNo)
		}
	}
	// 如果已经支付，则将支付的款项退回到账户
	if v := chanMap[payment.MWallet]; v > 0 {
		err = acc.Refund(member.AccountWallet,
			"订单退款", int(v), pv.TradeNo, "")
		if err != nil {
			log.Println("[ GO2O][ ERROR]: payment refund wallet failed,",
				err.Error(), "orderNo", p.value.TradeNo)
		}
	}
	if err == nil {
		//todo: 临时注释, 应该支付已去掉了订单的依赖
		// err = p.orderManager.Cancel(p.value.OutOrderNo,
		// 	p.value.SubOrder == 1,
		// 	pv.SubOrder == 1,
		// 	"超时未付款")
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
	if int64(cash+bank) > int64(p.value.FinalAmount) {
		return payment.ErrOutOfFinalAmount
	}
	if finalZero && int64(cash+bank) > int64(p.value.FinalAmount) {
		return payment.ErrNotMatchFinalAmount
	}
	p.value.DeductAmount += int(cash + bank)
	var err error
	if cash > 0 {
		p.value.FinalFlag |= payment.MCash
		err = p.saveTradeChan(cash, payment.MCash, "", "")
	}
	if bank > 0 {
		p.value.FinalFlag |= payment.MBankCard
		err = p.saveTradeChan(bank, payment.MBankCard, "", "")
	}
	if err == nil {
		err = p.saveOrder()
	}
	return err
}

// 交易完成
func (p *paymentOrderImpl) TradeFinish() error {
	if err := p.CheckPaymentState(); err != nil {
		return err
	}
	p.value.Status = payment.StateFinished
	p.value.OutTradeNo = ""
	p.value.PaidTime = int(time.Now().Unix())
	p.firstFinishPayment = true
	return p.saveOrder()
}

// 支付完成,传入第三名支付名称,以及外部的交易号
func (p *paymentOrderImpl) PaymentFinish(spName string, outerNo string) error {
	outerNo = strings.TrimSpace(outerNo)
	if len(outerNo) < 8 {
		return payment.ErrOuterNo
	}
	p.value.Status = 1
	if p.value.Status == payment.StateFinished {
		return payment.ErrOrderPayed
	}
	if p.value.Status == payment.StateRefunded {
		return payment.ErrOrderRefunded
	}
	if p.value.Status == payment.StateClosed {
		return payment.ErrOrderClosed
	}
	p.value.Status = payment.StateFinished
	p.value.OutTradeSp = spName
	p.value.OutTradeNo = outerNo
	p.value.PaidTime = int(time.Now().Unix())
	p.firstFinishPayment = true
	err := p.saveOrder()
	// 保存第三方支付记录
	if err == nil {
		err = p.saveTradeChan(int(p.value.FinalAmount), payment.MPaySP, spName, outerNo)
	}
	return err
}

// 更新订单状态, 需要注意,防止多次订单更新
func (p *paymentOrderImpl) applyPaymentFinish() error {
	if p.GetAggregateRootId() > 0 {
		// 发布支付成功事件，并进行其他业务处理
		eventbus.Publish(&payment.PaymentSuccessEvent{
			Order:         p,
			TradeChannels: p.TradeMethods(),
		})

	}
	return nil
}

// 优惠券抵扣

func (p *paymentOrderImpl) CouponDiscount(s string) (
	int, error) {
	var coupon promotion.ICouponPromotion
	//** todo:!!! 应该在订单除使用优惠券
	//todo: 如可以使用多张优惠券,那么初始化应该获取支付单的所有优惠券
	if p.coupons == nil {
		p.coupons = []promotion.ICouponPromotion{}
	}
	p.coupons = append(p.coupons, coupon)
	// 支付金额应减去立减和系统支付的部分
	fee := p.value.FinalAmount
	for _, v := range p.coupons {
		p.value.DiscountAmount += int64(v.GetCouponFee(int(fee)) * 100)
	}
	p.value.FinalFlag |= payment.MUserCoupon
	return int(p.value.DiscountAmount), nil
}

// 应用余额支付
func (p *paymentOrderImpl) getBalanceDeductAmount(acc member.IAccount) int64 {
	if p.value.FinalAmount <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.Balance >= int(p.value.FinalAmount) {
		return int64(p.value.FinalAmount)
	}
	return int64(acv.Balance)
}

// 获取可用于钱包抵扣的金额
func (p *paymentOrderImpl) getWalletDeductAmount(acc member.IAccount) int64 {
	if p.value.FinalAmount <= 0 {
		return 0
	}
	acv := acc.GetValue()
	if acv.WalletBalance >= int(p.value.FinalAmount) {
		return int64(p.value.FinalAmount)
	}
	return int64(acv.WalletBalance)
}

func (p *paymentOrderImpl) getPaymentUser() member.IMemberAggregateRoot {
	if p.paymentUser == nil && p.value.PayerId > 0 {
		p.paymentUser = p.memberRepo.GetMember(int64(p.value.PayerId))
	}
	return p.paymentUser
}

// 验证是否支持支付方式
func (p *paymentOrderImpl) andMethod(flag, method int) bool {
	return flag&method == method
}

// 使用余额抵扣
func (p *paymentOrderImpl) BalanceDeduct(remark string) error {
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
	amount := p.getBalanceDeductAmount(acc)
	if amount == 0 {
		return member.ErrAccountBalanceNotEnough
	}
	err := acc.PaymentDiscount(p.value.TradeNo, int(amount), remark)
	if err == nil {
		p.value.DeductAmount += int(amount) // 修改抵扣金额
		p.value.FinalFlag |= payment.MBalance
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(int(amount), payment.MBalance, "", "")
		}
	}
	return err
}

// 使用余额抵扣
func (p *paymentOrderImpl) WalletDeduct(remark string) error {
	if b := p.andMethod(p.value.PayFlag, payment.MWallet); !b { // 检查支付方式
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
	amount := p.getWalletDeductAmount(acc)
	if amount == 0 {
		return member.ErrAccountNotEnoughAmount
	}
	_, err := acc.Discount(member.AccountWallet, "订单抵扣",
		int(amount), p.value.OutOrderNo, remark)
	if err == nil {
		p.value.DeductAmount += int(amount) // 修改抵扣金额
		p.value.FinalFlag |= payment.MWallet
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(int(amount), payment.MWallet, "", "")
		}
	}
	return err
}

// 计算积分折算后的金额
func (p *paymentOrderImpl) getIntegralExchangeAmount(integral int) int {
	if integral > 0 {
		dic := p.registryRepo.Get(registry.IntegralDiscountQuantity).IntValue()
		if dic > 0 {
			return int(float32(integral)/float32(dic)) * 100
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
	if !ignoreAmount && int64(amount) > int64(p.value.FinalAmount) {
		amount = int(p.value.FinalAmount)
		dic := p.registryRepo.Get(registry.IntegralDiscountQuantity).IntValue()
		integral = int(float32(amount) * float32(dic))
	}
	if amount <= 0 {
		return 0, nil
	}
	acc := p.memberRepo.GetMember(int64(p.value.BuyerId)).GetAccount()
	//log.Println("----", p.value.BuyerId, acc.Value().Integral, "discount:", integral)
	//log.Printf("-----%#v\n", acc.Value())
	_, err = acc.Discount(member.AccountIntegral, "积分支付抵扣",
		integral, p.Get().TradeNo, "")
	// 抵扣积分
	if err == nil {
		p.value.DeductAmount += amount
		p.value.FinalFlag |= payment.MIntegral
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(amount, payment.MIntegral, "", "")
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
		p.value.FinalFlag |= payment.MSystemPay
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(fee, payment.MSystemPay, "", "")
		}
	}
	return err
}

// 保存支付信息
func (p *paymentOrderImpl) saveTradeChan(amount int, method int, code string, outTradeNo string) error {
	c := &payment.PayTradeData{
		TradeNo:      p.TradeNo(),
		OrderId:      p.GetAggregateRootId(),
		PayMethod:    method,
		Internal:     1,
		PayAmount:    amount,
		OutTradeCode: code,
		OutTradeNo:   outTradeNo,
		PayTime:      int(time.Now().Unix()),
	}
	_, err := p.repo.SavePaymentTradeChan(p.TradeNo(), c)
	return err
}

func (p *paymentOrderImpl) getBuyer() member.IMemberAggregateRoot {
	if p.buyer == nil {
		p.buyer = p.memberRepo.GetMember(int64(p.value.BuyerId))
	}
	return p.buyer
}

func (p *paymentOrderImpl) intAmount(a float32) int {
	return int(a * float32(enum.RATE_AMOUNT))
}

// HybridPayment 余额钱包混合支付，优先扣除余额。
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
	if acc.Balance >= int(v.FinalAmount) {
		return p.BalanceDeduct(remark)
	}
	// 判断是否能钱包支付
	if !p.andMethod(p.value.PayFlag, payment.MWallet) {
		return payment.ErrNotSupportPaymentChannel
	}
	// 判断是否余额不足
	if acc.Balance+acc.WalletBalance < int(v.FinalAmount) {
		return payment.ErrNotEnoughAmount
	}
	err := p.BalanceDeduct(remark)
	if err == nil {
		err = p.PaymentByWallet(remark)
	}
	return err
}

// PaymentByWallet 钱包账户支付
func (p *paymentOrderImpl) PaymentByWallet(remark string) error {
	if !p.andMethod(p.value.PayFlag, payment.MWallet) {
		return payment.ErrNotSupportPaymentChannel
	}
	buyer := p.getBuyer()
	if buyer == nil {
		return member.ErrNoSuchMember
	}
	amount := p.value.FinalAmount
	// 判断并从钱包里扣款
	acc := buyer.GetAccount()
	if acc.GetValue().WalletBalance < int(amount) {
		return payment.ErrNotEnoughAmount
	}
	_, err := acc.Consume(member.AccountWallet, "支付订单",
		int(amount), p.TradeNo(), remark)
	if err == nil {
		p.value.DeductAmount += amount
		p.value.FinalFlag |= payment.MWallet
		err = p.saveOrder()
		if err == nil { // 保存支付记录
			err = p.saveTradeChan(int(amount), payment.MWallet, "", "")
		}
	}
	return err
}

// PaymentWithCard 使用会员卡支付,cardCode:会员卡编码,amount:支付金额
func (p *paymentOrderImpl) PaymentWithCard(cardCode string, amount int) error {
	return errors.New("not support")
}

// 保存订单
func (p *paymentOrderImpl) saveOrder() error {
	// 检查支付单
	err := p.checkOrderFinalAmount()
	if err == nil {
		p.value.UpdateTime = int(time.Now().Unix())
		p.value.Id, err = p.repo.SavePaymentOrder(p.value)
	}
	//保存支付单后,通知支付成功。只通知一次
	if err == nil && p.firstFinishPayment {
		p.firstFinishPayment = false
		err = p.applyPaymentFinish()
	}
	return err
}

// 退款
func (p *paymentOrderImpl) Refund(amounts map[int]int, reason string) (err error) {
	if p.State() != payment.StateFinished {
		return errors.New("订单未支付不支持退款")
	}
	chanMap := p.TradeMethods()
	if amounts == nil {
		// 如果未指定退款金额，则按支付的可退金额退款
		amounts = make(map[int]int, 0)
		for _, v := range chanMap {
			amounts[v.PayMethod] = v.PayAmount - v.RefundAmount
		}
	} else {
		// 验证退款金额
		for _, v := range chanMap {
			amount, ok := amounts[v.PayMethod]
			if ok && v.PayAmount-v.RefundAmount < amount {
				return errors.New("退款金额超出可退款金额")
			}
		}
	}

	mm := p.getBuyer()
	if mm == nil {
		return member.ErrNoSuchMember
	}
	pv := p.Get()
	acc := mm.GetAccount()
	if len(reason) == 0 {
		reason = "订单退款"
	}

	// 检查退款是否超出分账和已退款之和
	tx := collections.FindArray(chanMap, func(e *payment.PayTradeData) bool {
		return e.PayMethod == payment.MPaySP
	})
	if tx != nil {
		amount, ok := amounts[tx.PayMethod]
		if ok {
			err = p.checkDivideAmount(tx.PayAmount, amount+tx.RefundAmount)
		}
		if err != nil {
			return err
		}
	}

	totalRefund := 0
	for _, v := range chanMap {
		amount, ok := amounts[v.PayMethod]
		if !ok {
			// 如果未使用该支付方式，则跳过
			continue
		}
		switch v.PayMethod {
		case payment.MBalance:
			err = acc.Refund(member.AccountBalance, reason, amount, pv.TradeNo, "")
		case payment.MWallet:
			err = acc.Refund(member.AccountWallet, reason, amount, pv.TradeNo, "")
		case payment.MPaySP:
			// 处理充值退款
			var txId int
			txId, err = p.handlePaymentOrderRefund(acc, totalRefund, reason)
			if err == nil {
				// 第三方支付原路退回, 异步发布退款事件
				go eventbus.Publish(&payment.PaymentProviderRefundEvent{
					Order:        p,
					Amount:       amount,
					Reason:       reason,
					OutTradeCode: v.OutTradeCode,
					OutTradeNo:   v.OutTradeNo,
					AccountTxId:  txId,
				})
			}
		}
		if err == nil {
			v.RefundAmount += amount
			_, err = p.repo.SavePaymentTradeChan(p.TradeNo(), v)
		}
		if err != nil {
			return err
		}
		totalRefund += amount
	}
	if totalRefund > 0 {
		// 更新支付单退单金额
		p.value.RefundAmount += totalRefund
		p.value.UpdateTime = int(time.Now().Unix())
		if p.value.FinalAmount == p.value.RefundAmount {
			// 如果全部退款，则标记支付单状态为：已退款
			p.value.Status = payment.StateRefunded
		}
		// 检查分账并更新状态
		p.checkDivideStatus(tx)
		_, err = p.repo.SavePaymentOrder(p.value)
	}
	return err
}

// 处理退款业务
func (p *paymentOrderImpl) handlePaymentOrderRefund(acc member.IAccount, totalRefund int, reason string) (txId int, err error) {
	if p.value.OrderType == payment.TypeRecharge {
		// 如果是充值订单，则需扣除充值的金额
		return acc.Consume(member.AccountWallet, reason, totalRefund, p.TradeNo(), "")
	}
	return 0, nil
}

// RefundAvail 请求退款全部可退金额，通常用于全额退款或消费后将剩余部分进行退款
func (p *paymentOrderImpl) RefundAvail(remark string) (amount int, err error) {
	if p.State() != payment.StateFinished {
		return 0, errors.New("订单未支付不支持退款")
	}
	if len(remark) == 0 {
		return 0, errors.New("缺少退款备注")
	}
	// 获取支付数据
	chanMap := p.TradeMethods()
	// 计算第三方支付可退金额
	spRefundAmount := 0
	tx := collections.FindArray(chanMap, func(e *payment.PayTradeData) bool {
		return e.PayMethod == payment.MPaySP
	})
	if tx != nil {
		divides := p.getDivides()
		divideAmount := 0
		for _, v := range divides {
			divideAmount += v.DivideAmount
		}
		amount := tx.PayAmount - tx.RefundAmount - divideAmount
		if amount < 0 {
			return 0, fmt.Errorf("第三方支付超出可退款金额: %.2f, 已处理金额: %.2f",
				float64(tx.PayAmount)/100, float64(divideAmount+tx.RefundAmount)/100)
		}
		spRefundAmount = amount
	}
	// 获取会员及账户
	mm := p.getBuyer()
	if mm == nil {
		return 0, member.ErrNoSuchMember
	}
	pv := p.Get()
	acc := mm.GetAccount()

	totalRefund := 0
	for _, v := range chanMap {
		amount := v.PayAmount - v.RefundAmount
		if amount <= 0 {
			continue
		}
		switch v.PayMethod {
		case payment.MBalance:
			err = acc.Refund(member.AccountBalance, remark, amount, pv.TradeNo, "")
		case payment.MWallet:
			err = acc.Refund(member.AccountWallet, remark, amount, pv.TradeNo, "")
		case payment.MPaySP:
			if spRefundAmount > 0 {
				// 第三方支付退款金额
				amount = spRefundAmount
				// 处理充值退款
				var txId int
				txId, err = p.handlePaymentOrderRefund(acc, totalRefund, remark)
				if err == nil {
					// 第三方支付原路退回, 异步发布退款事件
					go eventbus.Publish(&payment.PaymentProviderRefundEvent{
						Order:        p,
						Amount:       spRefundAmount,
						Reason:       remark,
						OutTradeCode: v.OutTradeCode,
						OutTradeNo:   v.OutTradeNo,
						AccountTxId:  txId,
					})
				}
			}
		}
		if err == nil {
			v.RefundAmount += amount
			_, err = p.repo.SavePaymentTradeChan(p.TradeNo(), v)
		}
		if err != nil {
			return 0, err
		}
		totalRefund += amount
	}
	if totalRefund > 0 {
		// 更新支付单退单金额
		p.value.RefundAmount += totalRefund
		p.value.UpdateTime = int(time.Now().Unix())
		if p.value.FinalAmount == p.value.RefundAmount {
			// 如果全部退款，则标记支付单状态为：已退款
			p.value.Status = payment.StateRefunded
		}
		// 检查分账并更新状态
		p.checkDivideStatus(tx)
		_, err = p.repo.SavePaymentOrder(p.value)
	}
	if err != nil {
		return 0, err
	}
	return totalRefund, nil
}

func (p *paymentOrderImpl) SupplementRefund(txId int) error {
	// 获取支付数据
	chanMap := p.TradeMethods()
	// 计算第三方支付可退金额
	tx := collections.FindArray(chanMap, func(e *payment.PayTradeData) bool {
		return e.PayMethod == payment.MPaySP
	})
	if tx == nil {
		return errors.New("支付单未使用第三方支付")
	}
	// 验证会员和交易状态
	m := p.getBuyer()
	if m == nil {
		return member.ErrNoSuchMember
	}
	acc := m.GetAccount()
	if acc == nil {
		return errors.New("会员账户不存在")
	}
	mtx := acc.GetWalletLog(int64(txId))
	if mtx.Id <= 0 {
		return errors.New("找不到交易记录")
	}
	if mtx.ReviewStatus == wallet.ReviewCompleted {
		return errors.New("交易已完成,无需补发")
	}

	// 第三方支付原路退回, 异步发布退款事件
	go eventbus.Publish(&payment.PaymentProviderRefundEvent{
		Order:        p,
		Amount:       mtx.ChangeValue - mtx.TransactionFee,
		Reason:       mtx.Subject,
		OutTradeCode: tx.OutTradeCode,
		OutTradeNo:   tx.OutTradeNo,
		AccountTxId:  txId,
	})
	return nil
}

// 检查分账状态，只有第三方支付的部分能参与分账
func (p *paymentOrderImpl) checkDivideStatus(tx *payment.PayTradeData) {
	if !p.isDivide() {
		return
	}
	dividedList := p.getDivides()
	// 统计分账的总金额
	dividedAmount := 0
	for _, k := range dividedList {
		dividedAmount += k.DivideAmount
	}
	if dividedAmount+tx.RefundAmount >= tx.PayAmount {
		// 如果分账金额+已退款金额等于订单金额，则自动标记为分账完成
		if dividedAmount == 0 {
			// 如果分账金额为0，则标记为分账成功
			p.value.DivideStatus = payment.DivideItemStatusSuccess
		} else {
			// 如果分账金额不为0，则标记为分账完成,等待下发分账指令
			p.value.DivideStatus = payment.DivideCompleted
		}
	}
}

// getDivides 获取分账列表
func (p *paymentOrderImpl) getDivides() []*payment.PayDivide {
	return p.repo.DivideRepo().FindList(nil, "pay_id = ?", p.GetAggregateRootId())
}

// checkDivideAmount 检查退款金额是否超出分账外的金额
func (p *paymentOrderImpl) checkDivideAmount(orderAmount int, refundTotalAmount int) error {
	if refundTotalAmount > orderAmount {
		return errors.New("退款金额超出可退款金额")
	}
	divides := p.getDivides()
	totalDiviedAmount := 0
	for _, v := range divides {
		totalDiviedAmount += v.DivideAmount
	}
	if orderAmount-totalDiviedAmount < refundTotalAmount {
		return errors.New("退款金额超出可退款金额(含分账)")
	}
	return nil
}

// 调整金额,如调整金额与实付金额相加小于等于零,则支付成功。
func (p *paymentOrderImpl) Adjust(amount int) error {
	p.value.AdjustAmount += amount
	p.value.FinalAmount += amount
	if p.value.FinalAmount <= 0 {
		return p.checkOrderFinalAmount()
	}
	return p.saveOrder()
}

// 获取支付途径支付信息字典
func (p *paymentOrderImpl) getPaymentChannelMap() map[int]int {
	mp := make(map[int]int, 0)
	arr := p.TradeMethods()
	for _, v := range arr {
		if v.PayAmount > 0 {
			c, ok := mp[v.PayMethod]
			if ok {
				mp[v.PayMethod] = c + v.PayAmount
			} else {
				mp[v.PayMethod] = v.PayAmount
			}
		}
	}
	return mp
}

func (p *paymentOrderImpl) ChanName(method int) string {
	switch method {
	case payment.MBalance:
		return "余额"
	case payment.MWallet:
		return "钱包"
	case payment.MIntegral:
		return "积分"
	case payment.MUserCard:
		return "会员卡"
	case payment.MUserCoupon:
		return "券"
	case payment.MCash:
		return "现金"
	case payment.MBankCard:
		return "刷卡"
	case payment.MPaySP:
		return "第三方"
	case payment.MSellerPay:
		return "卖家"
	case payment.MSystemPay:
		return "系统"
	}
	return fmt.Sprintf("未知的支付方式%d", method)
}

type RepoBase struct {
}

func (p *RepoBase) CreatePaymentOrder(v *payment.
	Order, repo payment.IPaymentRepo, mmRepo member.IMemberRepo,
	registryRepo registry.IRegistryRepo) payment.IPaymentOrder {
	return &paymentOrderImpl{
		repo:       repo,
		value:      v,
		memberRepo: mmRepo,
		//orderManager: orderManager,
		registryRepo: registryRepo,
	}
}

// IsDivide 是否支持分账
func (p *paymentOrderImpl) isDivide() bool {
	return p.value.AttrFlag&int(payment.FlagDivide) == int(payment.FlagDivide)
}

// Divide implements payment.IPaymentOrder.
func (p *paymentOrderImpl) Divide(outTxNo string, divides []*payment.DivideData) error {
	if !p.isDivide() {
		return errors.New("支付单不支持分账")
	}
	repo := p.repo.DivideRepo()
	if p.value.DivideStatus == payment.DivideCompleted {
		return errors.New("订单已分账完成")
	}
	if len(outTxNo) == 0 {
		return errors.New("缺少分账关联的外部交易单号")
	}
	dividedList := p.getDivides()
	// 统计将分账的总金额
	divideAmount := 0
	dividedAmount := 0
	for i, v := range divides {
		if v.DivideType != 1 && v.UserId <= 0 {
			return errors.New("分账用户编号不能为空")
		}
		divideAmount += v.DivideAmount
		for _, k := range dividedList {
			if i == 0 {
				dividedAmount += k.DivideAmount
			}
			if k.OutTxNo == outTxNo && k.UserId == v.UserId && k.DivideType == v.DivideType {
				return errors.New("同一交易不允许用户重复分账")
			}
		}
	}

	if divideAmount+dividedAmount > p.value.FinalAmount {
		return errors.New("超出订单可分账总额")
	}
	unix := int(time.Now().Unix())
	for _, v := range divides {
		pv := &payment.PayDivide{
			Id:           0,
			PayId:        p.GetAggregateRootId(),
			DivideType:   v.DivideType,
			UserId:       v.UserId,
			DivideAmount: v.DivideAmount,
			OutTxNo:      outTxNo,
			Remark:       "",
			SubmitStatus: 1,
			SubmitRemark: "",
			SubmitTime:   0,
			CreateTime:   unix,
		}
		if v.DivideType == payment.DivideUserPlatform {
			// 如果是平台，则默认已经分账
			pv.SubmitStatus = payment.DivideItemStatusSuccess
			pv.SubmitRemark = "系统通过"
		}
		ret, err := repo.Save(pv)
		if err != nil {
			return err
		}
		// 更新分账明细ID
		v.DivideItemId = ret.Id
	}
	var err error
	if divideAmount+dividedAmount+p.value.RefundAmount >= p.value.FinalAmount {
		// 如果分账金额+已退款金额等于订单金额，则自动标记为分账完成
		err = p.CompleteDivide()
	} else if p.value.DivideStatus == payment.DivideNoDivide {
		// 修改状态为待分账
		p.value.DivideStatus = payment.DividePending
		p.value.UpdateTime = unix
		_, err = p.repo.SavePaymentOrder(p.value)
	}
	if err == nil {
		// 发布分账事件
		eventbus.Publish(&payment.PaymentDivideEvent{
			Order:   p,
			Divides: divides,
		})
	}
	return err
}

// CompleteDivide implements payment.IPaymentOrder.
func (p *paymentOrderImpl) CompleteDivide() error {
	if !p.isDivide() {
		return errors.New("支付单不支持分账")
	}
	if p.value.DivideStatus == payment.DivideSuccess {
		return errors.New("支付平台已进行分账")
	}
	if p.value.DivideStatus != payment.DivideCompleted {
		// 如果金额全部分账，则自动标记为分账完成
		p.value.DivideStatus = payment.DivideCompleted
		p.value.UpdateTime = int(time.Now().Unix())
		_, err := p.repo.SavePaymentOrder(p.value)
		if err != nil {
			return err
		}
	}
	// 检查分账命令是否执行
	p.checkDivideCommandExecuted()
	return nil
}

// DivideSuccess 分账成功,支付平台已进行分账
func (p *paymentOrderImpl) DivideSuccess(outTxNo string) error {
	if !p.isDivide() {
		return errors.New("支付单不支持分账")
	}
	if p.value.DivideStatus == payment.DivideSuccess {
		return errors.New("支付平台已进行分账")
	}
	p.value.DivideStatus = payment.DivideSuccess
	p.value.UpdateTime = int(time.Now().Unix())
	_, err := p.repo.SavePaymentOrder(p.value)
	return err
}

// UpdateSubDivideStatus implements payment.IPaymentOrder.
func (p *paymentOrderImpl) UpdateSubDivideStatus(divideId int, success bool, divideNo string, remark string) error {
	divide := p.repo.DivideRepo().Get(divideId)
	if divide == nil {
		return errors.New("分账记录不存在")
	}
	if divide.PayId != p.GetAggregateRootId() {
		return errors.New("分账记录不属于当前订单")
	}
	if divide.SubmitStatus != payment.DividePending && divide.SubmitStatus != payment.DivideItemStatusReverted {
		// 只有待分账及撤销状态下允许更改状态
		return errors.New("分账记录状态错误")
	}
	divide.SubmitStatus = types.Ternary(success, payment.DivideItemStatusSuccess, payment.DivideItemStatusFailed)
	divide.SubmitRemark = remark
	divide.SubmitTime = int(time.Now().Unix())
	divide.SubmitDivideNo = divideNo
	_, err := p.repo.DivideRepo().Save(divide)
	if err == nil {
		// 检查分账命令是否执行
		p.checkDivideCommandExecuted()
	}
	return err
}

// checkDivideCommandExecuted 检查分账命令是否执行,如果执行，则执行分账完成指令(SP)
// 分账完成的条件： 所有分账都为成功状态, 且订单没有金额再用于分账
// 调用场景：
// 1. 手动完成分账
// 2. 更新分账子项状态时
func (p *paymentOrderImpl) checkDivideCommandExecuted() {
	if p.value.DivideStatus != payment.DivideCompleted {
		// 未标记为已经分账完成的订单，不下发分账完成指令
		return
	}
	isExecuted := true
	divides := p.getDivides()
	if len(divides) == 0 {
		// 没有分账记录
		return
	}
	for _, v := range divides {
		if v.SubmitStatus != payment.DivideItemStatusSuccess {
			isExecuted = false
			break
		}
	}
	if isExecuted {
		// 如果分账命令都已下发，则执行分账完成指令
		go eventbus.Publish(&payment.PaymentCompleteDivideEvent{
			Order: p,
		})
	}
}

// RevertSubDivide implements payment.IPaymentOrder.
func (p *paymentOrderImpl) RevertSubDivide(divideId int, remark string) error {
	divide := p.repo.DivideRepo().Get(divideId)
	if divide == nil {
		return errors.New("分账记录不存在")
	}
	if divide.PayId != p.GetAggregateRootId() {
		return errors.New("分账记录不属于当前订单")
	}
	if divide.SubmitStatus != payment.DivideItemStatusSuccess {
		return errors.New("分账未成功状态，不支持撤销")
	}
	raw := types.DeepClone(divide)
	divide.SubmitStatus = payment.DivideItemStatusReverted
	divide.SubmitRemark = remark
	divide.SubmitTime = int(time.Now().Unix())
	_, err := p.repo.DivideRepo().Save(divide)
	if err == nil {
		// 发布分账撤销事件
		eventbus.Publish(&payment.PaymentRevertSubDivideEvent{
			Order:   p,
			Divides: []*payment.PayDivide{raw},
		})
	}
	return err
}
