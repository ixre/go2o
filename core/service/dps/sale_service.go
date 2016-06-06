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
	"go2o/core/domain/interface/sale"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"strconv"
)

type saleService struct {
	_rep        sale.ISaleRep
	_goodsRep   sale.IGoodsRep
	_goodsQuery *query.GoodsQuery
	_cateRep    sale.ICategoryRep
}

func NewSaleService(r sale.ISaleRep, cateRep sale.ICategoryRep,
	goodsRep sale.IGoodsRep, goodsQuery *query.GoodsQuery) *saleService {
	return &saleService{
		_rep:        r,
		_goodsRep:   goodsRep,
		_goodsQuery: goodsQuery,
		_cateRep:    cateRep,
	}
}

// 获取产品值
func (this *saleService) GetValueItem(merchantId, itemId int) *sale.ValueItem {
	sl := this._rep.GetSale(merchantId)
	pro := sl.GetItem(itemId)
	v := pro.GetValue()
	return &v
}

// 获取商品值
func (this *saleService) GetValueGoods(merchantId, goodsId int) *valueobject.Goods {
	sl := this._rep.GetSale(merchantId)
	var goods sale.IGoods = sl.GetGoods(goodsId)
	if goods != nil {
		return goods.GetPackedValue()
	}
	return nil
}

// 根据SKU获取商品
func (this *saleService) GetGoodsBySku(merchantId int, itemId int, sku int) *valueobject.Goods {
	sl := this._rep.GetSale(merchantId)
	var goods sale.IGoods = sl.GetGoodsBySku(itemId, sku)
	return goods.GetPackedValue()
}

// 根据SKU获取商品
func (this *saleService) GetValueGoodsBySku(merchantId int, itemId int, sku int) *sale.ValueGoods {
	sl := this._rep.GetSale(merchantId)
	gs := sl.GetGoodsBySku(itemId, sku)
	if gs != nil {
		return gs.GetValue()
	}
	return nil
}

// 根据快照编号获取商品
func (this *saleService) GetGoodsBySnapshotId(snapshotId int) *sale.ValueGoods {
	snap := this._rep.GetGoodsSnapshot(snapshotId)
	if snap != nil {
		return this._goodsRep.GetValueGoodsById(snap.GoodsId)
	}
	return nil
}

