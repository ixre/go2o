// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.0
// source: wallet_service.proto

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
	WalletService_CreateWallet_FullMethodName        = "/WalletService/CreateWallet"
	WalletService_GetWalletId_FullMethodName         = "/WalletService/GetWalletId"
	WalletService_GetWallet_FullMethodName           = "/WalletService/GetWallet"
	WalletService_GetWalletByCode_FullMethodName     = "/WalletService/GetWalletByCode"
	WalletService_GetWalletLog_FullMethodName        = "/WalletService/GetWalletLog"
	WalletService_Adjust_FullMethodName              = "/WalletService/Adjust"
	WalletService_Discount_FullMethodName            = "/WalletService/Discount"
	WalletService_Freeze_FullMethodName              = "/WalletService/Freeze"
	WalletService_Unfreeze_FullMethodName            = "/WalletService/Unfreeze"
	WalletService_Charge_FullMethodName              = "/WalletService/Charge"
	WalletService_Transfer_FullMethodName            = "/WalletService/Transfer"
	WalletService_RequestWithdrawal_FullMethodName   = "/WalletService/RequestWithdrawal"
	WalletService_ReviewTakeOut_FullMethodName       = "/WalletService/ReviewTakeOut"
	WalletService_CompleteTransaction_FullMethodName = "/WalletService/CompleteTransaction"
	WalletService_PagingWalletLog_FullMethodName     = "/WalletService/PagingWalletLog"
)

// WalletServiceClient is the client API for WalletService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// * 钱包服务
type WalletServiceClient interface {
	// * 创建钱包，并返回钱包编号
	CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*Result, error)
	// * 获取钱包编号，如果钱包不存在，则返回0
	GetWalletId(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*Int64, error)
	// * 获取钱包账户,传入walletId
	GetWallet(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*SWallet, error)
	// * 获取钱包账户,传入walletCode
	GetWalletByCode(ctx context.Context, in *String, opts ...grpc.CallOption) (*SWallet, error)
	// * 获取钱包日志
	GetWalletLog(ctx context.Context, in *WalletLogIDRequest, opts ...grpc.CallOption) (*SWalletLog, error)
	// * 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
	Adjust(ctx context.Context, in *AdjustRequest, opts ...grpc.CallOption) (*Result, error)
	// * 支付抵扣,must是否必须大于0
	Discount(ctx context.Context, in *DiscountRequest, opts ...grpc.CallOption) (*TxResult, error)
	// * 冻结余额
	Freeze(ctx context.Context, in *FreezeRequest, opts ...grpc.CallOption) (*FreezeResponse, error)
	// * 解冻金额
	Unfreeze(ctx context.Context, in *UnfreezeRequest, opts ...grpc.CallOption) (*Result, error)
	// * 充值,kind: 业务类型
	Charge(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*Result, error)
	// * 转账,title如:转账给xxx， toTitle: 转账收款xxx
	Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*Result, error)
	// * 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,transactionFee为手续费
	RequestWithdrawal(ctx context.Context, in *RequestWithdrawalRequest, opts ...grpc.CallOption) (*Result, error)
	// * 确认提现
	ReviewTakeOut(ctx context.Context, in *ReviewWithdrawRequest, opts ...grpc.CallOption) (*Result, error)
	// * 完成提现
	CompleteTransaction(ctx context.Context, in *FinishWithdrawRequest, opts ...grpc.CallOption) (*Result, error)
	// * 获取分页钱包日志
	PagingWalletLog(ctx context.Context, in *PagingWalletLogRequest, opts ...grpc.CallOption) (*SPagingResult, error)
}

type walletServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletServiceClient(cc grpc.ClientConnInterface) WalletServiceClient {
	return &walletServiceClient{cc}
}

