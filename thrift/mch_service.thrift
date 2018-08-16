namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.mch_service
include "ttype.thrift"

// 商家
struct SComplexMerchant {
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


//商家服务
service MerchantService{
   // 获取商家符合的信息
   SComplexMerchant Complex(1:i32 mchId)
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