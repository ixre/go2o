package repos

import (
	"github.com/ixre/go2o/core/domain/interface/work/workorder"
	"github.com/ixre/go2o/core/infrastructure/fw"
	workorderImpl "github.com/ixre/go2o/core/domain/work/workorder"
)

var _ workorder.IWorkorderRepo = new(workorderRepo)

type workorderRepo struct {
	fw.BaseRepository[workorder.Workorder]
	commentRepo fw.Repository[workorder.WorkorderComment]
}

// NewWorkorderRepo 函数用于创建一个新的工作单仓库实例
//
// 参数：
//
//	db fw.ORM - 用于数据库操作的ORM对象
//
// 返回值：
//
//	workorder.IWorkorderRepo - 工作单仓库接口实例
func NewWorkorderRepo(db fw.ORM) workorder.IWorkorderRepo {
	r := &workorderRepo{}
	r.ORM = db
	r.commentRepo = &fw.BaseRepository[workorder.WorkorderComment]{
		ORM: db,
	}
	return r
}

// CommentRepo implements workorder.IWorkorderRepo.
func (w *workorderRepo) CommentRepo() fw.Repository[workorder.WorkorderComment] {
	return w.commentRepo
}

// CreateWorkorder implements workorder.IWorkorderRepo.
func (w *workorderRepo) CreateWorkorder(value *workorder.Workorder) workorder.IWorkorderAggregateRoot {
	return workorderImpl.NewWorkorder(value, w)
}

// GetWorkorder implements workorder.IWorkorderRepo.
func (w *workorderRepo) GetWorkorder(id int) workorder.IWorkorderAggregateRoot {
	v := w.Get(id)
	if v == nil {
		return nil
	}
	return w.CreateWorkorder(v)
}
