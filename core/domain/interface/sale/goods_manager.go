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
)
