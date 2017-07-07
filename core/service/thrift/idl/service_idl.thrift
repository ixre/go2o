namespace go define

//传输结果对象
struct Result{
   1:i32 ID
   2:bool Result
   3:string Code
   4:string Message
}
//传输结果对象
struct Result64{
   1:i64 ID
   2:bool Result
   3:string Code
   4:string Message
}
//传输结果对象(Double)
struct DResult{
   1:double Data
   2:bool Result
   3:string Code
   4:string Message
}

// 键值对
struct Pair{
   1:string Key
   2:string Value
}

// 商家
struct ComplexMerchant {
    1: i32 ID
    2: i64 MemberId
    3: string Usr
    4: string Pwd
    5: string Name
    6: i32 SelfSales
    7: i32 Level
    8: string Logo
    9: string CompanyName
    10: i32 Province
    11: i32 City
    12: i32 District
    13: i32 Enabled
    14: i64 ExpiresTime
    15: i64 JoinTime
    16: i64 UpdateTime
    17: i64 LoginTime
    18: i64 LastLoginTime
}


// 商店,需重构
struct Shop{
    1:i32 ID
    2:i32 VendorId
    3:i32 ShopType
    4:string Name
    5:i32 State
    6:i32 OpeningState
    7:map<string,string> Data
}

// 商铺
struct Store{
    1:i32 ID
    2:i32 VendorId
    3:string Name
    4:string Alias
    5:string Host
    6:string Logo
    7:i32 State
    8:i32 OpeningState
    9:string StorePhone
    10:string StoreTitle
    11:string StoreNotice
}

struct Level {
    1: i32 ID
    2: string Name
    3: i32 RequireExp
    4: string ProgramSignal
    5: i32 IsOfficial
    6: i32 Enabled
}

struct Member {
    1: i64 ID
    2: string Usr
    3: string Pwd
    4: string TradePwd
    5: i32 Exp
    6: i32 Level
    7: string InvitationCode
    // 高级用户类型
    8:i32   PremiumUser
    // 高级用户过期时间
    9:i64   PremiumExpires
    10: string RegFrom
    11: string RegIp
    12: i64 RegTime
    13: string CheckCode
    14: i64 CheckExpires
    15: i32 State
    16: i64 LoginTime
    17: i64 LastLoginTime
    18: i64 UpdateTime
    19: string DynamicToken
    20: i64 TimeoutTime
}

struct Profile {
    1: i64 MemberId
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
    1: i64 MemberId
    2: i64 Integral
    3: i64 FreezeIntegral
    4: double Balance
    5: double FreezeBalance
    6: double ExpiredBalance
    7: double WalletBalance
    8: double FreezeWallet
    9: double ExpiredPresent
    10: double TotalPresentFee
    11: double FlowBalance
    12: double GrowBalance
    13: double GrowAmount
    14: double GrowEarnings
    15: double GrowTotalEarnings
    16: double TotalExpense
    17: double TotalCharge
    18: double TotalPay
    19: i64 PriorityPay
    20: i64 UpdateTime
}

struct ComplexMember {
    1: i64 MemberId
    2: string Usr
    3: string Name
    4: string Avatar
    5: i32 Exp
    6: i32 Level
    7: string LevelName
    8: string LevelSign
    9: i32 LevelOfficial
    10: i32	PremiumUser
    11: i64	PremiumExpires
    12: string InvitationCode
    13: i32 TrustAuthState
    14: i32 State
    15: i64 Integral
    16: double Balance
    17: double WalletBalance
    18: double GrowBalance
    19: double GrowAmount
    20: double GrowEarnings
    21: double GrowTotalEarnings
    22: i64 UpdateTime
}

struct MemberRelation {
    1: i64 MemberId
    2: string CardId
    3: i64 InviterId
    4: string InviterStr
    5: i32 RegisterMchId
}

struct TrustedInfo {
    1: i64 MemberId
    2: string RealName
    3: string CardId
    4: string TrustImage
    5: i32 Reviewed
    6: i64 ReviewTime
    7: string Remark
    8: i64 UpdateTime
}


