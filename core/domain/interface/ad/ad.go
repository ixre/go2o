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
		GetAdByPositionKey(key string) IAd
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
		GetById(id int64) IAd
		// 根据KEY获取广告
		GetByPositionKey(key string) IAd
		// 创建广告对象
		CreateAd(*Ad) IAd
		// 设置广告
		SetAd(posId, adId int64) error
		QueryAdvertisement(keys []string) map[string]IAd
	}

	// IAd 广告接口
	IAd interface {
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

	// Position 广告位
	Position struct {
		// 编号
		Id int64 `db:"id" auto:"yes" pk:"yes"`
		// 引用键
		Key string `db:"key"`
		// 名称
		Name string `db:"name"`
		// 分组名称
		GroupName string `db:"group_name"`
		//todo:广告位类型限制
		// 广告类型限制,0为无限制
		TypeLimit int `db:"-"` //`db:"type_limit"`
		// 标志
		Flag int `db:"flag"`
		// 是否开放给外部
		Opened int `db:"opened"`
		// 是否启用
		Enabled int `db:"enabled"`
		// 默认广告编号
		PutAdId int64 `db:"put_aid"`
	}

	// 广告用户设置
	AdUserSet struct {
		// 编号
		Id int64 `db:"id"`

		// 广告位编号
		PosId int64 `db:"pos_id"`

		//广告用户编号
		AdUserId int64 `db:"user_id"`

		// 广告编号
		AdId int64 `db:"ad_id"`
	}

	// 广告
	Ad struct {
		// 编号
		Id int64 `db:"id" auto:"yes" pk:"yes"`

		//广告用户编号
		UserId int64 `db:"user_id"`

		// 名称
		Name string `db:"name"`

		// 广告类型
		AdType int `db:"type_id"`

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
		GetTextAdData(adId int64) *HyperLink

		// SaveTextAdData  保存超链接广告数据
		SaveTextAdData(value *HyperLink) (int64, error)

		// SaveImageAdData  保存广告图片
		SaveImageAdData(*Image) (int64, error)

		// GetAd  获取广告
		GetAd(id int64) *Ad

		// GetAdByKey 根据KEY获取广告
		GetAdByKey(userId int64, key string) *Ad

		// GetSwiperAd 获取轮播广告
		GetSwiperAd(adId int64) SwiperAd

		// GetSwiperAdImage 获取图片项
		GetSwiperAdImage(adId, id int64) *Image

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
