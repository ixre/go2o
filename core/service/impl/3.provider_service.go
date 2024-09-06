package impl

import (
	"context"

	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.ServiceProviderServiceServer = new(serviceProviderServiceImpl)

// 第三方服务提供者服务
type serviceProviderServiceImpl struct {
	proto.UnimplementedServiceProviderServiceServer
	serviceUtil
}

func NewServiceProviderService() proto.ServiceProviderServiceServer {
	return &serviceProviderServiceImpl{}
}

// GetOpenId implements proto.ServiceProviderServiceServer.
func (s *serviceProviderServiceImpl) GetOpenId(context.Context, *proto.GetUserOpenIdRequest) (*proto.UserOpenIdResponse, error) {
	panic("unimplemented")
}
