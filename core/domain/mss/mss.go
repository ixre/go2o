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
	"time"
)

var _ mss.IUserMessageManager = new(userMessageManagerImpl)

var _ mss.IMessageManager = new(messageManagerImpl)

type messageManagerImpl struct {
	_rep mss.IMssRep
}

func NewMessageManager(rep mss.IMssRep) mss.IMessageManager {
	return &messageManagerImpl{
		_rep: rep,
	}
}

// 获取所有的通知项
func (this *messageManagerImpl) GetAllNotifyItem() []mss.NotifyItem {
	return this._rep.GetAllNotifyItem()
}

// 获取通知项配置
func (this *messageManagerImpl) GetNotifyItem(key string) mss.NotifyItem {
	return *this._rep.GetNotifyItem(key)
}

// 保存通知项设置
func (this *messageManagerImpl) SaveNotifyItem(item *mss.NotifyItem) error {
	v := this._rep.GetNotifyItem(item.Key)
	if v == nil {
		return mss.ErrNoSuchNotifyItem
	}
	v.Content = item.Content
	v.TplId = item.TplId
	v.NotifyBy = item.NotifyBy
	return this._rep.SaveNotifyItem(v)
}

// 创建消息模版对象
func (this *messageManagerImpl) CreateMessage(msg *mss.Message,
	content interface{}) mss.IMessage {
	m := newMessage(msg, this._rep).(*messageImpl)
	if content != nil {
		switch m.Type() {
		case mss.TypeEmailMessage:
			return newMailMessage(m, content.(*mss.MailMessage), this._rep)
		case mss.TypeSiteMessage:
			return newSiteMessage(m, content.(*mss.SiteMessage), this._rep)
		case mss.TypePhoneMessage:
			return newPhoneMessage(m, content.(*mss.PhoneMessage), this._rep)
		}
	} else {
		if m.Type() == mss.TypeEmailMessage ||
			m.Type() == mss.TypeSiteMessage ||
			m.Type() == mss.TypePhoneMessage {
			return m
		}
	}
	panic(mss.ErrNotSupportMessageType)
}

// 创建消息模版对象
func (this *messageManagerImpl) GetMessage(id int) mss.IMessage {
	if msg := this._rep.GetMessage(id); msg != nil {
		con := this._rep.GetMessageContent(msg.Id)
		switch msg.Type {
		case mss.TypePhoneMessage:
			v := mss.PhoneMessage(con.Data)
			return this.CreateMessage(msg, &v)
		case mss.TypeEmailMessage:
			v := mss.MailMessage{}
			json.Unmarshal([]byte(con.Data), &v)
			return this.CreateMessage(msg, &v)
		case mss.TypeSiteMessage:
			v := mss.SiteMessage{}
			json.Unmarshal([]byte(con.Data), &v)
			return this.CreateMessage(msg, &v)
		}
	}
	return nil
}

type userMessageManagerImpl struct {
	_appUserId     int
	_userRole      int //todo: role
	_mssRep        mss.IMssRep
	_mailTemplates []*mss.MailTemplate
	_config        *mss.Config
}

func NewMssManager(appUserId int, rep mss.IMssRep) mss.IUserMessageManager {
	return &userMessageManagerImpl{
		_appUserId: appUserId,
		_mssRep:    rep,
	}
}

// 获取聚合根编号
func (this *userMessageManagerImpl) GetAggregateRootId() int {
	return this._appUserId
}

// 获取配置
func (this *userMessageManagerImpl) GetConfig() mss.Config {
	if this._config == nil {
		this._config = this._mssRep.GetConfig(this._appUserId)
	}
	return *this._config
}

// 保存消息设置
func (this *userMessageManagerImpl) SaveConfig(conf *mss.Config) error {
	err := this._mssRep.SaveConfig(this._appUserId, conf)
	if err == nil {
		this._config = nil
	}
	return err
}

// 获取邮箱模板
func (this *userMessageManagerImpl) GetMailTemplate(id int) *mss.MailTemplate {
	return this._mssRep.GetMailTemplate(this._appUserId, id)
}

// 保存邮箱模版
func (this *userMessageManagerImpl) SaveMailTemplate(v *mss.MailTemplate) (
	int, error) {
	v.MerchantId = this._appUserId
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return this._mssRep.SaveMailTemplate(v)
}

// 删除邮件模板
func (this *userMessageManagerImpl) DeleteMailTemplate(id int) error {
	//merchantId := this._partner.GetAggregateRootId()
	//if this._partnerRep.CheckKvContainValue(merchantId, "kvset", strconv.Itoa(id), "mail") > 0 {
	//	return mss.ErrTemplateUsed
	//}
	return this._mssRep.DeleteMailTemplate(this._appUserId, id)
}

// 获取所有的邮箱模版
func (this *userMessageManagerImpl) GetMailTemplates() []*mss.MailTemplate {
	if this._mailTemplates == nil {
		this._mailTemplates = this._mssRep.GetMailTemplates(this._appUserId)
	}
	return this._mailTemplates
}
