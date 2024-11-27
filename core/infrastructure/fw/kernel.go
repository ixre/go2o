package fw

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/typeconv"
	"gorm.io/gorm"
)

type (
	// QueryOption 列表查询参数
	QueryOption struct {
		// 跳过条数
		Skip int
		// 限制条数
		Limit int
		// 排序
		Order interface{}
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
		// Raw 获取原始ORM
		Raw() ORM
		// Get 获取实体
		Get(id interface{}) *M
		// FindBy 根据条件获取实体
		FindBy(where string, args ...interface{}) *M
		// FindList 查找列表
		FindList(opt *QueryOption, where string, args ...interface{}) []*M
		// Save 保存实体,如主键为空则新增
		Save(v *M) (*M, error)
		// Update 更新实体的非零字段
		Update(v *M) (*M, error)
		// 统计数量
		Count(where string, v ...interface{}) (int, error)
		// Delete 删除
		Delete(v *M) error
		// DeleteBy 根据条件删除
		DeleteBy(where string, v ...interface{}) (int, error)
		// QueryPaging 查询分页数据
		QueryPaging(p *PagingParams) (r *PagingResult, err error)
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
		FindList(opt *QueryOption, where string, args ...interface{}) []*M
		// 统计数量
		Count(where string, v ...interface{}) (int, error)
		// Delete 删除
		Delete(v *M) error
		// QueryPaging 查询分页数据
		QueryPaging(p *PagingParams) (r *PagingResult, err error)
	}
)

var _ Repository[any] = new(BaseRepository[any])

// NewRepository 创建仓储
func NewRepository[M any](orm ORM) Repository[M] {
	return &BaseRepository[M]{ORM: orm}
}

// 基础仓储
type BaseRepository[M any] struct {
	ORM
}

