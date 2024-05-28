package station

import (
	"github.com/ixre/go2o/core/domain/interface/delivery"
	"github.com/ixre/go2o/core/domain/interface/station"
)

var _ station.IStationManager = new(StationManagerImpl)

type StationManagerImpl struct {
	repo         station.IStationRepo
	deliveryRepo delivery.IDeliveryRepo
}

// SyncStations implements station.IStationManager.
func (s *StationManagerImpl) SyncStations() error {
	arr := s.deliveryRepo.GetAllCities()
	stations := s.repo.GetStations()
}
