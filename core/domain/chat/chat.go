package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ixre/go2o/core/domain/interface/chat"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/gof/crypto"
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
func (c *chatUserAggregateRoot) GetConversation(convId int) chat.IConversation {
	v := c.repo.Conversation().Get(convId)
	if v == nil {
		return nil
	}
	return newConversation(v, c, c.repo)
}

// BuildConversation implements chat.IChatUserAggregateRoot.
func (c *chatUserAggregateRoot) BuildConversation(rid int, chatType chat.ChatType, outOrderNo string) (chat.IConversation, error) {

	uid := []int{c.GetAggregateRootId(), rid}
	cr := c.repo.Conversation()
	conv := cr.FindBy("sid IN (?) AND rid IN(?)", uid, uid)
	if conv == nil {
		key := fmt.Sprintf("%d-%d#%d", uid[0], uid[1], chatType)
		key = crypto.Md5([]byte(key))[8:24]
		// outOrderNo只有在新增时才指定OutOrderNo，后面只能通过BindOutOrderNo更改
		v := &chat.ChatConversation{
			Sid:        c.GetAggregateRootId(),
			Rid:        rid,
			Key:        key,
			ChatType:   int(chatType),
			OutOrderNo: outOrderNo,
		}
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

// Get implements chat.IConversation.
func (c *converstationImpl) Get() chat.ChatConversation {
	return *types.DeepClone(c.value)
}

// GetMsg implements chat.IConversation.
func (c *converstationImpl) GetMsg(msgId int) *chat.ChatMsg {
	msg, err := c.getMsg(msgId)
	if err != nil || msg == nil {
		return nil
	}
	return types.DeepClone(msg)
}

// DeleteMsg implements chat.IConversation.
func (c *converstationImpl) DeleteMsg(msgId int) error {
	msg, err := c.getMsg(msgId)
	if err != nil || msg == nil || msg.Sid != c.user.GetAggregateRootId() {
		return errors.New("no such message")
	}
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

// BindOutOrderNo implements chat.IConversation.
func (c *converstationImpl) BindOutOrderNo(outOrderNo string) error {
	if len(outOrderNo) <= 0 {
		return errors.New("业务关联单号为空")
	}
	if c.value.OutOrderNo != outOrderNo {
		c.value.OutOrderNo = outOrderNo
		c.value.UpdateTime = int(time.Now().Unix())
		_, err := c.repo.Conversation().Save(c.value)
		return err
	}
	return nil
}

// FetchHistoryMsgList implements chat.IConversation.
func (c *converstationImpl) FetchHistoryMsgList(lastTime int, size int) []*chat.ChatMsg {
	return c.repo.Msg().FindList(&fw.QueryOption{
		Limit: size,
		Order: "create_time DESC",
	}, "create_time < ?", lastTime)
}

// FetchMsgList implements chat.IConversation.
func (c *converstationImpl) FetchMsgList(lastTime int, size int) []*chat.ChatMsg {
	return c.repo.Msg().FindList(&fw.QueryOption{
		Limit: size,
		Order: "create_time ASC",
	}, "create_time > ?", lastTime)
}

// Greet implements chat.IConversation.
func (c *converstationImpl) Greet(msg string) error {
	if c.value.Sid == c.user.GetAggregateRootId() {
		return errors.New("replay user can't send greet msg")
	}
	v := &chat.MsgBody{
		MsgType: 1,
		Content: msg,
		Extra:   nil,
	}
	_, err := c.sendMsg(v, chat.MsgFlagHint, 0)
	return err
}

func (c *converstationImpl) sendMsg(v *chat.MsgBody, flag int, expiresTime int) (int, error) {
	// 是否为发送人的消息
	sid := c.user.GetAggregateRootId()
	extra := ""
	if v.Extra != nil && len(v.Extra) > 0 {
		bytes, err := json.Marshal(v.Extra)
		if err != nil {
			logger.Error("chat msg marshal error", err)
		}
		extra = string(bytes)
	}
	msg := &chat.ChatMsg{
		ConvId:      c.GetDomainId(),
		MsgType:     int(v.MsgType),
		Sid:         sid,
		MsgFlag:     flag,
		Content:     v.Content,
		Extra:       extra,
		ExpiresTime: expiresTime,
	}
	return c.saveMsg(msg)
}

func (c *converstationImpl) saveMsg(msg *chat.ChatMsg) (int, error) {
	unix := int(time.Now().Unix())
	if msg.CreateTime <= 0 {
		msg.CreateTime = unix
	}
	isNewMsg := msg.Id == 0
	_, err := c.repo.Msg().Save(msg)
	if err == nil && isNewMsg {
		// 更新会话最后聊天时间和聊天内容
		c.value.LastChatTime = unix
		c.value.LastMsg = msg.Content
		_, err = c.repo.Conversation().Save(c.value)
	}
	return msg.Id, err
}

// RevertMsg implements chat.IConversation.
func (c *converstationImpl) RevertMsg(msgId int) error {
	msg, err := c.getMsg(msgId)
	if err == nil {
		msg.MsgFlag |= chat.MsgFlagRevert
		_, err = c.saveMsg(msg)
	}
	return err
}

// getMsg 获取消息并检测消息合法性
func (c *converstationImpl) getMsg(msgId int) (*chat.ChatMsg, error) {
	msg := c.repo.Msg().Get(&chat.ChatMsg{Id: msgId})
	if msg == nil {
		return msg, errors.New("msg not found")
	}
	if (msg.MsgFlag & chat.MsgFlagRevert) == chat.MsgFlagRevert {
		return msg, errors.New("msg has been revert")
	}
	if (msg.MsgFlag & chat.MsgFlagDelete) == chat.MsgFlagDelete {
		return msg, errors.New("msg has been deleted")
	}
	if c.user.GetAggregateRootId() != msg.Sid {
		// 非发送人的消息, 不能撤回
		return msg, errors.New("can't revert msg")
	}
	return msg, nil
}

// Send implements chat.IConversation.
func (c *converstationImpl) Send(msg *chat.MsgBody) (int, error) {
	// 默认消息30天过期
	expires := int(time.Now().Add(time.Hour * 24 * 30).Unix())
	return c.sendMsg(msg, 0, expires)
}

// UpdateMsgAttrs implements chat.IConversation.
func (c *converstationImpl) UpdateMsgAttrs(msgId int, attrs map[string]string) error {
	msg, err := c.getMsg(msgId)
	if err == nil {
		if attrs == nil {
			msg.Extra = ""
		} else {
			bytes, err := json.Marshal(attrs)
			if err != nil {
				logger.Error("json marshal error chat msg attribute", err)
			}
			msg.Extra = string(bytes)
		}
		_, err = c.saveMsg(msg)
	}
	return err
}
