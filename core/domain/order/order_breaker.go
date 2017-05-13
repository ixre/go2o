package order

import (
	"errors"
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/order"
)

// 订单拆单
type wholesaleOrderBreaker struct {
	repo order.IOrderRepo
}

func newWholesaleOrderBreaker(repo order.IOrderRepo) *wholesaleOrderBreaker {
	return &wholesaleOrderBreaker{
		repo: repo,
	}
}

func (w *wholesaleOrderBreaker) BreakUp(c cart.ICart,
	data map[string]string) ([]order.IOrder, error) {
	checked := cart.ParseCheckedMap(data["checked"])
	switch c.Kind() {
	case cart.KWholesale:
		return w.breakupWholesaleOrder(c, checked)
	}
	return []order.IOrder{}, errors.New("not support cart kind")
}

func (w *wholesaleOrderBreaker) breakupWholesaleOrder(c cart.ICart,
	checked map[int64][]int64) ([]order.IOrder, error) {
	wc := c.(cart.IWholesaleCart)
	items := wc.CheckedItems(checked)
	if len(items) == 0 {
		return []order.IOrder{}, order.ErrNoCheckedItem
	}
	// 将购物车的商品按运营商拆分
	vendorItemsMap := w.breakSellerItemMap(items)
	if l := len(vendorItemsMap); l == 0 {
		return []order.IOrder{}, cart.ErrNoChecked
	}
	list := []order.IOrder{}
	cc := c.(cart.ICart)
	buyerId := cc.BuyerId()
	for _, items := range vendorItemsMap {
		o := w.createWholesaleOrder(buyerId, items)
		list = append(list, o)
	}
	return list, nil
}

// 创建批发订单
func (w *wholesaleOrderBreaker) createWholesaleOrder(buyerId int64, items []*cart.WsCartItem) order.IOrder {
	v := &order.Order{
		BuyerId:   buyerId,
		OrderType: int32(order.TWholesale),
		State:     int32(order.StatAwaitingPayment),
	}
	o := w.repo.CreateOrder(v)
	wo := o.(order.IWholesaleOrder)
	list := make([]*order.MinifyItem, len(items))
	for i, v := range items {
		list[i] = &order.MinifyItem{
			ItemId:   v.ItemId,
			SkuId:    v.SkuId,
			Quantity: v.Quantity,
		}
	}
	wo.SetItems(list)
	return o
}

// 生成运营商与订单商品的映射
func (w *wholesaleOrderBreaker) breakSellerItemMap(items []*cart.WsCartItem) map[int32][]*cart.WsCartItem {
	mp := make(map[int32][]*cart.WsCartItem)
	for _, v := range items {
		list, ok := mp[v.SellerId]
		if !ok {
			list = []*cart.WsCartItem{}
		}
		mp[v.SellerId] = append(list, v)
	}
	return mp
}
