package dps

import (
	"go2o/core/domain/interface/delivery"
)

type deliveryService struct {
	_rep delivery.IDeliveryRep
}

func NewDeliveryService(r delivery.IDeliveryRep) *deliveryService {
	return &deliveryService{
		_rep: r,
	}
}

func (this *deliveryService) CreateCoverageArea(c *delivery.CoverageValue) (int, error) {
	return this._rep.SaveCoverageArea(c)
}
