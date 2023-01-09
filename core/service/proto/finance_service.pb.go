// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.4
// source: finance_service.proto

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

type TransferInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PersonId     int64   `protobuf:"zigzag64,1,opt,name=personId,proto3" json:"personId"`
	TransferWith int32   `protobuf:"zigzag32,2,opt,name=transferWith,proto3" json:"transferWith"`
	Amount       float64 `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount"`
}

func (x *TransferInRequest) Reset() {
	*x = TransferInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_finance_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferInRequest) ProtoMessage() {}

func (x *TransferInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_finance_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferInRequest.ProtoReflect.Descriptor instead.
func (*TransferInRequest) Descriptor() ([]byte, []int) {
	return file_finance_service_proto_rawDescGZIP(), []int{0}
}

func (x *TransferInRequest) GetPersonId() int64 {
	if x != nil {
		return x.PersonId
	}
	return 0
}

func (x *TransferInRequest) GetTransferWith() int32 {
	if x != nil {
		return x.TransferWith
	}
	return 0
}

func (x *TransferInRequest) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type PersonId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value"`
}

func (x *PersonId) Reset() {
	*x = PersonId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_finance_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersonId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersonId) ProtoMessage() {}

func (x *PersonId) ProtoReflect() protoreflect.Message {
	mi := &file_finance_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersonId.ProtoReflect.Descriptor instead.
func (*PersonId) Descriptor() ([]byte, []int) {
	return file_finance_service_proto_rawDescGZIP(), []int{1}
}

func (x *PersonId) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// 收益总记录
type SRiseInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 人员编号
	PersonId int64 `protobuf:"varint,1,opt,name=personId,proto3" json:"personId"`
	// 本金及收益的余额
	Balance float64 `protobuf:"fixed64,2,opt,name=balance,proto3" json:"balance"`
	// 结算金额
	SettlementAmount float64 `protobuf:"fixed64,3,opt,name=settlementAmount,proto3" json:"settlementAmount"`
	// 当前的收益
	Rise float64 `protobuf:"fixed64,4,opt,name=rise,proto3" json:"rise"`
	// 今日转入
	TransferIn float64 `protobuf:"fixed64,5,opt,name=transferIn,proto3" json:"transferIn"`
	// 总金额
	TotalAmount float64 `protobuf:"fixed64,6,opt,name=totalAmount,proto3" json:"totalAmount"`
	// 总收益
	TotalRise float64 `protobuf:"fixed64,7,opt,name=totalRise,proto3" json:"totalRise"`
	// 结算日期,用于筛选需要结算的数据
	SettledDate int64 `protobuf:"varint,8,opt,name=settledDate,proto3" json:"settledDate"`
	// 更新时间
	UpdateTime int64 `protobuf:"varint,9,opt,name=updateTime,proto3" json:"updateTime"`
}

func (x *SRiseInfo) Reset() {
	*x = SRiseInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_finance_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SRiseInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SRiseInfo) ProtoMessage() {}

func (x *SRiseInfo) ProtoReflect() protoreflect.Message {
	mi := &file_finance_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SRiseInfo.ProtoReflect.Descriptor instead.
func (*SRiseInfo) Descriptor() ([]byte, []int) {
	return file_finance_service_proto_rawDescGZIP(), []int{2}
}

func (x *SRiseInfo) GetPersonId() int64 {
	if x != nil {
		return x.PersonId
	}
	return 0
}

func (x *SRiseInfo) GetBalance() float64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

func (x *SRiseInfo) GetSettlementAmount() float64 {
	if x != nil {
		return x.SettlementAmount
	}
	return 0
}

func (x *SRiseInfo) GetRise() float64 {
	if x != nil {
		return x.Rise
	}
	return 0
}

func (x *SRiseInfo) GetTransferIn() float64 {
	if x != nil {
		return x.TransferIn
	}
	return 0
}

func (x *SRiseInfo) GetTotalAmount() float64 {
	if x != nil {
		return x.TotalAmount
	}
	return 0
}

func (x *SRiseInfo) GetTotalRise() float64 {
	if x != nil {
		return x.TotalRise
	}
	return 0
}

func (x *SRiseInfo) GetSettledDate() int64 {
	if x != nil {
		return x.SettledDate
	}
	return 0
}

func (x *SRiseInfo) GetUpdateTime() int64 {
	if x != nil {
		return x.UpdateTime
	}
	return 0
}

type RiseSettleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PersonId  int64   `protobuf:"varint,1,opt,name=personId,proto3" json:"personId"`
	SettleDay int64   `protobuf:"varint,2,opt,name=settleDay,proto3" json:"settleDay"`
	Ratio     float64 `protobuf:"fixed64,3,opt,name=ratio,proto3" json:"ratio"`
}

