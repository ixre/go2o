package handler

import (
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/repos"
)

// 订单分销处理
func (h EventHandler) HandleOrderAffiliteRebateEvent(*events.OrderAffiliteRebateEvent) {
	r := repos.Repo.GetRegistryRepo()
	s, _ := r.GetValue(registry.OrderPushAffiliteEvent)
	pushValue, _ := strconv.Atoi(s)
	//todo: 系统内处理分销，不推送分销事件
	if pushValue == 0 {

	}
	// 推送至外部系统，并由外部系统处理分销
	if pushValue == 1 {
		
	}

	//todo: 处理分销后将事件推送至外部系统
}
