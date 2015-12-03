/**
 * Copyright 2015 @ z3q.net.
 * name : mail_template
 * author : jarryliu
 * date : 2015-07-26 21:31
 * description :
 * history :
 */
package mss

// 邮件模版
type MailTemplate struct {
	// 编号
	Id int `db:"id" pk:"yes" auto:"yes"`
	// 商户编号
	PartnerId int `db:"partner_id"`
	// 名称
	Name string `db:"name"`
	// 主题
	Subject string `db:"subject"`
	// 内容
	Body string `db:"body"`

	// 是否启用
	Enabled int `db:"enabled"`

	// 创建时间
	CreateTime int64 `db:"create_time"`
	// 更新时间
	UpdateTime int64 `db:"update_time"`
}
