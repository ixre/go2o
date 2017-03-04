/**
 * Copyright 2015 @ z3q.net.
 * name : order_test.go
 * author : jarryliu
 * date : 2016-07-15 15:14
 * description :
 * history :
 */
package testing

import (
	"go2o/core/domain/interface/cart"
	"go2o/core/domain/interface/order"
	"go2o/core/domain/interface/payment"
	"go2o/core/testing/ti"
	"testing"
	"time"
)

/*

清理订单数据：

TRUNCATE `pay_order`;
TRUNCATE `ship_item`;
TRUNCATE `ship_order`;
TRUNCATE `sale_cart`;
TRUNCATE `sale_cart_item`;
TRUNCATE `sale_order`;
TRUNCATE `sale_sub_order`;
TRUNCATE `sale_order_item`;
TRUNCATE `sale_order_log`;
TRUNCATE `sale_refund`;
TRUNCATE `sale_return`;
TRUNCATE `sale_exchange`;
TRUNCATE `order_list`;
TRUNCATE `order_wholesale_item`;
TRUNCATE `order_wholesale_order`;

*/

func logState(t *testing.T, err error, o order.IOrder) {
	if err != nil {
		t.Log(err)
		t.FailNow()
	} else {
		t.Log(order.OrderState(o.State()).String())
	}
}

func TestOrderSetup(t *testing.T) {
	orderNo := "100000735578"
	orderRepo := ti.OrderRepo
	v := orderRepo.GetSubOrderByNo(orderNo)
	o := orderRepo.Manager().GetSubOrder(v.ID)

	t.Log("-[ 订单状态为:" + order.OrderState(o.GetValue().State).String())

	err := o.PaymentFinishByOnlineTrade()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderState(o.GetValue().State).String())
	}

	err = o.Confirm()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderState(o.GetValue().State).String())
	}

	err = o.PickUp()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderState(o.GetValue().State).String())
	}

	err = o.Ship(1, "100000")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderState(o.GetValue().State).String())
	}

	return
	err = o.BuyerReceived()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderState(o.GetValue().State).String())
	}
}

func TestCancelOrder(t *testing.T) {
	repo := ti.CartRepo
	var buyerId int32 = 1
	c := repo.GetMyCart(buyerId, cart.KRetail)
	joinItemsToCart(c, t)
	if c.Kind() == cart.KRetail {
		rc := c.(cart.IRetailCart)
		t.Log("购物车如下:")
		for _, v := range rc.Items() {
			t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
		}
	}
	_, err := c.Save()
	if err != nil {
		t.Error("保存购物车失败:", err.Error())
		t.FailNow()
	}

	orderRepo := ti.OrderRepo
	mmRepo := ti.MemberRepo
	manager := orderRepo.Manager()
	m := mmRepo.GetMember(buyerId)
	addressId := m.Profile().GetDefaultAddress().GetDomainId()
	o, err := manager.SubmitOrder(c, addressId, "", !true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	py := o.(order.INormalOrder).GetPaymentOrder()
	//err = py.PaymentByWallet("支付订单")
	pv := py.GetValue()
	payState := pv.State
	if payState == payment.StateFinishPayment {
		t.Logf("订单支付完成,金额：%.2f", pv.FinalAmount)
	} else {
		t.Logf("订单未完成支付,状态：%d;订单号：%s", pv.State, py.GetTradeNo())
	}
	t.Logf("支付单信息：%#v", pv)
	//t.Log("调价：",py.Adjust(-pv.FinalAmount))
	//t.Log(py.Cancel())
	//return
	time.Sleep(time.Second * 2)

	io := o.(order.INormalOrder)

	subs := io.GetSubOrders()
	for _, v := range subs {
		err = v.Cancel("买多了，不想要了!")
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}
	t.Log("退货成功")
}

// 测试批发订单
func TestWholesaleOrder(t *testing.T) {
	repo := ti.CartRepo
	var buyerId int32 = 1
	c := repo.GetMyCart(buyerId, cart.KWholesale)
	joinItemsToCart(c, t)
	rc := c.(cart.IWholesaleCart)

	t.Log("购物车如下:")
	for _, v := range rc.Items() {
		t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
	}
	if len(rc.GetValue().Items) == 0 {
		t.Log("购物车是空的")
		t.FailNow()
	}

	_, err := c.Save()
	if err != nil {
		t.Error("保存购物车失败:", err.Error())
		t.Fail()
	}

	orderRepo := ti.OrderRepo
	manager := orderRepo.Manager()

	orders, err := manager.PrepareWholesaleOrder(rc)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("批发单拆分数量：%d", len(orders))

	//m := mmRepo.GetMember(buyerId)
	//addressId := m.Profile().GetDefaultAddress().GetDomainId()
	//o, py, err := manager.SubmitOrder(c, addressId, "", "", true)
	////err = py.PaymentByWallet("支付订单")
	//pv := py.GetValue()
	//payState := pv.State
	//if payState == payment.StateFinishPayment {
	//    t.Logf("订单支付完成,金额：%.2f", pv.FinalAmount)
	//} else {
	//    t.Logf("订单未完成支付,状态：%d;订单号：%s", pv.State, py.GetTradeNo())
	//}
	//t.Logf("支付单信息：%#v", pv)
	//t.Log("调价：",py.Adjust(-pv.FinalAmount))
	//t.Log(py.Cancel())
	//return
	time.Sleep(time.Second * 2)

	addressId := orders[0].Buyer().Profile().GetDefaultAddress().GetDomainId()

	for _, o := range orders {
		io := o.(order.IWholesaleOrder)
		err = io.SetAddress(addressId)
		if err == nil {
			if err = o.Submit(); err == nil {
				//po := io.GetPaymentOrder()
				//err = po.PaymentFinish("alipay","1234567899000")
			}
		}
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		logState(t, io.Confirm(), o)
		logState(t, io.PickUp(), o)
		logState(t, io.Ship(1, "123456"), o)
		logState(t, io.BuyerReceived(), o)
	}
}