// 保存产品
func (this *saleService) SaveItem(merchantId int, v *sale.ValueItem) (int, error) {
	sl := this._rep.GetSale(merchantId)
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
func (this *saleService) SaveItemInfo(merchantId int, itemId int, info string) error {
	var err error
	sl := this._rep.GetSale(merchantId)
	pro := sl.GetItem(itemId)
	if pro == nil {
		err = errors.New("产品不存在")
	} else {
		v := pro.GetValue()
		v.Description = info
		pro.SetValue(&v)
	}
	_, err = pro.Save()
	return err
}

// 保存商品
func (this *saleService) SaveGoods(merchantId int, gs *sale.ValueGoods) (int, error) {
	sl := this._rep.GetSale(merchantId)
	if gs.Id > 0 {
		g := sl.GetGoods(gs.Id)
		g.SetValue(gs)
		return g.Save()
	}
	g := sl.CreateGoods(gs)
	return g.Save()
}

// 删除货品
func (this *saleService) DeleteItem(merchantId int, id int) error {
	sl := this._rep.GetSale(merchantId)
	return sl.DeleteItem(id)
}

// 获取分页上架的商品
func (this *saleService) GetPagedOnShelvesGoods(merchantId, categoryId, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	var sl sale.ISale = this._rep.GetSale(merchantId)

	if categoryId > 0 {
		var cate sale.ICategory = sl.CategoryManager().GetCategory(categoryId)
		var ids []int = cate.GetChildes()
		ids = append(ids, categoryId)
		total, list = this._goodsRep.GetPagedOnShelvesGoods(merchantId, ids, start, end, "", sortBy)
	} else {
		total = -1
		list = sl.GetOnShelvesGoods(start, end, sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取分页上架的商品
func (this *saleService) GetPagedOnShelvesGoodsByKeyword(merchantId,
	start, end int, word, sortQuery string) (int, []*valueobject.Goods) {
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

	return this._goodsQuery.GetPagedOnShelvesGoodsByKeyword(merchantId,
		start, end, word, where, orderBy)
}

// 删除产品
func (this *saleService) DeleteGoods(merchantId, goodsId int) error {
	sl := this._rep.GetSale(merchantId)
	return sl.DeleteGoods(goodsId)
}

func (this *saleService) GetCategory(merchantId, id int) (*sale.Category, domain.IOptionStore) {
	sl := this._rep.GetSale(merchantId)
	c := sl.CategoryManager().GetCategory(id)
	if c != nil {
		return c.GetValue(), c.GetOption()
	}
	return nil, nil
}

func (this *saleService) DeleteCategory(mchId, id int) error {
	sl := this._rep.GetSale(mchId)
	return sl.CategoryManager().DeleteCategory(id)
}

func (this *saleService) SaveCategory(merchantId int, v *sale.Category) (int, error) {
	sl := this._rep.GetSale(merchantId)
	var ca sale.ICategory
	if v.Id > 0 {
		ca = sl.CategoryManager().GetCategory(v.Id)
		if err := ca.SetValue(v); err != nil {
			return 0, err
		}
	} else {
		ca = sl.CategoryManager().CreateCategory(v)
	}

	return ca.Save()
}

func (this *saleService) GetCategoryTreeNode(merchantId int) *tree.TreeNode {
	sl := this._rep.GetSale(merchantId)
	cats := sl.CategoryManager().GetCategories()
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

func (this *saleService) GetCategories(merchantId int) []*sale.Category {
	sl := this._rep.GetSale(merchantId)
	cats := sl.CategoryManager().GetCategories()
	var list []*sale.Category = make([]*sale.Category, len(cats))
	for i, v := range cats {
		vv := v.GetValue()
		vv.Icon = format.GetResUrl(vv.Icon)
		list[i] = vv
	}
	return list
}

func (this *saleService) getCategoryManager(mchId int) sale.ICategoryManager {
	if mchId > 0 {
		sl := this._rep.GetSale(mchId)
		return sl.CategoryManager()
	}
	return this._cateRep.GetGlobManager()
}

func (this *saleService) GetBigCategories(mchId int) []*sale.Category {
	cats := this.getCategoryManager(mchId).GetCategories()
	list := []*sale.Category{}
	for _, v := range cats {
		if v2 := v.GetValue(); v2.ParentId == 0 && v2.Enabled == 1 {
			v2.Icon = format.GetResUrl(v2.Icon)
			list = append(list, v2)
		}
	}
	return list
}

func (this *saleService) GetChildCategories(mchId, parentId int) []*sale.Category {
	cats := this.getCategoryManager(mchId).GetCategories()
	list := []*sale.Category{}
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == parentId && vv.Enabled == 1 {
			vv.Icon = format.GetResUrl(vv.Icon)
			list = append(list, vv)
			this.setChild(cats, vv)
		}
	}
	return list
}

func (this *saleService) setChild(list []sale.ICategory, dst *sale.Category) {
	for _, v := range list {
		if vv := v.GetValue(); vv.ParentId == dst.Id && vv.Enabled == 1 {
			if dst.Child == nil {
				dst.Child = []*sale.Category{}
			}
			vv.Icon = format.GetResUrl(vv.Icon)
			dst.Child = append(dst.Child, vv)
		}
	}
}

func (this *saleService) GetAllSaleLabels(merchantId int) []*sale.SaleLabel {
	sl := this._rep.GetSale(merchantId)
	tags := sl.GetAllSaleLabels()

	var vtags []*sale.SaleLabel = make([]*sale.SaleLabel, len(tags))
	for i, v := range tags {
		vtags[i] = v.GetValue()
	}
	return vtags
}

// 获取销售标签
func (this *saleService) GetSaleLabel(merchantId, id int) *sale.SaleLabel {
	sl := this._rep.GetSale(merchantId)
	if tag := sl.GetSaleLabel(id); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 获取销售标签
func (this *saleService) GetSaleLabelByCode(merchantId int, code string) *sale.SaleLabel {
	sl := this._rep.GetSale(merchantId)
	if tag := sl.GetSaleLabelByCode(code); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 保存销售标签
func (this *saleService) SaveSaleLabel(merchantId int, v *sale.SaleLabel) (int, error) {
	sl := this._rep.GetSale(merchantId)
	if v.Id > 0 {
		tag := sl.GetSaleLabel(v.Id)
		tag.SetValue(v)
		return tag.Save()
	}
	return sl.CreateSaleLabel(v).Save()
}

// 获取商品的销售标签
func (this *saleService) GetItemSaleLabels(merchantId, itemId int) []*sale.SaleLabel {
	var list = make([]*sale.SaleLabel, 0)
	sl := this._rep.GetSale(merchantId)
	if goods := sl.GetItem(itemId); goods != nil {
		list = goods.GetSaleLabels()
	}
	return list
}

// 保存商品的销售标签
func (this *saleService) SaveItemSaleLabels(merchantId, itemId int, tagIds []int) error {
	var err error
	sl := this._rep.GetSale(merchantId)
	if goods := sl.GetItem(itemId); goods != nil {
		err = goods.SaveSaleLabels(tagIds)
	} else {
		err = errors.New("商品不存在")
	}
	return err
}

// 根据销售标签获取指定数目的商品
func (this *saleService) GetValueGoodsBySaleLabel(merchantId int,
	code, sortBy string, begin int, end int) []*valueobject.Goods {
	sl := this._rep.GetSale(merchantId)
	if tag := sl.GetSaleLabelByCode(code); tag != nil {
		list := tag.GetValueGoods(sortBy, begin, end)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
		}
		return list
	}
	return make([]*valueobject.Goods, 0)
}

// 根据分页销售标签获取指定数目的商品
func (this *saleService) GetPagedValueGoodsBySaleLabel(merchantId int,
	tagId int, sortBy string, begin int, end int) (int, []*valueobject.Goods) {
	sl := this._rep.GetSale(merchantId)
	tag := sl.CreateSaleLabel(&sale.SaleLabel{
		Id: tagId,
	})
	return tag.GetPagedValueGoods(sortBy, begin, end)
}

// 删除销售标签
func (this *saleService) DeleteSaleLabel(merchantId int, id int) error {
	return this._rep.GetSale(merchantId).DeleteSaleLabel(id)
}

// 获取商品的会员价
func (this *saleService) GetGoodsLevelPrices(merchantId, goodsId int) []*sale.MemberPrice {
	sl := this._rep.GetSale(merchantId)
	if goods := sl.GetGoods(goodsId); goods != nil {
		return goods.GetLevelPrices()
	}
	return make([]*sale.MemberPrice, 0)
}

// 保存商品的会员价
func (this *saleService) SaveMemberPrices(merchantId int, goodsId int, priceSet []*sale.MemberPrice) error {
	sl := this._rep.GetSale(merchantId)
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
func (this *saleService) GetGoodsDetails(merchantId, goodsId, mLevel int) (*valueobject.Goods, map[string]string) {
	sl := this._rep.GetSale(merchantId)
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
func (this *saleService) GetItemDescriptionByGoodsId(merchantId, goodsId int) string {
	sl := this._rep.GetSale(merchantId)
	var goods sale.IGoods = sl.GetGoods(goodsId)
	return goods.GetItem().GetValue().Description
}
