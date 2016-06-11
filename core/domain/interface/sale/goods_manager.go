/**
 * Copyright 2015 @ z3q.net.
 * name : sale_goods.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package sale

import (
    "go2o/core/domain/interface/promotion"
    "go2o/core/domain/interface/valueobject"
)

type(
    // 商品
    IGoods interface {
        // 获取领域对象编号
        GetDomainId() int

        // 获取货品
        GetItem() IItem

        // 设置值
        GetValue() *ValueGoods

        // 获取包装过的商品信息
        GetPackedValue() *valueobject.Goods

        // 获取促销信息
        GetPromotions() []promotion.IPromotion

        // 获取促销价
        GetPromotionPrice(level int) float32

        // 获取会员价销价,返回是否有会原价及价格
        GetLevelPrice(level int) (bool, float32)

        // 获取促销描述
        GetPromotionDescribe() map[string]string

        // 获取会员价
        GetLevelPrices() []*MemberPrice

        // 保存会员价
        SaveLevelPrice(*MemberPrice) (int, error)

        // 设置值
        SetValue(*ValueGoods) error

        // 保存
        Save() (int, error)

        // 更新销售数量,扣减库存
        AddSaleNum(quantity int) error

        // 取消销售
        CancelSale(quantity int, orderNo string) error

        // 生成快照
        GenerateSnapshot() (int, error)

        // 获取最新的快照
        GetLatestSnapshot() *GoodsSnapshot
    }

    // 商品服务
    IGoodsManager interface{
        // 创建商品
        CreateGoodsByItem(IItem, *ValueGoods) IGoods

        // 创建商品
        CreateGoods(*ValueGoods) IGoods

        // 根据产品编号获取商品
        GetGoods(int) IGoods

        // 根据产品SKU获取商品
        GetGoodsBySku(itemId, sku int) IGoods

        // 删除商品
        DeleteGoods(int) error

        // 获取指定的商品快照
        GetGoodsSnapshot(id int) *GoodsSnapshot

        // 根据Key获取商品快照
        GetGoodsSnapshotByKey(key string) *GoodsSnapshot

        // 获取指定数量已上架的商品
        GetOnShelvesGoods(start, end int, sortBy string) []*valueobject.Goods
    }

    // 商品仓储
    IGoodsRep interface {
        // 获取商品
        GetValueGoods(itemId int, sku int) *ValueGoods

        // 获取商品
        GetValueGoodsById(goodsId int) *ValueGoods

        // 根据SKU获取商品
        GetValueGoodsBySku(itemId, sku int) *ValueGoods

        // 保存商品
        SaveValueGoods(*ValueGoods) (int, error)

        // 获取在货架上的商品
        GetOnShelvesGoods(merchantId int, start, end int,
        sortBy string) []*valueobject.Goods

        // 获取在货架上的商品
        GetPagedOnShelvesGoods(merchantId int, catIds []int, start, end int,
        where, orderBy string) (total int, goods []*valueobject.Goods)

        // 根据编号获取商品
        GetGoodsByIds(ids ...int) ([]*valueobject.Goods, error)

        // 获取会员价
        GetGoodsLevelPrice(goodsId int) []*MemberPrice

        // 保存会员价
        SaveGoodsLevelPrice(*MemberPrice) (int, error)

        // 移除会员价
        RemoveGoodsLevelPrice(id int) error
    }

    // 商品
    ValueGoods struct {
        Id            int `db:"id" pk:"yes" auto:"yes"`

        // 货品编号
        ItemId        int `db:"item_id"`

        // 是否为赠品
        IsPresent     int `db:"is_present"`

        // 规格
        SkuId         int `db:"sku_id"`

        // 促销标志
        PromotionFlag int `db:"prom_flag"`

        // 库存
        StockNum      int `db:"stock_num"`

        // 已售件数
        SaleNum       int `db:"sale_num"`

        // 销售价
        SalePrice     float32 `db:"-"`

        // 促销价
        PromPrice     float32 `db:"-"`

        // 实际价
        Price         float32 `db:"-"`
    }


    // 商品快照
    GoodsSnapshot struct {
        Id           int    `db:"id" auto:"yes" pk:"yes"`
        Key          string `db:"snapshot_key"`
        ItemId       int    `db:"item_id"`
        GoodsId      int    `db:"goods_id"`
        GoodsName    string `db:"goods_name"`
        GoodsNo      string `db:"goods_no"`
        SmallTitle   string `db:"small_title"`
        CategoryName string `db:"category_name"`
        Image        string `db:"img"`

        //成本价
        Cost         float32 `db:"cost"`

        //定价
        Price        float32 `db:"price"`

        //销售价
        SalePrice    float32 `db:"sale_price"`
        CreateTime   int64   `db:"create_time"`
    }


    // 简单商品信息
    SimpleGoods struct {
        GoodsId    int    `json:"id"`
        GoodsImage string `json:"img"`
        Name       string `json:"name"`
        Quantity   string `json:"qty"`
    }

)
