/**
 * Copyright 2015 @ to2.net.
 * name : kv_keys
 * author : jarryliu
 * date : 2015-07-27 17:04
 * description :
 * history :
 */
package variable

const (
	KvPaymentOrderFinishQueue     = "go2o:mq:payment_success_notify" //支付单完成通知队列
	KvOrderBusinessQueue          = "go2o:mq:order_notify"           //订单业务队列(如已创建,已完成等只执行一次)
	KvNewMailTask                 = "go2o:mq:mail"
	KvTotalMembers                = "go2o:sys:total_members"
	KvMemberUpdateTime            = "go2o:mm:uptime_"
	KvAccountUpdateTime           = "go2o:acc:uptime_"
	KvMemberUpdateTcpNotifyQueue  = "go2o:mm:queue:t_up_notify"
	KvAccountUpdateTcpNotifyQueue = "go2o:mq:acc_tcp_notify"  //账户TCP更新对列
	KvMemberUpdateQueue           = "go2o:mq:mm_update"       //新加入会员队列
	KvOrderExpiresTime            = "go2o:order:timeout"      //订单过期时间
	KvOrderAutoReceive            = "go2o:order:auto_receive" //订单自动收货
)

const (
	//用户推荐DM页图片广告
	AdKeyInvitationDM = "UC_INVATION_DM"
)
