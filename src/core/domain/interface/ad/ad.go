/**
 * Copyright 2015 @ z3q.net.
 * name : IAdvertisement
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad

const (
	// 文字广告
	TypeHyperLink = 1
	// 图片广告
	TypeImage = 2
	// 图片轮播广告
	TypeGallery = 3
)

type (
	// 商户广告聚合根
	IMerchantAdvertisement interface {
		// 获取聚合根标识
		GetAggregateRootId() int

		// 初始化内置默认的广告
		InitInternalAdvertisements()

		// 删除广告
		DeleteAdvertisement(advertisementId int) error

		// 根据编号获取广告
		GetById(int) IAdvertisement

		// 根据名称获取广告
		GetByName(string) IAdvertisement

		// 创建广告对象
		CreateAdvertisement(*ValueAdvertisement) IAdvertisement
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

	ValueAdvertisement struct {
		// 编号
		Id int `db:"id" auto:"yes" pk:"yes"`

		MerchantId int `db:"merchant_id"`
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
)
