package sys

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/util/lbs"
)

var _ sys.ISystemAggregateRoot = new(systemAggregateRootImpl)

type systemAggregateRootImpl struct {
	_address  sys.ILocationManager
	_options  sys.IOptionManager
	_stations sys.IStationManager
	_log      sys.IApplicationManager
	_repo     sys.ISystemRepo
	_oauth    sys.IOAuthManager
}

func NewSystemAggregateRoot(repo sys.ISystemRepo) sys.ISystemAggregateRoot {
	return &systemAggregateRootImpl{_repo: repo}
}

// GetAggregateRootId implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) GetAggregateRootId() int {
	return 1
}

// Address implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) Location() sys.ILocationManager {
	if s._address == nil {
		s._address = &locationManagerImpl{
			s._repo.District(), nil}
	}
	return s._address
}

// Options implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) Options() sys.IOptionManager {
	if s._options == nil {
		s._options = &optionManagerImpl{
			Repository: s._repo.Option(),
		}
	}
	return s._options
}

// Stations implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) Stations() sys.IStationManager {
	if s._stations == nil {
		s._stations = NewStationManager(s._repo.Station(), s._repo, s)
	}
	return s._stations
}

// Application implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) Application() sys.IApplicationManager {
	if s._log == nil {
		s._log = newApplicationManager(s._repo)
	}
	return s._log
}

// OAuth 获取OAuth管理器
func (s *systemAggregateRootImpl) OAuth() sys.IOAuthManager {
	if s._oauth == nil {
		s._oauth = newOAuthManager()
	}
	return s._oauth
}

// FlushUpdateStatus implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) FlushUpdateStatus() {
	s._repo.FlushUpdateStatus()
}

// LastUpdateTime implements sys.ISystemAggregateRoot.
func (s *systemAggregateRootImpl) LastUpdateTime() int64 {
	return s._repo.LastUpdateTime()
}

// GetBanks 获取银行列表
func (s *systemAggregateRootImpl) GetBanks() []*sys.GeneralOption {
	return sys.BankCodes
}

var _ sys.ILocationManager = new(locationManagerImpl)

type locationManagerImpl struct {
	fw.Repository[sys.District]
	districtList []*sys.District
}

// FindCity 查找城市
func (a *locationManagerImpl) FindCity(name string) *sys.District {
	return collections.FindArray(a.GetAllCities(), func(d *sys.District) bool {
		return d.Name == name
	})
}

// GetDistrict 获取区域信息
func (a *locationManagerImpl) GetDistrict(id int) *sys.District {
	return collections.FindArray(a.getDistrictList(), func(d *sys.District) bool {
		return d.Code == id
	})
}

// getDistrictList 获取地区列表
func (a *locationManagerImpl) getDistrictList() []*sys.District {
	if a.districtList == nil {
		a.districtList = fw.ReduceFinds(func(opt *fw.QueryOption) []*sys.District {
			return a.FindList(opt, "")
		}, 1000)
	}
	return a.districtList
}

// getProvinces 获取省列表
func (a *locationManagerImpl) getProvinces() []*sys.District {
	return a.GetChildrenDistricts(0)
}

