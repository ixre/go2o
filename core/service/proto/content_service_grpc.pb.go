// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.0
// source: content_service.proto

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
	ContentService_GetPage_FullMethodName               = "/ContentService/getPage"
	ContentService_SavePage_FullMethodName              = "/ContentService/savePage"
	ContentService_DeletePage_FullMethodName            = "/ContentService/deletePage"
	ContentService_GetArticleCategories_FullMethodName  = "/ContentService/getArticleCategories"
	ContentService_GetArticleCategory_FullMethodName    = "/ContentService/getArticleCategory"
	ContentService_SaveArticleCategory_FullMethodName   = "/ContentService/saveArticleCategory"
	ContentService_DeleteArticleCategory_FullMethodName = "/ContentService/deleteArticleCategory"
	ContentService_GetArticle_FullMethodName            = "/ContentService/getArticle"
	ContentService_AddArticleViewsCount_FullMethodName  = "/ContentService/addArticleViewsCount"
	ContentService_LikeArticle_FullMethodName           = "/ContentService/likeArticle"
	ContentService_DeleteArticle_FullMethodName         = "/ContentService/deleteArticle"
	ContentService_SaveArticle_FullMethodName           = "/ContentService/saveArticle"
)

// ContentServiceClient is the client API for ContentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// * 内容服务
type ContentServiceClient interface {
	// 获取页面
	GetPage(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SPage, error)
	// 保存页面
	SavePage(ctx context.Context, in *SPage, opts ...grpc.CallOption) (*TxResult, error)
	// 删除页面
	DeletePage(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*TxResult, error)
	// 获取所有栏目
	GetArticleCategories(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ArticleCategoriesResponse, error)
	// 获取文章栏目,可传入ID或者别名
	GetArticleCategory(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SArticleCategory, error)
	// 保存文章栏目
	SaveArticleCategory(ctx context.Context, in *SArticleCategory, opts ...grpc.CallOption) (*TxResult, error)
	// 删除文章分类
	DeleteArticleCategory(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*TxResult, error)
	// 获取文章
	GetArticle(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SArticle, error)
	// 更新文章浏览次数
	AddArticleViewsCount(ctx context.Context, in *ArticleViewsRequest, opts ...grpc.CallOption) (*TxResult, error)
	// 喜欢/不喜欢文章
	LikeArticle(ctx context.Context, in *ArticleLikeRequest, opts ...grpc.CallOption) (*TxResult, error)
	// 删除文章
	DeleteArticle(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*TxResult, error)
	// 保存文章
	SaveArticle(ctx context.Context, in *SArticle, opts ...grpc.CallOption) (*TxResult, error)
}

type contentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewContentServiceClient(cc grpc.ClientConnInterface) ContentServiceClient {
	return &contentServiceClient{cc}
}

func (c *contentServiceClient) GetPage(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SPage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SPage)
	err := c.cc.Invoke(ctx, ContentService_GetPage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) SavePage(ctx context.Context, in *SPage, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_SavePage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) DeletePage(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_DeletePage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetArticleCategories(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ArticleCategoriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ArticleCategoriesResponse)
	err := c.cc.Invoke(ctx, ContentService_GetArticleCategories_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetArticleCategory(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SArticleCategory, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SArticleCategory)
	err := c.cc.Invoke(ctx, ContentService_GetArticleCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) SaveArticleCategory(ctx context.Context, in *SArticleCategory, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_SaveArticleCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) DeleteArticleCategory(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_DeleteArticleCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetArticle(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SArticle, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SArticle)
	err := c.cc.Invoke(ctx, ContentService_GetArticle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) AddArticleViewsCount(ctx context.Context, in *ArticleViewsRequest, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_AddArticleViewsCount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) LikeArticle(ctx context.Context, in *ArticleLikeRequest, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_LikeArticle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) DeleteArticle(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_DeleteArticle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) SaveArticle(ctx context.Context, in *SArticle, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, ContentService_SaveArticle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContentServiceServer is the server API for ContentService service.
// All implementations must embed UnimplementedContentServiceServer
// for forward compatibility.
//
// * 内容服务
type ContentServiceServer interface {
	// 获取页面
	GetPage(context.Context, *IdOrName) (*SPage, error)
	// 保存页面
	SavePage(context.Context, *SPage) (*TxResult, error)
	// 删除页面
	DeletePage(context.Context, *Int64) (*TxResult, error)
	// 获取所有栏目
	GetArticleCategories(context.Context, *Empty) (*ArticleCategoriesResponse, error)
	// 获取文章栏目,可传入ID或者别名
	GetArticleCategory(context.Context, *IdOrName) (*SArticleCategory, error)
	// 保存文章栏目
	SaveArticleCategory(context.Context, *SArticleCategory) (*TxResult, error)
	// 删除文章分类
	DeleteArticleCategory(context.Context, *Int64) (*TxResult, error)
	// 获取文章
	GetArticle(context.Context, *IdOrName) (*SArticle, error)
	// 更新文章浏览次数
	AddArticleViewsCount(context.Context, *ArticleViewsRequest) (*TxResult, error)
	// 喜欢/不喜欢文章
	LikeArticle(context.Context, *ArticleLikeRequest) (*TxResult, error)
	// 删除文章
	DeleteArticle(context.Context, *Int64) (*TxResult, error)
	// 保存文章
	SaveArticle(context.Context, *SArticle) (*TxResult, error)
	mustEmbedUnimplementedContentServiceServer()
}

// UnimplementedContentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedContentServiceServer struct{}

func (UnimplementedContentServiceServer) GetPage(context.Context, *IdOrName) (*SPage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPage not implemented")
}
func (UnimplementedContentServiceServer) SavePage(context.Context, *SPage) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SavePage not implemented")
}
func (UnimplementedContentServiceServer) DeletePage(context.Context, *Int64) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePage not implemented")
}
func (UnimplementedContentServiceServer) GetArticleCategories(context.Context, *Empty) (*ArticleCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleCategories not implemented")
}
func (UnimplementedContentServiceServer) GetArticleCategory(context.Context, *IdOrName) (*SArticleCategory, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleCategory not implemented")
}
func (UnimplementedContentServiceServer) SaveArticleCategory(context.Context, *SArticleCategory) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveArticleCategory not implemented")
}
func (UnimplementedContentServiceServer) DeleteArticleCategory(context.Context, *Int64) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArticleCategory not implemented")
}
func (UnimplementedContentServiceServer) GetArticle(context.Context, *IdOrName) (*SArticle, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticle not implemented")
}
func (UnimplementedContentServiceServer) AddArticleViewsCount(context.Context, *ArticleViewsRequest) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddArticleViewsCount not implemented")
}
func (UnimplementedContentServiceServer) LikeArticle(context.Context, *ArticleLikeRequest) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikeArticle not implemented")
}
func (UnimplementedContentServiceServer) DeleteArticle(context.Context, *Int64) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteArticle not implemented")
}
func (UnimplementedContentServiceServer) SaveArticle(context.Context, *SArticle) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveArticle not implemented")
}
func (UnimplementedContentServiceServer) mustEmbedUnimplementedContentServiceServer() {}
func (UnimplementedContentServiceServer) testEmbeddedByValue()                        {}

// UnsafeContentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContentServiceServer will
// result in compilation errors.
type UnsafeContentServiceServer interface {
	mustEmbedUnimplementedContentServiceServer()
}

func RegisterContentServiceServer(s grpc.ServiceRegistrar, srv ContentServiceServer) {
	// If the following call pancis, it indicates UnimplementedContentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ContentService_ServiceDesc, srv)
}

func _ContentService_GetPage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdOrName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetPage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetPage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetPage(ctx, req.(*IdOrName))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_SavePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SPage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).SavePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_SavePage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).SavePage(ctx, req.(*SPage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_DeletePage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int64)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).DeletePage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_DeletePage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).DeletePage(ctx, req.(*Int64))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetArticleCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetArticleCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetArticleCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetArticleCategories(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetArticleCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdOrName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetArticleCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetArticleCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetArticleCategory(ctx, req.(*IdOrName))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_SaveArticleCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SArticleCategory)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).SaveArticleCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_SaveArticleCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).SaveArticleCategory(ctx, req.(*SArticleCategory))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_DeleteArticleCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int64)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).DeleteArticleCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_DeleteArticleCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).DeleteArticleCategory(ctx, req.(*Int64))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdOrName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetArticle(ctx, req.(*IdOrName))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_AddArticleViewsCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleViewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).AddArticleViewsCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_AddArticleViewsCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).AddArticleViewsCount(ctx, req.(*ArticleViewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_LikeArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).LikeArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_LikeArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).LikeArticle(ctx, req.(*ArticleLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_DeleteArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int64)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).DeleteArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_DeleteArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).DeleteArticle(ctx, req.(*Int64))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_SaveArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SArticle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).SaveArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_SaveArticle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).SaveArticle(ctx, req.(*SArticle))
	}
	return interceptor(ctx, in, info, handler)
}

// ContentService_ServiceDesc is the grpc.ServiceDesc for ContentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ContentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ContentService",
	HandlerType: (*ContentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getPage",
			Handler:    _ContentService_GetPage_Handler,
		},
		{
			MethodName: "savePage",
			Handler:    _ContentService_SavePage_Handler,
		},
		{
			MethodName: "deletePage",
			Handler:    _ContentService_DeletePage_Handler,
		},
		{
			MethodName: "getArticleCategories",
			Handler:    _ContentService_GetArticleCategories_Handler,
		},
		{
			MethodName: "getArticleCategory",
			Handler:    _ContentService_GetArticleCategory_Handler,
		},
		{
			MethodName: "saveArticleCategory",
			Handler:    _ContentService_SaveArticleCategory_Handler,
		},
		{
			MethodName: "deleteArticleCategory",
			Handler:    _ContentService_DeleteArticleCategory_Handler,
		},
		{
			MethodName: "getArticle",
			Handler:    _ContentService_GetArticle_Handler,
		},
		{
			MethodName: "addArticleViewsCount",
			Handler:    _ContentService_AddArticleViewsCount_Handler,
		},
		{
			MethodName: "likeArticle",
			Handler:    _ContentService_LikeArticle_Handler,
		},
		{
			MethodName: "deleteArticle",
			Handler:    _ContentService_DeleteArticle_Handler,
		},
		{
			MethodName: "saveArticle",
			Handler:    _ContentService_SaveArticle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "content_service.proto",
}
