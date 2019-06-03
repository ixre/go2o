/**
 * Copyright 2015 @ to2.net.
 * name : notify
 * author : jarryliu
 * date : 2016-07-06 18:36
 * description :
 * history :
 */
package notify

import (
	"errors"
	"go2o/core/infrastructure/domain"
)

const (
	TypeSiteMessage = 1 + iota
	TypeEmailMessage
	TypePhoneMessage
)

var (
	ErrNoSuchNotifyItem *domain.DomainError = domain.NewError(
		"err_no_such_notify_item", "通知项不存在")

	// 类型字典
	NotifyTypeMap = map[int]string{
		TypeSiteMessage:  "站内信",
		TypeEmailMessage: "邮件",
		TypePhoneMessage: "短信",
	}

	// 类型顺序
	NotifyTypeIndex = []int{
		TypeSiteMessage,
		TypeEmailMessage,
		TypePhoneMessage,
	}

	// 默认通知项
	DefaultNotifyItems = NotifyItemSet{
		&NotifyItem{
			Key:      "注册通知",
			TplId:    -1,
			NotifyBy: TypeSiteMessage,
			Content:  "您好,恭喜您已注册成为会员!",
			Tags:     map[string]string{},
		},
		&NotifyItem{
			Key:        "验证手机",
			TplId:      -1,
			ReadonlyBy: true,
			NotifyBy:   TypePhoneMessage,
			Content:    "您好,本次{operation}验证码为{code},有效期为{minutes}分钟。",
			Tags: map[string]string{
				"operation": "操作,如找回密码,重置手机等",
				"code":      "验证码",
				"minutes":   "有效时间",
			},
		},
		&NotifyItem{
			Key:        "验证邮箱",
			TplId:      -1,
			ReadonlyBy: true,
			NotifyBy:   TypeEmailMessage,
			Content:    "您好,本次{operation}验证码为{code},有效期为{minutes}分钟。",
			Tags: map[string]string{
				"operation": "操作,如找回密码,重置手机等",
				"code":      "验证码",
				"minutes":   "有效时间",
			},
		},
	}
)

//可通过外部添加
func RegisterNotifyItem(key string, item *NotifyItem) {
	for _, v := range DefaultNotifyItems {
		if v.Key == key {
			panic(errors.New("通知项" + key + "已存在!"))
		}
	}
	DefaultNotifyItems = append(DefaultNotifyItems, item)
}

type (
	// 简讯
	PhoneMessage string

	// 邮件消息
	MailMessage struct {
		// 主题
		Subject string `json:"subject"`
		// 内容
		Body string `json:"body"`
	}

	// 站内信
	SiteMessage struct {
		// 主题
		Subject string `json:"subject"`
		// 信息内容
		Message string `json:"message"`
	}

	INotifyManager interface {
		// 获取所有的通知项
		GetAllNotifyItem() []NotifyItem
		// 获取通知项配置
		GetNotifyItem(key string) NotifyItem
		// 保存通知项设置
		SaveNotifyItem(item *NotifyItem) error
		// 发送手机短信
		SendPhoneMessage(phone string, msg PhoneMessage, data map[string]interface{}) error
		// 发送邮件
		SendEmail(to string, msg *MailMessage, data map[string]interface{}) error
	}

	INotifyRepo interface {
		// 获取所有的通知项
		GetAllNotifyItem() []NotifyItem

		// 获取通知项
		GetNotifyItem(key string) *NotifyItem

		// 保存通知项
		SaveNotifyItem(v *NotifyItem) error
	}

	// 通知项
	NotifyItem struct {
		Key string
		// 发送方式
		NotifyBy int
		// 不允许修改发送方式
		ReadonlyBy bool
		TplId      int
		Content    string
		Tags       map[string]string
	}

	// 通知项集合
	NotifyItemSet []*NotifyItem
)
