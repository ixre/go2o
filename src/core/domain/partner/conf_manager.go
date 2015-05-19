/**
 * Copyright 2015 @ S1N1 Team.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package partner

import (
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/domain/interface/valueobject"
)

var _ partner.IConfManager = new(ConfManager)

type ConfManager struct {
	_rep       partner.IPartnerRep
	_partnerId int
	_levelSet  []*valueobject.MemberLevel
}
