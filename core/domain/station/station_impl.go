package station

import (
	"errors"

	"github.com/ixre/go2o/core/domain/interface/station"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

var _ station.IStationAggregateRoot = new(StationImpl)

type StationImpl struct {
	value *station.SubStation
	repo  station.IStationRepo
}

// NewStation returns a station aggregate root.
func NewStation(value *station.SubStation, repo station.IStationRepo) *StationImpl {
	return &StationImpl{
		value: value,
		repo:  repo,
	}
}

// GetAggregateRootId implements station.IStationAggregateRoot.
func (s *StationImpl) GetAggregateRootId() int {
	if s.value != nil {
		return s.value.Id
	}
	return 0
}

func (s *StationImpl) GetValue() station.SubStation {
	return *types.DeepClone(s.value)
}

// SetValue implements station.IStationAggregateRoot.
func (s *StationImpl) SetValue(v station.SubStation) error {
	if s.GetAggregateRootId() <= 0 {
		if v.CityCode <= 0 {
			return errors.New("invalid city code")
		}
	}
	s.value.Status = v.Status
	return nil
}

// Save implements station.IStationAggregateRoot.
func (s *StationImpl) Save() error {
	id, err := s.repo.SaveStation(s.value)
	if s.GetAggregateRootId() == 0 {
		s.value.Id = id
	}
	return err
}
