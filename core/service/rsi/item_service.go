/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package rsi

import (
	"context"
	"fmt"
	"github.com/ixre/gof/crypto"
	"github.com/ixre/gof/storage"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/core/service/auto_gen/rpc/item_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/parser"
	"strconv"
)

var _ item_service.ItemService = new(itemService)

type itemService struct {
	*serviceUtil
	itemRepo  item.IGoodsItemRepo
	itemQuery *query.ItemQuery
	cateRepo  product.ICategoryRepo
	labelRepo item.ISaleLabelRepo
	promRepo  promodel.IProModelRepo
	mchRepo   merchant.IMerchantRepo
	valueRepo valueobject.IValueRepo
	sto       storage.Interface
}

func NewSaleService(sto storage.Interface, cateRepo product.ICategoryRepo,
	goodsRepo item.IGoodsItemRepo, goodsQuery *query.ItemQuery,
	labelRepo item.ISaleLabelRepo, promRepo promodel.IProModelRepo,
	mchRepo merchant.IMerchantRepo, valueRepo valueobject.IValueRepo) *itemService {
	return &itemService{
		sto:       sto,
		itemRepo:  goodsRepo,
		itemQuery: goodsQuery,
		cateRepo:  cateRepo,
		labelRepo: labelRepo,
		promRepo:  promRepo,
		mchRepo:   mchRepo,
		valueRepo: valueRepo,
	}
}

// 获取商品值
func (s *itemService) GetItemValue(itemId int64) *ttype.SOldItem {
	item := s.itemRepo.GetItem(itemId)
	if item != nil {
		return parser.ItemDto(item.GetValue())
	}
	return nil
}

// 获取SKU
func (s *itemService) GetSku(ctx context.Context, itemId, skuId int64) (r *ttype.SSku, err error) {
	item := s.itemRepo.GetItem(itemId)
	if item != nil {
		sku := item.GetSku(skuId)
		if sku != nil {
			return parser.SkuDto(sku), nil
		}
	}
	return nil, nil
}

// 获取SKU数组
func (s *itemService) GetSkuArray(itemId int64) []*item.Sku {
	it := s.itemRepo.GetItem(itemId)
	if it != nil {
		return it.SkuArray()
	}
	return []*item.Sku{}
}

// 获取商品规格HTML信息
func (s *itemService) GetSkuHtmOfItem(itemId int64) (specHtm string) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: itemId})
	specArr := it.SpecArray()
	return s.itemRepo.SkuService().GetSpecHtml(specArr)
}

// 获取商品详细数据
func (s *itemService) GetItemDetailData(ctx context.Context, itemId int64, iType int32) (r string, err error) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: itemId})
	switch iType {
	case item.ItemWholesale:
		data := it.Wholesale().GetJsonDetailData()
		return string(data), nil
	}
	return "不支持的商品类型", nil
}

// 获取商品的Sku-JSON格式
func (s *itemService) GetItemSkuJson(ctx context.Context, itemId int64) (r string, err error) {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: itemId})
	skuBytes := s.itemRepo.SkuService().GetSkuJson(it.SkuArray())
	return string(skuBytes), nil
}

// 保存商品
func (s *itemService) SaveItem(di *ttype.SOldItem, vendorId int32) (_ *ttype.Result_, err error) {
	var gi item.IGoodsItem
	it := parser.Item(di)
	if it.ID > 0 {
		gi = s.itemRepo.GetItem(it.ID)
		if gi == nil || gi.GetValue().VendorId != vendorId {
			err = item.ErrNoSuchItem
			goto R
		}
	} else {
		gi = s.itemRepo.CreateItem(it)
	}
	err = gi.SetValue(it)
	if err == nil {
		err = gi.SetSku(it.SkuArray)
		if err == nil {
			it.ID, err = gi.Save()
		}
	}
R:
	r := s.result(err)
	r.Data = map[string]string{
		"ItemId": strconv.Itoa(int(it.ID)),
	}
	return r, nil
}

// 获取上架商品数据（分页）
func (s *itemService) GetPagedOnShelvesItem(itemType int32, catId int32, start,
	end int32, where, sortBy string) (int32, []*ttype.SOldItem) {
	switch itemType {
	case item.ItemNormal:
		return s.getPagedOnShelvesItem(catId, start, end, where, sortBy)
	case item.ItemWholesale:
		return s.getPagedOnShelvesItemForWholesale(catId, start, end, where, sortBy)
	}
	return 0, []*ttype.SOldItem{}
}
func (s *itemService) getPagedOnShelvesItem(catId int32, start,
	end int32, where, sortBy string) (int32, []*ttype.SOldItem) {

	total, list := s.itemQuery.GetPagedOnShelvesItem(catId,
		start, end, where, sortBy)
	arr := make([]*ttype.SOldItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ItemDto(v)
	}
	return total, arr
}

