package testing

import (
	"go2o/core/domain/interface/cart"
	"go2o/core/testing/ti"
	"testing"
)

// 测试零售购物车
func TestRetailCart(t *testing.T) {
	repo := ti.CartRepo
	c := repo.GetMyCart(1, cart.KRetail)
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
		t.Fail()
	}
}

// 测试合并购物车
func TestCombineCart(t *testing.T) {
	repo := ti.CartRepo
	c := repo.GetMyCart(1, cart.KRetail)
	//c2 := repo.NewRetailCart()

	if c.Kind() == cart.KRetail {
		rc := c.(cart.IRetailCart)
		t.Log("购物车如下:")
		for _, v := range rc.Items() {
			t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
		}
	}
}

func joinItemsToCart(c cart.ICart, t *testing.T) {
	itemRepo := ti.ItemRepo
	gs := itemRepo.GetItem(3)
	arr := gs.SkuArray()
	itemId := gs.GetAggregateRootId()
	skuId := arr[0].Id
	err := c.Put(itemId, skuId, 1)
	if err != nil {
		t.Error("购物车加入失败:", err.Error())
		t.Fail()
	}
}

// 测试批发购物车
func TestWholesaleCart(t *testing.T) {
	repo := ti.CartRepo
	c := repo.GetMyCart(1, cart.KWholesale)
	joinItemsToCart(c, t)
	if c.Kind() == cart.KWholesale {
		rc := c.(cart.IWholesaleCart)
		t.Log("购物车如下:")
		for _, v := range rc.Items() {
			t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
		}
	}
	_, err := c.Save()
	if err != nil {
		t.Error("保存购物车失败:", err.Error())
		t.Fail()
	}
}
