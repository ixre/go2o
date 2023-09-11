// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: app_service.proto

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

type AppId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value"`
}

func (x *AppId) Reset() {
	*x = AppId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppId) ProtoMessage() {}

func (x *AppId) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppId.ProtoReflect.Descriptor instead.
func (*AppId) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{0}
}

func (x *AppId) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type AppVersionId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value"`
}

func (x *AppVersionId) Reset() {
	*x = AppVersionId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppVersionId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppVersionId) ProtoMessage() {}

func (x *AppVersionId) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppVersionId.ProtoReflect.Descriptor instead.
func (*AppVersionId) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{1}
}

func (x *AppVersionId) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// 检查版本请求
type CheckVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 版本号
	AppId int64 `protobuf:"varint,1,opt,name=appId,proto3" json:"appId"`
	// 更新通道, stable|beta|nightly
	Channel string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel"`
	// 当前版本
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version"`
}

func (x *CheckVersionRequest) Reset() {
	*x = CheckVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckVersionRequest) ProtoMessage() {}

func (x *CheckVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckVersionRequest.ProtoReflect.Descriptor instead.
func (*CheckVersionRequest) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{2}
}

func (x *CheckVersionRequest) GetAppId() int64 {
	if x != nil {
		return x.AppId
	}
	return 0
}

func (x *CheckVersionRequest) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *CheckVersionRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// 检测版本响应结果
type CheckVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 最新版本号
	LatestVersion string `protobuf:"bytes,1,opt,name=latestVersion,proto3" json:"latestVersion"`
	// App更新资源地址
	AppPkgURL string `protobuf:"bytes,2,opt,name=appPkgURL,proto3" json:"appPkgURL"`
	// 版本信息
	VersionInfo string `protobuf:"bytes,3,opt,name=versionInfo,proto3" json:"versionInfo"`
	// 是否为最新版本
	IsLatest bool `protobuf:"varint,4,opt,name=isLatest,proto3" json:"isLatest"`
	// 是否强制升级
	ForceUpdate bool `protobuf:"varint,5,opt,name=forceUpdate,proto3" json:"forceUpdate"`
	// 更新文件类型,如APK,EXE,ZIP等
	UpdateType string `protobuf:"bytes,6,opt,name=updateType,proto3" json:"updateType"`
	// 发布时间
	ReleaseTime int64 `protobuf:"varint,7,opt,name=releaseTime,proto3" json:"releaseTime"`
}

func (x *CheckVersionResponse) Reset() {
	*x = CheckVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckVersionResponse) ProtoMessage() {}

func (x *CheckVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckVersionResponse.ProtoReflect.Descriptor instead.
func (*CheckVersionResponse) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{3}
}

func (x *CheckVersionResponse) GetLatestVersion() string {
	if x != nil {
		return x.LatestVersion
	}
	return ""
}

func (x *CheckVersionResponse) GetAppPkgURL() string {
	if x != nil {
		return x.AppPkgURL
	}
	return ""
}

func (x *CheckVersionResponse) GetVersionInfo() string {
	if x != nil {
		return x.VersionInfo
	}
	return ""
}

func (x *CheckVersionResponse) GetIsLatest() bool {
	if x != nil {
		return x.IsLatest
	}
	return false
}

func (x *CheckVersionResponse) GetForceUpdate() bool {
	if x != nil {
		return x.ForceUpdate
	}
	return false
}

func (x *CheckVersionResponse) GetUpdateType() string {
	if x != nil {
		return x.UpdateType
	}
	return ""
}

func (x *CheckVersionResponse) GetReleaseTime() int64 {
	if x != nil {
		return x.ReleaseTime
	}
	return 0
}

// APP产品
type AppProdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 产品编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 产品名称
	ProdName string `protobuf:"bytes,2,opt,name=prodName,proto3" json:"prodName"`
	// 产品描述
	ProdDes string `protobuf:"bytes,3,opt,name=prodDes,proto3" json:"prodDes"`
	// Icon
	Icon string `protobuf:"bytes,4,opt,name=icon,proto3" json:"icon"`
	// 发布下载页面地址
	PublishURL string `protobuf:"bytes,5,opt,name=publishURL,proto3" json:"publishURL"`
	// 正式版文件地址
	StableFileURL string `protobuf:"bytes,6,opt,name=stableFileURL,proto3" json:"stableFileURL"`
	// 内测版文件地址
	AlphaFileURL string `protobuf:"bytes,8,opt,name=alphaFileURL,proto3" json:"alphaFileURL"`
	// 每夜版文件地址
	NightlyFileURL string `protobuf:"bytes,10,opt,name=nightlyFileURL,proto3" json:"nightlyFileURL"`
	// 更新方式,比如APK, EXE等
	UpdateType string `protobuf:"bytes,11,opt,name=updateType,proto3" json:"updateType"`
}

