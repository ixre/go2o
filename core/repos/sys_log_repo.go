package repos

import (
	"sync"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var _ sys.IApplicationRepository = new(SysAppRepoImpl)
var _singleton sys.IApplicationRepository
var _once sync.Once

type SysAppRepoImpl struct {
	_logRepo          fw.Repository[sys.SysLog]
	_versionRepo      fw.Repository[sys.SysAppVersion]
	_distributionRepo fw.Repository[sys.SysAppDistribution]
}

func NewSysAppRepo(db fw.ORM) sys.IApplicationRepository {
	_once.Do(func() {
		_singleton = &SysAppRepoImpl{
			_logRepo:          fw.NewRepository[sys.SysLog](db),
			_versionRepo:      fw.NewRepository[sys.SysAppVersion](db),
			_distributionRepo: fw.NewRepository[sys.SysAppDistribution](db),
		}
	})
	return _singleton
}

func (s *SysAppRepoImpl) Log() fw.Repository[sys.SysLog] {
	return s._logRepo
}
func (s *SysAppRepoImpl) Version() fw.Repository[sys.SysAppVersion] {
	return s._versionRepo
}

func (s *SysAppRepoImpl) Distribution() fw.Repository[sys.SysAppDistribution] {
	return s._distributionRepo
}
