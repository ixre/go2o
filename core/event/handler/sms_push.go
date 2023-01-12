package handler

import (
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/tool/sms"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

// 处理发送短信事件
func (e EventHandler) HandleSendSmsEvent(data interface{}) {
	v := data.(*events.SendSmsEvent)
	if v != nil {
		v.Template = sms.ResolveMessage(v.Template, v.Data)
		ev := &proto.EVSendSmsEventData{
			Provider:   int32(v.Provider),
			Phone:      v.Phone,
			Template:   v.Template,
			TemplateId: v.TemplateId,
			Data:       v.Data,
		}
		msq.Push(msq.SendSmsTopic, typeconv.MustJson(ev))
	}
}