func (x *AppProdRequest) Reset() {
	*x = AppProdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppProdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppProdRequest) ProtoMessage() {}

func (x *AppProdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppProdRequest.ProtoReflect.Descriptor instead.
func (*AppProdRequest) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{4}
}

func (x *AppProdRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AppProdRequest) GetProdName() string {
	if x != nil {
		return x.ProdName
	}
	return ""
}

func (x *AppProdRequest) GetProdDes() string {
	if x != nil {
		return x.ProdDes
	}
	return ""
}

func (x *AppProdRequest) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *AppProdRequest) GetPublishURL() string {
	if x != nil {
		return x.PublishURL
	}
	return ""
}

func (x *AppProdRequest) GetStableFileURL() string {
	if x != nil {
		return x.StableFileURL
	}
	return ""
}

func (x *AppProdRequest) GetAlphaFileURL() string {
	if x != nil {
		return x.AlphaFileURL
	}
	return ""
}

func (x *AppProdRequest) GetNightlyFileURL() string {
	if x != nil {
		return x.NightlyFileURL
	}
	return ""
}

func (x *AppProdRequest) GetUpdateType() string {
	if x != nil {
		return x.UpdateType
	}
	return ""
}

// APP版本
type AppVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 产品
	ProductId int64 `protobuf:"varint,2,opt,name=productId,proto3" json:"productId"`
	// 更新通道, stable:0|alpha:1|nightly:2
	Channel int32 `protobuf:"varint,3,opt,name=channel,proto3" json:"channel"`
	// 版本号
	Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version"`
	// 是否强制升级
	ForceUpdate bool `protobuf:"varint,5,opt,name=forceUpdate,proto3" json:"forceUpdate"`
	// 更新内容
	UpdateContent string `protobuf:"bytes,6,opt,name=updateContent,proto3" json:"updateContent"`
}

func (x *AppVersionRequest) Reset() {
	*x = AppVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppVersionRequest) ProtoMessage() {}

func (x *AppVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppVersionRequest.ProtoReflect.Descriptor instead.
func (*AppVersionRequest) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{5}
}

func (x *AppVersionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AppVersionRequest) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *AppVersionRequest) GetChannel() int32 {
	if x != nil {
		return x.Channel
	}
	return 0
}

func (x *AppVersionRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *AppVersionRequest) GetForceUpdate() bool {
	if x != nil {
		return x.ForceUpdate
	}
	return false
}

func (x *AppVersionRequest) GetUpdateContent() string {
	if x != nil {
		return x.UpdateContent
	}
	return ""
}

// APP产品
type SAppProd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 产品编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 产品名称
	ProdName string `protobuf:"bytes,2,opt,name=prodName,proto3" json:"prodName"`
	// 产品描述
	ProdDes string `protobuf:"bytes,3,opt,name=prodDes,proto3" json:"prodDes"`
	// 最新的版本ID
	LatestVid int64 `protobuf:"varint,4,opt,name=latestVid,proto3" json:"latestVid"`
	// 正式版文件hash值
	Md5Hash string `protobuf:"bytes,5,opt,name=md5Hash,proto3" json:"md5Hash"`
	// 发布下载页面地址
	PublishURL string `protobuf:"bytes,6,opt,name=publishURL,proto3" json:"publishURL"`
	// 正式版文件地址
	StableFileURL string `protobuf:"bytes,7,opt,name=stableFileURL,proto3" json:"stableFileURL"`
	// 内测版文件地址
	AlphaFileURL string `protobuf:"bytes,8,opt,name=alphaFileURL,proto3" json:"alphaFileURL"`
	// 每夜版文件地址
	NightlyFileURL string `protobuf:"bytes,9,opt,name=nightlyFileURL,proto3" json:"nightlyFileURL"`
	// 更新方式,比如APK, EXE等
	UpdateType string `protobuf:"bytes,10,opt,name=updateType,proto3" json:"updateType"`
	// 更新时间
	UpdateTime int64 `protobuf:"varint,11,opt,name=updateTime,proto3" json:"updateTime"`
}

func (x *SAppProd) Reset() {
	*x = SAppProd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SAppProd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SAppProd) ProtoMessage() {}

func (x *SAppProd) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SAppProd.ProtoReflect.Descriptor instead.
func (*SAppProd) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{6}
}

