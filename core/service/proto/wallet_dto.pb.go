// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v5.26.1
// source: message/wallet_dto.proto

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

// 提现方式
type EUserWithdrawalKind int32

const (
	EUserWithdrawalKind____WithdrawKind EUserWithdrawalKind = 0
	// 提现到银行卡
	EUserWithdrawalKind_WithdrawToBankCard EUserWithdrawalKind = 1
	// 提现到第三方账户
	EUserWithdrawalKind_WithdrawToPayWallet EUserWithdrawalKind = 2
	// 提现到自定义账户
	EUserWithdrawalKind_WithdrawCustom EUserWithdrawalKind = 3
	// 兑换为商城余额
	EUserWithdrawalKind_WithdrawByExchange EUserWithdrawalKind = 4
)

// Enum value maps for EUserWithdrawalKind.
var (
	EUserWithdrawalKind_name = map[int32]string{
		0: "___WithdrawKind",
		1: "WithdrawToBankCard",
		2: "WithdrawToPayWallet",
		3: "WithdrawCustom",
		4: "WithdrawByExchange",
	}
	EUserWithdrawalKind_value = map[string]int32{
		"___WithdrawKind":     0,
		"WithdrawToBankCard":  1,
		"WithdrawToPayWallet": 2,
		"WithdrawCustom":      3,
		"WithdrawByExchange":  4,
	}
)

func (x EUserWithdrawalKind) Enum() *EUserWithdrawalKind {
	p := new(EUserWithdrawalKind)
	*p = x
	return p
}

func (x EUserWithdrawalKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EUserWithdrawalKind) Descriptor() protoreflect.EnumDescriptor {
	return file_message_wallet_dto_proto_enumTypes[0].Descriptor()
}

func (EUserWithdrawalKind) Type() protoreflect.EnumType {
	return &file_message_wallet_dto_proto_enumTypes[0]
}

func (x EUserWithdrawalKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EUserWithdrawalKind.Descriptor instead.
func (EUserWithdrawalKind) EnumDescriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{0}
}

// * 账户入账请求
type UserWalletCarryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会员编号
	UserId int64 `protobuf:"zigzag64,1,opt,name=userId,proto3" json:"userId"`
	// 明细标题
	TransactionTitle string `protobuf:"bytes,3,opt,name=transactionTitle,proto3" json:"transactionTitle"`
	// 已扣除手续费的金额
	Amount int64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount"`
	// 手续费
	TransactionFee int64 `protobuf:"varint,5,opt,name=transactionFee,proto3" json:"transactionFee"`
	// 外部校译号
	OuterTransactionNo string `protobuf:"bytes,6,opt,name=outerTransactionNo,proto3" json:"outerTransactionNo"`
	// 备注
	TransactionRemark string `protobuf:"bytes,7,opt,name=transactionRemark,proto3" json:"transactionRemark"`
	// 是否先冻结
	Freeze bool `protobuf:"varint,8,opt,name=freeze,proto3" json:"freeze"`
}

func (x *UserWalletCarryRequest) Reset() {
	*x = UserWalletCarryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserWalletCarryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserWalletCarryRequest) ProtoMessage() {}

func (x *UserWalletCarryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserWalletCarryRequest.ProtoReflect.Descriptor instead.
func (*UserWalletCarryRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{0}
}

func (x *UserWalletCarryRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserWalletCarryRequest) GetTransactionTitle() string {
	if x != nil {
		return x.TransactionTitle
	}
	return ""
}

func (x *UserWalletCarryRequest) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *UserWalletCarryRequest) GetTransactionFee() int64 {
	if x != nil {
		return x.TransactionFee
	}
	return 0
}

func (x *UserWalletCarryRequest) GetOuterTransactionNo() string {
	if x != nil {
		return x.OuterTransactionNo
	}
	return ""
}

func (x *UserWalletCarryRequest) GetTransactionRemark() string {
	if x != nil {
		return x.TransactionRemark
	}
	return ""
}

