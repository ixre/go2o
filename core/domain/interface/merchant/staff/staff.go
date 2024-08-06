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
		// 创建时间
		CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
		// 服务总次数
		ServiceCount int `json:"serviceCount" db:"service_count" gorm:"column:service_count" bson:"serviceCount"`
		// 服务总人数
		CusCount int `json:"cusCount" db:"cus_count" gorm:"column:cus_count" bson:"cusCount"`
	}
)

func (s *Staff) TableName() string {
	return "mch_staff"
}
