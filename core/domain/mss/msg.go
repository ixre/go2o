/**
 * Copyright 2015 @ to2.net.
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
	"go2o/core/domain/interface/mss/notify"
	"regexp"
	"strconv"
	"time"
)

var reg = regexp.MustCompile("\\{([^\\}]+)\\}")

// 翻译标签
func Translate(c string, m map[string]string) string {
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
	rep  mss.IMssRepo
	msg  *mss.Message
	tpl  *mss.MailTemplate
	data mss.Data
}

func newMessage(msg *mss.Message, rep mss.IMssRepo) mss.IMessage {
	return &messageImpl{
		rep: rep,
		msg: msg,
	}
}

// 获取领域编号
func (m *messageImpl) GetDomainId() int32 {
	return m.msg.Id
}

func (m *messageImpl) Type() int {
	return m.msg.Type
}

// 是否向特定的人发送
func (m *messageImpl) SpecialTo() bool {
	return m.msg.To != nil && len(m.msg.To) > 0
}

// 检测是否有权限查看
func (m *messageImpl) CheckPerm(toUserId int32, toRole int) bool {
	if m.msg.AllUser == 1 || m.msg.ToRole == toRole {
		return true
	}
	if m.msg.To != nil {
		for _, v := range m.msg.To {
			if v.Id == toUserId && v.Role == toRole {
				return true
			}
		}
	}
	return false
}

// 获取消息
func (m *messageImpl) GetValue() mss.Message {
	return *m.msg
}

// 获取消息发送目标
func (m *messageImpl) GetTo(toUserId int32, toRole int) *mss.To {
	return m.rep.GetMessageTo(m.GetDomainId(), toUserId, toRole)
}

// 保存
//todo: 会出现保存后不发送的情况
func (m *messageImpl) Save() (int32, error) {
	if m.GetDomainId() > 0 {
		return m.msg.Id, mss.ErrMessageUpdate
	}
	// 检查消息用途,SenderRole不做检查
	if m.msg.UseFor != mss.UseForNotify &&
		m.msg.UseFor != mss.UseForService &&
		m.msg.UseFor != mss.UseForChat {
		return m.GetDomainId(), mss.ErrUnknownMessageUseFor
	}

	// 检查发送目标群体
	if m.msg.AllUser == 1 {
		if m.msg.ToRole > 0 ||
			(m.msg.To != nil && len(m.msg.To) > 0) {
			return 0, mss.ErrMessageAllUser
		}
	} else if m.msg.ToRole > 0 {
		//检验用户类型
		if m.msg.ToRole != mss.RoleMember &&
			m.msg.ToRole != mss.RoleMerchant &&
			m.msg.ToRole != mss.RoleSystem {
			return 0, mss.ErrUnknownRole
		}
		if len(m.msg.To) > 0 {
			return 0, mss.ErrMessageToRole
		}

	} else if len(m.msg.To) == 0 {
		return 0, mss.ErrNoSuchReceiveUser
	}
	m.msg.CreateTime = time.Now().Unix()
	id, err := m.rep.SaveMessage(m.msg)
	m.msg.Id = id
	return id, err
}

// 发送
func (m *messageImpl) Send(d mss.Data) error {
	if m.GetDomainId() <= 0 {
		return mss.ErrMessageNotSave
	}
	//todo: 检查是否已经发送
	return nil
}

// 保存消息内容
func (m *messageImpl) saveContent(v interface{}) (int32, error) {
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
		MsgId: m.GetDomainId(),
		Data:  content,
	}
	return m.rep.SaveMsgContent(co)
}

func (m *messageImpl) saveUserMsg(contentId int32, read int) (int32, error) {
	if len(m.msg.To) > 0 {
		for _, v := range m.msg.To {
			to := &mss.To{
				Id: 0,
				// 接收者编号
				ToId: v.Id,
				// 接收者角色
				ToRole: v.Role,
				// 消息编号
				MsgId: m.GetDomainId(),
				// 内容编号
				ContentId: contentId,
				// 是否阅读
				HasRead: read,
				// 阅读时间
				ReadTime: time.Now().Unix(),
			}
			return m.rep.SaveUserMsg(to)
		}
	}
	return -1, nil
}

var _ mss.IMailMessage = new(mailMessageImpl)
var _ mss.IMessage = new(mailMessageImpl)

type mailMessageImpl struct {
	*messageImpl
	_val *notify.MailMessage
	_rep mss.IMssRepo
}

func newMailMessage(m *messageImpl, v *notify.MailMessage,
	rep mss.IMssRepo) mss.IMessage {
	return &mailMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (m *mailMessageImpl) Value() *notify.MailMessage {
	return m._val
}

func (m *mailMessageImpl) Save() (int32, error) {
	return m.messageImpl.Save()
}

// 发送
func (m *mailMessageImpl) Send(d mss.Data) error {
	err := m.messageImpl.Send(d)
	if err == nil {
		v := m._val
		v.Body = Translate(v.Body, d)
		v.Subject = Translate(v.Subject, d)

		unix := time.Now().Unix()
		for _, t := range m.msg.To {
			task := &mss.MailTask{
				MerchantId: 0,
				Subject:    v.Subject,
				Body:       v.Body,
				SendTo:     strconv.Itoa(int(t.Id)), //todo: mail address
				CreateTime: unix,
			}
			m._rep.JoinMailTaskToQueen(task)
		}

		//var contentId int //内容编号
		//if contentId, err = m.saveContent(v);err == nil{
		//	m.saveUserMsg(contentId,1) //短信默认已读
		//}
	}
	return err
}

var _ mss.IPhoneMessage = new(phoneMessageImpl)
var _ mss.IMessage = new(phoneMessageImpl)

type phoneMessageImpl struct {
	*messageImpl
	_val *notify.PhoneMessage
	_rep mss.IMssRepo
}

func newPhoneMessage(m *messageImpl, v *notify.PhoneMessage,
	rep mss.IMssRepo) mss.IMessage {
	return &phoneMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (p *phoneMessageImpl) Value() *notify.PhoneMessage {
	return p._val
}

func (p *phoneMessageImpl) Save() (int32, error) {
	return p.messageImpl.Save()
}

// 发送
func (p *phoneMessageImpl) Send(d mss.Data) error {
	err := p.messageImpl.Send(d)
	if err == nil {
		v := *p._val
		v = notify.PhoneMessage(Translate(string(v), d))
		var contentId int32 //内容编号
		if contentId, err = p.saveContent(string(v)); err == nil {
			p.saveUserMsg(contentId, 1) //短信默认已读
		}
	}
	return err
}

var _ mss.ISiteMessage = new(siteMessageImpl)
var _ mss.IMessage = new(siteMessageImpl)

type siteMessageImpl struct {
	*messageImpl
	_val *notify.SiteMessage
	_rep mss.IMssRepo
}

func newSiteMessage(m *messageImpl, v *notify.SiteMessage,
	rep mss.IMssRepo) mss.IMessage {
	return &siteMessageImpl{
		messageImpl: m,
		_val:        v,
		_rep:        rep,
	}
}

func (s *siteMessageImpl) Value() *notify.SiteMessage {
	return s._val
}

func (s *siteMessageImpl) Save() (int32, error) {
	return s.messageImpl.Save()
}

// 发送
func (s *siteMessageImpl) Send(d mss.Data) error {
	err := s.messageImpl.Send(d)
	if err == nil {
		v := s._val
		v.Subject = Translate(v.Subject, d)
		v.Message = Translate(v.Message, d)
		var contentId int32 //内容编号
		if contentId, err = s.saveContent(v); err == nil {
			s.saveUserMsg(contentId, 0) //站内信默认未读
		}
	}
	return err
}
