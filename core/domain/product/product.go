/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:53
 * description :
 * history :
 */

package product

import (
	"fmt"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"strconv"
	"strings"
	"time"
)

var _ product.IProduct = new(productImpl)

type productImpl struct {
	value       *product.Product
	productRepo product.IProductRepo
	valueRepo   valueobject.IValueRepo
}

func NewProductImpl(v *product.Product,
	itemRepo product.IProductRepo,
	valRepo valueobject.IValueRepo) product.IProduct {
	return &productImpl{
		value:       v,
		productRepo: itemRepo,
		valueRepo:   valRepo,
	}
}

// 获取聚合根编号
func (i *productImpl) GetAggregateRootId() int32 {
	return i.value.Id
}

func (i *productImpl) GetValue() product.Product {
	return *i.value
}

func (i *productImpl) checkValue(v *product.Product) error {

	// 检测供应商
	if v.VendorId <= 0 || v.VendorId != i.value.VendorId {
		return product.ErrVendor
	}
	// 检测标题长度
	v.Name = strings.TrimSpace(v.Name)
	if len(v.Name) < 10 {
		return product.ErrItemNameLength
	}

	// 检测品牌
	if v.BrandId <= 0 {
		//todo: 检测是否有效，与模型是否匹配
		return product.ErrNoBrand
	}
	return nil

}

// 设置值
func (i *productImpl) SetValue(v *product.Product) error {
	if i.GetAggregateRootId() <= 0 {
		i.value.ShelveState = product.ShelvesDown
		i.value.ReviewState = enum.ReviewAwaiting
	}
	if i.value.ShelveState == product.ShelvesIncorrect {
		return product.ErrItemIncorrect
	}
	if err := i.checkValue(v); err != nil {
		return err
	}
	if v.Id == i.value.Id {
		//修改图片或标题后，要重新审核
		if i.value.Image != v.Image || i.value.Name != v.Name {
			i.resetReview()
		}
		i.value.Name = v.Name
		i.value.Code = v.Code
		i.value.BrandId = v.BrandId
		i.value.Image = v.Image
		if v.CategoryId > 0 {
			i.value.CategoryId = v.CategoryId
		}
		i.value.SortNum = v.SortNum
	}
	i.value.UpdateTime = time.Now().Unix()
	return nil
}

// 设置商品描述
func (i *productImpl) SetDescribe(describe string) error {
	if len(describe) < 20 {
		return product.ErrDescribeLength
	}
	if i.value.Description != describe {
		i.value.Description = describe
		i.resetReview()
		_, err := i.Save()
		return err
	}
	return nil
}

// 是否上架
func (i *productImpl) IsOnShelves() bool {
	return i.value.ShelveState == product.ShelvesOn
}

// 获取商品的销售标签
//func (i *itemImpl) GetSaleLabels() []*item.Label {
//    if i.saleLabels == nil {
//        i.saleLabels = i.saleLabelRepo.GetItemSaleLabels(i.GetAggregateRootId())
//    }
//    return i.saleLabels
//}
//
//// 保存销售标签
//func (i *itemImpl) SaveSaleLabels(tagIds []int) error {
//    err := i.saleLabelRepo.CleanItemSaleLabels(i.GetAggregateRootId())
//    if err == nil {
//        err = i.saleLabelRepo.SaveItemSaleLabels(i.GetAggregateRootId(), tagIds)
//        i.saleLabels = nil
//    }
//    return err
//}

// 重置审核状态
func (i *productImpl) resetReview() {
	i.value.ReviewState = enum.ReviewAwaiting
}

// 设置上架
func (i *productImpl) SetShelve(state int32, remark string) error {
	if state == product.ShelvesIncorrect && len(remark) == 0 {
		return product.ErrNilRejectRemark
	}
	i.value.ShelveState = state
	i.value.Remark = remark
	_, err := i.Save()
	return err
}

// 审核
func (i *productImpl) Review(pass bool, remark string) error {
	if pass {
		i.value.ReviewState = enum.ReviewPass

	} else {
		remark = strings.TrimSpace(remark)
		if remark == "" {
			return item.ErrEmptyReviewRemark
		}
		i.value.ReviewState = enum.ReviewReject
	}
	i.value.Remark = remark
	_, err := i.Save()
	return err
}

