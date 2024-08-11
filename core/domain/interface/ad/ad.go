/**
 * Copyright 2015 @ 56x.net.
 * name : IAdvertisement
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

//todo: 文字和图片广告待实现

const (
	// 文字广告
	TypeText = 1
	// 图片广告
	TypeImage = 2
	// 图片轮播广告
	TypeSwiper = 3
)

type (
	// IAdvertisementManager 广告管理
	IAdvertisementManager interface {
		// 根据编号获取广告位
		GetPosition(id int64) IAdPosition
		// 根据KEY获取广告位
		GetPositionByKey(key string) *Position
		// 根据广告位KEY获取默认广告
		GetAdByPositionKey(key string) IAdAggregateRoot
		// 获取用户的广告管理
		GetUserAd(adUserId int64) IUserAd
		// 获取广告分组
		GetGroups() []string
		// QueryAd 查询广告列表
		QueryAd(keyword string, size int) []*Ad
	}

	// IAdPosition 广告位
	IAdPosition interface {
		GetAggregateRootId() int64
		SetValue(v *Position) error
		Save() error
		// 设置广告
		PutAd(adId int64) error
		GetValue() Position
	}

	// IUserAd 商户广告聚合根
	IUserAd interface {
		// 获取聚合根标识
		GetAggregateRootId() int64
		// 删除广告
		DeleteAd(adId int64) error
		//获取广告关联的广告位
		GetAdPositionsByAdId(adId int64) []*Position
		// 根据编号获取广告
		GetById(id int64) IAdAggregateRoot
		// 根据KEY获取广告
		GetByPositionKey(key string) IAdAggregateRoot
		// 创建广告对象
		CreateAd(*Ad) IAdAggregateRoot
		// 设置广告
		SetAd(posId, adId int64) error
		QueryAdvertisement(keys []string) map[string]IAdAggregateRoot
	}

	// IAdAggregateRoot 广告接口
	IAdAggregateRoot interface {
		// 获取领域对象编号
		GetDomainId() int64
		// 是否为系统发布的广告
		System() bool
		// 广告类型
		Type() int
		// 广告名称
		Name() string
		// 设置值
		SetValue(*Ad) error
		// 获取值
		GetValue() *Ad
		// 保存广告
		Save() (int64, error)
		// 增加展现次数
		AddShowTimes(times int) error
		// 增加展现次数
		AddClickTimes(times int) error
		// 增加展现次数
		AddShowDays(times int) error
		// 转换为数据传输对象
		Dto() *AdDto
	}

	// 广告数据传输对象
	AdDto struct {
		Id     int64       `json:"id"`
		AdType int         `json:"type"`
		Data   interface{} `json:"data"`
	}

	// 广告仓储
	IAdRepo interface {
		// 获取广告管理器
		GetAdManager() IAdvertisementManager

		// 获取广告分组
		GetPosition(id int64) IAdPosition

		// 根据KEY获取广告位
		GetAdPositionByKey(key string) *Position

		// 根据ID获取广告位
		GetAdPositionById(adPosId int64) *Position

		// 获取广告位
		GetAdPositionsByGroupId(adGroupId int64) []*Position

		// 删除广告位
		DeleteAdPosition(id int64) error

		// 保存广告位
		SaveAdPosition(a *Position) (int64, error)

		// 设置用户的广告
		SetUserAd(adUserId, posId, adId int64) error

		// 根据名称获取广告编号
		GetIdByName(mchId int64, name string) int

		// 保存广告值
		SaveAdValue(*Ad) (int64, error)

		/* ===============  广告类型 ================*/

		// GetTextAdData  获取超链接广告数据
		GetTextAdData(adId int64) *Data

		// SaveTextAdData  保存超链接广告数据
		SaveTextAdData(value *Data) (int64, error)

		// SaveImageAdData  保存广告图片
		SaveImageAdData(*Data) (int64, error)

		// GetAd  获取广告
		GetAd(id int64) *Ad

		// GetAdByKey 根据KEY获取广告
		GetAdByKey(userId int64, key string) *Ad

		// GetSwiperAd 获取轮播广告
		GetSwiperAd(adId int64) SwiperAd

		// GetSwiperAdImage 获取图片项
		GetSwiperAdImage(adId, id int64) *Data

		// DeleteSwiperAdImage 删除图片项
		DeleteSwiperAdImage(adId, id int64) error

		// DeleteAd 删除广告
		DeleteAd(mchId, adId int64) error

		// DeleteImageAdData 删除广告的图片数据
		DeleteImageAdData(adId int64) error

		// DeleteTextAdData 删除广告的文字数据
		DeleteTextAdData(adId int64) error

		// GetGroups 获取广告分组
		GetGroups() []string
		CreateAdPosition(v *Position) IAdPosition
		QueryAdList(keyword string, size int) []*Ad
		// GetPositions 获取所有广告位
		GetPositions() []*Position
	}
)

