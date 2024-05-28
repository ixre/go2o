package station

import "github.com/ixre/go2o/core/domain/interface/station"

var _ station.IStationAggregateRoot = new(StationImpl)

type StationImpl struct {
	value *station.Station
	repo  station.IStationRepo
}

// GetAggregateRootId implements station.IStationAggregateRoot.
func (s *StationImpl) GetAggregateRootId() int {
	if s.value != nil {
		return s.value.Id
	}
	return 0
}

// Save implements station.IStationAggregateRoot.
func (s *StationImpl) Save() error {
	id, err := s.repo.SaveStation(s.value)
	if s.GetAggregateRootId() == 0 {
		s.value.Id = id
	}
	return err
}
