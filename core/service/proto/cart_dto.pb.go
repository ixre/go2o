// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.4
// source: message/cart_dto.proto

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

// todo: 废弃
type WsCartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MemberId int64             `protobuf:"zigzag64,1,opt,name=memberId,proto3" json:"memberId"`
	Action   string            `protobuf:"bytes,2,opt,name=action,proto3" json:"action"`
	Data     map[string]string `protobuf:"bytes,3,rep,name=data,proto3" json:"data" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *WsCartRequest) Reset() {
	*x = WsCartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WsCartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WsCartRequest) ProtoMessage() {}

func (x *WsCartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WsCartRequest.ProtoReflect.Descriptor instead.
func (*WsCartRequest) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{0}
}

func (x *WsCartRequest) GetMemberId() int64 {
	if x != nil {
		return x.MemberId
	}
	return 0
}

func (x *WsCartRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *WsCartRequest) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

type NormalCartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MemberId int64             `protobuf:"zigzag64,1,opt,name=memberId,proto3" json:"memberId"`
	Action   string            `protobuf:"bytes,2,opt,name=action,proto3" json:"action"`
	Data     map[string]string `protobuf:"bytes,3,rep,name=data,proto3" json:"data" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *NormalCartRequest) Reset() {
	*x = NormalCartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NormalCartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NormalCartRequest) ProtoMessage() {}

func (x *NormalCartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NormalCartRequest.ProtoReflect.Descriptor instead.
func (*NormalCartRequest) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{1}
}

func (x *NormalCartRequest) GetMemberId() int64 {
	if x != nil {
		return x.MemberId
	}
	return 0
}

func (x *NormalCartRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *NormalCartRequest) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

// 购物车
type SShoppingCart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 编号
	CartId int32 `protobuf:"zigzag32,1,opt,name=cartId,proto3" json:"cartId"`
	// 购物车KEY
	CartCode string `protobuf:"bytes,2,opt,name=cartCode,proto3" json:"cartCode"`
	// 店铺分组
	Sellers []*SShoppingCartGroup `protobuf:"bytes,3,rep,name=sellers,proto3" json:"sellers"`
}

func (x *SShoppingCart) Reset() {
	*x = SShoppingCart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SShoppingCart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SShoppingCart) ProtoMessage() {}

func (x *SShoppingCart) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SShoppingCart.ProtoReflect.Descriptor instead.
func (*SShoppingCart) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{2}
}

func (x *SShoppingCart) GetCartId() int32 {
	if x != nil {
		return x.CartId
	}
	return 0
}

func (x *SShoppingCart) GetCartCode() string {
	if x != nil {
		return x.CartCode
	}
	return ""
}

func (x *SShoppingCart) GetSellers() []*SShoppingCartGroup {
	if x != nil {
		return x.Sellers
	}
	return nil
}

// 购物车店铺分组
type SShoppingCartGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 供货商编号
	SellerId int64 `protobuf:"zigzag64,1,opt,name=sellerId,proto3" json:"sellerId"`
	// 店铺编号
	ShopId int64 `protobuf:"zigzag64,2,opt,name=shopId,proto3" json:"shopId"`
	// 店铺名称
	ShopName string `protobuf:"bytes,3,opt,name=shopName,proto3" json:"shopName"`
	// 是否结算
	Checked bool `protobuf:"varint,4,opt,name=checked,proto3" json:"checked"`
	// 商品
	Items []*SShoppingCartItem `protobuf:"bytes,5,rep,name=items,proto3" json:"items"`
}

func (x *SShoppingCartGroup) Reset() {
	*x = SShoppingCartGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SShoppingCartGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SShoppingCartGroup) ProtoMessage() {}

func (x *SShoppingCartGroup) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SShoppingCartGroup.ProtoReflect.Descriptor instead.
func (*SShoppingCartGroup) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{3}
}

func (x *SShoppingCartGroup) GetSellerId() int64 {
	if x != nil {
		return x.SellerId
	}
	return 0
}

