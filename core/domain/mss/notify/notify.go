/**
 * Copyright 2015 @ to2.net.
 * name : notify
 * author : jarryliu
 * date : 2016-07-06 18:41
 * description :
 * history :
 */
package notify

import (
	"encoding/json"
	"fmt"
	"github.com/ixre/gof/log"
	"go2o/core/domain/interface/mss/notify"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/tool/sms"
)

var _ notify.INotifyManager = new(notifyManagerImpl)

type notifyManagerImpl struct {
	repo         notify.INotifyRepo
	registryRepo registry.IRegistryRepo
	valueRepo    valueobject.IValueRepo
}

func NewNotifyManager(repo notify.INotifyRepo,registryRepo registry.IRegistryRepo) notify.INotifyManager {
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
func (n *notifyManagerImpl) SaveSmsApiPerm(provider string, v *notify.SmsApiPerm) error {
	err := sms.CheckSmsApiPerm(provider, v)
	if err == nil {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		key := "sms_api_" + provider
		if ir := n.registryRepo.Get(key); ir != nil {
			return ir.Update(string(data))
		}
		// 创建新的键
		data2, _ := json.Marshal(notify.SmsApiPerm{})
		ir := n.registryRepo.Create(&registry.Registry{
			Key:          key,
			Value:        string(data),
			DefaultValue: string(data2),
			Options:      "",
			Flag:         registry.FlagUserDefine,
			Description:  fmt.Sprintf("SMS-API(%s)", provider),
		})
		return ir.Save()
	}
	return err
}

// 获取短信API信息
func (n *notifyManagerImpl) GetSmsApiPerm(provider string) *notify.SmsApiPerm {
	key := "sms_api_" + provider
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
func (n *notifyManagerImpl) SendPhoneMessage(phone string,
	msg notify.PhoneMessage, data map[string]interface{}) error {
	provider := n.registryRepo.Get(registry.SmsDefaultProvider).StringValue()
	if provider == "" {
		return notify.ErrNotSettingSmsProvider
	}
	api := n.GetSmsApiPerm(provider)
	if api == nil {
		return notify.ErrNoSuchSmsProvider
	}
	return sms.SendSms(provider, api.ApiKey, api.ApiSecret, phone,
		api.ApiUrl, api.Encoding, api.SuccessChar, string(msg), data)
}

// 发送邮件
func (n *notifyManagerImpl) SendEmail(to string,
	msg *notify.MailMessage, data map[string]interface{}) error {
	return nil
}