func (x *UserWalletCarryRequest) GetFreeze() bool {
	if x != nil {
		return x.Freeze
	}
	return false
}

// * 账户调整请求
type UserWalletAdjustRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// * 会员编号
	UserId int64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId"`
	// * 调整金额/数量
	Value int64 `protobuf:"varint,3,opt,name=value,proto3" json:"value"`
	// * 是否人工调整
	ManualAdjust bool `protobuf:"varint,4,opt,name=manualAdjust,proto3" json:"manualAdjust"`
	// * 关联用户
	RelateUser int64 `protobuf:"varint,5,opt,name=relateUser,proto3" json:"relateUser"`
	// * 备注
	TransactionRemark string `protobuf:"bytes,6,opt,name=transactionRemark,proto3" json:"transactionRemark"`
}

func (x *UserWalletAdjustRequest) Reset() {
	*x = UserWalletAdjustRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserWalletAdjustRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserWalletAdjustRequest) ProtoMessage() {}

func (x *UserWalletAdjustRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserWalletAdjustRequest.ProtoReflect.Descriptor instead.
func (*UserWalletAdjustRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{1}
}

func (x *UserWalletAdjustRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserWalletAdjustRequest) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *UserWalletAdjustRequest) GetManualAdjust() bool {
	if x != nil {
		return x.ManualAdjust
	}
	return false
}

func (x *UserWalletAdjustRequest) GetRelateUser() int64 {
	if x != nil {
		return x.RelateUser
	}
	return 0
}

func (x *UserWalletAdjustRequest) GetTransactionRemark() string {
	if x != nil {
		return x.TransactionRemark
	}
	return ""
}

// * 冻结请求
type UserWalletFreezeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会员编号
	UserId int64 `protobuf:"zigzag64,1,opt,name=userId,proto3" json:"userId"`
	// * 标题
	TransactionTitle string `protobuf:"bytes,3,opt,name=transactionTitle,proto3" json:"transactionTitle"`
	// * 冻结金额
	Amount int64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount"`
	// 外部交易号
	OuterTransactionNo string `protobuf:"bytes,5,opt,name=outerTransactionNo,proto3" json:"outerTransactionNo"`
	// 交易流水编号,对冻结流水进行更新时,传递该参数
	TransactionId int64 `protobuf:"varint,6,opt,name=transactionId,proto3" json:"transactionId"`
	// 备注
	TransactionRemark string `protobuf:"bytes,7,opt,name=transactionRemark,proto3" json:"transactionRemark"`
}

func (x *UserWalletFreezeRequest) Reset() {
	*x = UserWalletFreezeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserWalletFreezeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserWalletFreezeRequest) ProtoMessage() {}

func (x *UserWalletFreezeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserWalletFreezeRequest.ProtoReflect.Descriptor instead.
func (*UserWalletFreezeRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{2}
}

func (x *UserWalletFreezeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserWalletFreezeRequest) GetTransactionTitle() string {
	if x != nil {
		return x.TransactionTitle
	}
	return ""
}

func (x *UserWalletFreezeRequest) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *UserWalletFreezeRequest) GetOuterTransactionNo() string {
	if x != nil {
		return x.OuterTransactionNo
	}
	return ""
}

func (x *UserWalletFreezeRequest) GetTransactionId() int64 {
	if x != nil {
		return x.TransactionId
	}
	return 0
}

func (x *UserWalletFreezeRequest) GetTransactionRemark() string {
	if x != nil {
		return x.TransactionRemark
	}
	return ""
}

// * 解冻请求
type UserWalletUnfreezeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会员编号
	UserId int64 `protobuf:"zigzag64,1,opt,name=userId,proto3" json:"userId"`
	// 交易流水号
	TransactionId int64 `protobuf:"zigzag64,2,opt,name=transactionId,proto3" json:"transactionId"`
	// * 标题
	TransactionTitle string `protobuf:"bytes,3,opt,name=transactionTitle,proto3" json:"transactionTitle"`
	// * 冻结金额
	Amount int64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount"`
	// 外部校译号
	OuterTransactionNo string `protobuf:"bytes,5,opt,name=outerTransactionNo,proto3" json:"outerTransactionNo"`
	// 备注
	TransactionRemark string `protobuf:"bytes,6,opt,name=transactionRemark,proto3" json:"transactionRemark"`
}

