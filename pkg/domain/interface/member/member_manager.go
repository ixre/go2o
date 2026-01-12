/**
 * Copyright 2015 @ 56x.net.
 * name : level_manager
 * author : jarryliu
 * date : 2016-05-26 11:33
 * description :
 * history :
 */
package member

const (
	RegisterModeNormal         = 1 // 正常注册
	RegisterModeClosed         = 2 // 关闭注册
	RegisterModeMustRedirect   = 3 // 必须直接注册
	RegisterModeMustInvitation = 4 // 必须邀请注册

)

type (
	// 会员服务
	IMemberManager interface {
		// 等级服务
		LevelManager() ILevelManager
		// 检查手机绑定,同时检查手机格式
		CheckPhoneBind(phone string, memberId int) error
		// 检查注册信息是否正确
		PrepareRegister(v *Member, pro *Profile, invitationCode string) (
			invitationId int64, err error)
		// 检查邀请注册
		CheckInviteRegister(code string) (inviterId int64, err error)
		// 获取所有买家分组
		GetAllBuyerGroups() []*BuyerGroup
		// 获取买家分组
		GetBuyerGroup(id int32) *BuyerGroup
		// 保存买家分组
		SaveBuyerGroup(*BuyerGroup) (int32, error)

		// IDocManager()IDocManager
	}

	// 买家（客户）分组
	BuyerGroup struct {
		// 编号
		ID        int32  `db:"id" pk:"yes" auto:"yes"`
		Name      string `db:"name"`
		IsDefault int32  `db:"is_default"`
	}

	ILevelManager interface {
		// 获取等级设置
		GetLevelSet() []*Level

		// 获取初始等级
		GetInitialLevel() *Level

		// 获取最高已启用的等级
		GetHighestLevel() *Level

		// 获取等级,todo:返回error
		GetLevelById(id int) *Level

		// 根据可编程字符获取会员等级
		GetLevelByProgramSign(sign string) *Level

		// 获取下一个等级
		GetNextLevelById(id int) *Level

		// 删除等级
		DeleteLevel(id int) error

		// 保存等级
		SaveLevel(*Level) (int, error)

		// 根据经验值获取等级值
		GetLevelIdByExp(exp int) int
	}
)

// MmLevel 会员等级
type Level struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 等级名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 需要经验值
	RequireExp int `json:"requireExp" db:"require_exp" gorm:"column:require_exp" bson:"requireExp"`
	// 编程符号
	ProgramSignal string `json:"programSignal" db:"program_signal" gorm:"column:program_signal" bson:"programSignal"`
	// 是否正式的会员等级
	IsOfficial int `json:"isOfficial" db:"is_official" gorm:"column:is_official" bson:"isOfficial"`
	// 是否启用
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 是否允许自动升级
	AllowUpgrade int `json:"allowUpgrade" db:"allow_upgrade" gorm:"column:allow_upgrade" bson:"allowUpgrade"`
}

func (m Level) TableName() string {
	return "mm_level"
}

// MmLevelup 会员升级日志表
type LevelUpLog struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 会员编号
	MemberId int `json:"memberId" db:"member_id" gorm:"column:member_id" bson:"memberId"`
	// 原来等级
	OriginLevel int `json:"originLevel" db:"origin_level" gorm:"column:origin_level" bson:"originLevel"`
	// 现在等级
	TargetLevel int `json:"targetLevel" db:"target_level" gorm:"column:target_level" bson:"targetLevel"`
	// 是否为免费升级的会员
	IsFree int `json:"isFree" db:"is_free" gorm:"column:is_free" bson:"isFree"`
	// 支付单编号
	PaymentId int `json:"paymentId" db:"payment_id" gorm:"column:payment_id" bson:"paymentId"`
	// UpgradeMode
	UpgradeMode int `json:"upgradeMode" db:"upgrade_mode" gorm:"column:upgrade_mode" bson:"upgradeMode"`
	// ReviewStatus
	ReviewStatus int `json:"reviewStatus" db:"review_status" gorm:"column:review_status" bson:"reviewStatus"`
	// 升级时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
}

func (m LevelUpLog) TableName() string {
	return "mm_levelup"
}
