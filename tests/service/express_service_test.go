package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

// TestGetProviderGroup 测试快递分组
func TestGetProviderGroup(t *testing.T) {
	ret, _ := inject.GetExpressService().GetProviderGroup(context.TODO(),
		&proto.Empty{})
	t.Log(typeconv.MustJson(ret.List))
}
