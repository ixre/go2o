/**
 * Copyright 2015 @ 56x.net.
 * name : order_test.go
 * author : jarryliu
 * date : 2016-07-15 15:14
 * description :
 * history :
 */
package domain

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/payment"
	"github.com/ixre/go2o/core/service/parser"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/storage"
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
		t.Log(order.OrderStatus(o.State()).String())
	}
}

func TestOrderSetup(t *testing.T) {
	orderNo := "1230213000412497"
	orderRepo := ti.Factory.GetOrderRepo()
	orderId := orderRepo.GetOrderId(orderNo, true)
	o := orderRepo.Manager().GetSubOrder(orderId)

	t.Logf("order:%#v", o.GetValue())
	t.Log("-[ 订单状态为:" + order.OrderStatus(o.GetValue().Status).String())

	err := o.PaymentFinishByOnlineTrade()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderStatus(o.GetValue().Status).String())
	}
	return
	err = o.PickUp()
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderStatus(o.GetValue().Status).String())
	}

	err = o.Ship(1, "100000")
	if err != nil {
		t.Log(err)
	} else {
		t.Log(order.OrderStatus(o.GetValue().Status).String())
	}

	err = o.BuyerReceived()
	if err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		t.Log(order.OrderStatus(o.GetValue().Status).String())
	}
	for {

	}
	time.Sleep(20000)
}

