namespace go define

//传输结果对象
struct Result{
   1:i32 Id
   2:bool Result
   3:i32 Code
   4:string Message
}

//传输结果对象(Double)
struct DResult{
   1:double Data
   2:bool Result
   3:string Message
}

// 键值对
struct Pair{
   1:string Key
   2:string Value
}

struct Member {
    1: i32 Id
    2: string Usr
    3: string Pwd
    4: string TradePwd
    5: i32 Exp
    6: i32 Level
    7: string InvitationCode
    8: string RegFrom
    9: string RegIp
    10: i64 RegTime
    11: string CheckCode
    12: i64 CheckExpires
    13: i32 State
    14: i64 LoginTime
    15: i64 LastLoginTime
    16: i64 UpdateTime
    17: string DynamicToken
    18: i64 TimeoutTime
}

struct Profile {
    1: i32 MemberId
    2: string Name
    3: string Avatar
    4: i32 Sex
    5: string BirthDay
    6: string Phone
    7: string Address
    8: string Im
    9: string Email
    10: i32 Province
    11: i32 City
    12: i32 District
    13: string Remark
    14: string Ext1
    15: string Ext2
    16: string Ext3
    17: string Ext4
    18: string Ext5
    19: string Ext6
    20: i64 UpdateTime
}

struct Account {
    1: i32 MemberId
    2: i64 Integral
    3: i64 FreezeIntegral
    4: double Balance
    5: double FreezeBalance
    6: double ExpiredBalance
    7: double PresentBalance
    8: double FreezeWallet
    9: double ExpiredPresent
    10: double TotalPresentFee
    11: double FlowBalance
    12: double GrowBalance
    13: double GrowAmount
    14: double GrowEarnings
    15: double GrowTotalEarnings
    16: double TotalConsumption
    17: double TotalCharge
    18: double TotalPay
    19: i64 PriorityPay
    20: i64 UpdateTime
}

struct MemberSummary {
    1: i32 MemberId
    2: string Usr
    3: string Name
    4: string Avatar
    5: i32 Exp
    6: i32 Level
    7: string LevelName
    8: string LevelSign
    9: i64 LevelOfficial
    10: string InvitationCode
    11: i64 Integral
    12: double Balance
    13: double PresentBalance
    14: double GrowBalance
    15: double GrowAmount
    16: double GrowEarnings
    17: double GrowTotalEarnings
    18: i64 UpdateTime
}

struct MemberRelation {
    1: i32 MemberId
    2: string CardId
    3: i32 InviterId
    4: string InviterStr
    5: i32 RegisterMchId
}

struct TrustedInfo {
    1: i32 MemberId
    2: string RealName
    3: string CardId
    4: string TrustImage
    5: i32 Reviewed
    6: i64 ReviewTime
    7: string Remark
    8: i64 UpdateTime
}


struct Address {
    1: i32 Id
    2: i32 MemberId
    3: string RealName
    4: string Phone
    5: i32 Province
    6: i32 City
    7: i32 District
    8: string Area
    9: string Address
    10: i32 IsDefault
}

//商品分类
struct Category {
    1: i32 ID
    2: i32 ParentId
    3: i32 ProModel
    4: string Name
    5: i32 Level
    6: string Icon
    7: string CatUrl
    8: i32 SortNum
    9: i32 FloorShow
    10: i32 Enabled
    11: i64 CreateTime
    12: list<Category> Children
}

struct Item {
    1: i32 ItemId
    2: i32 ProductId
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
    19: i32 SkuId
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
}

