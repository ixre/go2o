package staff

import "github.com/ixre/go2o/core/infrastructure/fw"

var (
	// 离线
	WorkStatusOffline = 1
	// 在线空闲
	WorkStatusFree = 2
	// 工作中
	WorkStatusBusy = 3
	// 离职
	WorkStatusOff = 4
)

type (
	// IStaffManager 员工管理接口
	IStaffManager interface {
		// Create 创建员工
		Create(memberId int) error
	}

	// IStaffRepo 员工数据访问接口
	IStaffRepo interface {
		fw.Repository[Staff]
		// GetStaffBy GetBy 商户代理人坐席(员工)
		GetStaffByMemberId(memberId int) *Staff
	}

	// Staff 商户代理人坐席(员工)
	Staff struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes" json:"id" bson:"id"`
		// 会员编号
		MemberId int `db:"member_id" json:"memberId" bson:"memberId"`
		// 站点编号
		StationId int `db:"station_id" json:"stationId" bson:"stationId"`
		// 商户编号
		MchId int `db:"mch_id" json:"mchId" bson:"mchId"`
		// 坐席标志
		Flag int `db:"flag" json:"flag" bson:"flag"`
		// 性别: 0: 未知 1:男 2:女
		Gender int `db:"gender" json:"gender" bson:"gender"`
		// 昵称
		Nickname string `db:"nickname" json:"nickname" bson:"nickname"`
		// 工作状态: 1: 离线 2:在线空闲 3: 工作中 4:离职
		WorkStatus int `db:"work_status" json:"workStatus" bson:"workStatus"`
		// 评分
		Grade int `db:"grade" json:"grade" bson:"grade"`
		// 状态: 1: 正常  2: 锁定
		Status int `db:"status" json:"status" bson:"status"`
		// 是否认证 0:否 1:是
		IsCertified int `db:"is_certified" json:"isCertified" bson:"isCertified"`
		// 认证姓名
		CertifiedName string `db:"certified_name" json:"certifiedName" bson:"certifiedName"`
		// 高级用户等级
		PremiumLevel int `db:"premium_level" json:"premiumLevel" bson:"premiumLevel"`
		// 创建时间
		CreateTime int `db:"create_time" json:"createTime" bson:"createTime"`
	}
)

func (s *Staff) TableName() string {
	return "mch_staff"
}