func (x *SShoppingCartGroup) GetShopId() int64 {
	if x != nil {
		return x.ShopId
	}
	return 0
}

func (x *SShoppingCartGroup) GetShopName() string {
	if x != nil {
		return x.ShopName
	}
	return ""
}

func (x *SShoppingCartGroup) GetChecked() bool {
	if x != nil {
		return x.Checked
	}
	return false
}

func (x *SShoppingCartGroup) GetItems() []*SShoppingCartItem {
	if x != nil {
		return x.Items
	}
	return nil
}

// 购物车商品勾选
type SCheckCartItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 商品编号
	ItemId int64 `protobuf:"zigzag64,1,opt,name=itemId,proto3" json:"itemId"`
	// SKU编号
	SkuId int64 `protobuf:"zigzag64,2,opt,name=skuId,proto3" json:"skuId"`
	// 是否勾选
	Checked bool `protobuf:"varint,3,opt,name=checked,proto3" json:"checked"`
}

func (x *SCheckCartItem) Reset() {
	*x = SCheckCartItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCheckCartItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCheckCartItem) ProtoMessage() {}

func (x *SCheckCartItem) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCheckCartItem.ProtoReflect.Descriptor instead.
func (*SCheckCartItem) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{4}
}

func (x *SCheckCartItem) GetItemId() int64 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *SCheckCartItem) GetSkuId() int64 {
	if x != nil {
		return x.SkuId
	}
	return 0
}

func (x *SCheckCartItem) GetChecked() bool {
	if x != nil {
		return x.Checked
	}
	return false
}

// 购物车商品操作响应
type CartItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrCode int32              `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode"`
	ErrMsg  string             `protobuf:"bytes,2,opt,name=errMsg,proto3" json:"errMsg"`
	Item    *SShoppingCartItem `protobuf:"bytes,3,opt,name=item,proto3" json:"item"`
}

func (x *CartItemResponse) Reset() {
	*x = CartItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CartItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CartItemResponse) ProtoMessage() {}

func (x *CartItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CartItemResponse.ProtoReflect.Descriptor instead.
func (*CartItemResponse) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{5}
}

func (x *CartItemResponse) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *CartItemResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *CartItemResponse) GetItem() *SShoppingCartItem {
	if x != nil {
		return x.Item
	}
	return nil
}

// 购物车商品
type SShoppingCartItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 商品编号
	ItemId int64 `protobuf:"zigzag64,1,opt,name=itemId,proto3" json:"itemId"`
	// SKU编号
	SkuId int64 `protobuf:"zigzag64,2,opt,name=skuId,proto3" json:"skuId"`
	// 商品标题
	Title string `protobuf:"bytes,3,opt,name=title,proto3" json:"title"`
	// 商品图片
	Image string `protobuf:"bytes,4,opt,name=image,proto3" json:"image"`
	// 规格文本
	SpecWord string `protobuf:"bytes,5,opt,name=specWord,proto3" json:"specWord"`
	// 商品编码
	Code string `protobuf:"bytes,6,opt,name=code,proto3" json:"code"`
	// 零售价
	RetailPrice float64 `protobuf:"fixed64,7,opt,name=retailPrice,proto3" json:"retailPrice"`
	// 销售价
	Price float64 `protobuf:"fixed64,8,opt,name=price,proto3" json:"price"`
	// 数量
	Quantity int32 `protobuf:"zigzag32,9,opt,name=quantity,proto3" json:"quantity"`
	// 是否结算
	Checked bool `protobuf:"varint,10,opt,name=checked,proto3" json:"checked"`
	// 库存文本
	StockText string `protobuf:"bytes,11,opt,name=stockText,proto3" json:"stockText"`
	// 店铺编号
	ShopId int64 `protobuf:"zigzag64,12,opt,name=shopId,proto3" json:"shopId"`
}

func (x *SShoppingCartItem) Reset() {
	*x = SShoppingCartItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SShoppingCartItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SShoppingCartItem) ProtoMessage() {}

