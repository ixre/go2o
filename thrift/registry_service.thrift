namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.registry_service
include "ttype.thrift"


// 基础服务
service RegistryService{
   /** 获取注册表键值 */
   string GetRegistry(1:string key)
   /** 获取键值存储数据字典 */
   map<string,string> GetRegistries(1:list<string> keys)
   /** 按键前缀获取键数据 */
   map<string,string> FindRegistries(1:string prefix)
   /** 更新注册表键值 */
   ttype.Result UpdateRegistry(1:map<string,string> registries)
   /** 搜索注册表 */
   list<SRegistry> SearchRegistry(1:string key)
   /** 创建自定义注册表项,@defaultValue 默认值,如需更改,使用UpdateRegistry方法  */
   ttype.Result CreateUserRegistry(1:string key,2:string defaultValue,3:string description)
   /** 获取键值存储数据 */
   list<string> GetRegistryV1(1:list<string> keys)
}

/** 注册表 */
struct SRegistry {
    /** 键 */
    1: string Key
    /** 值 */
    2: string Value
    /** 默认值 */
    3: string Default
    /** 可选值 */
    4: string Options
    /** 标志 */
    5: i32 Flag
    /** 描述 */
    6: string Description
}