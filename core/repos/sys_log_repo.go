package repos

import (
	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var _ sys.ILogRepository = new(SysLogRepoImpl)

type SysLogRepoImpl struct {
	_logRepo fw.Repository[sys.SysLog]
	_appRepo fw.Repository[sys.LogApp]
}

func NewSysLogRepo(db fw.ORM) sys.ILogRepository {
	return &SysLogRepoImpl{
		_logRepo: &fw.BaseRepository[sys.SysLog]{
			ORM: db,
		},
		_appRepo: &fw.BaseRepository[sys.LogApp]{
			ORM: db,
		},
	}
}

func (s *SysLogRepoImpl) App() fw.Repository[sys.LogApp] {
	return s._appRepo
}
func (s *SysLogRepoImpl) Log() fw.Repository[sys.SysLog] {
	return s._logRepo
}
