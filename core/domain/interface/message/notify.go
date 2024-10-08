/**
 * Copyright 2015 @ 56x.net.
 * name : notify
 * author : jarryliu
 * date : 2016-07-06 18:36
 * description :
 * history :
 */
package mss

import (
	"errors"

	"github.com/ixre/go2o/core/infrastructure/fw"
)

const (
	// 站内信
	TypeSiteMessage int = 1
	// 短信
	TypeSMS int = 2
	// 邮件
	TypeEmail int = 3
)

const (
	// 系统模板
	TplFlagSystem int = 1
)

var (
	// 类型字典
	NotifyTypeMap = map[int]string{
		TypeSiteMessage: "站内信",
		TypeEmail:       "邮件",
		TypeSMS:         "短信",
	}

	// 类型顺序
	NotifyTypeIndex = []int{
		TypeSiteMessage,
		TypeEmail,
		TypeSMS,
	}

	//todo: 已过期
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
			NotifyBy:   TypeSMS,
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
			NotifyBy:   TypeEmail,
			Content:    "您好,本次{operation}验证码为{code},有效期为{minutes}分钟。",
			Tags: map[string]string{
				"operation": "操作,如找回密码,重置手机等",
				"code":      "验证码",
				"minutes":   "有效时间",
			},
		},
	}
)

// 可通过外部添加
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
		// 保存通知模板
		SaveNotifyTemplate(tpl *NotifyTemplate) error
		// 获取所有的通知项
		GetAllNotifyItem() []NotifyItem
		// 获取通知项配置
		GetNotifyItem(key string) NotifyItem
		// 保存通知项设置
		SaveNotifyItem(item *NotifyItem) error
		// 保存短信API
		SaveSmsApiPerm(s *SmsApiPerm) error
		// 获取短信API信息
		GetSmsApiPerm(provider int) *SmsApiPerm
		// 发送手机短信
		SendPhoneMessage(phone string, msg PhoneMessage, data []string, templateId string) error
		// 发送邮件
		SendEmail(to string, msg *MailMessage, data []string, templateId string) error
	}

	INotifyRepo interface {
		// 获取通知模板仓储
		TemplateRepo() fw.Repository[NotifyTemplate]
		// 获取所有的通知项
		GetAllNotifyItem() []NotifyItem

		// 获取通知项
		GetNotifyItem(key string) *NotifyItem

		// 保存通知项
		SaveNotifyItem(v *NotifyItem) error

		// 保存或新增通知模板
		SaveNotifyTemplate(t *NotifyTemplate) (*NotifyTemplate, error)
		// GetNotifyTemplate Get 系统通知模板
		GetAllNotifyTemplate() []*NotifyTemplate
		// DeleteNotifyTemplate Delete 系统通知模板
		DeleteNotifyTemplate(primary interface{}) error
	}

	// 通知项
	NotifyItem struct {
		Key string
		// 发送方式
		NotifyBy int
		// 不允许修改发送方式
		ReadonlyBy bool
		// 模板编号
		TplId int
		// 内容
		Content string
		// 模板包含的标签
		Tags map[string]string
	}

	// 通知项集合
	NotifyItemSet []*NotifyItem
)

// NotifyTemplate 系统通知模板
type NotifyTemplate struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 模板编号
	TplCode string `json:"tplCode" db:"tpl_code" gorm:"column:tpl_code" bson:"tplCode"`
	// 模板类型,1:站内信 2:短信 3:邮件
	TplType int `json:"tplType" db:"tpl_type" gorm:"column:tpl_type" bson:"tplType"`
	// 模板标志,1:系统
	TplFlag int `json:"tplFlag" db:"tpl_flag" gorm:"column:tpl_flag" bson:"tplFlag"`
	// 模板名称
	TplName string `json:"tplName" db:"tpl_name" gorm:"column:tpl_name" bson:"tplName"`
	// 模板内容
	Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
	// 模板标签, 多个用,隔开
	Labels string `json:"labels" db:"labels" gorm:"column:labels" bson:"labels"`
	// 短信服务商代码
	SpCode string `json:"spCode" db:"sp_code" gorm:"column:sp_code" bson:"spCode"`
	// 短信服务商模板编号
	SpTid string `json:"spTid" db:"sp_tid" gorm:"column:sp_tid" bson:"spTid"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// UpdateTime
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
	// 是否删除,0:否 1:是
	IsDeleted int `json:"isDeleted" db:"is_deleted" gorm:"column:is_deleted" bson:"isDeleted"`
}

func (s NotifyTemplate) TableName() string {
	return "sys_notify_template"
}
