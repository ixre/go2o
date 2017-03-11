/**
 * Copyright 2015 @ z3q.net.
 * name : shop
 * author : jarryliu
 * date : 2016-05-29 11:12
 * description :
 * history :
 */
package shop

import (
	"fmt"
	"go2o/core/infrastructure/domain"
	"strconv"
)

var (
	ErrNoSuchShop *domain.DomainError = domain.NewDomainError(
		"err_shop_no_such_shop", "未指定店铺")

	ErrNotSetAlias *domain.DomainError = domain.NewDomainError(
		"err_not_set_alias", "请设置商城域名")

	ErrShopAliasFormat *domain.DomainError = domain.NewDomainError(
		"err_shop_alias_format", "域名前缀为3-11位的字母数字的组合")

	ErrShopAliasIncorrect *domain.DomainError = domain.NewDomainError(
		"err_shop_alias_incorrect", "域名前缀不允许使用")

	ErrShopAliasUsed *domain.DomainError = domain.NewDomainError(
		"err_shop_alias_used", "域名已被占用")

	ErrSupportSingleOnlineShop *domain.DomainError = domain.NewDomainError(
		"err_shop_support_only_online_shop", "当前商户仅支持1个店铺")
)

const (
	// 线上商店
	TypeOnlineShop = 1
	// 线下实体店
	TypeOfflineShop = 2
)

const (
	// 停用状态
	StateStopped = 0
	// 待审核状态
	StateAwaitingReview = 1
	// 正常状态
	StateNormal = 2
	// 暂停状态
	StateSuspend = 3
	// 暂停状态
	StatePause = 4
)

var (
	ErrSameNameShopExists *domain.DomainError = domain.NewDomainError(
		"err_same_name_shop_exists", "商店已经存在")

	// 商店状态字典
	StateTextMap = map[int]string{
		StateStopped:        "已停用",
		StateAwaitingReview: "待审核",
		StateNormal:         "正常",
		StateSuspend:        "系统挂起",
		StatePause:          "商户暂停",
	}

	// 状态顺序
	StateIndex = []int{
		StateStopped,
		StateAwaitingReview,
		StateNormal,
		StateSuspend,
		StatePause,
	}

	// 商店类型字电
	TypeTextMap = map[int]string{
		TypeOnlineShop:  "商店",
		TypeOfflineShop: "门店",
	}

	// 类型顺序
	TypeIndex = []int{
		TypeOnlineShop,
		TypeOfflineShop,
	}

	StateTextStrMap = map[string]string{
		strconv.Itoa(StateStopped):        StateTextMap[StateStopped],
		strconv.Itoa(StateAwaitingReview): StateTextMap[StateAwaitingReview],
		strconv.Itoa(StateNormal):         StateTextMap[StateNormal],
		strconv.Itoa(StateSuspend):        StateTextMap[StateSuspend],
		strconv.Itoa(StatePause):          StateTextMap[StatePause],
	}

	TypeTextStrMap = map[string]string{
		strconv.Itoa(TypeOnlineShop):  TypeTextMap[TypeOnlineShop],
		strconv.Itoa(TypeOfflineShop): TypeTextMap[TypeOfflineShop],
	}

	DefaultOnlineShop = OnlineShop{
		// 通讯地址
		Address: "",
		// 联系电话
		Tel: "",
		//别名,用于在商店域名
		Alias: "",
		//域名
		Host: "",
		//前台Logo
		Logo: "res/shop_logo.png",
		//首页标题
		IndexTitle: "",
		//子页面标题
		SubTitle: "",
		// Notice
		Notice: "",
	}

	DefaultOfflineShop = OfflineShop{
		// 联系电话
		Tel: "",
		// 通讯地址
		Address: "",
		// 经度
		Lng: 0,
		// 纬度
		Lat: 0,
		// 配送最大半径(公里)
		DeliverRadius: 0,
	}
)

type (
	IShop interface {
		GetDomainId() int32

		// 商店类型
		Type() int

		// 获取值
		GetValue() Shop

		// 设置值
		SetValue(*Shop) error

		// 保存
		Save() (int32, error)

		// 数据
		Data() *ShopDto
	}

	// 线上商城
	IOnlineShop interface {
		// 设置值
		SetShopValue(*OnlineShop) error

		// 获取值
		GetShopValue() OnlineShop
	}

	// 线下商店
	IOfflineShop interface {
		// 设置值
		SetShopValue(*OfflineShop) error

		// 获取值
		GetShopValue() OfflineShop

		// 获取经维度
		GetLngLat() (float64, float64)

		// 是否可以配送
		// 返回是否可以配送，以及距离(米)
		CanDeliver(lng, lat float64) (bool, int)

		// 是否可以配送
		// 返回是否可以配送，以及距离(米)
		CanDeliverTo(address string) (bool, int)
	}

	// 商店
	Shop struct {
		Id         int32  `db:"id" pk:"yes" auto:"yes"`
		MerchantId int32  `db:"mch_id"`
		ShopType   int    `db:"shop_type"`
		Name       string `db:"name"`
		State      int    `db:"state"`
		SortNum    int    `db:"sort_number"`
		CreateTime int64  `db:"create_time"`
	}

	// 商店数据传输对象
	ShopDto struct {
		Id         int32
		MerchantId int32
		ShopType   int
		Name       string
		State      int
		CreateTime int64
		// 线上/线下商店的数据
		Data interface{}
	}

	// 商城
	OnlineShop struct {
		// 商店编号
		ShopId int32 `db:"shop_id" pk:"yes" auto:"no"`
		// 通讯地址
		Address string `db:"addr"`
		// 联系电话
		Tel string `db:"tel"`

		//别名,用于在商店域名
		Alias string `db:"alias"`

		//域名
		Host string `db:"host"`

		//前台Logo
		Logo string `db:"logo"`

		//首页标题
		IndexTitle string `db:"index_tit"`

		//子页面标题
		SubTitle string `db:"sub_tit"`

		// Notice
		Notice string `db:"notice_html"`
	}

	// 门店
	OfflineShop struct {
		// 商店编号
		ShopId int32 `db:"shop_id" pk:"yes" auto:"no"`

		// 联系电话
		Tel string `db:"tel"`

		// 省
		Province int `db:"province"`

		// 市
		City int `db:"city"`

		// 区
		District int `db:"district"`

		// 通讯地址
		Address string `db:"addr"`

		// 经度
		Lng float32 `db:"lng"`

		// 纬度
		Lat float32 `db:"lat"`

		// 配送最大半径(公里)
		DeliverRadius int `db:"deliver_radius"`
	}
)

//位置(经度+"/"+纬度)
func (o OfflineShop) Location() string {
	return fmt.Sprintf("%f/%f", o.Lng, o.Lat)
}
