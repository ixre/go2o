package handler

import (
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/msq"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func (h *EventHandler) HandleMemberPushEvent(data interface{}) {
	v := data.(*events.MemberPushEvent)
	if v == nil {
		return
	}
	m := v.Member
	ev := &proto.EVMemberPushEventData{
		MemberId:      m.Id,
		IsNewMember:   v.IsCreate,
		UserCode:      m.UserCode,
		UserFlag:      int32(m.UserFlag),
		Username:      m.Username,
		Exp:           int64(m.Exp),
		Level:         int32(m.Level),
		Nickname:      m.Nickname,
		Portrait:      m.Portrait,
		Phone:         m.Phone,
		Email:         m.Email,
		RegFrom:       m.RegFrom,
		InviterId:     int64(v.InviterId),
		RealName:      m.RealName,
		LastLoginTime: m.LastLoginTime,
	}

	msq.Push(msq.MemberUpdated, typeconv.MustJson(ev))
}
