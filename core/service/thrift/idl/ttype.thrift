namespace go define

//传输结果对象
struct Result{
   /* 状态码,如为0表示成功 */
   1:i32 ErrCode
   /* 消息 */
   2:string ErrMsg
   /* 数据, 可以用来存放JSON字符串 */
   3:string Data
}

//传输结果对象
struct Result64{
   1:i64 ID
   2:bool Result
   3:string Code
   4:string ErrMsg
}

//传输结果对象(Double)
struct DResult{
   1:double Data
   2:bool Result
   3:string Code
   4:string ErrMsg
}

// 键值对
struct Pair{
   1:string Key
   2:string Value
}

/** 设置依据 */
enum  SettingBasis {
    /** 未设置 */
    None = 1,
    /** 使用全局 */
    Global = 2,
    /** 自定义 */
    Custom = 3
}

/** 价格计算方式 */
enum  PriceBasis{
    /** 原价 */
    Original = 1,
    /** 会员折扣价 */
    Discount = 2,
    /** 自定义价格 */
    Custom = 3,
}

/** 金额/提成依据 */
enum AmountBasis{
    /** 未设置 */
    NotSet = 1,
    /** 按金额 */
    Amount = 2,
    /** 按百分比 */
    Percent =3
}

/** 百分比比例放大倍数，保留3位小数;0.56 * 10000 = 560 */
const i32 RATE_PERCENT = 10000
/** 金额比例放大倍数;0.95 * 100 = 95  */
const i32 RATE_AMOUNT = 100
/** 折扣比例放大倍数; 0.9 * 1000 = 900 */
const i32 RATE_DISCOUNT = 1000

//商品分类
struct Category {
    1: i32 ID
    2: i32 ParentId
    3: i32 ProModel
    // 分类优先级
    4: i32 Priority
    // 分类名称
    5: string Name
    // 虚拟分类
    6: i32 VirtualCat
    // 分类目标地址
    7: string CatUrl
    8: string Icon
    // 图标（雪碧图）坐标
    9: string IconXY
    10: i32 SortNum
    11: i32 FloorShow
    12: i32 Enabled
    13: i32 Level
    14: i64 CreateTime
    15: list<Category> Children
}


struct Sku {
    1: i64 SkuId
    2: i64 ItemId
    3: i64 ProductId
    4: string Title
    5: string Image
    6: string SpecData
    7: string SpecWord
    8: string Code
    9: double RetailPrice
    10: double Price
    11: double Cost
    12: i32 Weight
    13: i32 Bulk
    14: i32 Stock
    15: i32 SaleNum
}


struct OldItem{
    1: i64 ItemId
    2: i64 ProductId
    3: i32 PromFlag
    4: i32 CatId
    5: i32 VendorId
    6: i32 BrandId
    7: i32 ShopId
    8: i32 ShopCatId
    9: i32 ExpressTid
    10: string Title
    11: string ShortTitle
    12: string Code
    13: string Image
    14: i32 IsPresent
    15: string PriceRange
    16: i32 StockNum
    17: i32 SaleNum
    18: i32 SkuNum
    19: i64 SkuId
    20: double Cost
    21: double Price
    22: double RetailPrice
    23: i32 Weight
    24: i32 Bulk
    25: i32 ShelveState
    26: i32 ReviewState
    27: string ReviewRemark
    28: i32 SortNum
    29: i64 CreateTime
    30: i64 UpdateTime
    31: double PromPrice
    32: list<Sku> SkuArray
    33: map<string,string> Data;
}

// 统一的商品对象
struct UnifiedItem{
    1: i64 ItemId
    2: i32 ProductId
    3: i32 CatId
    4: i32 VendorId
    5: i32 BrandId
    6: string Title
    7: string Code
    8: string Image
    9: double Price
    10: string PriceRange
    11: i32 StockNum
    12: i32 ShelveState
    13: i32 ReviewState
    14: i64 UpdateTime
    15: list<Sku> SkuArray
    16: map<string,string> Data
    // 3: i32 PromFlag
    // 7: i32 ShopId
    // 8: i32 ShopCatId
    // 9: i32 ExpressTid
    // 11: string ShortTitle
    // 14: i32 IsPresent
    // 23: i32 Weight
    // 24: i32 Bulk
    // 18: i32 SkuNum
    // 28: i32 SortNum
    // 29: i64 CreateTime
    // 31: double PromPrice
    // 22: double RetailPrice
    // 19: i32 SkuId
    // 20: double Cost
    // 17: i32 SaleNum

}



// 购物车
struct ShoppingCart {
    //编号
    1: i32 CartId
    //购物车KEY
    2: string Code
    //店铺分组
    3: list<ShoppingCartGroup> Shops
}
// 购物车商铺分组
struct ShoppingCartGroup {
    //商铺编号
    1: i32 ShopId
    //供货商编号
    2: i32 VendorId
    //商铺名称
    3: string ShopName
    //是否结算
    4: bool Checked
    //商品
    5: list<ShoppingCartItem> Items
}


// 购物车商品
struct ShoppingCartItem {
    //商品编号
    1: i64 ItemId
    //SKU编号
    2: i64 SkuId
    //商品标题
    3: string Title
    //商品图片
    4: string Image
    //规格文本
    5: string SpecWord
    //商品编码
    6: string Code
    //零售价
    7: double RetailPrice
    //销售价
    8: double Price
    //数量
    9: i32 Quantity
    //是否结算
    10: bool Checked
    //库存文本
    11: string StockText
    //店铺编号
    12: i32 ShopId
}



/** 分页参数 */
struct PagingParams{
    /** 参数 */
    1:map<string,string> Opt
    /** 排序字段 */
    2:string OrderField
    /** 是否倒序排列 */
    3:bool OrderDesc
    /** 开始记录数 */
    4:i32 Begin
    /** 结束记录数 */
    5:i32 Over
}

/** 分页结果 */
struct PagingResult{
    /** 代码 */
    1:i32 ErrCode
    /** 消息 */
    2:string ErrMsg
    /** 总数 */
    3:i32 Count
    /** 数据 */
    4:string Data
}