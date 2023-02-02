package domain

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/cart"
	"github.com/ixre/go2o/tests/ti"
	"github.com/ixre/gof/types/typeconv"
)

// 测试普通购物车
func TestNormalCart(t *testing.T) {
	repo := ti.Factory.GetCartRepo()
	c := repo.GetMyCart(1, cart.KNormal)
	_ = joinItemsToCart(c, 1)
	//_ = joinItemsToCart(c, 51)
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
		t.Fail()
	}
}

func TestCheckOnlyItem(t *testing.T) {
	repo := ti.Factory.GetCartRepo()
	ic := repo.NewTempNormalCart(0,"1234")
	err := ic.Put(1,7,1,true)
	if err != nil{
		t.Error(err)
		t.FailNow()
	}
 	in :=	ic.(cart.INormalCart)
	t.Log(typeconv.MustJson(in.Items()))
}

// 测试合并购物车
func TestCombineCart(t *testing.T) {
	repo := ti.Factory.GetCartRepo()
	c := repo.GetMyCart(1, cart.KNormal)
	//c2 := repo.NewNormalCart()

	if c.Kind() == cart.KNormal {
		rc := c.(cart.INormalCart)
		t.Log("购物车如下:")
		for _, v := range rc.Items() {
			t.Logf("商品：%d-%d 数量：%d\n", v.ItemId, v.SkuId, v.Quantity)
		}
	}
}

func joinItemsToCart(c cart.ICart, itemId int64) error {
	itemRepo := ti.Factory.GetItemRepo()
	gs := itemRepo.GetItem(itemId)
	arr := gs.SkuArray()
	var skuId int64 = 0
	if len(arr) > 0 {
		skuId = arr[0].Id
	}
	return c.Put(itemId, skuId, 1, false)
}

// 生成购物车全部结算的数据
func GetCartCheckedData(c cart.ICart) string {
	mp := make(map[string][]string)
	if c.Kind() == cart.KWholesale {
		wc := c.(cart.IWholesaleCart)
		for _, v := range wc.Items() {
			id := strconv.Itoa(int(v.ItemId))
			if _, ok := mp[id]; !ok {
				mp[id] = []string{}
			}
			mp[id] = append(mp[id], strconv.Itoa(int(v.SkuId)))
		}
	} else {
		rc := c.(cart.INormalCart)
		for _, v := range rc.Items() {
			id := strconv.Itoa(int(v.ItemId))
			if _, ok := mp[id]; !ok {
				mp[id] = []string{}
			}
			mp[id] = append(mp[id], strconv.Itoa(int(v.SkuId)))
		}
	}
	b, err := json.Marshal(mp)
	if err != nil {
		log.Println("--- parse cart checked data error :", err)
	}
	return string(b)
}

// 测试批发购物车
func TestWholesaleCart(t *testing.T) {
	repo := ti.Factory.GetCartRepo()
	c := repo.GetMyCart(1, cart.KWholesale)
	_ = joinItemsToCart(c, 1)
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
