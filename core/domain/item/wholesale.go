package item

import (
	"encoding/json"
	"github.com/ixre/gof/math"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/infrastructure/format"
	"strconv"
	"strings"
)

var _ item.IWholesaleItem = new(wholesaleItemImpl)

type wholesaleItemImpl struct {
	itemId   int64
	value    *item.WsItem
	it       item.IGoodsItem
	itemRepo item.IGoodsItemRepo
	repo     item.IItemWholesaleRepo
}

func newWholesaleItem(itemId int64, it item.IGoodsItem,
	itemRepo item.IGoodsItemRepo, repo item.IItemWholesaleRepo) item.IWholesaleItem {
	return (&wholesaleItemImpl{
		itemId:   itemId,
		it:       it,
		itemRepo: itemRepo,
		repo:     repo,
	}).init()
}

func (w *wholesaleItemImpl) init() item.IWholesaleItem {
	v := w.repo.GetWsItem(w.itemId)
	if v == nil {
		iv := w.it.GetValue()
		v = &item.WsItem{
			ItemId:      w.itemId,
			VendorId:    iv.VendorId,
			Price:       float64(iv.Price),
			PriceRange:  iv.PriceRange,
			ShelveState: item.ShelvesInWarehouse,
			ReviewState: iv.ReviewState,
		}
		w.repo.SaveWsItem(v, true)
	}
	w.value = v
	return w
}

// 获取领域编号
func (w *wholesaleItemImpl) GetDomainId() int64 {
	return w.itemId
}

// 是否允许批发
func (w *wholesaleItemImpl) CanWholesale() bool {
	return w.IsOnShelves() && w.value.ReviewState == enum.ReviewPass
}

// 保存
func (w *wholesaleItemImpl) Save() (int32, error) {
	return util.I32Err(w.repo.SaveWsItem(w.value, false))
}

// 是否上架
func (w *wholesaleItemImpl) IsOnShelves() bool {
	return w.value.ShelveState == item.ShelvesOn
}

// 设置上架
func (w *wholesaleItemImpl) SetShelve(state int32, remark string) error {
	if state == item.ShelvesIncorrect && len(remark) == 0 {
		return product.ErrNilRejectRemark
	}
	if state == item.ShelvesOn && w.value.Price <= 0 {
		return item.ErrNotSetWholesalePrice
	}
	w.value.ShelveState = state
	w.value.ReviewRemark = remark
	_, err := w.Save()
	return err
}

// 标记为违规
func (w *wholesaleItemImpl) Incorrect(remark string) error {
	w.value.ShelveState = item.ShelvesIncorrect
	w.value.ReviewRemark = remark
	_, err := w.Save()
	return err
}

// 审核
func (w *wholesaleItemImpl) Review(pass bool, remark string) error {
	if pass {
		w.value.ReviewState = enum.ReviewPass

	} else {
		remark = strings.TrimSpace(remark)
		if remark == "" {
			return item.ErrEmptyReviewRemark
		}
		w.value.ReviewState = enum.ReviewReject
	}
	w.value.ReviewRemark = remark
	_, err := w.Save()
	return err
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
func (w *wholesaleItemImpl) GetWholesalePrice(skuId int64, quantity int32) float64 {
	var price float64 = 0
	arr := w.GetSkuPrice(skuId)
	if len(arr) > 0 {
		var compare int32
		for _, v := range arr {
			if quantity < v.RequireQuantity {
				continue
			}
			if v.RequireQuantity > compare {
				compare = v.RequireQuantity
				price = v.WholesalePrice
			}
		}
	}
	return price
}

// 根据SKU获取价格设置
func (w *wholesaleItemImpl) GetSkuPrice(skuId int64) []*item.WsSkuPrice {
	return w.repo.SelectWsSkuPrice("item_id=? AND sku_id=?",
		w.value.ItemId, skuId)
}

// 保存批发SKU价格设置
func (w *wholesaleItemImpl) SaveSkuPrice(skuId int64, arr []*item.WsSkuPrice) error {
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
	var min, max float64
	for _, v := range arr {
		if min == 0 || max == 0 {
			min = v.WholesalePrice
			max = v.WholesalePrice
		}
		if v.WholesalePrice > max {
			max = v.WholesalePrice
		}
		if v.WholesalePrice < min {
			min = v.WholesalePrice
		}
		// 保存SKU批发价格
		v.ItemId = w.itemId
		v.SkuId = skuId
		i, err := util.I32Err(w.repo.SaveWsSkuPrice(v))
		if err == nil {
			v.ID = i
		}
	}
	// 更新商品批发价格
	if min > 0 && max > 0 {
		w.value.Price = min
		if min == max {
			w.value.PriceRange = format.DecimalToString(min)
		} else {
			w.value.PriceRange = format.DecimalToString(min) +
				"~" + format.DecimalToString(max)
		}
		_, err := w.Save()
		return err
	}
	return nil
}

type itemDetailData struct {
	SpecArray []specJdo `json:"specArray"`
	SkuArray  []skuJdo  `json:"skuArray"`
}

// 获取详细信息
func (w *wholesaleItemImpl) GetJsonDetailData() []byte {
	skuArr := w.it.SkuArray()
	okSkuArr := []*item.Sku{}
	skuJdoArr := []skuJdo{}
	for _, v := range skuArr {
		pArr := w.GetSkuPrice(v.ID)
		if len(pArr) == 0 {
			continue
		}
		okSkuArr = append(okSkuArr, v)
		jdo := skuJdo{
			SkuId:            strconv.Itoa(int(v.ID)),
			SpecData:         v.SpecData,
			SpecWord:         v.SpecWord,
			Price:            math.Round(float64(v.Price), 2),
			DiscountPrice:    math.Round(float64(v.Price), 2),
			CanSalesQuantity: v.Stock,
			SalesCount:       v.SaleNum,
			PriceArray:       []skuPriceJdo{},
		}
		// 如果只包含一个价格，则不返回价格数组
		for j, p := range pArr {
			if j == 0 {
				jdo.Price = p.WholesalePrice
				jdo.DiscountPrice = p.WholesalePrice
				if len(pArr) == 1 {
					break
				}
			}
			jdo.PriceArray = append(jdo.PriceArray, skuPriceJdo{
				Quantity: p.RequireQuantity,
				Price:    p.WholesalePrice,
			})
		}
		skuJdoArr = append(skuJdoArr, jdo)
	}

	spec := w.itemRepo.SkuService().GetSpecArray(okSkuArr)

	i := &itemDetailData{
		SpecArray: iJsonUtil.getSpecJdo(spec),
		SkuArray:  skuJdoArr,
	}
	data, _ := json.MarshalIndent(i, "", " ")
	return data
}
