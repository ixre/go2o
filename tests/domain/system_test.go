package domain

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

func TestGetOptions(t *testing.T) {
	ia := inject.GetSystemRepo().GetSystemAggregateRoot() // TODO: write test code here
	arr := ia.Options().GetChildOptions(0, "BIZ")
	t.Logf("options = %#v \n", arr)
}

func TestGetAllCities(t *testing.T) {
	ia := inject.GetSystemRepo().GetSystemAggregateRoot() // TODO: write test code here
	cities := ia.Address().GetAllCities()
	if len(cities) == 0 {
		t.Error("No cities found")
	}
	for _, city := range cities {
		t.Log(city)
	}
}

// 测试根据城市获取站点
func TestGetStationByCityCode(t *testing.T) {
	city := 110100
	ia := inject.GetSystemRepo().GetSystemAggregateRoot() // TODO: write test code here
	d := ia.Address().GetDistrict(city)
	t.Logf("district = %s \n", d.Name)
	station := ia.Stations().FindStationByCity(city)
	t.Logf("station = %#v \n", station.GetValue())
}
