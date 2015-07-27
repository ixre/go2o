/**
 * Copyright 2015 @ S1N1 Team.
 * name : mss_rep
 * author : jarryliu
 * date : 2015-07-27 08:46
 * description :
 * history :
 */
package mss

type IMssRep interface {
	// 获取邮箱模板
	GetMailTemplate(partnerId, id int) *MailTemplate
	// 保存邮箱模版
	SaveMailTemplate(*MailTemplate) (int, error)
	// 获取所有的邮箱模版
	GetMailTemplates(partnerId int) []*MailTemplate
	// 删除邮件模板
	DeleteMailTemplate(partnerId, id int) error
	// 加入到发送对列
	JoinMailTaskToQueen(*MailTask) error
}
