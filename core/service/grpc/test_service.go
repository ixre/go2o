package grpc

import (
	"context"
	"go2o/core/service/proto"
)
var _ proto.GreeterServiceServer = new(TestServiceImpl)
type TestServiceImpl struct{
}

func (t *TestServiceImpl) Hello(ctx context.Context, user *proto.User)(response *proto.UserResponse,err error) {
	rsp := &proto.UserResponse{
		Name:  user.Name,
		State: proto.EState_Normal,
	}
	return rsp, nil
}

