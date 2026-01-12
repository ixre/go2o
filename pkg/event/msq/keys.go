package msq

const (
	// MemberUpdated 会员创建或更新
	MemberUpdated = "go2o-member-updated"
	// MemberAccountUpdated 会员账户更新
	MemberAccountUpdated = "go2o-member-account-updated"
	// MemberAccountLogTopic 会员账户日志订阅
	MemberAccountLogTopic = "go2o-member-account-log-topic"
	// MemberWithdrawalTopic 用户提现
	MembertWithdrawalTopic = "go2o-member-withdrawal-topic"
	// NormalOrderStatusTopic 普通订单状态变更
	NormalOrderStatusTopic = "go2o-normal-order-status-topic"
	// NormalOrderAffiliateTopic 订单分销
	NormalOrderAffiliateTopic = "go2o-normal-order-affiliate-topic"
	// RegistryTopic 自定义键值变更订阅
	RegistryTopic = "go2o-registry-topic"
	// SendSmsTopic 发送短信
	SendSmsTopic = "go2o-send-sms-topic"
)
