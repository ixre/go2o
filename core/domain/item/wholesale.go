package item

import (
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/item"
)

var _ item.IWholesaleItem = new(wholesaleItemImpl)

type wholesaleItemImpl struct {
	itemId int32
	value  *item.WsItem
	it     item.IGoodsItem
	repo   item.IItemWholesaleRepo
}

func newWholesaleItem(itemId int32, it item.IGoodsItem,
	repo item.IItemWholesaleRepo) item.IWholesaleItem {
	return (&wholesaleItemImpl{
		itemId: itemId,
		it:     it,
		repo:   repo,
	}).init()
}

func (w *wholesaleItemImpl) init() item.IWholesaleItem {
	v := w.repo.GetWsItem(w.itemId)
	if v == nil {
		v = &item.WsItem{
			ItemId:          w.itemId,
			EnableWholesale: 0,
		}
		w.repo.SaveWsItem(v, true)
	}
	w.value = v
	return w
}

// 获取领域编号
func (w *wholesaleItemImpl) GetDomainId() int32 {
	return w.itemId
}

// 是否允许批发
func (w *wholesaleItemImpl) Wholesale() bool {
	return w.value.EnableWholesale == 1
}

// 开启批发功能
func (w *wholesaleItemImpl) TurnWholesale(on bool) error {
	var iv int32 = util.BoolExt.TInt32(on, 1, 0)
	if w.value.EnableWholesale != iv {
		w.value.EnableWholesale = iv
		_, err := w.Save()
		return err
	}
	return nil
}

// 保存
func (w *wholesaleItemImpl) Save() (int32, error) {
	return util.I32Err(w.repo.SaveWsItem(w.value, false))
}

// 根据商品金额获取折扣
func (w *wholesaleItemImpl) GetWholesaleDiscount(groupId int32, amount int32) float64 {
	var rate float64 = 0
	arr := w.GetItemDiscount(groupId)
	if len(arr) > 0 {
		var maxRequire int32
		for _, v := range arr {
			if v.RequireAmount > maxRequire && amount >= v.RequireAmount {
				maxRequire = v.RequireAmount
				rate = v.DiscountRate
			}
		}
	}
	return rate
}

// 获取全部批发折扣
func (w *wholesaleItemImpl) GetItemDiscount(groupId int32) []*item.WsItemDiscount {
	return w.repo.SelectWsItemDiscount("item_id=? AND buyer_gid=?",
		w.value.ItemId, groupId)
}

// 保存批发折扣
func (w *wholesaleItemImpl) SaveItemDiscount(groupId int32, arr []*item.WsItemDiscount) error {
	// 获取存在的项
	old := w.GetItemDiscount(groupId)
	// 分析当前数据并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*item.WsItemDiscount, len(arr))
	for _, v := range arr {
		currMap[v.RequireAmount] = v
	}
	// 筛选出要删除的项,如存在，则赋予ID
	for _, v := range old {
		new := currMap[v.RequireAmount]
		if new == nil {
			delList = append(delList, v.ID)
		} else {
			new.ID = v.ID
		}
	}
	// 删除项
	for _, v := range delList {
		w.repo.BatchDeleteWsItemDiscount("id=?", v)
	}
	// 保存项
	for _, v := range arr {
		v.ItemId = w.itemId
		v.BuyerGid = groupId
		i, err := util.I32Err(w.repo.SaveWsItemDiscount(v))
		if err == nil {
			v.ID = i
		}
	}
	return nil
}

// 获取批发价格
func (w *wholesaleItemImpl) GetWholesalePrice(skuId, quantity int32) float64 {
	var price float64 = 0
	arr := w.GetSkuPrice(skuId)
	if len(arr) > 0 {
		var maxRequire int32
		for _, v := range arr {
			if v.RequireQuantity > maxRequire && quantity >= v.RequireQuantity {
				maxRequire = v.RequireQuantity
				price = v.WholesalePrice
			}
		}
	}
	return price
}

// 根据SKU获取价格设置
func (w *wholesaleItemImpl) GetSkuPrice(skuId int32) []*item.WsSkuPrice {
	return w.repo.SelectWsSkuPrice("item_id=? AND sku_id=?",
		w.value.ItemId, skuId)
}

// 保存批发SKU价格设置
func (w *wholesaleItemImpl) SaveSkuPrice(skuId int32, arr []*item.WsSkuPrice) error {
	// 获取存在的项
	old := w.GetSkuPrice(skuId)
	// 分析当前数据并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*item.WsSkuPrice, len(arr))
	for _, v := range arr {
		currMap[v.RequireQuantity] = v
	}
	// 筛选出要删除的项,如存在，则赋予ID
	for _, v := range old {
		new := currMap[v.RequireQuantity]
		if new == nil {
			delList = append(delList, v.ID)
		} else {
			new.ID = v.ID
		}
	}
	// 删除项
	for _, v := range delList {
		w.repo.BatchDeleteWsSkuPrice("id=?", v)
	}
	// 保存项
	for _, v := range arr {
		v.ItemId = w.itemId
		v.SkuId = skuId
		i, err := util.I32Err(w.repo.SaveWsSkuPrice(v))
		if err == nil {
			v.ID = i
		}
	}
	return nil
}
