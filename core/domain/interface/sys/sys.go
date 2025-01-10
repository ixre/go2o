/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: sys.go
 * author: jarrysix (jarrysix#gmail.com)
 * date: 2024-05-11 19:10:56
 * description:
 * history:
*
 * 行政区划数据参考： https://open.yeepay.com/docs/v2/products/fwssfk/others/5f59fc1720289f001ba82528/5f59fcd020289f001ba82529/index.html
*/

package sys

import (
	"sort"

	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

// ISystemAggregateRoot 系统聚合根
type (
	ISystemAggregateRoot interface {
		domain.IAggregateRoot

		// Location 获取地址管理器
		Location() ILocationManager

		// Options 获取选项管理器
		Options() IOptionManager

		// GetBanks 获取银行列表
		GetBanks() []*GeneralOption

		// Stations 获取站点管理器
		Stations() IStationManager

		// Application 获取应用管理器
		Application() IApplicationManager

		// 标记更新状态,通常监听数据变更或事件后调用
		FlushUpdateStatus()

		// 获取最后更新时间,用于对比系统设置是否已经变更
		LastUpdateTime() int64

		// OAuth 获取OAuth管理器
		OAuth() IOAuthManager
	}

	// ILocationManager 地址管理器
	ILocationManager interface {
		// GetAllCities 获取所有城市
		GetAllCities() []*District
		// GetChildrenDistricts 获取区域信息
		GetChildrenDistricts(parentId int) []*District
		// GetDistricts 获取区域名称
		GetDistrictNames(code ...int) map[int]string
		// FindCity 查找城市
		FindCity(name string) *District
		// GetDistrict 获取区域信息
		GetDistrict(id int) *District
		// FindRegionByIp 根据IP查找区域信息
		FindRegionByIp(ip string) (*District, error)
	}

	// IOptionManager 选项管理器
	IOptionManager interface {
		// SaveOption 保存选项
		SaveOption(option *GeneralOption) error
		// GetOptionNames 获取选项名称
		GetOptionNames(code ...int) map[int]string
		// 获取下级选项
		GetChildOptions(parentId int, typeName string) []*GeneralOption
		// 是否为叶子节点
		IsLeaf(n *GeneralOption) bool
		// 删除选项
		Delete(option *GeneralOption) error
	}
)

// ISystemRepo 系统仓储
type ISystemRepo interface {
	// GetSystemAggregateRoot 获取系统聚合根
	GetSystemAggregateRoot() ISystemAggregateRoot
	// 标记已更新
	FlushUpdateStatus()
	// 获取最后更新时间
	LastUpdateTime() int64
	// District 获取区域仓储
	District() fw.Repository[District]
	// Option 获取选项仓储
	Option() fw.Repository[GeneralOption]
	// Station 获取站点仓储
	Station() IStationRepo
	// Log 获取日志仓储
	App() IApplicationRepository
}

type (
	// District ChinaArea
	District struct {
		// Code
		Code int `db:"code" pk:"yes" json:"code" bson:"code"`
		// Name
		Name string `db:"name" json:"name" bson:"name"`
		// Parent
		Parent int `db:"parent" json:"parent" bson:"parent"`
	}
)

func (a District) TableName() string {
	return "sys_district"
}

// GeneralOption 系统通用选项(用于存储分类和选项等数据)

type GeneralOption struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 类型
	Type string `json:"type" db:"type" gorm:"column:type" bson:"type"`
	// 上级编号
	Pid int `json:"pid" db:"pid" gorm:"column:pid" bson:"pid"`
	// 名称
	Label string `json:"label" db:"label" gorm:"column:label" bson:"label"`
	// 值
	Value string `json:"value" db:"value" gorm:"column:value" bson:"value"`
	// 排列序号
	SortNum int `json:"sortNum" db:"sort_num" gorm:"column:sort_num" bson:"sortNum"`
	// 是否启用
	Enabled int `json:"enabled" db:"enabled" gorm:"column:enabled" bson:"enabled"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 子选项
	Children []*GeneralOption `db:"-" gorm:"-:all" json:"children" bson:"-"`
	// 是否为叶子节点
	IsLeaf bool `db:"-" json:"isLeaf" gorm:"-:all" bson:"-"`
}

func (s GeneralOption) TableName() string {
	return "sys_general_option"
}

// GeneralOptions 通用选项列表排序
var _ sort.Interface = GeneralOptions{}

type GeneralOptions []*GeneralOption

// Less implements sort.Interface.
func (s GeneralOptions) Less(i int, j int) bool {
	return s[i].SortNum < s[j].SortNum
}

// Swap implements sort.Interface.
func (s GeneralOptions) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GeneralOptions) Len() int {
	return len(s)
}
