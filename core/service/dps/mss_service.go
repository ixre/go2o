/**
 * Copyright 2015 @ z3q.net.
 * name : mss_service.go
 * author : jarryliu
 * date : 2016-06-11 20:51
 * description :
 * history :
 */
package dps

import (
	"go2o/core/domain/interface/mss"
)

type mssService struct {
	_rep mss.IMssRep
}

func NewMssService(rep mss.IMssRep) *mssService {
	return &mssService{
		_rep: rep,
	}
}

// 获取邮件模版
func (this *mssService) GetMailTemplate(id int) *mss.MailTemplate {
	return this._rep.GetProvider().GetMailTemplate(id)
}

// 保存邮件模板
func (this *mssService) SaveMailTemplate(v *mss.MailTemplate) (int, error) {
	return this._rep.GetProvider().SaveMailTemplate(v)
}

// 获取邮件模板
func (this *mssService) GetMailTemplates() []*mss.MailTemplate {
	return this._rep.GetProvider().GetMailTemplates()
}

// 删除邮件模板
func (this *mssService) DeleteMailTemplate(id int) error {
	return this._rep.GetProvider().DeleteMailTemplate(id)
}

// 获取邮件绑定
func (this *mssService) GetConfig() mss.Config {
	return this._rep.GetProvider().GetConfig()
}

// 保存邮件
func (this *mssService) SaveConfig(conf *mss.Config) error {
	return this._rep.GetProvider().SaveConfig(conf)
}

//可通过外部添加
func (this *mssService) RegisterNotifyItem(key string, item *mss.NotifyItem) {
	mss.RegisterNotifyItem(key, item)
}

func (this *mssService) GetAllNotifyItem() []mss.NotifyItem {
	return this._rep.GetManager().GetAllNotifyItem()
}

// 获取通知项配置
func (this *mssService) GetNotifyItem(key string) mss.NotifyItem {
	return this._rep.GetManager().GetNotifyItem(key)
}

// 保存通知项设置
func (this *mssService) SaveNotifyItem(item *mss.NotifyItem) error {
	return this._rep.GetManager().SaveNotifyItem(item)
}

//todo: 考虑弄一个,确定后再发送.这样可以先在系统,然后才发送
// 发送站内通知信息,
// toRole: 为-1时发送给所有用户
// sendNow: 是否马上发送
func (this *mssService) SendSiteNotifyMessage(senderId int, toRole int,
	msg *mss.SiteMessage, sendNow bool) error {
	v := &mss.Message{
		Id: 0,
		// 消息类型
		Type: mss.TypeSiteMessage,
		// 消息用途
		UseFor: mss.UseForNotify,
		// 发送人角色
		SenderRole: mss.RoleSystem,
		// 发送人编号
		SenderId: senderId,
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收
		AllUser: 1,
		// 是否只能阅读
		Readonly: 1,
	}

	if toRole == -1 {
		v.AllUser = 1
	} else {
		v.ToRole = toRole
	}
	var err error
	m := this._rep.GetManager().CreateMessage(v, msg)
	if _, err = m.Save(); err == nil {
		err = m.Send(nil)
	}
	return err
}
