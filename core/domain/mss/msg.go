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
	"encoding/json"
	"go2o/core/domain/interface/mss"
	"regexp"
	"strconv"
	"time"
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
	_data mss.Data
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
//todo: 会出现保存后不发送的情况
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
	this._msg.CreateTime = time.Now().Unix()
	id, err := this._rep.SaveMessage(this._msg)
	this._msg.Id = id
	return id, err
}

// 发送
func (this *messageImpl) Send(d mss.Data) error {
	if this.GetDomainId() <= 0 {
		return mss.ErrMessageNotSave
	}
	//todo: 检查是否已经发送
	return nil
}

// 保存消息内容
func (this *messageImpl) saveContent(v interface{}) (int, error) {
	content, ok := v.(string)
	if !ok {
		if d, err := json.Marshal(v); err != nil {
			return -1, err
		} else {
			content = string(d)
		}
	}
	co := &mss.Content{
		Id:    0,
		MsgId: this.GetDomainId(),
		Data:  content,
	}
	return this._rep.SaveMsgContent(co)
}

func (this *messageImpl) saveUserMsg(contentId int, read int) (int, error) {
	if len(this._msg.To) > 0 {
		for _, v := range this._msg.To {
			to := &mss.To{
				Id: 0,
				// 接收者编号
				ToId: v.Id,
				// 接收者角色
				ToRole: v.Role,
				// 内容编号
				ContentId: contentId,
				// 是否阅读
				HasRead: read,
				// 阅读时间
				ReadTime: time.Now().Unix(),
			}
			this._rep.SaveUserMsg(to)
		}
	}
	return -1, nil
}

var _ mss.IMailMessage = new(mailMessageImpl)
var _ mss.IMessage = new(mailMessageImpl)

type mailMessageImpl struct {
	*messageImpl
	_val *mss.MailMessage
	_rep mss.IMssRep
}

func newMailMessage(m *messageImpl, v *mss.MailMessage,
	rep mss.IMssRep) mss.IMessage {
	return &mailMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (this *mailMessageImpl) Value() *mss.MailMessage {
	return this._val
}

func (this *mailMessageImpl) Save() (int, error) {
	return this.messageImpl.Save()
}

// 发送
func (this *mailMessageImpl) Send(d mss.Data) error {
	err := this.messageImpl.Send(d)
	if err == nil {
		v := this._val
		v.Body = Transplate(v.Body, d)
		v.Subject = Transplate(v.Subject, d)

		unix := time.Now().Unix()
		for _, t := range this._msg.To {
			task := &mss.MailTask{
				MerchantId: 0,
				Subject:    v.Subject,
				Body:       v.Body,
				SendTo:     strconv.Itoa(t.Id), //todo: mail address
				CreateTime: unix,
			}
			this._rep.JoinMailTaskToQueen(task)
		}

		//var contentId int //内容编号
		//if contentId, err = this.saveContent(v);err == nil{
		//	this.saveUserMsg(contentId,1) //短信默认已读
		//}
	}
	return err
}

var _ mss.IPhoneMessage = new(phoneMessageImpl)
var _ mss.IMessage = new(phoneMessageImpl)

type phoneMessageImpl struct {
	*messageImpl
	_val *mss.PhoneMessage
	_rep mss.IMssRep
}

func newPhoneMessage(m *messageImpl, v *mss.PhoneMessage,
	rep mss.IMssRep) mss.IMessage {
	return &phoneMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (this *phoneMessageImpl) Value() *mss.PhoneMessage {
	return this._val
}

func (this *phoneMessageImpl) Save() (int, error) {
	return this.messageImpl.Save()
}

// 发送
func (this *phoneMessageImpl) Send(d mss.Data) error {
	err := this.messageImpl.Send(d)
	if err == nil {
		v := *this._val
		v = mss.PhoneMessage(Transplate(string(v), d))
		var contentId int //内容编号
		if contentId, err = this.saveContent(string(v)); err == nil {
			this.saveUserMsg(contentId, 1) //短信默认已读
		}
	}
	return err
}

var _ mss.ISiteMessage = new(siteMessageImpl)
var _ mss.IMessage = new(siteMessageImpl)

type siteMessageImpl struct {
	*messageImpl
	_val *mss.SiteMessage
	_rep mss.IMssRep
}

func newSiteMessage(m *messageImpl, v *mss.SiteMessage,
	rep mss.IMssRep) mss.IMessage {
	return &siteMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (this *siteMessageImpl) Value() *mss.SiteMessage {
	return this._val
}

func (this *siteMessageImpl) Save() (int, error) {
	return this.messageImpl.Save()
}

// 发送
func (this *siteMessageImpl) Send(d mss.Data) error {
	err := this.messageImpl.Send(d)
	if err == nil {
		v := this._val
		v.Subject = Transplate(v.Subject, d)
		v.Message = Transplate(v.Message, d)
		var contentId int //内容编号
		if contentId, err = this.saveContent(v); err == nil {
			this.saveUserMsg(contentId, 0) //站内信默认未读
		}
	}
	return err
}
