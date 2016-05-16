/**
 * Copyright 2015 @ z3q.net.
 * name : conf_manager
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package merchant

import (
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/domain/interface/valueobject"
)

var _ merchant.IConfManager = new(ConfManager)

type ConfManager struct {
	_rep        merchant.IMerchantRep
	_merchantId int
	_levelSet   []*valueobject.MemberLevel
}
