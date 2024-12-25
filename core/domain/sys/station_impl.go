/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: station_impl.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-08-09 20:02:09
 * description: 站点领域模型
 * history:
 */
package sys

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/go2o/core/infrastructure/util"
)

var _ sys.IStationDomain = new(StationImpl)

type StationImpl struct {
	value *sys.SubStation
	repo  sys.IStationRepo
}

// NewStation returns a station aggregate root.
func NewStation(value *sys.SubStation, repo sys.IStationRepo) *StationImpl {
	return &StationImpl{
		value: value,
		repo:  repo,
	}
}

// GetAggregateRootId implements sys.IStationAggregateRoot.
func (s *StationImpl) GetDomainId() int {
	if s.value != nil {
		return s.value.Id
	}
	return 0
}

func (s *StationImpl) GetValue() sys.SubStation {
	return *types.DeepClone(s.value)
}

// SetValue implements sys.IStationAggregateRoot.
func (s *StationImpl) SetValue(v sys.SubStation) error {
	if s.GetDomainId() <= 0 {
		if v.CityCode <= 0 {
			return errors.New("invalid city code")
		}
	}
	s.value.Status = v.Status
	return nil
}

// Save implements sys.IStationAggregateRoot.
func (s *StationImpl) Save() error {
	id, err := s.repo.SaveStation(s.value)
	if s.GetDomainId() == 0 {
		s.value.Id = id
	}
	return err
}

var _ sys.IStationManager = new(stationManagerImpl)

type stationManagerImpl struct {
	repo        sys.IStationRepo
	sysRepo     sys.ISystemRepo
	stationList []*sys.SubStation
	sysRootImpl *systemAggregateRootImpl
}

func NewStationManager(repo sys.IStationRepo, sysRepo sys.ISystemRepo, sysRootImpl *systemAggregateRootImpl) sys.IStationManager {
	return &stationManagerImpl{
		repo:        repo,
		sysRepo:     sysRepo,
		sysRootImpl: sysRootImpl,
	}
}

// SyncStations implements sys.IStationManager.
func (s *stationManagerImpl) SyncStations() error {
	is := s.sysRepo.GetSystemAggregateRoot()
	arr := is.Location().GetAllCities()
	stations := s.repo.GetStations()
	syncArray := make([]*sys.District, 0)
	for _, v := range arr {
		exists := false
		for _, s := range stations {
			if s.CityCode == v.Code {
				exists = true
				break
			}
		}
		if !exists {
			syncArray = append(syncArray, v)
		}
	}
	if l := len(syncArray); l > 0 {
		log.Printf("[ Go2o][ Info]: will sync %d stations \n", l)
		for _, v := range syncArray {
			s.createSubStation(v)
		}
	}
	// 清空站点列表
	s.stationList = nil
	s.getStationList()
	return nil
}

func (s *stationManagerImpl) createSubStation(city *sys.District) {
	i := s.repo.CreateStation(&sys.SubStation{
		CityCode:   city.Code,
		Status:     0,
		Letter:     util.GetHansFirstLetter(city.Name),
		IsHot:      s.isHot(city.Name),
		CreateTime: int(time.Now().Unix()),
	})
	i.Save()
}

// 是否为热门城市
func (s *stationManagerImpl) isHot(name string) int {
	v := []rune(name)[0:2]
	cityName := string(v)
	in := collections.InArray([]string{"北京", "上海", "广州", "深圳", "深圳", "佛山", "厦门", "重庆", "杭州"}, cityName)
	return types.Ternary(in, 1, 0)
}

// getStationList 获取站点列表
func (s *stationManagerImpl) getStationList() []*sys.SubStation {
	if s.stationList == nil {
		s.stationList = fw.ReduceFinds(func(opt *fw.QueryOption) []*sys.SubStation {
			return s.repo.FindList(opt, "")
		}, 1000)
	}
	return s.stationList
}

// FindStationByCity 查找站点
func (s *stationManagerImpl) FindStationByCity(cityCode int) sys.IStationDomain {
	if cityCode <= 0 {
		return nil
	}
	arr := s.getStationList()
	// 查找城市站点
	dst := collections.FindArray(arr, func(v *sys.SubStation) bool {
		return v.CityCode == cityCode
	})
	if dst != nil {
		return s.repo.GetStation(dst.Id)
	}
	// 查找直辖市站点
	d := s.sysRootImpl.Location().GetDistrict(cityCode)
	if d == nil {
		// 地区不存在
		return nil
	}
	if !strings.Contains(d.Name, "市辖区") &&
		!strings.Contains(d.Name, "市辖县") &&
		!strings.Contains(d.Name, "县") {
		// 非直辖市和区县,则返回空
		return nil
	}
	// 查找直辖市站点
	provinceCode := d.Parent
	dst = collections.FindArray(arr, func(v *sys.SubStation) bool {
		return v.CityCode == provinceCode
	})
	if dst != nil {
		return s.repo.GetStation(dst.Id)
	}
	return nil
}
