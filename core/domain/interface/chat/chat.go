package chat

import (
	"reflect"

	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

// ChatType 聊天类型
type ChatType int

// 消息标志
type MsgFlag int

const (
	// MsgFlagHint 提示
	MsgFlagHint = 1
	// MsgFlagRevert 撤回
	MsgFlagRevert = 2
	// MsgFlagDelete 删除
	MsgFlagDelete = 4
)

var (
	// ChatTypeNormal 用户聊天
	ChatTypeNormal = 0
	// ChatTypeService 客服
	ChatTypeService = 1
)

type (
	IChatUserAggregateRoot interface {
		domain.IAggregateRoot
		// GetConversation 获取聊天会话
		GetConversation(convId int) IConversation
		// BuildConversation 生成聊天会话
		BuildConversation(rid int, chatType ChatType) (IConversation, error)
	}

	IConversation interface {
		domain.IDomain
		// Get 获取值
		Get() ChatConversation
		// GetMsg 获取消息
		GetMsg(msgId int) *ChatMsg
		// Destroy 删除会话
		Destroy() error
		// Greet 打招呼
		Greet(msg string) error
		// Send 发送消息，并返回消息编号
		Send(msg *MsgBody) (int, error)
		// FetchHistoryMsgs 获取历史消息
		FetchHistoryMsgList(lastTime int, size int) []*ChatMsg
		// FetchMsgList 获取最近的消息
		FetchMsgList(lastTime int, size int) []*ChatMsg
		// UpdateMsgAttrs 更新消息扩展数据
		UpdateMsgAttrs(msgId int, attrs map[string]string) error
		// RevertMsg 撤回消息
		RevertMsg(msgId int) error
		// DeleteMsg 删除消息
		DeleteMsg(msgId int) error
	}

	// IChatRepository 聊天仓储
	IChatRepository interface {
		// GetChatUser 获取聊天人聚合
		GetChatUser(sid int) IChatUserAggregateRoot
		// Conversation 获取聊天仓储
		Conversation() IChatConversationRepo
		// Msg 获取消息仓储
		Msg() IChatMsgRepo
	}

	// IChatConversationRepo 聊天会话仓储
	IChatConversationRepo interface {
		fw.Repository[ChatConversation]
	}
	// IChatMsgRepo 消息消息仓储
	IChatMsgRepo interface {
		fw.Repository[ChatMsg]
	}
)

var _ domain.IValueObject = new(ChatMsg)

// ChatConversation 聊天会话
type ChatConversation struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 编码
	Key string `json:"key" db:"key" gorm:"column:key" bson:"key"`
	// 会话发起人
	Sid int `json:"sid" db:"sid" gorm:"column:sid" bson:"sid"`
	// 会话回复人
	Rid int `json:"rid" db:"rid" gorm:"column:rid" bson:"rid"`
	// 预留标志
	Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
	// 聊天类型,1:用户  2:客服
	ChatType int `json:"chatType" db:"chat_type" gorm:"column:chat_type" bson:"chatType"`
	// 打招呼内容
	GreetWord string `json:"greetWord" db:"greet_word" gorm:"column:greet_word" bson:"greetWord"`
	// 最后聊天时间
	LastChatTime int `json:"lastChatTime" db:"last_chat_time" gorm:"column:last_chat_time" bson:"lastChatTime"`
	// LastMsg
	LastMsg string `json:"lastMsg" db:"last_msg" gorm:"column:last_msg" bson:"lastMsg"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (c ChatConversation) TableName() string {
	return "chat_conversation"
}

// MsgBody 消息内容
type MsgBody struct {
	// 消息类型, 1: 文本  2: 图片  3: 表情  4: 文件  5:语音  6:位置  7:语音  8:红包  9:名片
	MsgType int `json:"msgType" db:"msg_type" gorm:"column:msg_type" bson:"msgType"`
	// 消息内容
	Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
	// 扩展数据
	Extra map[string]string `json:"extra" db:"extra" gorm:"column:extra" bson:"extra"`
}

// ChatMsg 消息消息
type ChatMsg struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会话编号
	ConvId int `json:"convId" db:"conv_id" gorm:"column:conv_id" bson:"convId"`
	// 发送人编号
	Sid int `json:"sid" db:"sid" gorm:"column:sid" bson:"sid"`
	// 消息类型, 1: 文本  2: 图片  3: 表情  4: 文件  5:语音  6:位置  7:语音  8:红包  9:名片  11: 委托申请
	MsgType int `json:"msgType" db:"msg_type" gorm:"column:msg_type" bson:"msgType"`
	// 消息标志: 1:撤回 2:删除
	MsgFlag int `json:"msgFlag" db:"msg_flag" gorm:"column:msg_flag" bson:"msgFlag"`
	// 消息内容
	Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
	// 扩展数据
	Extra string `json:"extra" db:"extra" gorm:"column:extra" bson:"extra"`
	// 过期时间
	ExpiresTime int `json:"expiresTime" db:"expires_time" gorm:"column:expires_time" bson:"expiresTime"`
	// 消息清理时间,0表示永不清理
	PurgeTime int `json:"purgeTime" db:"purge_time" gorm:"column:purge_time" bson:"purgeTime"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

// Equal implements domain.IValueObject.
func (c *ChatMsg) Equal(v interface{}) bool {
	return reflect.DeepEqual(c, v)
}

func (c ChatMsg) TableName() string {
	return "chat_msg"
}