func (x *SAppProd) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SAppProd) GetProdName() string {
	if x != nil {
		return x.ProdName
	}
	return ""
}

func (x *SAppProd) GetProdDes() string {
	if x != nil {
		return x.ProdDes
	}
	return ""
}

func (x *SAppProd) GetLatestVid() int64 {
	if x != nil {
		return x.LatestVid
	}
	return 0
}

func (x *SAppProd) GetMd5Hash() string {
	if x != nil {
		return x.Md5Hash
	}
	return ""
}

func (x *SAppProd) GetPublishURL() string {
	if x != nil {
		return x.PublishURL
	}
	return ""
}

func (x *SAppProd) GetStableFileURL() string {
	if x != nil {
		return x.StableFileURL
	}
	return ""
}

func (x *SAppProd) GetAlphaFileURL() string {
	if x != nil {
		return x.AlphaFileURL
	}
	return ""
}

func (x *SAppProd) GetNightlyFileURL() string {
	if x != nil {
		return x.NightlyFileURL
	}
	return ""
}

func (x *SAppProd) GetUpdateType() string {
	if x != nil {
		return x.UpdateType
	}
	return ""
}

func (x *SAppProd) GetUpdateTime() int64 {
	if x != nil {
		return x.UpdateTime
	}
	return 0
}

// APP版本
type SAppVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 编号
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	// 产品
	ProductId int64 `protobuf:"varint,2,opt,name=productId,proto3" json:"productId"`
	// 更新通道, 0:stable|1:beta|2:nightly
	Channel int32 `protobuf:"varint,3,opt,name=channel,proto3" json:"channel"`
	// 版本号
	Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version"`
	// 数字版本
	VersionCode int32 `protobuf:"varint,5,opt,name=versionCode,proto3" json:"versionCode"`
	// 是否强制升级
	ForceUpdate bool `protobuf:"varint,6,opt,name=forceUpdate,proto3" json:"forceUpdate"`
	// 更新内容
	UpdateContent string `protobuf:"bytes,7,opt,name=updateContent,proto3" json:"updateContent"`
	// 发布时间
	CreateTime int64 `protobuf:"varint,8,opt,name=createTime,proto3" json:"createTime"`
}

func (x *SAppVersion) Reset() {
	*x = SAppVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SAppVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SAppVersion) ProtoMessage() {}

func (x *SAppVersion) ProtoReflect() protoreflect.Message {
	mi := &file_app_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SAppVersion.ProtoReflect.Descriptor instead.
func (*SAppVersion) Descriptor() ([]byte, []int) {
	return file_app_service_proto_rawDescGZIP(), []int{7}
}

func (x *SAppVersion) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SAppVersion) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *SAppVersion) GetChannel() int32 {
	if x != nil {
		return x.Channel
	}
	return 0
}

func (x *SAppVersion) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *SAppVersion) GetVersionCode() int32 {
	if x != nil {
		return x.VersionCode
	}
	return 0
}

func (x *SAppVersion) GetForceUpdate() bool {
	if x != nil {
		return x.ForceUpdate
	}
	return false
}

func (x *SAppVersion) GetUpdateContent() string {
	if x != nil {
		return x.UpdateContent
	}
	return ""
}

func (x *SAppVersion) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

var File_app_service_proto protoreflect.FileDescriptor

