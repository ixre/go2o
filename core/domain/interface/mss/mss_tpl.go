/**
 * Copyright 2015 @ to2.net.
 * name : mss_tpl
 * author : jarryliu
 * date : 2016-06-14 19:29
 * description :
 * history :
 */
package mss

// 邮件模版
type MailTemplate struct {
	// 编号
	Id int32 `db:"id" pk:"yes" auto:"yes"`
	// 商户编号
	MerchantId int32 `db:"merchant_id"`
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

type MailTask struct {
	// 编号
	Id int32 `db:"id" pk:"yes" auto:"yes"`
	// 任务编号,无任务为0
	TaskId int32 `db:"task_id"`
	// 商户编号
	MerchantId int32 `db:"merchant_id"`
	// 发送至
	SendTo string `db:"send_to"`
	// 主题
	Subject string `db:"subject"`
	// 内容
	Body string `db:"body"`
	// 是否发送(0,1)
	IsSend int `db:"is_send"`
	// 是否失败(0,1)
	IsFailed int `db:"is_failed"`
	// 创建时间
	CreateTime int64 `db:"create_time"`
	// 发送时间
	SendTime int64 `db:"update_time`
}
