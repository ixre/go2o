package fw

import (
	"gorm.io/gorm"
)

// Repository 仓储接口
type (
	// ORM 数据库关系映
	ORM = *gorm.DB

	// 日志
	ILogger interface {
		// 输出调试日志
		Debug(format string, args ...interface{})
		// 输出普通日志
		Info(format string, args ...interface{})
		// 输出错误日志
		Error(format string, args ...interface{})
		// 输出警告日志
		Warn(format string, args ...interface{})
		// 输出致命日志
		Fatal(format string, args ...interface{})
	}

	// Repository 仓储接口
	Repository[M any] interface {
		// Get 获取实体
		Get(id interface{}) *M
		// FindBy 根据条件获取实体
		FindBy(where string, v ...interface{}) *M
		// FindList 查找列表
		FindList(where string, v ...interface{}) []*M
		// Save 保存
		Save(v *M) (*M, error)
		// Delete 删除
		Delete(v *M) error
		// DeleteBy 根据条件删除
		DeleteBy(where string, v ...interface{}) (int, error)
		// PagingQuery 查询分页数据
		PagingQuery(begin, end int, orderBy string, where string, args ...interface{}) (total int, rows []*M, err error)
		// Count 统计条数
		//Count(where string, v ...interface{}) (int, error)
	}
	Service[M any] interface {
		// Get 获取实体
		Get(id interface{}) *M
		// Save 保存
		Save(v *M) (*M, error)
		// FindList 查找列表
		FindList(where string, args ...interface{}) []*M

		// Delete 删除
		Delete(v *M) error
		// PagingQuery 查询分页数据
		PagingQuery(begin, end int, orderBy, where string, args ...interface{}) (total int, rows []*M, err error)
	}
)

var _ Repository[any] = new(BaseRepository[any])

// 基础仓储
type BaseRepository[M any] struct {
	ORM ORM
}

func (r *BaseRepository[M]) Get(id interface{}) *M {
	var e M
	r.ORM.First(&e, id)
	return &e
}

func (r *BaseRepository[M]) joinQueryParams(where string, v ...interface{}) []interface{} {
	var params []interface{}
	if len(where) != 0 {
		params = append([]interface{}{where}, v...)
	} else {
		params = v
	}
	return params
}

func (r *BaseRepository[M]) FindBy(where string, v ...interface{}) *M {
	var m []M
	r.ORM.Limit(1).Find(&m, r.joinQueryParams(where, v...)...)
	if len(m) > 0 {
		return &m[0]
	}
	return nil
}

func (r *BaseRepository[M]) FindList(where string, v ...interface{}) []*M {
	list := make([]*M, 0)
	r.ORM.Select(&list, r.joinQueryParams(where, v...)...)
	return list
}

// SaveStaffExtent Save 商户坐席(员工)扩展表
func (r *BaseRepository[M]) Save(v *M) (*M, error) {
	ctx := r.ORM.Save(v)
	return v, ctx.Error
}

// DeleteStaffExtent Delete 商户坐席(员工)扩展表
func (r *BaseRepository[M]) Delete(v *M) error {
	tx := r.ORM.Delete(v)
	return tx.Error
}

// BatchDeleteStaffExtent Batch Delete 商户坐席(员工)扩展表
func (r *BaseRepository[M]) DeleteBy(where string, v ...interface{}) (int, error) {
	var m M
	tx := r.ORM.Delete(&m, r.joinQueryParams(where, v...)...)
	return int(tx.RowsAffected), nil
}

// PagingQueryStaffExtent Query paging data
func (r *BaseRepository[M]) PagingQuery(begin, end int, orderBy string, where string, args ...interface{}) (total int, rows []*M, err error) {
	var m M
	var list []*M
	var t int64
	wh := func(tx *gorm.DB) *gorm.DB {
		if len(where) > 0 {
			tx.Where(where, args...)
		}
		return tx
	}
	err = wh(r.ORM.Model(&m)).Count(&t).Error
	if err == nil {
		if t > 0 {
			err = wh(r.ORM.Limit(end - begin).Offset(begin)).Find(&list).Error
		}
	}
	return int(t), list, err
}

var _ Service[any] = new(BaseService[any])

type BaseService[M any] struct {
	Repo Repository[M]
}

func (m *BaseService[M]) Save(r *M) (*M, error) {
	return m.Repo.Save(r)
}

func (m *BaseService[M]) Get(id interface{}) *M {
	return m.Repo.Get(id)
}

func (m *BaseService[M]) FindList(where string, args ...interface{}) []*M {
	return m.Repo.FindList(where, args...)
}

func (m *BaseService[M]) Delete(v *M) error {
	return m.Repo.Delete(v)
}

func (m *BaseService[M]) PagingQuery(begin, end int, orderBy, where string, args ...interface{}) (total int, rows []*M, err error) {
	return m.Repo.PagingQuery(begin, end, orderBy, where, args...)
}