func (s *itemService) getPagedOnShelvesItemForWholesale(catId int32, start,
	end int32, where, sortBy string) (int32, []*ttype.SOldItem) {

	total, list := s.itemQuery.GetPagedOnShelvesItemForWholesale(catId,
		start, end, where, sortBy)
	arr := make([]*ttype.SOldItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		dto := parser.ItemDto(v)
		s.attachWholesaleItemData(dto)
		arr[i] = dto
	}
	return total, arr
}

// 获取上架商品数据（分页）
func (s *itemService) SearchOnShelvesItem(itemType int32, word string, start,
	end int32, where, sortBy string) (int32, []*ttype.SOldItem) {

	switch itemType {
	case item.ItemNormal:
		return s.searchOnShelveItem(word, start, end, where, sortBy)
	case item.ItemWholesale:
		return s.searchOnShelveItemForWholesale(word, start, end, where, sortBy)
	}
	return 0, []*ttype.SOldItem{}
}

func (s itemService) searchOnShelveItem(word string, start,
	end int32, where, sortBy string) (int32, []*ttype.SOldItem) {
	total, list := s.itemQuery.SearchOnShelvesItem(word,
		start, end, where, sortBy)
	arr := make([]*ttype.SOldItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		arr[i] = parser.ItemDto(v)
	}
	return total, arr
}

func (s itemService) searchOnShelveItemForWholesale(word string, start,
	end int32, where, sortBy string) (int32, []*ttype.SOldItem) {
	total, list := s.itemQuery.SearchOnShelvesItemForWholesale(word,
		start, end, where, sortBy)
	arr := make([]*ttype.SOldItem, len(list))
	for i, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
		dto := parser.ItemDto(v)
		s.attachWholesaleItemData(dto)
		arr[i] = dto

	}
	return total, arr
}

// 附加批发商品的信息
func (s *itemService) attachWholesaleItemData(dto *ttype.SOldItem) {
	dto.Data = make(map[string]string)
	vendor := s.mchRepo.GetMerchant(dto.VendorId)
	if vendor != nil {
		vv := vendor.GetValue()
		pStr := s.valueRepo.GetAreaName(vv.Province)
		cStr := s.valueRepo.GetAreaName(vv.City)
		dto.Data["VendorName"] = vv.CompanyName
		dto.Data["ShipArea"] = pStr + cStr
		// 认证信息
		ei := vendor.ProfileManager().GetEnterpriseInfo()
		if ei != nil && ei.Reviewed == enum.ReviewPass {
			dto.Data["Authorized"] = "true"
		} else {
			dto.Data["Authorized"] = "false"
		}
		// 品牌
		b := s.promRepo.BrandService().Get(dto.BrandId)
		if b != nil {
			dto.Data["BrandName"] = b.Name
			dto.Data["BrandImage"] = b.Image
			dto.Data["BrandId"] = strconv.Itoa(int(b.ID))
		}
	}
}

// 获取上架商品数据（分页）
func (s *itemService) GetRandomItem(catId int32, quantity int32, where string) []*ttype.SOldItem {
	hash := fmt.Sprintf("%d-%d-%s", catId, quantity, where)
	hash = crypto.Md5([]byte(hash))
	key := "go2o:query:cache:rd-item:" + hash
	var arr []*ttype.SOldItem

	fn := func() interface{} {
		list := s.itemQuery.GetRandomItem(catId, quantity, where)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
			arr = append(arr, parser.ItemDto(v))
		}
		return arr
	}
	s.sto.RWJson(key, &arr, fn, 600)
	return arr
}

// 获取上架商品数据（分页）
func (s *itemService) GetBigCatItems(catId, quantity int32, where string) []*ttype.SOldItem {
	c := s.cateRepo.GlobCatService().GetCategory(catId)
	if c != nil {
		ids := c.GetChildes()
		list := s.itemQuery.GetOnShelvesItem(ids, 0, quantity, where)
		arr := make([]*ttype.SOldItem, len(list))
		for i, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
			arr[i] = parser.ItemDto(v)
		}
		return arr
	}
	return []*ttype.SOldItem{}
}

