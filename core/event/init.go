package event

import (
	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/event/handler"
	"github.com/ixre/gof/domain/eventbus"
)

type EventSource struct {
	*handler.EventHandler
	*handler.PaymentEventHandler
	*handler.MerchantEventHandler
	*handler.InvoiceEventHandler
}

func NewEventSource(h *handler.EventHandler,
	p *handler.PaymentEventHandler,
	m *handler.MerchantEventHandler,
	i *handler.InvoiceEventHandler,
) *EventSource {
	return &EventSource{
		EventHandler:         h,
		PaymentEventHandler:  p,
		MerchantEventHandler: m,
		InvoiceEventHandler:  i,
	}
}

func (e *EventSource) Bind() {
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

	// 注册商户事件
	e.initMchEvents()
	// 注册发票事件
	e.initInvoiceEvents()

	// 注册审批事件
	eventbus.Subscribe(approval.ApprovalProcessEvent{}, h.OnApprovalProcess)
	// 注册支付成功事件
	eventbus.Subscribe(payment.PaymentSuccessEvent{}, e.HandlePaymentSuccessEvent)
	// 注册支付分账撤销事件
	eventbus.Subscribe(payment.PaymentRevertSubDivideEvent{}, e.HandlePaymentSubDivideRevertEvent)
	// 注册支付分账事件
	eventbus.Subscribe(payment.PaymentDivideEvent{}, e.HandlePaymentDivideEvent)
	// 注册第三方支付退款事件
	eventbus.Subscribe(payment.PaymentProviderRefundEvent{}, e.HandlePaymentProviderRefundEvent)
	// 注册支付完成分账事件
	eventbus.Subscribe(payment.PaymentCompleteDivideEvent{}, e.HandlePaymentCompleteDivideEvent)
	// 注册支付商户入网事件
	eventbus.Subscribe(payment.PaymentMerchantRegistrationEvent{}, e.HandlePaymentMerchantRegistrationEvent)

	// 注册员工IM初始化事件
	eventbus.Subscribe(staff.StaffRequireImInitEvent{}, e.HandleStaffRequireImInitEvent)
}

// 注册商户事件
func (e *EventSource) initMchEvents() {
	// 注册商户结算事件
	eventbus.Subscribe(merchant.MerchantBillSettleEvent{}, e.HandleMerchantBillSettleEvent)
	// 注册员工转移审批通过事件
	eventbus.Subscribe(staff.StaffTransferApprovedEvent{}, e.HandleStaffTransferApprovedEvent)
}

// 注册发票事件
func (e *EventSource) initInvoiceEvents() {
	// 注册发票撤销事件
	eventbus.Subscribe(invoice.InvoiceRevertEvent{}, e.HandleInvoiceRevertEvent)
}
