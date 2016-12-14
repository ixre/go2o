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
    "go2o/core/domain/interface/express"
    "go2o/core/domain/interface/item"
    "go2o/core/domain/interface/product"
    "go2o/core/domain/interface/promotion"
    "go2o/core/domain/interface/shipment"
    "go2o/core/domain/interface/valueobject"
    "strconv"
)

var _ item.IGoodsItem = new(tmpGoodsImpl)

// 临时的商品实现  todo: 要与item分开
type tmpGoodsImpl struct {
    pro           product.IProduct
    value         *item.GoodsItem
    goodsRepo     item.IGoodsItemRepo
    productRepo   product.IProductRepo
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
goodsRepo item.IGoodsItemRepo, expressRepo express.IExpressRepo,
promRepo promotion.IPromotionRepo) item.IGoodsItem {
    v := &tmpGoodsImpl{
        pro:         pro,
        value:       value,
        productRepo: itemRepo,
        goodsRepo:   goodsRepo,
        promRepo:    promRepo,
        valRepo:     valRepo,
        expressRepo: expressRepo,
    }
    return v.init()
}

func (g *tmpGoodsImpl) init() item.IGoodsItem {
    if g.pro != nil {
        g.value.PromPrice = g.value.Price
    }
    return g
}

//获取聚合根编号
func (g *tmpGoodsImpl) GetAggregateRootId() int32 {
    return g.value.Id
}

// 商品快照
func (g *tmpGoodsImpl) SnapshotManager() item.ISnapshotManager {
    if g.snapManager == nil {
        var item *product.Product
        gi := g.GetItem()
        if gi != nil {
            v := gi.GetValue()
            item = &v
        }
        g.snapManager = NewSnapshotManagerImpl(g.GetAggregateRootId(),
            g.goodsRepo, g.productRepo, g.GetValue(), item)
    }
    return g.snapManager
}

// 获取货品
func (g *tmpGoodsImpl) GetItem() product.IProduct {
    return g.pro
}

// 设置值
func (g *tmpGoodsImpl) GetValue() *item.GoodsItem {
    return g.value
}

