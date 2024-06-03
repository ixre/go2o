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
