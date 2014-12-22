// 合作商户的接口
package server

import (
	"com/ording/dao"
	"com/ording/entity"
	"ops/cf/net/jsv"
)

type Share struct{}

func (this *Share) GetShoppingCart(m *jsv.Args, r *jsv.Result) error {
	var cartData string = (*m)["cart"].(string)

	cart, err := dao.ParseShoppingCart(cartData)
	if err != nil {
		return err
	}

	totalFee, orderFee := cart.GetFee()

	r.Data = entity.ShoppingCart{
		TotalFee: totalFee,
		OrderFee: orderFee,
		Summary:  cart.GetSummary(),
	}
	r.Result = true

	return nil
}
