package station

type IStationAggregateRoot interface {
	// 获取聚合根编号
	GetAggregateRootId() int
	// 保存站点
	Save() error
}

type IStationManager interface {
	// SyncStations 同步所有站
	SyncStations() error
}

type IStationRepo interface {
	// GetStation 获取站点
	GetStation(id int) IStationAggregateRoot
	// SaveStation 保存站点
	SaveStation(s *Station) (int, error)
}

// 站点
type Station struct {
	// 编号
	Id int `db:"id" pk:"yes" auto:"yes"`
	// 城市编码
	CityCode int `db:"city_code"`
}
