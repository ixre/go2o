/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: station.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-08-09 16:24:37
 * description: 提供地方站点的领域模型
 * history:
 */

package sys

import (
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

// 站点聚合根
type IStationDomain interface {
	// 获取领域编号
	domain.IDomain
	// 获取站点
	GetValue() SubStation
	// 设置值
	SetValue(v SubStation) error
	// 保存站点
	Save() error
}

type IStationManager interface {
	// SyncStations 同步所有站
	SyncStations() error
	// FindStationByCity 根据城市代码查找站点,如果为直辖市或区县,则自动查找上级
	FindStationByCity(cityCode int) IStationDomain
}

// 站点仓库
type IStationRepo interface {
	// 仓储基类型
	fw.Repository[SubStation]
	// CreateStation 创建站点
	CreateStation(v *SubStation) IStationDomain
	// GetStation 获取站点
	GetStation(id int) IStationDomain
	// 获取所有的站点
	GetStations() []*SubStation
	// SaveStation 保存站点
	SaveStation(s *SubStation) (int, error)
}

var (
	ErrNoSuchStation = domain.NewError("no_such_station", "没有找到站点:%s")
)

// SysSubStation 地区子站
type SubStation struct {
	// 编号
	Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
	// 城市代码
	CityCode int `json:"cityCode" db:"city_code" gorm:"column:city_code" bson:"cityCode"`
	// 状态: 0: 待开通  1: 已开通  2: 已关闭
	Status int `json:"status" db:"status" gorm:"column:status" bson:"status"`
	// 创建时间
	CreateTime int `json:"createTime" db:"create_time" gorm:"column:create_time" bson:"createTime"`
	// 首字母
	Letter string `json:"letter" db:"letter" gorm:"column:letter" bson:"letter"`
	// 是否热门
	IsHot int `json:"isHot" db:"is_hot" gorm:"column:is_hot" bson:"isHot"`
}

func (s SubStation) TableName() string {
	return "sys_sub_station"
}
