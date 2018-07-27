namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.shipment_service
include "ttype.thrift"

// 发货服务
service ShipmentService{
    /** 物流追踪 */
    SShipOrderTrack GetLogisticFlowTrack(1:string shipperCode,2:string logisticCode,3:bool invert)
    /** 获取发货单的物流追踪信息,$shipOrderId:发货单编号 $invert:是否倒序排列 */
    SShipOrderTrack ShipOrderLogisticTrack(1:i64 shipOrderId,2:bool invert)
}

// 发货单追踪
struct SShipOrderTrack {
    // 返回状态码
    1:i32 Code
    // 返回错误信息
    2:string Message
    // 物流单号
    3:string LogisticCode
    // 承运商代码
    4:string ShipperCode
    /** 承运商名称 */
    5:string ShipperName
    // 发货状态
    6:string ShipState
    // 更新时间
    7:i64 UpdateTime
    // 包含发货单流
    8:list<SShipFlow> Flows
}
// 发货流
struct SShipFlow  {
    // 记录标题
    1:string Subject
    // 记录时间
    2:string CreateTime
    // 备注
    3:string Remark
}