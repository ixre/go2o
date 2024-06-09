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

	mss "github.com/ixre/go2o/core/domain/interface/message"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/util/collections"
	"github.com/ixre/go2o/core/infrastructure/util/sms"
	"github.com/ixre/gof/domain/eventbus"
	"github.com/ixre/gof/log"
)

var _ mss.INotifyManager = new(notifyManagerImpl)

type notifyManagerImpl struct {
	repo         mss.INotifyRepo
	registryRepo registry.IRegistryRepo
	mssRepo      mss.IMessageRepo
	valueRepo    valueobject.IValueRepo
}

func NewNotifyManager(repo mss.INotifyRepo,
	mssRepo mss.IMessageRepo,
	registryRepo registry.IRegistryRepo) mss.INotifyManager {
	return &notifyManagerImpl{
		repo:         repo,
		mssRepo:      mssRepo,
		registryRepo: registryRepo,
	}
}

// 获取所有的通知项
func (n *notifyManagerImpl) GetAllNotifyItem() []mss.NotifyItem {
	return n.repo.GetAllNotifyItem()
}

// 获取通知项配置
func (n *notifyManagerImpl) GetNotifyItem(key string) mss.NotifyItem {
	return *n.repo.GetNotifyItem(key)
}

// 保存通知项设置
func (n *notifyManagerImpl) SaveNotifyItem(item *mss.NotifyItem) error {
	v := n.repo.GetNotifyItem(item.Key)
	if v == nil {
		return mss.ErrNoSuchNotifyItem
	}
	v.Content = item.Content
	v.TplId = item.TplId
	v.NotifyBy = item.NotifyBy
	return n.repo.SaveNotifyItem(v)
}

// 保存短信API
func (n *notifyManagerImpl) SaveSmsApiPerm(v *mss.SmsApiPerm) error {
	if v.Provider == int(mss.CUSTOM) {
		return errors.New("can't setting for custom sms")
	}
	err := sms.CheckSmsApiPerm(v)
	if err == nil {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("sms_api_%d", v.Provider)
		if ir := n.registryRepo.Get(key); ir != nil {
			err = ir.Update(string(data))
			if err == nil {
				return ir.Save()
			}
		}
		// 创建新的键
		data2, _ := json.Marshal(mss.SmsApiPerm{})
		ir := n.registryRepo.Create(&registry.Registry{
			Key:          key,
			Value:        string(data),
			DefaultValue: string(data2),
			Options:      "",
			//Flag:         registry.FlagUserDefine,
			Description: fmt.Sprintf("SMS-API-%d", v.Provider),
		})
		return ir.Save()
	}
	return err
}

// 获取短信API信息
func (n *notifyManagerImpl) GetSmsApiPerm(provider int) *mss.SmsApiPerm {
	key := fmt.Sprintf("sms_api_%d", provider)
	ir := n.registryRepo.Get(key)
	if ir != nil {
		perm := &mss.SmsApiPerm{}
		if err := json.Unmarshal([]byte(ir.StringValue()), perm); err != nil {
			log.Println("[ GO2O][ Sms]: unmarshal api perm failed!", err)
			return nil
		}
		return perm
	}
	return nil
}

// 发送手机短信
func (n *notifyManagerImpl) SendPhoneMessage(phone string, msg mss.PhoneMessage,
	data []string, templateId string) error {
	provider := n.registryRepo.Get(registry.SmsDefaultProvider).IntValue()

	tpl := n.getSmsTemplate(templateId)
	if tpl == nil {
		return fmt.Errorf(mss.ErrNoSuchTemplate.Error(), templateId)
	}
	// 通过外部系统发送短信
	if provider == int(mss.CUSTOM) {
		eventbus.Publish(&events.SendSmsEvent{
			Provider:     provider,
			Phone:        phone,
			Template:     string(msg),
			TemplateCode: templateId,
			SpTemplateId: templateId,
			Data:         data,
		})
		return nil
	}
	// if provider <= 0 {
	// 	return mss.ErrNotSettingSmsProvider
	// }
	setting := n.GetSmsApiPerm(provider)
	if setting == nil {
		//return mss.ErrNotSettingSmsProvider
		setting = &mss.SmsApiPerm{
			Provider: provider,
		}
	}
	return sms.SendSms(setting, phone, string(msg), data)
}

func (n *notifyManagerImpl) getSmsTemplate(templateId string) *mss.NotifyTemplate {
	arr := n.mssRepo.GetAllNotifyTemplate()
	return collections.FindArray(arr, func(t *mss.NotifyTemplate) bool {
		return t.TempType == 2 && t.Code == templateId
	})
}

// 发送邮件
func (n *notifyManagerImpl) SendEmail(to string,
	msg *mss.MailMessage, data []string) error {
	return errors.New("not implement message via mail")
}
