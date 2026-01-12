package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/pkg/inject"
	"github.com/ixre/go2o/pkg/service/proto"
)

// 测试查询自营店铺列表
func TestQuerySelfSupportShopList(t *testing.T) {
	rsp, _ := inject.GetShopService().GetSelfSupportShops(context.TODO(),
		&proto.SelfSupportShopRequest{})
	t.Log("shop list", rsp)
}
