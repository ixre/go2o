/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-08 10:49
 * description :
 * history :
 */

package sale

type(
    // 物品
    IItem interface {
        GetDomainId() int

        // 获取商品的值
        GetValue() Item

        // 是否上架
        IsOnShelves() bool

        // 获取销售标签
        GetSaleLabels() []*Label

        // 保存销售标签
        SaveSaleLabels([]int) error

        // 设置商品值
        SetValue(*Item) error

        // 保存
        Save() (int, error)
    }

    // 货品服务
    IItemManager interface {

    }

    // 商品值
    Item struct {
        Id          int    `db:"id" auto:"yes" pk:"yes"`
        CategoryId  int    `db:"category_id"`
        Name        string `db:"name"`
        //供应商编号(暂时同mch_id)
        SupplierId    int    `db:"supplier_id"`
        // 货号
        GoodsNo     string `db:"goods_no"`
        SmallTitle  string `db:"small_title"`
        Image       string `db:"img"`
        //成本价
        Cost        float32 `db:"cost"`
        //定价
        Price       float32 `db:"price"`
        //参考销售价
        SalePrice   float32 `db:"sale_price"`
        ApplySubs   string  `db:"apply_subs"`

        //简单备注,如:(限时促销)
        Remark      string `db:"remark"`
        Description string `db:"description"`

        // 是否上架,1为上架
        OnShelves   int `db:"on_shelves"`

        State       int   `db:"state"`
        CreateTime  int64 `db:"create_time"`
        UpdateTime  int64 `db:"update_time"`
    }

)
