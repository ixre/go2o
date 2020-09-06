/**
 * Copyright 2015 @ to2.net.
 * name : message_service.go
 * author : jarryliu
 * date : 2016-06-11 20:51
 * description :
 * history :
 */
package impl

import (
	"context"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/dto"
	"go2o/core/service/proto"
)

var _ proto.MessageServiceServer = new(messageService)

type messageService struct {
	_rep mss.IMssRepo
	serviceUtil
}



func NewMessageService(rep mss.IMssRepo) *messageService {
	return &messageService{
		_rep: rep,
	}
}

// 获取通知项配置
func (m *messageService) GetNotifyItem(_ context.Context,key *proto.String) (*proto.SNotifyItem, error) {
	it := m._rep.NotifyManager().GetNotifyItem(key.Value)
	return &proto.SNotifyItem{
		Key:        it.Key,
		NotifyBy:   int32(it.NotifyBy),
		ReadonlyBy: it.ReadonlyBy,
		TplId:      int32(it.TplId),
		Content:    it.Content,
		Tags:       it.Tags,
	}, nil
}

// 发送短信
func (m *messageService) SendPhoneMessage(_ context.Context,r *proto.SendMessageRequest) (*proto.Result, error) {
	mg := m._rep.NotifyManager()
	extra := make(map[string]interface{})
	if r.Data != nil {
		for k, v := range r.Data {
			extra[k] = v
		}
	}
	err := mg.SendPhoneMessage(r.Account, notify.PhoneMessage(r.Message), extra)
	if err != nil {
		return m.error(err), nil
	}
	return m.success(nil), nil
}

// 获取邮件模版
func (m *messageService) GetMailTemplate(id int32) *mss.MailTemplate {
	return m._rep.GetProvider().GetMailTemplate(id)
}

// 保存邮件模板
func (m *messageService) SaveMailTemplate(v *mss.MailTemplate) (int32, error) {
	return m._rep.GetProvider().SaveMailTemplate(v)
}

// 获取邮件模板
func (m *messageService) GetMailTemplates() []*mss.MailTemplate {
	return m._rep.GetProvider().GetMailTemplates()
}

// 删除邮件模板
func (m *messageService) DeleteMailTemplate(id int32) error {
	return m._rep.GetProvider().DeleteMailTemplate(id)
}

// 获取邮件绑定
func (m *messageService) GetConfig() mss.Config {
	return m._rep.GetProvider().GetConfig()
}

// 保存邮件
func (m *messageService) SaveConfig(conf *mss.Config) error {
	return m._rep.GetProvider().SaveConfig(conf)
}

//可通过外部添加
func (m *messageService) RegisterNotifyItem(key string, item *notify.NotifyItem) {
	notify.RegisterNotifyItem(key, item)
}

func (m *messageService) GetAllNotifyItem() []notify.NotifyItem {
	return m._rep.NotifyManager().GetAllNotifyItem()
}

// 保存通知项设置
func (m *messageService) SaveNotifyItem(item *notify.NotifyItem) error {
	return m._rep.NotifyManager().SaveNotifyItem(item)
}

//todo: 考虑弄一个,确定后再发送.这样可以先在系统,然后才发送
// 发送站内通知信息,
// toRole: 为-1时发送给所有用户
// sendNow: 是否马上发送
func (ms *messageService) SendSiteNotifyMessage(senderId int32, toRole int,
	msg *notify.SiteMessage, sendNow bool) error {
	v := &mss.Message{
		Id: 0,
		// 消息类型
		Type: notify.TypeSiteMessage,
		// 消息用途
		UseFor: mss.UseForNotify,
		// 发送人角色
		SenderRole: mss.RoleSystem,
		// 发送人编号
		SenderId: senderId,
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收
		AllUser: -1,
		// 是否只能阅读
		Readonly: 1,
	}

	if toRole <= 0 {
		v.AllUser = 1
	} else {
		v.ToRole = toRole
	}
	var err error
	m := ms._rep.MessageManager().CreateMessage(v, msg)
	if _, err = m.Save(); err == nil {
		err = m.Send(nil)
	}
	return err
}

// 对会用户发送站内信
func (ms *messageService) SendSiteMessageToUser(senderId int32, toRole int, toUser int64,
	msg *notify.SiteMessage, sendNow bool) error {
	v := &mss.Message{
		Id: 0,
		// 消息类型
		Type: notify.TypeSiteMessage,
		// 消息用途
		UseFor: mss.UseForNotify,
		// 发送人角色
		SenderRole: mss.RoleSystem,
		// 发送人编号
		SenderId: senderId,
		To: []mss.User{
			{Id: int32(toUser), Role: toRole},
		},
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收
		AllUser: -1,
		// 是否只能阅读
		Readonly: 1,
	}

	var err error
	m := ms._rep.MessageManager().CreateMessage(v, msg)
	if _, err = m.Save(); err == nil {
		err = m.Send(nil)
	}
	return err
}

// 获取站内信
func (m *messageService) GetSiteMessage(id, toUserId int32, toRole int) *dto.SiteMessage {
	msg := m._rep.MessageManager().GetMessage(id)
	if msg != nil && msg.CheckPerm(toUserId, toRole) {
		val := msg.GetValue()
		dto := &dto.SiteMessage{
			Id:           val.Id,
			Type:         val.Type,
			UseFor:       val.UseFor,
			SenderUserId: 0,
			SenderName:   "系统",
			Readonly:     val.Readonly,
			CreateTime:   val.CreateTime,
			ToId:         toUserId,
			ToRole:       toRole,
		}

		switch msg.Type() {
		case notify.TypePhoneMessage:
			dto.Data = msg.(mss.IPhoneMessage).Value()
		case notify.TypeEmailMessage:
			dto.Data = msg.(mss.IMailMessage).Value()
		case notify.TypeSiteMessage:
			dto.Data = msg.(mss.ISiteMessage).Value()
		}

		if msg.SpecialTo() {
			if to := msg.GetTo(toUserId, toRole); to != nil {
				dto.HasRead = to.HasRead
				dto.ReadTime = to.ReadTime
			}
		}
		return dto
	}
	return nil
}

// 获取聊天会话编号
func (m *messageService) GetChatSessionId(senderRole int, senderId int32, toRole int, toId int32) int32 {
	return m._rep.MessageManager().GetChatSessionId(senderRole, senderId, toRole, toId)
}

// 创建聊天会话
func (m *messageService) CreateChatSession(senderRole int, senderId int32, toRole int, toId int32) (mss.Message, error) {
	return m._rep.MessageManager().CreateChatSession(senderRole, senderId, toRole, toId)
}
