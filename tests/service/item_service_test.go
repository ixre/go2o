package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestGetItem(t *testing.T) {
	goods, _ := impl.ItemService.GetItem(context.TODO(), &proto.GetItemRequest{
		ItemId: 1,
	})
	t.Log(typeconv.MustJson(goods))
}
