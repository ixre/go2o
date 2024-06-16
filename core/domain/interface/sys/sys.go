package sys

import (
	"github.com/ixre/go2o/core/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

// ISystemAggregateRoot 系统聚合根
type (
	ISystemAggregateRoot interface {
		domain.IAggregateRoot

		// 获取地址管理器
		Address() IAddressManager

		// 获取选项管理器
		Options() IOptionManager
	}

	// IAddressManager 地址管理器
	IAddressManager interface {
		// GetAllCities 获取所有城市
		GetAllCities() []*Region
		// GetRegionList 获取区域信息
		GetRegionList(parentId int) []*Region
		// GetRegions 获取区域名称
		GetRegionNames(code ...int) map[int]string
	}

	// IOptionManager 选项管理器
	IOptionManager interface {
	}
)

// ISystemRepo 系统仓储
type ISystemRepo interface {
	// GetSystemAggregateRoot 获取系统聚合根
	GetSystemAggregateRoot() ISystemAggregateRoot
	// 获取区域仓储
	Region() fw.Repository[Region]
	// GetAllCities 获取所有城市
	GetAllCities() []*Region
	// GetRegionList 获取区域信息
	GetRegionList(parentId int) []*Region
}

type (
	// Region ChinaArea
	Region struct {
		// Code
		Code int `db:"code" pk:"yes" json:"code" bson:"code"`
		// Name
		Name string `db:"name" json:"name" bson:"name"`
		// Parent
		Parent int `db:"parent" json:"parent" bson:"parent"`
	}
)

func (a Region) TableName() string {
	return "sys_region"
}
