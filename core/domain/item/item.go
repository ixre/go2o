/**
 * Copyright 2015 @ to2.net.
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
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/express"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/pro_model"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/promotion"
	"go2o/core/domain/interface/registry"
	"go2o/core/domain/interface/shipment"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/format"
	"strconv"
	"strings"
	"time"
)

var _ item.IGoodsItem = new(itemImpl)

// 商品实现
type itemImpl struct {
	pro           product.IProduct
	value         *item.GoodsItem
	wholesale     item.IWholesaleItem
	snapshot      *item.Snapshot
	repo          item.IGoodsItemRepo
	catRepo       product.ICategoryRepo
	productRepo   product.IProductRepo
	itemWsRepo    item.IItemWholesaleRepo
	proMRepo      promodel.IProModelRepo
	promRepo      promotion.IPromotionRepo
	levelPrices   []*item.MemberPrice
	promDescribes map[string]string
	registryRepo  registry.IRegistryRepo
	expressRepo   express.IExpressRepo
}

//todo:??? 去掉依赖promotion.IPromotionRepo

func NewItem(
	itemRepo product.IProductRepo, catRepo product.ICategoryRepo,
	pro product.IProduct, value *item.GoodsItem, registryRepo registry.IRegistryRepo,
	goodsRepo item.IGoodsItemRepo, proMRepo promodel.IProModelRepo,
	itemWsRepo item.IItemWholesaleRepo, expressRepo express.IExpressRepo,
	promRepo promotion.IPromotionRepo) item.IGoodsItem {
	v := &itemImpl{
		pro:          pro,
		value:        value,
		catRepo:      catRepo,
		productRepo:  itemRepo,
		repo:         goodsRepo,
		proMRepo:     proMRepo,
		itemWsRepo:   itemWsRepo,
		promRepo:     promRepo,
		registryRepo: registryRepo,
		expressRepo:  expressRepo,
	}
	return v.init()
}

func (g *itemImpl) init() item.IGoodsItem {
	if g.pro != nil {
		g.value.PromPrice = g.value.Price
	}
	return g
}

//获取聚合根编号
func (g *itemImpl) GetAggregateRootId() int64 {
	return g.value.ID
}

// 商品快照
func (g *itemImpl) Snapshot() *item.Snapshot {
	if g.snapshot == nil {
		g.snapshot = g.repo.SnapshotService().GetLatestSnapshot(
			g.GetAggregateRootId())
	}
	return g.snapshot
}

// 获取货品
func (g *itemImpl) Product() product.IProduct {
	if g.pro == nil && g.value.ProductId > 0 {
		g.pro = g.productRepo.GetProduct(g.value.ProductId)
	}
	return g.pro
}

// 批发
func (i *itemImpl) Wholesale() item.IWholesaleItem {
	if i.wholesale == nil {
		i.wholesale = newWholesaleItem(i.GetAggregateRootId(),
			i, i.repo, i.itemWsRepo)
	}
	return i.wholesale
}

// 设置值
func (g *itemImpl) GetValue() *item.GoodsItem {
	return g.value
}

// 获取包装过的商品信息
func (g *itemImpl) GetPackedValue() *valueobject.Goods {
	//item := g.GetItem().Value()
	gv := g.GetValue()
	goods := &valueobject.Goods{
		ProductId:     gv.ProductId,
		CategoryId:    gv.CatId,
		Title:         gv.Title,
		GoodsNo:       gv.Code,
		Image:         gv.Image,
		RetailPrice:   gv.RetailPrice,
		Price:         gv.Price,
		PriceRange:    gv.PriceRange,
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
func (g *itemImpl) SetValue(v *item.GoodsItem) error {
	err := g.checkItemValue(v)
	if err == nil {
		err = g.copyFromProduct(v)
		if err == nil {
			// 创建商品时，设为已下架
			if g.GetAggregateRootId() <= 0 {
				g.value.ShelveState = item.ShelvesDown
				// 分类在创建后，不允许再进行修改。并且分类不能为虚拟分类
				// 如果修改，则所有SKU和属性应删除。
				c := g.catRepo.GlobCatService().GetCategory(v.CatId)
				if c == nil {
					return item.ErrIncorrectProductCategory
				}
				cv := c.GetValue()
				if cv.VirtualCat == 1 {
					return item.ErrIncorrectProductCategory
				}
				g.value.CatId = v.CatId
			}
			g.value.ShopId = v.ShopId
			g.value.IsPresent = v.IsPresent
			g.value.ShopCatId = v.ShopCatId
			g.value.IsPresent = v.IsPresent
			g.value.ProductId = v.ProductId
			g.value.PromFlag = v.PromFlag
			g.value.ShopCatId = v.ShopId
			g.value.ExpressTid = v.ExpressTid
			g.value.Title = v.Title
			g.value.ShortTitle = v.ShortTitle
			g.value.Code = v.Code
			g.value.SaleNum = v.SaleNum
			g.value.StockNum = v.StockNum
			g.value.Cost = v.Cost
			g.value.RetailPrice = v.RetailPrice
			g.value.Price = v.Price
			g.value.Weight = v.Weight
			g.value.Bulk = v.Bulk
			//设置默认的价格区间
			if g.value.PriceRange == "0" || g.value.PriceRange == "" {
				g.value.PriceRange = format.FormatFloat(v.Price)
			}
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

// 设置SKU
func (g *itemImpl) SetSku(arr []*item.Sku) error {
	g.value.SkuArray = arr
	return nil
}

// ========== [# SKU处理开始 ]  ===========//

// 保存商品SKU
func (g *itemImpl) saveItemSku(arrPtr *[]*item.Sku) (err error) {
	arr := *arrPtr
	pk := g.GetAggregateRootId()
	ss := g.repo.SkuService()
	// 格式化数据
	err = ss.RebuildSkuArray(&arr, g.value)
	if err == nil {
		// 获取之前的SKU设置
		old := g.repo.SelectItemSku("item_id= $1", pk)
		// 合并SKU
		ss.Merge(old, &arr)
		// 分析当前项目并加入到MAP中
		delList := []int64{}
		currMap := make(map[int64]*item.Sku, len(arr))
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
			g.repo.DeleteItemSku(v)
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
				v.ID, err = util.I64Err(g.repo.SaveItemSku(v))
			}
		}
	}
	return err
}

// 获取SKU数组
func (g *itemImpl) SkuArray() []*item.Sku {
	if g.value.SkuArray == nil {
		g.value.SkuArray = g.repo.SelectItemSku("item_id= $1",
			g.GetAggregateRootId())
	}
	return g.value.SkuArray
}

// 获取商品的规格
func (g *itemImpl) SpecArray() promodel.SpecList {
	return g.repo.SkuService().GetSpecArray(g.SkuArray())
}

// 获取SKU
func (g *itemImpl) GetSku(skuId int64) *item.Sku {
	if g.value.SkuArray != nil {
		for _, v := range g.value.SkuArray {
			if v.ID == skuId {
				return v
			}
		}
	}
	return g.repo.GetItemSku(skuId)
}

// ========== [/ SKU处理结束 ] ===========//

// 从产品中拷贝信息
//todo: 如后期弄成公共产品，则应保持产品与商品的数据独立。
func (g *itemImpl) copyFromProduct(v *item.GoodsItem) error {
	pro := g.productRepo.GetProductValue(v.ProductId)
	if pro == nil {
		return product.ErrNoSuchProduct
	}
	//g.value.CatId = pro.CatId
	g.value.VendorId = pro.VendorId
	g.value.BrandId = pro.BrandId
	if g.value.Title == "" {
		g.value.Title = pro.Name
	}
	if g.value.Code == "" {
		g.value.Code = pro.Code
	}
	g.value.Image = pro.Image
	g.value.SortNum = pro.SortNum
	g.value.CreateTime = pro.CreateTime
	g.value.UpdateTime = pro.UpdateTime
	return nil
}

// 重置审核状态
func (i *itemImpl) resetReview() {
	i.value.ReviewState = enum.ReviewAwaiting
}

// 检查商品数据是否正确
func (g *itemImpl) checkItemValue(v *item.GoodsItem) error {
	defaultImage := g.registryRepo.Get(registry.GoodsDefaultImage).StringValue()
	// 检测是否上传图片
	if v.Image == "" || v.Image == defaultImage {
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
func (i *itemImpl) checkPrice(v *item.GoodsItem) error {
	rate := (v.Price - v.Cost) / v.Price
	minRate := i.registryRepo.Get(registry.GoodsMinProfitRate).FloatValue()
	// 如果未设定最低利润率，则可以与供货价一致
	if minRate != 0 && float64(rate) < minRate {
		return errors.New(fmt.Sprintf(item.ErrGoodsMinProfitRate.Error(),
			strconv.Itoa(int(minRate*100))+"%"))
	}
	return nil
}

// 保存
func (g *itemImpl) Save() (_ int64, err error) {
	ss := g.repo.SkuService()
	// 保存SKU
	if g.value.SkuArray != nil {
		err = ss.UpgradeBySku(g.value, g.value.SkuArray)
		if err == nil {
			// 创建商品
			if g.GetAggregateRootId() <= 0 {
				g.value.ID, err = g.repo.SaveValueGoods(g.value)
			}
			// 保存商品SKU
			if err == nil {
				err = g.saveItemSku(&g.value.SkuArray)
				// 设置默认SKU
				g.value.SkuId = 0
				if l := len(g.value.SkuArray); l > 0 && err == nil {
					g.value.SkuId = g.value.SkuArray[0].ID
				}
			}
		}
		if err != nil {
			return g.value.ID, err
		}
	}

	// 保存商品
	g.value.ID, err = g.repo.SaveValueGoods(g.value)
	if err == nil {
		g.snapshot = nil
		// 保存商品快照
		_, err = g.repo.SnapshotService().GenerateSnapshot(g.value)
	}
	return g.value.ID, err
}

// 获取促销信息
func (g *itemImpl) GetPromotions() []promotion.IPromotion {
	//todo: 商品促销
	return []promotion.IPromotion{}

	var vp []*promotion.PromotionInfo = g.promRepo.GetPromotionOfGoods(
		g.GetAggregateRootId())
	var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
	for i, v := range vp {
		proms[i] = g.promRepo.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (g *itemImpl) GetLevelPrice(level int) (bool, float32) {
	lvp := g.GetLevelPrices()
	for _, v := range lvp {
		if level == v.Level && v.Price < g.value.Price {
			return true, v.Price
		}
	}
	return false, g.value.Price
}

// 获取促销价
func (g *itemImpl) GetPromotionPrice(level int) float32 {
	b, price := g.GetLevelPrice(level)
	if b {
		return price
	}
	return g.value.Price
}

// 获取促销描述
func (g *itemImpl) GetPromotionDescribe() map[string]string {
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
			//					g._promDescribes[key] = v.Value().ShortName
			//				} else {
			//					g._promDescribes[key] = txt + ";" + v.Value().ShortName
			//				}
			//			} else if v.Type() == promotion.TypeFlagCoupon {
			//				if txt, ok := g._promDescribes[key]; !ok {
			//					g._promDescribes[key] = v.Value().ShortName
			//				} else {
			//					g._promDescribes[key] = txt + ";" + v.Value().ShortName
			//				}
			//			}

			//todo: other promotion implement
		}
	}
	return g.promDescribes
}

// 获取会员价
func (g *itemImpl) GetLevelPrices() []*item.MemberPrice {
	if g.levelPrices == nil {
		g.levelPrices = g.repo.GetGoodsLevelPrice(g.GetAggregateRootId())
	}
	return g.levelPrices
}

// 保存会员价
func (g *itemImpl) SaveLevelPrice(v *item.MemberPrice) (int32, error) {
	v.GoodsId = g.GetAggregateRootId()
	if g.value.Price == v.Price {
		if v.Id > 0 {
			g.repo.RemoveGoodsLevelPrice(v.Id)
		}
		return -1, nil
	}
	return g.repo.SaveGoodsLevelPrice(v)
}

// 是否上架
func (g *itemImpl) IsOnShelves() bool {
	return g.value.ShelveState == item.ShelvesOn
}

// 设置上架
func (g *itemImpl) SetShelve(state int32, remark string) error {
	if state == item.ShelvesIncorrect && len(remark) == 0 {
		return product.ErrNilRejectRemark
	}
	g.value.ShelveState = state
	if g.value.ReviewState != enum.ReviewPass {
		g.value.ReviewState = enum.ReviewAwaiting
	}
	g.value.ReviewRemark = remark
	_, err := g.Save()
	return err
}

// 标记为违规
func (g *itemImpl) Incorrect(remark string) error {
	g.value.ShelveState = item.ShelvesIncorrect
	g.value.ReviewRemark = remark
	_, err := g.Save()
	return err
}

// 审核
func (g *itemImpl) Review(pass bool, remark string) error {
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
func (g *itemImpl) AddSalesNum(skuId int64, quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	//log.Println("--商品：",g.value.ID,"; 库存：",
	// g.value.StockNum,"; 数量:",quantity)
	if quantity > g.value.StockNum {
		return item.ErrOutOfStock
	}
	g.value.SaleNum += quantity
	_, err := g.Save()
	if err == nil {
		if sku := g.GetSku(skuId); sku != nil {
			sku.SaleNum += quantity
			_, err = g.saveSku(sku)
		}
	}
	return err
}

// 取消销售
func (g *itemImpl) CancelSale(skuId int64, quantity int32, orderNo string) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	g.value.SaleNum -= quantity
	_, err := g.Save()
	if err == nil {
		if sku := g.GetSku(skuId); sku != nil {
			sku.SaleNum -= quantity
			_, err = g.saveSku(sku)
		}
	}
	return err
}

// 占用库存
func (g *itemImpl) TakeStock(skuId int64, quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	if quantity > g.value.StockNum {
		return item.ErrOutOfStock
	}
	g.value.StockNum -= quantity
	_, err := g.Save()
	if err == nil {
		if sku := g.GetSku(skuId); sku != nil {
			sku.Stock -= quantity
			_, err = g.saveSku(sku)
		}
	}
	return err
}

func (g *itemImpl) saveSku(sku *item.Sku) (_ int64, err error) {
	sku.ID, err = util.I64Err(g.repo.SaveItemSku(sku))
	return sku.ID, err
}

// 释放库存
func (g *itemImpl) FreeStock(skuId int64, quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	g.value.StockNum += quantity
	_, err := g.Save()
	if err == nil {
		if sku := g.GetSku(skuId); sku != nil {
			sku.Stock += quantity
			_, err = g.saveSku(sku)
		}
	}
	return err
}

// 删除商品
func (g *itemImpl) Destroy() error {
	//g.goodsRepo.
	return nil
}
