namespace go define

//传输结果对象
struct Result{
   1:i32 Id
   2:bool Result
   3:i32 Code
   4:string Message
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

//会员服务
service MemberService{
    // 登陆，返回结果(Result)和会员编号(Id);
    // Result值为：-1:会员不存在; -2:账号密码不正确; -3:账号被停用
    Result Login(1:string user,2:string pwd,3:bool update),
    // 根据会员编号获取会员信息
    Member GetMember(1:i32 id),
    // 根据用户名获取会员信息
    Member GetMemberByUser(1:string user),
    // 根据会员编号获取会员资料
    Profile GetProfile(1:i32 id),
    // 获取会员的会员Token,reset表示是否重置token
    string GetToken(1:i32 memberId,2:bool reset)
    // 检查会员的会话Token是否正确，如正确返回: 1
    bool CheckToken(1:i32 memberId,2:string token)
    // 移除会员的Token
    void RemoveToken(1:i32 memberId)
    // 获取地址，如果addrId为0，则返回默认地址
    Address GetAddress(1:i32 memberId,2:i32 addrId)
}

struct PlatformConf {
    1: string Name
    2: string Logo
    3: bool Suspend
    4: string SuspendMessage
    5: bool MchGoodsCategory
    6: bool MchPageCategory
}

// 单点登陆应用
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

   // 注册单点登陆应用,返回值：
   //   -  1. 成功，并返回token
   //   - -1. 接口地址不正确
   //   - -2. 已经注册
   string RegisterSsoApp(1:SsoApp app)
   // 获取单点登陆应用
   list<string> GetAllSsoApp()
   // 验证超级用户账号和密码
   bool ValidateSuper(1:string user,2:string pwd)
   // 保存超级用户账号和密码
   void FlushSuperPwd(1:string user,2:string pwd)
   // 创建同步登陆的地址
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
}
