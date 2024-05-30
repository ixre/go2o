package station

import (
	"time"

	"github.com/ixre/go2o/core/domain/interface/station"
)

var _ station.IStationManager = new(stationManagerImpl)

type stationManagerImpl struct {
	repo station.IStationRepo
}

func NewStationManager(repo station.IStationRepo) station.IStationManager {
	return &stationManagerImpl{
		repo: repo,
	}
}

// SyncStations implements station.IStationManager.
func (s *stationManagerImpl) SyncStations() error {
	arr := s.repo.GetAllCities()
	stations := s.repo.GetStations()
	for _, v := range arr {
		exists := false
		for _, s := range stations {
			if s.CityCode == v.Code {
				exists = true
				break
			}
		}
		if !exists {
			s.createSubStation(v)
		}
	}
	return nil
}

func (s *stationManagerImpl) createSubStation(city *station.Area) {
	i := s.repo.CreateStation(&station.SubStation{
		CityCode:   city.Code,
		Status:     0,
		CreateTime: time.Now().Unix(),
	})
	i.Save()
}
