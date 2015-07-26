/**
 * Copyright 2015 @ S1N1 Team.
 * name : msg_manager
 * author : jarryliu
 * date : 2015-07-26 21:30
 * description :
 * history :
 */
package mss

// Message send manager
type IMssManager interface {
	// 发送消息
	Send(IMsgTemplate, MsgData) error

	// 获取邮箱模板
	GetMailTemplate(int) *MailTemplate

	// 保存邮箱模版
	SaveMailTemplate(*MailTemplate) error

	// 获取所有的邮箱模版
	GetMailTemplates() []*MailTemplate
}
