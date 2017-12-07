namespace go define
include "ttype.thrift"

// 发货服务
service ShipmentService{
    // 物流追踪
    SShipOrderTrace GetLogisticFlowTrace(1:string shipperCode,2:string logisticCode)
}

// 发货单追踪
struct SShipOrderTrace {
    // 返回状态码
    1:i32 Code
    // 返回错误信息
    2:string Message
    // 物流单号
    3:string LogisticCode
    // 承运商代码
    4:string ShipperCode
    // 发货状态
    5:string ShipState
    // 更新时间
    6:i64 UpdateTime
    // 包含发货单流
    7:list<SShipFlow> Flows
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