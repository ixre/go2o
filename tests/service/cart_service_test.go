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
		CartCode:    "b61d09ec4a3782cd",
		IsWholesale: false,
	}
	ret, _ := impl.CartService.GetShoppingCart(context.TODO(), id)
	impl.CartService.PutItems(context.TODO(), &proto.CartItemRequest{
		CartId: &proto.ShoppingCartId{
			UserId:      id.UserId,
			CartCode:    ret.CartCode,
			IsWholesale: false,
		},
		Items: []*proto.RCartItem{
			{
				ItemId:   187,
				SkuId:    0,
				Quantity: 1,
			}},
	})
	t.Log(id.CartCode)
	t.Log(ret.CartCode)
	time.Sleep(1000)
	//t.Log(typeconv.MustJson(ret))
}
