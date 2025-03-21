//*
// Copyright (C) 2007-2024 fze.NET, All rights reserved.
//
// name: provider_service.proto
// author: jarrysix (jarrysix#gmail.com)
// date: 2024-09-06 09:17:13
// description: 第三方服务接口
// history:

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0
// source: provider_service.proto

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

// 获取用户OpenId请求
type GetUserOpenIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 应用ID
	AppId int64 `protobuf:"varint,1,opt,name=appId,proto3" json:"appId"`
	// 用户授权码
	Code string `protobuf:"bytes,2,opt,name=code,proto3" json:"code"`
	// 类型
	Type string `protobuf:"bytes,3,opt,name=type,proto3" json:"type"`
}

func (x *GetUserOpenIdRequest) Reset() {
	*x = GetUserOpenIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provider_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserOpenIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserOpenIdRequest) ProtoMessage() {}

func (x *GetUserOpenIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_provider_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserOpenIdRequest.ProtoReflect.Descriptor instead.
func (*GetUserOpenIdRequest) Descriptor() ([]byte, []int) {
	return file_provider_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserOpenIdRequest) GetAppId() int64 {
	if x != nil {
		return x.AppId
	}
	return 0
}

func (x *GetUserOpenIdRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *GetUserOpenIdRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

// 用户OpenId响应
type UserOpenIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 状态码
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	// 状态信息
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message"`
	// 用户OpenId
	OpenId string `protobuf:"bytes,3,opt,name=openId,proto3" json:"openId"`
	// 用户UnionId
	UnionId string `protobuf:"bytes,4,opt,name=unionId,proto3" json:"unionId"`
	// 返回第三方应用Id, 如微信小程序的appId
	AppId string `protobuf:"bytes,5,opt,name=appId,proto3" json:"appId"`
}

func (x *UserOpenIdResponse) Reset() {
	*x = UserOpenIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provider_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserOpenIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserOpenIdResponse) ProtoMessage() {}

func (x *UserOpenIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_provider_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserOpenIdResponse.ProtoReflect.Descriptor instead.
func (*UserOpenIdResponse) Descriptor() ([]byte, []int) {
	return file_provider_service_proto_rawDescGZIP(), []int{1}
}

func (x *UserOpenIdResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UserOpenIdResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UserOpenIdResponse) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

func (x *UserOpenIdResponse) GetUnionId() string {
	if x != nil {
		return x.UnionId
	}
	return ""
}

func (x *UserOpenIdResponse) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

// 小程序二维码请求
type MPCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 应用ID
	AppId int64 `protobuf:"varint,1,opt,name=appId,proto3" json:"appId"`
	// 页面路径
	Page string `protobuf:"bytes,2,opt,name=page,proto3" json:"page"`
	// 场景值,可用来传递参数，并在页面上获取，如: orderNo=20260102
	Scene string `protobuf:"bytes,3,opt,name=scene,proto3" json:"scene"`
	// 是否保存到本地
	SaveLocal bool `protobuf:"varint,4,opt,name=saveLocal,proto3" json:"saveLocal"`
	// 所有者Key,如果保存到本地则需设置
	OwnerKey string `protobuf:"bytes,5,opt,name=ownerKey,proto3" json:"ownerKey"`
}

func (x *MPCodeRequest) Reset() {
	*x = MPCodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provider_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MPCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MPCodeRequest) ProtoMessage() {}

func (x *MPCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_provider_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MPCodeRequest.ProtoReflect.Descriptor instead.
func (*MPCodeRequest) Descriptor() ([]byte, []int) {
	return file_provider_service_proto_rawDescGZIP(), []int{2}
}

func (x *MPCodeRequest) GetAppId() int64 {
	if x != nil {
		return x.AppId
	}
	return 0
}

func (x *MPCodeRequest) GetPage() string {
	if x != nil {
		return x.Page
	}
	return ""
}

func (x *MPCodeRequest) GetScene() string {
	if x != nil {
		return x.Scene
	}
	return ""
}

func (x *MPCodeRequest) GetSaveLocal() bool {
	if x != nil {
		return x.SaveLocal
	}
	return false
}

func (x *MPCodeRequest) GetOwnerKey() string {
	if x != nil {
		return x.OwnerKey
	}
	return ""
}

// 小程序二维码响应
type MPQrCodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 状态码
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	// 状态信息
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message"`
	// 二维码图片
	QrCodeUrl string `protobuf:"bytes,3,opt,name=qrCodeUrl,proto3" json:"qrCodeUrl"`
}

func (x *MPQrCodeResponse) Reset() {
	*x = MPQrCodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provider_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MPQrCodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MPQrCodeResponse) ProtoMessage() {}

func (x *MPQrCodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_provider_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MPQrCodeResponse.ProtoReflect.Descriptor instead.
func (*MPQrCodeResponse) Descriptor() ([]byte, []int) {
	return file_provider_service_proto_rawDescGZIP(), []int{3}
}

func (x *MPQrCodeResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *MPQrCodeResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *MPQrCodeResponse) GetQrCodeUrl() string {
	if x != nil {
		return x.QrCodeUrl
	}
	return ""
}

var File_provider_service_proto protoreflect.FileDescriptor

var file_provider_service_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x8a,
	0x01, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x75,
	0x6e, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x6e,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x22, 0x89, 0x01, 0x0a, 0x0d,
	0x4d, 0x50, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x61, 0x70,
	0x70, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x65, 0x6e, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x73, 0x61, 0x76, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x73, 0x61, 0x76, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x22, 0x5e, 0x0a, 0x10, 0x4d, 0x50, 0x51, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x71, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x55, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x55, 0x72, 0x6c, 0x32, 0x81, 0x01, 0x0a, 0x16, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x12,
	0x15, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x70, 0x65,
	0x6e, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x4d, 0x50, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x2e, 0x4d, 0x50, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x4d, 0x50, 0x51, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1f, 0x0a, 0x13, 0x63,
	0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72,
	0x70, 0x63, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_provider_service_proto_rawDescOnce sync.Once
	file_provider_service_proto_rawDescData = file_provider_service_proto_rawDesc
)

func file_provider_service_proto_rawDescGZIP() []byte {
	file_provider_service_proto_rawDescOnce.Do(func() {
		file_provider_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_provider_service_proto_rawDescData)
	})
	return file_provider_service_proto_rawDescData
}

