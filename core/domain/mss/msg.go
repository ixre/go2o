/**
 * Copyright 2015 @ z3q.net.
 * name : mail_template
 * author : jarryliu
 * date : 2015-07-27 09:19
 * description :
 * history :
 */
package mss

import (
	"go2o/core/domain/interface/mss"
	"regexp"
	"time"
	"strconv"
)

var reg = regexp.MustCompile("\\{([^\\}]+)\\}")

// 翻译标签
func Transplate(c string, m map[string]string) string {
	return reg.ReplaceAllStringFunc(c, func(k string) string {
		key := k[1 : len(k)-1]
		if v, ok := m[key]; ok {
			return v
		}
		return k
	})
}

var _ mss.IMessage = new(messageImpl)

type messageImpl struct {
	_rep  mss.IMssRep
	_msg  *mss.Message
	_tpl  *mss.MailTemplate
	_data mss.MessageData
}

func newMessage(msg *mss.Message, rep mss.IMssRep) mss.IMessage {
	return &messageImpl{
		_rep: rep,
		_msg: msg,
	}
}

// 获取领域编号
func (this *messageImpl) GetDomainId() int {
	return this._msg.Id
}

func (this *messageImpl) Type() int {
	return this._msg.Type
}

// 保存
func (this *messageImpl) Save() (int, error) {
	if this.GetDomainId() > 0 {
		return this._msg.Id, mss.ErrMessageUpdate
	}
	// 检查消息用途,SenderRole不做检查
	if this._msg.UseFor != mss.UseForNotify &&
		this._msg.UseFor != mss.UseForService &&
		this._msg.UseFor != mss.UserForChat {
		return this.GetDomainId(), mss.ErrUnknownMessageUseFor
	}

	// 检查发送目标群体
	if this._msg.AllUser == 1 {
		if this._msg.ToRole > 0 ||
			(this._msg.To != nil && len(this._msg.To) > 0) {
			return 0, mss.ErrMessageAllUser
		}
	} else if this._msg.ToRole > 0 {
		//检验用户类型
		if this._msg.ToRole != mss.RoleMember &&
			this._msg.ToRole != mss.RoleMerchant &&
			this._msg.ToRole != mss.RoleSystem {
			return 0, mss.ErrUnknownRole
		}
		if len(this._msg.To) > 0 {
			return 0, mss.ErrMessageToRole
		}

	} else if len(this._msg.To) == 0 {
		return 0, mss.ErrNoSuchReceiveUser
	}

	id, err := this._rep.SaveMessage(this._msg)
	this._msg.Id = id
	return id, err
}

// 发送
func (this *messageImpl) Send(d mss.MessageData) error {
	if this.GetDomainId() <= 0 {
		return mss.ErrMessageNotSave
	}

	//todo: 检查是否已经发送
	return nil
}

var _ mss.IMailMessage = new(mailMessageImpl)
var _ mss.IMessage = new(mailMessageImpl)

type mailMessageImpl struct {
	*messageImpl
	_val *mss.ValueMailMessage
	_rep mss.IMssRep
}

func newMailMessage(m *messageImpl, v *mss.ValueMailMessage,
	rep mss.IMssRep) mss.IMessage {
	return &mailMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (this *mailMessageImpl) Value() *mss.ValueMailMessage {
	return this._val
}

func (this *mailMessageImpl) Save() (int, error) {
	return this.messageImpl.Save()
}

// 发送
func (this *mailMessageImpl) Send(d mss.MessageData) error {
	if err := this.messageImpl.Send(d);err != nil{
		return err
	}
	v := this._val
	v.Body = Transplate(v.Body, d)
	v.Subject = Transplate(v.Subject, d)

	unix := time.Now().Unix()
	for _, t := range this._msg.To {
		task := &mss.MailTask{
			MerchantId: 0,
			Subject:  v.Subject,
			Body:  v.Body,
			SendTo:    strconv.Itoa(t.Id),  //todo: mail address
			CreateTime: unix,
		}
		this._rep.JoinMailTaskToQueen(task)
	}
	return this.messageImpl.Send(d)
}

var _ mss.IPhoneMessage = new(phoneMessageImpl)
var _ mss.IMessage = new(phoneMessageImpl)

type phoneMessageImpl struct {
	*messageImpl
	_val *mss.ValuePhoneMessage
	_rep mss.IMssRep
}

func newPhoneMessage(m *messageImpl, v *mss.ValuePhoneMessage,
	rep mss.IMssRep) mss.IMessage {
	return &phoneMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (this *phoneMessageImpl) Value() *mss.ValuePhoneMessage {
	return this._val
}

func (this *phoneMessageImpl) Save() (int, error) {
	return this.messageImpl.Save()
}

// 发送
func (this *phoneMessageImpl) Send(d mss.MessageData) error {
	if err := this.messageImpl.Send(d);err != nil{
		return err
	}
	v := *this._val
	v = mss.ValuePhoneMessage(Transplate(string(v), d))
	return this.messageImpl.Send(d)
}

var _ mss.ISiteMessage = new(siteMessageImpl)
var _ mss.IMessage = new(siteMessageImpl)

type siteMessageImpl struct {
	*messageImpl
	_val *mss.ValueSiteMessage
	_rep mss.IMssRep
}

func newSiteMessage(m *messageImpl, v *mss.ValueSiteMessage,
	rep mss.IMssRep) mss.IMessage {
	return &siteMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (this *siteMessageImpl) Value() *mss.ValueSiteMessage {
	return this._val
}

func (this *siteMessageImpl) Save() (int, error) {
	return this.messageImpl.Save()
}

// 发送
func (this *siteMessageImpl) Send(d mss.MessageData) error {
	if err := this.messageImpl.Send(d);err != nil{
		return err
	}
	v := this._val
	v.Subject = Transplate(v.Subject, d)
	v.Message = Transplate(v.Message, d)

	return nil
}
