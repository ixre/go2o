package order

import (
	"fmt"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/order"
	id "go2o/core/infrastructure/domain"
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

func (p *postedDataImpl) GetComment(sellerId int32) string {
	k := fmt.Sprintf("seller_comment_%d", sellerId)
	return p.data[k]
}
