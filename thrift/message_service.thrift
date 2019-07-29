namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.message_service
include "ttype.thrift"

/** 消息服务 */
service MessageService{
    /** 获取通知项 */
    SNotifyItem GetNotifyItem(1:string key)
    /** 发送短信 */
    ttype.Result SendPhoneMessage(1:string phone,2:string message,3:map<string,string> data)
}

/** 消息方式 */
enum EMessageChannel{
    /** 站内信 */
    SiteMessage = 1,
    /** 邮件 */
	EmailMessage = 2,
	/** 短信 */
	SmsMessage = 3,
}

/** 通知项 */
struct SNotifyItem {
    /** 键 */
    1: string Key
    /** 发送方式 */
    2: i32 NotifyBy
    /** 不允许修改发送方式 */
    3: bool ReadonlyBy
    /** 模板编号 */
    4: i32 TplId
    /** 内容 */
    5: string Content
    /** 模板包含的标签 */
    6: map<string,string> Tags
}