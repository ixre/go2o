package repos

import (
	"github.com/ixre/go2o/core/domain/interface/chat"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var _ chat.IChatRepo = new(chatRepoImpl)

type chatRepoImpl struct {
	o fw.ORM
}

func NewChatRepo(o fw.ORM) chat.IChatRepo {
	return &chatRepoImpl{
		o: o,
	}
}

// Conversation implements chat.IChatRepo.
func (c *chatRepoImpl) Conversation() chat.IChatConversationRepo {
	panic("unimplemented")
}

// Msg implements chat.IChatRepo.
func (c *chatRepoImpl) Msg() chat.IChatMsgRepo {
	panic("unimplemented")
}

var _ chat.IChatConversationRepo = new(chatConversationRepoImpl)

type chatConversationRepoImpl struct {
	fw.BaseRepository[chat.ChatConversation]
}

// NewChatConversationRepo 创建聊天会话仓储
func NewChatConversationRepo(o fw.ORM) chat.IChatConversationRepo {
	r := &chatConversationRepoImpl{}
	r.ORM = o
	return r
}

var _ chat.IChatMsgRepo = new(chatMsgRepoImpl)

type chatMsgRepoImpl struct {
	fw.BaseRepository[chat.ChatMsg]
}

// NewChatMsgRepo 创建消息消息仓储
func NewChatMsgRepo(o fw.ORM) chat.IChatMsgRepo {
	r := &chatMsgRepoImpl{}
	r.ORM = o
	return r
}
