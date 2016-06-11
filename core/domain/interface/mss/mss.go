/**
 * Copyright 2015 @ z3q.net.
 * name : msg_manager
 * author : jarryliu
 * date : 2015-07-26 21:30
 * description :
 * history :
 */
package mss

// Message manager
type IMessageProvider interface {
	// 获取聚合根编号
	GetAggregateRootId() int

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
