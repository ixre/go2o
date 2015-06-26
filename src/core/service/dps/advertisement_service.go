/**
 * Copyright 2015 @ S1N1 Team.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dps
import (
	"go2o/src/core/domain/interface/ad"
)


type advertisementService struct {
	_rep ad.IAdvertisementRep
	//_query     *query.ContentQuery
}

func NewAdvertisementService(rep ad.IAdvertisementRep) *advertisementService {
	return &advertisementService{
		_rep: rep,
	}
}

func (this *advertisementService) GetAdvertisement(partnerId,id int)*ad.ValueAdvertisement{
	pa := this._rep.GetPartnerAdvertisement(partnerId)
	if adv := pa.GetById(id);adv != nil{
		return adv.GetValue()
	}
	return nil
}