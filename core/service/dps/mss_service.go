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
	return this._rep.GetManager().GetMailTemplate(id)
}

// 保存邮件模板
func (this *mssService) SaveMailTemplate(v *mss.MailTemplate) (int, error) {
	return this._rep.GetManager().SaveMailTemplate(v)
}

// 获取邮件模板
func (this *mssService) GetMailTemplates() []*mss.MailTemplate {
	return this._rep.GetManager().GetMailTemplates()
}

// 删除邮件模板
func (this *mssService) DeleteMailTemplate(id int) error {
	return this._rep.GetManager().DeleteMailTemplate(id)
}

// 获取邮件绑定
func (this *mssService) GetConfig() mss.Config {
	return this._rep.GetManager().GetConfig()
}

// 保存邮件
func (this *mssService) SaveConfig(conf *mss.Config) error {
	return this._rep.GetManager().SaveConfig(conf)
}
