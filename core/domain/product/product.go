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
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
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
	pmRepo    promodel.IProModelRepo
	valueRepo valueobject.IValueRepo
}

func NewProductImpl(v *product.Product,
	itemRepo product.IProductRepo, pmRepo promodel.IProModelRepo,
	valRepo valueobject.IValueRepo) product.IProduct {
	return &productImpl{
		value:     v,
		repo:      itemRepo,
		pmRepo:    pmRepo,
		valueRepo: valRepo,
	}
}

// 获取聚合根编号
func (p *productImpl) GetAggregateRootId() int64 {
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
	unix := time.Now().Unix()
	if p.value.CreateTime <= 0 {
		p.value.CreateTime = unix
	}
	p.value.UpdateTime = unix
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
func (p *productImpl) Save() (i int64, err error) {
	if p.value.Attr != nil {
		if p.GetAggregateRootId() <= 0 {
			p.value.Id, err = util.I64Err(p.repo.SaveProduct(p.value))
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
		unix := time.Now().Unix()
		cs := strconv.Itoa(int(p.value.CatId))
		us := strconv.Itoa(int(unix))
		l := len(cs)
		p.value.Code = fmt.Sprintf("%s%s", cs, us[4+l:])
	}
	p.value.Id, err = util.I64Err(p.repo.SaveProduct(p.value))
R:
	return p.value.Id, err
}

// 合并属性
func (p *productImpl) mergeAttr(src []*product.Attr, dst *[]*product.Attr) {
	if src == nil || dst == nil || len(src) == 0 || len(*dst) == 0 {
		return
	}
	to := *dst
	sMap := make(map[int32]int32, len(src))
	for _, v := range src {
		sMap[v.AttrId] = v.ID
	}
	for _, v := range to {
		if id, ok := sMap[v.AttrId]; ok {
			v.ID = id
		}
	}
}

// 重建Attr数组，将信息附加
func (p *productImpl) RebuildAttrArray(arr *[]*product.Attr) error {
	for _, v := range *arr {
		vArr := util.StrExt.I32Slice(v.AttrData, ",")
		for i, v2 := range vArr {
			if i != 0 {
				v.AttrWord += ","
			}
			item := p.pmRepo.GetAttrItem(v2)
			if item != nil {
				v.AttrWord += item.Value
			}
		}
	}
	return nil
}

// 保存属性
func (p *productImpl) saveAttr(arr []*product.Attr) (err error) {
	pk := p.GetAggregateRootId()
	// 获取之前的SKU设置
	old := p.repo.SelectAttr("product_id=?", pk)
	// 合并属性
	p.mergeAttr(old, &p.value.Attr)
	// 设置属性值
	if err = p.RebuildAttrArray(&arr); err != nil {
		return err
	}
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*product.Attr, len(arr))
	for _, v := range arr {
		currMap[v.ID] = v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.ID] == nil {
			delList = append(delList, v.ID)
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
		if v.ProductId == pk && v.AttrData != "" {
			v.ID, err = util.I32Err(p.repo.SaveAttr(v))
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
//	if v.ID <= 0 {
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
//		Key:          fmt.Sprintf("%d-g%d-%d", mchId, v.ID, unix),
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