func (c *walletServiceClient) CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_CreateWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetWalletId(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*Int64, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Int64)
	err := c.cc.Invoke(ctx, WalletService_GetWalletId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetWallet(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*SWallet, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SWallet)
	err := c.cc.Invoke(ctx, WalletService_GetWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetWalletByCode(ctx context.Context, in *String, opts ...grpc.CallOption) (*SWallet, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SWallet)
	err := c.cc.Invoke(ctx, WalletService_GetWalletByCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetWalletLog(ctx context.Context, in *WalletLogIDRequest, opts ...grpc.CallOption) (*SWalletLog, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SWalletLog)
	err := c.cc.Invoke(ctx, WalletService_GetWalletLog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) Adjust(ctx context.Context, in *AdjustRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_Adjust_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) Discount(ctx context.Context, in *DiscountRequest, opts ...grpc.CallOption) (*TxResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TxResult)
	err := c.cc.Invoke(ctx, WalletService_Discount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) Freeze(ctx context.Context, in *FreezeRequest, opts ...grpc.CallOption) (*FreezeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FreezeResponse)
	err := c.cc.Invoke(ctx, WalletService_Freeze_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) Unfreeze(ctx context.Context, in *UnfreezeRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_Unfreeze_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) Charge(ctx context.Context, in *ChargeRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_Charge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_Transfer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) RequestWithdrawal(ctx context.Context, in *RequestWithdrawalRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_RequestWithdrawal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) ReviewTakeOut(ctx context.Context, in *ReviewWithdrawRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_ReviewTakeOut_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) CompleteTransaction(ctx context.Context, in *FinishWithdrawRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, WalletService_CompleteTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) PagingWalletLog(ctx context.Context, in *PagingWalletLogRequest, opts ...grpc.CallOption) (*SPagingResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SPagingResult)
	err := c.cc.Invoke(ctx, WalletService_PagingWalletLog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServiceServer is the server API for WalletService service.
// All implementations must embed UnimplementedWalletServiceServer
// for forward compatibility.
//
// * 钱包服务
type WalletServiceServer interface {
	// * 创建钱包，并返回钱包编号
	CreateWallet(context.Context, *CreateWalletRequest) (*Result, error)
	// * 获取钱包编号，如果钱包不存在，则返回0
	GetWalletId(context.Context, *GetWalletRequest) (*Int64, error)
	// * 获取钱包账户,传入walletId
	GetWallet(context.Context, *Int64) (*SWallet, error)
	// * 获取钱包账户,传入walletCode
	GetWalletByCode(context.Context, *String) (*SWallet, error)
	// * 获取钱包日志
	GetWalletLog(context.Context, *WalletLogIDRequest) (*SWalletLog, error)
	// * 调整余额，可能存在扣为负数的情况，需传入操作人员编号或操作人员名称
	Adjust(context.Context, *AdjustRequest) (*Result, error)
	// * 支付抵扣,must是否必须大于0
	Discount(context.Context, *DiscountRequest) (*TxResult, error)
	// * 冻结余额
	Freeze(context.Context, *FreezeRequest) (*FreezeResponse, error)
	// * 解冻金额
	Unfreeze(context.Context, *UnfreezeRequest) (*Result, error)
	// * 充值,kind: 业务类型
	Charge(context.Context, *ChargeRequest) (*Result, error)
	// * 转账,title如:转账给xxx， toTitle: 转账收款xxx
	Transfer(context.Context, *TransferRequest) (*Result, error)
	// * 申请提现,kind：提现方式,返回info_id,交易号 及错误,value为提现金额,transactionFee为手续费
	RequestWithdrawal(context.Context, *RequestWithdrawalRequest) (*Result, error)
	// * 确认提现
	ReviewTakeOut(context.Context, *ReviewWithdrawRequest) (*Result, error)
	// * 完成提现
	CompleteTransaction(context.Context, *FinishWithdrawRequest) (*Result, error)
	// * 获取分页钱包日志
	PagingWalletLog(context.Context, *PagingWalletLogRequest) (*SPagingResult, error)
	mustEmbedUnimplementedWalletServiceServer()
}

// UnimplementedWalletServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedWalletServiceServer struct{}

func (UnimplementedWalletServiceServer) CreateWallet(context.Context, *CreateWalletRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWallet not implemented")
}
func (UnimplementedWalletServiceServer) GetWalletId(context.Context, *GetWalletRequest) (*Int64, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletId not implemented")
}
func (UnimplementedWalletServiceServer) GetWallet(context.Context, *Int64) (*SWallet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWallet not implemented")
}
func (UnimplementedWalletServiceServer) GetWalletByCode(context.Context, *String) (*SWallet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletByCode not implemented")
}
func (UnimplementedWalletServiceServer) GetWalletLog(context.Context, *WalletLogIDRequest) (*SWalletLog, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletLog not implemented")
}
func (UnimplementedWalletServiceServer) Adjust(context.Context, *AdjustRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Adjust not implemented")
}
func (UnimplementedWalletServiceServer) Discount(context.Context, *DiscountRequest) (*TxResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discount not implemented")
}
func (UnimplementedWalletServiceServer) Freeze(context.Context, *FreezeRequest) (*FreezeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Freeze not implemented")
}
func (UnimplementedWalletServiceServer) Unfreeze(context.Context, *UnfreezeRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unfreeze not implemented")
}
func (UnimplementedWalletServiceServer) Charge(context.Context, *ChargeRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Charge not implemented")
}
func (UnimplementedWalletServiceServer) Transfer(context.Context, *TransferRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
}
func (UnimplementedWalletServiceServer) RequestWithdrawal(context.Context, *RequestWithdrawalRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestWithdrawal not implemented")
}
func (UnimplementedWalletServiceServer) ReviewTakeOut(context.Context, *ReviewWithdrawRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReviewTakeOut not implemented")
}
func (UnimplementedWalletServiceServer) CompleteTransaction(context.Context, *FinishWithdrawRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteTransaction not implemented")
}
func (UnimplementedWalletServiceServer) PagingWalletLog(context.Context, *PagingWalletLogRequest) (*SPagingResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PagingWalletLog not implemented")
}
func (UnimplementedWalletServiceServer) mustEmbedUnimplementedWalletServiceServer() {}
func (UnimplementedWalletServiceServer) testEmbeddedByValue()                       {}

// UnsafeWalletServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WalletServiceServer will
// result in compilation errors.
type UnsafeWalletServiceServer interface {
	mustEmbedUnimplementedWalletServiceServer()
}

func RegisterWalletServiceServer(s grpc.ServiceRegistrar, srv WalletServiceServer) {
	// If the following call pancis, it indicates UnimplementedWalletServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&WalletService_ServiceDesc, srv)
}

func _WalletService_CreateWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).CreateWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_CreateWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).CreateWallet(ctx, req.(*CreateWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetWalletId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWalletId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_GetWalletId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWalletId(ctx, req.(*GetWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int64)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_GetWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWallet(ctx, req.(*Int64))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetWalletByCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWalletByCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_GetWalletByCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWalletByCode(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetWalletLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalletLogIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWalletLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_GetWalletLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWalletLog(ctx, req.(*WalletLogIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_Adjust_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdjustRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Adjust(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_Adjust_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Adjust(ctx, req.(*AdjustRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_Discount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Discount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_Discount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Discount(ctx, req.(*DiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_Freeze_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FreezeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Freeze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_Freeze_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Freeze(ctx, req.(*FreezeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_Unfreeze_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnfreezeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Unfreeze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_Unfreeze_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Unfreeze(ctx, req.(*UnfreezeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_Charge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChargeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Charge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_Charge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Charge(ctx, req.(*ChargeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_Transfer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).Transfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_RequestWithdrawal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithdrawalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).RequestWithdrawal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_RequestWithdrawal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).RequestWithdrawal(ctx, req.(*RequestWithdrawalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_ReviewTakeOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReviewWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).ReviewTakeOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_ReviewTakeOut_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).ReviewTakeOut(ctx, req.(*ReviewWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_CompleteTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishWithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).CompleteTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_CompleteTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).CompleteTransaction(ctx, req.(*FinishWithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_PagingWalletLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PagingWalletLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).PagingWalletLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_PagingWalletLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).PagingWalletLog(ctx, req.(*PagingWalletLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WalletService_ServiceDesc is the grpc.ServiceDesc for WalletService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WalletService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "WalletService",
	HandlerType: (*WalletServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWallet",
			Handler:    _WalletService_CreateWallet_Handler,
		},
		{
			MethodName: "GetWalletId",
			Handler:    _WalletService_GetWalletId_Handler,
		},
		{
			MethodName: "GetWallet",
			Handler:    _WalletService_GetWallet_Handler,
		},
		{
			MethodName: "GetWalletByCode",
			Handler:    _WalletService_GetWalletByCode_Handler,
		},
		{
			MethodName: "GetWalletLog",
			Handler:    _WalletService_GetWalletLog_Handler,
		},
		{
			MethodName: "Adjust",
			Handler:    _WalletService_Adjust_Handler,
		},
		{
			MethodName: "Discount",
			Handler:    _WalletService_Discount_Handler,
		},
		{
			MethodName: "Freeze",
			Handler:    _WalletService_Freeze_Handler,
		},
		{
			MethodName: "Unfreeze",
			Handler:    _WalletService_Unfreeze_Handler,
		},
		{
			MethodName: "Charge",
			Handler:    _WalletService_Charge_Handler,
		},
		{
			MethodName: "Transfer",
			Handler:    _WalletService_Transfer_Handler,
		},
		{
			MethodName: "RequestWithdrawal",
			Handler:    _WalletService_RequestWithdrawal_Handler,
		},
		{
			MethodName: "ReviewTakeOut",
			Handler:    _WalletService_ReviewTakeOut_Handler,
		},
		{
			MethodName: "CompleteTransaction",
			Handler:    _WalletService_CompleteTransaction_Handler,
		},
		{
			MethodName: "PagingWalletLog",
			Handler:    _WalletService_PagingWalletLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet_service.proto",
}
