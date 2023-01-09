package handler

import (
	"encoding/json"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/go2o/core/service/proto"
)

// 订单分销处理
func (h EventHandler) HandleOrderAffiliteRebateEvent(data interface{}) {
	v := data.(*events.OrderAffiliteRebateEvent)
	r := repos.Repo.GetRegistryRepo()
	s, _ := r.GetValue(registry.OrderPushAffiliteEvent)
	pushValue, _ := strconv.Atoi(s)
	//todo: 系统内处理分销，不推送分销事件
	if pushValue == 0 {

	}
	ev := &proto.EVOrderAffiliteRebateOrder{
		OrderNo:       v.OrderNo,
		OrderAmount:   v.OrderAmount,
		AffiliteItems: []*proto.EVOrderAffiliteItem{},
	}
	for _, v := range v.AffiliteItems {
		//todo: 实现商品自定义分销比例
		ev.AffiliteItems = append(ev.AffiliteItems, &proto.EVOrderAffiliteItem{
			ItemId:      v.ItemId,
			SkuId:       v.SkuId,
			Quantity:    v.Quantity,
			Amount:      v.Amount,
			FinalAmount: v.FinalAmount,
			Params:      []*proto.EVItemAffiliteConfig{},
		})
	}
	bytes, _ := json.Marshal(ev)
	// 推送至外部系统，并由外部系统处理分销
	if pushValue == 1 {
		msq.Push(msq.OrderAffiliteTopic, string(bytes))
		return
	}

	//todo: 处理分销后将事件推送至外部系统
}