func TestCancelOrder(t *testing.T) {
	repo := ti.Factory.GetCartRepo()
	var buyerId int64 = 1
	c := repo.GetMyCart(buyerId, cart.KNormal)
	_ = joinItemsToCart(c, 1)
	if c.Kind() == cart.KNormal {
		rc := c.(cart.INormalCart)
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
	orderRepo := ti.Factory.GetOrderRepo()
	mmRepo := ti.Factory.GetMemberRepo()
	manager := orderRepo.Manager()
	m := mmRepo.GetMember(buyerId)
	addressId := m.Profile().GetDefaultAddress().GetDomainId()
	data := order.SubmitOrderData{
		BuyerId:       buyerId,
		Type:          order.TRetail,
		Subject:       "",
		AddressId:     addressId,
		CouponCode:    "",
		BalanceDeduct: true,
		AffiliateCode: "",
		PostedData:    nil,
	}
	o, rd, err := manager.SubmitOrder(data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("订单金额为:%d", rd.TradeAmount)
	o = manager.GetOrderById(o.GetAggregateRootId())
	py := o.GetPaymentOrder()
	err = py.PaymentByWallet("支付订单")
	pv := py.Get()
	payState := pv.State
	if payState == payment.StateFinished {
		t.Logf("订单支付完成,金额：%d", pv.FinalAmount)
	} else {
		t.Logf("订单未完成支付,状态：%d;订单号：%s", pv.State, py.TradeNo())
	}
	t.Logf("支付单信息：%#v", pv)

	no := o.(order.INormalOrder)
	for _, v := range no.GetSubOrders() {
		err = v.Cancel(true, "买多了，不想要了")
		if err != nil {
			t.Log("取消失败：", err.Error())
			t.FailNow()
		}
	}
	if err == nil {
		t.Log("订单取消成功")
	}
}

func TestCancelSubOrderByOrderNo(t *testing.T) {
	var orderId int64 = 117
	orderRepo := ti.Factory.GetOrderRepo()
	manager := orderRepo.Manager()
	is := manager.GetSubOrder(orderId)
	err := is.Cancel(true, "不想要了")
	if err != nil {
		t.Log("取消失败：", err.Error())
		t.FailNow()
	}
}

// 测试提交普通订单,并完成付款
func TestSubmitNormalOrder(t *testing.T) {
	var buyerId int64 = 1
	cartRepo := ti.Factory.GetCartRepo()
	c := cartRepo.GetMyCart(buyerId, cart.KNormal)
	_ = joinItemsToCart(c, 47)
	err := joinItemsToCart(c, 51)
	if err != nil {
		t.Error("购物车加入失败:", err.Error())
		t.FailNow()
	}
	rc := c.(cart.INormalCart)
	if len(rc.Value().Items) == 0 {
		t.Log("购物车是空的")
		t.FailNow()
	}
	t.Log("购物车如下:")
	for _, v := range rc.Items() {
		t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
	}
	_, err = c.Save()
	if err != nil {
		t.Error("保存购物车失败:", err.Error())
		t.Fail()
	}
	orderRepo := ti.Factory.GetOrderRepo()
	manager := orderRepo.Manager()
	buyer := ti.Factory.GetMemberRepo().GetMember(buyerId)
	addressId := buyer.Profile().GetDefaultAddress().GetDomainId()

	data := order.SubmitOrderData{
		BuyerId:       buyer.GetValue().Id,
		Type:          order.TRetail,
		Subject:       "",
		AddressId:     addressId,
		CouponCode:    "",
		BalanceDeduct: true,
		AffiliateCode: "",
		PostedData:    nil,
	}
	o, _, err := manager.SubmitOrder(data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	//ro := o.(order.INormalOrder)
	//_ = ro.OnlinePaymentTradeFinish()
	time.Sleep(time.Second * 2)
	t.Log("提交成功，订单号：", o.OrderNo())
}

// 测试从订单重新创建订单并提交付款
func TestRebuildSubmitNormalOrder(t *testing.T) {
	orderNo := "1230124001810642"
	repo := ti.Factory.GetOrderRepo()
	memRepo := ti.Factory.GetMemberRepo()
	payRepo := ti.Factory.GetPaymentRepo()
	so := repo.GetSubOrderByOrderNo(orderNo)
	io := so.ParentOrder()
	ic := io.BuildCart()
	ic.Save()
	memberId := io.Buyer().GetAggregateRootId()
	addressId := memRepo.GetDeliverAddress(memberId)[0].Id

	data := order.SubmitOrderData{
		BuyerId:       memberId,
		Type:          io.Type(),
		Subject:       "",
		AddressId:     addressId,
		CouponCode:    "",
		BalanceDeduct: true,
		AffiliateCode: "",
		PostedData:    nil,
	}
	nio, _, err := repo.Manager().SubmitOrder(data)
	if err != nil {
		t.Log("提交订单", err.Error())
		t.FailNow()
	}
	t.Logf("提交的订单号为：%s", nio.OrderNo())
	ipo := payRepo.GetPaymentOrderByOrderNo(int(order.TRetail), nio.OrderNo())
	err = ipo.PaymentFinish("alipay", "1233535080808wr")
	if err == nil {
		t.Logf("支付的交易号为：%s,最终金额:%d", nio.OrderNo(), ipo.Get().FinalAmount)
	} else {
		t.Log("支付订单", err.Error())
		t.FailNow()
	}
	// 开始完成发货流程并收货\
	newOrder := repo.Manager().GetOrderById(nio.GetAggregateRootId())
	ino := newOrder.(order.INormalOrder)
	for _, v := range ino.GetSubOrders() {
		v.Confirm()
		err = v.PickUp()
		if err == nil {
			err = v.Ship(1, "12345345")
			// if err == nil {
			// 	err = v.BuyerReceived()
			// }
		}
		if err != nil {
			t.Log("收货不成功：", err)
			t.FailNow()
		}
	}
	time.Sleep(2000)
}

func TestFinishSubOrder(t *testing.T) {
	repo := ti.Factory.GetOrderRepo()
	io := repo.GetSubOrderByOrderNo("1230115001702364")
	err := io.BuyerReceived()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	for {
	}
	time.Sleep(5000)
}

func TestFinishNormalOrder(t *testing.T) {
	repo := ti.Factory.GetOrderRepo()
	io := repo.Manager().GetOrderByNo("1230115001702364")
	subOrders := io.(order.INormalOrder).GetSubOrders()
	for _, o := range subOrders {
		//err := o.Confirm()
		//err := o.Ship(1,"33")
		err := o.BuyerReceived()
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}
	time.Sleep(5000)
}

// 测试批发订单,并完成付款
func TestWholesaleOrder(t *testing.T) {
	var buyerId int64 = 1
	cartRepo := ti.Factory.GetCartRepo()
	c := cartRepo.GetMyCart(buyerId, cart.KWholesale)
	_ = joinItemsToCart(c, 1)
	rc := c.(cart.IWholesaleCart)
	if len(rc.GetValue().Items) == 0 {
		t.Log("购物车是空的")
		t.FailNow()
	}
	t.Log("购物车如下:")
	for _, v := range rc.Items() {
		t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
	}
	_, err := c.Save()
	if err != nil {
		t.Error("保存购物车失败:", err.Error())
		t.Fail()
	}

	orderRepo := ti.Factory.GetOrderRepo()
	manager := orderRepo.Manager()

	buyer := ti.Factory.GetMemberRepo().GetMember(buyerId)
	addressId := buyer.Profile().GetDefaultAddress().GetDomainId()

	data := map[string]string{
		"address_id":       strconv.Itoa(int(addressId)),
		"seller_comment_1": "测试留言",
		"checked":          GetCartCheckedData(c),
	}

	log.Println("----", fmt.Sprintf("%#v", data))
	data1 := order.SubmitOrderData{
		BuyerId:       buyerId,
		Type:          order.TWholesale,
		Subject:       "",
		AddressId:     addressId,
		CouponCode:    "",
		BalanceDeduct: false,
		AffiliateCode: "",
		PostedData:    parser.NewPostedData(data, nil),
	}
	_, rd, err := manager.SubmitOrder(data1)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	arr := strings.Split(rd.OrderNo, ",")
	t.Logf("批发单拆分数量：%d , 订单号：%s", len(arr), rd.OrderNo)

	for _, orderNo := range arr {
		if orderNo != "" {
			// 重新获取订单
			o := manager.GetOrderByNo(orderNo)
			io := o.(order.IWholesaleOrder)
			// 付款操作
			io.OnlinePaymentTradeFinish()
			time.Sleep(time.Second * 5)
			// 可能会自动完成
			//logState(t, io.Confirm(), o)
			logState(t, io.PickUp(), o)
			logState(t, io.Ship(1, "123456"), o)
			//logState(t, io.BuyerReceived(), o)
		}
	}
}

func TestTradeOrder(t *testing.T) {
	var rate = 0.8 // 结算给商家80%
	var storeId = 1
	repo := ti.Factory.GetOrderRepo()
	manager := repo.Manager()
	cashPay := true
	requireTicket := true
	if requireTicket {
		//repos.DefaultGlobMchSaleConf.TradeOrderRequireTicket = true
	}
	c := order.SubmitOrderData{
		BuyerId:       397,
		Type:          0,
		Subject:       "万宁佛山祖庙店",
		AddressId:     0,
		CouponCode:    "",
		BalanceDeduct: false,
		AffiliateCode: "",
		PostedData: parser.NewPostedData(nil, &proto.SubmitOrderRequest{
			TradeOrder: &proto.TradeOrderRequest{
				StoreId:     int64(storeId),
				TradeAmount: 100,
				Discount:    float32(rate),
			},
		}),
	}
	o, _, err := manager.SubmitOrder(c)
	if err != nil {
		t.Errorf("提交订单错误：%s", err.Error())
		t.FailNow()
	}
	io := o.(order.ITradeOrder)
	// 使用现金支付或者使用钱包支付
	if cashPay {
		err = io.CashPay()
		if err != nil {
			t.Errorf("现金支付错误：%s", err.Error())
			t.FailNow()
		}
	} else {
		py := o.GetPaymentOrder()
		err = py.PaymentByWallet("订单支付")
		if err != nil {
			t.Errorf("钱包支付错误：%s", err.Error())
			t.FailNow()
		}
	}
	time.Sleep(time.Second * 2)
	o = manager.GetOrderById(o.GetAggregateRootId())
	t.Log("订单状态为：", o.State().String())
	if requireTicket {
		t.Log("上传发票")
		io = o.(order.ITradeOrder)
		err := io.UpdateTicket("//img.ts.com/res/nopic.gif")
		if err != nil {
			t.Errorf("上传发票出错：%s", err.Error())
			t.FailNow()
		}
		t.Log("订单状态为：", o.State().String())
	}
}

func TestMergePaymentOrder(t *testing.T) {
	repo := ti.Factory.GetOrderRepo()
	memRepo := ti.Factory.GetMemberRepo()
	io := repo.Manager().GetOrderByNo("1180517000262166")
	ic := io.BuildCart()
	ic.Save()
	memberId := io.Buyer().GetAggregateRootId()
	addressId := memRepo.GetDeliverAddress(memberId)[0].Id
	data := order.SubmitOrderData{
		BuyerId:       io.Buyer().GetValue().Id,
		Type:          io.Type(),
		Subject:       "",
		AddressId:     addressId,
		CouponCode:    "",
		BalanceDeduct: false,
		AffiliateCode: "",
		PostedData:    nil,
	}
	_, rd, err := repo.Manager().SubmitOrder(data)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	print(fmt.Sprintf("%#v", rd))
}

// 通知交易单
func TestNotifyTradeOrder(t *testing.T) {
	orderNo := "1180518115439092"
	sub := true
	rds := ti.GetApp().Storage().(storage.IRedisStorage)
	conn := rds.GetConn()
	defer conn.Close()
	value := orderNo
	if sub {
		value = "sub!" + orderNo
	}
	conn.Do("RPUSH", variable.KvOrderBusinessQueue, value)
}
