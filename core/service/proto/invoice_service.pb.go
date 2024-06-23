//*
// This file is auto generated by tto v0.5.8 !
// If you want to modify this code, please read the guide
// to modify code template.
//
// Get started: https://github.com/ixre/tto
//
// Copyright (C) 2009-2024 56X.NET, All rights reserved.
//
// name : invoice_tenant_service.proto
// author : jarrysix
// date : 2024/06/24 15:22:54
// description :
// history :

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v5.27.0
// source: invoice_service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 保存发票租户请求
type SaveTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 租户类型,1:会员  2:商户
	TenantType int32 `protobuf:"varint,2,opt,name=tenantType,proto3" json:"tenantType"`
	// 租户用户编号
	TenantUid int64 `protobuf:"varint,3,opt,name=tenantUid,proto3" json:"tenantUid"`
	// 创建时间
	CreateTime int64 `protobuf:"varint,4,opt,name=createTime,proto3" json:"createTime"`
}

func (x *SaveTenantRequest) Reset() {
	*x = SaveTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveTenantRequest) ProtoMessage() {}

func (x *SaveTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveTenantRequest.ProtoReflect.Descriptor instead.
func (*SaveTenantRequest) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{0}
}

func (x *SaveTenantRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SaveTenantRequest) GetTenantType() int32 {
	if x != nil {
		return x.TenantType
	}
	return 0
}

func (x *SaveTenantRequest) GetTenantUid() int64 {
	if x != nil {
		return x.TenantUid
	}
	return 0
}

func (x *SaveTenantRequest) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

// 保存发票租户响应
type SaveTenantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode int32  `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode"`
	ErrMsg  string `protobuf:"bytes,2,opt,name=errMsg,proto3" json:"errMsg"`
	Id      int64  `protobuf:"varint,3,opt,name=id,proto3" json:"id"`
}

func (x *SaveTenantResponse) Reset() {
	*x = SaveTenantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveTenantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveTenantResponse) ProtoMessage() {}

func (x *SaveTenantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveTenantResponse.ProtoReflect.Descriptor instead.
func (*SaveTenantResponse) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{1}
}

func (x *SaveTenantResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *SaveTenantResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *SaveTenantResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// 发票租户编号
type InvoiceTenantId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value"`
}

func (x *InvoiceTenantId) Reset() {
	*x = InvoiceTenantId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvoiceTenantId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvoiceTenantId) ProtoMessage() {}

func (x *InvoiceTenantId) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvoiceTenantId.ProtoReflect.Descriptor instead.
func (*InvoiceTenantId) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{2}
}

func (x *InvoiceTenantId) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// 发票租户
type STenant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 租户类型,1:会员  2:商户
	TenantType int32 `protobuf:"varint,2,opt,name=tenantType,proto3" json:"tenantType"`
	// 租户用户编号
	TenantUid int64 `protobuf:"varint,3,opt,name=tenantUid,proto3" json:"tenantUid"`
	// 创建时间
	CreateTime int64 `protobuf:"varint,4,opt,name=createTime,proto3" json:"createTime"`
}

func (x *STenant) Reset() {
	*x = STenant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *STenant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*STenant) ProtoMessage() {}

func (x *STenant) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use STenant.ProtoReflect.Descriptor instead.
func (*STenant) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{3}
}

func (x *STenant) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *STenant) GetTenantType() int32 {
	if x != nil {
		return x.TenantType
	}
	return 0
}

func (x *STenant) GetTenantUid() int64 {
	if x != nil {
		return x.TenantUid
	}
	return 0
}

func (x *STenant) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

// 查询发票租户请求
type QueryTenantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *QueryTenantRequest) Reset() {
	*x = QueryTenantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryTenantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTenantRequest) ProtoMessage() {}

func (x *QueryTenantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTenantRequest.ProtoReflect.Descriptor instead.
func (*QueryTenantRequest) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{4}
}

// 查询发票租户响应
type QueryTenantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []*STenant `protobuf:"bytes,1,rep,name=value,proto3" json:"value"`
}

func (x *QueryTenantResponse) Reset() {
	*x = QueryTenantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryTenantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTenantResponse) ProtoMessage() {}

func (x *QueryTenantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTenantResponse.ProtoReflect.Descriptor instead.
func (*QueryTenantResponse) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{5}
}

func (x *QueryTenantResponse) GetValue() []*STenant {
	if x != nil {
		return x.Value
	}
	return nil
}

// 发票租户分页数据
type PagingTenant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 租户类型,1:会员  2:商户
	TenantType int32 `protobuf:"varint,2,opt,name=tenantType,proto3" json:"tenantType"`
	// 租户用户编号
	TenantUid int64 `protobuf:"varint,3,opt,name=tenantUid,proto3" json:"tenantUid"`
	// 创建时间
	CreateTime int64 `protobuf:"varint,4,opt,name=createTime,proto3" json:"createTime"`
}

func (x *PagingTenant) Reset() {
	*x = PagingTenant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PagingTenant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagingTenant) ProtoMessage() {}

func (x *PagingTenant) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagingTenant.ProtoReflect.Descriptor instead.
func (*PagingTenant) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{6}
}

func (x *PagingTenant) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PagingTenant) GetTenantType() int32 {
	if x != nil {
		return x.TenantType
	}
	return 0
}

func (x *PagingTenant) GetTenantUid() int64 {
	if x != nil {
		return x.TenantUid
	}
	return 0
}

func (x *PagingTenant) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

// 发票租户分页请求
type TenantPagingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 分页参数
	Params *SPagingParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (x *TenantPagingRequest) Reset() {
	*x = TenantPagingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantPagingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantPagingRequest) ProtoMessage() {}

