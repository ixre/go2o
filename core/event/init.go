package event

import (
	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/event/handler"
	"github.com/ixre/gof/domain/eventbus"
)

type EventSource struct {
	*handler.EventHandler
}

func NewEventSource(h *handler.EventHandler) *EventSource {
	return &EventSource{
		EventHandler: h,
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

	eventbus.Subscribe(approval.ApprovalProcessEvent{}, h.OnApprovalProcess)
}
