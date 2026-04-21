/**
 * Copyright 2015 @ 56x.net.
 * name : queue
 * author : jarryliu
 * date : 2015-07-27 17:04
 * description : Message queue key constants
 * history :
 */
package constants

// Message queue keys
const (
	QueuePaymentOrderFinish     = "go2o:mq:payment_success_notify" // Payment order completion notification queue
	QueueOrderBusiness          = "go2o:mq:order_notify"           // Order business queue (created, completed, etc. executed once)
	QueueNewMailTask            = "go2o:mq:mail"                   // New mail task queue
	QueueMemberUpdateTime       = "go2o:mm:uptime_"                // Member update time
	QueueAccountUpdateTime      = "go2o:acc:uptime_"               // Account update time
	QueueMemberUpdateTcpNotify  = "go2o:mm:queue:t_up_notify"      // Member TCP update notification queue
	QueueAccountUpdateTcpNotify = "go2o:mq:acc_tcp_notify"         // Account TCP update notification queue
	QueueMemberUpdate           = "go2o:mq:mm_update"              // New member queue
	QueueOrderExpiresTime       = "go2o:order:timeout"             // Order expiration time
	QueueMemberAutoUnlock       = "go2o:order:unlock"              // Member auto unlock
	QueueOrderAutoReceive       = "go2o:order:auto_receive"        // Order auto receive
)

// Advertisement keys
const (
	AdKeyInvitationDM = "UC_INVATION_DM" // User invitation DM page image ad
)
