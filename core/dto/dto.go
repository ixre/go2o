/**
 * Copyright 2015 @ z3q.net.
 * name : message_result
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

type (
	//操作消息结果
	MessageResult struct {
		Result  bool   `json:"result"`
		Message string `json:"message"`
		Tag     int    `json:"tag"`
	}

	// 站内信
	SiteMessage struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 消息类型
		Type int `db:"msg_type"`
		// 消息用途
		UseFor       int `db:"use_for"`
		SenderUserId int
		SenderName   string
		// 是否只能阅读
		Readonly int `db:"read_only"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
		// 数据
		Data interface{}
		// 接收者编号
		ToId int `db:"to_id"`
		// 接收者角色
		ToRole int `db:"to_role"`
		// 是否阅读
		HasRead int `db:"has_read"`
		// 阅读时间
		ReadTime int64 `db:"read_time"`
	}
)
