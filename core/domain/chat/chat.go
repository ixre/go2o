package chat

import (
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain/interface/chat"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

var _ chat.IChatUserAggregateRoot = new(chatUserAggregateRoot)

type chatUserAggregateRoot struct {
	sid  int
	repo chat.IChatRepository
}

func CreateChatUser(sid int, repo chat.IChatRepository) chat.IChatUserAggregateRoot {
	return &chatUserAggregateRoot{
		sid:  sid,
		repo: repo,
	}
}

// GetAggregateRootId implements chat.IChatUserAggregateRoot.
func (c *chatUserAggregateRoot) GetAggregateRootId() int {
	return c.sid
}

// GetConversation implements chat.IChatUserAggregateRoot.
func (c *chatUserAggregateRoot) GetConversation(rid int, chatType chat.ChatType) (chat.IConversation, error) {
	v := &chat.ChatConversation{
		Sid:      c.GetAggregateRootId(),
		Rid:      rid,
		ChatType: int(chatType),
	}
	uid := []int{rid, c.GetAggregateRootId()}
	cr := c.repo.Conversation()
	conv := cr.FindBy("sid IN (?) AND rid IN(?)", uid, uid)
	if conv == nil {
		_, err := cr.Save(v)
		if err != nil {
			return nil, err
		}
		conv = v
	}
	return newConversation(conv, c, c.repo), nil
}

var _ chat.IConversation = new(converstationImpl)

// 会话实现
type converstationImpl struct {
	value *chat.ChatConversation
	repo  chat.IChatRepository
	user  *chatUserAggregateRoot
}

func newConversation(value *chat.ChatConversation, user *chatUserAggregateRoot, repo chat.IChatRepository) chat.IConversation {
	return &converstationImpl{
		value: value,
		repo:  repo,
		user:  user,
	}
}

// GetDomainId implements chat.IConversation.
func (c *converstationImpl) GetDomainId() int {
	return c.value.Id
}

// DeleteMsg implements chat.IConversation.
func (c *converstationImpl) DeleteMsg(msgId int) error {
	return c.repo.Msg().Delete(&chat.ChatMsg{Id: msgId})
}

// Destroy implements chat.IConversation.
func (c *converstationImpl) Destroy() error {
	_, err := c.repo.Msg().DeleteBy("conv_id = ?", c.GetDomainId())
	if err == nil {
		err = c.repo.Conversation().Delete(&chat.ChatConversation{Id: c.GetDomainId()})
	}
	return err
}

// FetchHistoryMsgList implements chat.IConversation.
func (c *converstationImpl) FetchHistoryMsgList(lastTime int, size int) {
	panic("unimplemented")
}

// FetchMsgList implements chat.IConversation.
func (c *converstationImpl) FetchMsgList(lastTime int, size int) []*chat.ChatMsg {
	panic("unimplemented")
}

// Greet implements chat.IConversation.
func (c *converstationImpl) Greet(msg string) error {
	if c.value.Sid == c.user.GetAggregateRootId() {
		return errors.New("replay user can't send greet msg")
	}
	v := &chat.MsgBody{
		MsgType: 1,
		Content: msg,
		Extrta:  "{\"greet\":1}",
	}
	_, err := c.sendMsg(v, 0)
	return err
}

func (c *converstationImpl) sendMsg(v *chat.MsgBody, expiresTime int) (int, error) {
	v.CreateTime = int(time.Now().Unix())
	v.ExpiresTime = expiresTime
	v.ConvId = c.GetDomainId()
	// 是否为发送人的消息
	v.SidMsg = types.Ternary(c.value.Sid == c.user.GetAggregateRootId(), 1, 0)
	_, err := c.repo.Msg().Save(v)
	return v.Id, err
}

// RevertMsg implements chat.IConversation.
func (c *converstationImpl) RevertMsg(msgId int) error {
	panic("unimplemented")
}

// Send implements chat.IConversation.
func (c *converstationImpl) Send(msg *chat.MsgBody) (int, error) {
	panic("unimplemented")
}

// UpdateMsgAttrs implements chat.IConversation.
func (c *converstationImpl) UpdateMsgAttrs(msgId int, attrs map[string]interface{}) error {
	panic("unimplemented")
}
