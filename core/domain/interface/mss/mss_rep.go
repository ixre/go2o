/**
 * Copyright 2015 @ z3q.net.
 * name : mss_rep
 * author : jarryliu
 * date : 2015-07-27 08:46
 * description :
 * history :
 */
package mss

type IMssRep interface {
	GetManager() IMessageProvider
	// 获取邮箱模板
	GetMailTemplate(merchantId, id int) *MailTemplate
	// 保存邮箱模版
	SaveMailTemplate(*MailTemplate) (int, error)
	// 获取所有的邮箱模版
	GetMailTemplates(merchantId int) []*MailTemplate
	// 删除邮件模板
	DeleteMailTemplate(merchantId, id int) error
	// 加入到发送对列
	JoinMailTaskToQueen(*MailTask) error
}
