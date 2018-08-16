namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.item_service
include "ttype.thrift"


// 商品服务
service ItemService{
    // 获取SKU
    ttype.SSku GetSku(1:i64 itemId,2:i64 skuId)
    // 获取商品的Sku-JSON格式
    string GetItemSkuJson(1:i64 itemId)
    // 获取商品详细数据
    string GetItemDetailData(1:i64 itemId,2:i32 iType)
}