func (x *UserWalletUnfreezeRequest) Reset() {
	*x = UserWalletUnfreezeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserWalletUnfreezeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserWalletUnfreezeRequest) ProtoMessage() {}

func (x *UserWalletUnfreezeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserWalletUnfreezeRequest.ProtoReflect.Descriptor instead.
func (*UserWalletUnfreezeRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{3}
}

func (x *UserWalletUnfreezeRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserWalletUnfreezeRequest) GetTransactionId() int64 {
	if x != nil {
		return x.TransactionId
	}
	return 0
}

func (x *UserWalletUnfreezeRequest) GetTransactionTitle() string {
	if x != nil {
		return x.TransactionTitle
	}
	return ""
}

func (x *UserWalletUnfreezeRequest) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *UserWalletUnfreezeRequest) GetOuterTransactionNo() string {
	if x != nil {
		return x.OuterTransactionNo
	}
	return ""
}

func (x *UserWalletUnfreezeRequest) GetTransactionRemark() string {
	if x != nil {
		return x.TransactionRemark
	}
	return ""
}

// 提现申请
type UserWithdrawRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会员编号
	UserId int64 `protobuf:"zigzag64,1,opt,name=userId,proto3" json:"userId"`
	// 提现金额
	Amount int64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount"`
	// 手续费
	TransactionFee int64 `protobuf:"varint,3,opt,name=transactionFee,proto3" json:"transactionFee"`
	// 提现方式,21:提现并兑换到余额  22:提现到银行卡(人工提现) 23:第三方钱包
	WithdrawalKind EUserWithdrawalKind `protobuf:"varint,4,opt,name=withdrawalKind,proto3,enum=EUserWithdrawalKind" json:"withdrawalKind"`
	// 银行账号或第三方支付钱包
	AccountNo string `protobuf:"bytes,5,opt,name=accountNo,proto3" json:"accountNo"`
}

func (x *UserWithdrawRequest) Reset() {
	*x = UserWithdrawRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserWithdrawRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserWithdrawRequest) ProtoMessage() {}

func (x *UserWithdrawRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserWithdrawRequest.ProtoReflect.Descriptor instead.
func (*UserWithdrawRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{4}
}

func (x *UserWithdrawRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *UserWithdrawRequest) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *UserWithdrawRequest) GetTransactionFee() int64 {
	if x != nil {
		return x.TransactionFee
	}
	return 0
}

func (x *UserWithdrawRequest) GetWithdrawalKind() EUserWithdrawalKind {
	if x != nil {
		return x.WithdrawalKind
	}
	return EUserWithdrawalKind____WithdrawKind
}

func (x *UserWithdrawRequest) GetAccountNo() string {
	if x != nil {
		return x.AccountNo
	}
	return ""
}

// 申请提现请求
type ReviewUserWithdrawalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会员编号
	UserId int64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId"`
	// 提现申请流水Id
	TransactionId int64 `protobuf:"varint,2,opt,name=transactionId,proto3" json:"transactionId"`
	// 是否通过审核
	Pass bool `protobuf:"varint,3,opt,name=pass,proto3" json:"pass"`
	// 备注
	TransactionRemark string `protobuf:"bytes,4,opt,name=transactionRemark,proto3" json:"transactionRemark"`
}

func (x *ReviewUserWithdrawalRequest) Reset() {
	*x = ReviewUserWithdrawalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReviewUserWithdrawalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReviewUserWithdrawalRequest) ProtoMessage() {}

