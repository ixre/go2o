/**
 * Copyright 2015 @ z3q.net.
 * name : msg_manager
 * author : jarryliu
 * date : 2015-07-26 21:30
 * description :
 * history :
 */
package mss

const (
	NotifyByMessage = 1 + iota
	NotifyByEMail
	NotifyByPhoneMessage
)

var (
	// 类型字典
	NotifyTypeMap = map[int]string{
		NotifyByMessage:      "站内信",
		NotifyByEMail:        "邮件",
		NotifyByPhoneMessage: "短信",
	}

	// 类型顺序
	NotifyTypeIndex = []int{
		NotifyByMessage,
		NotifyByEMail,
		NotifyByPhoneMessage,
	}
)

type (
	// Message manager
	IMessageProvider interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 获取配置
		GetConfig() Config

		// 保存消息设置
		SaveConfig(conf *Config) error

		// 发送消息
		Send(tpl IMsgTemplate, d MsgData, to []string) error

		// 获取邮箱模板
		GetMailTemplate(int) *MailTemplate

		// 保存邮箱模版
		SaveMailTemplate(*MailTemplate) (int, error)

		// 获取所有的邮箱模版
		GetMailTemplates() []*MailTemplate

		// 删除邮件模板
		DeleteMailTemplate(int) error

		// 创建消息模版对象
		CreateMsgTemplate(v interface{}) (IMsgTemplate, error)
	}

	// 系统消息发送配置
	Config struct {
		// 资料完成
		ProfileCompleteNotifyEnabled bool
		// 通知类型
		ProfileCompleteNotifyType int
		// 通知模板选择
		ProfileCompleteNotifyTpl int

		//注册完成
		RegisterNotifyEnabled bool
		RegisterNotifyType    int
		RegisterNotifyTpl     int
	}
)
