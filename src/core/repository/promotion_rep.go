/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:12
 * description :
 * history :
 */

package repository

import (
	"github.com/atnet/gof/db"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/promotion"
)

var _ promotion.IPromotionRep = new(promotionRep)

type promotionRep struct {
	db.Connector
	_memberRep member.IMemberRep
}

func NewPromotionRep(c db.Connector, memberRep member.IMemberRep) promotion.IPromotionRep {
	return &promotionRep{
		Connector:  c,
		_memberRep: memberRep,
	}
}