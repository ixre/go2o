package event

import (
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/event/handler"
	"github.com/ixre/gof/domain/eventbus"
)

func InitEvent() {
	h := &handler.EventHandler{}
	eventbus.SubscribeAsync(events.WalletLogClickhouseUpdateEvent{}, h.HandleWalletLogWriteEvent)
	eventbus.SubscribeAsync(events.OrderAffiliteRebateEvent{},h.HandleOrderAffiliteRebateEvent)
}
