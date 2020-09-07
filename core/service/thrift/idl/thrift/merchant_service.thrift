namespace java com.github.jsix.go2o.rpc
namespace netstd com.github.jsix.go2o.rpc
namespace go go2o.core.service.thrift.auto_gen.rpc.merchant_service
include "ttype.thrift"


// 商家
struct SMerchant {
    /** 编号 */
    1:i32 Id
    /** 会员编号 */
    2:i64 MemberId
    /** 登录用户 */
    3:string LoginUser
    /** 登录密码 */
    4:string LoginPwd
    /** 名称 */
    5:string Name
    /** 公司名称 */
    6:string CompanyName
    /** 是否字营 */
    7:i16 SelfSales
    /** 商户等级 */
    8:i32 Level
    /** 标志 */
    9:string Logo
    /** 省 */
    10:i32 Province
    /** 市 */
    11:i32 City
    /** 区 */
    12:i32 District
    /** 标志 */
    13:i32 Flag
    /** 是否启用 */
    14:i16 Enabled
    /** 最后登录时间 */
    15:i32 LastLoginTime
}


// 商家
struct SMerchantPack {
    /** 登录用户 */
    1:string LoginUser
    /** 登录密码 */
    2:string LoginPwd
    /** 名称 */
    3:string Name
    /** 是否字营 */
    4:i16 SelfSales
    /** 店铺名称 */
    5:string ShopName
    /** 标志 */
    6:string ShopLogo
    /** 电话 */
    7:string Tel
    /** 地址 */
    8:string Addr
}


//商家服务
service MerchantService{
   // 获取商家符合的信息
   SMerchant GetMerchant(1:i32 mchId)
   // 注册商户并开店
   ttype.Result CreateMerchant(1:SMerchantPack mch,2:i64 relMemberId)

   // 验证用户密码,并返回编号。可传入商户或会员的账号密码
   ttype.Result CheckLogin(1:string usr,2:string oriPwd)
   // 验证商户状态
   ttype.Result Stat(1:i32 mchId)
   // 同步批发商品
   map<string,i32> SyncWholesaleItem(1:i32 mchId)
   // 获取所有的交易设置
   list<STradeConf> GetAllTradeConf(1:i32 mchId)
   // 根据交易类型获取交易设置
   STradeConf GetTradeConf(1:i32 mchId,2:i32 tradeType)
   // 保存交易设置
   ttype.Result SaveTradeConf(1:i32 mchId,2:list<STradeConf> arr)
}

// 商户交易设置
struct STradeConf  {
	// 商户编号
	1:i64 MchId
	// 交易类型
	2:i32 TradeType
	// 交易方案，根据方案来自动调整比例
	3:i64 PlanId
	// 交易标志
	4:i32 Flag
	// 交易手续费依据,1:按金额 2:按比例
	5:i32 AmountBasis
	// 交易费，按单笔收取
	6:i32 TradeFee
	// 交易手续费比例
	7:i32 TradeRate
}