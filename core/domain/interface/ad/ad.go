/**
 * Copyright 2015 @ z3q.net.
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
	TypeHyperLink = 1
	// 图片广告
	TypeImage = 2
	// 图片轮播广告
	TypeGallery = 3
)

type (
	// 广告管理
	IAdManager interface {
		// 获取广告分组
		GetAdGroups() []IAdGroup
		// 获取单个广告分组
		GetAdGroup(id int32) IAdGroup
		// 删除广告组
		DelAdGroup(id int32) error
		// 创建广告组
		CreateAdGroup(name string) IAdGroup
		// 根据编号获取广告位
		GetAdPositionById(id int32) *AdPosition
		// 根据KEY获取广告位
		GetAdPositionByKey(key string) *AdPosition
		// 根据广告位KEY获取默认广告
		GetAdByPositionKey(key string) IAd
		// 获取用户的广告管理
		GetUserAd(adUserId int32) IUserAd
	}

	// 广告分组
	IAdGroup interface {
		// 获取领域编号
		GetDomainId() int32
		// 获取值
		GetValue() AdGroup
		// 设置值
		SetValue(v *AdGroup) error
		// 获取广告位
		GetPositions() []*AdPosition
		// 根据Id获取广告位
		GetPosition(id int32) *AdPosition
		// 删除广告位
		DelPosition(id int32) error
		// 保存广告位
		SavePosition(a *AdPosition) (int32, error)
		// 保存,需调用Save()保存
		Save() (int32, error)
		// 开放,需调用Save()保存
		Open() error
		// 关闭,需调用Save()保存
		Close() error
		// 启用,需调用Save()保存
		Enabled() error
		// 禁用,需调用Save()保存
		Disabled() error
		// 设置默认广告
		SetDefault(posId int32, adId int32) error
	}

	// 商户广告聚合根
	IUserAd interface {
		// 获取聚合根标识
		GetAggregateRootId() int32
		// 删除广告
		DeleteAd(adId int32) error
		//获取广告关联的广告位
		GetAdPositionsByAdId(adId int32) []*AdPosition
		// 根据编号获取广告
		GetById(id int32) IAd
		// 根据KEY获取广告
		GetByPositionKey(key string) IAd
		// 创建广告对象
		CreateAd(*Ad) IAd
		// 设置广告
		SetAd(posId, adId int32) error
	}

	// 广告接口
	IAd interface {
		// 获取领域对象编号
		GetDomainId() int32
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
		Save() (int32, error)
		// 增加展现次数
		AddShowTimes(times int) error
		// 增加展现次数
		AddClickTimes(times int) error
		// 增加展现次数
		AddShowDays(times int) error
		// 转换为数据传输对象
		Dto() *AdDto
	}

	// 广告分组
	AdGroup struct {
		ID      int32  `db:"id" auto:"yes" pk:"yes"`
		Name    string `db:"name"`
		Opened  int    `db:"opened"`
		Enabled int    `db:"enabled"`
	}

	// 广告位
	AdPosition struct {
		// 编号
		ID int32 `db:"id" auto:"yes" pk:"yes"`
		// 分组编号
		GroupId int32 `db:"group_id"`
		// 引用键
		Key string `db:"key"`
		// 名称
		Name string `db:"name"`
		//todo:广告位类型限制
		// 广告类型限制,0为无限制
		TypeLimit int `db:"-"` //`db:"type_limit"`
		// 是否开放给外部
		Opened int `db:"opened"`
		// 是否启用
		Enabled int `db:"enabled"`
		// 默认广告编号
		DefaultId int32 `db:"default_id"`
	}

	// 广告用户设置
	AdUserSet struct {
		// 编号
		Id int32 `db:"id"`

		// 广告位编号
		PosId int32 `db:"pos_id"`

		//广告用户编号
		AdUserId int32 `db:"user_id"`

		// 广告编号
		AdId int32 `db:"ad_id"`
	}

	// 广告
	Ad struct {
		// 编号
		Id int32 `db:"id" auto:"yes" pk:"yes"`

		//广告用户编号
		UserId int32 `db:"user_id"`

		// 名称
		Name string `db:"name"`

		// 广告类型
		Type int `db:"type_id"`

		// 展现次数
		ShowTimes int `db:"show_times" json:"-"`

		// 点击次数
		ClickTimes int `db:"click_times" json:"-"`

		// 展现天数
		ShowDays int `db:"show_days" json:"-"`

		// 修改时间
		UpdateTime int64 `db:"update_time" json:"-"`
	}

	// 广告数据传输对象
	AdDto struct {
		Id   int32       `json:"id"`
		Type int         `json:"type"`
		Data interface{} `json:"data"`
	}

	// 广告仓储
	IAdRepo interface {
		// 获取广告管理器
		GetAdManager() IAdManager

		// 获取广告分组
		GetAdGroups() []*AdGroup

		// 删除广告组
		DelAdGroup(id int32) error

		// 根据KEY获取广告位
		GetAdPositionByKey(key string) *AdPosition

		// 根据ID获取广告位
		GetAdPositionById(adPosId int32) *AdPosition

		// 获取广告位
		GetAdPositionsByGroupId(adGroupId int32) []*AdPosition

		// 删除广告位
		DelAdPosition(id int32) error

		// 保存广告位
		SaveAdPosition(a *AdPosition) (int32, error)

		// 保存
		SaveAdGroup(value *AdGroup) (int32, error)

		// 设置用户的广告
		SetUserAd(adUserId, posId, adId int32) error

		// 根据名称获取广告编号
		GetIdByName(mchId int32, name string) int

		// 保存广告值
		SaveAdValue(*Ad) (int32, error)

		/* ===============  广告类型 ================*/

		// 获取超链接广告数据
		GetHyperLinkData(adId int32) *HyperLink

		// 保存超链接广告数据
		SaveHyperLinkData(value *HyperLink) (int32, error)

		// 保存广告图片
		SaveAdImageValue(*Image) (int32, error)

		// 获取广告
		GetValueAd(id int32) *Ad

		// 根据KEY获取广告
		GetAdByKey(userId int32, name string) *Ad

		// 获取轮播广告
		GetValueGallery(adId int32) ValueGallery

		// 获取图片项
		GetValueAdImage(adId, id int32) *Image

		// 删除图片项
		DelAdImage(adId, id int32) error

		// 删除广告
		DelAd(mchId, adId int32) error

		// 删除广告的图片数据
		DelImageDataForAdvertisement(adId int32) error

		// 删除广告的文字数据
		DelTextDataForAdvertisement(adId int32) error
	}
)
