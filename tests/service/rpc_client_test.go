package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/service"
	"github.com/ixre/go2o/core/service/proto"
)

// 测试连接rpc服务
func TestRpcStatus(t *testing.T) {
	service.ConfigureClient(nil, "115.55.65.123:1427")
	trans, cli, err := service.StatusServiceClient()
	if err != nil {
		t.Error(err)
	} else {
		v, err := cli.Ping(context.TODO(), &proto.Empty{})
		if err != nil {
			t.Error(err)
		} else {
			t.Log("---result:", v.Value)
		}
		trans.Close()
	}
}
