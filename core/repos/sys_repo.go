package repos

import (
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/sys"
	impl "github.com/ixre/go2o/core/domain/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/logger"
	"github.com/ixre/gof/storage"
)

var (
	_sysRepo      = new(systemRepoImpl)
	_sysAggregate sys.ISystemAggregateRoot
	_repoOnce     sync.Once
)

type systemRepoImpl struct {
	fw.ORM
	st storage.Interface
	// 最后更新时间,用于跟踪系统参数变更
	lastUpdateTime int64
	areaRepo       fw.Repository[sys.District]
	optRepo        fw.Repository[sys.GeneralOption]
	_stationRepo   sys.IStationRepo
}

func NewSystemRepo(o fw.ORM, st storage.Interface) sys.ISystemRepo {
	_repoOnce.Do(func() {
		_sysRepo = &systemRepoImpl{
			ORM:            o,
			st:             st,
			lastUpdateTime: 0,
		}
		_sysAggregate = impl.NewSystemAggregateRoot(_sysRepo)
	})
	return _sysRepo
}

// FlushUpdateStatus implements sys.ISystemRepo.
func (s *systemRepoImpl) FlushUpdateStatus() {
	// 更新缓存时间
	unix := time.Now().Unix()
	err := s.st.Set("go2o:sys:last_update_time", unix)
	if err != nil {
		logger.Error("set last update time error: %s", err.Error())
	}
	// 清除缓存
	_sysAggregate = impl.NewSystemAggregateRoot(_sysRepo)
}

// LastUpdateTime implements sys.ISystemRepo.
func (s *systemRepoImpl) LastUpdateTime() int64 {
	if s.lastUpdateTime <= 0 {
		key := "go2o:sys:last_update_time"
		updateTime, err := s.st.GetInt64(key)
		if err != nil {
			if s.st.Exists(key) {
				logger.Error("get last update time error: %s", err.Error())
			}

			s.lastUpdateTime = time.Now().Unix()
			s.st.Set(key, s.lastUpdateTime)

		}
		s.lastUpdateTime = updateTime
	}
	return s.lastUpdateTime
}

// Option implements sys.ISystemRepo.
func (s *systemRepoImpl) Option() fw.Repository[sys.GeneralOption] {
	if s.optRepo == nil {
		s.optRepo = newGeneralOptionRepoImpl(s.ORM)
	}
	return s.optRepo
}

// Station implements sys.ISystemRepo.
func (s *systemRepoImpl) Station() sys.IStationRepo {
	if s._stationRepo == nil {
		s._stationRepo = NewStationRepo(s.ORM, s)
	}
	return s._stationRepo
}

// AreaRepo implements sys.ISystemRepo.
func (s *systemRepoImpl) District() fw.Repository[sys.District] {
	if s.areaRepo == nil {
		s.areaRepo = newDistrictRepository(s.ORM)
	}
	return s.areaRepo
}

// GetSystemAggregateRoot implements sys.ISystemRepo.
func (s *systemRepoImpl) GetSystemAggregateRoot() sys.ISystemAggregateRoot {
	return _sysAggregate
}

type districtRepository struct {
	fw.BaseRepository[sys.District]
}

func newDistrictRepository(o fw.ORM) fw.Repository[sys.District] {
	s := &districtRepository{}
	s.ORM = o
	return s
}

type sysGeneralOptionRepoImpl struct {
	fw.BaseRepository[sys.GeneralOption]
}

// NewSysGeneralOptionRepo Create new SysGeneralOptionRepo
func newGeneralOptionRepoImpl(o fw.ORM) fw.Repository[sys.GeneralOption] {
	r := &sysGeneralOptionRepoImpl{}
	r.ORM = o
	return r
}
