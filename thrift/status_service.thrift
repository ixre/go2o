namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.status_service
include "ttype.thrift"

/** 状态服务 */
service StatusService{
    /** 尝试连接 */
    string Ping()
}