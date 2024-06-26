/**
 * This file is auto generated by tto v0.5.7 !
 * If you want to modify this code, please read the guide
 * to modify code template.
 *
 * Get started: https://github.com/ixre/tto
 *
 * Copyright (C) 2009-2024 56X.NET, All rights reserved.
 *
 * name : sys_sub_station_dao_impl.go
 * author : root
 * date : 2024/05/29 17:31:31
 * description :
 * history :
 */
package repos

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/ixre/go2o/core/domain/interface/station"
	"github.com/ixre/go2o/core/domain/interface/sys"
	si "github.com/ixre/go2o/core/domain/station"
	"github.com/ixre/gof/db/orm"
)

var _ station.IStationRepo = new(stationRepoImpl)

// 站点仓储
type stationRepoImpl struct {
	_orm     orm.Orm
	_manager station.IStationManager
	sysRepo  sys.ISystemRepo
}

var modelIsMapped = false

// NewStationRepo Create new StationRepo
func NewStationRepo(o orm.Orm, sysRepo sys.ISystemRepo) station.IStationRepo {
	if !modelIsMapped {
		_ = o.Mapping(station.SubStation{}, "sys_sub_station")
		modelIsMapped = true
	}
	return &stationRepoImpl{
		_orm:    o,
		sysRepo: sysRepo,
	}
}

// CreateStation implements station.IStationRepo.
func (s *stationRepoImpl) CreateStation(v *station.SubStation) station.IStationAggregateRoot {
	return si.NewStation(v, s)
}

// GetManager implements station.IStationRepo.
func (s *stationRepoImpl) GetManager() station.IStationManager {
	if s._manager == nil {
		s._manager = si.NewStationManager(s, s.sysRepo)
	}
	return s._manager
}

// GetSubStation Get 地区子站
func (s *stationRepoImpl) GetStation(id int) station.IStationAggregateRoot {
	e := station.SubStation{}
	err := s._orm.Get(id, &e)
	if err == nil {
		return s.CreateStation(&e)
	}
	if err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:SubStation\n", err.Error())
	}
	return nil
}

// GetStationByCity implements station.IStationRepo.
func (s *stationRepoImpl) GetStationByCity(cityCode int) station.IStationAggregateRoot {
	e := station.SubStation{}
	err := s._orm.GetBy(&e, "city_code = $1", strconv.Itoa(cityCode))
	if err == nil {
		return s.CreateStation(&e)
	}
	if err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:SubStation\n", err.Error())
	}
	return nil
}

// GetStations implements station.IStationRepo.
func (s *stationRepoImpl) GetStations() []*station.SubStation {
	list := make([]*station.SubStation, 0)
	err := s._orm.Select(&list, "")
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:SubStation\n", err.Error())
	}
	return list
}

// SaveSubStation Save 地区子站
func (s *stationRepoImpl) SaveStation(v *station.SubStation) (int, error) {
	id, err := orm.Save(s._orm, v, int(v.Id))
	if err != nil && err != sql.ErrNoRows {
		log.Printf("[ Orm][ Error]: %s; Entity:SubStation\n", err.Error())
	}
	return id, err
}
