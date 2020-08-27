package grpc

import (
	"context"
	"go2o/core/service/proto"
)
var _ proto.GreeterServiceHandler = new(TestServiceImpl)
type TestServiceImpl struct{
}

func (t TestServiceImpl) Hello(ctx context.Context, user *proto.User, response *proto.UserResponse) error {
	response.Name = user.Name
	response.State = proto.EState_Normal
	return nil
}


