package handler

import (
	"github.com/ixre/go2o/pkg/event/msq"
	"github.com/ixre/go2o/pkg/interface/domain/registry"
	"github.com/ixre/go2o/pkg/interface/service/proto"
	"github.com/ixre/gof/typeconv"
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