// 根据SKU获取商品
func (s *itemService) GetGoodsBySku(mchId int32, itemId, sku int64) *valueobject.Goods {
	v := s.itemRepo.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		return s.itemRepo.CreateItem(v).GetPackedValue()
	}
	return nil
}

// 根据SKU获取商品
func (s *itemService) GetValueGoodsBySku(mchId int32, itemId, sku int64) *item.GoodsItem {
	v := s.itemRepo.GetValueGoodsBySku(itemId, sku)
	if v != nil {
		return s.itemRepo.CreateItem(v).GetValue()
	}
	return nil
}

// 根据快照编号获取商品
func (s *itemService) GetGoodsBySnapshotId(snapshotId int64) *item.GoodsItem {
	snap := s.itemRepo.GetSalesSnapshot(snapshotId)
	if snap != nil {
		return s.itemRepo.GetValueGoodsById(snap.SkuId)
	}
	return nil
}

// 根据快照编号获取商品
func (s *itemService) GetSaleSnapshotById(snapshotId int64) *item.TradeSnapshot {
	return s.itemRepo.GetSalesSnapshot(snapshotId)
}

// 获取分页上架的商品
func (s *itemService) GetShopPagedOnShelvesGoods(shopId, categoryId int32, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cat := s.cateRepo.GlobCatService().GetCategory(categoryId)
		ids := cat.GetChildes()
		ids = append(ids, categoryId)
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId, ids, start, end, "", sortBy)
	} else {
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId, nil, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取上架商品数据（分页）
func (s *itemService) GetPagedOnShelvesGoods__(shopId int32, categoryId int32, start, end int,
	sortBy string) (total int, list []*valueobject.Goods) {
	if categoryId > 0 {
		cate := s.cateRepo.GlobCatService().GetCategory(categoryId)
		var ids = cate.GetChildes()
		ids = append(ids, categoryId)
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId,
			ids, start, end, "", sortBy)
	} else {
		total, list = s.itemRepo.GetPagedOnShelvesGoods(shopId,
			[]int32{}, start, end, "", sortBy)
	}
	for _, v := range list {
		v.Image = format.GetGoodsImageUrl(v.Image)
	}
	return total, list
}

// 获取分页上架的商品
func (s *itemService) GetPagedOnShelvesGoodsByKeyword(shopId int32, start, end int,
	word, sortQuery string) (int, []*valueobject.Goods) {
	var where string
	var orderBy string
	switch sortQuery {
	case "price_0":
		where = ""
		orderBy = "it.price ASC"
	case "price_1":
		where = ""
		orderBy = "it.price DESC"
	case "sale_0":
		where = ""
		orderBy = "it.sale_num ASC"
	case "sale_1":
		where = ""
		orderBy = "it.sale_num DESC"
	case "rate_0":
		//todo:
	case "rate_1":
		//todo:
	}
	return s.itemQuery.GetPagedOnShelvesGoodsByKeyword(shopId,
		start, end, word, where, orderBy)
}

// 删除产品
func (s *itemService) DeleteGoods(mchId int32, itemId int64) error {
	gi := s.itemRepo.GetItem(itemId)
	if gi == nil || gi.GetValue().VendorId != mchId {
		return item.ErrNoSuchItem
	}
	return gi.Destroy()
}

func (s *itemService) GetAllSaleLabels() []*item.Label {
	tags := s.labelRepo.LabelService().GetAllSaleLabels()
	lbs := make([]*item.Label, len(tags))
	for i, v := range tags {
		lbs[i] = v.GetValue()
	}
	return lbs
}

