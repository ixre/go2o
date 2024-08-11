// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: advertisement_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AdvertisementService_GetGroups_FullMethodName              = "/AdvertisementService/GetGroups"
	AdvertisementService_GetPosition_FullMethodName            = "/AdvertisementService/GetPosition"
	AdvertisementService_SaveAdPosition_FullMethodName         = "/AdvertisementService/SaveAdPosition"
	AdvertisementService_DeleteAdPosition_FullMethodName       = "/AdvertisementService/DeleteAdPosition"
	AdvertisementService_PutDefaultAd_FullMethodName           = "/AdvertisementService/PutDefaultAd"
	AdvertisementService_QueryAd_FullMethodName                = "/AdvertisementService/QueryAd"
	AdvertisementService_QueryAdvertisementData_FullMethodName = "/AdvertisementService/QueryAdvertisementData"
	AdvertisementService_SetUserAd_FullMethodName              = "/AdvertisementService/SetUserAd"
	AdvertisementService_GetAdvertisement_FullMethodName       = "/AdvertisementService/GetAdvertisement"
	AdvertisementService_SaveAd_FullMethodName                 = "/AdvertisementService/SaveAd"
	AdvertisementService_DeleteAd_FullMethodName               = "/AdvertisementService/DeleteAd"
)

// AdvertisementServiceClient is the client API for AdvertisementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdvertisementServiceClient interface {
	// * 获取广告组
	GetGroups(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AdGroupResponse, error)
	// * 获取广告位
	GetPosition(ctx context.Context, in *AdPositionId, opts ...grpc.CallOption) (*SAdPosition, error)
	// * 更新广告位
	SaveAdPosition(ctx context.Context, in *SAdPosition, opts ...grpc.CallOption) (*Result, error)
	// * 删除广告位
	DeleteAdPosition(ctx context.Context, in *AdPositionId, opts ...grpc.CallOption) (*Result, error)
	// 投放广告位的默认广告
	PutDefaultAd(ctx context.Context, in *SetDefaultAdRequest, opts ...grpc.CallOption) (*Result, error)
	// 查询广告
	QueryAd(ctx context.Context, in *QueryAdRequest, opts ...grpc.CallOption) (*QueryAdResponse, error)
	// 查询广告并返回广告数据
	QueryAdvertisementData(ctx context.Context, in *QueryAdvertisementDataRequest, opts ...grpc.CallOption) (*QueryAdvertisementDataResponse, error)
	// 用户投放广告
	SetUserAd(ctx context.Context, in *SetUserAdRequest, opts ...grpc.CallOption) (*Result, error)
	// 获取广告,returnData=true返回数据传输对象
	GetAdvertisement(ctx context.Context, in *AdIdRequest, opts ...grpc.CallOption) (*SAdDto, error)
	// 保存广告,更新时不允许修改类型
	SaveAd(ctx context.Context, in *SaveAdRequest, opts ...grpc.CallOption) (*TxResult, error)
	// 删除广告
	DeleteAd(ctx context.Context, in *AdIdRequest, opts ...grpc.CallOption) (*Result, error)
}

type advertisementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdvertisementServiceClient(cc grpc.ClientConnInterface) AdvertisementServiceClient {
	return &advertisementServiceClient{cc}
}

