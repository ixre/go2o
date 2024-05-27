package handler

import (
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/msq"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func (h EventHandler) HandleRegistryPushEvent(data interface{}) {
	v := data.(*registry.RegistryPushEvent)
	if v != nil {
		ev := &proto.EVRegistryPushEventData{
			Key:     v.Key,
			Value:   v.Value,
			IsUsers: v.IsUser,
		}
		msq.Push(msq.RegistryTopic, typeconv.MustJson(ev))
	}
}
