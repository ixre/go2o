/**
 * Copyright 2015 @ S1N1 Team.
 * name : MssManager
 * author : jarryliu
 * date : 2015-07-26 23:08
 * description :
 * history :
 */
package mss

import (
	"errors"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/partner/mss"
	"strconv"
	"time"
)

var _ mss.IMssManager = new(MssManager)

type MssManager struct {
	_partner       partner.IPartner
	_mssRep        mss.IMssRep
	_partnerRep    partner.IPartnerRep
	_mailTemplates []*mss.MailTemplate
}

func NewMssManager(p partner.IPartner, rep mss.IMssRep, partnerRep partner.IPartnerRep) mss.IMssManager {
	return &MssManager{
		_partner:    p,
		_mssRep:     rep,
		_partnerRep: partnerRep,
	}
}

// 创建消息模版对象
func (this *MssManager) CreateMsgTemplate(v interface{}) (mss.IMsgTemplate, error) {
	//todo: other message type
	var err error
	partnerId := this._partner.GetAggregateRootId()
	switch v.(type) {
	case *mss.MailTemplate:
		tpl := v.(*mss.MailTemplate)
		if tpl.Enabled == 0 {
			err = mss.ErrNotEnabled
		}
		return newMailTemplate(partnerId, this._mssRep, tpl), err
	}
	return nil, mss.ErrNotSupportMessageType
}

// 发送消息
func (this *MssManager) Send(tpl mss.IMsgTemplate, data mss.MsgData, to []string) error {
	if tpl != nil {
		tpl.ApplyData(data)
		return tpl.JoinQueen(to)
	}
	return errors.New("template is nil")
}

// 获取邮箱模板
func (this *MssManager) GetMailTemplate(id int) *mss.MailTemplate {
	return this._mssRep.GetMailTemplate(this._partner.GetAggregateRootId(), id)
}

// 保存邮箱模版
func (this *MssManager) SaveMailTemplate(v *mss.MailTemplate) (int, error) {
	v.PartnerId = this._partner.GetAggregateRootId()
	v.UpdateTime = time.Now().Unix()
	if v.CreateTime == 0 {
		v.CreateTime = v.UpdateTime
	}
	return this._mssRep.SaveMailTemplate(v)
}

// 删除邮件模板
func (this *MssManager) DeleteMailTemplate(id int) error {
	partnerId := this._partner.GetAggregateRootId()
	if this._partnerRep.CheckKvContainValue(partnerId, "kvset",strconv.Itoa(id), "mail") > 0 {
		return mss.ErrTemplateUsed
	}
	return this._mssRep.DeleteMailTemplate(partnerId, id)
}

// 获取所有的邮箱模版
func (this *MssManager) GetMailTemplates() []*mss.MailTemplate {
	if this._mailTemplates == nil {
		this._mailTemplates = this._mssRep.GetMailTemplates(this._partner.GetAggregateRootId())
	}
	return this._mailTemplates
}
