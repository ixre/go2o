// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.0
// source: status_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	StatusService_Ping_FullMethodName  = "/StatusService/Ping"
	StatusService_Hello_FullMethodName = "/StatusService/Hello"
)

// StatusServiceClient is the client API for StatusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// * 状态服务
type StatusServiceClient interface {
	// 发送心跳
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*String, error)
	// 打招呼
	Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
}

type statusServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStatusServiceClient(cc grpc.ClientConnInterface) StatusServiceClient {
	return &statusServiceClient{cc}
}

func (c *statusServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*String, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(String)
	err := c.cc.Invoke(ctx, StatusService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusServiceClient) Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(String)
	err := c.cc.Invoke(ctx, StatusService_Hello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatusServiceServer is the server API for StatusService service.
// All implementations must embed UnimplementedStatusServiceServer
// for forward compatibility.
//
// * 状态服务
type StatusServiceServer interface {
	// 发送心跳
	Ping(context.Context, *Empty) (*String, error)
	// 打招呼
	Hello(context.Context, *String) (*String, error)
	mustEmbedUnimplementedStatusServiceServer()
}

// UnimplementedStatusServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStatusServiceServer struct{}

func (UnimplementedStatusServiceServer) Ping(context.Context, *Empty) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedStatusServiceServer) Hello(context.Context, *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedStatusServiceServer) mustEmbedUnimplementedStatusServiceServer() {}
func (UnimplementedStatusServiceServer) testEmbeddedByValue()                       {}

// UnsafeStatusServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatusServiceServer will
// result in compilation errors.
type UnsafeStatusServiceServer interface {
	mustEmbedUnimplementedStatusServiceServer()
}

func RegisterStatusServiceServer(s grpc.ServiceRegistrar, srv StatusServiceServer) {
	// If the following call pancis, it indicates UnimplementedStatusServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StatusService_ServiceDesc, srv)
}

func _StatusService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusService_Hello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusServiceServer).Hello(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

// StatusService_ServiceDesc is the grpc.ServiceDesc for StatusService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatusService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "StatusService",
	HandlerType: (*StatusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _StatusService_Ping_Handler,
		},
		{
			MethodName: "Hello",
			Handler:    _StatusService_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "status_service.proto",
}
