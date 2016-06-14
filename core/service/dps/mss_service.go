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
	mss.RegisterNotifyItem(key,item)
}

func (this *mssService) GetAllNotifyItem()[]mss.NotifyItem{
	return this._rep.GetManager().GetAllNotifyItem()
}

// 获取通知项配置
func (this *mssService) GetNotifyItem(key string) mss.NotifyItem{
	return this._rep.GetManager().GetNotifyItem(key)
}
// 保存通知项设置
func (this *mssService) SaveNotifyItem(item *mss.NotifyItem) error{
	return this._rep.GetManager().SaveNotifyItem(item)
}