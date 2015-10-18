/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package dps

import (
	"errors"
	"fmt"
	"github.com/jsix/gof/web/ui/tree"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/infrastructure/format"
	"strconv"
)

type saleService struct {
	_rep      sale.ISaleRep
	_goodsRep sale.IGoodsRep
}

func NewSaleService(r sale.ISaleRep, goodsRep sale.IGoodsRep) *saleService {
	return &saleService{
		_rep:      r,
		_goodsRep: goodsRep,
	}
}

// 获取产品值
func (this *saleService) GetValueItem(partnerId, itemId int) *sale.ValueItem {
	sl := this._rep.GetSale(partnerId)
	pro := sl.GetItem(itemId)
	v := pro.GetValue()
	return &v
}

// 获取商品值
func (this *saleService) GetValueGoods(partnerId, goodsId int) *valueobject.Goods {
	sl := this._rep.GetSale(partnerId)
	var goods sale.IGoods = sl.GetGoods(goodsId)
	return goods.GetPackedValue()
}

// 根据SKU获取商品
func (this *saleService) GetGoodsBySku(partnerId int, itemId int, sku int) *valueobject.Goods {
	sl := this._rep.GetSale(partnerId)
	var goods sale.IGoods = sl.GetGoodsBySku(itemId, sku)
	return goods.GetPackedValue()
}

// 保存产品
func (this *saleService) SaveItem(partnerId int, v *sale.ValueItem) (int, error) {
	sl := this._rep.GetSale(partnerId)
	var pro sale.IItem
	if v.Id > 0 {
		pro = sl.GetItem(v.Id)
		if pro == nil {
			return 0, errors.New("产品不存在")
		}

		// 修改货品时，不会修改详情
		v.Description = pro.GetValue().Description

		if err := pro.SetValue(v); err != nil {
			return 0, err
		}
	} else {
		pro = sl.CreateItem(v)
	}
	return pro.Save()
}

// 保存货品描述
func (this *saleService) SaveItemInfo(partnerId int,itemId int,info string)error {
	var err error
	sl := this._rep.GetSale(partnerId)
	pro := sl.GetItem(itemId)
	if pro == nil {
		err = errors.New("产品不存在")
	}else {
		v := pro.GetValue()
		v.Description = info
		pro.SetValue(&v)
	}
	_,err = pro.Save()
	return err
}
// 删除货品
func (this *saleService) DeleteItem(partnerId int, id int) error {
	sl := this._rep.GetSale(partnerId)
	return sl.DeleteItem(id)
}

// 获取分页上架的商品
func (this *saleService) GetPagedOnShelvesGoods(partnerId, categoryId, start, end int,
	sortQuery string) (int, []*valueobject.Goods) {
	var sl sale.ISale = this._rep.GetSale(partnerId)
	var cate sale.ICategory = sl.GetCategory(categoryId)
	var ids []int = cate.GetChildId()
	ids = append(ids, categoryId)
	//todo: cache

	var where string
	var orderBy string
	switch sortQuery {
	case "price_0":
		where = ""
		orderBy = "gs_item.sale_price ASC"
	case "price_1":
		where = ""
		orderBy = "gs_item.sale_price DESC"
	case "sale_0":
		where = ""
		orderBy = "gs_goods.sale_num ASC"
	case "sale_1":
		where = ""
		orderBy = "gs_goods.sale_num DESC"
	case "rate_0":
		//todo:
	case "rate_1":
		//todo:
	}

	return this._goodsRep.GetPagedOnShelvesGoods(partnerId, ids, start, end, where, orderBy)

}

// 删除产品
func (this *saleService) DeleteGoods(partnerId, goodsId int) error {
	sl := this._rep.GetSale(partnerId)
	return sl.DeleteGoods(goodsId)
}

func (this *saleService) GetCategory(partnerId, id int) *sale.ValueCategory {
	sl := this._rep.GetSale(partnerId)
	c := sl.GetCategory(id)
	if c != nil {
		cv := c.GetValue()
		return &cv
	}
	return nil
}

func (this *saleService) DeleteCategory(partnerId, id int) error {
	sl := this._rep.GetSale(partnerId)
	return sl.DeleteCategory(id)
}

func (this *saleService) SaveCategory(partnerId int, v *sale.ValueCategory) (int, error) {
	sl := this._rep.GetSale(partnerId)
	var ca sale.ICategory
	if v.Id > 0 {
		ca = sl.GetCategory(v.Id)
		if err := ca.SetValue(v); err != nil {
			return 0, err
		}
	} else {
		ca = sl.CreateCategory(v)
	}

	return ca.Save()
}

func (this *saleService) GetCategoryTreeNode(partnerId int) *tree.TreeNode {
	sl := this._rep.GetSale(partnerId)
	cats := sl.GetCategories()
	rootNode := &tree.TreeNode{
		Text:   "根节点",
		Value:  "",
		Url:    "",
		Icon:   "",
		Open:   true,
		Childs: nil}
	this.walkCategoryTree(rootNode, 0, cats)
	return rootNode
}