func (x *RiseSettleRequest) Reset() {
	*x = RiseSettleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_finance_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RiseSettleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RiseSettleRequest) ProtoMessage() {}

func (x *RiseSettleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_finance_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RiseSettleRequest.ProtoReflect.Descriptor instead.
func (*RiseSettleRequest) Descriptor() ([]byte, []int) {
	return file_finance_service_proto_rawDescGZIP(), []int{3}
}

func (x *RiseSettleRequest) GetPersonId() int64 {
	if x != nil {
		return x.PersonId
	}
	return 0
}

func (x *RiseSettleRequest) GetSettleDay() int64 {
	if x != nil {
		return x.SettleDay
	}
	return 0
}

func (x *RiseSettleRequest) GetRatio() float64 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

type RiseTransferOutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PersonId     int64   `protobuf:"varint,1,opt,name=personId,proto3" json:"personId"`
	TransferWith int64   `protobuf:"varint,2,opt,name=transferWith,proto3" json:"transferWith"`
	Amount       float64 `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount"`
	// 提现银行账号
	BankAccountNo string `protobuf:"bytes,4,opt,name=bankAccountNo,proto3" json:"bankAccountNo"`
}

func (x *RiseTransferOutRequest) Reset() {
	*x = RiseTransferOutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_finance_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RiseTransferOutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RiseTransferOutRequest) ProtoMessage() {}

func (x *RiseTransferOutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_finance_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RiseTransferOutRequest.ProtoReflect.Descriptor instead.
func (*RiseTransferOutRequest) Descriptor() ([]byte, []int) {
	return file_finance_service_proto_rawDescGZIP(), []int{4}
}

func (x *RiseTransferOutRequest) GetPersonId() int64 {
	if x != nil {
		return x.PersonId
	}
	return 0
}

func (x *RiseTransferOutRequest) GetTransferWith() int64 {
	if x != nil {
		return x.TransferWith
	}
	return 0
}

func (x *RiseTransferOutRequest) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *RiseTransferOutRequest) GetBankAccountNo() string {
	if x != nil {
		return x.BankAccountNo
	}
	return ""
}

type CommitTransferRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PersonId int64 `protobuf:"varint,1,opt,name=personId,proto3" json:"personId"`
	LogId    int64 `protobuf:"varint,2,opt,name=logId,proto3" json:"logId"`
}

func (x *CommitTransferRequest) Reset() {
	*x = CommitTransferRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_finance_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitTransferRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitTransferRequest) ProtoMessage() {}

