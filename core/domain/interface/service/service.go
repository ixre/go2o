package service

// 领域服务仓储
type IDomainServiceRep interface {
	// Get ProBrand
	GetProBrand(primary interface{}) *ProBrand
	// Save ProBrand
	SaveProBrand(v *ProBrand) (int, error)
	// Delete ProBrand
	DeleteProBrand(primary interface{}) error
	// Select ProBrand
	SelectProBrand(where string, v ...interface{}) []*ProBrand
}
