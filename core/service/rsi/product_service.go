package rsi

import "go2o/core/domain/interface/pro_model"

// 产品服务
type productService struct {
	pmRep promodel.IProModelRepo
}

func NewProService(pmRep promodel.IProModelRepo) *productService {
	return &productService{
		pmRep: pmRep,
	}
}
