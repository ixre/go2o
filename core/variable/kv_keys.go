/**
 * Copyright 2015 @ z3q.net.
 * name : kv_keys
 * author : jarryliu
 * date : 2015-07-27 17:04
 * description :
 * history :
 */
package variable

const (
	KvNewMailTask                 = "go2o:q:mail"
	KvTotalMembers                = "go2o:sys:total_members"
	KvMemberUpdateTime            = "go2o:mm:uptime_"
	KvAccountUpdateTime           = "go2o:acc:uptime_"
	KvMemberUpdateTcpNotifyQueue  = "go2o:mm:queue:t_up_notify"
	KvAccountUpdateTcpNotifyQueue = "go2o:q:acc_tcp_notify"  //账户TCP更新对列
	KvMemberUpdateQueue           = "go2o:q:mm_update"       //新加入会员队列
	KvPaymentOrderFinishQueue     = "go2o:q:pay_order"       //支付单完成通知队列
	KvOrderBusinessQueue          = "go2o:q:sa_order_busi"   //订单业务队列(如已创建,已完成等只执行一次)
	KvOrderExpiresTime            = "go2o:order:timeout"     //订单过期时间
	KvOrderAutoReceive            = "go2o:order:autoreceive" //订单自动收货
)

const (
	//用户推荐DM页图片广告
	AdKeyInvitationDM = "UC_INVATION_DM"
)
