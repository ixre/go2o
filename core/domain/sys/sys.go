package sys

import (
	"strings"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/util/collections"
)

var _ sys.ISystemAggregateRoot = new(systemAggregateRootImpl)

type systemAggregateRootImpl struct {
	_address sys.IAddressManager
	_options sys.IOptionManager
	_repo    sys.ISystemRepo
}

func NewSystemAggregateRoot(repo sys.ISystemRepo) sys.ISystemAggregateRoot {
	return &systemAggregateRootImpl{_repo: repo}
}

// GetAggregateRootId implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) GetAggregateRootId() int {
	return 1
}

// Address implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) Address() sys.IAddressManager {
	if s._address == nil {
		s._address = &addressManagerImpl{
			s._repo.Region(), nil}
	}
	return s._address
}

// Options implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) Options() sys.IOptionManager {
	if s._options == nil {
		//s._options = &optionManagerImpl{s._repo.Options()}
	}
	return s._options
}

var _ sys.IAddressManager = new(addressManagerImpl)

type addressManagerImpl struct {
	fw.Repository[sys.Region]
	areaList []*sys.Region
}

// getAreaList 获取地区列表
func (a *addressManagerImpl) getAreaList() []*sys.Region {
	if a.areaList == nil {
		a.areaList = a.FindList("")
	}
	return a.areaList
}

// getProvinces 获取省列表
func (a *addressManagerImpl) getProvinces() []*sys.Region {
	return collections.FilterArray(a.getAreaList(), func(a *sys.Region) bool {
		return a.Parent == 0 && a.Code != 0
	})
}

// GetAllCities 获取所有城市列表
func (a *addressManagerImpl) GetAllCities() []*sys.Region {
	provinceList := a.getProvinces()
	provinceCodes := collections.MapList(provinceList,
		func(s *sys.Region) int {
			return s.Code
		})
	cityList := collections.FilterArray(a.getAreaList(), func(s *sys.Region) bool {
		return s.Parent != 0 && collections.AnyArray(provinceCodes, func(c int) bool {
			return c == s.Parent
		})
	})
	ret := make([]*sys.Region, 0)
	for _, c := range cityList {
		c.Name = strings.TrimSpace(c.Name)
		if c.Name == "市辖区" || c.Name == "县" || c.Name == "区" {
			// 将直辖市加入到城市列表中
			if collections.AnyArray(ret, func(a *sys.Region) bool {
				return a.Code == c.Parent
			}) {
				// 已添加直辖市到城市列表中
				continue
			}
			parent := collections.FindArray(provinceList, func(a *sys.Region) bool {
				return a.Code == c.Parent
			})
			if parent != nil {
				ret = append(ret, parent)
			}
		} else {
			ret = append(ret, c)
		}
	}
	return ret
}

// GetRegionList implements sys.IAddressManager.
func (a *addressManagerImpl) GetRegionList(parentId int) []*sys.Region {
	return a.FindList("parent=$1", parentId)
}

var _ sys.IOptionManager = new(optionManagerImpl)

type optionManagerImpl struct {
}
