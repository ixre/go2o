package order

import (
	"fmt"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/order"
	"strconv"
)

var _ order.IPostedData = new(postedDataImpl)

type postedDataImpl struct {
	data       map[string]string
	_addressId int64
	_checked   map[int64][]int64
}

func NewPostedData(data map[string]string) order.IPostedData {
	return &postedDataImpl{data: data}
}

func (p *postedDataImpl) CheckedData() map[int64][]int64 {
	if p._checked == nil {
		p._checked = cart.ParseCheckedMap(p.data["checked"])
	}
	return p._checked
}

func (p *postedDataImpl) AddressId() int64 {
	if p._addressId == 0 {
		p._addressId, _ = util.I64Err(strconv.Atoi(p.data["address_id"]))
	}
	return p._addressId
}

func (p *postedDataImpl) GetComment(sellerId int32) string {
	k := fmt.Sprintf("seller_comment_%d", sellerId)
	return p.data[k]
}
