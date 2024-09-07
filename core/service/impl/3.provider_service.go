package impl

import (
	"context"

	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/sp/tencent"
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
func (s *serviceProviderServiceImpl) GetOpenId(_ context.Context, req *proto.GetUserOpenIdRequest) (*proto.UserOpenIdResponse, error) {
	ret, err := tencent.WECHAT.GetOpenId(req.Code, nil)
	if err != nil {
		return &proto.UserOpenIdResponse{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &proto.UserOpenIdResponse{
		OpenId:  ret.OpenID,
		UnionId: ret.UnionID,
		AppId:   ret.AppId,
	}, nil
}
