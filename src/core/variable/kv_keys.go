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
	KvNewMailTask                 = "go2o:sys:mss:"
	KvHaveNewCreatedOrder         = "go2o:sys:queue:have_new_create_order"
	KvHaveNewCompletedOrder       = "go2o:sys:queue:have_new_completed_order"
	KvHaveNewMember               = "go2o:sys:queue:have_new_member"
	KvTotalMembers                = "go2o:sys:total_members"
	KvMemberUpdateTime            = "go2o:mm:uptime_"
	KvAccountUpdateTime           = "go2o:acc:uptime_"
	KvMemberUpdateTcpNotifyQueue  = "go2o:mm:queue:t_up_notify"
	KvAccountUpdateTcpNotifyQueue = "go2o:q:acc_tcp_notify" //账户TCP更新对列
	KvMemberUpdateQueue           = "go2o:q:mm_update"      //新加入会员队列
	KvOrderCreatedQueue           = "go2o:sa:q:order_new"   //新订单队列
	KvOrderBusinessQueue          = "go2o:q:sa_order_busi"  //订单业务队列(如已创建,已完成等只执行一次)
	KvOrderExpiresTime            = "go2o:o:expires:"       //订单过期时间
)
