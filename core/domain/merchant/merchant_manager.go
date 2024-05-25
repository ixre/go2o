package merchant

import (
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
)

var _ merchant.IMerchantManager = new(merchantManagerImpl)

type merchantManagerImpl struct {
	rep     merchant.IMerchantRepo
	valRepo valueobject.IValueRepo
}

func NewMerchantManager(rep merchant.IMerchantRepo,
	valRepo valueobject.IValueRepo) merchant.IMerchantManager {
	return &merchantManagerImpl{
		rep:     rep,
		valRepo: valRepo,
	}
}

// GetMerchantByMemberId 获取会员关联的商户
func (m *merchantManagerImpl) GetMerchantByMemberId(memberId int) merchant.IMerchant {
	return m.rep.GetMerchantByMemberId(memberId)
}
