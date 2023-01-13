package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/service"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestRegistryUpateValueMap(t *testing.T) {
	mp := map[string]string{
		"order_push_sub_order_enabled":    "0",
		"member_withdrawal_push_enabled1": "1",
	}
	ret, _ := impl.RegistryService.UpdateValues(context.TODO(), &proto.StringMap{
		Value: mp,
	})
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
	}
	t.Log(typeconv.MustJson(ret))
}

func TestRegistryClientUpateValueMap(t *testing.T) {
	mp := map[string]string{
		"order_push_sub_order_enabled": "0",
		"member_withdrawal_push_enabled1": "1",
	}
	service.ConfigureClient(nil, "192.168.0.159:1427")
	trans, cli, _ := service.RegistryServiceClient()
	defer trans.Close()
	ret, err := cli.UpdateValues(context.TODO(), &proto.StringMap{
		Value: mp,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if ret.ErrCode > 0 {
		t.Error(ret.ErrMsg)
	}
	t.Log(typeconv.MustJson(ret))
}
