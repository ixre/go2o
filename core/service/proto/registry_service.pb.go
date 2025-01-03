// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0
// source: registry_service.proto

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

type RegistriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []*SRegistry `protobuf:"bytes,1,rep,name=value,proto3" json:"value"`
}

func (x *RegistriesResponse) Reset() {
	*x = RegistriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistriesResponse) ProtoMessage() {}

func (x *RegistriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistriesResponse.ProtoReflect.Descriptor instead.
func (*RegistriesResponse) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{0}
}

func (x *RegistriesResponse) GetValue() []*SRegistry {
	if x != nil {
		return x.Value
	}
	return nil
}

type RegistryPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value"`
}

func (x *RegistryPair) Reset() {
	*x = RegistryPair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistryPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistryPair) ProtoMessage() {}

func (x *RegistryPair) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistryPair.ProtoReflect.Descriptor instead.
func (*RegistryPair) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{1}
}

func (x *RegistryPair) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *RegistryPair) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type RegistryValueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value  string `protobuf:"bytes,1,opt,name=value,proto3" json:"value"`
	ErrMsg string `protobuf:"bytes,2,opt,name=errMsg,proto3" json:"errMsg"`
}

func (x *RegistryValueResponse) Reset() {
	*x = RegistryValueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistryValueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistryValueResponse) ProtoMessage() {}

func (x *RegistryValueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistryValueResponse.ProtoReflect.Descriptor instead.
func (*RegistryValueResponse) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{2}
}

func (x *RegistryValueResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *RegistryValueResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

type RegistryCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 键
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key"`
	// 分组
	Group string `protobuf:"bytes,2,opt,name=group,proto3" json:"group"`
	// 默认值
	DefaultValue string `protobuf:"bytes,3,opt,name=defaultValue,proto3" json:"defaultValue"`
	// 描述
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description"`
}

func (x *RegistryCreateRequest) Reset() {
	*x = RegistryCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistryCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistryCreateRequest) ProtoMessage() {}

func (x *RegistryCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistryCreateRequest.ProtoReflect.Descriptor instead.
func (*RegistryCreateRequest) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{3}
}

func (x *RegistryCreateRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *RegistryCreateRequest) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *RegistryCreateRequest) GetDefaultValue() string {
	if x != nil {
		return x.DefaultValue
	}
	return ""
}

func (x *RegistryCreateRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type RegistrySearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key"`
}

func (x *RegistrySearchRequest) Reset() {
	*x = RegistrySearchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrySearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrySearchRequest) ProtoMessage() {}

func (x *RegistrySearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrySearchRequest.ProtoReflect.Descriptor instead.
func (*RegistrySearchRequest) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{4}
}

func (x *RegistrySearchRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// * 注册表
type SRegistry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * 键
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key"`
	// * 值
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value"`
	// * 分组
	Group string `protobuf:"bytes,3,opt,name=group,proto3" json:"group"`
	// * 默认值
	DefaultValue string `protobuf:"bytes,4,opt,name=defaultValue,proto3" json:"defaultValue"`
	// * 可选值
	Options string `protobuf:"bytes,5,opt,name=options,proto3" json:"options"`
	// * 标志
	Flag int32 `protobuf:"zigzag32,6,opt,name=flag,proto3" json:"flag"`
	// * 描述
	Description string `protobuf:"bytes,7,opt,name=description,proto3" json:"description"`
}

func (x *SRegistry) Reset() {
	*x = SRegistry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SRegistry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SRegistry) ProtoMessage() {}

func (x *SRegistry) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SRegistry.ProtoReflect.Descriptor instead.
func (*SRegistry) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{5}
}

func (x *SRegistry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SRegistry) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *SRegistry) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *SRegistry) GetDefaultValue() string {
	if x != nil {
		return x.DefaultValue
	}
	return ""
}

func (x *SRegistry) GetOptions() string {
	if x != nil {
		return x.Options
	}
	return ""
}

func (x *SRegistry) GetFlag() int32 {
	if x != nil {
		return x.Flag
	}
	return 0
}

func (x *SRegistry) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type RegistryGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value"`
}

func (x *RegistryGroupResponse) Reset() {
	*x = RegistryGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_registry_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistryGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistryGroupResponse) ProtoMessage() {}

