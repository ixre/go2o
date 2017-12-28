package merchant

import "go2o/core/domain/interface/merchant"

var _ merchant.ISaleManager = new(SaleManagerImpl)

type SaleManagerImpl struct {
	mch   *merchantImpl
	mchId int
}

func newSaleManagerImpl(id int, m *merchantImpl) merchant.ISaleManager {
	return &SaleManagerImpl{
		mchId: id,
		mch:   m,
	}
}

func (s *SaleManagerImpl) MathTradeFee(tradeType int, amount int) (int, error) {
	panic("implement me")
}