var file_app_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x1d, 0x0a, 0x05, 0x41, 0x70, 0x70, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x24, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x5f, 0x0a, 0x13, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x61, 0x70,
	0x70, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xfc, 0x01, 0x0a, 0x14, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x70, 0x70, 0x50, 0x6b, 0x67,
	0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x70, 0x70, 0x50, 0x6b,
	0x67, 0x55, 0x52, 0x4c, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x4c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x72, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x9c, 0x02, 0x0a, 0x0e, 0x41, 0x70, 0x70, 0x50, 0x72,
	0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f,
	0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f,
	0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x44, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x44, 0x65, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69,
	0x63, 0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x55, 0x52,
	0x4c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x55, 0x52, 0x4c, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x55, 0x52, 0x4c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x26, 0x0a,
	0x0e, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x6c, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x6c, 0x79, 0x46, 0x69,
	0x6c, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0xbd, 0x01, 0x0a, 0x11, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a,
	0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x24, 0x0a, 0x0d, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0xda, 0x02, 0x0a, 0x08, 0x53, 0x41, 0x70, 0x70, 0x50, 0x72,
	0x6f, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x44, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x70, 0x72, 0x6f, 0x64, 0x44, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x56, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x56, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x64, 0x35, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x64, 0x35, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x55, 0x52, 0x4c, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x55, 0x52, 0x4c,
	0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52,
	0x4c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x46,
	0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x26, 0x0a, 0x0e, 0x6e, 0x69,
	0x67, 0x68, 0x74, 0x6c, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x52, 0x4c, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x6e, 0x69, 0x67, 0x68, 0x74, 0x6c, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x55,
	0x52, 0x4c, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x22, 0xf9, 0x01, 0x0a, 0x0b, 0x53, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x66, 0x6f, 0x72,
	0x63, 0x65, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1e,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x32, 0xba,
	0x02, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x26, 0x0a,
	0x08, 0x53, 0x61, 0x76, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x12, 0x0f, 0x2e, 0x41, 0x70, 0x70, 0x50,
	0x72, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x0b, 0x53, 0x61, 0x76, 0x65, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x2e, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x00, 0x12, 0x1e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x12, 0x06,
	0x2e, 0x41, 0x70, 0x70, 0x49, 0x64, 0x1a, 0x09, 0x2e, 0x53, 0x41, 0x70, 0x70, 0x50, 0x72, 0x6f,
	0x64, 0x22, 0x00, 0x12, 0x2b, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x0d, 0x2e, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x1a, 0x0c, 0x2e, 0x53, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x00,
	0x12, 0x1f, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x12, 0x06,
	0x2e, 0x41, 0x70, 0x70, 0x49, 0x64, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22,
	0x00, 0x12, 0x29, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x0d, 0x2e, 0x41, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0c,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x2e, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x15, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1f, 0x0a, 0x13, 0x63,
	0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72,
	0x70, 0x63, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_service_proto_rawDescOnce sync.Once
	file_app_service_proto_rawDescData = file_app_service_proto_rawDesc
)

func file_app_service_proto_rawDescGZIP() []byte {
	file_app_service_proto_rawDescOnce.Do(func() {
		file_app_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_service_proto_rawDescData)
	})
	return file_app_service_proto_rawDescData
}

var file_app_service_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_app_service_proto_goTypes = []interface{}{
	(*AppId)(nil),                // 0: AppId
	(*AppVersionId)(nil),         // 1: AppVersionId
	(*CheckVersionRequest)(nil),  // 2: CheckVersionRequest
	(*CheckVersionResponse)(nil), // 3: CheckVersionResponse
	(*AppProdRequest)(nil),       // 4: AppProdRequest
	(*AppVersionRequest)(nil),    // 5: AppVersionRequest
	(*SAppProd)(nil),             // 6: SAppProd
	(*SAppVersion)(nil),          // 7: SAppVersion
	(*Result)(nil),               // 8: Result
}
var file_app_service_proto_depIdxs = []int32{
	4, // 0: AppService.SaveProd:input_type -> AppProdRequest
	5, // 1: AppService.SaveVersion:input_type -> AppVersionRequest
	0, // 2: AppService.GetProd:input_type -> AppId
	1, // 3: AppService.GetVersion:input_type -> AppVersionId
	0, // 4: AppService.DeleteProd:input_type -> AppId
	1, // 5: AppService.DeleteVersion:input_type -> AppVersionId
	2, // 6: AppService.CheckVersion:input_type -> CheckVersionRequest
	8, // 7: AppService.SaveProd:output_type -> Result
	8, // 8: AppService.SaveVersion:output_type -> Result
	6, // 9: AppService.GetProd:output_type -> SAppProd
	7, // 10: AppService.GetVersion:output_type -> SAppVersion
	8, // 11: AppService.DeleteProd:output_type -> Result
	8, // 12: AppService.DeleteVersion:output_type -> Result
	3, // 13: AppService.CheckVersion:output_type -> CheckVersionResponse
	7, // [7:14] is the sub-list for method output_type
	0, // [0:7] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_app_service_proto_init() }
func file_app_service_proto_init() {
	if File_app_service_proto != nil {
		return
	}
	file_global_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_app_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppId); i {
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
		file_app_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppVersionId); i {
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
		file_app_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckVersionRequest); i {
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
		file_app_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckVersionResponse); i {
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
		file_app_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppProdRequest); i {
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
		file_app_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppVersionRequest); i {
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
		file_app_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SAppProd); i {
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
		file_app_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SAppVersion); i {
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
			RawDescriptor: file_app_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_app_service_proto_goTypes,
		DependencyIndexes: file_app_service_proto_depIdxs,
		MessageInfos:      file_app_service_proto_msgTypes,
	}.Build()
	File_app_service_proto = out.File
	file_app_service_proto_rawDesc = nil
	file_app_service_proto_goTypes = nil
	file_app_service_proto_depIdxs = nil
}
