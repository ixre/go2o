package parser

import (
	"go2o/core/domain/interface/merchant"
	"go2o/core/service/auto_gen/rpc/mch_service"
)

func TradeConfDto(conf *merchant.TradeConf) *mch_service.STradeConf {
	return &mch_service.STradeConf{
		MchId:       conf.MchId,
		TradeType:   int32(conf.TradeType),
		PlanId:      conf.PlanId,
		Flag:        int32(conf.Flag),
		AmountBasis: int32(conf.AmountBasis),
		TradeFee:    int32(conf.TradeFee),
		TradeRate:   int32(conf.TradeRate),
	}
}

func TradeConf(conf *mch_service.STradeConf) *merchant.TradeConf {
	return &merchant.TradeConf{
		MchId:       conf.MchId,
		TradeType:   int(conf.TradeType),
		PlanId:      conf.PlanId,
		Flag:        int(conf.Flag),
		AmountBasis: int(conf.AmountBasis),
		TradeFee:    int(conf.TradeFee),
		TradeRate:   int(conf.TradeRate),
	}
}
