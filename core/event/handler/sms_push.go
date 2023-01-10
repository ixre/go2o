package handler

import (
	"encoding/json"

	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/service/proto"
)

func (e EventHandler) HandleSendSmsEvent(data interface{}) {
	v := data.(*events.SendSmsEvent)
	ev := &proto.EVSendSmsEventData{
		Provider:   int32(v.Provider),
		Phone:      v.Phone,
		Template:   v.Template,
		TemplateId: v.TemplateId,
		Data:       v.Data,
	}
	bytes, _ := json.Marshal(ev)
	msq.Push(msq.SystemSendSmsTopic, string(bytes))
}
