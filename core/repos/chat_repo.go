package repos

import (
	"github.com/ixre/go2o/core/domain/interface/chat"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var _ chat.IChatRepository = new(chatRepoImpl)

type chatRepoImpl struct {
	o        fw.ORM
	convRepo chat.IChatConversationRepo
	msgRepo  chat.IChatMsgRepo
}

func NewChatRepo(o fw.ORM) chat.IChatRepository {
	return &chatRepoImpl{
		o: o,
	}
}

// Conversation implements chat.IChatRepo.
func (c *chatRepoImpl) Conversation() chat.IChatConversationRepo {
	if c.convRepo == nil {
		c.convRepo = newChatConversationRepo(c.o)
	}
	return c.convRepo
}

// Msg implements chat.IChatRepo.
func (c *chatRepoImpl) Msg() chat.IChatMsgRepo {
	if c.msgRepo == nil {
		c.msgRepo = newChatMsgRepo(c.o)
	}
	return c.msgRepo
}

var _ chat.IChatConversationRepo = new(chatConversationRepoImpl)

type chatConversationRepoImpl struct {
	fw.BaseRepository[chat.ChatConversation]
}

// newChatConversationRepo 创建聊天会话仓储
func newChatConversationRepo(o fw.ORM) chat.IChatConversationRepo {
	r := &chatConversationRepoImpl{}
	r.ORM = o
	return r
}

var _ chat.IChatMsgRepo = new(chatMsgRepoImpl)

type chatMsgRepoImpl struct {
	fw.BaseRepository[chat.ChatMsg]
}

// newChatMsgRepo 创建消息消息仓储
func newChatMsgRepo(o fw.ORM) chat.IChatMsgRepo {
	r := &chatMsgRepoImpl{}
	r.ORM = o
	return r
}
