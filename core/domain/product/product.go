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
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"strconv"
	"strings"
	"time"
)

var _ product.IProduct = new(productImpl)

type productImpl struct {
	value     *product.Product
	repo      product.IProductRepo
	valueRepo valueobject.IValueRepo
}

func NewProductImpl(v *product.Product,
	itemRepo product.IProductRepo,
	valRepo valueobject.IValueRepo) product.IProduct {
	return &productImpl{
		value:     v,
		repo:      itemRepo,
		valueRepo: valRepo,
	}
}

// 获取聚合根编号
func (p *productImpl) GetAggregateRootId() int32 {
	return p.value.Id
}

func (p *productImpl) GetValue() product.Product {
	return *p.value
}

func (p *productImpl) checkValue(v *product.Product) error {

	// 检测供应商
	if v.VendorId <= 0 || v.VendorId != p.value.VendorId {
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
func (p *productImpl) SetValue(v *product.Product) error {
	//if p.GetAggregateRootId() <= 0 {
	//    p.value.ShelveState = item.ShelvesDown
	//    p.value.ReviewState = enum.ReviewAwaiting
	//}
	//if p.value.ShelveState == item.ShelvesIncorrect {
	//    return product.ErrItemIncorrect
	//}
	if err := p.checkValue(v); err != nil {
		return err
	}
	if v.Id == p.value.Id {
		p.value.Name = v.Name
		p.value.Code = v.Code
		p.value.BrandId = v.BrandId
		p.value.Image = v.Image
		if v.CatId > 0 {
			p.value.CatId = v.CatId
		}
		p.value.SortNum = v.SortNum
	}
	p.value.UpdateTime = time.Now().Unix()
	return nil
}

// 设置产品属性
func (p *productImpl) SetAttr(attrs []*product.Attr) error {
	if attrs == nil {
		return product.ErrNoSuchAttr
	}
	p.value.Attr = attrs
	return nil
}

// 获取属性
func (p *productImpl) Attr() []*product.Attr {
	if p.value.Attr == nil {
		p.value.Attr = p.repo.SelectAttr("product_id=?",
			p.GetAggregateRootId())
	}
	return p.value.Attr
}

// 设置商品描述
func (p *productImpl) SetDescribe(describe string) error {
	if len(describe) < 20 {
		return product.ErrDescribeLength
	}
	if p.value.Description != describe {
		p.value.Description = describe
		_, err := p.Save()
		return err
	}
	return nil
}

// 获取商品的销售标签
//func (i *itemImpl) GetSaleLabels() []*item.Label {
//    if i.saleLabels == nil {
//        i.saleLabels = i.saleLabelRepo.GetItemSaleLabels(p.GetAggregateRootId())
//    }
//    return i.saleLabels
//}
//
//// 保存销售标签
//func (i *itemImpl) SaveSaleLabels(tagIds []int) error {
//    err := i.saleLabelRepo.CleanItemSaleLabels(p.GetAggregateRootId())
//    if err == nil {
//        err = i.saleLabelRepo.SaveItemSaleLabels(p.GetAggregateRootId(), tagIds)
//        i.saleLabels = nil
//    }
//    return err
//}

// 保存
func (p *productImpl) Save() (i int32, err error) {
	unix := time.Now().Unix()
	p.value.UpdateTime = unix
	if p.value.Attr != nil {
		if p.GetAggregateRootId() <= 0 {
			p.value.CreateTime = unix
			p.value.Id, err = util.I32Err(p.repo.SaveProduct(p.value))
			if err != nil {
				goto R
			}
		}
		if err = p.saveAttr(p.value.Attr); err != nil {
			goto R
		}
	}
	// 自动生成货号
	if p.value.Code == "" {
		cs := strconv.Itoa(int(p.value.CatId))
		us := strconv.Itoa(int(unix))
		l := len(cs)
		p.value.Code = fmt.Sprintf("%s%s", cs, us[4+l:])
	}
	p.value.Id, err = util.I32Err(p.repo.SaveProduct(p.value))
R:
	return p.value.Id, err
}

// 保存属性
func (p *productImpl) saveAttr(arr []*product.Attr) (err error) {
	pk := p.GetAggregateRootId()
	// 获取之前的SKU设置
	old := p.repo.SelectAttr("product_id=?", pk)
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*product.Attr, len(arr))
	for _, v := range arr {
		currMap[v.Id] = v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.Id] == nil {
			delList = append(delList, v.Id)
		}
	}
	// 删除项
	for _, v := range delList {
		p.repo.DeleteAttr(v)
	}
	// 保存项
	for _, v := range arr {
		if v.ProductId == 0 {
			v.ProductId = pk
		}
		if v.ProductId == pk {
			v.Id, err = util.I32Err(p.repo.SaveAttr(v))
		}
	}
	return err
}

// 销毁产品
func (p *productImpl) Destroy() error {
	num := p.repo.GetProductSaleNum(p.GetAggregateRootId())
	if num > 0 {
		return item.ErrCanNotDeleteItem
	}
	return p.repo.DeleteProduct(p.GetAggregateRootId())
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
//		GoodsId:      p.GetAggregateRootId(),
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
//		i._latestSnapshot = i._saleRepo.GetLatestGoodsSnapshot(p.GetAggregateRootId())
//	}
//	return i._latestSnapshot
//}
