package service

import (
	"context"
	"testing"
	"time"

	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
)

func TestCombileCart(t *testing.T) {
	id := &proto.ShoppingCartId{
		UserId:      1,
		CartCode:    "af7d4093-9645-11ed-b327-50ebf6326d4d",
		IsWholesale: false,
	}
	ret, _ := impl.CartService.GetShoppingCart(context.TODO(), id)
	impl.CartService.PutInCart(context.TODO(), &proto.CartItemRequest{
		Id: &proto.ShoppingCartId{
			UserId:      id.UserId,
			CartCode:    ret.CartCode,
			IsWholesale: false,
		},
		Item: &proto.RCartItem{
			ItemId:    187,
			SkuId:     0,
			Quantity:  1,
			CheckOnly: true,
		},
	})
	t.Log(id.CartCode)
	t.Log(ret.CartCode)
	time.Sleep(1000)
	//t.Log(typeconv.MustJson(ret))
}
