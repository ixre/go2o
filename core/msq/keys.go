package msq

const (
	// 会员创建或更新, 消息: create|1
	MemberUpdated = "go2o-member-updated"
	// 会员实名信息通过, 消息: 会员ID|证件类型|证件号码|姓名
	MemberTrustInfoPassed = "go2o-member-trusted-info-passed"
	// 会员账户更新, 消息: 1
	MemberAccountUpdated = "go2o-member-account-updated"
	// 会员资料更新, 消息: 1
	MemberProfileUpdated = "go2o-member-profile-updated"
	// 会员关系更新, 消息: 1
	MemberRelationUpdated = "go2o-member-relation-updated"
)
