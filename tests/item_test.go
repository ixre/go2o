package tests

import (
	"go2o/core/domain/interface/item"
	"go2o/tests/ti"
	"testing"
)

// 检查商品库存
func TestItemSkuStock(t *testing.T) {
	var itemId int64 = 66
	var skuId int64 = 454
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	if it == nil {
		t.Error(item.ErrNoSuchItem)
	}
	sku := it.GetSku(skuId)
	if sku == nil {
		t.Error(item.ErrNoSuchSku)
	}
	t.Log("ItemId:", itemId, " SkuId:", skuId,
		" stock is ", sku.Stock)
}

// 测试同步批发商品
func TestSyncWholesaleItem(t *testing.T) {
	venRepo := ti.Factory.GetMerchantRepo()
	vd := venRepo.GetMerchant(1)
	mp := vd.Wholesaler().SyncItems(true)
	t.Logf("sync finished, add:%d,del:%d", mp["add"], mp["del"])
}

// 测试批发折扣
func TestItemWholesaleDiscount(t *testing.T) {
	repo := ti.Factory.GetItemRepo()
	mmRepo := ti.Factory.GetMemberRepo()
	var itemId int64 = 6      //商品编号
	var disRate float64 = 0.9 //折扣率
	var disAmount int32 = 50  //折扣金额下限
	it := repo.GetItem(itemId)
	wsIt := it.Wholesale()
	groups := mmRepo.GetManager().GetAllBuyerGroups()
	for _, g := range groups {
		arr := []*item.WsItemDiscount{
			{
				RequireAmount: disAmount,
				DiscountRate:  disRate,
			},
		}
		err := wsIt.SaveItemDiscount(g.ID, arr)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}
	groupId := groups[0].ID
	//计算商品订单的折扣率
	rate1 := wsIt.GetWholesaleDiscount(groupId, 120)
	if rate1 != disRate {
		t.Error("大于目标金额,折扣率不正确")
		t.Fail()
	}
	rate2 := wsIt.GetWholesaleDiscount(groupId, 0)
	if rate2 != 0 {
		t.Error("低于目标金额，折扣率不为零")
		t.Fail()
	}
}

// 测试批发SKU价格
func TestItemWholesaleSkuPrice(t *testing.T) {
	repo := ti.Factory.GetItemRepo()
	var itemId int64 = 6 //商品编号
	it := repo.GetItem(itemId)
	wsIt := it.Wholesale()
	// 保存SKU价格
	for _, sku := range it.SkuArray() {
		arr := []*item.WsSkuPrice{
			{
				RequireQuantity: 2,
				WholesalePrice:  float64(sku.Price - 0.1),
			},
			{
				RequireQuantity: 10,
				WholesalePrice:  float64(sku.Price - 0.2),
			},
		}
		err := wsIt.SaveSkuPrice(sku.ID, arr)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	sku0 := it.SkuArray()[0]
	skuId := sku0.ID
	skuPrice := float64(sku0.Price)

	//计算商品订单的折扣率
	price1 := wsIt.GetWholesalePrice(skuId, 1)
	if price1 != skuPrice {
		t.Error("购买1件,价格应为原价")
		t.Fail()
	}
	price2 := wsIt.GetWholesalePrice(skuId, 2)
	if price2 != skuPrice-0.1 {
		t.Error("购买1件,价格不正确")
		t.Fail()
	}
	price3 := wsIt.GetWholesalePrice(skuId, 11)
	if price3 != skuPrice-0.2 {
		t.Error("购买1件,价格不正确")
		t.Fail()
	}
}