func (x *SShoppingCartItem) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SShoppingCartItem.ProtoReflect.Descriptor instead.
func (*SShoppingCartItem) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{6}
}

func (x *SShoppingCartItem) GetItemId() int64 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *SShoppingCartItem) GetSkuId() int64 {
	if x != nil {
		return x.SkuId
	}
	return 0
}

func (x *SShoppingCartItem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *SShoppingCartItem) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *SShoppingCartItem) GetSpecWord() string {
	if x != nil {
		return x.SpecWord
	}
	return ""
}

func (x *SShoppingCartItem) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *SShoppingCartItem) GetRetailPrice() float64 {
	if x != nil {
		return x.RetailPrice
	}
	return 0
}

func (x *SShoppingCartItem) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *SShoppingCartItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *SShoppingCartItem) GetChecked() bool {
	if x != nil {
		return x.Checked
	}
	return false
}

func (x *SShoppingCartItem) GetStockText() string {
	if x != nil {
		return x.StockText
	}
	return ""
}

func (x *SShoppingCartItem) GetShopId() int64 {
	if x != nil {
		return x.ShopId
	}
	return 0
}

// 购物车商品请求
type SCartItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 商品编号
	ItemId int64 `protobuf:"zigzag64,2,opt,name=itemId,proto3" json:"itemId"`
	// SKU编号
	SkuId int64 `protobuf:"zigzag64,3,opt,name=skuId,proto3" json:"skuId"`
	// 数量
	Quantity int32 `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity"`
	// 是否只勾选
	Checked bool `protobuf:"varint,5,opt,name=checked,proto3" json:"checked"`
}

func (x *SCartItem) Reset() {
	*x = SCartItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_cart_dto_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SCartItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SCartItem) ProtoMessage() {}

func (x *SCartItem) ProtoReflect() protoreflect.Message {
	mi := &file_message_cart_dto_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SCartItem.ProtoReflect.Descriptor instead.
func (*SCartItem) Descriptor() ([]byte, []int) {
	return file_message_cart_dto_proto_rawDescGZIP(), []int{7}
}

func (x *SCartItem) GetItemId() int64 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *SCartItem) GetSkuId() int64 {
	if x != nil {
		return x.SkuId
	}
	return 0
}

func (x *SCartItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *SCartItem) GetChecked() bool {
	if x != nil {
		return x.Checked
	}
	return false
}

var File_message_cart_dto_proto protoreflect.FileDescriptor

var file_message_cart_dto_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x63, 0x61, 0x72, 0x74, 0x5f, 0x64,
	0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x01, 0x0a, 0x0d, 0x57, 0x73, 0x43,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x08, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2c,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x57,
	0x73, 0x43, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09,
	0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xb2, 0x01, 0x0a, 0x11, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c,
	0x43, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x08, 0x6d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x30, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x43, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x72, 0x0a, 0x0d, 0x53, 0x53,
	0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63,
	0x61, 0x72, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x11, 0x52, 0x06, 0x63, 0x61, 0x72,
	0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x2d, 0x0a, 0x07, 0x73, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x53, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x72, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x07, 0x73, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x73, 0x22, 0xa8,
	0x01, 0x0a, 0x12, 0x53, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x72, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x08, 0x73, 0x65, 0x6c, 0x6c, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x12, 0x52, 0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x68, 0x6f,
	0x70, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f,
	0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64, 0x12,
	0x28, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x53, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x72, 0x74, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x58, 0x0a, 0x0e, 0x53, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x43, 0x61, 0x72, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x69,
	0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x69, 0x74, 0x65,
	0x6d, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x6b, 0x75, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x12, 0x52, 0x05, 0x73, 0x6b, 0x75, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x65, 0x64, 0x22, 0x6c, 0x0a, 0x10, 0x43, 0x61, 0x72, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x26, 0x0a, 0x04, 0x69, 0x74, 0x65,
	0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x53, 0x53, 0x68, 0x6f, 0x70, 0x70,
	0x69, 0x6e, 0x67, 0x43, 0x61, 0x72, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x69, 0x74, 0x65,
	0x6d, 0x22, 0xc1, 0x02, 0x0a, 0x11, 0x53, 0x53, 0x68, 0x6f, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43,
	0x61, 0x72, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x6b, 0x75, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x12, 0x52, 0x05,
	0x73, 0x6b, 0x75, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x70, 0x65, 0x63, 0x57, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x70, 0x65, 0x63, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x72, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x11, 0x52, 0x08, 0x71, 0x75, 0x61,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x54, 0x65, 0x78, 0x74, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x54, 0x65, 0x78, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x12, 0x52, 0x06, 0x73,
	0x68, 0x6f, 0x70, 0x49, 0x64, 0x22, 0x6f, 0x0a, 0x09, 0x53, 0x43, 0x61, 0x72, 0x74, 0x49, 0x74,
	0x65, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x12, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x6b,
	0x75, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x12, 0x52, 0x05, 0x73, 0x6b, 0x75, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x65, 0x64, 0x42, 0x1f, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72, 0x70, 0x63, 0x5a, 0x08, 0x2e,
	0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_cart_dto_proto_rawDescOnce sync.Once
	file_message_cart_dto_proto_rawDescData = file_message_cart_dto_proto_rawDesc
)

