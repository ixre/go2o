package merchant

import (
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/merchant"
)

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

// 计算交易手续费
func (s *SaleManagerImpl) MathTradeFee(tradeType int, amount int) (int, error) {
	cm := s.mch.ConfManager()
	conf := cm.GetTradeConf(tradeType)
	if conf == nil {
		//todo: 应使用系统默认的比例进行手续费
		return 0, nil
	}
	// 免费
	if conf.Flag&merchant.TFlagFree == merchant.TFlagFree {
		return 0, nil
	}
	switch conf.AmountBasis {
	case enum.AmountBasisNotSet: // 免费
		return 0, nil
	case enum.AmountBasisByAmount: // 按订单单数，收取金额
		return conf.TradeFee, nil
	case enum.AmountBasisByPercent: // 按订单金额，收取百分比
		return int(float64(amount*conf.TradeRate) / enum.RATE_PERCENT), nil
	default:
		panic("not support amount basis")
	}
	return 0, nil
}
