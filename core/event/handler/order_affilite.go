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
func (h EventHandler) HandleOrderAffiliateRebateEvent(data interface{}) {
	v := data.(*events.OrderAffiliateRebateEvent)
	r := repos.Repo.GetRegistryRepo()
	s, _ := r.GetValue(registry.OrderPushAffiliateEvent)
	pushValue, _ := strconv.Atoi(s)
	//todo: 系统内处理分销，不推送分销事件
	if pushValue == 0 {

	}
	ev := &proto.EVOrderAffiliateEventData{
		BuyerId:        v.BuyerId,
		SubOrder:       true,
		OrderNo:        v.OrderNo,
		OrderAmount:    v.OrderAmount,
		AffiliateItems: []*proto.EVOrderAffiliateItem{},
	}
	for _, v := range v.AffiliateItems {
		//todo: 实现商品自定义分销比例
		ev.AffiliateItems = append(ev.AffiliateItems, &proto.EVOrderAffiliateItem{
			ItemId:      v.ItemId,
			SkuId:       v.SkuId,
			Quantity:    v.Quantity,
			Amount:      v.Amount,
			FinalAmount: v.FinalAmount,
			Params:      []*proto.EVItemAffiliateConfig{},
		})
	}
	bytes, _ := json.Marshal(ev)
	// 推送至外部系统，并由外部系统处理分销
	if pushValue == 1 {
		msq.Push(msq.OrderAffiliateTopic, string(bytes))
		return
	}

	//todo: 处理分销后将事件推送至外部系统
}
