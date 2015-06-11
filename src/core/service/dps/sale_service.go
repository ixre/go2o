/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package dps

import (
	"errors"
	"github.com/atnet/gof/web/ui/tree"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/dto"
	"strconv"
)

type saleService struct {
	_rep sale.ISaleRep
}

func NewSaleService(r sale.ISaleRep) *saleService {
	return &saleService{
		_rep: r,
	}
}

func (this *saleService) GetValueGoods(partnerId, goodsId int) *sale.ValueGoods {
	sl := this._rep.GetSale(partnerId)
	pro := sl.GetGoods(goodsId)
	v := pro.GetValue()
	return &v
}

func (this *saleService) SaveGoods(partnerId int, v *sale.ValueGoods) (int, error) {
	sl := this._rep.GetSale(partnerId)
	var pro sale.IGoods
	if v.Id > 0 {
		pro = sl.GetGoods(v.Id)
		if pro == nil {
			return 0, errors.New("产品不存在")
		}
		if err := pro.SetValue(v); err != nil {
			return 0, err
		}
	} else {
		pro = sl.CreateGoods(v)
	}
	return pro.Save()
}

func (this *saleService) GetPagedOnShelvesGoods(partnerId, categoryId, start,end int)(int,[]*dto.ListGoods){
	var sl sale.ISale = this._rep.GetSale(partnerId)
	var cate sale.ICategory = sl.GetCategory(categoryId)
	var ids []int = cate.GetChildId()
	ids = append(ids, categoryId)

	//todo: cache

	total,goods := this._rep.GetPagedOnShelvesGoods(partnerId, ids, start,end)
	var listGoods []*dto.ListGoods = make([]*dto.ListGoods, len(goods))
	for i, v := range goods {
		listGoods[i] = &dto.ListGoods{
			Id:         v.Id,
			Name:       v.Name,
			SmallTitle: v.SmallTitle,
			Image:      v.Image,
			Price:      v.Price,
			SalePrice:  v.SalePrice,
		}
	}
	return total,listGoods
}

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
	this.iterCategoryTree(rootNode, 0, cats)
	return rootNode
}

func (this *saleService) iterCategoryTree(node *tree.TreeNode, parentId int, categories []sale.ICategory) {
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
			this.iterCategoryTree(cNode, cate.Id, categories)
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

// 初始化销售标签
func (this *saleService) InitSaleTags(partnerId int) error {
	sl := this._rep.GetSale(partnerId)
	return sl.InitSaleTags()
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
func (this *saleService) GetGoodsSaleTags(partnerId, goodsId int) []*sale.ValueSaleTag {
	var list = make([]*sale.ValueSaleTag, 0)
	sl := this._rep.GetSale(partnerId)
	if goods := sl.GetGoods(goodsId); goods != nil {
		list = goods.GetSaleTags()
	}
	return list
}

// 保存商品的销售标签
func (this *saleService) SaveGoodsSaleTags(partnerId, goodsId int, tagIds []int) error {
	var err error
	sl := this._rep.GetSale(partnerId)
	if goods := sl.GetGoods(goodsId); goods != nil {
		err = goods.SaveSaleTags(tagIds)
	} else {
		err = errors.New("商品不存在")
	}
	return err
}

// 根据销售标签获取指定数目的商品
func (this *saleService) GetValueGoodsBySaleTag(partnerId int, code string, begin int, end int) []*sale.ValueGoods {
	sl := this._rep.GetSale(partnerId)
	if tag := sl.GetSaleTagByCode(code); tag != nil {
		return tag.GetValueGoods(begin, end)
	}
	return make([]*sale.ValueGoods, 0)
}
