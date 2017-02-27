package testing

import (
	"testing"
	//"go2o/core/testing/ti"
	"go2o/core/domain/interface/cart"
	"go2o/core/testing/ti"
)

func TestRetailCart(t *testing.T) {
	repo := ti.CartRepo
	itemRepo := ti.ItemRepo
	c := repo.GetMemberCurrentCart(1)
	if c == nil {
		c = repo.CreateCart(&cart.ValueCart{
			BuyerId: 1,
		})
		c.Save()
	}
	gs := itemRepo.GetItem(3)
	arr := gs.SkuArray()
	itemId := gs.GetAggregateRootId()
	skuId := arr[0].Id
	err := c.Put(itemId, skuId, 1)
	if err != nil {
		t.Error("购物车加入失败:", err.Error())
		t.Fail()
	}
	if c.Kind() == cart.KRetail {
		rc := c.(cart.IRetailCart)
		t.Log("购物车有：", len(rc.Items()), "件商品")
	}
}

func TestWholesaleCart(t *testing.T) {
	repo := ti.CartRepo
	itemRepo := ti.ItemRepo
	c := repo.GetMyCart(1, cart.KWholesale)
	gs := itemRepo.GetItem(3)
	arr := gs.SkuArray()
	itemId := gs.GetAggregateRootId()
	skuId := arr[0].Id
	err := c.Put(itemId, skuId, 1)
	if err != nil {
		t.Error("购物车加入失败:", err.Error())
		t.Fail()
	}
	if c.Kind() == cart.KWholesale {
		rc := c.(cart.IWholesaleCart)
		t.Log("购物车有：", len(rc.Items()), "件商品")
	}
	c.Save()
}
