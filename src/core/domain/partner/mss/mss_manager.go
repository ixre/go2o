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
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/partner/mss"
)

var _ mss.IMssManager = new(MssManager)

type MssManager struct {
	_partner       partner.IPartner
	_mssRep        mss.IMssRep
	_mailTemplates []*mss.MailTemplate
}

func NewMssManager(p partner.IPartner, rep mss.IMssRep) mss.IMssManager {
	return &MssManager{
		_partner: p,
		_mssRep:  rep,
	}
}

// 发送消息
func (this *MssManager) Send(tpl mss.IMsgTemplate, data mss.MsgData) error {
	//todo:
	return nil
}

// 获取邮箱模板
func (this *MssManager) GetMailTemplate(id int) *mss.MailTemplate {
	return this._mssRep.GetMailTemplate(this._partner.GetAggregateRootId(), id)
}

// 保存邮箱模版
func (this *MssManager) SaveMailTemplate(v *mss.MailTemplate) error {
	v.PartnerId = this._partner.GetAggregateRootId()
	_, err := this._mssRep.SaveMailTemplate(v)
	return err
}

// 获取所有的邮箱模版
func (this *MssManager) GetMailTemplates() []*mss.MailTemplate {
	if this._mailTemplates == nil {
		this._mailTemplates = this._mssRep.GetMailTemplates(this._partner.GetAggregateRootId())
	}
	return this._mailTemplates
}
