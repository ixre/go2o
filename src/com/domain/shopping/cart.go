package shopping

import (
	"bytes"
	"com/domain/interface/shopping"
	"strconv"
)

type Cart struct {
	val     *shopping.ValueCart
	summary string
}

func newCart(val *shopping.ValueCart) shopping.ICart {
	return &Cart{
		val: val,
	}
}

func (this *Cart) GetValue() shopping.ValueCart {
	return *this.val
}

// todo: 可能将购物车存放于会员信息中
func (this *Cart) GetDomainId() int {
	return -1
}

// 获取总览信息
func (this *Cart) GetSummary() string {
	if len(this.summary) != 0 {
		return this.summary
	}
	buf := bytes.NewBufferString("")
	length := len(this.val.Items)
	for i, pro := range this.val.Items {
		v := pro.GetValue()
		buf.WriteString(v.Name)
		if len(v.Note) != 0 {
			buf.WriteString("(" + v.Note + ")")
		}
		buf.WriteString("*" + strconv.Itoa(this.val.Quantities[v.Id]))
		if i < length-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// 获取订单金额,返回totalFee为总额，
// orderFee为实际订单的金额(扣去促销优惠等后的金额)
func (this *Cart) GetFee() (totalFee float32, orderFee float32) {
	var qua float32
	for _, pro := range this.val.Items {
		v := pro.GetValue()
		qua = float32(this.val.Quantities[v.Id])
		totalFee = totalFee + v.Price*qua
		orderFee = orderFee + v.SalePrice*qua
	}
	return totalFee, orderFee
}
