/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"go2o/src/core/service/dps"
)

var (
	partnerIds []int
)

func getPartners() []int {
	if partnerIds == nil {
		partnerIds = dps.PartnerService.GetPartnersId()
	}
	return partnerIds
}


