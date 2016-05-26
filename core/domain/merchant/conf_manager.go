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
	"go2o/core/domain/interface/merchant"
)

var _ merchant.IConfManager = new(ConfManager)

type ConfManager struct {
	_rep        merchant.IMerchantRep
	_merchantId int
}
