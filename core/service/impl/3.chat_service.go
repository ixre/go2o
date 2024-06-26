package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/chat"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.ChatServiceServer = new(chatServiceImpl)

type chatServiceImpl struct {
	repo chat.IChatRepository
	proto.UnimplementedChatServiceServer
	serviceUtil
}

func NewChatService(repo chat.IChatRepository) proto.ChatServiceServer {
	return &chatServiceImpl{repo: repo}
}

// GetConversation implements proto.ChatServiceServer.
func (c *chatServiceImpl) GetConversation(_ context.Context,req *proto.ChatConversationRequest) (*proto.ChatConversationResponse, error) {
	iu := c.repo.GetChatUser(int(req.Uid))
	ic,err := iu.BuildConversation(int(req.Rid),chat.ChatType(req.ChatType))
	if err != nil{
		return &proto.ChatConversationResponse{
			ErrCode: 1,
			LastMsg: err.Error(),
		},nil
	}
	v := ic.Get()
	return &proto.ChatConversationResponse{
		ConvId:       int64(ic.GetDomainId()),
		Key:          v.Key,
		Sid:          int64(v.Sid),
		Rid:          int64(v.Rid),
		ChatType:     int32(v.ChatType),
		LastMsg:      v.LastMsg,
		LastChatTime: int64(v.LastChatTime),
	},nil
}


// DeleteMsg implements proto.ChatServiceServer.
func (c *chatServiceImpl) DeleteMsg(_ context.Context,req *proto.MsgIdRequest) (*proto.Result, error) {
	var err error
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		err = errors.New("no such conversation")
	}else{
		err = ic.DeleteMsg(int(req.MsgId))
	}
	return c.error(err),nil
}

// DestroyConversation implements proto.ChatServiceServer.
func (c *chatServiceImpl) DestroyConversation(_ context.Context, req *proto.ConversationIdRequest) (*proto.Result, error) {
	var err error
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		err = errors.New("no such conversation")
	}else{
		err = ic.Destroy()
	}
	return c.error(err),nil
}

// FetchHistoryMsgList implements proto.ChatServiceServer.
func (c *chatServiceImpl) FetchMsgList(_ context.Context, req *proto.FetchMsgRequest) (*proto.FetchMsgResponse, error) {
	var err error
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		return &proto.FetchMsgResponse{}, err
	}
	var list []*chat.ChatMsg
	if req.HistoryMsg {
	list = ic.FetchHistoryMsgList(int(req.LastTime),int(req.Size))
	}else{
		list = ic.FetchMsgList(int(req.LastTime),int(req.Size))
	}
	arr := make([]*proto.SMsg,len(list))
	for i,v := range list{
		arr[i] = c.parseChatMsg(v)
	}
	return &proto.FetchMsgResponse{
		MsgList: arr,
		IsOver:  false,
	},nil
}

func (c *chatServiceImpl) parseChatMsg(v *chat.ChatMsg)*proto.SMsg{
	var mp map[string]string
		json.Unmarshal(v.Extra, &mp)
	return &proto.SMsg{
		MsgId:          int64(v.Id),
		Sid:         int64(v.Sid),
		MsgType:     int32(v.MsgType),
		MsgFlag:     int32(v.MsgFlag),
		Content:     v.Content,
		Extra:       mp,
		CreateTime:  v.CreateTime,
	}
}


// GetMsg implements proto.ChatServiceServer.
func (c *chatServiceImpl) GetMsg(_ context.Context,req *proto.MsgIdRequest) (*proto.SMsg, error) {
	var err error
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		return &proto.SMsg{}, err
	}
	msg := ic.GetMsg(req.MsgId)
	return c.parseChatMsg(msg),nil
}

// RevertMsg implements proto.ChatServiceServer.
func (c *chatServiceImpl) RevertMsg(_ context.Context,req *proto.MsgIdRequest) (*proto.Result, error) {
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		err = errors.New("no such conversation")
	}else{
		err = ic.RevertMsg(int(req.MsgId))
	}
	return c.error(err),nil
}

// Send implements proto.ChatServiceServer.
func (c *chatServiceImpl) Send(_ context.Context,req *proto.SendMsgRequest) (*proto.SendMsgResponse, error) {
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		err = errors.New("no such conversation")
	}else{
		err = ic.Send(&chat.MsgBody{
			MsgType: int(req.MsgType),
			Content: req.Content,
			Extrta:  req.Extrta,
		})
	}
	return c.error(err),nil
}

// UpdateMsgAttrs implements proto.ChatServiceServer.
func (c *chatServiceImpl) UpdateMsgAttrs(_ context.Context, *proto.UpdateMsgAttrRequest) (*proto.Result, error) {
	iu := c.repo.GetChatUser(int(req.Uid))
	ic := iu.GetConversation(int(req.ConvId))
	if ic == nil{
		err = errors.New("no such conversation")
	}else{
		err = ic.RevertMsg(int(req.MsgId))
	}
	return c.error(err),nil
}
