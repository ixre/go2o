package impl

import (
	"context"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.StatusServiceServer = new(statusServiceImpl)

type statusServiceImpl struct {
}

func NewStatusService() *statusServiceImpl {
	return &statusServiceImpl{}
}

func (s *statusServiceImpl) Ping(_ context.Context, empty *proto.Empty) (*proto.String, error) {
	return &proto.String{Value: "pong"}, nil
}
