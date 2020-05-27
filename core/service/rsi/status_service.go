package rsi

import (
	"context"
	"go2o/core/service/auto_gen/rpc/status_service"
)

var _ status_service.StatusService = new(statusServiceImpl)

type statusServiceImpl struct {
}

func NewStatusService() *statusServiceImpl {
	return &statusServiceImpl{}
}

func (s *statusServiceImpl) Ping(ctx context.Context) (r string, err error) {
	return "pong", nil
}
