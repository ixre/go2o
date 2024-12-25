package domain

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

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
	err := inject.GetSystemRepo().GetSystemAggregateRoot().Stations().SyncStations()
	if err != nil {
		t.Error(err)
	}
}
