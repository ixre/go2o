package impl

import (
	"context"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.GreeterServiceServer = new(TestServiceImpl)

type TestServiceImpl struct {
}

func (t *TestServiceImpl) Hello(_ context.Context, user *proto.User1) (response *proto.UserResponse, err error) {
	rsp := &proto.UserResponse{
		Name:  user.Name,
		State: proto.EState_Normal,
	}
	return rsp, nil
}
