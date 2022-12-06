package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
)

// 测试查询自营店铺列表
func TestQuerySelfSupportShopList(t *testing.T) {
	rsp, _ := impl.ShopService.GetSelfSupportShops(context.TODO(),
		&proto.SelfSupportShopRequest{})
	t.Log("shop list", rsp)
}
