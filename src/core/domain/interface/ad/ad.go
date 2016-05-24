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
		GetAdGroup(id int) IAdGroup

		// 删除广告组
		DelAdGroup(id int) error

		// 创建广告组
		CreateAdGroup(name string) IAdGroup

		// 获取用户的广告管理
		GetUserAd(adUserId int) IUserAd
	}

	// 广告分组
	IAdGroup interface {
		// 获取领域编号
		GetDomainId() int
		// 获取值
		GetValue() AdGroup
		// 设置值
		SetValue(v *AdGroup) error
		// 获取广告位
		GetPositions() []*AdPosition
		// 根据Id获取广告位
		GetPosition(id int) *AdPosition
		// 删除广告位
		DelPosition(id int) error
		// 保存广告位
		SavePosition(a *AdPosition) (int, error)
		// 保存,需调用Save()保存
		Save() (int, error)
		// 开放,需调用Save()保存
		Open() error
		// 关闭,需调用Save()保存
		Close() error
		// 启用,需调用Save()保存
		Enabled() error
		// 禁用,需调用Save()保存
		Disabled() error
		// 设置默认广告
		SetDefault(adPositionId int, adId int) error
	}

	// 商户广告聚合根
	IUserAd interface {
		// 获取聚合根标识
		GetAggregateRootId() int

		// 删除广告
		DeleteAdvertisement(advertisementId int) error

		// 根据编号获取广告
		GetById(int) IAdvertisement

		// 根据名称获取广告
		GetByName(string) IAdvertisement

		// 创建广告对象
		CreateAdvertisement(*ValueAdvertisement) IAdvertisement

		// 设置广告
		SetAd(posId, adId int) error
	}

	IAdvertisement interface {
		// 获取领域对象编号
		GetDomainId() int

		// 是否为系统内置的广告
		System() bool

		// 广告类型
		Type() int

		// 广告名称
		Name() string

		// 设置值
		SetValue(*ValueAdvertisement) error

		// 获取值
		GetValue() *ValueAdvertisement

		// 保存广告
		Save() (int, error)
	}

	// 广告分组
	AdGroup struct {
		Id      int    `db:"id" auto:"yes" pk:"yes"`
		Name    string `db:"name"`
		Opened  int    `db:"opened"`
		Enabled int    `db:"enabled"`
	}

	// 广告位
	AdPosition struct {
		// 编号
		Id int `db:"id" auto:"yes" pk:"yes"`
		// 分组编号
		GroupId int `db:"group_id"`
		// 名称
		Name string `db:"name"`
		// 描述
		Description string `db:"description"`
		// 是否开放给外部
		Opened int `db:"opened"`
		// 是否启用
		Enabled int `db:"enabled"`
		// 默认广告
		Default int `db:"default"`
	}

	// 广告用户设置
	AdUserSet struct {
		// 编号
		Id int `db:"id"`

		// 广告位编号
		PosId int `db:"pos_id"`

		//广告用户编号
		AdUserId int `db:"user_id"`

		// 广告编号
		AdId int `db:"ad_id"`
	}

	// 广告
	ValueAdvertisement struct {
		// 编号
		Id int `db:"id" auto:"yes" pk:"yes"`

		//广告用户编号
		AdUserId int `db:"user_id"`

		// 名称
		Name string `db:"name"`

		// 是否内部
		IsInternal int `db:"is_internal"`

		// 广告类型
		Type int `db:"type_id"`

		// 是否启用
		Enabled int `db:"enabled"`

		// 修改时间
		UpdateTime int64 `db:"update_time"`
	}

	ValueText struct {
		// 编号
		Id int `db:"id" auto:"yes" pk:"true"`

		// 广告编号
		AdvertisementId int `db:"ad_id"`

		// 标题
		Title string `db:"title"`

		// 链接
		LinkUrl string `db:"link_url"`

		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 广告图片
	ValueImage struct {
		// 图片编号
		Id int `db:"id" auto:"yes" pk:"true"`

		// 广告编号
		AdvertisementId int `db:"ad_id"`

		// 图片标题
		Title string `db:"title"`

		// 链接
		LinkUrl string `db:"link_url"`

		// 图片地址
		ImageUrl string `db:"image_url"`

		// 排列序号
		SortNumber int `db:"sort_number"`

		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 广告仓储
	IAdRep interface {
		// 获取广告管理器
		GetAdManager() IAdManager

		// 获取广告分组
		GetAdGroups() []*AdGroup

		// 删除广告组
		DelAdGroup(id int) error

		// 获取广告位
		GetAdPositionsByGroupId(adGroupId int) []*AdPosition

		// 删除广告位
		DelAdPosition(id int) error

		// 保存广告位
		SaveAdPosition(a *AdPosition) (int, error)

		// 保存
		SaveAdGroup(value *AdGroup) (int, error)

		// 设置用户的广告
		SetUserAd(adUserId, posId, adId int) error

		// 根据名称获取广告编号
		GetIdByName(merchantId int, name string) int

		// 保存广告值
		SaveAdvertisementValue(*ValueAdvertisement) (int, error)

		// 保存广告图片
		SaveAdImageValue(*ValueImage) (int, error)

		// 获取广告
		GetValueAdvertisement(id int) *ValueAdvertisement

		// 根据名称获取广告
		GetValueAdvertisementByName(merchantId int, name string) *ValueAdvertisement

		// 获取轮播广告
		GetValueGallery(advertisementId int) ValueGallery

		// 获取图片项
		GetValueAdImage(advertisementId, id int) *ValueImage

		// 删除图片项
		DelAdImage(advertisementId, id int) error

		// 删除广告
		DelAdvertisement(merchantId, advertisementId int) error

		// 删除广告的图片数据
		DelImageDataForAdvertisement(advertisementId int) error

		// 删除广告的文字数据
		DelTextDataForAdvertisement(advertisementId int) error
	}
)