// 标记为违规
func (i *productImpl) Incorrect(remark string) error {
	i.value.ShelveState = product.ShelvesIncorrect
	i.value.Remark = remark
	_, err := i.Save()
	return err
}

// 保存
func (i *productImpl) Save() (int32, error) {
	unix := time.Now().Unix()
	i.value.UpdateTime = unix
	if i.GetAggregateRootId() <= 0 {
		i.value.CreateTime = unix
	}
	if i.value.Code == "" {
		cs := strconv.Itoa(int(i.value.CategoryId))
		us := strconv.Itoa(int(unix))
		l := len(cs)
		i.value.Code = fmt.Sprintf("%s%s", cs, us[4+l:])
	}

	id, err := i.productRepo.SaveProductValue(i.value)
	if err == nil {
		i.value.Id = id

		//todo: 保存商品
		i.saveGoods()

		// 创建快照
		//_, err = i.GenerateSnapshot()
	}
	return id, err
}

//todo: 过渡方法,应有SKU,不根据Item生成Goods
func (i *productImpl) saveGoods() {
	//val := i.goodsRepo.GetValueGoods(i.GetAggregateRootId(), 0)
	//if val == nil {
	//    val = &item.ItemGoods{
	//        Id:            0,
	//        ProductId:     i.GetAggregateRootId(),
	//        IsPresent:     0,
	//        SkuId:         0,
	//        PromotionFlag: 0,
	//        StockNum:      100,
	//        SaleNum:       100,
	//    }
	//}

	panic(" save goods ")
	//goods := NewSaleItem(i.itemRepo, i, val,
	//    i.saleRepo, i.goodsRepo, i.promRepo)
	//goods.Save()
}

// 销毁产品
func (i *productImpl) Destroy() error {
	num := i.productRepo.GetProductSaleNum(i.GetAggregateRootId())
	if num > 0 {
		return item.ErrCanNotDeleteItem
	}
	return i.productRepo.DeleteProduct(i.GetAggregateRootId())
}

//// 生成快照
//func (i *Goods) GenerateSnapshot() (int64, error) {
//	v := i._value
//	if v.Id <= 0 {
//		return 0, item.ErrNoSuchGoods
//	}
//
//	if v.OnShelves == 0 {
//		return 0, item.ErrNotOnShelves
//	}
//
//	mchId := i._sale.GetAggregateRootId()
//	unix := time.Now().Unix()
//	cate := i._saleRepo.GetCategory(mchId, v.CategoryId)
//	var gsn *goods.GoodsSnapshot = &goods.GoodsSnapshot{
//		Key:          fmt.Sprintf("%d-g%d-%d", mchId, v.Id, unix),
//		GoodsId:      i.GetAggregateRootId(),
//		GoodsName:    v.Name,
//		GoodsNo:      v.GoodsNo,
//		SmallTitle:   v.SmallTitle,
//		CategoryName: cate.Name,
//		Image:        v.Image,
//		Cost:         v.Cost,
//		Price:        v.Price,
//		SalePrice:    v.SalePrice,
//		CreateTime:   unix,
//	}
//
//	if i.isNewSnapshot(gsn) {
//		i._latestSnapshot = gsn
//		return i._saleRepo.SaveSnapshot(gsn)
//	}
//	return 0, item.ErrLatestSnapshot
//}
//
//// 是否为新快照,与旧有快照进行数据对比
//func (i *Goods) isNewSnapshot(gsn *goods.GoodsSnapshot) bool {
//	latestGsn := i.GetLatestSnapshot()
//	if latestGsn != nil {
//		return latestGsn.GoodsName != gsn.GoodsName ||
//			latestGsn.SmallTitle != gsn.SmallTitle ||
//			latestGsn.CategoryName != gsn.CategoryName ||
//			latestGsn.Image != gsn.Image ||
//			latestGsn.Cost != gsn.Cost ||
//			latestGsn.Price != gsn.Price ||
//			latestGsn.SalePrice != gsn.SalePrice
//	}
//	return true
//}
//
//// 获取最新的快照
//func (i *Goods) GetLatestSnapshot() *goods.GoodsSnapshot {
//	if i._latestSnapshot == nil {
//		i._latestSnapshot = i._saleRepo.GetLatestGoodsSnapshot(i.GetAggregateRootId())
//	}
//	return i._latestSnapshot
//}