func (this *saleService) walkCategoryTree(node *tree.TreeNode, parentId int, categories []sale.ICategory) {
	node.Childs = []*tree.TreeNode{}
	for _, v := range categories {
		cate := v.GetValue()
		if cate.ParentId == parentId {
			cNode := &tree.TreeNode{
				Text:   cate.Name,
				Value:  strconv.Itoa(cate.Id),
				Url:    "",
				Icon:   "",
				Open:   true,
				Childs: nil}
			node.Childs = append(node.Childs, cNode)
			this.walkCategoryTree(cNode, cate.Id, categories)
		}
	}
}

func (this *saleService) GetCategories(partnerId int) []*sale.ValueCategory {
	sl := this._rep.GetSale(partnerId)
	cats := sl.GetCategories()
	var list []*sale.ValueCategory = make([]*sale.ValueCategory, len(cats))
	for i, v := range cats {
		vv := v.GetValue()
		list[i] = &vv
	}
	return list
}

func (this *saleService) GetAllSaleTags(partnerId int) []*sale.ValueSaleTag {
	sl := this._rep.GetSale(partnerId)
	tags := sl.GetAllSaleTags()

	var vtags []*sale.ValueSaleTag = make([]*sale.ValueSaleTag, len(tags))
	for i, v := range tags {
		vtags[i] = v.GetValue()
	}
	return vtags
}

// 获取销售标签
func (this *saleService) GetSaleTag(partnerId, id int) *sale.ValueSaleTag {
	sl := this._rep.GetSale(partnerId)
	if tag := sl.GetSaleTag(id); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 获取销售标签
func (this *saleService) GetSaleTagByCode(partnerId int, code string) *sale.ValueSaleTag {
	sl := this._rep.GetSale(partnerId)
	if tag := sl.GetSaleTagByCode(code); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 保存销售标签
func (this *saleService) SaveSaleTag(partnerId int, v *sale.ValueSaleTag) (int, error) {
	sl := this._rep.GetSale(partnerId)
	if v.Id > 0 {
		tag := sl.GetSaleTag(v.Id)
		tag.SetValue(v)
		return tag.Save()
	}
	return sl.CreateSaleTag(v).Save()
}

// 获取商品的销售标签
func (this *saleService) GetItemSaleTags(partnerId, itemId int) []*sale.ValueSaleTag {
	var list = make([]*sale.ValueSaleTag, 0)
	sl := this._rep.GetSale(partnerId)
	if goods := sl.GetItem(itemId); goods != nil {
		list = goods.GetSaleTags()
	}
	return list
}

// 保存商品的销售标签
func (this *saleService) SaveItemSaleTags(partnerId, itemId int, tagIds []int) error {
	var err error
	sl := this._rep.GetSale(partnerId)
	if goods := sl.GetItem(itemId); goods != nil {
		err = goods.SaveSaleTags(tagIds)
	} else {
		err = errors.New("商品不存在")
	}
	return err
}

// 根据销售标签获取指定数目的商品
func (this *saleService) GetValueGoodsBySaleTag(partnerId int, code string, begin int, end int) []*valueobject.Goods {
	sl := this._rep.GetSale(partnerId)
	if tag := sl.GetSaleTagByCode(code); tag != nil {
		return tag.GetValueGoods(begin, end)
	}
	return make([]*valueobject.Goods, 0)
}

// 根据分页销售标签获取指定数目的商品
func (this *saleService) GetPagedValueGoodsBySaleTag(partnerId int, tagId int, begin int, end int) (int, []*valueobject.Goods) {
	sl := this._rep.GetSale(partnerId)
	tag := sl.CreateSaleTag(&sale.ValueSaleTag{
		Id: tagId,
	})
	return tag.GetPagedValueGoods(begin, end)
}

// 删除销售标签
func (this *saleService) DeleteSaleTag(partnerId int, id int) error {
	return this._rep.GetSale(partnerId).DeleteSaleTag(id)
}

// 获取商品的会员价
func (this *saleService) GetGoodsLevelPrices(partnerId, goodsId int) []*sale.MemberPrice {
	sl := this._rep.GetSale(partnerId)
	if goods := sl.GetGoods(goodsId); goods != nil {
		return goods.GetLevelPrices()
	}
	return make([]*sale.MemberPrice, 0)
}

// 保存商品的会员价
func (this *saleService) SaveMemberPrices(partnerId int, goodsId int, priceSet []*sale.MemberPrice) error {
	sl := this._rep.GetSale(partnerId)
	var err error
	if goods := sl.GetGoods(goodsId); goods != nil {
		for _, v := range priceSet {
			if _, err = goods.SaveLevelPrice(v); err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取商品详情
func (this *saleService) GetGoodsDetails(partnerId, goodsId, mLevel int) (*valueobject.Goods, map[string]string) {
	sl := this._rep.GetSale(partnerId)
	var goods sale.IGoods = sl.GetGoods(goodsId)
	gv := goods.GetPackedValue()
	proMap := goods.GetPromotionDescribe()
	if b, price := goods.GetLevelPrice(mLevel); b {
		gv.PromPrice = price
		proMap["会员专享"] = fmt.Sprintf("会员优惠,仅需<b>￥%s</b>",
			format.FormatFloat(price))
	}
	return gv, proMap
}

// 获取货品描述
func (this *saleService) GetItemDescriptionByGoodsId(partnerId,goodsId int)string{
	sl := this._rep.GetSale(partnerId)
	var goods sale.IGoods = sl.GetGoods(goodsId)
	return goods.GetItem().GetValue().Description
}