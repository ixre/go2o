package station

import (
	"log"
	"time"

	"github.com/ixre/go2o/core/domain/interface/station"
	"github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/go2o/core/infrastructure/util/collections"
	"github.com/ixre/go2o/core/infrastructure/util/types"
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
	syncArray := make([]*station.Area, 0)
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
	return nil
}

func (s *stationManagerImpl) createSubStation(city *station.Area) {
	i := s.repo.CreateStation(&station.SubStation{
		CityCode:   city.Code,
		Status:     0,
		Letter:     util.GetHansFirstLetter(city.Name),
		IsHot:      s.isHot(city.Name),
		CreateTime: time.Now().Unix(),
	})
	i.Save()
}

// 是否为热门城市
func (s *stationManagerImpl) isHot(name string) int {
	v := name[0:2]
	in := collections.InArray([]string{"北京", "上海", "广州", "深圳", "深圳", "佛山", "厦门", "重庆", "杭州"}, v)
	return types.Ternary(in, 1, 0)
}
