package parser

import (
	"fmt"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/order"
	id "github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/util"
)

var _ order.IPostedData = new(postedDataImpl)

type postedDataImpl struct {
	data       map[string]string
	_addressId int64
	_checked   map[int64][]int64
	_req       *proto.SubmitOrderRequest
}

func NewPostedData(data map[string]string, req *proto.SubmitOrderRequest) order.IPostedData {
	return &postedDataImpl{data: data, _req: req}
}

// 转换勾选字典,数据如：{"1":["10","11"],"2":["20","21"]}
func (p *postedDataImpl) CheckedData() map[int64][]int64 {
	if p._checked == nil {
		p._checked = id.ParseCartCheckedMap(p.data["checked"])
	}
	return p._checked
}

func (p *postedDataImpl) AddressId() int64 {
	if p._addressId == 0 {
		p._addressId, _ = util.I64Err(strconv.Atoi(p.data["address_id"]))
	}
	return p._addressId
}

func (p *postedDataImpl) GetComment(sellerId int64) string {
	k := fmt.Sprintf("seller_comment_%d", sellerId)
	return p.data[k]
}

// TradeOrderAmount implements order.IPostedData
func (p *postedDataImpl) TradeOrderAmount() int64 {
	if p._req.TradeOrder != nil {
		return p._req.TradeOrder.TradeAmount
	}
	return 0
}

// TradeOrderDiscount implements order.IPostedData
func (p *postedDataImpl) TradeOrderDiscount() float32 {
	if p._req.TradeOrder != nil {
		return p._req.TradeOrder.Discount
	}
	return 0
}

// TradeOrderStoreId implements order.IPostedData
func (p *postedDataImpl) TradeOrderStoreId() int64 {
	if p._req.TradeOrder != nil {
		return p._req.TradeOrder.StoreId
	}
	return 0
}
