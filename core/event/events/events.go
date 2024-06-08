package events

import (
	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/order"
)

// 应用初始化事件
type AppInitialEvent struct {
}

// 账户日志推送事件
type AccountLogPushEvent struct {
	// 账户类型
	Account int
	// 是否为更新日志事件
	IsUpdateEvent bool
	// 会员编号
	MemberId int
	// 编号
	LogId int
	// 业务类型
	LogKind int
	// 标题
	Subject string
	// 外部订单号
	OuterNo string
	// 变动金额
	ChangeValue int
	// 余额
	Balance int
	// 交易手续费
	ProcedureFee int
	// 审核状态
	ReviewStatus int
	// 创建时间
	CreateTime int
}

// 订单分销事件
type OrderAffiliateRebateEvent struct {
	// 订单号
	OrderNo string
	// 子订单
	SubOrder bool
	// 买家编号
	BuyerId int64
	// 订单金额
	OrderAmount int64
	// 分销商品
	AffiliateItems []*order.SubOrderItem
}

// 发送短信事件
type SendSmsEvent struct {
	// 短信服务商
	Provider int
	// 手机号
	Phone string
	// 短信内容
	Template string
	// 模板代码
	TemplateCode string
	// 短信模板ID
	SpTemplateId string
	// 数据
	Data []string
}

// 会员推送事件
type MemberPushEvent struct {
	// 是否新会员
	IsCreate bool
	// 会员信息
	Member *member.Member
	// 邀请人编号
	InviterId int
}

// 会员账户推送事件
type MemberAccountPushEvent struct {
	member.Account
}

// 订单推送事件
type SubOrderPushEvent struct {
	// 订单号
	OrderNo string
	// 订单金额
	OrderAmount int
	// 收货人
	ConsigneeName string
	// 收货电话
	ConsigneePhone string
	// 收货地址
	ConsigneeAddress string
	// 状态
	OrderState int
}

// 提现申请推送事件
type WithdrawalPushEvent struct {
	// 会员编号
	MemberId int64
	// 流水号Id
	RequestId int
	// 提现金额
	Amount int
	// 手续费
	ProcedureFee int
	// 是否为已审核通过的事件
	IsReviewEvent bool
	// 是否审核通过
	ReviewResult bool
	// 提现账号
	AccountNo string
	// 提现账户名称
	AccountName string
	// 提现银行名称
	BankName string
}
