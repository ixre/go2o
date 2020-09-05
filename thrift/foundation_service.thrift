namespace java com.github.jsix.go2o.rpc
namespace netstd com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.foundation_service
include "ttype.thrift"



// 基础服务
service FoundationService{
   /** 获取短信API凭据, @provider 短信服务商, 默认:http */
   SSmsApi GetSmsApi(1:string provider)
   /** 保存短信API凭据,@provider 短信服务商, 默认:http */
   ttype.Result SaveSmsApi(1:string provider,2:SSmsApi api)
   /** 保存面板HOOK数据,这通常是在第三方应用中初始化或调用,参见文档：BoardHooks */
   ttype.Result SaveBoardHook(1:string hookURL,2:string token)

   // 格式化资源地址并返回
   string ResourceUrl(1:string url)
   // 设置键值
   ttype.Result SetValue(1:string key,2:string value)
   // 删除值
   ttype.Result DeleteValue(1:string key)
   // 根据前缀获取值
   map<string,string> GetValuesByPrefix(1:string prefix)
   // 注册单点登录应用,返回值：
   //   -  1. 成功，并返回token
   //   - -1. 接口地址不正确
   //   - -2. 已经注册
   string RegisterApp(1:SSsoApp app)
   // 获取应用信息
   SSsoApp GetApp(1:string name)
   // 获取单点登录应用
   list<string> GetAllSsoApp()
   // 验证超级用户账号和密码
   bool SuperValidate(1:string user,2:string pwd)
   // 保存超级用户账号和密码
   void FlushSuperPwd(1:string user,2:string pwd)
   // 创建同步登录的地址
   string GetSyncLoginUrl(1:string returnUrl)
   // 获取地区名称
   list<string> GetAreaNames(1:list<i32> codes)
   // 获取下级区域
   list<SArea> GetChildAreas(1:i32 code)
}


// 单点登录应用
struct SSsoApp{
    // 编号
    1: i32 ID
    // 应用名称
    2: string Name
    // API地址
    3: string ApiUrl
    // 密钥
    4: string Token
}

/** 行政区域 */
struct SArea  {
    1:i32 Code
    2:i32 Parent
    3:string Name
}

/** 短信接口信息 */
struct SSmsApi {
    /** 接口地址 */
    1:string ApiUrl
    /** 接口KEY */
    2:string Key
    /** 接口密钥 */
    3:string Secret
	/** 请求数据,如: phone={phone}&content={content}*/
	4:string Params
	/** 请求方式, GET或POST */
	5:string Method
    /** 编码 */
    6:string Charset
    /** 签名 */
   	7:string Signature
   	/** 发送成功，包含的字符，用于检测是否发送成功 */
   	8:string SuccessChar
}
