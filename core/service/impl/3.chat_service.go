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

// ApplyItem implements proto.CartServiceServer.
func (c *chatServiceImpl) ApplyItem(context.Context, *proto.CartItemOpRequest) (*proto.CartItemResponse, error) {
	panic("unimplemented")
}

// CheckCart implements proto.CartServiceServer.
func (c *chatServiceImpl) CheckCart(context.Context, *proto.CheckCartRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// GetShoppingCart implements proto.CartServiceServer.
func (c *chatServiceImpl) GetShoppingCart(context.Context, *proto.ShoppingCartId) (*proto.SShoppingCart, error) {
	panic("unimplemented")
}

// WholesaleCartV1 implements proto.CartServiceServer.
func (c *chatServiceImpl) WholesaleCartV1(context.Context, *proto.WsCartRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedCartServiceServer implements proto.CartServiceServer.
func (c *chatServiceImpl) mustEmbedUnimplementedCartServiceServer() {
	panic("unimplemented")
}
