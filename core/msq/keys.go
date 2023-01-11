package msq

const (
	// MemberUpdated 会员创建或更新, 消息: create|1
	MemberUpdated = "go2o-member-updated"
	// MemberTrustInfoPassed 会员实名信息通过, 消息: 会员ID|证件类型|证件号码|姓名
	MemberTrustInfoPassed = "go2o-member-trusted-info-passed"
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
	MemberRequestWithdrawal = "go2o-member-request-withdrawal"
	// MemberWithdrawalAudited 用户提现申请已通过
	MemberWithdrawalAudited = "go2o-member-withdrawal-audited"
	// SystemSendSmsTopic 发送短信
	SystemSendSmsTopic = "go2o-system-send-sms-topic"
)
