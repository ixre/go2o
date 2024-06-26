package chat

import (
	"reflect"

	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type ChatType int

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
		GetConversation(rid int, chatType ChatType) (IConversation, error)
	}

	IConversation interface {
		domain.IDomain
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
		UpdateMsgAttrs(msgId int, attrs map[string]interface{}) error
		// RevertMsg 撤回消息
		RevertMsg(msgId int) error
		// DeleteMsg 删除消息
		DeleteMsg(msgId int) error
	}

	// IChatRepository 聊天仓储
	IChatRepository interface {
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
	Extrta string `json:"extra" db:"extra" gorm:"column:extra" bson:"extra"`
}

// ChatMsg 消息消息
type ChatMsg struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会话编号
	ConvId int `json:"convId" db:"conv_id" gorm:"column:conv_id" bson:"convId"`
	// 消息类型, 1: 文本  2: 图片  3: 表情  4: 文件  5:语音  6:位置  7:语音  8:红包  9:名片
	MsgType int `json:"msgType" db:"msg_type" gorm:"column:msg_type" bson:"msgType"`
	// 是否为发起人的消息, 0:否 1:是
	SidMsg int `json:"sidMsg" db:"sid_msg" gorm:"column:sid_msg" bson:"sidMsg"`
	// 消息内容
	Content string `json:"content" db:"content" gorm:"column:content" bson:"content"`
	// 扩展数据
	Extra string `json:"extra" db:"extra" gorm:"column:extra" bson:"extra"`
	// 是否撤回 0:否 1:是, 撤回的消息对方不可见
	IsRevert int `json:"isRevert" db:"is_revert" gorm:"column:is_revert" bson:"isRevert"`
	// 是否删除, 删除的消息对方可见,自己不可见
	IsDeleted int `json:"isDeleted" db:"is_deleted" gorm:"column:is_deleted" bson:"isDeleted"`
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
