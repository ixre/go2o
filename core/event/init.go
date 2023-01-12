package event

import (
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/event/handler"
	"github.com/ixre/gof/domain/eventbus"
)

func InitEvent() {
	h := &handler.EventHandler{}
	eventbus.SubscribeAsync(registry.RegistryPushEvent{}, h.HandleRegistryPushEvent)
	eventbus.SubscribeAsync(events.WalletLogClickhouseUpdateEvent{}, h.HandleWalletLogWriteEvent)
	eventbus.SubscribeAsync(events.OrderAffiliateRebateEvent{}, h.HandleOrderAffiliateRebateEvent)
	eventbus.SubscribeAsync(events.SendSmsEvent{}, h.HandleSendSmsEvent)
	eventbus.SubscribeAsync(events.MemberPushEvent{}, h.HandleMemberPushEvent)
	eventbus.SubscribeAsync(events.MemberAccountPushEvent{}, h.HandleMemberAccountPushEvent)
	eventbus.SubscribeAsync(events.WithdrawalPushEvent{},h.HandleWithdrawalPushEvent)
}
