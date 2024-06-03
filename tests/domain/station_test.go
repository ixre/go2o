package domain

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

func TestGetAllCities(t *testing.T) {
	cities := inject.GetStationRepo().GetAllCities()
	if len(cities) == 0 {
		t.Error("No cities found")
	}
	for _, city := range cities {
		t.Log(city)
	}
}

func TestGetGroupStations(t *testing.T) {
	cities := inject.GetStationQueryService().QueryGroupStations(0)
	if len(cities) == 0 {
		t.Error("No cities found")
	}
	for _, city := range cities {
		t.Log(city)
	}
}

// 测试同步站点
func TestSyncStations(t *testing.T) {
	err := inject.GetStationRepo().GetManager().SyncStations()
	if err != nil {
		t.Error(err)
	}
}
