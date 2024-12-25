package lbs

var lbsProvider LbsProvider

// 位置信息
type LocationInfo struct {
	// 纬度
	Lat float64 `json:"lat"`
	// 经度
	Lng float64 `json:"lng"`
	// 地址
	Location string `json:"location"`
	// 城市
	City string `json:"city"`
	// 位置代码
	PlaceCode int `json:"code"`
}

// LbsProvider 位置服务提供者接口
type LbsProvider interface {
	// 根据IP地址查询位置信息
	QueryLocation(ip string) (*LocationInfo, error)
}

// GetProvider 获取位置服务提供者实例
func GetProvider() LbsProvider {
	return lbsProvider
}

// Configure 配置位置服务提供者
func Configure(p LbsProvider) {
	lbsProvider = p
}
