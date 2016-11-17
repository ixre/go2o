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
	"go2o/core/domain/interface/sale/goods"
	"go2o/core/domain/interface/sale/item"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/dto"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"strconv"
)

type saleService struct {
	_rep        sale.ISaleRep
	_goodsRep   goods.IGoodsRep
	_goodsQuery *query.GoodsQuery
	_cateRep    sale.ICategoryRep
}

func NewSaleService(r sale.ISaleRep, cateRep sale.ICategoryRep,
	goodsRep goods.IGoodsRep, goodsQuery *query.GoodsQuery) *saleService {
	return &saleService{
		_rep:        r,
		_goodsRep:   goodsRep,
		_goodsQuery: goodsQuery,
		_cateRep:    cateRep,
	}
}

// 获取产品值
func (s *saleService) GetValueItem(supplierId, itemId int) *item.Item {
	sl := s._rep.GetSale(supplierId)
	pro := sl.ItemManager().GetItem(itemId)
	if pro != nil {
		v := pro.GetValue()
		return &v
	}
	return nil
}

// 获取商品值
func (s *saleService) GetValueGoods(merchantId, goodsId int) *valueobject.Goods {
	sl := s._rep.GetSale(merchantId)
	var goods sale.IGoods = sl.GoodsManager().GetGoods(goodsId)
	if goods != nil {
		return goods.GetPackedValue()
	}
	return nil
}

// 根据SKU获取商品
func (s *saleService) GetGoodsBySku(merchantId int, itemId int, sku int) *valueobject.Goods {
	sl := s._rep.GetSale(merchantId)
	var goods sale.IGoods = sl.GoodsManager().GetGoodsBySku(itemId, sku)
	return goods.GetPackedValue()
}

// 根据SKU获取商品
func (s *saleService) GetValueGoodsBySku(merchantId int, itemId int, sku int) *goods.ValueGoods {
	sl := s._rep.GetSale(merchantId)
	gs := sl.GoodsManager().GetGoodsBySku(itemId, sku)
	if gs != nil {
		return gs.GetValue()
	}
	return nil
}

// 根据快照编号获取商品
func (s *saleService) GetGoodsBySnapshotId(snapshotId int) *goods.ValueGoods {
	snap := s._goodsRep.GetSaleSnapshot(snapshotId)
	if snap != nil {
		return s._goodsRep.GetValueGoodsById(snap.SkuId)
	}
	return nil
}

// 根据快照编号获取商品
func (s *saleService) GetSaleSnapshotById(snapshotId int) *goods.SalesSnapshot {
	return s._goodsRep.GetSaleSnapshot(snapshotId)
}

// 保存产品
func (s *saleService) SaveItem(vendorId int, v *item.Item) (int64, error) {
	sl := s._rep.GetSale(vendorId)
	var pro sale.IItem
	v.VendorId = vendorId //设置供应商编号
	if v.Id > 0 {
		pro = sl.ItemManager().GetItem(v.Id)
		if pro == nil || pro.GetValue().VendorId != vendorId {
			return 0, errors.New("产品不存在")
		}
		// 修改货品时，不会修改详情
		v.Description = pro.GetValue().Description
	} else {
		pro = sl.ItemManager().CreateItem(v)
	}
	if err := pro.SetValue(v); err != nil {
		return 0, err
	}
	return pro.Save()
}

// 保存货品描述
func (s *saleService) SaveItemInfo(merchantId int, itemId int, info string) error {
	sl := s._rep.GetSale(merchantId)
	pro := sl.ItemManager().GetItem(itemId)
	if pro == nil {
		return goods.ErrNoSuchGoods
	}
	return pro.SetDescribe(info)
}

// 保存商品
func (s *saleService) SaveGoods(merchantId int, gs *goods.ValueGoods) (int, error) {
	sl := s._rep.GetSale(merchantId)
	if gs.Id > 0 {
		g := sl.GoodsManager().GetGoods(gs.Id)
		g.SetValue(gs)
		return g.Save()
	}
	g := sl.GoodsManager().CreateGoods(gs)
	return g.Save()
}

