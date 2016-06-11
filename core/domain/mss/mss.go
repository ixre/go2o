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
	"errors"
	"go2o/core/domain/interface/mss"
	"time"
)

var _ mss.IMessageProvider = new(messageProviderImpl)

type messageProviderImpl struct {
	_appUserId     int
	_mssRep        mss.IMssRep
	_mailTemplates []*mss.MailTemplate
}

func NewMssManager(appUserId int, rep mss.IMssRep) mss.IMessageProvider {
	return &messageProviderImpl{
		_appUserId: appUserId,
		_mssRep:    rep,
	}
}

// 获取聚合根编号
func (this *messageProviderImpl) GetAggregateRootId() int {
	return this._appUserId
}

// 创建消息模版对象
func (this *messageProviderImpl) CreateMsgTemplate(v interface{}) (
	mss.IMsgTemplate, error) {
	//todo: other message type
	var err error
	switch v.(type) {
	case *mss.MailTemplate:
		tpl := v.(*mss.MailTemplate)
		if tpl.Enabled == 0 {
			err = mss.ErrNotEnabled
		}
		return newMailTemplate(this._appUserId, this._mssRep, tpl), err
	}
	return nil, mss.ErrNotSupportMessageType
}

// 发送消息
func (this *messageProviderImpl) Send(tpl mss.IMsgTemplate,
	data mss.MsgData, to []string) error {
	if tpl != nil {
		tpl.ApplyData(data)
		return tpl.JoinQueen(to)
	}
	return errors.New("template is nil")
}

// 获取邮箱模板
func (this *messageProviderImpl) GetMailTemplate(id int) *mss.MailTemplate {
	return this._mssRep.GetMailTemplate(this._appUserId, id)
}

// 保存邮箱模版
func (this *messageProviderImpl) SaveMailTemplate(v *mss.MailTemplate) (
	int, error) {
	v.MerchantId = this._appUserId
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return this._mssRep.SaveMailTemplate(v)
}

// 删除邮件模板
func (this *messageProviderImpl) DeleteMailTemplate(id int) error {
	//merchantId := this._partner.GetAggregateRootId()
	//if this._partnerRep.CheckKvContainValue(merchantId, "kvset", strconv.Itoa(id), "mail") > 0 {
	//	return mss.ErrTemplateUsed
	//}
	return this._mssRep.DeleteMailTemplate(this._appUserId, id)
}

// 获取所有的邮箱模版
func (this *messageProviderImpl) GetMailTemplates() []*mss.MailTemplate {
	if this._mailTemplates == nil {
		this._mailTemplates = this._mssRep.GetMailTemplates(this._appUserId)
	}
	return this._mailTemplates
}