func (r *BaseRepository[M]) Raw() ORM {
	return r.ORM
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

func (r *BaseRepository[M]) FindList(opt *QueryOption, where string, v ...interface{}) []*M {
	list := make([]*M, 0)
	tx := r.ORM
	if opt != nil {
		if opt.Limit > 0 {
			tx = tx.Limit(opt.Limit).Offset(opt.Skip)
		}
		if opt.Order != nil {
			if v, ok := opt.Order.(string); ok {
				// "id desc"
				if len(v) > 0 {
					tx = tx.Order(v)
				}
			} else {
				tx = tx.Order(opt.Order)
			}
		}
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

func (r *BaseRepository[M]) Count(where string, v ...interface{}) (int, error) {
	var count int64
	var m M
	tx := r.ORM.Model(&m).Where(where, v...).Count(&count)
	return int(count), tx.Error
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

func (r *BaseRepository[M]) QueryPaging(p *PagingParams) (ret *PagingResult, err error) {
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

func (m *BaseService[M]) FindList(opt *QueryOption, where string, args ...interface{}) []*M {
	return m.Repo.FindList(opt, where, args...)
}

func (m *BaseService[M]) Count(where string, args ...interface{}) (int, error) {
	return m.Repo.Count(where, args...)
}

func (m *BaseService[M]) Delete(v *M) error {
	return m.Repo.Delete(v)
}

func (m *BaseService[M]) QueryPaging(p *PagingParams) (ret *PagingResult, err error) {
	return m.Repo.QueryPaging(p)
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
func (p *PagingParams) where(field string, exp string, value ...interface{}) *PagingParams {
	buf := bytes.NewBuffer(nil)
	isBlank := len(p.Arguments) == 0
	if !isBlank {
		buf.WriteString(p.Arguments[0].(string))
		buf.WriteString(" AND ")
	}
	buf.WriteString(fmt.Sprintf("%s %s", field, exp))
	if isBlank {
		p.Arguments = []interface{}{buf.String()}
	} else {
		p.Arguments[0] = buf.String()
	}
	p.Arguments = append(p.Arguments, value...)
	return p
}

// Equal 方法为 PagingParams 结构体添加了一个 Equal 功能
// 用于在查询条件中添加等于某个字段值的条件
//
// 参数：
// field string - 要进行比较的字段名
// value interface{} - 要进行比较的值，可以是任意类型
//
// 返回值：
// *PagingParams - 返回一个新的 PagingParams 指针，用于链式调用
func (p *PagingParams) Equal(field string, value interface{}) *PagingParams {
	return p.where(field, "= ?", value)
}

// NotEqual 方法用于在 PagingParams 结构体上添加一个不等于条件的查询参数
//
// 参数：
//
//	p *PagingParams - PagingParams 结构体指针，表示当前分页参数对象
//	field string - 要进行不等于条件判断的字段名
//	value interface{} - 要进行不等于条件判断的值
//
// 返回值：
//
//	*PagingParams - 返回修改后的 PagingParams 结构体指针
func (p *PagingParams) NotEqual(field string, value interface{}) *PagingParams {
	return p.where(field, " <> ?", value)
}

// In 函数向PagingParams结构体中添加一个IN查询条件
//
// 参数：
//
//	p *PagingParams - PagingParams结构体指针，用于存储查询条件
//	field string - 要查询的字段名
//	value ...interface{} - 要查询的值列表，可以为单个值或切片
//
// 返回值：
//
//	*PagingParams - 添加了IN查询条件的PagingParams结构体指针
func (p *PagingParams) In(field string, value ...interface{}) *PagingParams {
	l := len(value)
	if l == 0 {
		panic("value is empty")
	}
	if reflect.TypeOf(value[0]).Kind() == reflect.Slice {
		if l > 1 {
			panic("value is slice,but give more than one value")
		}
		return p.where(field, " IN ?", value[0])
	}
	return p.where(field, " IN ?", value)
}
func (p *PagingParams) Like(field string, value interface{}) *PagingParams {
	return p.where(field, "LIKE ?", value)
}

// Gt 大于
func (p *PagingParams) Gt(field string, value interface{}) *PagingParams {
	return p.where(field, "> ?", value)
}

// Lt 小于
func (p *PagingParams) Lt(field string, value interface{}) *PagingParams {
	return p.where(field, "< ?", value)
}

func (p *PagingParams) And(where string, values ...interface{}) *PagingParams {
	buf := bytes.NewBuffer(nil)
	isBlank := len(p.Arguments) == 0
	if !isBlank {
		buf.WriteString(p.Arguments[0].(string))
		buf.WriteString(" AND ")
	}
	buf.WriteString(where)
	if isBlank {
		p.Arguments = []interface{}{buf.String()}
	} else {
		p.Arguments[0] = buf.String()
	}
	p.Arguments = append(p.Arguments, values...)
	return p
}

func (p *PagingParams) Between(field string, arr []string) (*PagingParams, error) {
	if len(arr) != 2 {
		return nil, errors.New("between need two value")
	}
	return p.where(field, "BETWEEN ? AND ?", arr[0], arr[1]), nil
}

func (p *PagingParams) BetweenInts(field string, arr []int) (*PagingParams, error) {
	if len(arr) != 2 {
		return nil, errors.New("between need two value")
	}
	return p.where(field, "BETWEEN ? AND ?", arr[0], arr[1]), nil
}

// OrderBy 添加排序条件
func (p *PagingParams) OrderBy(order string) *PagingParams {
	p.Order = order
	return p
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
func ReduceFinds[T any](fn func(opt *QueryOption) []*T, size int) (arr []*T) {
	begin := 0
	for {
		list := fn(&QueryOption{
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

// UnifinedQueryPaging 通用查询
//
//	tables Like: mm_member m ON m.id = o.member_id INNER JOIN mch_merchant mch ON mch.id = o.mch_id
//	fields like: s.gender,m.nickname,certified_name
func UnifinedQueryPaging(o ORM, p *PagingParams, tables string, fields string) (_ *PagingResult, err error) {
	var ret PagingResult
	from := `FROM ` + tables
	// 查询条件
	where, args := "", []interface{}{}
	if len(p.Arguments) > 0 {
		where = " WHERE " + p.Arguments[0].(string)
		args = p.Arguments[1:]
	}
	// 查询条数
	sql := fmt.Sprintf("SELECT COUNT(*) %s %s", from, where)
	o.Raw(sql, args...).Scan(&ret.Total)
	// 查询行数
	order := types.Ternary(p.Order != "", " ORDER BY "+p.Order, "")
	if ret.Total > 0 {
		skipper := GetSkipperSQL(o, p)
		sql = strings.Join([]string{"SELECT", fields, from, where, order, skipper}, " ")
		rows, err := o.Raw(sql, args...).Rows()
		if err != nil {
			log.Println("paging query rows error: %s", err.Error())
		} else {
			for _, v := range db.RowsToMarshalMap(rows) {
				ret.Rows = append(ret.Rows, v)
			}
			rows.Close()
		}
	}
	return &ret, nil
}

// 生成分页条件
func GetSkipperSQL(o ORM, p *PagingParams) string {
	if p.Size <= 0 {
		return ""
	}
	switch o.Dialector.Name() {
	case "mysql":
		return fmt.Sprintf(" LIMIT %d,%d", p.Begin, p.Size)
	case "postgres":
		return fmt.Sprintf(" OFFSET %d LIMIT %d", p.Begin, p.Size)
	case "sqlite3":
		return fmt.Sprintf(" LIMIT %d OFFSET %d", p.Size, p.Begin)
	case "mssql":
		return fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", p.Begin, p.Size)
	}
	panic("not support dialect")
}

// 分页行
type EffectRow struct {
	v map[string]interface{}
}

func ParseRow(v interface{}) *EffectRow {
	return &EffectRow{v: v.(map[string]interface{})}
}

func (p *EffectRow) AsInt(keys ...string) {
	for _, key := range keys {
		v, ok := p.v[key].([]uint8)
		if ok {
			f := typeconv.MustInt(string(v))
			p.v[key] = f
		}
	}
}

// 转换为float类型
func (p *EffectRow) AsFloat(keys ...string) {
	for _, key := range keys {
		v, ok := p.v[key].([]uint8)
		if ok {
			f := typeconv.MustFloat(string(v))
			p.v[key] = f
		}
	}
}

// Excludes 排除字段
func (p *EffectRow) Excludes(keys ...string) {
	for _, key := range keys {
		delete(p.v, key)
	}
}

// Put 添加/更新字段
func (p *EffectRow) Put(key string, v interface{}) {
	p.v[key] = v
}

// Get 获取字段
func (p *EffectRow) Get(key string) interface{} {
	return p.v[key]
}

// GetInt 获取int类型字段
func (p *EffectRow) GetInt(key string) int {
	return typeconv.Int(p.Get(key))
}

/** 错误处理 */
type Error struct {
	// 错误码
	Code int `json:"code"`
	// 错误信息
	Message string `json:"message"`
}

// 解析错误
func ParseError(err error) *Error {
	if err != nil {
		return &Error{
			Code:    1,
			Message: err.Error(),
		}
	}
	return nil
}

// NewError 创建错误
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// 断言错误
func AssertError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
