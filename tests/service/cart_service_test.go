package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestCombileCart(t *testing.T) {
	ret, _ := impl.CartService.GetShoppingCart(context.TODO(), &proto.ShoppingCartId{
		UserId:      1,
		CartCode:    "0742c193-9642-11ed-9649-0242ac1a0003",
		IsWholesale: false,
	})
	t.Log(typeconv.MustJson(ret))
}