func file_message_cart_dto_proto_rawDescGZIP() []byte {
	file_message_cart_dto_proto_rawDescOnce.Do(func() {
		file_message_cart_dto_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_cart_dto_proto_rawDescData)
	})
	return file_message_cart_dto_proto_rawDescData
}

var file_message_cart_dto_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_message_cart_dto_proto_goTypes = []interface{}{
	(*WsCartRequest)(nil),      // 0: WsCartRequest
	(*NormalCartRequest)(nil),  // 1: NormalCartRequest
	(*SShoppingCart)(nil),      // 2: SShoppingCart
	(*SShoppingCartGroup)(nil), // 3: SShoppingCartGroup
	(*SCheckCartItem)(nil),     // 4: SCheckCartItem
	(*CartItemResponse)(nil),   // 5: CartItemResponse
	(*SShoppingCartItem)(nil),  // 6: SShoppingCartItem
	(*SCartItem)(nil),          // 7: SCartItem
	nil,                        // 8: WsCartRequest.DataEntry
	nil,                        // 9: NormalCartRequest.DataEntry
}
var file_message_cart_dto_proto_depIdxs = []int32{
	8, // 0: WsCartRequest.data:type_name -> WsCartRequest.DataEntry
	9, // 1: NormalCartRequest.data:type_name -> NormalCartRequest.DataEntry
	3, // 2: SShoppingCart.sellers:type_name -> SShoppingCartGroup
	6, // 3: SShoppingCartGroup.items:type_name -> SShoppingCartItem
	6, // 4: CartItemResponse.item:type_name -> SShoppingCartItem
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_message_cart_dto_proto_init() }
func file_message_cart_dto_proto_init() {
	if File_message_cart_dto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_cart_dto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WsCartRequest); i {
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
		file_message_cart_dto_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NormalCartRequest); i {
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
		file_message_cart_dto_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SShoppingCart); i {
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
		file_message_cart_dto_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SShoppingCartGroup); i {
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
		file_message_cart_dto_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCheckCartItem); i {
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
		file_message_cart_dto_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CartItemResponse); i {
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
		file_message_cart_dto_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SShoppingCartItem); i {
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
		file_message_cart_dto_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SCartItem); i {
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
			RawDescriptor: file_message_cart_dto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_cart_dto_proto_goTypes,
		DependencyIndexes: file_message_cart_dto_proto_depIdxs,
		MessageInfos:      file_message_cart_dto_proto_msgTypes,
	}.Build()
	File_message_cart_dto_proto = out.File
	file_message_cart_dto_proto_rawDesc = nil
	file_message_cart_dto_proto_goTypes = nil
	file_message_cart_dto_proto_depIdxs = nil
}