// GetDistrictNames implements sys.IAddressManager.
func (a *locationManagerImpl) GetDistrictNames(code ...int) map[int]string {
	mp := make(map[int]string)
	for _, v := range a.getDistrictList() {
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
func (a *locationManagerImpl) GetAllCities() []*sys.District {
	provinceList := a.getProvinces()
	provinceCodes := collections.MapList(provinceList,
		func(s *sys.District) int {
			return s.Code
		})
	cityList := collections.FilterArray(a.getDistrictList(), func(s *sys.District) bool {
		return s.Parent != 0 && collections.AnyArray(provinceCodes, func(c int) bool {
			return c == s.Parent
		})
	})
	ret := make([]*sys.District, 0)
	for _, c := range cityList {
		c.Name = strings.TrimSpace(c.Name)
		if c.Name == "市辖区" || c.Name == "县" || c.Name == "区" {
			// 将直辖市加入到城市列表中
			if collections.AnyArray(ret, func(a *sys.District) bool {
				return a.Code == c.Parent
			}) {
				// 已添加直辖市到城市列表中
				continue
			}
			parent := collections.FindArray(provinceList, func(a *sys.District) bool {
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

// GetChildrenDistricts implements sys.IAddressManager.
func (a *locationManagerImpl) GetChildrenDistricts(parentId int) []*sys.District {
	return collections.FilterArray(a.getDistrictList(), func(a *sys.District) bool {
		return a.Parent == parentId && a.Code != 0
	})
}

// FindRegionByIp 根据IP查找区域信息
func (a *locationManagerImpl) FindRegionByIp(ip string) (*sys.District, error) {
	if strings.HasPrefix(ip, "127.") ||
		strings.HasPrefix(ip, "172.") ||
		strings.HasPrefix(ip, "192.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "0.") ||
		strings.HasPrefix(ip, "::1") {
		// 本地网络IP不更新位置
		return &sys.District{Code: 0, Name: "局域网", Parent: 0}, nil
	}
	provider := lbs.GetProvider()
	if provider == nil {
		return nil, errors.New("未配置位置服务提供者")
	}
	info, err := provider.QueryLocation(ip)
	if err != nil {
		return nil, err
	}
	return a.FindCity(info.City), nil
}

var _ sys.IOptionManager = new(optionManagerImpl)

type optionManagerImpl struct {
	fw.Repository[sys.GeneralOption]
	allList []*sys.GeneralOption
	rwlock  sync.RWMutex
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

// getOption 获取选项
func (o *optionManagerImpl) getOption(id int) *sys.GeneralOption {
	return collections.FindArray(o.getList(), func(s *sys.GeneralOption) bool {
		return s.Id == id
	})
}

// getDistrictList 获取地区列表
func (o *optionManagerImpl) getList() []*sys.GeneralOption {
	o.rwlock.RLock()
	if o.allList == nil {
		o.rwlock.RUnlock()
		o.rwlock.Lock()
		o.allList = fw.ReduceFinds(func(opt *fw.QueryOption) []*sys.GeneralOption {
			return o.FindList(opt, "")
		}, 1000)
		o.rwlock.Unlock()
	} else {
		o.rwlock.RUnlock()
	}
	return o.allList
}

// IsLeaf implements sys.IOptionManager.
func (o *optionManagerImpl) IsLeaf(g *sys.GeneralOption) bool {
	return collections.FindArray(o.getList(), func(n *sys.GeneralOption) bool {
		return n.Pid == g.Id && n.Enabled == 1
	}) == nil
}

func (o *optionManagerImpl) SaveOption(option *sys.GeneralOption) error {
	if option.Pid != 0 {
		pidOption := o.getOption(option.Pid)
		if pidOption == nil {
			return errors.New("上级节点不存在")
		}
	}
	if option.SortNum == 0 {
		option.SortNum = o.getMaxSortNum(option.Pid) + 1
	}
	if len(option.Label) == 0 {
		return errors.New("标签不能为空")
	}
	if option.Pid == 0 && len(option.Type) == 0 {
		return errors.New("类型不能为空")
	}
	var err error
	if option.Id > 0 {
		// 更新
		origin := o.getOption(option.Id)
		if origin == nil {
			return errors.New("选项数据不存在")
		}
		origin.Enabled = option.Enabled
		origin.Label = option.Label
		origin.Value = option.Value
		origin.SortNum = option.SortNum
		_, err = o.Save(origin)
	} else {
		option.CreateTime = int(time.Now().Unix())
		// 新增
		_, err = o.Save(option)
	}
	if err == nil {
		o.flushCache()
	}
	return err
}

// getMaxSortNum 获取最大排序号
func (o *optionManagerImpl) getMaxSortNum(pid int) int {
	var max int
	for _, v := range o.getList() {
		if v.Pid == pid {
			if v.SortNum > max {
				max = v.SortNum
			}
		}
	}
	return max
}

// Delete 删除选项
func (o *optionManagerImpl) Delete(option *sys.GeneralOption) error {
	var opt *sys.GeneralOption
	for _, v := range o.getList() {
		if v.Pid == option.Id {
			return errors.New("存在下级选项,无法删除")
		}
		if v.Id == option.Id {
			opt = v
		}
	}
	if opt == nil {
		return errors.New("no such data")
	}
	err := o.Repository.Delete(opt)
	if err == nil {
		o.flushCache()
	}
	return err
}

// flushCache 刷新缓存
func (o *optionManagerImpl) flushCache() {
	o.rwlock.Lock()
	o.allList = nil
	o.rwlock.Unlock()
	// 重新加载
	o.getList()
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
			mp[v.Id] = v.Label
		}
	}
	return mp
}
