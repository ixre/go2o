/**
 * Copyright 2015 @ to2.net.
 * name : msg_template
 * author : jarryliu
 * date : 2015-07-26 21:57
 * description :
 * history :
 */
package mss

import "go2o/core/domain/interface/mss/notify"

//todo: 客服消息
var (
	RoleSystem   = 0
	RoleMember   = 1
	RoleMerchant = 2
)

var (
	// 用于通知
	UseForNotify = 1
	// 用于好友交流
	UseForChat = 2
	// 用于客服
	UseForService = 3
)

var (
	// 站内信用途表
	UseForMap = map[int]string{
		1: "站内信",
		2: "系统公告",
		3: "系统通知",
	}
)

type (
	// 消息数据
	Data map[string]string

	User struct {
		Id   int32
		Role int
	}

	// 消息,优先级为: AllUser ->  ToRole  ->  To
	Message struct {
		// 消息编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 消息类型
		Type int `db:"msg_type"`
		// 消息用途
		UseFor int `db:"use_for"`
		// 发送人角色
		SenderRole int `db:"sender_role"`
		// 发送人编号
		SenderId int32 `db:"sender_id"`
		// 发送的目标
		To []User `db:"-"`
		// 内容
		Content *Content `db:"-"`
		// 发送的用户角色
		ToRole int `db:"to_role"`
		// 全系统接收,1为是,0为否
		AllUser int `db:"all_user"`
		// 是否只能阅读
		Readonly int `db:"read_only"`
		// 创建时间
		CreateTime int64 `db:"create_time"`
	}

	// 消息内容
	Content struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 消息编号
		MsgId int32 `db:"msg_id"`
		// 数据
		Data string `db:"msg_data"`
	}

	// 用户消息绑定
	To struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 接收者编号
		ToId int32 `db:"to_id"`
		// 接收者角色
		ToRole int `db:"to_role"`
		// 消息编号
		MsgId int32 `db:"msg_id"`
		// 内容编号
		ContentId int32 `db:"content_id"`
		// 是否阅读
		HasRead int `db:"has_read"`
		// 阅读时间
		ReadTime int64 `db:"read_time"`
	}

	// 回复
	Replay struct {
		// 编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 关联回复编号
		ReferId int32 `db:"refer_id"`
		// 发送者编号
		SenderId int32 `db:"sender_id"`
		// 发送者角色
		SenderRole int `db:"sender_role"`
		// 内容
		Content string `db:"content"`
	}

	IMessage interface {
		// 应用数据
		//Parse(MessageData) string
		// 加入到发送对列
		//JoinQueen(to []string) error

		// 获取领域编号
		GetDomainId() int32

		// 消息类型
		Type() int

		// 检测是否有权限查看
		CheckPerm(toUserId int32, toRole int) bool

		// 是否向特定的人发送
		SpecialTo() bool

		// 获取消息
		GetValue() Message

		// 获取消息发送目标
		GetTo(toUserId int32, toRole int) *To

		// 保存
		Save() (int32, error)

		// 发送
		Send(data Data) error
	}

	ISiteMessage interface {
		Value() *notify.SiteMessage
	}

	IMailMessage interface {
		Value() *notify.MailMessage
	}

	IPhoneMessage interface {
		Value() *notify.PhoneMessage
	}
)
