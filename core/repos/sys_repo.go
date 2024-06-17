package repos

import (
	"sync"

	"github.com/ixre/go2o/core/domain/interface/sys"
	impl "github.com/ixre/go2o/core/domain/sys"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

var (
	_sysRepo      = new(systemRepoImpl)
	_sysAggregate sys.ISystemAggregateRoot
	_repoOnce     sync.Once
)

type systemRepoImpl struct {
	fw.ORM
	areaRepo fw.Repository[sys.Region]
	optRepo  fw.Repository[sys.GeneralOption]
}

func NewSystemRepo(o fw.ORM) sys.ISystemRepo {
	_repoOnce.Do(func() {
		_sysRepo = &systemRepoImpl{
			ORM: o,
		}
		_sysAggregate = impl.NewSystemAggregateRoot(_sysRepo)
	})
	return _sysRepo
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

// GetAllCities implements sys.ISystemRepo.
func (s *systemRepoImpl) GetAllCities() []*sys.Region {
	panic("unimplemented")
}

// GetRegionList implements sys.ISystemRepo.
func (s *systemRepoImpl) GetRegionList(parentId int) []*sys.Region {
	panic("unimplemented")
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
