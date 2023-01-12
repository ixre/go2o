package handler

import (
	"log"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/repos"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

// 子订单推送
func (h EventHandler) HandleSubOrderPushEvent(data interface{}) {
	v := data.(*events.SubOrderPushEvent)
	if v == nil {
		return
	}
	r := repos.Repo.GetRegistryRepo()
	s, _ := r.GetValue(registry.OrderPushAffiliateEvent)
	pushValue, _ := strconv.Atoi(s)
}

// 订单分销处理
func (h EventHandler) HandleOrderAffiliateRebateEvent(data interface{}) {
	v := data.(*events.OrderAffiliateRebateEvent)
	if v == nil {
		return
	}
	r := repos.Repo.GetRegistryRepo()
	s, _ := r.GetValue(registry.OrderPushAffiliateEvent)
	pushValue, _ := strconv.Atoi(s)
	//todo: 系统内处理分销，不推送分销事件
	// 0:不推送(内部处理),1:仅推送(内部处理),2:推送并处理(外部处理分销)
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
	// 推送至外部系统，并由外部系统处理分销
	if pushValue == 2 {
		err := msq.Push(msq.OrderAffiliateTopic, typeconv.MustJson(ev))
		if err != nil {
			log.Println("[ go2o][ event]: push order affiliate event failed, error: ", err.Error())
		}
		return
	}

	//todo: 处理分销后将事件推送至外部系统
}