func (x *ReviewUserWithdrawalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReviewUserWithdrawalRequest.ProtoReflect.Descriptor instead.
func (*ReviewUserWithdrawalRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{5}
}

func (x *ReviewUserWithdrawalRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReviewUserWithdrawalRequest) GetTransactionId() int64 {
	if x != nil {
		return x.TransactionId
	}
	return 0
}

func (x *ReviewUserWithdrawalRequest) GetPass() bool {
	if x != nil {
		return x.Pass
	}
	return false
}

func (x *ReviewUserWithdrawalRequest) GetTransactionRemark() string {
	if x != nil {
		return x.TransactionRemark
	}
	return ""
}

// 完成提现
type FinishUserWithdrawalRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 会员编号
	UserId int64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId"`
	// 提现记录编号
	TransactionId int64 `protobuf:"varint,2,opt,name=transactionId,proto3" json:"transactionId"`
	// 汇款/交易单号
	OuterTransactionNo string `protobuf:"bytes,3,opt,name=outerTransactionNo,proto3" json:"outerTransactionNo"`
}

func (x *FinishUserWithdrawalRequest) Reset() {
	*x = FinishUserWithdrawalRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_wallet_dto_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FinishUserWithdrawalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FinishUserWithdrawalRequest) ProtoMessage() {}

func (x *FinishUserWithdrawalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_wallet_dto_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FinishUserWithdrawalRequest.ProtoReflect.Descriptor instead.
func (*FinishUserWithdrawalRequest) Descriptor() ([]byte, []int) {
	return file_message_wallet_dto_proto_rawDescGZIP(), []int{6}
}

func (x *FinishUserWithdrawalRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *FinishUserWithdrawalRequest) GetTransactionId() int64 {
	if x != nil {
		return x.TransactionId
	}
	return 0
}

func (x *FinishUserWithdrawalRequest) GetOuterTransactionNo() string {
	if x != nil {
		return x.OuterTransactionNo
	}
	return ""
}

var File_message_wallet_dto_proto protoreflect.FileDescriptor

var file_message_wallet_dto_proto_rawDesc = []byte{
	0x0a, 0x18, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x5f, 0x64, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x92, 0x02, 0x0a, 0x16, 0x55,
	0x73, 0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x61, 0x72, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2a, 0x0a,
	0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x26, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x65, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x65, 0x65, 0x12, 0x2e, 0x0a, 0x12, 0x6f, 0x75, 0x74,
	0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x12, 0x2c, 0x0a, 0x11, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x72, 0x65, 0x65, 0x7a,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x66, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x22,
	0xb9, 0x01, 0x0a, 0x17, 0x55, 0x73, 0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x41, 0x64,
	0x6a, 0x75, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6d, 0x61, 0x6e,
	0x75, 0x61, 0x6c, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0c, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x41, 0x64, 0x6a, 0x75, 0x73, 0x74, 0x12, 0x1e, 0x0a,
	0x0a, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2c, 0x0a,
	0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61,
	0x72, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0xf9, 0x01, 0x0a, 0x17,
	0x55, 0x73, 0x65, 0x72, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x46, 0x72, 0x65, 0x65, 0x7a, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x2a, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x12, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x12, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x4e, 0x6f, 0x12, 0x24, 0x0a, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0xfb, 0x01, 0x0a, 0x19, 0x55, 0x73, 0x65, 0x72,
	0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x55, 0x6e, 0x66, 0x72, 0x65, 0x65, 0x7a, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a,
	0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x12, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x12, 0x6f, 0x75, 0x74, 0x65, 0x72,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x12, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x12, 0x2c, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0xc9, 0x01, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72, 0x57, 0x69,
	0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x26, 0x0a,
	0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x65, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x65, 0x65, 0x12, 0x3c, 0x0a, 0x0e, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61,
	0x77, 0x61, 0x6c, 0x4b, 0x69, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e,
	0x45, 0x55, 0x73, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x4b,
	0x69, 0x6e, 0x64, 0x52, 0x0e, 0x77, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x4b,
	0x69, 0x6e, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x6f,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e,
	0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x1b, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x55, 0x73, 0x65, 0x72,
	0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x70,
	0x61, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6d, 0x61, 0x72,
	0x6b, 0x22, 0x8b, 0x01, 0x0a, 0x1b, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x55, 0x73, 0x65, 0x72,
	0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x2e, 0x0a, 0x12, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x6f, 0x75, 0x74,
	0x65, 0x72, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x6f, 0x2a,
	0x87, 0x01, 0x0a, 0x13, 0x45, 0x55, 0x73, 0x65, 0x72, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61,
	0x77, 0x61, 0x6c, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x13, 0x0a, 0x0f, 0x5f, 0x5f, 0x5f, 0x57, 0x69,
	0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x4b, 0x69, 0x6e, 0x64, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12,
	0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x54, 0x6f, 0x42, 0x61, 0x6e, 0x6b, 0x43, 0x61,
	0x72, 0x64, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77,
	0x54, 0x6f, 0x50, 0x61, 0x79, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x10, 0x02, 0x12, 0x12, 0x0a,
	0x0e, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x10,
	0x03, 0x12, 0x16, 0x0a, 0x12, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x42, 0x79, 0x45,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x10, 0x04, 0x42, 0x1f, 0x0a, 0x13, 0x63, 0x6f, 0x6d,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72, 0x70, 0x63,
	0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_message_wallet_dto_proto_rawDescOnce sync.Once
	file_message_wallet_dto_proto_rawDescData = file_message_wallet_dto_proto_rawDesc
)

func file_message_wallet_dto_proto_rawDescGZIP() []byte {
	file_message_wallet_dto_proto_rawDescOnce.Do(func() {
		file_message_wallet_dto_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_wallet_dto_proto_rawDescData)
	})
	return file_message_wallet_dto_proto_rawDescData
}

