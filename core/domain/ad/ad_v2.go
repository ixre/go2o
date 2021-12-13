package ad

import (
	"errors"
	"go2o/core/domain/interface/ad"
)

type adPositionImpl struct {
	repo  ad.IAdRepo
	value *ad.Position
}

func (a adPositionImpl) GetAggregateRootId() int64 {
	return a.value.Id
}

func (a adPositionImpl) SetValue(v *ad.Position) error {
	if len(v.Name) == 0 {
		return errors.New("name is empty")
	}
	a.value.GroupName = v.GroupName
	a.value.Key = v.Key
	a.value.Name = v.Name
	return nil
}

func (a adPositionImpl) Save() error {
	id, err := a.repo.SaveAdPosition(a.value)
	if err == nil {
		a.value.Id = id
	}
	return err
}

func NewAdPosition(repo ad.IAdRepo, v *ad.Position) ad.IAdPosition {
	return &adPositionImpl{
		repo:  repo,
		value: v,
	}
}
