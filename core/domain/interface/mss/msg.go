/**
 * Copyright 2015 @ z3q.net.
 * name : msg_template
 * author : jarryliu
 * date : 2015-07-26 21:57
 * description :
 * history :
 */
package mss

var (
	RoleSystem   = 0
	RoleMember   = 1
	RoleMerchant = 2
)

var (
	// 用于通知
	UseForNotify = 1
	// 用于好友交流
	UserForChat = 2
	// 用于客服
	UseForService = 3
)

type (

	// 消息数据
	MessageData map[string]string

	// 简讯
	ValuePhoneMessage string

	// 邮件消息
	ValueMailMessage struct {
		// 主题
		Subject string `json:"subject"`
		// 内容
		Body string `json:"body"`
	}

	// 站内信
	ValueSiteMessage struct {
		// 主题
		Subject string `json:"subject"`
		// 信息内容
		Message string `json:"message"`
	}

	User struct {
		Id   int
		Role int
	}

	// 消息,优先级为: AllUser ->  ToRole  ->  To
	Message struct {
		// 消息编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 消息类型
		Type int `db:"msg_type"`
		// 消息用途
		UseFor int `db:"use_for"`
		// 发送人角色
		SenderRole int `db:"sender_role"`
		// 发送人类型
		SenderId int `db:"sender_int"`
		// 发送的目标
		To []User `db:"-"`
		// 发送的用户角色
		ToRole int `db:"to_role"`
		// 全系统接收
		AllUser int `db:"all_user"`
		// 是否只能阅读
		Readonly int `db:"read_only"`
	}

	// 消息内容
	Content struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 消息编号
		MsgId int `db:"msg_id"`
		// 数据
		Data string `db:"msg_data"`
	}

	// 用户消息绑定
	UserMessage struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 接收者编号
		ToId int `db:"to_id"`
		// 接收者角色
		ToRole int `db:"to_role"`
		// 内容编号
		ContentId int `db:"content_id"`
		// 是否阅读
		HasRead int `db:"has_read"`
	}

	// 回复
	Replay struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 关联回复编号
		ReplayId int `db:"from_id"`
		// 发送者编号
		SenderId int `db:"sender_id"`
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
		GetDomainId() int

		// 消息类型
		Type() int

		// 保存
		Save() (int, error)

		// 发送
		Send(data MessageData) error
	}

	ISiteMessage interface {
		Value() *ValueSiteMessage
	}

	IMailMessage interface {
		Value() *ValueMailMessage
	}

	IPhoneMessage interface {
		Value() *ValuePhoneMessage
	}
)
