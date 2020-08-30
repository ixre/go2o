package grpc

import (
	"context"
	"go2o/core/service/proto"
)
var _ proto.GreeterServiceHandler = new(TestServiceImpl)
type TestServiceImpl struct{
}

func (t *TestServiceImpl) Hello(ctx context.Context, user *proto.User, response *proto.UserResponse) error {
	response.Name = user.Name
	response.State = proto.EState_Normal
	return nil
}


type Helloworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Helloworld) Call(ctx context.Context, req *helloworld.Request, rsp *helloworld.Response) error {
	return nil
}

