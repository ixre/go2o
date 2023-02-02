package domain

import (
	"fmt"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/tests/ti"
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
	var itemId int64 = 6     //商品编号
	var disRate = 0.9        //折扣率
	var disAmount int32 = 50 //折扣金额下限
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
				WholesalePrice:  sku.Price - 10,
			},
			{
				RequireQuantity: 10,
				WholesalePrice:  sku.Price - 20,
			},
		}
		err := wsIt.SaveSkuPrice(sku.Id, arr)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}

	sku0 := it.SkuArray()[0]
	skuId := sku0.Id
	skuPrice := sku0.Price

	//计算商品订单的折扣率
	price1 := wsIt.GetWholesalePrice(skuId, 1)
	if price1 != skuPrice {
		t.Error("购买1件,价格应为原价")
		t.Fail()
	}
	price2 := wsIt.GetWholesalePrice(skuId, 2)
	if price2 != skuPrice-10 {
		t.Error("购买1件,价格不正确")
		t.Fail()
	}
	price3 := wsIt.GetWholesalePrice(skuId, 11)
	if price3 != skuPrice-20 {
		t.Error("购买1件,价格不正确")
		t.Fail()
	}
}

// 测试修改非SKU商品价格
func TestUpdateItemNoSkuPrice(t *testing.T) {
	var itemId int64 = 50
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	if it == nil {
		t.Error(item.ErrNoSuchItem)
		t.Failed()
	}
	v := it.GetValue()
	t.Log(v.Price)
	v.Price += 100
	err := it.SetValue(v)
	if err == nil {
		_, err = it.Save()
	}
	if err != nil {
		t.Error(err)
		t.Failed()
	}
}

// 测试保存SKU
func TestSaveItemSku(t *testing.T) {
	var itemId int64 = 50
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	err := it.SetSku(it.SkuArray())
	if err == nil {
		_, err = it.Save()
	}
	if err != nil {
		t.Error(err)
		t.Failed()
	}
}

// 测试保存商品图片
func TestSaveItemImages(t *testing.T) {
	var itemId int64 = 50
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	println(fmt.Sprintf("%#v", it.Images()))
	images := make([]string, 0)
	images = append(images, "https://img14.360buyimg.com/ceco/s300x300_jfs/t1/159722/38/5682/268261/601a43e3E78cbacb6/60bdf8c1c170c8ae.jpg!q70.jpg.webp#1")
	//images = append(images, "https://img14.360buyimg.com/ceco/s300x300_jfs/t1/159722/38/5682/268261/601a43e3E78cbacb6/60bdf8c1c170c8ae.jpg!q70.jpg.webp#2")
	//images = append(images,"https://img14.360buyimg.com/ceco/s300x300_jfs/t1/159722/38/5682/268261/601a43e3E78cbacb6/60bdf8c1c170c8ae.jpg!q70.jpg.webp#3")
	//images = append(images, "https://img14.360buyimg.com/ceco/s300x300_jfs/t1/159722/38/5682/268261/601a43e3E78cbacb6/60bdf8c1c170c8ae.jpg!q70.jpg.webp#4")
	//images = append(images, "https://img14.360buyimg.com/ceco/s300x300_jfs/t1/159722/38/5682/268261/601a43e3E78cbacb6/60bdf8c1c170c8ae.jpg!q70.jpg.webp#5")
	err := it.SetImages(images)
	if err == nil {
		_, err = it.Save()
	}
	if err != nil {
		t.Error(err)
		t.Failed()
	}
}

func TestSaveAffiliateItemFlag(t *testing.T) {
	var itemId int64 = 47
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	err := it.GrantFlag(item.FlagAffiliate)
	if err == nil {
		_, err = it.Save()
	}
	if err != nil {
		t.Error(err)
		t.Failed()
	}
}
func TestSaveItemFreeDeliveryFlag(t *testing.T) {
	var itemId int64 = 1
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	iv := it.GetValue()
	iv.ExpressTid = 0
	err := it.SetValue(iv)
	if err == nil {
		_, err = it.Save()
	}
	if err != nil {
		t.Error(err)
	}
}

func TestCheckContainerItemFlag(t *testing.T) {
	t.Log(-1&item.FlagNewOnShelve == item.FlagNewOnShelve)
}

func TestAuditItem(t *testing.T) {
	var itemId int64 = 1
	repo := ti.Factory.GetItemRepo()
	it := repo.GetItem(itemId)
	it.Review(true, "")
}
