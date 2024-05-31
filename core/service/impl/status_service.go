package impl

import (
	"context"

	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.StatusServiceServer = new(StatusServiceImpl)

type StatusServiceImpl struct {
	proto.UnimplementedStatusServiceServer
}

func NewStatusService() proto.StatusServiceServer {
	return &StatusServiceImpl{}
}

func (s *StatusServiceImpl) Ping(_ context.Context, empty *proto.Empty) (*proto.String, error) {
	return &proto.String{Value: "pong"}, nil
}

func (t *StatusServiceImpl) Hello(_ context.Context, user *proto.String) (response *proto.String, err error) {
	return &proto.String{
		Value: "hello " + user.Value,
	}, nil
}