struct Address {
    1: i64 ID
    2: i64 MemberId
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
    // 分类优先级
    4: i32 Priority
    // 分类名称
    5: string Name
    // 虚拟分类
    6: i32 VirtualCat
    // 分类目标地址
    7: string CatUrl
    8: string Icon
    9: i32 SortNum
    10: i32 FloorShow
    11: i32 Enabled
    12: i32 Level
    13: i64 CreateTime
    14: list<Category> Children
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
    33: map<string,string> Data;
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
    1: i32 ID
    2: string TradeNo
    3: i32 VendorId
    4: i32 Type
    5: i32 OrderId
    6: string Subject
    7: i64 BuyUser
    8: i64 PaymentUser
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
struct ComplexItem {
    1: i64 ID
    2: i64 OrderId
    3: i64 ItemId
    4: i64 SkuId
    5: i64 SnapshotId
    6: i32 Quantity
    7: i32 ReturnQuantity
    8: double Amount
    9: double FinalAmount
    10: i32 IsShipped
    11: map<string,string> Data
}

// 子订单
struct ComplexOrder {
    1: i64 OrderId
    2: i64 SubOrderId
    3: i32 OrderType
    4: string OrderNo
    5: i64 BuyerId
    6: i32 VendorId
    7: i32 ShopId
    8: string Subject
    9: double ItemAmount
    10: double DiscountAmount
    11: double ExpressFee
    12: double PackageFee
    13: double FinalAmount
    14: string ConsigneePerson
    15: string ConsigneePhone
    16: string ShippingAddress
    17: string BuyerRemark
    18: i32 IsBreak
    19: i32 State
    20: i64 CreateTime
    21: i64 UpdateTime
    22: list<ComplexItem> Items
    // 扩展信息
    23: map<string,string> Data
}


struct PlatformConf {
   // 1: string Name
   // 2: string Logo
   // 3: string Telephone
    1: bool Suspend
    2: string SuspendMessage
    3: bool MchGoodsCategory
    4: bool MchPageCategory
}

// 单点登录应用
struct SsoApp{
    // 编号
    1: i32 ID
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
   // 获取键值存储数据
   list<string> GetRegistryV1(1:list<string> keys)
   // 获取键值存储数据字典
   map<string,string> GetRegistryMapV1(1:list<string> keys)
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
   bool SuperValidate(1:string user,2:string pwd)
   // 保存超级用户账号和密码
   void FlushSuperPwd(1:string user,2:string pwd)
   // 创建同步登录的地址
   string GetSyncLoginUrl(1:string returnUrl)
}



//会员服务
service MemberService{
    // 登录，返回结果(Result)和会员编号(Id);
    // Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
    Result64 CheckLogin(1:string user,2:string pwd,3:bool update)
    // 检查交易密码
    Result CheckTradePwd(1:i64 id,2:string tradePwd)
    // 等级列表
    list<Level> LevelList()
    // 获取实名信息
    TrustedInfo GetTrustInfo(1:i64 id)
    // 获取等级信息
    Level GetLevel(1:i32 id)
    // 根据会员编号获取会员信息
    Member GetMember(1:i64 id)
    // 根据用户名获取会员信息
    Member GetMemberByUser(1:string user)
    // 根据会员编号获取会员资料
    Profile GetProfile(1:i64 id)
    // 获取会员汇总信息
    ComplexMember Complex(1:i64 memberId)
     // 升级为高级会员
    Result Premium(1:i64 memberId,2:i32 v,3:i64 expires)
    // 获取会员的会员Token,reset表示是否重置token
    string GetToken(1:i64 memberId,2:bool reset)
    // 检查会员的会话Token是否正确，如正确返回: 1
    bool CheckToken(1:i64 memberId,2:string token)
    // 移除会员的Token
    void RemoveToken(1:i64 memberId)
    // 获取地址，如果addrId为0，则返回默认地址
    Address GetAddress(1:i64 memberId,2:i64 addrId)
    // 获取会员账户信息
    Account GetAccount(1:i64 memberId)
    // 获取邀请人会员编号数组
    list<i64> InviterArray(1:i64 memberId,2:i32 depth)
    // 获取从指定时间到现在推荐指定等级会员的数量
    i32 GetInviterQuantity(1:i64 memberId,2:map<string,string> data)

    // 账户充值
    Result ChargeAccount(1:i64 memberId ,2:i32 account,3:i32 kind,
      4:string title,5:string outerNo,6:double amount,7:i64 relateUser)
    // 抵扣账户
    Result DiscountAccount(1:i64 memberId,2:i32 account,3:string title,
      4:string outerNo,5:double amount,6:i64 relateUser,7:bool mustLargeZero)
}

// 支付服务
service PaymentService{
    // 创建支付单并提交
    Result SubmitPaymentOrder(1:PaymentOrder o)
    // 根据支付单号获取支付单
    PaymentOrder GetPaymentOrder(1:string paymentNo)
    // 根据交易号获取支付单编号
    i32 GetPaymentOrderId(1:string tradeNo)
    // 根据编号获取支付单
    PaymentOrder GetPaymentOrderById(1:i32 id)
    // 调整支付单金额
    Result AdjustOrder(1:string paymentNo, 2:double amount)
    // 余额抵扣
    Result DiscountByBalance(1:i32 orderId,2:string remark )
   // 积分抵扣支付单
    DResult DiscountByIntegral(1:i32 orderId,2:i64 integral,3:bool ignoreOut)
    // 钱包账户支付
    Result PaymentByWallet(1:i32 orderId,2:string remark)
    // 余额钱包混合支付，优先扣除余额。
    Result HybridPayment(1:i32 orderId,2:string remark)
    // 完成支付单支付，并传入支付方式及外部订单号
    Result FinishPayment(1:string tradeNo ,2:string spName,3:string outerNo)
}



// 销售服务
service SaleService {
  // 批发购物车接口
  Result WholesaleCartV1(1:i64 memberId,2:string action,3:map<string,string> data)
  // 提交订单
  map<string,string> SubmitOrderV1(1:i64 buyerId,2:i32 cartType,3:map<string,string> data)
  // 获取订单信息
  ComplexOrder GetOrder(1:string order_no,2:bool sub_order)
  // 获取订单和商品项信息
  ComplexOrder GetOrderAndItems(1:string order_no,2:bool sub_order)

  // 获取子订单
  ComplexOrder GetSubOrder(1:i64 id)
  // 根据订单号获取子订单
  ComplexOrder GetSubOrderByNo(1:string orderNo)
  // 获取订单商品项
  list<ComplexItem> GetSubOrderItems(1:i64 subOrderId)
  // 提交交易订单
  Result64 SubmitTradeOrder(1:ComplexOrder o,2:double rate)
  // 交易单现金支付
  Result64 TradeOrderCashPay(1:i64 orderId)
  // 上传交易单发票
  Result64 TradeOrderUpdateTicket(1:i64 orderId,2:string img)
}

// 商品服务
service ItemService{
    // 获取SKU
    Sku GetSku(1:i32 itemId,2:i32 skuId)
    // 获取商品的Sku-JSON格式
    string GetItemSkuJson(1:i32 itemId)
    // 获取商品详细数据
    string GetItemDetailData(1:i32 itemId,2:i32 iType)
}

//商家服务
service MerchantService{
   // 获取商家符合的信息
   ComplexMerchant Complex(1:i32 mchId)
   // 验证用户密码,并返回编号。可传入商户或会员的账号密码
   Result CheckLogin(1:string usr,2:string oriPwd)
   // 验证商户状态
   Result Stat(1:i32 mchId)
   // 同步批发商品
   map<string,i32> SyncWholesaleItem(1:i32 mchId)
}

// 商店服务
service ShopService{
    // 获取店铺
    Store GetStore(1:i32 venderId)
    // 获取店铺
    Store GetStoreById(1:i32 shopId)
    // 获取门店
    //Shop GetOfflineShop(1:i32 shopId)
    // 打开或关闭商店
    Result TurnShop(1:i32 shopId,2:bool on,3:string reason)
    // 设置商店是否营业
    Result OpenShop(1:i32 shopId,2:bool opening,3:string reason)
}