package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/typeconv"
)

func TestGetItem(t *testing.T) {
	goods, _ := inject.GetItemService().GetItem(context.TODO(), &proto.GetItemRequest{
		ItemId: 1,
	})
	t.Log(typeconv.MustJson(goods))
}
func TestGetItemSku(t *testing.T) {
	goods, _ := inject.GetItemService().GetSku(context.TODO(), &proto.SkuId{
		ItemId: 3272,
		SkuId:  0,
	})
	t.Log(typeconv.MustJson(goods))

	goods2, _ := inject.GetItemService().GetItemBySku(context.TODO(), &proto.ItemBySkuRequest{
		ProductId: 3272,
		SkuId:     0,
	})
	t.Log(typeconv.MustJson(goods2))
}
