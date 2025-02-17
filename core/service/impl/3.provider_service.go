package impl

import (
	"context"
	"encoding/base64"

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
	if len(req.Type) == 0 {
		return &proto.UserOpenIdResponse{
			Code:    1,
			Message: "缺少参数: clientType",
		}, nil
	}
	if len(req.Code) == 0 {
		return &proto.UserOpenIdResponse{
			Code:    1,
			Message: "缺少参数: clientCode",
		}, nil
	}
	ret, err := tencent.WECHAT.GetMiniProgramOpenId("", req.Code)
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

// GetMPCode 获取小程序二维码
func (s *serviceProviderServiceImpl) GetMPCode(_ context.Context, req *proto.MPCodeRequest) (*proto.MPQrCodeResponse, error) {
	bytes, err := tencent.WECHAT.GetMiniProgramUnlimitCode("", req.Page, req.Scene,
		req.SaveLocal, req.OwnerKey)
	if err != nil {
		return &proto.MPQrCodeResponse{
			Code:    1001,
			Message: err.Error(),
		}, nil
	}
	base64Img := base64.StdEncoding.EncodeToString(bytes)
	return &proto.MPQrCodeResponse{
		QrCodeUrl: "data:image/png;base64," + base64Img,
	}, nil
}
