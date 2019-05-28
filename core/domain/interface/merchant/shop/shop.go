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
	ErrNoSuchShop *domain.DomainError = domain.NewError(
		"err_shop_no_such_shop", "未指定店铺")

	ErrShopNoLogo *domain.DomainError = domain.NewError(
		"err_shop_no_logo", "请上传店铺Logo")

	ErrShopAliasFormat *domain.DomainError = domain.NewError(
		"err_shop_alias_format", "域名前缀为3-11位的字母数字的组合")

	ErrShopAliasIncorrect *domain.DomainError = domain.NewError(
		"err_shop_alias_incorrect", "域名前缀不允许使用")

	ErrShopAliasUsed *domain.DomainError = domain.NewError(
		"err_shop_alias_used", "域名已被占用")

	ErrSupportSingleOnlineShop *domain.DomainError = domain.NewError(
		"err_shop_support_only_online_shop", "当前商户仅支持1个店铺")
)

const (
	// 线上商店
	TypeOnlineShop int32 = 1
	// 线下实体店
	TypeOfflineShop int32 = 2
)

const (
	// 正常状态
	StateNormal int32 = 1
	// 停用状态
	StateStopped int32 = 2
	// 营业状态-正常
	OStateNormal int32 = 1
	// 营业状态-暂停营业
	OStatePause int32 = 2
)

var (
	ErrSameNameShopExists *domain.DomainError = domain.NewError(
		"err_same_name_shop_exists", "商店已经存在")

	// 商店状态字典
	StateTextMap = map[int32]string{
		StateNormal:  "正常",
		StateStopped: "已关闭",
	}

	// 状态顺序
	StateIndex = []int32{
		StateNormal,
		StateStopped,
	}

	// 商店类型字电
	TypeTextMap = map[int32]string{
		TypeOnlineShop:  "商店",
		TypeOfflineShop: "门店",
	}

	StateTextStrMap = map[string]string{
		strconv.Itoa(int(StateNormal)):  StateTextMap[StateNormal],
		strconv.Itoa(int(StateStopped)): StateTextMap[StateStopped],
	}

	TypeTextStrMap = map[string]string{
		strconv.Itoa(int(TypeOnlineShop)):  TypeTextMap[TypeOnlineShop],
		strconv.Itoa(int(TypeOfflineShop)): TypeTextMap[TypeOfflineShop],
	}

	DefaultOnlineShop = OnlineShop{
		// 通讯地址
		Address: "",
		// 联系电话
		ServiceTel: "",
		//别名,用于在商店域名
		Alias: "",
		//域名
		Host: "",
		//前台Logo
		Logo: "res/shop_logo.png",
		//首页标题
		ShopTitle: "",
		// ShopNotice
		ShopNotice: "",
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
		CoverRadius: 0,
	}
)

type (
	IShop interface {
		GetDomainId() int32
		// 商店类型
		Type() int32
		// 获取值
		GetValue() Shop
		// 设置值
		SetValue(*Shop) error
		// 开启店铺
		TurnOn() error
		// 停用店铺
		TurnOff(reason string) error
		// 商店营业
		Opening() error
		// 商店暂停营业
		Pause() error
		// 保存
		Save() (int32, error)
		// 获取商店信息
		Data() *ComplexShop
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
		//商店编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		//运营商编号
		VendorId int32 `db:"vendor_id"`
		//商店类型
		ShopType int32 `db:"shop_type"`
		//商店名称
		Name string `db:"name"`
		//商店状态
		State int32 `db:"state"`
		//营业状态
		OpeningState int32 `db:"opening_state"`
		//排序
		SortNum int32 `db:"sort_num"`
		//创建时间
		CreateTime int64 `db:"create_time"`
	}

	// 商店数据传输对象
	ComplexShop struct {
		ID       int32
		VendorId int32
		ShopType int32
		Name     string
		State    int32
		// 线上/线下商店的数据
		Data map[string]string
	}

	// 商城
	OnlineShop struct {
		// 商店编号
		ShopId int32 `db:"shop_id" pk:"yes" auto:"no"`
		// 通讯地址
		Address string `db:"addr"`
		// 联系电话
		ServiceTel string `db:"tel"`
		//别名,用于在商店域名
		Alias string `db:"alias"`
		//域名
		Host string `db:"host"`
		//前台Logo
		Logo string `db:"logo"`
		//首页标题
		ShopTitle string `db:"shop_title"`
		// ShopNotice
		ShopNotice string `db:"shop_notice"`
	}

	// 门店
	OfflineShop struct {
		// 商店编号
		ShopId int32 `db:"shop_id" pk:"yes" auto:"no"`
		// 联系电话
		Tel string `db:"tel"`
		// 省
		Province int32 `db:"province"`
		// 市
		City int32 `db:"city"`
		// 区
		District int32 `db:"district"`
		// 通讯地址
		Address string `db:"addr"`
		// 经度
		Lng float32 `db:"lng"`
		// 纬度
		Lat float32 `db:"lat"`
		// 配送最大半径(公里)
		CoverRadius int `db:"deliver_radius"`
	}
)

//位置(经度+"/"+纬度)
func (o OfflineShop) Location() string {
	return fmt.Sprintf("%f/%f", o.Lng, o.Lat)
}
