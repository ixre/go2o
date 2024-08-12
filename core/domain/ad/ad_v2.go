package ad

import (
	"errors"

	"github.com/ixre/go2o/core/domain/interface/ad"
)

var _ ad.IAdPosition = new(adPositionImpl)

type adPositionImpl struct {
	repo  ad.IAdRepo
	value *ad.Position
}

func NewAdPosition(repo ad.IAdRepo, v *ad.Position) ad.IAdPosition {
	return &adPositionImpl{
		repo:  repo,
		value: v,
	}
}

// PutAd 投放广告
func (a adPositionImpl) PutAd(adId int64) error {
	ia := a.repo.GetAd(adId)
	if ia == nil {
		return ad.ErrNoSuchAd
	}
	a.value.PutAid = int(adId)
	return a.Save()
}

func (a adPositionImpl) GetValue() ad.Position {
	return *a.value
}

func (a adPositionImpl) GetAggregateRootId() int {
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
		a.value.Id = int(id)
	}
	return err
}