// 获取销售标签
func (s *itemService) GetSaleLabel(labelId int32) *item.Label {
	tag := s.labelRepo.LabelService().GetSaleLabel(labelId)
	if tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 获取销售标签
func (s *itemService) GetSaleLabelByCode(code string) *item.Label {
	tag := s.labelRepo.LabelService().GetSaleLabelByCode(code)
	if tag != nil {
		return tag.GetValue()
	}
	return nil
}

// 保存销售标签
func (s *itemService) SaveSaleLabel(v *item.Label) (int32, error) {
	ls := s.labelRepo.LabelService()
	var lb item.ISaleLabel

	if v.Id > 0 {
		lb = ls.GetSaleLabel(v.Id)
		if lb == nil {
			panic("没有销售标签")
		}
	} else {
		lb = ls.CreateSaleLabel(v)
	}
	err := lb.SetValue(v)
	if err == nil {
		return lb.Save()
	}
	return v.Id, err
}

// 根据销售标签获取指定数目的商品
func (s *itemService) GetValueGoodsBySaleLabel(code, sortBy string, begin int, end int) []*valueobject.Goods {
	tag := s.labelRepo.LabelService().GetSaleLabelByCode(code)
	if tag != nil {
		list := tag.GetValueGoods(sortBy, begin, end)
		for _, v := range list {
			v.Image = format.GetGoodsImageUrl(v.Image)
		}
		return list
	}
	return make([]*valueobject.Goods, 0)
}

// 根据分页销售标签获取指定数目的商品
func (s *itemService) GetPagedValueGoodsBySaleLabel(shopId int32, tagId int32, sortBy string, begin int, end int) (int,
	[]*valueobject.Goods) {
	tag := s.labelRepo.LabelService().CreateSaleLabel(&item.Label{
		Id: tagId,
	})
	return tag.GetPagedValueGoods(sortBy, begin, end)
}

// 删除销售标签
func (s *itemService) DeleteSaleLabel(id int32) error {
	return s.labelRepo.LabelService().DeleteSaleLabel(id)
}

// 获取商品的会员价
func (s *itemService) GetGoodsLevelPrices(itemId int64) []*item.MemberPrice {
	gi := s.itemRepo.GetItem(itemId)
	if gi != nil {
		return gi.GetLevelPrices()
	}
	return make([]*item.MemberPrice, 0)
}

// 保存商品的会员价
func (s *itemService) SaveMemberPrices(mchId int32, itemId int64,
	priceSet []*item.MemberPrice) (err error) {
	gi := s.itemRepo.GetItem(itemId)
	if gi != nil {
		for _, v := range priceSet {
			if _, err = gi.SaveLevelPrice(v); err != nil {
				return err
			}
		}
	}
	return err
}

//func (s *saleService) GetGoodsComplexInfo(goodsId int32) *dto.GoodsComplex {
//	return s._goodsQuery.GetGoodsComplex(goodsId)
//}

// 获取商品详情
func (s *itemService) GetGoodsDetails(itemId int64, mLevel int32) (
	*valueobject.Goods, map[string]string) {
	goods := s.itemRepo.GetItem(itemId)
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
func (s *itemService) GetItemDescriptionByGoodsId(itemId int64) string {
	it := s.itemRepo.CreateItem(&item.GoodsItem{ID: itemId})
	pro := it.Product()
	if pro != nil {
		return pro.GetValue().Description
	}
	return ""
}

// 获取商品快照
func (s *itemService) GetSnapshot(skuId int64) *item.Snapshot {
	return s.itemRepo.GetLatestSnapshot(skuId)
}

// 设置商品货架状态
func (s *itemService) SetShelveState(vendorId int32, itemId int64,
	itemType int32, state int32, remark string) (_ *ttype.Result_, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchItem
	} else {
		switch itemType {
		case item.ItemWholesale:
			err = it.Wholesale().SetShelve(state, remark)
		default:
			err = it.SetShelve(state, remark)
		}
	}
	return s.result(err), nil
}

// 设置商品货架状态
func (s *itemService) ReviewItem(vendorId int32, itemId int64,
	pass bool, remark string) (_ *ttype.Result_, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchItem
	} else {
		err = it.Review(pass, remark)
	}
	return s.result(err), nil
}

// 标记为违规
func (s *itemService) SignIncorrect(vendorId int32, itemId int64,
	remark string) (_ *ttype.Result_, err error) {
	it := s.itemRepo.GetItem(itemId)
	if it == nil || it.GetValue().VendorId != vendorId {
		err = item.ErrNoSuchItem
	} else {
		err = it.Incorrect(remark)
	}
	return s.result(err), nil
}

// 获取批发价格数组
func (s *itemService) GetWholesalePriceArray(itemId, skuId int64) []*item.WsSkuPrice {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().GetSkuPrice(skuId)
}

// 保存批发价格
func (s *itemService) SaveWholesalePrice(itemId, skuId int64, arr []*item.WsSkuPrice) error {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().SaveSkuPrice(skuId, arr)
}

// 获取批发折扣数组
func (s *itemService) GetWholesaleDiscountArray(itemId int64, groupId int32) []*item.WsItemDiscount {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().GetItemDiscount(groupId)
}

// 保存批发折扣
func (s *itemService) SaveWholesaleDiscount(itemId int64, groupId int32, arr []*item.WsItemDiscount) error {
	it := s.itemRepo.GetItem(itemId)
	return it.Wholesale().SaveItemDiscount(groupId, arr)
}
