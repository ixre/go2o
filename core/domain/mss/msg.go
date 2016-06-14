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
	"errors"
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
	_rep        mss.IMssRep
	_msg        *mss.Message
	_tpl        *mss.MailTemplate
	_data       mss.MessageData
}

func newMailTemplate(msg *mss.Message,rep mss.IMssRep) mss.IMessage {
	return &messageImpl{
		_rep:        rep,
		_msg :msg,
	}
}

// 解析消息模板内容
func (this *messageImpl) parse(val interface{},d mss.MessageData)interface{} {
	this._data = d
	switch this._msg.Type{
	case mss.TypeEmailMessage:
		v := val.(*mss.ValueMailMessage)
		v.Body = Transplate(v.Body,d)
		v.Subject = Transplate(v.Subject,d)
	case mss.TypeSiteMessage:
		v := val.(*mss.ValueSiteMessage)
		v.Subject = Transplate(v.Subject,d)
		v.Message = Transplate(v.Message,d)
	case mss.TypePhoneMessage:
		v := val.(*mss.ValuePhoneMessage)
		Transplate(string(*v),d)
	default:
		panic(errors.New("Unkown message type"))
	}
	return val
}

//todo: 修改邮箱信息
// 加入到发送对列
func (this *messageImpl) sendMailMessage(v *mss.ValueMailMessage) error {
	//unix := time.Now().Unix()
	//for _, _ := range this._msg.To {
	//	task := &mss.MailTask{
	//		MerchantId: 0,
	//		Subject:  v.Subject,
	//		Body:  v.Body,
	//		//SendTo:     t.Id,
	//		CreateTime: unix,
	//	}
	//	this._rep.JoinMailTaskToQueen(task)
	//}
	return nil
}

// 获取领域编号
func (this *messageImpl) GetDomainId()int{
	return this._msg.Id
}

// 保存
func (this *messageImpl)Save()(int,error){
	if this.GetDomainId() > 0{
		return this._msg.Id,mss.ErrMessageUpdate
	}
	id,err := this._rep.SaveMessage(this._msg)
	this._msg.Id = id
	return id,err
}

// 发送
func (this *messageImpl) Send(msgContent interface{}, d mss.MessageData) error{
	if this.GetDomainId() <= 0 {
		return mss.ErrMessageNotSave
	}
	switch this._msg.Type{
	case mss.TypeEmailMessage:
		v := msgContent.(*mss.ValueMailMessage)
		v.Body = Transplate(v.Body, d)
		v.Subject = Transplate(v.Subject, d)
		return this.sendMailMessage(v)
	case mss.TypeSiteMessage:
		v := msgContent.(*mss.ValueSiteMessage)
		v.Subject = Transplate(v.Subject, d)
		v.Message = Transplate(v.Message, d)
		return this.sendSiteMessage(v)
	case mss.TypePhoneMessage:
		v := msgContent.(*mss.ValuePhoneMessage)
		*v = mss.ValuePhoneMessage(Transplate(string(*v), d))
		return this.sendPhoneMessage(v)
	}
	return mss.ErrNotSupportMessageType
}

func (this *messageImpl) sendSiteMessage(v *mss.ValueSiteMessage)error{
	//todo:
	return nil
}


func (this *messageImpl) sendPhoneMessage(v *mss.ValuePhoneMessage)error{
	//todo:
	return nil
}


