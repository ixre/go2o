/**
 * Copyright 2015 @ 56x.net.
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
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/express"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/merchant/shop"
	promodel "github.com/ixre/go2o/core/domain/interface/pro_model"
	"github.com/ixre/go2o/core/domain/interface/product"
	"github.com/ixre/go2o/core/domain/interface/promotion"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/go2o/core/infrastructure/regex"
	"github.com/ixre/gof/util"
)

var _ item.IGoodsItem = new(itemImpl)

// 商品实现
type itemImpl struct {
	pro             product.IProduct
	value           *item.GoodsItem
	wholesale       item.IWholesaleItem
	snapshot        *item.Snapshot
	repo            item.IItemRepo
	catRepo         product.ICategoryRepo
	productRepo     product.IProductRepo
	itemWsRepo      item.IItemWholesaleRepo
	proMRepo        promodel.IProductModelRepo
	promRepo        promotion.IPromotionRepo
	levelPrices     []*item.MemberPrice
	images          []string
	promDescribes   map[string]string
	registryRepo    registry.IRegistryRepo
	expressRepo     express.IExpressRepo
	shopRepo        shop.IShopRepo
	awaitSaveImages []*item.Image
}

//todo:??? 去掉依赖promotion.IPromotionRepo

func NewItem(
	itemRepo product.IProductRepo, catRepo product.ICategoryRepo,
	pro product.IProduct, value *item.GoodsItem, registryRepo registry.IRegistryRepo,
	goodsRepo item.IItemRepo, proMRepo promodel.IProductModelRepo,
	itemWsRepo item.IItemWholesaleRepo, expressRepo express.IExpressRepo,
	shopRepo shop.IShopRepo,
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
		shopRepo:     shopRepo,
		expressRepo:  expressRepo,
	}
	return v.init()
}

func (i *itemImpl) Images() []string {
	if i.images == nil {
		arr := i.repo.GetItemImages(i.GetAggregateRootId())
		i.images = make([]string, len(arr))
		for k, v := range arr {
			i.images[k] = v.ImageUrl
		}
	}
	return i.images
}

// GrantFlag 添加商品标志
func (i *itemImpl) GrantFlag(flag int) error {
	f := int(math.Abs(float64(flag)))
	if f&(f-1) != 0 {
		return errors.New("not right flag value")
	}
	if flag > 0 { // 添加标志
		if i.value.ItemFlag&f != f {
			i.value.ItemFlag |= flag
		}
	} else { // 去除标志
		if i.value.ItemFlag&f == f {
			i.value.ItemFlag ^= f
		}
	}
	_, err := i.Save()
	return err
}

func (i *itemImpl) SetImages(images []string) error {
	if len(images) == 0 {
		return errors.New("images not set")
	}
	// 构造新数组并合并旧有数据
	old := i.repo.GetItemImages(i.GetAggregateRootId())
	oldMap := make(map[string]*item.Image, len(old))
	delArr := make([]int64, 0)
	for _, v := range old {
		exists := false
		for _, v2 := range images {
			if v2 == v.ImageUrl {
				exists = true
			}
		}
		if !exists || oldMap[v.ImageUrl] != nil {
			delArr = append(delArr, v.Id)
			continue
		}
		oldMap[v.ImageUrl] = v
	}
	arr := make([]*item.Image, len(images))
	for k, v := range images {
		arr[k] = &item.Image{
			Id:         0,
			ItemId:     i.GetAggregateRootId(),
			ImageUrl:   v,
			SortNum:    k,
			CreateTime: time.Now().Unix(),
		}
		if it, ok := oldMap[v]; ok {
			arr[k].Id = it.Id
		}
	}
	// 删除项
	for _, v := range delArr {
		i.repo.DeleteItemImage(v)
	}
	// 设置商品主图
	if len(images) > 0 {
		i.value.Image = images[0]
	}
	i.awaitSaveImages = arr

	// 清除图片数据
	i.images = nil
	return nil
}

func (i *itemImpl) init() item.IGoodsItem {
	if i.pro != nil {
		i.value.PromPrice = i.value.Price
	}
	return i
}

// GetAggregateRootId 获取聚合根编号
func (i *itemImpl) GetAggregateRootId() int64 {
	return i.value.Id
}

// Snapshot 商品快照
func (i *itemImpl) Snapshot() *item.Snapshot {
	if i.snapshot == nil {
		i.snapshot = i.repo.SnapshotService().GetLatestSnapshot(
			i.GetAggregateRootId())
	}
	return i.snapshot
}

// Product 获取货品
func (i *itemImpl) Product() product.IProduct {
	if i.pro == nil && i.value.ProductId > 0 {
		i.pro = i.productRepo.GetProduct(i.value.ProductId)
	}
	return i.pro
}

// Wholesale 批发
func (i *itemImpl) Wholesale() item.IWholesaleItem {
	if i.wholesale == nil {
		i.wholesale = newWholesaleItem(i.GetAggregateRootId(),
			i, i.repo, i.itemWsRepo)
	}
	return i.wholesale
}

// GetValue 获取值
func (i *itemImpl) GetValue() *item.GoodsItem {
	return i.value
}

// GetPackedValue 获取包装过的商品信息
func (i *itemImpl) GetPackedValue() *valueobject.Goods {
	//item := i.GetItem().Value()
	gv := i.GetValue()
	goods := &valueobject.Goods{
		ProductId:   gv.ProductId,
		CategoryId:  gv.CategoryId,
		Title:       gv.Title,
		GoodsNo:     gv.Code,
		Image:       gv.Image,
		IntroVideo:  gv.IntroVideo,
		RetailPrice: gv.RetailPrice,
		Price:       gv.Price,
		PriceRange:  gv.PriceRange,
		PromPrice:   gv.Price,
		ItemId:      i.GetAggregateRootId(),
		SkuId:       gv.SkuId,
		ItemFlag:    gv.ItemFlag,
		StockNum:    gv.StockNum,
		SaleNum:     gv.SaleNum,
	}
	return goods
}

// SetValue 设置值
func (i *itemImpl) SetValue(v *item.GoodsItem) error {
	err := i.checkItemValue(v)
	if err == nil {
		// 检测店铺
		if v.ShopId <= 0 {
			return item.ErrIncorrectShopOfItem
		}
		isp := i.shopRepo.GetShop(v.ShopId)
		if isp == nil {
			return shop.ErrNoSuchShop
		}
		i.value.ShopId = v.ShopId
		i.value.VendorId = isp.GetValue().VendorId
		// 创建商品时，设为已下架
		if i.GetAggregateRootId() <= 0 {
			i.value.ShelveState = item.ShelvesDown
			// 分类在创建后，不允许再进行修改。并且分类不能为虚拟分类
			// 如果修改，则所有SKU和属性应删除。
			c := i.catRepo.GlobCatService().GetCategory(int(v.CategoryId))
			if c == nil {
				return item.ErrIncorrectProductCategory
			}
			cv := c.GetValue()
			if cv.VirtualCat == 1 {
				return item.ErrIncorrectProductCategory
			}
			i.value.CategoryId = v.CategoryId
		} else {
			if err = i.copyFromProduct(v); err != nil {
				return err
			}
		}
		i.value.ShopCatId = v.ShopCatId
		i.value.ProductId = v.ProductId
		i.value.ShopCatId = v.ShopCatId
		i.value.ExpressTid = v.ExpressTid
		i.value.Title = v.Title
		i.value.ItemFlag = v.ItemFlag
		i.value.ShortTitle = v.ShortTitle
		i.value.Code = v.Code
		i.value.SaleNum = v.SaleNum
		i.value.StockNum = v.StockNum
		i.value.Cost = v.Cost
		i.value.RetailPrice = v.RetailPrice
		i.value.Price = v.Price
		i.value.Weight = v.Weight
		i.value.Bulk = v.Bulk
		// 更新视频
		i.value.IntroVideo = v.IntroVideo
		// 如果零售价和成本未填写,则默认等于价格
		if i.value.RetailPrice == 0 {
			i.value.RetailPrice = v.Price
		}
		if i.value.Cost == 0 {
			i.value.Cost = v.Price
		}
		//设置默认的价格区间
		if len(i.SkuArray()) == 0 {
			i.value.PriceRange = format.FormatFloat64(float64(v.Price) / 100)
		}
		if i.value.CreateTime == 0 {
			i.value.CreateTime = time.Now().Unix()
		}
		//修改图片或标题后，要重新审核
		if i.value.Image != v.Image || i.value.Title != v.Title {
			i.resetReview()
		}
		// 更新包邮标志
		i.autoSetFlagValues(isp)
	}
	return err
}

// 自动设置一些标志值
func (i *itemImpl) autoSetFlagValues(isp shop.IShop) {
	// 更新包邮标志
	if i.value.ExpressTid <= 0 {
		i.GrantFlag(item.FlagFreeDelivery)
	} else {
		i.GrantFlag(-item.FlagFreeDelivery)
	}
	// 更新自营标志
	iol := isp.(shop.IOnlineShop)
	if iol != nil && domain.TestFlag(iol.GetShopValue().Flag, shop.FlagSelfSales) {
		i.GrantFlag(item.FlagSelfSales)
	} else {
		i.GrantFlag(-item.FlagSelfSales)
	}
}

// SetSku 设置SKU
func (i *itemImpl) SetSku(arr []*item.Sku) error {
	for _, v := range arr {
		if v.RetailPrice == 0 {
			v.RetailPrice = v.Price
		}
		if v.Cost == 0 {
			v.Cost = v.Price
		}
	}
	i.value.SkuArray = arr
	return nil
}

// ========== [# SKU处理开始 ]  ===========//

// 保存商品SKU
func (i *itemImpl) saveItemSku(arrPtr *[]*item.Sku) (err error) {
	arr := *arrPtr
	pk := i.GetAggregateRootId()
	ss := i.repo.SkuService()
	// 格式化数据
	err = ss.RebuildSkuArray(&arr, i.value)
	if err == nil {
		// 获取之前的SKU设置
		old := i.repo.SelectItemSku("item_id= $1", pk)
		// 合并SKU
		ss.Merge(old, &arr)
		// 分析当前项目并加入到MAP中
		var delList []int64
		currMap := make(map[int64]*item.Sku, len(arr))
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
			i.repo.DeleteItemSku(v)
		}
		// 保存项
		for _, v := range arr {
			if v.ItemId == 0 {
				v.ItemId = pk
			}
			if proId := i.value.ProductId; v.ProductId != proId {
				v.ProductId = proId
			}
			if v.ItemId == pk {
				v.Id, err = util.I64Err(i.repo.SaveItemSku(v))
			}
		}
	}
	return err
}

// SkuArray 获取SKU数组
func (i *itemImpl) SkuArray() []*item.Sku {
	if i.value.SkuArray == nil {
		i.value.SkuArray = i.repo.SelectItemSku("item_id= $1",
			i.GetAggregateRootId())
	}
	return i.value.SkuArray
}

// SpecArray 获取商品的规格
func (i *itemImpl) SpecArray() promodel.SpecList {
	return i.repo.SkuService().GetItemSpecArray(i.SkuArray())
}

// GetSku 获取SKU
func (i *itemImpl) GetSku(skuId int64) *item.Sku {
	if i.value.SkuArray != nil {
		for _, v := range i.value.SkuArray {
			if v.Id == skuId {
				return v
			}
		}
	}
	return i.repo.GetItemSku(skuId)
}

// ========== [/ SKU处理结束 ] ===========//

// 从产品中拷贝信息
// todo: 如后期弄成公共产品，则应保持产品与商品的数据独立。
func (i *itemImpl) copyFromProduct(v *item.GoodsItem) error {
	pro := i.productRepo.GetProductValue(v.ProductId)
	if pro == nil {
		return product.ErrNoSuchProduct
	}
	//i.value.CategoryId = pro.CategoryId
	i.value.VendorId = pro.VendorId
	i.value.BrandId = pro.BrandId
	if i.value.Title == "" {
		i.value.Title = pro.Name
	}
	if i.value.Code == "" {
		i.value.Code = pro.Code
	}
	//i.value.Image = pro.Image
	i.value.SortNum = pro.SortNum
	i.value.CreateTime = pro.CreateTime
	i.value.UpdateTime = pro.UpdateTime
	return nil
}

// 重置审核状态
func (i *itemImpl) resetReview() {
	i.value.ReviewState = enum.ReviewAwaiting
}

// 检查商品数据是否正确
func (i *itemImpl) checkItemValue(v *item.GoodsItem) error {
	if b, _ := regex.ContainInvalidChars(v.Title); b {
		return item.ErrInvalidTitle
	}
	defaultImage := i.registryRepo.Get(registry.GoodsDefaultImage).StringValue()
	// 检测是否上传图片
	if v.Image == "" || v.Image == defaultImage {
		return product.ErrNotUploadImage
	}
	// 检测运费模板
	if v.ExpressTid > 0 {
		ve := i.expressRepo.GetUserExpress(int(v.VendorId))
		tpl := ve.GetTemplate(int(v.ExpressTid))
		if tpl == nil {
			return express.ErrNoSuchTemplate
		}
		if !tpl.Enabled() {
			return express.ErrTemplateNotEnabled
		}
	}
	// 检测价格
	return i.checkPrice(v)
}

// 判断价格是否正确
func (i *itemImpl) checkPrice(v *item.GoodsItem) error {
	if v.Price == 0 {
		return nil
	}
	rate := (v.Price - v.Cost) / v.Price
	minRate := i.registryRepo.Get(registry.GoodsMinProfitRate).FloatValue()
	// 如果未设定最低利润率，则可以与供货价一致
	if minRate != 0 && float64(rate) < minRate {
		return fmt.Errorf(item.ErrGoodsMinProfitRate.Error(),
			strconv.Itoa(int(minRate*100))+"%")
	}
	return nil
}

// Save 保存
func (i *itemImpl) Save() (_ int64, err error) {
	// 保存SKU
	if i.value.SkuArray != nil {
		ss := i.repo.SkuService()
		err = ss.UpgradeBySku(i.value, i.value.SkuArray)
		if err == nil {
			// 创建商品
			if i.GetAggregateRootId() <= 0 {
				i.value.Id, err = i.repo.SaveValueGoods(i.value)
			}
			// 保存商品SKU
			if err == nil {
				err = i.saveItemSku(&i.value.SkuArray)
				// if err == nil {
				// 	// 设置默认SKU
				// 	i.value.SkuId = 0
				// 	// 按价格排序
				// 	list := item.SkuList(i.value.SkuArray)
				// 	sort.Sort(list)
				// 	// 如果SKU不为空
				// 	if l := len(i.value.SkuArray); l > 0 && err == nil {
				// 		i.value.SkuId = i.value.SkuArray[0].Id
				// 	}
				// }
			}
		}
		if err != nil {
			return i.value.Id, err
		}
	}
	// 新增自营商品自动上架(待审核)
	if i.GetAggregateRootId() == 0 {
		log.Println("[ GO2O][ LOG]: new item", i.value.ItemFlag)
		if domain.TestFlag(i.value.ItemFlag, item.FlagSelfSales) {
			i.value.ShelveState = item.ShelvesOn
			i.value.ReviewState = enum.ReviewAwaiting
		}
	}

	// 保存商品
	i.value.Id, err = i.repo.SaveValueGoods(i.value)
	if err == nil {
		i.snapshot = nil
		// 保存商品快照
		_, err = i.repo.SnapshotService().GenerateSnapshot(i.value)
		if err == nil {
			// 保存商品图片
			if i.awaitSaveImages != nil {
				for _, v := range i.awaitSaveImages {
					v.ItemId = i.GetAggregateRootId()
					i.repo.SaveItemImage(v)
				}
			}
		}
	}
	return i.value.Id, err
}

// 获取促销信息
func (i *itemImpl) GetPromotions() []promotion.IPromotion {
	//todo: 商品促销
	return []promotion.IPromotion{}

	var vp = i.promRepo.GetPromotionOfGoods(
		i.GetAggregateRootId())
	var proms = make([]promotion.IPromotion, len(vp))
	for j, v := range vp {
		proms[j] = i.promRepo.CreatePromotion(v)
	}
	return proms
}

// 获取会员价销价
func (i *itemImpl) GetLevelPrice(level int) (bool, int64) {
	lvp := i.GetLevelPrices()
	for _, v := range lvp {
		if level == v.LevelId && v.Price < i.value.Price {
			return true, v.Price
		}
	}
	return false, i.value.Price
}

// 获取促销价
func (i *itemImpl) GetPromotionPrice(level int) int64 {
	b, price := i.GetLevelPrice(level)
	if b {
		return price
	}
	return i.value.Price
}

// 获取促销描述
func (i *itemImpl) GetPromotionDescribe() map[string]string {
	if i.promDescribes == nil {
		proms := i.GetPromotions()
		i.promDescribes = make(map[string]string, len(proms))
		for _, v := range proms {
			key := v.TypeName()
			if txt, ok := i.promDescribes[key]; !ok {
				i.promDescribes[key] = v.GetValue().ShortName
			} else {
				i.promDescribes[key] = txt + ";" + v.GetValue().ShortName
			}

			//			if v.AdType() == promotion.TypeFlagCashBack {
			//				if txt, ok := i._promDescribes[key]; !ok {
			//					i._promDescribes[key] = v.Value().ShortName
			//				} else {
			//					i._promDescribes[key] = txt + ";" + v.Value().ShortName
			//				}
			//			} else if v.AdType() == promotion.TypeFlagCoupon {
			//				if txt, ok := i._promDescribes[key]; !ok {
			//					i._promDescribes[key] = v.Value().ShortName
			//				} else {
			//					i._promDescribes[key] = txt + ";" + v.Value().ShortName
			//				}
			//			}

			//todo: other promotion implement
		}
	}
	return i.promDescribes
}

// 获取会员价
func (i *itemImpl) GetLevelPrices() []*item.MemberPrice {
	if i.levelPrices == nil {
		i.levelPrices = i.repo.GetGoodSMemberLevelPrice(i.GetAggregateRootId())
	}
	return i.levelPrices
}

// 保存会员价
func (i *itemImpl) SaveLevelPrice(v *item.MemberPrice) (int32, error) {
	v.GoodsId = i.GetAggregateRootId()
	if i.value.Price == v.Price {
		if v.Id > 0 {
			i.repo.RemoveGoodSMemberLevelPrice(v.Id)
		}
		return -1, nil
	}
	return i.repo.SaveGoodSMemberLevelPrice(v)
}

// 是否上架
func (i *itemImpl) IsOnShelves() bool {
	return i.value.ShelveState == item.ShelvesOn
}

// 设置上架
func (i *itemImpl) SetShelve(state int32, remark string) error {
	if state == item.ShelvesIncorrect && len(remark) == 0 {
		return product.ErrNilRejectRemark
	}
	i.value.ShelveState = state
	if i.value.ReviewState != enum.ReviewPass {
		i.value.ReviewState = enum.ReviewAwaiting
	}
	i.value.ReviewRemark = remark
	_, err := i.Save()
	return err
}

// 标记为违规
func (i *itemImpl) Incorrect(remark string) error {
	i.value.ShelveState = item.ShelvesIncorrect
	i.value.ReviewRemark = remark
	_, err := i.Save()
	return err
}

// 审核
func (i *itemImpl) Review(pass bool, remark string) error {
	if pass {
		i.value.ReviewState = enum.ReviewPass
	} else {
		remark = strings.TrimSpace(remark)
		if remark == "" {
			return item.ErrEmptyReviewRemark
		}
		i.value.ShelveState = item.ShelvesDown
		i.value.ReviewState = enum.ReviewReject
	}
	i.value.ReviewRemark = remark
	_, err := i.Save()
	return err
}

// 更新销售数量
func (i *itemImpl) AddSalesNum(skuId int64, quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	//log.Println("--商品：",i.value.Id,"; 库存：",
	// i.value.StockNum,"; 数量:",quantity)
	if quantity > i.value.StockNum {
		return item.ErrOutOfStock
	}
	i.value.SaleNum += quantity
	_, err := i.Save()
	if err == nil {
		if sku := i.GetSku(skuId); sku != nil {
			sku.SaleNum += quantity
			_, err = i.saveSku(sku)
		}
	}
	return err
}

// 取消销售
func (i *itemImpl) CancelSale(skuId int64, quantity int32, orderNo string) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	i.value.SaleNum -= quantity
	_, err := i.Save()
	if err == nil {
		if sku := i.GetSku(skuId); sku != nil {
			sku.SaleNum -= quantity
			_, err = i.saveSku(sku)
		}
	}
	return err
}

// 占用库存
func (i *itemImpl) TakeStock(skuId int64, quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	if quantity > i.value.StockNum {
		return item.ErrOutOfStock
	}
	i.value.StockNum -= quantity
	_, err := i.Save()
	if err == nil {
		if sku := i.GetSku(skuId); sku != nil {
			sku.Stock -= quantity
			_, err = i.saveSku(sku)
		}
	}
	return err
}

func (i *itemImpl) saveSku(sku *item.Sku) (_ int64, err error) {
	sku.Id, err = util.I64Err(i.repo.SaveItemSku(sku))
	return sku.Id, err
}

// 释放库存
func (i *itemImpl) ReleaseStock(skuId int64, quantity int32) error {
	if quantity <= 0 {
		return item.ErrGoodsNum
	}
	i.value.StockNum += quantity
	_, err := i.Save()
	if err == nil {
		if sku := i.GetSku(skuId); sku != nil {
			sku.Stock += quantity
			_, err = i.saveSku(sku)
		}
	}
	return err
}

// 回收商品
func (i *itemImpl) Recycle() error {
	if i.value.IsRecycle == 0 {
		i.value.IsRecycle = 1
		_, err := i.Save()
		return err
	}
	return nil
}

// 从回收站中撤回
func (i *itemImpl) RecycleRevert() error {
	if i.value.IsRecycle == 1 {
		i.value.IsRecycle = 0
		_, err := i.Save()
		return err
	}
	return nil
}

// 删除商品
func (i *itemImpl) Destroy() error {
	if i.value.IsRecycle == 0 {
		return errors.New("商品非回收状态，不允许删除")
	}
	//i.goodsRepo.
	return nil
}
