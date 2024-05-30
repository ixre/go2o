package handler

import (
	"github.com/ixre/go2o/core/event/events"
)

// 子订单推送
func (h EventHandler) HandleAppInitialEvent(data interface{}) {
	v := data.(*events.AppInitialEvent)
	if v == nil {
		return
	}
}
