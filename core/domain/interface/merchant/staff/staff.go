package staff

import (
	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var (
	// 离线
	WorkStatusOffline = 1
	// 在线空闲
	WorkStatusIdle = 2
	// 工作中
	WorkStatusBusy = 3
	// 离职
	WorkStatusOff = 4
)

var (
	// 审核状态
	ReviewStatusPending = 1
	// 审核拒绝
	ReviewStatusRejected = 2
	// 审核通过
	ReviewStatusApproved = 3
)

type (
	// IStaffManager 员工管理接口
	IStaffManager interface {
		// Create 创建员工
		Create(memberId int) error
		// RequestTransfer 请求转商户
		RequestTransfer(staffId, mchId int) (int, error)
		// TransferApproval 处理转商户审批
		TransferApproval(trans *StaffTransfer, event *approval.ApprovalProcessEvent) error
		// UpdateWorkStatus 更新员工工作状态
		UpdateWorkStatus(staffId int, workStatus int, isKeepOnline bool) error
		// IsKeepOnline 是否保持上线
		IsKeepOnline(staffId int) bool
	}

	// IStaffRepo 员工数据访问接口
	IStaffRepo interface {
		fw.Repository[Staff]
		// GetStaffBy GetBy 商户代理人坐席(员工)
		GetStaffByMemberId(memberId int) *Staff
		// TransferRepo 员工转商户仓储
		TransferRepo() fw.Repository[StaffTransfer]
	}

	// StaffRequireImInitEvent 员工IM初始化事件
	StaffRequireImInitEvent struct {
		// 员工信息
		Staff Staff
	}

	// StaffTransferApprovedEvent 员工转移审批通过事件
	StaffTransferApprovedEvent struct {
		// 员工信息
		Staff Staff
		// 原商户
		OriginMchId int
		// 转移商户
		TransferMchId int
	}

	// MchStaff 商户代理人坐席(员工)
	Staff struct {
		// 编号
		Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
		// 会员编号
		MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
		// 站点编号
		StationId int `json:"stationId" db:"station_id" gorm:"column:station_id" bson:"stationId"`
		// 商户编号
		MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" bson:"mchId"`
		// 坐席标志
		Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
		// 性别: 0: 未知 1:男 2:女
		Gender int `json:"gender" db:"gender" gorm:"column:gender" bson:"gender"`
		// 昵称
		Nickname string `json:"nickname" db:"nickname" gorm:"column:nickname" bson:"nickname"`
		// 工作状态: 1: 离线 2:在线空闲 3: 工作中
		WorkStatus int `json:"workStatus" db:"work_status" gorm:"column:work_status" bson:"workStatus"`
		// 评分
		Grade int `json:"grade" db:"grade" gorm:"column:grade" bson:"grade"`
		// 状态: 1: 正常  2: 锁定
		Status int `json:"status" db:"status" gorm:"column:status" bson:"status"`
		// 是否认证 0:否 1:是
		IsCertified int `json:"isCertified" db:"is_certified" gorm:"column:is_certified" bson:"isCertified"`
		// 认证姓名
		CertifiedName string `json:"certifiedName" db:"certified_name" gorm:"column:certified_name" bson:"certifiedName"`
		// 高级用户等级
		PremiumLevel int `json:"premiumLevel" db:"premium_level" gorm:"column:premium_level" bson:"premiumLevel"`
		// 最后在线时间
		LastOnlineTime int `json:"lastOnlineTime" db:"last_online_time" gorm:"column:last_online_time" bson:"lastOnlineTime"`
		// 创建时间
		CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
		// 服务总次数
		ServiceCount int `json:"serviceCount" db:"service_count" gorm:"column:service_count" bson:"serviceCount"`
		// 服务总人数
		CusCount int `json:"cusCount" db:"cus_count" gorm:"column:cus_count" bson:"cusCount"`
		// IM是否注册
		ImInitialized int `json:"imInitialized" db:"im_initialized" gorm:"column:im_initialized" bson:"imInitialized"`
	}
)

func (s *Staff) TableName() string {
	return "mch_staff"
}

// StaffTransfer 员工转商户
type StaffTransfer struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 员工编号
	StaffId int `json:"staffId" db:"staff_id" gorm:"column:staff_id" bson:"staffId"`
	// 原商户
	OriginMchId int `json:"originMchId" db:"origin_mch_id" gorm:"column:origin_mch_id" bson:"originMchId"`
	// 转移商户
	TransferMchId int `json:"transferMchId" db:"transfer_mch_id" gorm:"column:transfer_mch_id" bson:"transferMchId"`
	// 审批编号
	ApprovalId int `json:"approvalId" db:"approval_id" gorm:"column:approval_id" bson:"approvalId"`
	// 审核状态
	ReviewStatus int `json:"reviewStatus" db:"review_status" gorm:"column:review_status" bson:"reviewStatus"`
	// 审核备注
	ReviewRemark string `json:"reviewRemark" db:"review_remark" gorm:"column:review_remark" bson:"reviewRemark"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (m StaffTransfer) TableName() string {
	return "mch_staff_transfer"
}
