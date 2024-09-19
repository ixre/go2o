/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: payment_event.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-09-08 10:45:16
 * description: 支付事件处理
 * history:
 */

package handler

import (
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/logger"
)

// PaymentEventHandler 支付事件处理
type PaymentEventHandler struct {
	_orderRepo  order.IOrderRepo
	_memberRepo member.IMemberRepo
}

func NewPaymentEventHandler(memberRepo member.IMemberRepo, orderRepo order.IOrderRepo) *PaymentEventHandler {
	return &PaymentEventHandler{
		_orderRepo:  orderRepo,
		_memberRepo: memberRepo,
	}
}

// 处理支付成功事件
func (p *PaymentEventHandler) HandlePaymentSuccessEvent(event interface{}) {
	e := event.(*payment.PaymentSuccessEvent)
	ov := e.Order.Get()
	switch ov.OrderType {
	case payment.TypeOrder:
		p.handleOrderSuccessEvent(e)
	case payment.TypeRecharge:
		p.handleRechargeSuccessEvent(e)
	}
}

// 处理支付分账事件
func (p *PaymentEventHandler) HandlePaymentDivideEvent(event interface{}) {
	//e := event.(*payment.PaymentDivideEvent)
	// note: 支付分账事件由具体的支付渠道通过订阅事件处理，这里不作处理
}

// 处理支付分账撤销事件
func (p *PaymentEventHandler) HandlePaymentSubDivideRevertEvent(event interface{}) {
	//e := event.(*payment.PaymentSubDivideRevertEvent)
	// note: 支付分账撤销事件由具体的支付渠道通过订阅事件处理，这里不作处理
}

// 处理第三方支付退款事件
func (p *PaymentEventHandler) HandlePaymentProviderRefundEvent(event interface{}) {
	//e := event.(*payment.PaymentProviderRefundEvent)
	// note: 第三方支付退款事件由具体的支付渠道通过订阅事件处理，这里不作处理
}

// 处理支付完成分账事件
func (p *PaymentEventHandler) HandlePaymentCompleteDivideEvent(event interface{}) {
	//e := event.(*payment.PaymentCompleteDivideEvent)
	// note: 支付完成分账事件由具体的支付渠道通过订阅事件处理，这里不作处理
}

// 处理支付商户入网事件
func (p *PaymentEventHandler) HandlePaymentMerchantRegistrationEvent(event interface{}) {
	//e := event.(*payment.PaymentMerchantRegistrationEvent)
	// note: 支付商户入网事件由具体的支付渠道通过订阅事件处理，这里不作处理
}

// 处理商城订单支付完成
func (p *PaymentEventHandler) handleOrderSuccessEvent(e *payment.PaymentSuccessEvent) {
	ov := e.Order.Get()
	// 通知订单支付完成
	if ov.OutOrderNo != "" {
		subOrder := ov.SubOrder == 1
		err := p._orderRepo.Manager().NotifyOrderTradeSuccess(ov.OutOrderNo, subOrder)
		domain.HandleError(err, "domain")
	}
}

// 处理充值订单支付完成
func (p *PaymentEventHandler) handleRechargeSuccessEvent(e *payment.PaymentSuccessEvent) {
	ov := e.Order.Get()
	m := p._memberRepo.GetMember(int64(ov.BuyerId))
	if m == nil {
		logger.Error("recharge success but member is nil, payment order id: %d", e.Order.GetAggregateRootId())
		return
	}
	err := m.GetAccount().Charge(member.AccountWallet,
		ov.Subject,
		int(ov.TotalAmount),
		ov.TradeNo,
		"支付充值",
	)
	if err != nil {
		logger.Error("处理用户充值失败，错误信息: %s, 支付单ID:%d", err.Error(), e.Order.GetAggregateRootId())
		return
	}
	logger.Debug("处理用户充值成功, 支付单ID:%d", e.Order.GetAggregateRootId())
}