func (x *CommitTransferRequest) ProtoReflect() protoreflect.Message {
	mi := &file_finance_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitTransferRequest.ProtoReflect.Descriptor instead.
func (*CommitTransferRequest) Descriptor() ([]byte, []int) {
	return file_finance_service_proto_rawDescGZIP(), []int{5}
}

func (x *CommitTransferRequest) GetPersonId() int64 {
	if x != nil {
		return x.PersonId
	}
	return 0
}

func (x *CommitTransferRequest) GetLogId() int64 {
	if x != nil {
		return x.LogId
	}
	return 0
}

var File_finance_service_proto protoreflect.FileDescriptor

var file_finance_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6b, 0x0a, 0x11, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x08, 0x70, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x11, 0x52, 0x0c, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x20, 0x0a, 0x08, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0xa3, 0x02, 0x0a, 0x09, 0x53, 0x52, 0x69, 0x73, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x73, 0x65, 0x74, 0x74,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x10, 0x73, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x69, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x04, 0x72, 0x69, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x49, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x52, 0x69, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x52, 0x69, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x65, 0x74, 0x74,
	0x6c, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x73,
	0x65, 0x74, 0x74, 0x6c, 0x65, 0x64, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x63, 0x0a, 0x11, 0x52, 0x69,
	0x73, 0x65, 0x53, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x65, 0x74, 0x74, 0x6c, 0x65, 0x44, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x73, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x44, 0x61, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x22,
	0x96, 0x01, 0x0a, 0x16, 0x52, 0x69, 0x73, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x4f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x61, 0x6e, 0x6b, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x4e, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x62, 0x61, 0x6e, 0x6b, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x6f, 0x22, 0x49, 0x0a, 0x15, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x6f, 0x67, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x6f,
	0x67, 0x49, 0x64, 0x32, 0xb0, 0x02, 0x0a, 0x0e, 0x46, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x26, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x69, 0x73,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x09, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64,
	0x1a, 0x0a, 0x2e, 0x53, 0x52, 0x69, 0x73, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x00, 0x12, 0x2f,
	0x0a, 0x0e, 0x52, 0x69, 0x73, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x6e,
	0x12, 0x12, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x49, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12,
	0x35, 0x0a, 0x0f, 0x52, 0x69, 0x73, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x4f,
	0x75, 0x74, 0x12, 0x17, 0x2e, 0x52, 0x69, 0x73, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x4f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x07, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x30, 0x0a, 0x0f, 0x52, 0x69, 0x73, 0x65, 0x53, 0x65,
	0x74, 0x74, 0x6c, 0x65, 0x42, 0x79, 0x44, 0x61, 0x79, 0x12, 0x12, 0x2e, 0x52, 0x69, 0x73, 0x65,
	0x53, 0x65, 0x74, 0x74, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x07, 0x2e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x07, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x12, 0x27, 0x0a,
	0x0f, 0x4f, 0x70, 0x65, 0x6e, 0x52, 0x69, 0x73, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x09, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x1a, 0x07, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x42, 0x1f, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72, 0x70, 0x63, 0x5a, 0x08, 0x2e,
	0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_finance_service_proto_rawDescOnce sync.Once
	file_finance_service_proto_rawDescData = file_finance_service_proto_rawDesc
)

func file_finance_service_proto_rawDescGZIP() []byte {
	file_finance_service_proto_rawDescOnce.Do(func() {
		file_finance_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_finance_service_proto_rawDescData)
	})
	return file_finance_service_proto_rawDescData
}

var file_finance_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_finance_service_proto_goTypes = []interface{}{
	(*TransferInRequest)(nil),      // 0: TransferInRequest
	(*PersonId)(nil),               // 1: PersonId
	(*SRiseInfo)(nil),              // 2: SRiseInfo
	(*RiseSettleRequest)(nil),      // 3: RiseSettleRequest
	(*RiseTransferOutRequest)(nil), // 4: RiseTransferOutRequest
	(*CommitTransferRequest)(nil),  // 5: CommitTransferRequest
	(*Result)(nil),                 // 6: Result
}
var file_finance_service_proto_depIdxs = []int32{
	1, // 0: FinanceService.GetRiseInfo:input_type -> PersonId
	0, // 1: FinanceService.RiseTransferIn:input_type -> TransferInRequest
	4, // 2: FinanceService.RiseTransferOut:input_type -> RiseTransferOutRequest
	3, // 3: FinanceService.RiseSettleByDay:input_type -> RiseSettleRequest
	5, // 4: FinanceService.CommitTransfer:input_type -> CommitTransferRequest
	1, // 5: FinanceService.OpenRiseService:input_type -> PersonId
	2, // 6: FinanceService.GetRiseInfo:output_type -> SRiseInfo
	6, // 7: FinanceService.RiseTransferIn:output_type -> Result
	6, // 8: FinanceService.RiseTransferOut:output_type -> Result
	6, // 9: FinanceService.RiseSettleByDay:output_type -> Result
	6, // 10: FinanceService.CommitTransfer:output_type -> Result
	6, // 11: FinanceService.OpenRiseService:output_type -> Result
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_finance_service_proto_init() }
func file_finance_service_proto_init() {
	if File_finance_service_proto != nil {
		return
	}
	file_global_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_finance_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferInRequest); i {
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
		file_finance_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersonId); i {
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
		file_finance_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SRiseInfo); i {
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
		file_finance_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RiseSettleRequest); i {
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
		file_finance_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RiseTransferOutRequest); i {
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
		file_finance_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitTransferRequest); i {
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
			RawDescriptor: file_finance_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_finance_service_proto_goTypes,
		DependencyIndexes: file_finance_service_proto_depIdxs,
		MessageInfos:      file_finance_service_proto_msgTypes,
	}.Build()
	File_finance_service_proto = out.File
	file_finance_service_proto_rawDesc = nil
	file_finance_service_proto_goTypes = nil
	file_finance_service_proto_depIdxs = nil
}
