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
	// MemberRelationUpdated 会员关系更新, 消息: 1
	MemberRelationUpdated = "go2o-member-relation-updated"
	// WalletLogTopic 会员钱包日志订阅
	WalletLogTopic = "go2o-wallet-log-topic"
	// RegistryTopic 自定义键值变更订阅
	RegistryTopic = "go2o-registry-topic"
)
