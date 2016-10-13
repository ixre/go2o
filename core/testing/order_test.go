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
	"go2o/core/domain/interface/order"
	"go2o/core/repository"
	"go2o/core/testing/include"
	"testing"
)

func getRep() order.IOrderRep {
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
	mchRep := repository.NewMerchantRep(db, sto, shopRep, userRep,
		memberRep, mssRep, valRep)
	//personFinanceRep := repository.NewPersonFinanceRepository(db, memberRep)
	deliveryRep := repository.NewDeliverRep(db)
	//contentRep := repository.NewContentRep(db)
	//adRep := repository.NewAdvertisementRep(db)
	return repository.NewOrderRep(app.Storage(), db, mchRep, nil, saleRep, cartRep, goodsRep,
		promRep, memberRep, deliveryRep, expressRep, shipRep, valRep)
}

func TestOrderSetup(t *testing.T) {
	orderNo := "100000735578"
	orderRep := getRep()
	v := orderRep.GetSubOrderByNo(orderNo)
	o := orderRep.Manager().GetSubOrder(v.Id)

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