func (x *RegistryGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_registry_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistryGroupResponse.ProtoReflect.Descriptor instead.
func (*RegistryGroupResponse) Descriptor() ([]byte, []int) {
	return file_registry_service_proto_rawDescGZIP(), []int{6}
}

func (x *RegistryGroupResponse) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_registry_service_proto protoreflect.FileDescriptor

var file_registry_service_proto_rawDesc = []byte{
	0x0a, 0x16, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x12, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x53, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x36,
	0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x50, 0x61, 0x69, 0x72, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x45, 0x0a, 0x15, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x22, 0x85, 0x01,
	0x0a, 0x15, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x22, 0x0a, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x29, 0x0a, 0x15, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x22, 0xbd, 0x01, 0x0a, 0x09, 0x53, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x22, 0x0a, 0x0c,
	0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c,
	0x61, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x11, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x2d, 0x0a, 0x15, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32,
	0xd4, 0x03, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x2d, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73,
	0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x24, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x12, 0x07, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x0a, 0x2e, 0x53, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x22, 0x00, 0x12, 0x2d, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x07, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x16, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0d, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x50, 0x61, 0x69, 0x72, 0x1a, 0x09, 0x2e, 0x54, 0x78, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x00, 0x12, 0x27, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12,
	0x0c, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x41, 0x72, 0x72, 0x61, 0x79, 0x1a, 0x0a, 0x2e,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x22, 0x00, 0x12, 0x27, 0x0a, 0x0c, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x0a, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x1a, 0x09, 0x2e, 0x54, 0x78, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x00, 0x12, 0x2e, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x16,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4d,
	0x61, 0x70, 0x22, 0x00, 0x12, 0x27, 0x0a, 0x0e, 0x46, 0x69, 0x6e, 0x64, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x07, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a,
	0x0a, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x70, 0x22, 0x00, 0x12, 0x30, 0x0a,
	0x0e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12,
	0x07, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x13, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x35, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x12, 0x16, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x54, 0x78, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x42, 0x1f, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72, 0x70, 0x63, 0x5a, 0x08, 0x2e,
	0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_registry_service_proto_rawDescOnce sync.Once
	file_registry_service_proto_rawDescData = file_registry_service_proto_rawDesc
)

func file_registry_service_proto_rawDescGZIP() []byte {
	file_registry_service_proto_rawDescOnce.Do(func() {
		file_registry_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_registry_service_proto_rawDescData)
	})
	return file_registry_service_proto_rawDescData
}

var file_registry_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_registry_service_proto_goTypes = []interface{}{
	(*RegistriesResponse)(nil),    // 0: RegistriesResponse
	(*RegistryPair)(nil),          // 1: RegistryPair
	(*RegistryValueResponse)(nil), // 2: RegistryValueResponse
	(*RegistryCreateRequest)(nil), // 3: RegistryCreateRequest
	(*RegistrySearchRequest)(nil), // 4: RegistrySearchRequest
	(*SRegistry)(nil),             // 5: SRegistry
	(*RegistryGroupResponse)(nil), // 6: RegistryGroupResponse
	(*Empty)(nil),                 // 7: Empty
	(*String)(nil),                // 8: String
	(*StringArray)(nil),           // 9: StringArray
	(*StringMap)(nil),             // 10: StringMap
	(*TxResult)(nil),              // 11: TxResult
}
var file_registry_service_proto_depIdxs = []int32{
	5,  // 0: RegistriesResponse.value:type_name -> SRegistry
	7,  // 1: RegistryService.GetGroups:input_type -> Empty
	8,  // 2: RegistryService.GetRegistry:input_type -> String
	8,  // 3: RegistryService.GetValue:input_type -> String
	1,  // 4: RegistryService.UpdateValue:input_type -> RegistryPair
	9,  // 5: RegistryService.GetValues:input_type -> StringArray
	10, // 6: RegistryService.UpdateValues:input_type -> StringMap
	4,  // 7: RegistryService.Search:input_type -> RegistrySearchRequest
	8,  // 8: RegistryService.FindRegistries:input_type -> String
	8,  // 9: RegistryService.SearchRegistry:input_type -> String
	3,  // 10: RegistryService.CreateRegistry:input_type -> RegistryCreateRequest
	6,  // 11: RegistryService.GetGroups:output_type -> RegistryGroupResponse
	5,  // 12: RegistryService.GetRegistry:output_type -> SRegistry
	2,  // 13: RegistryService.GetValue:output_type -> RegistryValueResponse
	11, // 14: RegistryService.UpdateValue:output_type -> TxResult
	10, // 15: RegistryService.GetValues:output_type -> StringMap
	11, // 16: RegistryService.UpdateValues:output_type -> TxResult
	10, // 17: RegistryService.Search:output_type -> StringMap
	10, // 18: RegistryService.FindRegistries:output_type -> StringMap
	0,  // 19: RegistryService.SearchRegistry:output_type -> RegistriesResponse
	11, // 20: RegistryService.CreateRegistry:output_type -> TxResult
	11, // [11:21] is the sub-list for method output_type
	1,  // [1:11] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_registry_service_proto_init() }
func file_registry_service_proto_init() {
	if File_registry_service_proto != nil {
		return
	}
	file_global_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_registry_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistriesResponse); i {
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
		file_registry_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistryPair); i {
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
		file_registry_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistryValueResponse); i {
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
		file_registry_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistryCreateRequest); i {
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
		file_registry_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrySearchRequest); i {
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
		file_registry_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SRegistry); i {
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
		file_registry_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistryGroupResponse); i {
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
			RawDescriptor: file_registry_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_registry_service_proto_goTypes,
		DependencyIndexes: file_registry_service_proto_depIdxs,
		MessageInfos:      file_registry_service_proto_msgTypes,
	}.Build()
	File_registry_service_proto = out.File
	file_registry_service_proto_rawDesc = nil
	file_registry_service_proto_goTypes = nil
	file_registry_service_proto_depIdxs = nil
}