// 获取包装过的商品信息
func (g *tmpGoodsImpl) GetPackedValue() *valueobject.Goods {
    item := g.GetItem().GetValue()
    gv := g.GetValue()
    goods := &valueobject.Goods{
        ProductId:     item.Id,
        CategoryId:    item.CategoryId,
        Name:          item.Name,
        GoodsNo:       item.Code,
        Image:         item.Image,
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

// 获取促销信息
func (g *tmpGoodsImpl) GetPromotions() []promotion.IPromotion {
    var vp []*promotion.PromotionInfo = g.promRepo.GetPromotionOfGoods(
        g.GetAggregateRootId())
    var proms []promotion.IPromotion = make([]promotion.IPromotion, len(vp))
    for i, v := range vp {
        proms[i] = g.promRepo.CreatePromotion(v)
    }
    return proms
}

// 获取会员价销价
func (g *tmpGoodsImpl) GetLevelPrice(level int32) (bool, float32) {
    lvp := g.GetLevelPrices()
    for _, v := range lvp {
        if level == v.Level && v.Price < g.value.Price {
            return true, v.Price
        }
    }
    return false, g.value.Price
}

// 获取促销价
func (g *tmpGoodsImpl) GetPromotionPrice(level int32) float32 {
    b, price := g.GetLevelPrice(level)
    if b {
        return price
    }
    return g.value.Price
}

// 获取促销描述
func (g *tmpGoodsImpl) GetPromotionDescribe() map[string]string {
    if g.promDescribes == nil {
        proms := g.GetPromotions()
        g.promDescribes = make(map[string]string, len(proms))
        for _, v := range proms {
            key := v.TypeName()
            if txt, ok := g.promDescribes[key]; !ok {
                g.promDescribes[key] = v.GetValue().ShortName
            } else {
                g.promDescribes[key] = txt + "；" + v.GetValue().ShortName
            }

            //			if v.Type() == promotion.TypeFlagCashBack {
            //				if txt, ok := g._promDescribes[key]; !ok {
            //					g._promDescribes[key] = v.GetValue().ShortName
            //				} else {
            //					g._promDescribes[key] = txt + "；" + v.GetValue().ShortName
            //				}
            //			} else if v.Type() == promotion.TypeFlagCoupon {
            //				if txt, ok := g._promDescribes[key]; !ok {
            //					g._promDescribes[key] = v.GetValue().ShortName
            //				} else {
            //					g._promDescribes[key] = txt + "；" + v.GetValue().ShortName
            //				}
            //			}

            //todo: other promotion implement
        }
    }
    return g.promDescribes
}

// 获取会员价
func (g *tmpGoodsImpl) GetLevelPrices() []*item.MemberPrice {
    if g.levelPrices == nil {
        g.levelPrices = g.goodsRepo.GetGoodsLevelPrice(g.GetAggregateRootId())
    }
    return g.levelPrices
}

// 保存会员价
func (g *tmpGoodsImpl) SaveLevelPrice(v *item.MemberPrice) (int32, error) {
    v.GoodsId = g.GetAggregateRootId()
    if g.value.Price == v.Price {
        if v.Id > 0 {
            g.goodsRepo.RemoveGoodsLevelPrice(v.Id)
        }
        return -1, nil
    }
    return g.goodsRepo.SaveGoodsLevelPrice(v)
}

// 判断价格是否正确
func (i *tmpGoodsImpl) checkPrice(v *item.GoodsItem) error {
    rate := (v.Price - v.Cost) / v.Price
    conf := i.valRepo.GetRegistry()
    minRate := conf.GoodsMinProfitRate
    // 如果未设定最低利润率，则可以与供货价一致
    if minRate != 0 && rate < minRate {
        return errors.New(fmt.Sprintf(item.ErrGoodsMinProfitRate.Error(),
            strconv.Itoa(int(minRate * 100)) + "%"))
    }
    return nil
}

func (g *tmpGoodsImpl) checkItemValue(v *item.GoodsItem) error {
    registry := g.valRepo.GetRegistry()
    // 检测是否上传图片
    if v.Image == registry.GoodsDefaultImage {
        return product.ErrNotUploadImage
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

// 设置值
func (g *tmpGoodsImpl) SetValue(v *item.GoodsItem) error {
    g.value.IsPresent = v.IsPresent
    g.value.SaleNum = v.SaleNum
    g.value.StockNum = v.StockNum
    g.value.SkuId = v.SkuId
    //g._value.PromotionFlag = v.PromotionFlag
    err := g.checkItemValue(v)
    return err
}

// 保存
func (g *tmpGoodsImpl) Save() (int32, error) {
    id, err := g.goodsRepo.SaveValueGoods(g.value)
    if err == nil {
        g.value.Id = id
        _, err = g.SnapshotManager().GenerateSnapshot()
    }
    //todo: save promotion
    return id, err
}

// 更新销售数量
func (g *tmpGoodsImpl) AddSalesNum(quantity int32) error {
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
func (g *tmpGoodsImpl) CancelSale(quantity int32, orderNo string) error {
    if quantity <= 0 {
        return item.ErrGoodsNum
    }
    g.value.SaleNum -= quantity
    _, err := g.Save()
    return err
}

// 占用库存
func (g *tmpGoodsImpl) TakeStock(quantity int32) error {
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
func (g *tmpGoodsImpl) FreeStock(quantity int32) error {
    if quantity <= 0 {
        return item.ErrGoodsNum
    }
    g.value.StockNum += quantity
    _, err := g.Save()
    return err
}

// 删除商品
func (g *tmpGoodsImpl) Destroy() error {
    //g.goodsRepo.
    return nil
}
