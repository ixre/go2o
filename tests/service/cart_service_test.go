package service

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
)

// 测试立即结算
func TestCheckoutItem(t *testing.T) {
	ret2, _ := inject.GetCartService().ApplyItem(context.TODO(), &proto.CartItemOpRequest{
		CartId: &proto.ShoppingCartId{
			CartCode:    "2f7d4093-9645-11ed-b327-50ebf6326d4d",
			IsWholesale: false,
		},
		Item: &proto.RCartItem{
			ItemId:   3274,
			SkuId:    78,
			Quantity: 1,
		},
		Op: proto.ECartItemOp_CHECKOUT,
	})
	log.Println("---", ret2.ErrCode)
	time.Sleep(1000)
	//t.Log(typeconv.MustJson(ret))
}

// 测试合并购物车
func TestCombileCart(t *testing.T) {
	id := &proto.ShoppingCartId{
		UserId:      1,
		CartCode:    "b61d09ec4a3782cd",
		IsWholesale: false,
	}
	ret, _ := inject.GetCartService().GetShoppingCart(context.TODO(), id)
	inject.GetCartService().ApplyItem(context.TODO(), &proto.CartItemOpRequest{
		CartId: &proto.ShoppingCartId{
			UserId:      id.UserId,
			CartCode:    ret.CartCode,
			IsWholesale: false,
		},
		Item: &proto.RCartItem{
			ItemId:   187,
			SkuId:    0,
			Quantity: 1,
		},
		Op: proto.ECartItemOp_PUT,
	})
	t.Log(id.CartCode)
	t.Log(ret.CartCode)
	time.Sleep(1000)
	//t.Log(typeconv.MustJson(ret))
}
