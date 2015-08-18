/**
 * Copyright 2015 @ z3q.net.
 * name : mail_task
 * author : jarryliu
 * date : 2015-07-27 09:36
 * description :
 * history :
 */
package mss

type MailTask struct {
	// 编号
	Id int `db:"id" pk:"yes" auto:"yes"`
	// 任务编号,无任务为0
	TaskId int `db:"task_id"`
	// 商户编号
	PartnerId int `db:"partner_id"`
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
