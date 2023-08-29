package order

import (
	"errors"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/order"
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

func (w *wholesaleOrderBreaker) BreakUp(c cart.ICartAggregateRoot,
	data order.IPostedData) ([]order.IOrder, error) {
	switch c.Kind() {
	case cart.KWholesale:
		return w.breakupWholesaleOrder(c, data)
	}
	return []order.IOrder{}, errors.New("not support cart kind")
}

func (w *wholesaleOrderBreaker) breakupWholesaleOrder(c cart.ICartAggregateRoot,
	data order.IPostedData) ([]order.IOrder, error) {
	checked := data.CheckedData()
	items := c.CheckedItems(checked)
	if len(items) == 0 {
		return []order.IOrder{}, order.ErrNoCheckedItem
	}
	// 将购物车的商品按运营商拆分
	vendorItemsMap := w.breakSellerItemMap(items)
	if l := len(vendorItemsMap); l == 0 {
		return []order.IOrder{}, cart.ErrNoChecked
	}
	var list []order.IOrder
	cc := c.(cart.ICartAggregateRoot)
	buyerId := cc.BuyerId()
	for sellerId, items := range vendorItemsMap {
		o := w.createWholesaleOrder(sellerId, buyerId, items, data)
		list = append(list, o)
	}
	return list, nil
}

// 创建批发订单
func (w *wholesaleOrderBreaker) createWholesaleOrder(sellerId int64,
	buyerId int64, items []*cart.ItemPair, data order.IPostedData) order.IOrder {
	v := &order.Order{
		BuyerId:   buyerId,
		OrderType: int(order.TWholesale),
		Status:    order.StatAwaitingPayment,
	}
	o := w.repo.CreateOrder(v)
	wo := o.(order.IWholesaleOrder)
	wo.SetItems(items)
	wo.SetComment(data.GetComment(sellerId))
	_ = o.SetShipmentAddress(data.AddressId())
	return o
}

// 生成运营商与订单商品的映射
func (w *wholesaleOrderBreaker) breakSellerItemMap(items []*cart.ItemPair) map[int64][]*cart.ItemPair {
	mp := make(map[int64][]*cart.ItemPair)
	for _, v := range items {
		list, ok := mp[v.SellerId]
		if !ok {
			list = []*cart.ItemPair{}
		}
		mp[v.SellerId] = append(list, v)
	}
	return mp
}