struct Sku {
    1: i32 SkuId
    2: i32 ItemId
    3: i32 ProductId
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

// 购物车
struct ShoppingCart {
    //编号
    1: i32 CartId
    //购物车KEY
    2: string Key
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
    1: i32 ItemId
    //SKU编号
    2: i32 SkuId
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

//支付单
struct PaymentOrder {
    1: i32 Id
    2: string TradeNo
    3: i32 VendorId
    4: i32 Type
    5: i32 OrderId
    6: string Subject
    7: i32 BuyUser
    8: i32 PaymentUser
    9: double TotalFee
    10: double BalanceDiscount
    11: double  IntegralDiscount
    12: double SystemDiscount
    13: double CouponDiscount
    14: double SubAmount
    15: double AdjustmentAmount
    16: double FinalAmount
    17: i32 PaymentOptFlag
    18: i32 PaymentSign
    19: string OuterNo
    20: i64 CreateTime
    21: i64 PaidTime
    22: i32 State
}

// 订单项
struct OrderItem {
    1: i32 Id
    2: i32 OrderId
    3: i32 ItemId
    4: i32 SkuId
    5: i32 SnapshotId
    6: i32 Quantity
    7: i32 ReturnQuantity
    8: double Amount
    9: double FinalAmount
    10: i64 IsShipped
    11: i64 UpdateTime
}

// 子订单
struct SubOrder {
    1: i32 Id
    2: string OrderNo
    3: i32 ParentId
    4: i32 BuyerId
    5: i32 VendorId
    6: i32 ShopId
    7: string Subject
    8: double GoodsAmount
    9: double DiscountAmount
    10: double ExpressFee
    11: double PackageFee
    12: double FinalAmount
    13: i64 IsPaid
    14: i64 IsSuspend
    15: string Note
    16: string Remark
    17: i64 CreateTime
    18: i64 UpdateTime
    19: i32 State
    20: list<OrderItem> Items
}



//会员服务
service MemberService{
    // 登录，返回结果(Result)和会员编号(Id);
    // Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
    Result Login(1:string user,2:string pwd,3:bool update),
    // 根据会员编号获取会员信息
    Member GetMember(1:i32 id),
    // 根据用户名获取会员信息
    Member GetMemberByUser(1:string user),
    // 根据会员编号获取会员资料
    Profile GetProfile(1:i32 id),
    // 获取会员汇总信息
    MemberSummary Summary(1:i32 memberId)
    // 获取会员的会员Token,reset表示是否重置token
    string GetToken(1:i32 memberId,2:bool reset)
    // 检查会员的会话Token是否正确，如正确返回: 1
    bool CheckToken(1:i32 memberId,2:string token)
    // 移除会员的Token
    void RemoveToken(1:i32 memberId)
    // 获取地址，如果addrId为0，则返回默认地址
    Address GetAddress(1:i32 memberId,2:i32 addrId)
    // 获取会员账户信息
    Account GetAccount(1:i32 memberId)
    // 获取邀请人会员编号数组
    list<i32> InviterArray(1:i32 memberId,2:i32 depth)
    // 账户充值
    Result ChargeAccount(1:i32 memberId ,2:i32 account,3:i32 kind,
      4:string title,5:string outerNo,6:double amount,7:i32 relateUser)
    // 抵扣账户
    Result DiscountAccount(1:i32 memberId,2:i32 account,3:string title,
      4:string outerNo,5:double amount,6:i32 relateUser,7:bool mustLargeZero)
}

struct PlatformConf {
    1: string Name
    2: string Logo
    3: string Telephone
    4: bool Suspend
    5: string SuspendMessage
    6: bool MchGoodsCategory
    7: bool MchPageCategory
}

// 单点登录应用
struct SsoApp{
    // 编号
    1: i32 Id
    // 应用名称
    2: string Name
    // API地址
    3: string ApiUrl
    // 密钥
    4: string Token
}

// 基础服务
service FoundationService{
   // 格式化资源地址并返回
   string ResourceUrl(1:string url)
   // 获取平台设置
   PlatformConf GetPlatformConf()
   // 根据键获取值
   string GetValue(1:string key)
   // 设置键值
   Result SetValue(1:string key,2:string value)
   // 删除值
   Result DeleteValue(1:string key)
   // 根据前缀获取值
   map<string,string> GetValuesByPrefix(1:string prefix)
   // 注册单点登录应用,返回值：
   //   -  1. 成功，并返回token
   //   - -1. 接口地址不正确
   //   - -2. 已经注册
   string RegisterApp(1:SsoApp app)
   // 获取应用信息
   SsoApp GetApp(1:string name)
   // 获取单点登录应用
   list<string> GetAllSsoApp()
   // 验证超级用户账号和密码
   bool ValidateSuper(1:string user,2:string pwd)
   // 保存超级用户账号和密码
   void FlushSuperPwd(1:string user,2:string pwd)
   // 创建同步登录的地址
   string GetSyncLoginUrl(1:string returnUrl)
}

// 支付服务
service PaymentService{
    // 创建支付单
    Result CreatePaymentOrder(1:PaymentOrder o)
    // 根据支付单号获取支付单
    PaymentOrder GetPaymentOrder(1:string paymentNo)
    // 根据编号获取支付单
    PaymentOrder GetPaymentOrderById(1:i32 id)
    // 调整支付单金额
    Result AdjustOrder(1:string paymentNo, 2:double amount)
    // 余额抵扣
    Result DiscountByBalance(1:i32 orderId,2:string remark )
   // 积分抵扣支付单
    DResult DiscountByIntegral(1:i32 orderId,2:i32 integral,3:bool ignoreOut)
    // 钱包账户支付
    Result PaymentByWallet(1:i32 orderId,2:string remark)
    // 完成支付单支付，并传入支付方式及外部订单号
    Result FinishPayment(1:string tradeNo ,2:string spName,3:string outerNo)
}

// 销售服务
service SaleService {
  // 获取子订单
  SubOrder GetSubOrder(1:i32 id)
  // 根据订单号获取子订单
  SubOrder GetSubOrderByNo(1:string orderNo)
  // 获取订单商品项
  list<OrderItem> GetSubOrderItems(1:i32 subOrderId)
}