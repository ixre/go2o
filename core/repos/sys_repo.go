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
	areaRepo       fw.Repository[sys.Region]
	optRepo        fw.Repository[sys.GeneralOption]
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
		logger.Error("sys", "set last update time error: %s", err.Error())
	}
	// 清除缓存
	_sysAggregate = impl.NewSystemAggregateRoot(_sysRepo)
}

// LastUpdateTime implements sys.ISystemRepo.
func (s *systemRepoImpl) LastUpdateTime() int64 {
	if s.lastUpdateTime <= 0 {
		updateTime, err := s.st.GetInt64("go2o:sys:last_update_time")
		if err != nil {
			logger.Error("sys", "get last update time error: %s", err.Error())
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

// AreaRepo implements sys.ISystemRepo.
func (s *systemRepoImpl) Region() fw.Repository[sys.Region] {
	if s.areaRepo == nil {
		s.areaRepo = newAreaRepository(s.ORM)
	}
	return s.areaRepo
}

// GetSystemAggregateRoot implements sys.ISystemRepo.
func (s *systemRepoImpl) GetSystemAggregateRoot() sys.ISystemAggregateRoot {
	return _sysAggregate
}

type areaRepository struct {
	fw.BaseRepository[sys.Region]
}

func newAreaRepository(o fw.ORM) fw.Repository[sys.Region] {
	s := &areaRepository{}
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