var file_provider_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_provider_service_proto_goTypes = []interface{}{
	(*GetUserOpenIdRequest)(nil), // 0: GetUserOpenIdRequest
	(*UserOpenIdResponse)(nil),   // 1: UserOpenIdResponse
	(*MPCodeRequest)(nil),        // 2: MPCodeRequest
	(*MPQrCodeResponse)(nil),     // 3: MPQrCodeResponse
}
var file_provider_service_proto_depIdxs = []int32{
	0, // 0: ServiceProviderService.GetOpenId:input_type -> GetUserOpenIdRequest
	2, // 1: ServiceProviderService.GetMPCode:input_type -> MPCodeRequest
	1, // 2: ServiceProviderService.GetOpenId:output_type -> UserOpenIdResponse
	3, // 3: ServiceProviderService.GetMPCode:output_type -> MPQrCodeResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_provider_service_proto_init() }
func file_provider_service_proto_init() {
	if File_provider_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_provider_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserOpenIdRequest); i {
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
		file_provider_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserOpenIdResponse); i {
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
		file_provider_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MPCodeRequest); i {
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
		file_provider_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MPQrCodeResponse); i {
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
			RawDescriptor: file_provider_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_provider_service_proto_goTypes,
		DependencyIndexes: file_provider_service_proto_depIdxs,
		MessageInfos:      file_provider_service_proto_msgTypes,
	}.Build()
	File_provider_service_proto = out.File
	file_provider_service_proto_rawDesc = nil
	file_provider_service_proto_goTypes = nil
	file_provider_service_proto_depIdxs = nil
}
