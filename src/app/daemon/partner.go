/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
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
