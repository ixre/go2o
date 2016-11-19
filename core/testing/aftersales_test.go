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
	"github.com/jsix/gof/log"
	"go2o/core/domain/interface/after-sales"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/order"
	"go2o/core/repository"
	"go2o/core/testing/include"
	"testing"
)

func getRep22() order.IOrderRep {
	app := include.GetApp()
	db := app.Db()
	sto := app.Storage()
	goodsRep := repository.NewGoodsRep(db)
	valRep := repository.NewValueRep(db, sto)
	userRep := repository.NewUserRep(db)
	notifyRep := repository.NewNotifyRep(db)
	mssRep := repository.NewMssRep(db, notifyRep, valRep)
	expressRep := repository.NewExpressRep(db, valRep)
	shipRep := repository.NewShipmentRep(db, expressRep)
	memberRep := repository.NewMemberRep(app.Storage(), db, mssRep, valRep)
	itemRep := repository.NewItemRep(db)
	tagSaleRep := repository.NewTagSaleRep(db)
	promRep := repository.NewPromotionRep(db, goodsRep, memberRep)
	cateRep := repository.NewCategoryRep(db, valRep, sto)
	saleRep := repository.NewSaleRep(db, cateRep, valRep, tagSaleRep,
		itemRep, expressRep, goodsRep, promRep)
	cartRep := repository.NewCartRep(db, memberRep, goodsRep)
	shopRep := repository.NewShopRep(db, sto)
	mchRep := repository.NewMerchantRep(db, sto, shopRep, userRep, memberRep, mssRep, valRep)
	//personFinanceRep := repository.NewPersonFinanceRepository(db, memberRep)
	deliveryRep := repository.NewDeliverRep(db)
	//contentRep := repository.NewContentRep(db)
	//adRep := repository.NewAdvertisementRep(db)
	return repository.NewOrderRep(app.Storage(), db, mchRep, nil, saleRep, cartRep, goodsRep,
		promRep, memberRep, deliveryRep, expressRep, shipRep, valRep)
}

func getMemberRep() member.IMemberRep {
	app := include.GetApp()
	db := app.Db()
	sto := app.Storage()
	valRep := repository.NewValueRep(db, sto)
	notifyRep := repository.NewNotifyRep(db)
	mssRep := repository.NewMssRep(db, notifyRep, valRep)
	return repository.NewMemberRep(app.Storage(), db, mssRep, valRep)
}

func getAfterSalesRep() afterSales.IAfterSalesRep {
	db := include.GetApp().Db()
	sto := include.GetApp().Storage()
	memberRep := getMemberRep()
	orderRep := getRep22()
	valRep := repository.NewValueRep(db, sto)
	paymentRep := repository.NewPaymentRep(sto, db, memberRep, orderRep, valRep)
	return repository.NewAfterSalesRep(db, getRep22(), getMemberRep(), paymentRep)
}

// 测试退款
func TestOrderRefund(t *testing.T) {
	subOrderNo := "100000160304"
	orderRep := getRep22()
	rep := getAfterSalesRep()
	v := orderRep.GetSubOrderByNo(subOrderNo)
	od := orderRep.Manager().GetSubOrder(v.Id)
	ro := rep.CreateAfterSalesOrder(&afterSales.AfterSalesOrder{
		Id: 0,
		// 订单编号
		OrderId: v.Id,
		// 类型，退货、换货、维修
		Type: afterSales.TypeRefund,
		// 售后原因
		Reason: "不想要了,我想推掉",
	})
	//err := ro.Submit()
	//if err != nil{
	//	t.Log("提交退货单(未设定产品):",err.Error())
	//}
	item := od.Items()[0]
	err := ro.SetItem(item.SnapshotId, item.Quantity+1)
	if err != nil {
		t.Log("设定退货产品(超出数量):", err.Error())
	}
	err = ro.SetItem(item.SnapshotId, item.Quantity)
	if err != nil {
		t.Log("设定退货产品(正常数量):", err.Error())
	}
	_, err = ro.Submit()
	if err != nil {
		t.Log("提交售后单", err.Error())
	}

	err = ro.Agree()
	if err != nil {
		t.Log("运营商同意:", err.Error())
	}

	err = ro.Confirm()
	if err != nil {
		t.Log("系统确认:", err.Error())
	}

	err = ro.Reject("系统退回")
	if err != nil {
		t.Log("系统退回:", err.Error())
	}

	err = ro.Process()
	if err != nil {
		t.Log("系统处理:", err.Error())
	}

	log.Println("售后单状态为:", ro.Value().State, ro.Value().State == afterSales.StatCompleted)
	log.Printf("%#v", ro.Value().Data)
}