// 删除货品
func (s *saleService) DeleteItem(merchantId int, id int) error {
	sl := s._rep.GetSale(merchantId)
	return sl.ItemManager().DeleteItem(id)
}

// 获取分页上架的商品
func (s *saleService) GetShopPagedOnShelvesGoods(shopId, categoryId, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cat := s._cateRep.GetGlobManager().GetCategory(categoryId)
		ids := cat.GetChildes()
		ids = append(ids, categoryId)
		total, list = s._goodsRep.GetPagedOnShelvesGoods(shopId, ids, start, end, "", sortBy)
	} else {
		total, list = s._goodsRep.GetPagedOnShelvesGoods(shopId, nil, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取上架商品数据（分页）
func (s *saleService) GetPagedOnShelvesGoods(shopId int, categoryId, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cate := s._cateRep.GetGlobManager().GetCategory(categoryId)
		var ids []int = cate.GetChildes()
		ids = append(ids, categoryId)
		total, list = s._goodsRep.GetPagedOnShelvesGoods(shopId,
			ids, start, end, "", sortBy)
	} else {
		total, list = s._goodsRep.GetPagedOnShelvesGoods(shopId,
			[]int{}, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取分页上架的商品
func (s *saleService) GetPagedOnShelvesGoodsByKeyword(shopId, start, end int,
	word, sortQuery string) (int, []*valueobject.Goods) {
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
	return s._goodsQuery.GetPagedOnShelvesGoodsByKeyword(shopId,
		start, end, word, where, orderBy)
}

// 删除产品
func (s *saleService) DeleteGoods(merchantId, goodsId int) error {
	sl := s._rep.GetSale(merchantId)
	return sl.GoodsManager().DeleteGoods(goodsId)
}

// 获取商品分类
func (s *saleService) GetCategory(merchantId, id int) *sale.Category {
	sl := s._rep.GetSale(merchantId)
	c := sl.CategoryManager().GetCategory(id)
	if c != nil {
		return c.GetValue()
	}
	return nil
}

// 获取商品分类和选项
func (s *saleService) GetCategoryAndOptions(merchantId, id int) (*sale.Category,
	domain.IOptionStore) {
	sl := s._rep.GetSale(merchantId)
	c := sl.CategoryManager().GetCategory(id)
	if c != nil {
		return c.GetValue(), c.GetOption()
	}
	return nil, nil
}

func (s *saleService) DeleteCategory(mchId, id int) error {
	sl := s._rep.GetSale(mchId)
	return sl.CategoryManager().DeleteCategory(id)
}

func (s *saleService) SaveCategory(merchantId int, v *sale.Category) (int, error) {
	sl := s._rep.GetSale(merchantId)
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

func (s *saleService) GetCategoryTreeNode(merchantId int) *tree.TreeNode {
	sl := s._rep.GetSale(merchantId)
	cats := sl.CategoryManager().GetCategories()
	rootNode := &tree.TreeNode{
		Text:   "根节点",
		Value:  "",
		Url:    "",
		Icon:   "",
		Open:   true,
		Childs: nil}
	s.walkCategoryTree(rootNode, 0, cats)
	return rootNode
}

func (s *saleService) walkCategoryTree(node *tree.TreeNode, parentId int, categories []sale.ICategory) {
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
			s.walkCategoryTree(cNode, cate.Id, categories)
		}
	}
}

func (s *saleService) GetCategories(merchantId int) []*sale.Category {
	sl := s._rep.GetSale(merchantId)
	cats := sl.CategoryManager().GetCategories()
	var list []*sale.Category = make([]*sale.Category, len(cats))
	for i, v := range cats {
		vv := v.GetValue()
		vv.Icon = format.GetResUrl(vv.Icon)
		list[i] = vv
	}
	return list
}

// 根据上级编号获取分类列表
func (s *saleService) GetCategoriesByParentId(mchId, parentId int) []*sale.Category {
	cats := s.getCategoryManager(mchId).GetCategories()
	list := []*sale.Category{}
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == parentId && vv.Enabled == 1 {
			v2 := *vv
			v2.Icon = format.GetResUrl(v2.Icon)
			list = append(list, &v2)
		}
	}
	return list
}

func (s *saleService) getCategoryManager(mchId int) sale.ICategoryManager {
	if mchId > 0 {
		sl := s._rep.GetSale(mchId)
		return sl.CategoryManager()
	}
	return s._cateRep.GetGlobManager()
}

func (s *saleService) GetBigCategories(mchId int) []dto.Category {
	cats := s.getCategoryManager(mchId).GetCategories()
	list := []dto.Category{}
	for _, v := range cats {
		if v2 := v.GetValue(); v2.ParentId == 0 && v2.Enabled == 1 {
			v2.Icon = format.GetResUrl(v2.Icon)
			dv := dto.Category{}
			CopyCategory(v2, &dv)
			list = append(list, dv)
		}
	}
	return list
}

func (s *saleService) GetChildCategories(mchId, parentId int) []dto.Category {
	cats := s.getCategoryManager(mchId).GetCategories()
	list := []dto.Category{}
	for _, v := range cats {
		if vv := v.GetValue(); vv.ParentId == parentId && vv.Enabled == 1 {
			vv.Icon = format.GetResUrl(vv.Icon)
			dv := dto.Category{}
			CopyCategory(vv, &dv)
			s.setChild(cats, &dv)
			list = append(list, dv)
		}
	}
	return list
}

func CopyCategory(src *sale.Category, dst *dto.Category) {
	dst.Id = src.Id
	dst.Name = src.Name
	dst.Level = src.Level
	dst.Icon = src.Icon
	dst.Url = src.Url
}

func (s *saleService) setChild(list []sale.ICategory, dst *dto.Category) {
	for _, v := range list {
		if vv := v.GetValue(); vv.ParentId == dst.Id && vv.Enabled == 1 {
			if dst.Child == nil {
				dst.Child = []dto.Category{}
			}
			vv.Icon = format.GetResUrl(vv.Icon)
			dv := dto.Category{}
			CopyCategory(vv, &dv)
			dst.Child = append(dst.Child, dv)
		}
	}
}

func (s *saleService) GetAllSaleLabels(merchantId int) []*sale.Label {
	sl := s._rep.GetSale(merchantId)
	tags := sl.LabelManager().GetAllSaleLabels()

	lbs := make([]*sale.Label, len(tags))
	for i, v := range tags {
		lbs[i] = v.GetValue()
	}
	return lbs
}

// 获取销售标签
func (s *saleService) GetSaleLabel(merchantId, id int) *sale.Label {
	sl := s._rep.GetSale(merchantId)
	if tag := sl.LabelManager().GetSaleLabel(id); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 获取销售标签
func (s *saleService) GetSaleLabelByCode(merchantId int, code string) *sale.Label {
	sl := s._rep.GetSale(merchantId)
	if tag := sl.LabelManager().GetSaleLabelByCode(code); tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 保存销售标签
func (s *saleService) SaveSaleLabel(merchantId int, v *sale.Label) (int, error) {
	sl := s._rep.GetSale(merchantId)
	if v.Id > 0 {
		tag := sl.LabelManager().GetSaleLabel(v.Id)
		tag.SetValue(v)
		return tag.Save()
	}
	return sl.LabelManager().CreateSaleLabel(v).Save()
}

// 获取商品的销售标签
func (s *saleService) GetItemSaleLabels(merchantId, itemId int) []*sale.Label {
	var list = make([]*sale.Label, 0)
	sl := s._rep.GetSale(merchantId)
	if goods := sl.ItemManager().GetItem(itemId); goods != nil {
		list = goods.GetSaleLabels()
	}
	return list
}

// 保存商品的销售标签
func (s *saleService) SaveItemSaleLabels(merchantId, itemId int, tagIds []int) error {
	var err error
	sl := s._rep.GetSale(merchantId)
	if goods := sl.ItemManager().GetItem(itemId); goods != nil {
		err = goods.SaveSaleLabels(tagIds)
	} else {
		err = errors.New("商品不存在")
	}
	return err
}

// 根据销售标签获取指定数目的商品
func (s *saleService) GetValueGoodsBySaleLabel(merchantId int,
	code, sortBy string, begin int, end int) []*valueobject.Goods {
	sl := s._rep.GetSale(merchantId)
	if tag := sl.LabelManager().GetSaleLabelByCode(code); tag != nil {
		list := tag.GetValueGoods(sortBy, begin, end)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
		}
		return list
	}
	return make([]*valueobject.Goods, 0)
}

// 根据分页销售标签获取指定数目的商品
func (s *saleService) GetPagedValueGoodsBySaleLabel(merchantId int,
	tagId int, sortBy string, begin int, end int) (int, []*valueobject.Goods) {
	sl := s._rep.GetSale(merchantId)
	tag := sl.LabelManager().CreateSaleLabel(&sale.Label{
		Id: tagId,
	})
	return tag.GetPagedValueGoods(sortBy, begin, end)
}

// 删除销售标签
func (s *saleService) DeleteSaleLabel(merchantId int, id int) error {
	return s._rep.GetSale(merchantId).LabelManager().DeleteSaleLabel(id)
}

// 获取商品的会员价
func (s *saleService) GetGoodsLevelPrices(merchantId, goodsId int) []*goods.MemberPrice {
	sl := s._rep.GetSale(merchantId)
	if goods := sl.GoodsManager().GetGoods(goodsId); goods != nil {
		return goods.GetLevelPrices()
	}
	return make([]*goods.MemberPrice, 0)
}

// 保存商品的会员价
func (s *saleService) SaveMemberPrices(merchantId int, goodsId int,
	priceSet []*goods.MemberPrice) error {
	sl := s._rep.GetSale(merchantId)
	var err error
	if goods := sl.GoodsManager().GetGoods(goodsId); goods != nil {
		for _, v := range priceSet {
			if _, err = goods.SaveLevelPrice(v); err != nil {
				return err
			}
		}
	}
	return nil
}

//func (s *saleService) GetGoodsComplexInfo(goodsId int) *dto.GoodsComplex {
//	return s._goodsQuery.GetGoodsComplex(goodsId)
//}

// 获取商品详情
func (s *saleService) GetGoodsDetails(mchId, goodsId,
	mLevel int) (*valueobject.Goods, map[string]string) {
	sl := s._rep.GetSale(mchId)
	var goods sale.IGoods = sl.GoodsManager().GetGoods(goodsId)
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
func (s *saleService) GetItemDescriptionByGoodsId(merchantId, goodsId int) string {
	sl := s._rep.GetSale(merchantId)
	var goods sale.IGoods = sl.GoodsManager().GetGoods(goodsId)
	return goods.GetItem().GetValue().Description
}

// 获取商品快照
func (s *saleService) GetSnapshot(skuId int) *goods.Snapshot {
	return s._goodsRep.GetLatestSnapshot(skuId)
}

// 设置商品货架状态
func (s *saleService) SetShelveState(mchId int, itemId int, state int, remark string) error {
	sl := s._rep.GetSale(mchId)
	gi := sl.ItemManager().GetItem(itemId)
	if gi == nil {
		return goods.ErrNoSuchGoods
	}
	return gi.SetShelve(state, remark)
}

// 设置商品货架状态
func (s *saleService) ReviewItem(mchId int, itemId int, pass bool, remark string) error {
	sl := s._rep.GetSale(mchId)
	gi := sl.ItemManager().GetItem(itemId)
	if gi == nil {
		return goods.ErrNoSuchGoods
	}
	return gi.Review(pass, remark)
}

// 标记为违规
func (s *saleService) SignIncorrect(supplierId int, itemId int, remark string) error {
	sl := s._rep.GetSale(supplierId)
	gi := sl.ItemManager().GetItem(itemId)
	if gi == nil {
		return goods.ErrNoSuchGoods
	}
	return gi.Incorrect(remark)
}
