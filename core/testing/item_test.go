package testing

import (
	"go2o/core/domain/interface/item"
	"go2o/core/testing/ti"
	"testing"
)

// 测试批发折扣
func TestItemWholesaleDiscount(t *testing.T) {
	repo := ti.ItemRepo
	mmRepo := ti.MemberRepo
	var itemId int32 = 6      //商品编号
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
	repo := ti.ItemRepo
	var itemId int32 = 6 //商品编号
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
