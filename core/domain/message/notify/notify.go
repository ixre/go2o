/**
 * Copyright 2015 @ 56x.net.
 * name : notify
 * author : jarryliu
 * date : 2016-07-06 18:41
 * description :
 * history :
 */
package notify

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ixre/go2o/core/domain/interface/message/notify"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/tool/sms"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/log"
)

var _ notify.INotifyManager = new(notifyManagerImpl)

type notifyManagerImpl struct {
	repo         notify.INotifyRepo
	registryRepo registry.IRegistryRepo
	valueRepo    valueobject.IValueRepo
}

func NewNotifyManager(repo notify.INotifyRepo, registryRepo registry.IRegistryRepo) notify.INotifyManager {
	return &notifyManagerImpl{
		repo:         repo,
		registryRepo: registryRepo,
	}
}

// 获取所有的通知项
func (n *notifyManagerImpl) GetAllNotifyItem() []notify.NotifyItem {
	return n.repo.GetAllNotifyItem()
}

// 获取通知项配置
func (n *notifyManagerImpl) GetNotifyItem(key string) notify.NotifyItem {
	return *n.repo.GetNotifyItem(key)
}

// 保存通知项设置
func (n *notifyManagerImpl) SaveNotifyItem(item *notify.NotifyItem) error {
	v := n.repo.GetNotifyItem(item.Key)
	if v == nil {
		return notify.ErrNoSuchNotifyItem
	}
	v.Content = item.Content
	v.TplId = item.TplId
	v.NotifyBy = item.NotifyBy
	return n.repo.SaveNotifyItem(v)
}

// 保存短信API
func (n *notifyManagerImpl) SaveSmsApiPerm(v *notify.SmsApiPerm) error {
	err := sms.CheckSmsApiPerm(v)
	if err == nil {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("sms_api_" ,v.Provider)
		if ir := n.registryRepo.Get(key); ir != nil {
			err = ir.Update(string(data))
			if err == nil {
				return ir.Save()
			}
		}
		// 创建新的键
		data2, _ := json.Marshal(notify.SmsApiPerm{})
		ir := n.registryRepo.Create(&registry.Registry{
			Key:          key,
			Value:        string(data),
			DefaultValue: string(data2),
			Options:      "",
			Flag:         registry.FlagUserDefine,
			Description:  fmt.Sprintf("SMS-API(%s)", v.Provider),
		})
		return ir.Save()
	}
	return err
}

// 获取短信API信息
func (n *notifyManagerImpl) GetSmsApiPerm(provider int) *notify.SmsApiPerm {
	key := fmt.Sprintf("sms_api_%d", provider)
	ir := n.registryRepo.Get(key)
	if ir != nil {
		perm := &notify.SmsApiPerm{}
		if err := json.Unmarshal([]byte(ir.StringValue()), perm); err != nil {
			log.Println("[ Go2o][ Sms]: unmarshal api perm failed!", err)
			return nil
		}
		return perm
	}
	return nil
}

// 发送手机短信
func (n *notifyManagerImpl) SendPhoneMessage(phone string, msg notify.PhoneMessage,
	data []string, templateId string) error {
	pushEvent := n.registryRepo.Get(registry.SmsPushSendEvent).BoolValue()
	provider := n.registryRepo.Get(registry.SmsDefaultProvider).IntValue()
	if provider <= 0 {
		return notify.ErrNotSettingSmsProvider
	}
	// 通过外部系统发送短信
	if pushEvent {
		eventbus.Publish(&events.SendSmsEvent{
			Provider:   provider,
			Phone:      phone,
			Template:   string(msg),
			TemplateId: templateId,
			Data:       data,
		})
		return nil
	}
	setting := n.GetSmsApiPerm(provider)
	if setting == nil {
		return notify.ErrNotSettingSmsProvider
	}
	setting.Provider = provider
	return sms.SendSms(setting, phone, string(msg), data)
}

// 发送邮件
func (n *notifyManagerImpl) SendEmail(to string,
	msg *notify.MailMessage, data []string) error {
	return errors.New("not implement message via mail")
}
