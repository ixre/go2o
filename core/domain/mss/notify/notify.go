/**
 * Copyright 2015 @ z3q.net.
 * name : notify
 * author : jarryliu
 * date : 2016-07-06 18:41
 * description :
 * history :
 */
package notify

import (
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/tool/sms"
)

var _ notify.INotifyManager = new(notifyManagerImpl)

type notifyManagerImpl struct {
	rep      notify.INotifyRep
	valueRep valueobject.IValueRep
}

func NewNotifyManager(rep notify.INotifyRep,
	valueRep valueobject.IValueRep) notify.INotifyManager {
	return &notifyManagerImpl{
		rep:      rep,
		valueRep: valueRep,
	}
}

// 获取所有的通知项
func (n *notifyManagerImpl) GetAllNotifyItem() []notify.NotifyItem {
	return n.rep.GetAllNotifyItem()
}

// 获取通知项配置
func (n *notifyManagerImpl) GetNotifyItem(key string) notify.NotifyItem {
	return *n.rep.GetNotifyItem(key)
}

// 保存通知项设置
func (n *notifyManagerImpl) SaveNotifyItem(item *notify.NotifyItem) error {
	v := n.rep.GetNotifyItem(item.Key)
	if v == nil {
		return notify.ErrNoSuchNotifyItem
	}
	v.Content = item.Content
	v.TplId = item.TplId
	v.NotifyBy = item.NotifyBy
	return n.rep.SaveNotifyItem(v)
}

// 发送手机短信
func (n *notifyManagerImpl) SendPhoneMessage(phone string,
	msg notify.PhoneMessage, data map[string]interface{}) error {
	i, api := n.valueRep.GetDefaultSmsApiPerm()
	return sms.SendSms(i, api.ApiKey, api.ApiSecret, phone,
		api.ApiUrl, api.Encoding, api.SuccessChar, string(msg), data)
}

// 发送邮件
func (n *notifyManagerImpl) SendEmail(to string,
	msg *notify.MailMessage, data map[string]interface{}) error {
	return nil
}
