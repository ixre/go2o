/**
 * Copyright 2015 @ z3q.net.
 * name : sale_goods
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package item

import (
	"errors"
	"fmt"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"strconv"
	"strings"
	"time"
)

var _ item.IGoodsItem = new(goodsItemImpl)

// 临时的商品实现  todo: 要与item分开
type goodsItemImpl struct {
	pro           product.IProduct
	value         *item.GoodsItem
	goodsRepo     item.IGoodsItemRepo
	productRepo   product.IProductRepo
	proMRepo      promodel.IProModelRepo
	promRepo      promotion.IPromotionRepo
	levelPrices   []*item.MemberPrice
	promDescribes map[string]string
	snapManager   item.ISnapshotManager
	valRepo       valueobject.IValueRepo
	expressRepo   express.IExpressRepo
}

//todo:??? 去掉依赖promotion.IPromotionRepo

func NewSaleItem(
	itemRepo product.IProductRepo, pro product.IProduct,
	value *item.GoodsItem, valRepo valueobject.IValueRepo,
	goodsRepo item.IGoodsItemRepo, proMRepo promodel.IProModelRepo,
	expressRepo express.IExpressRepo,
	promRepo promotion.IPromotionRepo) item.IGoodsItem {
	v := &goodsItemImpl{
		pro:         pro,
		value:       value,
		productRepo: itemRepo,
		goodsRepo:   goodsRepo,
		proMRepo:    proMRepo,
		promRepo:    promRepo,
		valRepo:     valRepo,
		expressRepo: expressRepo,
	}
	return v.init()
}

func (g *goodsItemImpl) init() item.IGoodsItem {
	if g.pro != nil {
		g.value.PromPrice = g.value.Price
	}
	return g
}

//获取聚合根编号
func (g *goodsItemImpl) GetAggregateRootId() int32 {
	return g.value.Id
}

// 商品快照
func (g *goodsItemImpl) SnapshotManager() item.ISnapshotManager {
	if g.snapManager == nil {
		g.snapManager = NewSnapshotManagerImpl(g.GetAggregateRootId(),
			g.goodsRepo, g.GetValue())
	}
	return g.snapManager
}

// 获取货品
func (g *goodsItemImpl) Product() product.IProduct {
	return g.pro
}

// 设置值
func (g *goodsItemImpl) GetValue() *item.GoodsItem {
	return g.value
}

// 获取包装过的商品信息
func (g *goodsItemImpl) GetPackedValue() *valueobject.Goods {
	//item := g.GetItem().GetValue()
	gv := g.GetValue()
	goods := &valueobject.Goods{
		ProductId:     gv.ProductId,
		CategoryId:    gv.CatId,
		Name:          gv.Title,
		GoodsNo:       gv.Code,
		Image:         gv.Image,
		Price:         gv.RetailPrice,
		SalePrice:     gv.Price,
		PromPrice:     gv.Price,
		GoodsId:       g.GetAggregateRootId(),
		SkuId:         gv.SkuId,
		IsPresent:     gv.IsPresent,
		PromotionFlag: gv.PromFlag,
		StockNum:      gv.StockNum,
		SaleNum:       gv.SaleNum,
	}
	return goods
}

// 设置值
func (g *goodsItemImpl) SetValue(v *item.GoodsItem) error {
	err := g.checkItemValue(v)
	if err == nil {
		err = g.copyFromProduct(v)
		if err == nil {
			g.value.ShopId = v.ShopId
			g.value.IsPresent = v.IsPresent
			//g.value.Title = v.Title
			g.value.ShopCatId = v.ShopCatId
			g.value.IsPresent = v.IsPresent
			g.value.ProductId = v.ProductId
			g.value.PromFlag = v.PromFlag
			g.value.ShopCatId = v.ShopId
			g.value.ExpressTid = v.ExpressTid
			g.value.ShortTitle = v.ShortTitle
			g.value.Code = v.Code
			g.value.SaleNum = v.SaleNum
			g.value.StockNum = v.StockNum
			g.value.Cost = v.Cost
			g.value.RetailPrice = v.RetailPrice
			g.value.Price = v.Price
			g.value.Weight = v.Weight
			g.value.Bulk = v.Bulk
			if g.value.CreateTime == 0 {
				g.value.CreateTime = time.Now().Unix()
			}
			//修改图片或标题后，要重新审核
			if g.value.Image != v.Image || g.value.Title != v.Title {
				g.resetReview()
			}
		}
	}
	return err
}

// 从产品中拷贝信息
//todo: 如后期弄成公共产品，则应保持产品与商品的数据独立。
func (g *goodsItemImpl) copyFromProduct(v *item.GoodsItem) error {
	pro := g.productRepo.GetProductValue(v.ProductId)
	if pro == nil {
		return product.ErrNoSuchProduct
	}
	g.value.CatId = pro.CatId
	g.value.VendorId = pro.VendorId
	g.value.BrandId = pro.BrandId
	g.value.Title = pro.Name
	g.value.Code = pro.Code
	g.value.Image = pro.Image
	g.value.SortNum = pro.SortNum
	g.value.CreateTime = pro.CreateTime
	g.value.UpdateTime = pro.UpdateTime
	return nil
}

// 重置审核状态
func (i *goodsItemImpl) resetReview() {
	i.value.ReviewState = enum.ReviewAwaiting
}

// 检查商品数据是否正确
func (g *goodsItemImpl) checkItemValue(v *item.GoodsItem) error {
	registry := g.valRepo.GetRegistry()
	// 检测是否上传图片
	if v.Image == registry.GoodsDefaultImage {
		return product.ErrNotUploadImage
	}
	if v.ShopId <= 0 {
		return item.ErrNotBindShop
	} else {
		//todo: 判断商铺是否为本商家的
		// return item.ErrIncorrectShopOfItem
	}

	// 检测运费模板
	if v.ExpressTid > 0 {
		ve := g.expressRepo.GetUserExpress(v.VendorId)
		tpl := ve.GetTemplate(v.ExpressTid)
		if tpl == nil {
			return express.ErrNoSuchTemplate
		}
		if !tpl.Enabled() {
			return express.ErrTemplateNotEnabled
		}
	} else {
		return shipment.ErrNotSetExpressTemplate
	}
	// 检测价格
	return g.checkPrice(v)
}

// 判断价格是否正确
func (i *goodsItemImpl) checkPrice(v *item.GoodsItem) error {
	rate := (v.Price - v.Cost) / v.Price
	conf := i.valRepo.GetRegistry()
	minRate := conf.GoodsMinProfitRate
	// 如果未设定最低利润率，则可以与供货价一致
	if minRate != 0 && rate < minRate {
		return errors.New(fmt.Sprintf(item.ErrGoodsMinProfitRate.Error(),
			strconv.Itoa(int(minRate*100))+"%"))
	}
	return nil
}

// 设置SKU
func (g *goodsItemImpl) SetSku(arr []*item.Sku) error {
	g.value.SkuArray = arr
	return nil
}

// 保存
func (g *goodsItemImpl) Save() (_ int32, err error) {
	// 创建商品
	if g.GetAggregateRootId() <= 0 {
		g.value.Id, err = g.goodsRepo.SaveValueGoods(g.value)
		if err != nil {
			return g.value.Id, err
		}
	}
	// 保存商品SKU
	if g.value.SkuArray != nil {
		g.saveItemSku(g.value.SkuArray)
		g.value.SkuNum = int32(len(g.value.SkuArray))
	}
	// 保存商品
	g.value.Id, err = g.goodsRepo.SaveValueGoods(g.value)
	if err == nil {
		// 保存商品快照
		_, err = g.SnapshotManager().GenerateSnapshot()
	}
	return g.value.Id, err
}

// ========== [# SKU处理开始 ]  ===========//

// 保存商品SKU
func (g *goodsItemImpl) saveItemSku(arr []*item.Sku) (err error) {
	pk := g.GetAggregateRootId()
	ss := g.goodsRepo.SkuService()
	// 格式化数据
	err = ss.RebuildSkuArray(&arr, g.value)
	if err == nil {
		// 获取之前的SKU设置
		old := g.goodsRepo.SelectItemSku("item_id=?", pk)
		// 合并SKU
		ss.Merge(old, &arr)
		// 分析当前项目并加入到MAP中
		delList := []int32{}
		currMap := make(map[int32]*item.Sku, len(arr))
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
			g.goodsRepo.DeleteItemSku(v)
		}
		// 保存项
		for _, v := range arr {
			if v.ItemId == 0 {
				v.ItemId = pk
			}
			if proId := g.value.ProductId; v.ProductId != proId {
				v.ProductId = proId
			}
			if v.ItemId == pk {
				v.Id, err = util.I32Err(g.goodsRepo.SaveItemSku(v))
			}
		}
	}
	return err
}

// 获取SKU数组
func (g *goodsItemImpl) SkuArray() []*item.Sku {
	if g.value.SkuArray == nil {
		g.value.SkuArray = g.goodsRepo.SelectItemSku("item_id=?",
			g.GetAggregateRootId())
	}
	return g.value.SkuArray
}

// ========== [/ SKU处理结束 ] ===========//

// 获取促销信息
func (g *goodsItemImpl) GetPromotions() []promotion.IPromotion {
	var vp []*promotion.PromotionInfo = g.promRepo.GetPromotionOfGoods(
		g.GetAggregateRootId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = g.promRepo.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (g *goodsItemImpl) GetLevelPrice(level int32) (bool, float32) {
	lvp := g.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < g.value.Price {
			return true, v.Price
		}
	}
	return false, g.value.Price
}

// 获取促销价
func (g *goodsItemImpl) GetPromotionPrice(level int32) float32 {
	b, price := g.GetLevelPrice(level)
	if b {
		return price
	}
	return g.value.Price
}

// 获取促销描述
func (g *goodsItemImpl) GetPromotionDescribe() map[string]string {
	if g.promDescribes == nil {
		proms := g.GetPromotions()
		g.promDescribes = make(map[string]string, len(proms))
		for _, v := range proms {
			key := v.TypeName()
			if txt, ok := g.promDescribes[key]; !ok {
				g.promDescribes[key] = v.GetValue().ShortName
			} else {
				g.promDescribes[key] = txt + ";" + v.GetValue().ShortName
			}

			//			if v.Type() == promotion.TypeFlagCashBack {
			//				if txt, ok := g._promDescribes[key]; !ok {
			//					g._promDescribes[key] = v.GetValue().ShortName
			//				} else {
			//					g._promDescribes[key] = txt + ";" + v.GetValue().ShortName
			//				}
			//			} else if v.Type() == promotion.TypeFlagCoupon {
			//				if txt, ok := g._promDescribes[key]; !ok {
			//					g._promDescribes[key] = v.GetValue().ShortName
			//				} else {
			//					g._promDescribes[key] = txt + ";" + v.GetValue().ShortName
			//				}
			//			}

			//todo: other promotion implement
		}
	}
	return g.promDescribes
}

// 获取会员价
func (g *goodsItemImpl) GetLevelPrices() []*item.MemberPrice {
	if g.levelPrices == nil {
		g.levelPrices = g.goodsRepo.GetGoodsLevelPrice(g.GetAggregateRootId())
	}
	return g.levelPrices
}

// 保存会员价
func (g *goodsItemImpl) SaveLevelPrice(v *item.MemberPrice) (int32, error) {
	v.GoodsId = g.GetAggregateRootId()
	if g.value.Price == v.Price {
		if v.Id > 0 {
			g.goodsRepo.RemoveGoodsLevelPrice(v.Id)
		}
		return -1, nil
	}
	return g.goodsRepo.SaveGoodsLevelPrice(v)
}

// 是否上架
func (g *goodsItemImpl) IsOnShelves() bool {
	return g.value.ShelveState == item.ShelvesOn
}

// 设置上架
func (g *goodsItemImpl) SetShelve(state int32, remark string) error {
	if state == item.ShelvesIncorrect && len(remark) == 0 {
		return product.ErrNilRejectRemark
	}
	g.value.ShelveState = state
	g.value.ReviewRemark = remark
	_, err := g.Save()
	return err
}

// 标记为违规
func (g *goodsItemImpl) Incorrect(remark string) error {
	g.value.ShelveState = item.ShelvesIncorrect
	g.value.ReviewRemark = remark
	_, err := g.Save()
	return err
}

// 审核
func (g *goodsItemImpl) Review(pass bool, remark string) error {
	if pass {
		g.value.ReviewState = enum.ReviewPass

	} else {
		remark = strings.TrimSpace(remark)
		if remark == "" {
			return item.ErrEmptyReviewRemark
		}
		g.value.ReviewState = enum.ReviewReject
	}
	g.value.ReviewRemark = remark
	_, err := g.Save()
	return err
}

// 更新销售数量
func (g *goodsItemImpl) AddSalesNum(quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return item.ErrOutOfStock
	}
	g.value.SaleNum += quantity
	_, err := g.Save()
	return err
}

// 取消销售
func (g *goodsItemImpl) CancelSale(quantity int32, orderNo string) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	g.value.SaleNum -= quantity
	_, err := g.Save()
	return err
}

// 占用库存
func (g *goodsItemImpl) TakeStock(quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return item.ErrOutOfStock
	}
	g.value.StockNum -= quantity
	_, err := g.Save()
	return err
}

// 释放库存
func (g *goodsItemImpl) FreeStock(quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	g.value.StockNum += quantity
	_, err := g.Save()
	return err
}

// 删除商品
func (g *goodsItemImpl) Destroy() error {
	//g.goodsRepo.
	return nil
}
