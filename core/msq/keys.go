package msq

const (
	// MemberUpdated 会员创建或更新, 消息: create|1
	MemberUpdated = "go2o-member-updated"
	// MemberAccountUpdated 会员账户更新, 消息: 1
	MemberAccountUpdated = "go2o-member-account-updated"
	// MemberProfileUpdated 会员资料更新, 消息: 1
	MemberProfileUpdated = "go2o-member-profile-updated"
	// 普通订单状态变更
	ORDER_NormalOrderStatusChange = "go2o-normal-order-status-change"
	// OrderAffiliateTopic 订单分销
	OrderAffiliateTopic = "go2o-order-affiliate-topic"
	// WalletLogTopic 会员钱包日志订阅
	WalletLogTopic = "go2o-wallet-log-topic"
	// RegistryTopic 自定义键值变更订阅
	RegistryTopic = "go2o-registry-topic"
	// MemberRequestWithdrawal 用户发起提现申请
	MembertWithdrawalTopic = "go2o-member-withdrawal-topic"
	// SendSmsTopic 发送短信
	SendSmsTopic = "go2o-send-sms-topic"
)