// 测试退货
func TestOrderReturn(t *testing.T) {
	subOrderNo := "100000160304"
	orderRep := getRep22()
	rep := getAfterSalesRep()
	v := orderRep.GetSubOrderByNo(subOrderNo)
	od := orderRep.Manager().GetSubOrder(v.Id)
	ro := rep.CreateAfterSalesOrder(&afterSales.AfterSalesOrder{
		Id: 0,
		// 订单编号
		OrderId: v.Id,
		// 类型，退货、换货、维修
		Type: afterSales.TypeReturn,
		// 售后原因
		Reason: "不想要了,我想推掉",
	})
	//err := ro.Submit()
	//if err != nil{
	//	t.Log("提交退货单(未设定产品):",err.Error())
	//}
	item := od.Items()[0]
	err := ro.SetItem(item.SnapshotId, item.Quantity+1)
	if err != nil {
		t.Log("设定退货产品(超出数量):", err.Error())
	}
	err = ro.SetItem(item.SnapshotId, item.Quantity)
	if err != nil {
		t.Log("设定退货产品(正常数量):", err.Error())
	}
	_, err = ro.Submit()
	if err != nil {
		t.Log("提交售后单", err.Error())
	}
	err = ro.Agree()
	if err != nil {
		t.Log("运营商同意:", err.Error())
	}
	err = ro.Confirm()
	if err != nil {
		t.Log("确认退货:", err.Error())
	}
	err = ro.ReturnShip("顺风快递", "1000", "")
	if err != nil {
		t.Log("快递货物:", err.Error())
	}
	err = ro.ReturnReceive()
	if err != nil {
		t.Log("接收退货:", err.Error())
	}
	err = ro.Process()
	if err != nil {
		t.Log("处理退货:", err.Error())
	}
	log.Println("售后单状态为:", ro.Value().State, ro.Value().State == afterSales.StatCompleted)
	log.Printf("%#v", ro.Value().Data)
}

// 测试换货
func TestOrderExchange(t *testing.T) {
	subOrderNo := "100000160304"
	orderRep := getRep22()
	rep := getAfterSalesRep()
	v := orderRep.GetSubOrderByNo(subOrderNo)
	od := orderRep.Manager().GetSubOrder(v.Id)
	ro := rep.CreateAfterSalesOrder(&afterSales.AfterSalesOrder{
		// 订单编号
		OrderId: v.Id,
		// 类型，退货、换货、维修
		Type: afterSales.TypeExchange,
		// 售后原因
		Reason: "产品中间有瑕疵,请帮我换货!",
	})
	//err := ro.Submit()
	//if err != nil{
	//	t.Log("提交退货单(未设定产品):",err.Error())
	//}
	item := od.Items()[0]
	err := ro.SetItem(item.SnapshotId, item.Quantity+1)
	if err != nil {
		t.Log("设定换货产品(超出数量):", err.Error())
	}
	err = ro.SetItem(item.SnapshotId, item.Quantity)
	if err != nil {
		t.Log("设定换货产品(正常数量):", err.Error())
	}
	_, err = ro.Submit()
	if err != nil {
		t.Log("提交售后单", err.Error())
	}
	err = ro.Agree()
	if err != nil {
		t.Log("运营商同意:", err.Error())
	}

	err = ro.Confirm()
	if err != nil {
		t.Log("确认退货:", err.Error())
	}

	err = ro.ReturnShip("顺风快递", "10007927972432", "")
	if err != nil {
		t.Log("快递货物:", err.Error())
	}
	err = ro.ReturnReceive()
	if err != nil {
		t.Log("接收退货:", err.Error())
	}

	eo := ro.(afterSales.IExchangeOrder)
	err = eo.ExchangeShip("申通快递", "10989789274234")
	if err != nil {
		t.Log("配送货物:", err.Error())
	}

	err = eo.LongReceive()
	if err != nil {
		t.Log("延长配送时间:", err.Error())
	}

	err = eo.ExchangeReceive()
	if err != nil {
		t.Log("收到换货:", err.Error())
	}

	log.Println("售后单状态为:", ro.Value().State, ro.Value().State == afterSales.StatCompleted)
	log.Printf("%#v", ro.Value().Data)
}
