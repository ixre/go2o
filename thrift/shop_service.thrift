namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.shop_service
include "ttype.thrift"

// 商店,需重构
struct SShop2{
    1:i32 ID
    2:i32 VendorId
    3:i32 ShopType
    4:string Name
    5:i32 State
    6:i32 OpeningState
    7:map<string,string> Data
}


// 店铺
struct SShop{
  /** 店铺编号 */
  1:i32 Id
  /** 商户编号 */
  2:i32 VendorId
  /** 店铺名称 */
  3:string ShopName
  /** 店铺标志 */
  4:string Logo
  /** 自定义 域名 */
  5:string Host
  /** 个性化域名 */
  6:string Alias
  /** 电话 */
  7:string Tel
  /** 地址 */
  8:string Addr
  /** 店铺标题 */
  9:string ShopTitle
  /** 店铺公告 */
  10:string ShopNotice
  /** 标志 */
  11:i32 Flag
  /** 状态 */
  12:i16 State
}

// 店铺
struct SStore{
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

// 商店服务
service ShopService{
    /** 查询店铺ID,根据Host */
    i32 QueryShopIdByHost(1:string host)
    /** 获取店铺 */
    SShop GetShop(1:i32 shopId)

    // 获取店铺
    SStore GetStore(1:i32 venderId)
    // 获取店铺
    SStore GetStoreById(1:i32 shopId)
    /** 根据主机头获取店铺编号 */
    i32 QueryStoreByHost(1:string host)
    // 获取门店
    //Shop GetOfflineShop(1:i32 shopId)
    // 打开或关闭商店
    ttype.Result TurnShop(1:i32 shopId,2:bool on,3:string reason)
    // 设置商店是否营业
    ttype.Result OpenShop(1:i32 shopId,2:bool opening,3:string reason)
}