var file_message_wallet_dto_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_wallet_dto_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_message_wallet_dto_proto_goTypes = []interface{}{
	(EUserWithdrawalKind)(0),            // 0: EUserWithdrawalKind
	(*UserWalletCarryRequest)(nil),      // 1: UserWalletCarryRequest
	(*UserWalletAdjustRequest)(nil),     // 2: UserWalletAdjustRequest
	(*UserWalletFreezeRequest)(nil),     // 3: UserWalletFreezeRequest
	(*UserWalletUnfreezeRequest)(nil),   // 4: UserWalletUnfreezeRequest
	(*UserWithdrawRequest)(nil),         // 5: UserWithdrawRequest
	(*ReviewUserWithdrawalRequest)(nil), // 6: ReviewUserWithdrawalRequest
	(*FinishUserWithdrawalRequest)(nil), // 7: FinishUserWithdrawalRequest
}
var file_message_wallet_dto_proto_depIdxs = []int32{
	0, // 0: UserWithdrawRequest.withdrawalKind:type_name -> EUserWithdrawalKind
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_message_wallet_dto_proto_init() }
func file_message_wallet_dto_proto_init() {
	if File_message_wallet_dto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_wallet_dto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserWalletCarryRequest); i {
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
		file_message_wallet_dto_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserWalletAdjustRequest); i {
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
		file_message_wallet_dto_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserWalletFreezeRequest); i {
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
		file_message_wallet_dto_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserWalletUnfreezeRequest); i {
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
		file_message_wallet_dto_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserWithdrawRequest); i {
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
		file_message_wallet_dto_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReviewUserWithdrawalRequest); i {
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
		file_message_wallet_dto_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FinishUserWithdrawalRequest); i {
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
			RawDescriptor: file_message_wallet_dto_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_wallet_dto_proto_goTypes,
		DependencyIndexes: file_message_wallet_dto_proto_depIdxs,
		EnumInfos:         file_message_wallet_dto_proto_enumTypes,
		MessageInfos:      file_message_wallet_dto_proto_msgTypes,
	}.Build()
	File_message_wallet_dto_proto = out.File
	file_message_wallet_dto_proto_rawDesc = nil
	file_message_wallet_dto_proto_goTypes = nil
	file_message_wallet_dto_proto_depIdxs = nil
}
