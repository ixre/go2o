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

func getRepo() order.IOrderRepo {
	app := include.GetApp()
	db := app.Db()
	sto := app.Storage()
	goodsRepo := repository.NewGoodsItemRepo(db, productRepo, expressRepo, valRepo)
	valRepo := repository.NewValueRepo(db, sto)
	userRepo := repository.NewUserRepo(db)
	notifyRepo := repository.NewNotifyRepo(db)
	mssRepo := repository.NewMssRepo(db, notifyRepo, valRepo)
	expressRepo := repository.NewExpressRepo(db, valRepo)
	shipRepo := repository.NewShipmentRepo(db, expressRepo)
	memberRepo := repository.NewMemberRepo(app.Storage(), db, mssRepo, valRepo)
	itemRepo := repository.NewProductRepo(db)
	tagSaleRepo := repository.NewTagSaleRepo(db)
	promRepo := repository.NewPromotionRepo(db, goodsRepo, memberRepo)
	cateRepo := repository.NewCategoryRepo(db, valRepo, sto)
	saleRepo := repository.NewSaleRepo(db, cateRepo, valRepo, tagSaleRepo,
		itemRepo, expressRepo, goodsRepo, promRepo)
	cartRepo := repository.NewCartRepo(db, memberRepo, goodsRepo)
	shopRepo := repository.NewShopRepo(db, sto)
	mchRepo := repository.NewMerchantRepo(db, sto, shopRepo, userRepo,
		memberRepo, mssRepo, valRepo)
	//personFinanceRepo := repository.NewPersonFinanceRepository(db, memberRepo)
	deliveryRepo := repository.NewDeliverRepo(db)
	//contentRepo := repository.NewContentRepo(db)
	//adRepo := repository.NewAdvertisementRepo(db)
	return repository.NewOrderRepo(app.Storage(), db, mchRepo, nil, saleRepo, cartRepo, goodsRepo,
		promRepo, memberRepo, deliveryRepo, expressRepo, shipRepo, valRepo)
}

func TestOrderSetup(t *testing.T) {
	orderNo := "100000735578"
	orderRepo := getRepo()
	v := orderRepo.GetSubOrderByNo(orderNo)
	o := orderRepo.Manager().GetSubOrder(v.Id)

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
