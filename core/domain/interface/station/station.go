package station

// 站点聚合根
type IStationAggregateRoot interface {
	// 获取聚合根编号
	GetAggregateRootId() int
	// 设置值
	SetValue(v SubStation) error
	// 保存站点
	Save() error
}

type IStationManager interface {
	// SyncStations 同步所有站
	SyncStations() error
}

// 站点仓库
type IStationRepo interface {
	// 获取站点管理器
	GetManager() IStationManager
	// CreateStation 创建站点
	CreateStation(v *SubStation) IStationAggregateRoot
	// GetStation 获取站点
	GetStation(id int) IStationAggregateRoot
	// 获取所有的站点
	GetStations() []*SubStation
	// SaveStation 保存站点
	SaveStation(s *SubStation) (int, error)
	// GetAllCities 获取所有城市
	GetAllCities() []*Area
	// GetAreaList 获取区域信息
	GetAreaList(parentId int) []*Area
}

// SubStation 地区子站
type SubStation struct {
	// 编号
	Id int `db:"id" pk:"yes" auto:"yes" json:"id" bson:"id"`
	// 城市代码
	CityCode int `db:"city_code" json:"cityCode" bson:"cityCode"`
	// 状态: 0: 待开通  1: 已开通  2: 已关闭
	Status int `db:"status" json:"status" bson:"status"`
	// 首字母
	Letter string `db:"letter" json:"letter" bson:"letter"`
	// 是否热门
	IsHot int `db:"is_hot" json:"isHot" bson:"isHot"`
	// 创建时间
	CreateTime int64 `db:"create_time" json:"createTime" bson:"createTime"`
}

// Area ChinaArea
type Area struct {
	// Code
	Code int `db:"code" pk:"yes" json:"code" bson:"code"`
	// Name
	Name string `db:"name" json:"name" bson:"name"`
	// Parent
	Parent int `db:"parent" json:"parent" bson:"parent"`
}
