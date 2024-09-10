package event

import (
	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/event/handler"
	"github.com/ixre/gof/domain/eventbus"
)

type EventSource struct {
	*handler.EventHandler
	*handler.PaymentEventHandler
}

func NewEventSource(h *handler.EventHandler, p *handler.PaymentEventHandler) *EventSource {
	return &EventSource{
		EventHandler:        h,
		PaymentEventHandler: p,
	}
}

func (e *EventSource) Init() {
	h := e.EventHandler
	eventbus.SubscribeAsync(events.AppInitialEvent{}, h.HandleAppInitialEvent)
	eventbus.SubscribeAsync(registry.RegistryPushEvent{}, h.HandleRegistryPushEvent)
	eventbus.SubscribeAsync(events.AccountLogPushEvent{}, h.HandleMemberAccountLogPushEvent)
	eventbus.SubscribeAsync(events.OrderAffiliateRebateEvent{}, h.HandleOrderAffiliateRebateEvent)
	eventbus.SubscribeAsync(events.SendSmsEvent{}, h.HandleSendSmsEvent)
	eventbus.SubscribeAsync(events.SubOrderPushEvent{}, h.HandleSubOrderPushEvent)
	eventbus.SubscribeAsync(events.MemberPushEvent{}, h.HandleMemberPushEvent)
	eventbus.SubscribeAsync(events.MemberAccountPushEvent{}, h.HandleMemberAccountPushEvent)
	eventbus.SubscribeAsync(events.WithdrawalPushEvent{}, h.HandleWithdrawalPushEvent)

	// 注册审批事件
	eventbus.Subscribe(approval.ApprovalProcessEvent{}, h.OnApprovalProcess)
	// 注册支付成功事件
	eventbus.Subscribe(payment.PaymentSuccessEvent{}, e.HandlePaymentSuccessEvent)
	// 注册支付分账撤销事件
	eventbus.Subscribe(payment.PaymentSubDivideRevertEvent{}, e.HandlePaymentSubDivideRevertEvent)
	// 注册支付分账事件
	eventbus.Subscribe(payment.PaymentDivideEvent{}, e.HandlePaymentDivideEvent)
	// 注册第三方支付退款事件
	eventbus.Subscribe(payment.PaymentProviderRefundEvent{}, e.HandlePaymentProviderRefundEvent)
}
