package station

import (
	"log"
	"time"

	"github.com/ixre/go2o/core/domain/interface/station"
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/go2o/core/infrastructure/util/collections"
	"github.com/ixre/go2o/core/infrastructure/util/types"
)

var _ station.IStationManager = new(stationManagerImpl)

type stationManagerImpl struct {
	repo    station.IStationRepo
	sysRepo sys.ISystemRepo
}

func NewStationManager(repo station.IStationRepo, sysRepo sys.ISystemRepo) station.IStationManager {
	return &stationManagerImpl{
		repo:    repo,
		sysRepo: sysRepo,
	}
}

// SyncStations implements station.IStationManager.
func (s *stationManagerImpl) SyncStations() error {
	is := s.sysRepo.GetSystemAggregateRoot()
	arr := is.Address().GetAllCities()
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
	return nil
}

func (s *stationManagerImpl) createSubStation(city *sys.District) {
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
	v := []rune(name)[0:2]
	cityName := string(v)
	in := collections.InArray([]string{"北京", "上海", "广州", "深圳", "深圳", "佛山", "厦门", "重庆", "杭州"}, cityName)
	return types.Ternary(in, 1, 0)
}
