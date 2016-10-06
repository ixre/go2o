/**
 * Copyright 2015 @ z3q.net.
 * name : mss_test
 * author : jarryliu
 * date : 2016-06-15 08:21
 * description :
 * history :
 */
package mss

import (
	"github.com/jsix/gof"
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/mss"
	"go2o/core/domain/interface/mss/notify"
	_ "go2o/core/testing"
	"testing"
)

var _ mss.IMssRep = new(MssRep)

type MssRep struct {
	_conn        db.Connector
	_globMss     mss.IUserMessageManager
	_notifyItems map[string]*notify.NotifyItem
	_sysManger   mss.IMessageManager
}

func NewMssRep(conn db.Connector) mss.IMssRep {
	return &MssRep{
		_conn: conn,
	}
}

// 系统消息服务
func (this *MssRep) GetManager() mss.IMessageManager {
	if this._sysManger == nil {
		this._sysManger = NewMessageManager(this)
	}
	return this._sysManger
}

func (this *MssRep) GetProvider() mss.IUserMessageManager {
	if this._globMss == nil {
		this._globMss = NewMssManager(0, this)
	}
	return this._globMss
}

// 获取短信配置
func (this *MssRep) GetConfig(userId int) *mss.Config {
	return nil
}

// 保存消息设置
func (this *MssRep) SaveConfig(userId int, conf *mss.Config) error {
	return nil
}

// 获取所有的通知项
func (this *MssRep) GetAllNotifyItem() []notify.NotifyItem {
	return []notify.NotifyItem{}
}

// 获取通知项
func (this *MssRep) GetNotifyItem(key string) *notify.NotifyItem {
	return nil
}

// 保存通知项
func (this *MssRep) SaveNotifyItem(v *notify.NotifyItem) error {
	return nil
}

// 获取邮箱模板
func (this *MssRep) GetMailTemplate(merchantId, id int) *mss.MailTemplate {
	var e mss.MailTemplate
	if err := this._conn.GetOrm().Get(id, &e); err == nil {
		return &e
	}
	return nil
}

// 保存邮箱模版
func (this *MssRep) SaveMailTemplate(v *mss.MailTemplate) (int, error) {
	return v.Id, nil
}

// 获取所有的邮箱模版
func (this *MssRep) GetMailTemplates(merchantId int) []*mss.MailTemplate {
	return []*mss.MailTemplate{}
}

// 删除邮件模板
func (this *MssRep) DeleteMailTemplate(merchantId, id int) error {
	return nil
}

// 加入到发送对列
func (this *MssRep) JoinMailTaskToQueen(v *mss.MailTask) error {
	return nil
}

// 保存消息
func (this *MssRep) SaveMessage(v *mss.Message) (int, error) {
	var err error
	if v.Id > 0 {
		_, _, err = this._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = this._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

// 获取消息
func (this *MssRep) GetMessage(id int) *mss.Message {
	//todo:
	msg := mss.Message{}
	return &msg
}

// 保存用户消息关联
func (this *MssRep) SaveUserMsg(v *mss.To) (int, error) {
	return orm.Save(this._conn.GetOrm(), v, v.Id)
}

// 保存消息内容
func (this *MssRep) SaveMsgContent(v *mss.Content) (int, error) {
	var err error
	if v.Id > 0 {
		_, _, err = this._conn.GetOrm().Save(v.Id, v)
	} else {
		var id int64
		_, id, err = this._conn.GetOrm().Save(nil, v)
		v.Id = int(id)
	}
	return v.Id, err
}

func TestMessageManagerImpl_SendMessage(t *testing.T) {
	mgr := NewMessageManager(NewMssRep(gof.CurrentApp.Db()))
	v := &mss.Message{
		Id: 0,
		// 消息类型
		Type: notify.TypeEmailMessage,
		// 消息用途
		UseFor: mss.UseForNotify,
		// 发送人角色
		SenderRole: mss.RoleSystem,
		// 发送人编号
		SenderId: 0,
		// 发送的目标
		To: []mss.User{
			mss.User{
				Role: mss.RoleMember,
				Id:   1,
			},
		},
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收
		AllUser: -1,
		// 是否只能阅读
		Readonly: 1,
	}
	val := &notify.MailMessage{
		Subject: "邮件",
		Body:    "您好,邮件{Name}",
	}

	msg := mgr.CreateMessage(v, val)
	var data = map[string]string{
		"Name": "GO2O",
	}
	var err error
	if _, err = msg.Save(); err != nil {
		t.Fatal(err)
	}
	if err = msg.Send(data); err != nil {
		t.Fatal(err)
	}

	t.Log("--- mail sending ok")

	v = &mss.Message{
		Id: 0,
		// 消息类型
		Type: notify.TypePhoneMessage,
		// 消息用途
		UseFor: mss.UseForNotify,
		// 发送人角色
		SenderRole: mss.RoleSystem,
		// 发送人编号
		SenderId: 0,
		// 发送的目标
		To: []mss.User{
			mss.User{
				Role: mss.RoleMember,
				Id:   1,
			},
		},
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收
		AllUser: -1,
		// 是否只能阅读
		Readonly: 1,
	}
	pv := mss.PhoneMessage("您好短信{Name}")
	msg = mgr.CreateMessage(v, &pv)
	if msg.GetDomainId() == 0 {
		_, err = msg.Save()
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := msg.Send(data); err != nil {
		t.Fatal(err)
	}

	t.Log("--- phone message sending ok")

	v = &mss.Message{
		Id: 0,
		// 消息类型
		Type: notify.TypeSiteMessage,
		// 消息用途
		UseFor: mss.UseForNotify,
		// 发送人角色
		SenderRole: mss.RoleSystem,
		// 发送人编号
		SenderId: 0,
		// 发送的用户角色
		ToRole: -1,
		// 全系统接收
		AllUser: 1,
		// 是否只能阅读
		Readonly: 1,
	}
	sm := notify.SiteMessage{
		Subject: "站内信",
		Message: "您好短信{Name}",
	}
	msg = mgr.CreateMessage(v, &sm)
	if msg.GetDomainId() == 0 {
		_, err = msg.Save()
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := msg.Send(data); err != nil {
		t.Fatal(err)
	}

	t.Log("--- site message sending ok")

}