func (c *advertisementServiceClient) GetGroups(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AdGroupResponse, error) {
	out := new(AdGroupResponse)
	err := c.cc.Invoke(ctx, AdvertisementService_GetGroups_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) GetPosition(ctx context.Context, in *AdPositionId, opts ...grpc.CallOption) (*SAdPosition, error) {
	out := new(SAdPosition)
	err := c.cc.Invoke(ctx, AdvertisementService_GetPosition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) SaveAdPosition(ctx context.Context, in *SAdPosition, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, AdvertisementService_SaveAdPosition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) DeleteAdPosition(ctx context.Context, in *AdPositionId, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, AdvertisementService_DeleteAdPosition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) PutDefaultAd(ctx context.Context, in *SetDefaultAdRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, AdvertisementService_PutDefaultAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) QueryAd(ctx context.Context, in *QueryAdRequest, opts ...grpc.CallOption) (*QueryAdResponse, error) {
	out := new(QueryAdResponse)
	err := c.cc.Invoke(ctx, AdvertisementService_QueryAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) QueryAdvertisementData(ctx context.Context, in *QueryAdvertisementDataRequest, opts ...grpc.CallOption) (*QueryAdvertisementDataResponse, error) {
	out := new(QueryAdvertisementDataResponse)
	err := c.cc.Invoke(ctx, AdvertisementService_QueryAdvertisementData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) SetUserAd(ctx context.Context, in *SetUserAdRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, AdvertisementService_SetUserAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) GetAdvertisement(ctx context.Context, in *AdIdRequest, opts ...grpc.CallOption) (*SAdDto, error) {
	out := new(SAdDto)
	err := c.cc.Invoke(ctx, AdvertisementService_GetAdvertisement_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) SaveAd(ctx context.Context, in *SaveAdRequest, opts ...grpc.CallOption) (*TxResult, error) {
	out := new(TxResult)
	err := c.cc.Invoke(ctx, AdvertisementService_SaveAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *advertisementServiceClient) DeleteAd(ctx context.Context, in *AdIdRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, AdvertisementService_DeleteAd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdvertisementServiceServer is the server API for AdvertisementService service.
// All implementations must embed UnimplementedAdvertisementServiceServer
// for forward compatibility
type AdvertisementServiceServer interface {
	// * 获取广告组
	GetGroups(context.Context, *Empty) (*AdGroupResponse, error)
	// * 获取广告位
	GetPosition(context.Context, *AdPositionId) (*SAdPosition, error)
	// * 更新广告位
	SaveAdPosition(context.Context, *SAdPosition) (*Result, error)
	// * 删除广告位
	DeleteAdPosition(context.Context, *AdPositionId) (*Result, error)
	// 投放广告位的默认广告
	PutDefaultAd(context.Context, *SetDefaultAdRequest) (*Result, error)
	// 查询广告
	QueryAd(context.Context, *QueryAdRequest) (*QueryAdResponse, error)
	// 查询广告并返回广告数据
	QueryAdvertisementData(context.Context, *QueryAdvertisementDataRequest) (*QueryAdvertisementDataResponse, error)
	// 用户投放广告
	SetUserAd(context.Context, *SetUserAdRequest) (*Result, error)
	// 获取广告,returnData=true返回数据传输对象
	GetAdvertisement(context.Context, *AdIdRequest) (*SAdDto, error)
	// 保存广告,更新时不允许修改类型
	SaveAd(context.Context, *SaveAdRequest) (*TxResult, error)
	// 删除广告
	DeleteAd(context.Context, *AdIdRequest) (*Result, error)
	mustEmbedUnimplementedAdvertisementServiceServer()
}

// UnimplementedAdvertisementServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdvertisementServiceServer struct {
}

func (UnimplementedAdvertisementServiceServer) GetGroups(context.Context, *Empty) (*AdGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroups not implemented")
}
func (UnimplementedAdvertisementServiceServer) GetPosition(context.Context, *AdPositionId) (*SAdPosition, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPosition not implemented")
}
func (UnimplementedAdvertisementServiceServer) SaveAdPosition(context.Context, *SAdPosition) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveAdPosition not implemented")
}
func (UnimplementedAdvertisementServiceServer) DeleteAdPosition(context.Context, *AdPositionId) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAdPosition not implemented")
}
func (UnimplementedAdvertisementServiceServer) PutDefaultAd(context.Context, *SetDefaultAdRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutDefaultAd not implemented")
}
func (UnimplementedAdvertisementServiceServer) QueryAd(context.Context, *QueryAdRequest) (*QueryAdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAd not implemented")
}
func (UnimplementedAdvertisementServiceServer) QueryAdvertisementData(context.Context, *QueryAdvertisementDataRequest) (*QueryAdvertisementDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAdvertisementData not implemented")
}
func (UnimplementedAdvertisementServiceServer) SetUserAd(context.Context, *SetUserAdRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserAd not implemented")
}
func (UnimplementedAdvertisementServiceServer) GetAdvertisement(context.Context, *AdIdRequest) (*SAdDto, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdvertisement not implemented")
}
func (UnimplementedAdvertisementServiceServer) SaveAd(context.Context, *SaveAdRequest) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveAd not implemented")
}
func (UnimplementedAdvertisementServiceServer) DeleteAd(context.Context, *AdIdRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAd not implemented")
}
func (UnimplementedAdvertisementServiceServer) mustEmbedUnimplementedAdvertisementServiceServer() {}

// UnsafeAdvertisementServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdvertisementServiceServer will
// result in compilation errors.
type UnsafeAdvertisementServiceServer interface {
	mustEmbedUnimplementedAdvertisementServiceServer()
}

func RegisterAdvertisementServiceServer(s grpc.ServiceRegistrar, srv AdvertisementServiceServer) {
	s.RegisterService(&AdvertisementService_ServiceDesc, srv)
}

func _AdvertisementService_GetGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).GetGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_GetGroups_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).GetGroups(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_GetPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdPositionId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).GetPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_GetPosition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).GetPosition(ctx, req.(*AdPositionId))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_SaveAdPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SAdPosition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).SaveAdPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_SaveAdPosition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).SaveAdPosition(ctx, req.(*SAdPosition))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_DeleteAdPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdPositionId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).DeleteAdPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_DeleteAdPosition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).DeleteAdPosition(ctx, req.(*AdPositionId))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_PutDefaultAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDefaultAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).PutDefaultAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_PutDefaultAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).PutDefaultAd(ctx, req.(*SetDefaultAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_QueryAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).QueryAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_QueryAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).QueryAd(ctx, req.(*QueryAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_QueryAdvertisementData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAdvertisementDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).QueryAdvertisementData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_QueryAdvertisementData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).QueryAdvertisementData(ctx, req.(*QueryAdvertisementDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_SetUserAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUserAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).SetUserAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_SetUserAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).SetUserAd(ctx, req.(*SetUserAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_GetAdvertisement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).GetAdvertisement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_GetAdvertisement_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).GetAdvertisement(ctx, req.(*AdIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_SaveAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveAdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).SaveAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_SaveAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).SaveAd(ctx, req.(*SaveAdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdvertisementService_DeleteAd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdvertisementServiceServer).DeleteAd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdvertisementService_DeleteAd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdvertisementServiceServer).DeleteAd(ctx, req.(*AdIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AdvertisementService_ServiceDesc is the grpc.ServiceDesc for AdvertisementService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdvertisementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AdvertisementService",
	HandlerType: (*AdvertisementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGroups",
			Handler:    _AdvertisementService_GetGroups_Handler,
		},
		{
			MethodName: "GetPosition",
			Handler:    _AdvertisementService_GetPosition_Handler,
		},
		{
			MethodName: "SaveAdPosition",
			Handler:    _AdvertisementService_SaveAdPosition_Handler,
		},
		{
			MethodName: "DeleteAdPosition",
			Handler:    _AdvertisementService_DeleteAdPosition_Handler,
		},
		{
			MethodName: "PutDefaultAd",
			Handler:    _AdvertisementService_PutDefaultAd_Handler,
		},
		{
			MethodName: "QueryAd",
			Handler:    _AdvertisementService_QueryAd_Handler,
		},
		{
			MethodName: "QueryAdvertisementData",
			Handler:    _AdvertisementService_QueryAdvertisementData_Handler,
		},
		{
			MethodName: "SetUserAd",
			Handler:    _AdvertisementService_SetUserAd_Handler,
		},
		{
			MethodName: "GetAdvertisement",
			Handler:    _AdvertisementService_GetAdvertisement_Handler,
		},
		{
			MethodName: "SaveAd",
			Handler:    _AdvertisementService_SaveAd_Handler,
		},
		{
			MethodName: "DeleteAd",
			Handler:    _AdvertisementService_DeleteAd_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "advertisement_service.proto",
}