// AdGroup AdGroup
type AdGroup struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" bson:"id"`
	// Name
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// Opened
	Opened int `json:"opened" db:"opened" gorm:"column:opened" bson:"opened"`
	// Enabled
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
}

func (a AdGroup) TableName() string {
	return "ad_group"
}

// AdList 广告列表
type Ad struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 用户编号
	UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
	// 广告名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 广告类型
	TypeId int `json:"typeId" db:"type_id" gorm:"column:type_id" bson:"typeId"`
	// 展现次数
	ShowTimes int `json:"showTimes" db:"show_times" gorm:"column:show_times" bson:"showTimes"`
	// 点击次数
	ClickTimes int `json:"clickTimes" db:"click_times" gorm:"column:click_times" bson:"clickTimes"`
	// 显示天数
	ShowDays int `json:"showDays" db:"show_days" gorm:"column:show_days" bson:"showDays"`
	// 更新时间
	UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
}

func (a Ad) TableName() string {
	return "ad_list"
}

// AdPosition 广告位
type Position struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 广告位编码
	Key string `json:"key" db:"key" gorm:"column:key" bson:"key"`
	// 广告位名称
	Name string `json:"name" db:"name" gorm:"column:name" bson:"name"`
	// 投放的广告编号
	PutAid int `json:"putAid" db:"put_aid" gorm:"column:put_aid" bson:"putAid"`
	// Opened
	Opened int `json:"opened" db:"opened" gorm:"column:opened" bson:"opened"`
	// Enabled
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 标志
	Flag int `json:"flag" db:"flag" gorm:"column:flag" bson:"flag"`
	// 分组名称
	GroupName string `json:"groupName" db:"group_name" gorm:"column:group_name" bson:"groupName"`
}

func (a Position) TableName() string {
	return "ad_position"
}

// AdUserset AdUserset
type AdUserSet struct {
	// Id
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// PosId
	PosId int `json:"posId" db:"pos_id" gorm:"column:pos_id" bson:"posId"`
	// UserId
	UserId int `json:"userId" db:"user_id" gorm:"column:user_id" bson:"userId"`
	// AdId
	AdId int `json:"adId" db:"ad_id" gorm:"column:ad_id" bson:"adId"`
}

func (a AdUserSet) TableName() string {
	return "ad_userset"
}

// AdData 广告图片
type Data struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 广告编号
	AdId int `json:"adId" db:"ad_id" gorm:"column:ad_id" bson:"adId"`
	// 标题
	Title string `json:"title" db:"title" gorm:"column:title" bson:"title"`
	// 链接地址
	LinkUrl string `json:"linkUrl" db:"link_url" gorm:"column:link_url" bson:"linkUrl"`
	// 图片地址
	ImageUrl string `json:"imageUrl" db:"image_url" gorm:"column:image_url" bson:"imageUrl"`
	// 排列序号
	SortNum int `json:"sortNum" db:"sort_num" gorm:"column:sort_num" bson:"sortNum"`
	// 是否启用
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
}

func (a Data) TableName() string {
	return "ad_data"
}