func (x *TenantPagingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantPagingRequest.ProtoReflect.Descriptor instead.
func (*TenantPagingRequest) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{7}
}

func (x *TenantPagingRequest) GetParams() *SPagingParams {
	if x != nil {
		return x.Params
	}
	return nil
}

// 发票租户分页响应
type TenantPagingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 总数
	Total int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total"`
	// 数据列表
	Value []*PagingTenant `protobuf:"bytes,2,rep,name=value,proto3" json:"value"`
}

func (x *TenantPagingResponse) Reset() {
	*x = TenantPagingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invoice_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TenantPagingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TenantPagingResponse) ProtoMessage() {}

func (x *TenantPagingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_invoice_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TenantPagingResponse.ProtoReflect.Descriptor instead.
func (*TenantPagingResponse) Descriptor() ([]byte, []int) {
	return file_invoice_service_proto_rawDescGZIP(), []int{8}
}

func (x *TenantPagingResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *TenantPagingResponse) GetValue() []*PagingTenant {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_invoice_service_proto protoreflect.FileDescriptor

var file_invoice_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x69, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x01, 0x0a, 0x11, 0x53, 0x61, 0x76, 0x65, 0x54, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x55, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x55, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x56, 0x0a, 0x12, 0x53, 0x61, 0x76,
	0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72,
	0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73,
	0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x27, 0x0a, 0x0f, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x77, 0x0a, 0x07, 0x53, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x55,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x55, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x35, 0x0a, 0x13, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x1e, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x08, 0x2e, 0x53, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x7c, 0x0a, 0x0c, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x55, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x55, 0x69, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x3d,
	0x0a, 0x13, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x53, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x51, 0x0a,
	0x14, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x23, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x67, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x32, 0x9e, 0x02, 0x0a, 0x0e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x53, 0x61, 0x76, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x12, 0x12, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x54, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x09,
	0x47, 0x65, 0x74, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x10, 0x2e, 0x49, 0x6e, 0x76, 0x6f,
	0x69, 0x63, 0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0x08, 0x2e, 0x53, 0x54,
	0x65, 0x6e, 0x61, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0f, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x13, 0x2e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x12, 0x10, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63,
	0x65, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0c, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x54, 0x65,
	0x6e, 0x61, 0x6e, 0x74, 0x12, 0x14, 0x2e, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x50, 0x61, 0x67,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x54, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x36, 0x0a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x69, 0x78, 0x72, 0x65, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x2e, 0x6c, 0x73, 0x70, 0x2e, 0x72, 0x70, 0x63, 0x5a,
	0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_invoice_service_proto_rawDescOnce sync.Once
	file_invoice_service_proto_rawDescData = file_invoice_service_proto_rawDesc
)

func file_invoice_service_proto_rawDescGZIP() []byte {
	file_invoice_service_proto_rawDescOnce.Do(func() {
		file_invoice_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_invoice_service_proto_rawDescData)
	})
	return file_invoice_service_proto_rawDescData
}

var file_invoice_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_invoice_service_proto_goTypes = []interface{}{
	(*SaveTenantRequest)(nil),    // 0: SaveTenantRequest
	(*SaveTenantResponse)(nil),   // 1: SaveTenantResponse
	(*InvoiceTenantId)(nil),      // 2: InvoiceTenantId
	(*STenant)(nil),              // 3: STenant
	(*QueryTenantRequest)(nil),   // 4: QueryTenantRequest
	(*QueryTenantResponse)(nil),  // 5: QueryTenantResponse
	(*PagingTenant)(nil),         // 6: PagingTenant
	(*TenantPagingRequest)(nil),  // 7: TenantPagingRequest
	(*TenantPagingResponse)(nil), // 8: TenantPagingResponse
	(*SPagingParams)(nil),        // 9: SPagingParams
	(*Result)(nil),               // 10: Result
}
var file_invoice_service_proto_depIdxs = []int32{
	3,  // 0: QueryTenantResponse.value:type_name -> STenant
	9,  // 1: TenantPagingRequest.params:type_name -> SPagingParams
	6,  // 2: TenantPagingResponse.value:type_name -> PagingTenant
	0,  // 3: InvoiceService.SaveTenant:input_type -> SaveTenantRequest
	2,  // 4: InvoiceService.GetTenant:input_type -> InvoiceTenantId
	4,  // 5: InvoiceService.QueryTenantList:input_type -> QueryTenantRequest
	2,  // 6: InvoiceService.DeleteTenant:input_type -> InvoiceTenantId
	7,  // 7: InvoiceService.PagingTenant:input_type -> TenantPagingRequest
	1,  // 8: InvoiceService.SaveTenant:output_type -> SaveTenantResponse
	3,  // 9: InvoiceService.GetTenant:output_type -> STenant
	5,  // 10: InvoiceService.QueryTenantList:output_type -> QueryTenantResponse
	10, // 11: InvoiceService.DeleteTenant:output_type -> Result
	8,  // 12: InvoiceService.PagingTenant:output_type -> TenantPagingResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_invoice_service_proto_init() }
func file_invoice_service_proto_init() {
	if File_invoice_service_proto != nil {
		return
	}
	file_global_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_invoice_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveTenantRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveTenantResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvoiceTenantId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*STenant); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryTenantRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryTenantResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PagingTenant); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantPagingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invoice_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TenantPagingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_invoice_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_invoice_service_proto_goTypes,
		DependencyIndexes: file_invoice_service_proto_depIdxs,
		MessageInfos:      file_invoice_service_proto_msgTypes,
	}.Build()
	File_invoice_service_proto = out.File
	file_invoice_service_proto_rawDesc = nil
	file_invoice_service_proto_goTypes = nil
	file_invoice_service_proto_depIdxs = nil
}
