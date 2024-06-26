package fw

import (
	"bytes"
	"fmt"

	"gorm.io/gorm"
)

type (
	// ListOption 列表查询参数
	ListOption struct {
		// 跳过条数
		Skip int
		// 限制条数
		Limit int
	}
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
		FindBy(where string, args ...interface{}) *M
		// FindList 查找列表
		FindList(opt *ListOption, where string, args ...interface{}) []*M
		// Save 保存实体,如主键为空则新增
		Save(v *M) (*M, error)
		// Update 更新实体的非零字段
		Update(v *M) (*M, error)
		// Delete 删除
		Delete(v *M) error
		// DeleteBy 根据条件删除
		DeleteBy(where string, v ...interface{}) (int, error)
		// PagingQuery 查询分页数据
		PagingQuery(p *PagingParams) (r *PagingResult, err error)
		// Count 统计条数
		//Count(where string, v ...interface{}) (int, error)
	}
	Service[M any] interface {
		// Get 获取实体
		Get(id interface{}) *M
		// FindBy 根据条件获取实体
		FindBy(where string, v ...interface{}) *M
		// Save 保存
		Save(v *M) (*M, error)
		// FindList 查找列表
		FindList(opt *ListOption, where string, args ...interface{}) []*M

		// Delete 删除
		Delete(v *M) error
		// PagingQuery 查询分页数据
		PagingQuery(p *PagingParams) (r *PagingResult, err error)
	}
)

var _ Repository[any] = new(BaseRepository[any])

// 基础仓储
type BaseRepository[M any] struct {
	ORM
}

func (r *BaseRepository[M]) Get(id interface{}) *M {
	var e M
	r.ORM.First(&e, id)
	return &e
}

func (r *BaseRepository[M]) joinQueryParams(where string, v ...interface{}) []interface{} {
	var params []interface{}
	if len(where) == 0 {
		return params
	}
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

func (r *BaseRepository[M]) FindList(opt *ListOption, where string, v ...interface{}) []*M {
	list := make([]*M, 0)
	var tx *gorm.DB
	if opt != nil && opt.Limit > 0 {
		tx = r.ORM.Limit(opt.Limit).Offset(opt.Skip)
	} else {
		tx = r.ORM
	}
	tx.Find(&list, r.joinQueryParams(where, v...)...)
	return list
}

// Save 保存实体
func (r *BaseRepository[M]) Save(v *M) (*M, error) {
	ctx := r.ORM.Save(v)
	return v, ctx.Error
}

// Update 更新实体的非零字段
func (r *BaseRepository[M]) Update(v *M) (*M, error) {
	ctx := r.ORM.Model(v).Updates(v)
	return v, ctx.Error
}

func (r *BaseRepository[M]) Delete(v *M) error {
	tx := r.ORM.Delete(v)
	return tx.Error
}

func (r *BaseRepository[M]) DeleteBy(where string, v ...interface{}) (int, error) {
	var m M
	tx := r.ORM.Delete(&m, r.joinQueryParams(where, v...)...)
	return int(tx.RowsAffected), nil
}

func (r *BaseRepository[M]) PagingQuery(p *PagingParams) (ret *PagingResult, err error) {
	var m M
	var t int64
	wh := func(tx *gorm.DB) *gorm.DB {
		if len(p.Arguments) > 0 {
			tx.Where(p.Arguments[0], p.Arguments[1:]...)
		}
		return tx
	}
	var list []*M
	err = wh(r.ORM.Model(&m)).Count(&t).Error
	if err == nil && t > 0 {
		tx := r.ORM.Limit(p.Size).Offset(p.Begin)
		if len(p.Order) > 0 {
			// 排序
			tx = tx.Order(p.Order)
		}
		err = wh(tx).Find(&list).Error
	}
	var arr = make([]interface{}, 0)
	for _, v := range list {
		arr = append(arr, v)
	}
	return &PagingResult{
		Total: int(t),
		Rows:  arr,
		Extra: nil,
	}, err
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

func (m *BaseService[M]) FindBy(where string, args ...interface{}) *M {
	return m.Repo.FindBy(where, args...)
}

func (m *BaseService[M]) FindList(opt *ListOption, where string, args ...interface{}) []*M {
	return m.Repo.FindList(opt, where, args...)
}

func (m *BaseService[M]) Delete(v *M) error {
	return m.Repo.Delete(v)
}

func (m *BaseService[M]) PagingQuery(p *PagingParams) (ret *PagingResult, err error) {
	return m.Repo.PagingQuery(p)
}

// 分页参数
// // Arguments: 第一个参数为SQL条件语句,后面跟参数,如:
//
//	&PagingParams{
//		Arguments: []interface{}{"id=? and name=?", "test"},
//	}
type PagingParams struct {
	// 开始数量
	Begin int `json:"begin"`
	// 单页数量
	Size int `json:"size"`
	// 排序条件
	Order string `json:"order"`
	// SQL条件及参数
	Arguments []interface{} `json:"arguments"`
}

// Where 添加条件
func (p *PagingParams) where(field string, exp string, value interface{}) *PagingParams {
	buf := bytes.NewBuffer(nil)
	isBlank := len(p.Arguments) == 0
	if !isBlank {
		buf.WriteString(p.Arguments[0].(string))
		buf.WriteString(" AND ")
	}
	buf.WriteString(fmt.Sprintf("%s %s", field, exp))
	if isBlank {
		p.Arguments = []interface{}{buf.String(), value}
	} else {
		p.Arguments[0] = buf.String()
		p.Arguments = append(p.Arguments, value)
	}
	return p
}

func (p *PagingParams) Equal(field string, value interface{}) *PagingParams {
	return p.where(field, "=?", value)
}

func (p *PagingParams) Like(field string, value interface{}) *PagingParams {
	return p.where(field, "LIKE ?", value)
}

// 分页结果
type PagingResult struct {
	// 总数量
	Total int `json:"total"`
	// 当前页数据
	Rows []interface{} `json:"rows"`
	// 额外信息
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// ReduceFinds 分次查询合并数组,用于分次查询出数量较多的数据
func ReduceFinds[T any](fn func(opt *ListOption) []*T, size int) (arr []*T) {
	begin := 0
	for {
		list := fn(&ListOption{
			Skip:  begin,
			Limit: size,
		})
		l := len(list)
		if l != 0 {
			arr = append(arr, list...)
		}
		begin += l
		if l < size {
			break
		}
	}
	return arr
}
