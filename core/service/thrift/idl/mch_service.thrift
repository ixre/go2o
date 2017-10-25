namespace go define

include "type.thrift"

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


//商家服务
service MerchantService{
   // 获取商家符合的信息
   ComplexMerchant Complex(1:i32 mchId)
   // 验证用户密码,并返回编号。可传入商户或会员的账号密码
   type.Result CheckLogin(1:string usr,2:string oriPwd)
   // 验证商户状态
   type.Result Stat(1:i32 mchId)
   // 同步批发商品
   map<string,i32> SyncWholesaleItem(1:i32 mchId)
}