// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.0
// source: item_service.proto

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
	ItemService_GetItem_FullMethodName                        = "/ItemService/GetItem"
	ItemService_SaveItem_FullMethodName                       = "/ItemService/SaveItem"
	ItemService_GetItemBySku_FullMethodName                   = "/ItemService/GetItemBySku"
	ItemService_GetItemAndSnapshot_FullMethodName             = "/ItemService/GetItemAndSnapshot"
	ItemService_GetTradeSnapshot_FullMethodName               = "/ItemService/GetTradeSnapshot"
	ItemService_GetSku_FullMethodName                         = "/ItemService/GetSku"
	ItemService_ReviewItem_FullMethodName                     = "/ItemService/ReviewItem"
	ItemService_RecycleItem_FullMethodName                    = "/ItemService/RecycleItem"
	ItemService_SaveLevelPrices_FullMethodName                = "/ItemService/SaveLevelPrices"
	ItemService_SignAsIllegal_FullMethodName                  = "/ItemService/SignAsIllegal"
	ItemService_SetShelveState_FullMethodName                 = "/ItemService/SetShelveState"
	ItemService_GetItemDetailData_FullMethodName              = "/ItemService/GetItemDetailData"
	ItemService_GetItems_FullMethodName                       = "/ItemService/GetItems"
	ItemService_GetWholesalePriceArray_FullMethodName         = "/ItemService/GetWholesalePriceArray"
	ItemService_SaveWholesalePrice_FullMethodName             = "/ItemService/SaveWholesalePrice"
	ItemService_GetWholesaleDiscountArray_FullMethodName      = "/ItemService/GetWholesaleDiscountArray"
	ItemService_SaveWholesaleDiscount_FullMethodName          = "/ItemService/SaveWholesaleDiscount"
	ItemService_GetAllSaleLabels_FullMethodName               = "/ItemService/GetAllSaleLabels"
	ItemService_GetSaleLabel_FullMethodName                   = "/ItemService/GetSaleLabel"
	ItemService_SaveSaleLabel_FullMethodName                  = "/ItemService/SaveSaleLabel"
	ItemService_DeleteSaleLabel_FullMethodName                = "/ItemService/DeleteSaleLabel"
	ItemService_GetPagedValueGoodsBySaleLabel__FullMethodName = "/ItemService/GetPagedValueGoodsBySaleLabel_"
)

// ItemServiceClient is the client API for ItemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 商品服务
type ItemServiceClient interface {
	// 获取商品的数据
	GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*SItemDataResponse, error)
	// 保存商品的数据
	SaveItem(ctx context.Context, in *SaveItemRequest, opts ...grpc.CallOption) (*SaveItemResponse, error)
	// 根据SKU获取商品
	GetItemBySku(ctx context.Context, in *ItemBySkuRequest, opts ...grpc.CallOption) (*SUnifiedViewItem, error)
	// 获取商品用于销售的快照和信息
	GetItemAndSnapshot(ctx context.Context, in *GetItemAndSnapshotRequest, opts ...grpc.CallOption) (*ItemSnapshotResponse, error)
	// 获取商品交易快照,参数传递:snapshotId
	GetTradeSnapshot(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*STradeSnapshot, error)
	// 获取SKU
	GetSku(ctx context.Context, in *SkuId, opts ...grpc.CallOption) (*SSku, error)
	// 设置商品货架状态
	ReviewItem(ctx context.Context, in *ItemReviewRequest, opts ...grpc.CallOption) (*Result, error)
	// 回收商品
	RecycleItem(ctx context.Context, in *RecycleItemRequest, opts ...grpc.CallOption) (*Result, error)
	// 保存商品的会员价
	SaveLevelPrices(ctx context.Context, in *SaveLevelPriceRequest, opts ...grpc.CallOption) (*Result, error)
	// 商品标记为违规
	SignAsIllegal(ctx context.Context, in *ItemIllegalRequest, opts ...grpc.CallOption) (*Result, error)
	// 设置商品货架状态
	SetShelveState(ctx context.Context, in *ShelveStateRequest, opts ...grpc.CallOption) (*Result, error)
	// 获取商品详细数据
	GetItemDetailData(ctx context.Context, in *ItemDetailRequest, opts ...grpc.CallOption) (*String, error)
	// 获取上架商品数据
	GetItems(ctx context.Context, in *GetItemsRequest, opts ...grpc.CallOption) (*PagingGoodsResponse, error)
	// 获取批发价格数组
	GetWholesalePriceArray(ctx context.Context, in *SkuId, opts ...grpc.CallOption) (*SWsSkuPriceListResponse, error)
	// 保存批发价格
	SaveWholesalePrice(ctx context.Context, in *SaveSkuPricesRequest, opts ...grpc.CallOption) (*Result, error)
	// 获取批发折扣数组
	GetWholesaleDiscountArray(ctx context.Context, in *GetWsDiscountRequest, opts ...grpc.CallOption) (*SWsItemDiscountListResponse, error)
	// 保存批发折扣
	SaveWholesaleDiscount(ctx context.Context, in *SaveItemDiscountRequest, opts ...grpc.CallOption) (*Result, error)
	// 获取所有的商品标签
	GetAllSaleLabels(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ItemLabelListResponse, error)
	// 获取销售标签
	GetSaleLabel(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SItemLabel, error)
	// 保存销售标签
	SaveSaleLabel(ctx context.Context, in *SItemLabel, opts ...grpc.CallOption) (*Result, error)
	// 删除销售标签
	DeleteSaleLabel(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*Result, error)
	// 根据分页销售标签获取指定数目的商品
	GetPagedValueGoodsBySaleLabel_(ctx context.Context, in *SaleLabelItemsRequest_, opts ...grpc.CallOption) (*PagingGoodsResponse, error)
}

type itemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewItemServiceClient(cc grpc.ClientConnInterface) ItemServiceClient {
	return &itemServiceClient{cc}
}

func (c *itemServiceClient) GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*SItemDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SItemDataResponse)
	err := c.cc.Invoke(ctx, ItemService_GetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SaveItem(ctx context.Context, in *SaveItemRequest, opts ...grpc.CallOption) (*SaveItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveItemResponse)
	err := c.cc.Invoke(ctx, ItemService_SaveItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetItemBySku(ctx context.Context, in *ItemBySkuRequest, opts ...grpc.CallOption) (*SUnifiedViewItem, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SUnifiedViewItem)
	err := c.cc.Invoke(ctx, ItemService_GetItemBySku_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetItemAndSnapshot(ctx context.Context, in *GetItemAndSnapshotRequest, opts ...grpc.CallOption) (*ItemSnapshotResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ItemSnapshotResponse)
	err := c.cc.Invoke(ctx, ItemService_GetItemAndSnapshot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetTradeSnapshot(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*STradeSnapshot, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(STradeSnapshot)
	err := c.cc.Invoke(ctx, ItemService_GetTradeSnapshot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetSku(ctx context.Context, in *SkuId, opts ...grpc.CallOption) (*SSku, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SSku)
	err := c.cc.Invoke(ctx, ItemService_GetSku_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) ReviewItem(ctx context.Context, in *ItemReviewRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_ReviewItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) RecycleItem(ctx context.Context, in *RecycleItemRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_RecycleItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SaveLevelPrices(ctx context.Context, in *SaveLevelPriceRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_SaveLevelPrices_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SignAsIllegal(ctx context.Context, in *ItemIllegalRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_SignAsIllegal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SetShelveState(ctx context.Context, in *ShelveStateRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_SetShelveState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetItemDetailData(ctx context.Context, in *ItemDetailRequest, opts ...grpc.CallOption) (*String, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(String)
	err := c.cc.Invoke(ctx, ItemService_GetItemDetailData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetItems(ctx context.Context, in *GetItemsRequest, opts ...grpc.CallOption) (*PagingGoodsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PagingGoodsResponse)
	err := c.cc.Invoke(ctx, ItemService_GetItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetWholesalePriceArray(ctx context.Context, in *SkuId, opts ...grpc.CallOption) (*SWsSkuPriceListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SWsSkuPriceListResponse)
	err := c.cc.Invoke(ctx, ItemService_GetWholesalePriceArray_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SaveWholesalePrice(ctx context.Context, in *SaveSkuPricesRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_SaveWholesalePrice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetWholesaleDiscountArray(ctx context.Context, in *GetWsDiscountRequest, opts ...grpc.CallOption) (*SWsItemDiscountListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SWsItemDiscountListResponse)
	err := c.cc.Invoke(ctx, ItemService_GetWholesaleDiscountArray_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SaveWholesaleDiscount(ctx context.Context, in *SaveItemDiscountRequest, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_SaveWholesaleDiscount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetAllSaleLabels(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ItemLabelListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ItemLabelListResponse)
	err := c.cc.Invoke(ctx, ItemService_GetAllSaleLabels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetSaleLabel(ctx context.Context, in *IdOrName, opts ...grpc.CallOption) (*SItemLabel, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SItemLabel)
	err := c.cc.Invoke(ctx, ItemService_GetSaleLabel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) SaveSaleLabel(ctx context.Context, in *SItemLabel, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_SaveSaleLabel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) DeleteSaleLabel(ctx context.Context, in *Int64, opts ...grpc.CallOption) (*Result, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Result)
	err := c.cc.Invoke(ctx, ItemService_DeleteSaleLabel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemServiceClient) GetPagedValueGoodsBySaleLabel_(ctx context.Context, in *SaleLabelItemsRequest_, opts ...grpc.CallOption) (*PagingGoodsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PagingGoodsResponse)
	err := c.cc.Invoke(ctx, ItemService_GetPagedValueGoodsBySaleLabel__FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemServiceServer is the server API for ItemService service.
// All implementations must embed UnimplementedItemServiceServer
// for forward compatibility.
//
// 商品服务
type ItemServiceServer interface {
	// 获取商品的数据
	GetItem(context.Context, *GetItemRequest) (*SItemDataResponse, error)
	// 保存商品的数据
	SaveItem(context.Context, *SaveItemRequest) (*SaveItemResponse, error)
	// 根据SKU获取商品
	GetItemBySku(context.Context, *ItemBySkuRequest) (*SUnifiedViewItem, error)
	// 获取商品用于销售的快照和信息
	GetItemAndSnapshot(context.Context, *GetItemAndSnapshotRequest) (*ItemSnapshotResponse, error)
	// 获取商品交易快照,参数传递:snapshotId
	GetTradeSnapshot(context.Context, *Int64) (*STradeSnapshot, error)
	// 获取SKU
	GetSku(context.Context, *SkuId) (*SSku, error)
	// 设置商品货架状态
	ReviewItem(context.Context, *ItemReviewRequest) (*Result, error)
	// 回收商品
	RecycleItem(context.Context, *RecycleItemRequest) (*Result, error)
	// 保存商品的会员价
	SaveLevelPrices(context.Context, *SaveLevelPriceRequest) (*Result, error)
	// 商品标记为违规
	SignAsIllegal(context.Context, *ItemIllegalRequest) (*Result, error)
	// 设置商品货架状态
	SetShelveState(context.Context, *ShelveStateRequest) (*Result, error)
	// 获取商品详细数据
	GetItemDetailData(context.Context, *ItemDetailRequest) (*String, error)
	// 获取上架商品数据
	GetItems(context.Context, *GetItemsRequest) (*PagingGoodsResponse, error)
	// 获取批发价格数组
	GetWholesalePriceArray(context.Context, *SkuId) (*SWsSkuPriceListResponse, error)
	// 保存批发价格
	SaveWholesalePrice(context.Context, *SaveSkuPricesRequest) (*Result, error)
	// 获取批发折扣数组
	GetWholesaleDiscountArray(context.Context, *GetWsDiscountRequest) (*SWsItemDiscountListResponse, error)
	// 保存批发折扣
	SaveWholesaleDiscount(context.Context, *SaveItemDiscountRequest) (*Result, error)
	// 获取所有的商品标签
	GetAllSaleLabels(context.Context, *Empty) (*ItemLabelListResponse, error)
	// 获取销售标签
	GetSaleLabel(context.Context, *IdOrName) (*SItemLabel, error)
	// 保存销售标签
	SaveSaleLabel(context.Context, *SItemLabel) (*Result, error)
	// 删除销售标签
	DeleteSaleLabel(context.Context, *Int64) (*Result, error)
	// 根据分页销售标签获取指定数目的商品
	GetPagedValueGoodsBySaleLabel_(context.Context, *SaleLabelItemsRequest_) (*PagingGoodsResponse, error)
	mustEmbedUnimplementedItemServiceServer()
}

// UnimplementedItemServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedItemServiceServer struct{}

func (UnimplementedItemServiceServer) GetItem(context.Context, *GetItemRequest) (*SItemDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItem not implemented")
}
func (UnimplementedItemServiceServer) SaveItem(context.Context, *SaveItemRequest) (*SaveItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveItem not implemented")
}
func (UnimplementedItemServiceServer) GetItemBySku(context.Context, *ItemBySkuRequest) (*SUnifiedViewItem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemBySku not implemented")
}
func (UnimplementedItemServiceServer) GetItemAndSnapshot(context.Context, *GetItemAndSnapshotRequest) (*ItemSnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemAndSnapshot not implemented")
}
func (UnimplementedItemServiceServer) GetTradeSnapshot(context.Context, *Int64) (*STradeSnapshot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTradeSnapshot not implemented")
}
func (UnimplementedItemServiceServer) GetSku(context.Context, *SkuId) (*SSku, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSku not implemented")
}
func (UnimplementedItemServiceServer) ReviewItem(context.Context, *ItemReviewRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReviewItem not implemented")
}
func (UnimplementedItemServiceServer) RecycleItem(context.Context, *RecycleItemRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecycleItem not implemented")
}
func (UnimplementedItemServiceServer) SaveLevelPrices(context.Context, *SaveLevelPriceRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveLevelPrices not implemented")
}
func (UnimplementedItemServiceServer) SignAsIllegal(context.Context, *ItemIllegalRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignAsIllegal not implemented")
}
func (UnimplementedItemServiceServer) SetShelveState(context.Context, *ShelveStateRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetShelveState not implemented")
}
func (UnimplementedItemServiceServer) GetItemDetailData(context.Context, *ItemDetailRequest) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemDetailData not implemented")
}
func (UnimplementedItemServiceServer) GetItems(context.Context, *GetItemsRequest) (*PagingGoodsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItems not implemented")
}
func (UnimplementedItemServiceServer) GetWholesalePriceArray(context.Context, *SkuId) (*SWsSkuPriceListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWholesalePriceArray not implemented")
}
func (UnimplementedItemServiceServer) SaveWholesalePrice(context.Context, *SaveSkuPricesRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveWholesalePrice not implemented")
}
func (UnimplementedItemServiceServer) GetWholesaleDiscountArray(context.Context, *GetWsDiscountRequest) (*SWsItemDiscountListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWholesaleDiscountArray not implemented")
}
func (UnimplementedItemServiceServer) SaveWholesaleDiscount(context.Context, *SaveItemDiscountRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveWholesaleDiscount not implemented")
}
func (UnimplementedItemServiceServer) GetAllSaleLabels(context.Context, *Empty) (*ItemLabelListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllSaleLabels not implemented")
}
func (UnimplementedItemServiceServer) GetSaleLabel(context.Context, *IdOrName) (*SItemLabel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSaleLabel not implemented")
}
func (UnimplementedItemServiceServer) SaveSaleLabel(context.Context, *SItemLabel) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveSaleLabel not implemented")
}
func (UnimplementedItemServiceServer) DeleteSaleLabel(context.Context, *Int64) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSaleLabel not implemented")
}
func (UnimplementedItemServiceServer) GetPagedValueGoodsBySaleLabel_(context.Context, *SaleLabelItemsRequest_) (*PagingGoodsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPagedValueGoodsBySaleLabel_ not implemented")
}
func (UnimplementedItemServiceServer) mustEmbedUnimplementedItemServiceServer() {}
func (UnimplementedItemServiceServer) testEmbeddedByValue()                     {}

// UnsafeItemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ItemServiceServer will
// result in compilation errors.
type UnsafeItemServiceServer interface {
	mustEmbedUnimplementedItemServiceServer()
}

func RegisterItemServiceServer(s grpc.ServiceRegistrar, srv ItemServiceServer) {
	// If the following call pancis, it indicates UnimplementedItemServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ItemService_ServiceDesc, srv)
}

func _ItemService_GetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetItem(ctx, req.(*GetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SaveItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SaveItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SaveItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SaveItem(ctx, req.(*SaveItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetItemBySku_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemBySkuRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetItemBySku(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetItemBySku_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetItemBySku(ctx, req.(*ItemBySkuRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetItemAndSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemAndSnapshotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetItemAndSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetItemAndSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetItemAndSnapshot(ctx, req.(*GetItemAndSnapshotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetTradeSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int64)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetTradeSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetTradeSnapshot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetTradeSnapshot(ctx, req.(*Int64))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetSku_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SkuId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetSku(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetSku_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetSku(ctx, req.(*SkuId))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_ReviewItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).ReviewItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_ReviewItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).ReviewItem(ctx, req.(*ItemReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_RecycleItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecycleItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).RecycleItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_RecycleItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).RecycleItem(ctx, req.(*RecycleItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SaveLevelPrices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveLevelPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SaveLevelPrices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SaveLevelPrices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SaveLevelPrices(ctx, req.(*SaveLevelPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SignAsIllegal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemIllegalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SignAsIllegal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SignAsIllegal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SignAsIllegal(ctx, req.(*ItemIllegalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SetShelveState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShelveStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SetShelveState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SetShelveState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SetShelveState(ctx, req.(*ShelveStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetItemDetailData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetItemDetailData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetItemDetailData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetItemDetailData(ctx, req.(*ItemDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetItems(ctx, req.(*GetItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetWholesalePriceArray_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SkuId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetWholesalePriceArray(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetWholesalePriceArray_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetWholesalePriceArray(ctx, req.(*SkuId))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SaveWholesalePrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveSkuPricesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SaveWholesalePrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SaveWholesalePrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SaveWholesalePrice(ctx, req.(*SaveSkuPricesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetWholesaleDiscountArray_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWsDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetWholesaleDiscountArray(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetWholesaleDiscountArray_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetWholesaleDiscountArray(ctx, req.(*GetWsDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SaveWholesaleDiscount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveItemDiscountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SaveWholesaleDiscount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SaveWholesaleDiscount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SaveWholesaleDiscount(ctx, req.(*SaveItemDiscountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetAllSaleLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetAllSaleLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetAllSaleLabels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetAllSaleLabels(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetSaleLabel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdOrName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetSaleLabel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetSaleLabel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetSaleLabel(ctx, req.(*IdOrName))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_SaveSaleLabel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SItemLabel)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).SaveSaleLabel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_SaveSaleLabel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).SaveSaleLabel(ctx, req.(*SItemLabel))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_DeleteSaleLabel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Int64)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).DeleteSaleLabel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_DeleteSaleLabel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).DeleteSaleLabel(ctx, req.(*Int64))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemService_GetPagedValueGoodsBySaleLabel__Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaleLabelItemsRequest_)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemServiceServer).GetPagedValueGoodsBySaleLabel_(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemService_GetPagedValueGoodsBySaleLabel__FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemServiceServer).GetPagedValueGoodsBySaleLabel_(ctx, req.(*SaleLabelItemsRequest_))
	}
	return interceptor(ctx, in, info, handler)
}

// ItemService_ServiceDesc is the grpc.ServiceDesc for ItemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ItemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ItemService",
	HandlerType: (*ItemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetItem",
			Handler:    _ItemService_GetItem_Handler,
		},
		{
			MethodName: "SaveItem",
			Handler:    _ItemService_SaveItem_Handler,
		},
		{
			MethodName: "GetItemBySku",
			Handler:    _ItemService_GetItemBySku_Handler,
		},
		{
			MethodName: "GetItemAndSnapshot",
			Handler:    _ItemService_GetItemAndSnapshot_Handler,
		},
		{
			MethodName: "GetTradeSnapshot",
			Handler:    _ItemService_GetTradeSnapshot_Handler,
		},
		{
			MethodName: "GetSku",
			Handler:    _ItemService_GetSku_Handler,
		},
		{
			MethodName: "ReviewItem",
			Handler:    _ItemService_ReviewItem_Handler,
		},
		{
			MethodName: "RecycleItem",
			Handler:    _ItemService_RecycleItem_Handler,
		},
		{
			MethodName: "SaveLevelPrices",
			Handler:    _ItemService_SaveLevelPrices_Handler,
		},
		{
			MethodName: "SignAsIllegal",
			Handler:    _ItemService_SignAsIllegal_Handler,
		},
		{
			MethodName: "SetShelveState",
			Handler:    _ItemService_SetShelveState_Handler,
		},
		{
			MethodName: "GetItemDetailData",
			Handler:    _ItemService_GetItemDetailData_Handler,
		},
		{
			MethodName: "GetItems",
			Handler:    _ItemService_GetItems_Handler,
		},
		{
			MethodName: "GetWholesalePriceArray",
			Handler:    _ItemService_GetWholesalePriceArray_Handler,
		},
		{
			MethodName: "SaveWholesalePrice",
			Handler:    _ItemService_SaveWholesalePrice_Handler,
		},
		{
			MethodName: "GetWholesaleDiscountArray",
			Handler:    _ItemService_GetWholesaleDiscountArray_Handler,
		},
		{
			MethodName: "SaveWholesaleDiscount",
			Handler:    _ItemService_SaveWholesaleDiscount_Handler,
		},
		{
			MethodName: "GetAllSaleLabels",
			Handler:    _ItemService_GetAllSaleLabels_Handler,
		},
		{
			MethodName: "GetSaleLabel",
			Handler:    _ItemService_GetSaleLabel_Handler,
		},
		{
			MethodName: "SaveSaleLabel",
			Handler:    _ItemService_SaveSaleLabel_Handler,
		},
		{
			MethodName: "DeleteSaleLabel",
			Handler:    _ItemService_DeleteSaleLabel_Handler,
		},
		{
			MethodName: "GetPagedValueGoodsBySaleLabel_",
			Handler:    _ItemService_GetPagedValueGoodsBySaleLabel__Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "item_service.proto",
}
