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
		s._options = &optionManagerImpl{s._repo.Option(), nil}
	}
	return s._options
}

// FlushUpdateStatus implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) FlushUpdateStatus() {
	s._repo.FlushUpdateStatus()
}

// LastUpdateTime implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) LastUpdateTime() int64 {
	return s._repo.LastUpdateTime()
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

// GetRegionNames implements sys.IAddressManager.
func (a *addressManagerImpl) GetRegionNames(code ...int) map[int]string {
	mp := make(map[int]string)
	for _, v := range a.getAreaList() {
		if len(mp) == len(code) {
			break
		}
		if collections.AnyArray(code, func(c int) bool {
			return c == v.Code
		}) {
			mp[v.Code] = v.Name
		}
	}
	return mp
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
	fw.Repository[sys.GeneralOption]
	allList []*sys.GeneralOption
}

// GetChildOptions implements sys.IOptionManager.
func (o *optionManagerImpl) GetChildOptions(parentId int, typeName string) []*sys.GeneralOption {
	l := len(typeName)
	if parentId == 0 && l == 0 {
		// 无法根据参数获取数据
		return []*sys.GeneralOption{}
	}
	if parentId == 0 && l > 0 {
		// 返回顶级节点的下级数据
		t := collections.FindArray(o.getList(), func(s *sys.GeneralOption) bool {
			return s.Type == typeName
		})
		if t == nil {
			return []*sys.GeneralOption{}
		}
		parentId = t.Id
		typeName = ""
	}
	return collections.FilterArray(o.getList(), func(s *sys.GeneralOption) bool {
		return s.Pid == parentId && (s.Type == typeName || typeName == "")
	})
}

// getAreaList 获取地区列表
func (o *optionManagerImpl) getList() []*sys.GeneralOption {
	if o.allList == nil {
		o.allList = o.FindList("")
	}
	return o.allList
}

// GetOptionNames implements sys.IOptionManager.
func (o *optionManagerImpl) GetOptionNames(code ...int) map[int]string {
	mp := make(map[int]string)
	for _, v := range o.getList() {
		if len(mp) == len(code) {
			break
		}
		if collections.AnyArray(code, func(c int) bool {
			return c == v.Id
		}) {
			mp[v.Id] = v.Name
		}
	}
	return mp
}
