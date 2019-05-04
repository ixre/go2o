/**
 * Copyright 2015 @ z3q.net.
 * name : MssManager
 * author : jarryliu
 * date : 2015-07-26 23:08
 * description :
 * history :
 */
package mss

import (
	"encoding/json"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/tmp"
	"time"
)

var _ mss.IUserMessageManager = new(userMessageManagerImpl)

var _ mss.IMessageManager = new(messageManagerImpl)

type messageManagerImpl struct {
	rep mss.IMssRepo
}

func NewMessageManager(rep mss.IMssRepo) mss.IMessageManager {
	return &messageManagerImpl{
		rep: rep,
	}
}

// 创建消息模版对象
func (mm *messageManagerImpl) CreateMessage(msg *mss.Message,
	content interface{}) mss.IMessage {
	m := newMessage(msg, mm.rep).(*messageImpl)
	if content != nil {
		switch m.Type() {
		case notify.TypeEmailMessage:
			return newMailMessage(m, content.(*notify.MailMessage), m.rep)
		case notify.TypeSiteMessage:
			return newSiteMessage(m, content.(*notify.SiteMessage), m.rep)
		case notify.TypePhoneMessage:
			return newPhoneMessage(m, content.(*notify.PhoneMessage), m.rep)
		}
	} else {
		if m.Type() == notify.TypeEmailMessage ||
			m.Type() == notify.TypeSiteMessage ||
			m.Type() == notify.TypePhoneMessage {
			return m
		}
	}
	panic(mss.ErrNotSupportMessageType)
}

// 创建消息模版对象
func (m *messageManagerImpl) GetMessage(id int32) mss.IMessage {
	if msg := m.rep.GetMessage(id); msg != nil {
		con := m.rep.GetMessageContent(msg.Id)
		switch msg.Type {
		case notify.TypePhoneMessage:
			v := notify.PhoneMessage(con.Data)
			return m.CreateMessage(msg, &v)
		case notify.TypeEmailMessage:
			v := notify.MailMessage{}
			json.Unmarshal([]byte(con.Data), &v)
			return m.CreateMessage(msg, &v)
		case notify.TypeSiteMessage:
			v := notify.SiteMessage{}
			json.Unmarshal([]byte(con.Data), &v)
			return m.CreateMessage(msg, &v)
		}
	}
	return nil
}

// 创建用于会员通知的消息对象
func (m *messageManagerImpl) CreateMemberNotifyMessage(memberId int64, msgType int,
	content interface{}) mss.IMessage {
	msg := &mss.Message{
		Type:       msgType,
		UseFor:     mss.UseForNotify,
		SenderRole: mss.RoleSystem,
		SenderId:   0,
		To: []mss.User{
			{
				Id:   int32(memberId),
				Role: mss.RoleMember,
			},
		},
		CreateTime: time.Now().Unix(),
	}
	return m.CreateMessage(msg, content)
}

// 获取聊天会话编号
func (m *messageManagerImpl) GetChatSessionId(senderRole int,
	senderId int32, toRole int, toId int32) int32 {
	var msgId int32 = 0
	tmp.Db().ExecScalar(`SELECT msg_list.id FROM msg_list INNER JOIN msg_to
        ON msg_to.msg_id = msg_list.id WHERE use_for= $1 AND msg_type= $2 AND sender_role= $3
        AND sender_id= $4 AND to_role= $5 AND to_id= $6`, &msgId, mss.UseForChat,
		notify.TypeSiteMessage, senderRole, senderId, toRole, toId)
	return msgId
}

// 创建聊天会话
func (m *messageManagerImpl) CreateChatSession(senderRole int, senderId int32,
	toRole int, toId int32) (mss.Message, error) {
	msgId := m.GetChatSessionId(senderRole, senderId, toRole, toId)
	if msgId > 0 {
		return m.GetMessage(msgId).GetValue(), nil
	}
	msg := mss.Message{
		Type: notify.TypeSiteMessage,
		// 消息用途
		UseFor: mss.UseForChat,
		// 发送人角色
		SenderRole: senderRole,
		// 发送人编号
		SenderId: senderId,
		// 发送的目标
		To: []mss.User{
			{
				Id:   toId,
				Role: toRole,
			},
		},
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收,1为是,0为否
		AllUser: -1,
		// 是否只能阅读
		Readonly: 0,
		// 创建时间
		CreateTime: time.Now().Unix(),
	}
	im := m.CreateMessage(&msg, notify.SiteMessage{
		Subject: "chat",
		Message: "chat",
	})
	id, err := im.Save()
	msg.Id = id
	return msg, err
}

type userMessageManagerImpl struct {
	_appUserId     int32
	_userRole      int //todo: role
	_mssRepo       mss.IMssRepo
	_mailTemplates []*mss.MailTemplate
	_config        *mss.Config
}

func NewMssManager(appUserId int32, rep mss.IMssRepo) mss.IUserMessageManager {
	return &userMessageManagerImpl{
		_appUserId: appUserId,
		_mssRepo:   rep,
	}
}

// 获取聚合根编号
func (u *userMessageManagerImpl) GetAggregateRootId() int32 {
	return u._appUserId
}

// 获取配置
func (u *userMessageManagerImpl) GetConfig() mss.Config {
	if u._config == nil {
		u._config = u._mssRepo.GetConfig(u._appUserId)
	}
	return *u._config
}

// 保存消息设置
func (u *userMessageManagerImpl) SaveConfig(conf *mss.Config) error {
	err := u._mssRepo.SaveConfig(u._appUserId, conf)
	if err == nil {
		u._config = nil
	}
	return err
}

// 获取邮箱模板
func (u *userMessageManagerImpl) GetMailTemplate(id int32) *mss.MailTemplate {
	return u._mssRepo.GetMailTemplate(u._appUserId, id)
}

// 保存邮箱模版
func (u *userMessageManagerImpl) SaveMailTemplate(v *mss.MailTemplate) (int32, error) {
	v.MerchantId = u._appUserId
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return u._mssRepo.SaveMailTemplate(v)
}

// 删除邮件模板
func (u *userMessageManagerImpl) DeleteMailTemplate(id int32) error {
	//mchId := this._partner.GetAggregateRootId()
	//if this._partnerRepo.CheckKvContainValue(mchId, "kvset", strconv.Itoa(id), "mail") > 0 {
	//	return mss.ErrTemplateUsed
	//}
	return u._mssRepo.DeleteMailTemplate(u._appUserId, id)
}

// 获取所有的邮箱模版
func (u *userMessageManagerImpl) GetMailTemplates() []*mss.MailTemplate {
	if u._mailTemplates == nil {
		u._mailTemplates = u._mssRepo.GetMailTemplates(u._appUserId)
	}
	return u._mailTemplates
}
